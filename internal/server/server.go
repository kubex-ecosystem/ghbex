// Package server provides an HTTP server for the application.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
	githubx "github.com/rafa-mori/ghbex/internal/client"
	config "github.com/rafa-mori/ghbex/internal/config"
	"github.com/rafa-mori/ghbex/internal/defs"
	"github.com/rafa-mori/ghbex/internal/manager"
	notify "github.com/rafa-mori/ghbex/internal/notifiers"
	"github.com/rafa-mori/ghbex/internal/operators/analytics"
)

type GHServerEngine interface {
	Start(context.Context) error
	Stop(context.Context) error
	Status(context.Context) error
}

type ghServerEngine struct {
	MainConfig config.MainConfig
	// ghc Is the GitHub client used for interacting with GitHub APIs.
	ghc *github.Client
}

func NewGHServerEngine(cfg config.MainConfig) GHServerEngine {
	return &ghServerEngine{
		MainConfig: cfg,
		ghc:        nil,
	}
}

func NewGithubClient(ctx context.Context, cfg config.MainConfig) (*github.Client, error) {
	var err error
	var ghc *github.Client
	gh := cfg.GetGitHub()
	expandedToken := os.ExpandEnv(gh.Auth.Token)
	log.Printf("DEBUG: Auth kind: %s, Token template: %s, Expanded token length: %d",
		gh.Auth.Kind, gh.Auth.Token, len(expandedToken))
	switch strings.ToLower(gh.Auth.Kind) {
	case "pat":
		ghc, err = githubx.NewPAT(ctx, githubx.PATConfig{
			Token:     expandedToken,
			BaseURL:   gh.Auth.BaseURL,
			UploadURL: gh.Auth.UploadURL,
		})
	case "app":
		ghc, err = githubx.NewApp(ctx, githubx.AppConfig{
			AppID:          gh.Auth.AppID,
			InstallationID: gh.Auth.InstallationID,
			PrivateKeyPath: gh.Auth.PrivateKeyPath,
			BaseURL:        gh.Auth.BaseURL,
			UploadURL:      gh.Auth.UploadURL,
		})
	default:
		err = errors.New("github.auth.kind must be pat|app")
	}

	if err != nil {
		return nil, err
	}

	return ghc, nil
}

func (g *ghServerEngine) Start(ctx context.Context) error {
	var err error

	// build github client
	g.ghc, err = NewGithubClient(ctx, g.MainConfig)
	if err != nil {
		log.Fatal(err)
	}

	// notifiers
	var notifierz []defs.Notifier
	for _, n := range *g.MainConfig.GetNotifiers() {
		switch n.Type {
		case "discord":
			notifierz = append(notifierz, notify.Discord{Webhook: os.ExpandEnv(n.Webhook)})
		case "stdout":
			notifierz = append(notifierz, &notify.Stdout{})
		}
	}

	// service
	svc := manager.New(
		g.ghc,
		g.MainConfig,
	)

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		response := map[string]interface{}{
			"status":       "ok",
			"version":      "0.0.1",
			"github_auth":  g.ghc != nil,
			"config_repos": len(g.MainConfig.GetGitHub().Repos),
		}
		_ = json.NewEncoder(w).Encode(response)
	})

	// List configured repositories endpoint
	http.HandleFunc("/repos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		repos := make([]map[string]interface{}, 0)
		for _, repo := range g.MainConfig.GetGitHub().Repos {
			repoInfo := map[string]interface{}{
				"owner": repo.Owner,
				"name":  repo.Name,
				"url":   "https://github.com/" + repo.Owner + "/" + repo.Name,
				"rules": map[string]interface{}{
					"runs": map[string]interface{}{
						"max_age_days":      repo.Rules.Runs.MaxAgeDays,
						"keep_success_last": repo.Rules.Runs.KeepSuccessLast,
					},
					"artifacts": map[string]interface{}{
						"max_age_days": repo.Rules.Artifacts.MaxAgeDays,
					},
					"monitoring": map[string]interface{}{
						"inactive_days_threshold": repo.Rules.Monitoring.InactiveDaysThreshold,
					},
				},
			}
			repos = append(repos, repoInfo)
		}

		response := map[string]interface{}{
			"total":        len(repos),
			"repositories": repos,
		}
		_ = json.NewEncoder(w).Encode(response)
	})

	// Bulk sanitize endpoint for multiple repositories
	http.HandleFunc("/admin/sanitize/bulk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST", http.StatusMethodNotAllowed)
			return
		}

		dry := r.URL.Query().Get("dry_run")
		dryRun := dry == "1" || strings.EqualFold(dry, "true")

		var bulkResults []map[string]interface{}
		totalRuns := 0
		totalArtifacts := 0
		startTime := time.Now()

		log.Printf("🚀 BULK SANITIZATION STARTED - DRY_RUN: %v", dryRun)

		for _, repoConfig := range g.MainConfig.GetGitHub().Repos {
			log.Printf("📊 Processing %s/%s...", repoConfig.Owner, repoConfig.Name)

			rpt, err := svc.SanitizeRepo(r.Context(), repoConfig.Owner, repoConfig.Name, repoConfig.Rules, dryRun)
			if err != nil {
				log.Printf("❌ Error processing %s/%s: %v", repoConfig.Owner, repoConfig.Name, err)
				continue
			}

			totalRuns += rpt.Runs.Deleted
			totalArtifacts += rpt.Artifacts.Deleted

			result := map[string]interface{}{
				"owner":     rpt.Owner,
				"repo":      rpt.Repo,
				"runs":      rpt.Runs.Deleted,
				"artifacts": rpt.Artifacts.Deleted,
				"releases":  rpt.Releases.DeletedDrafts,
				"success":   true,
			}
			bulkResults = append(bulkResults, result)

			log.Printf("✅ %s/%s - Runs: %d, Artifacts: %d", repoConfig.Owner, repoConfig.Name, rpt.Runs.Deleted, rpt.Artifacts.Deleted)
		}

		duration := time.Since(startTime)

		response := map[string]interface{}{
			"bulk_operation":          true,
			"dry_run":                 dryRun,
			"started_at":              startTime.Format("2006-01-02 15:04:05"),
			"duration_ms":             duration.Milliseconds(),
			"total_repos":             len(bulkResults),
			"total_runs_cleaned":      totalRuns,
			"total_artifacts_cleaned": totalArtifacts,
			"productivity_summary": map[string]interface{}{
				"estimated_storage_saved_mb": (totalRuns * 10) + (totalArtifacts * 50), // Estimativa
				"estimated_time_saved_min":   (totalRuns + totalArtifacts) * 2,         // Estimativa
			},
			"repositories": bulkResults,
		}

		log.Printf("🎉 BULK SANITIZATION COMPLETED - Duration: %v, Total Runs: %d, Total Artifacts: %d",
			duration, totalRuns, totalArtifacts)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(response)
	})

	// Analytics endpoint for repository insights
	http.HandleFunc("/analytics/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}

		// Parse path: /analytics/{owner}/{repo}
		path := strings.TrimPrefix(r.URL.Path, "/analytics/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			http.Error(w, "missing owner/repo in path", http.StatusBadRequest)
			return
		}

		owner, repo := parts[0], parts[1]

		// Get analysis days from query param (default 90)
		analysisDays := 90
		if days := r.URL.Query().Get("days"); days != "" {
			if parsed, err := time.ParseDuration(days + "h"); err == nil {
				analysisDays = int(parsed.Hours() / 24)
			}
		}

		log.Printf("🔍 ANALYTICS REQUEST - %s/%s - Analysis Days: %d", owner, repo, analysisDays)
		startTime := time.Now()

		// Perform analytics
		insights, err := analytics.AnalyzeRepository(r.Context(), g.ghc, owner, repo, analysisDays)
		if err != nil {
			log.Printf("❌ Analytics error for %s/%s: %v", owner, repo, err)
			http.Error(w, fmt.Sprintf("Analytics failed: %v", err), http.StatusInternalServerError)
			return
		}

		duration := time.Since(startTime)
		log.Printf("✅ ANALYTICS COMPLETED - %s/%s - Duration: %v, Health Score: %.1f (%s)",
			owner, repo, duration, insights.HealthScore.Overall, insights.HealthScore.Grade)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(insights)
	})

	// route: POST /admin/repos/{owner}/{repo}/sanitize?dry_run=1
	http.HandleFunc("/admin/repos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST", http.StatusMethodNotAllowed)
			return
		}
		// naive parse
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/admin/repos/"), "/")
		if len(parts) < 3 || parts[2] != "sanitize" {
			http.NotFound(w, r)
			return
		}
		owner, repo := parts[0], parts[1]
		dry := r.URL.Query().Get("dry_run")
		dryRun := dry == "1" || strings.EqualFold(dry, "true")

		log.Printf("🎯 INDIVIDUAL SANITIZATION - %s/%s - DRY_RUN: %v", owner, repo, dryRun)
		startTime := time.Now()

		// find rules (optional override via cfg)
		var rules defs.Rules
		for _, rc := range g.MainConfig.GetGitHub().Repos {
			if rc.Owner == owner && rc.Name == repo {
				rules = rc.Rules
				break
			}
		}

		var dummy defs.Rules
		dummy.Runs.MaxAgeDays = 30
		dummy.Artifacts.MaxAgeDays = 7
		dummy.Releases.DeleteDrafts = true

		if rules.Artifacts == dummy.Artifacts &&
			rules.Runs.MaxAgeDays == dummy.Runs.MaxAgeDays &&
			rules.Releases == dummy.Releases {
			// default sane rules
			rules.Runs.MaxAgeDays = 30
			rules.Artifacts.MaxAgeDays = 7
			rules.Releases.DeleteDrafts = true
		}

		rpt, err := svc.SanitizeRepo(r.Context(), owner, repo, rules, dryRun)
		if err != nil {
			log.Printf("❌ Error sanitizing %s/%s: %v", owner, repo, err)
			http.Error(w, err.Error(), 500)
			return
		}

		duration := time.Since(startTime)
		log.Printf("✅ SANITIZATION COMPLETED - %s/%s - Duration: %v, Runs: %d, Artifacts: %d",
			owner, repo, duration, rpt.Runs.Deleted, rpt.Artifacts.Deleted)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(rpt)
	})

	srv := &http.Server{
		Addr:              g.MainConfig.GetServer().Addr,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", g.MainConfig.GetServer().Addr)
	log.Fatal(srv.ListenAndServe())

	return nil
}

func (g *ghServerEngine) Stop(ctx context.Context) error {
	// Implement stop logic
	return nil
}

func (g *ghServerEngine) Status(ctx context.Context) error {
	// Implement status logic
	return nil
}
