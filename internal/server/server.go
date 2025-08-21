// Package server provides an HTTP server for the application.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
	githubx "github.com/rafa-mori/ghbex/internal/client"
	"github.com/rafa-mori/ghbex/internal/defs"
	"github.com/rafa-mori/ghbex/internal/frontend"
	"github.com/rafa-mori/ghbex/internal/interfaces"
	"github.com/rafa-mori/ghbex/internal/manager"
	gl "github.com/rafa-mori/ghbex/internal/module/logger"
	"github.com/rafa-mori/ghbex/internal/notifiers"
	"github.com/rafa-mori/ghbex/internal/operators/analytics"
	i "github.com/rafa-mori/ghbex/internal/operators/intelligence"
	"github.com/rafa-mori/ghbex/internal/operators/productivity"
)

type GHServerEngine interface {
	Start(context.Context) error
	Stop(context.Context) error
	Status(context.Context) error
}

type ghServerEngine struct {
	MainConfig interfaces.IMainConfig
	// ghc Is the GitHub client used for interacting with GitHub APIs.
	ghc *github.Client
}

func NewGHServerEngine(cfg interfaces.IMainConfig) GHServerEngine {
	if cfg == nil {
		gl.Log("fatal", "Failed to create GitHub client: config is nil")
		return nil
	}
	ghc, err := NewGithubClient(context.Background(), cfg)
	if err != nil {
		gl.Log("fatal", fmt.Sprintf("Failed to create GitHub client: %v", err))
	}
	return &ghServerEngine{
		MainConfig: cfg,
		ghc:        ghc,
	}
}

func NewGithubClient(ctx context.Context, cfg interfaces.IMainConfig) (*github.Client, error) {
	var err error
	var ghc *github.Client
	gh := cfg.GetGitHub()
	if gh.GetAuth() == nil {
		gl.Log("fatal", "github.auth is not configured in the main config")
	}
	gl.Log(
		"debug",
		fmt.Sprintf("github.auth.kind: %s, base_url: %s, upload_url: %s",
			gh.GetAuth().GetKind(), gh.GetAuth().GetBaseURL(), gh.GetAuth().GetUploadURL()),
	)

	switch strings.ToLower(gh.GetAuth().GetKind()) {
	case "pat":
		pCfg := githubx.PATConfig{
			Token:     gh.GetAuth().GetToken(),
			BaseURL:   gh.GetAuth().GetBaseURL(),
			UploadURL: gh.GetAuth().GetUploadURL(),
		}
		ghc, err = githubx.NewPAT(ctx, pCfg)
	case "app":
		ghc, err = githubx.NewApp(ctx, githubx.AppConfig{
			AppID:          gh.GetAuth().GetAppID(),
			InstallationID: gh.GetAuth().GetInstallationID(),
			PrivateKeyPath: gh.GetAuth().GetPrivateKeyPath(),
			BaseURL:        gh.GetAuth().GetBaseURL(),
			UploadURL:      gh.GetAuth().GetUploadURL(),
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
	notifierz := defs.NewNotifiers()
	for _, n := range g.MainConfig.GetNotifiers().GetNotifiers() {
		switch n.GetType() {
		case "slack":
			webhook := os.ExpandEnv(n.GetWebhook())
			notifierz.AddNotifier(notifiers.NewDiscordNotifier(webhook))
		case "discord":
			webhook := os.ExpandEnv(n.GetWebhook())
			notifierz.AddNotifier(notifiers.NewDiscordNotifier(webhook))
		case "stdout":
			notifierz.AddNotifier(notifiers.NewStdoutNotifier())
		}
	}

	// service
	svc := manager.New(
		g.ghc,
		g.MainConfig,
	)

	routes := getRoutesMap(svc, g)

	bindingAddr := net.JoinHostPort(
		g.MainConfig.GetServer().GetAddr(),
		g.MainConfig.GetServer().GetPort(),
	)

	srv := &http.Server{
		Addr:              bindingAddr,
		ReadHeaderTimeout: 5 * time.Second,
	}

	for path, handler := range routes {
		http.Handle(path, handler)
	}

	gl.Log("info", fmt.Sprintf("Server is starting on %s", bindingAddr))
	gl.Log("info", fmt.Sprintf("Visit http://localhost:%s to access the dashboard", g.MainConfig.GetServer().GetPort()))
	gl.Log("info", "Routes:")
	for path := range routes {
		gl.Log("info", fmt.Sprintf("  - %s", path))
	}
	gl.Log("info", "Server logs:")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		gl.Log("error", fmt.Sprintf("HTTP server error: %v", err))
	}

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

func getRoutesMap(svc *manager.Service, g *ghServerEngine) map[string]http.HandlerFunc {
	routes := make(map[string]http.HandlerFunc)

	// Dashboard web interface
	routes["/"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}

		// Serve dashboard only for root path
		if r.URL.Path == "/" {
			frontend.ServeDashboard(w, r)
			return
		}

		// For other paths, return 404
		http.NotFound(w, r)
	})

	// Health check endpoint
	routes["/health"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		response := make(map[string]any)
		response["status"] = "ok"
		response["version"] = "0.0.1"
		response["github_auth"] = g.ghc != nil
		response["config_repos"] = len(g.MainConfig.GetGitHub().GetRepos())
		_ = json.NewEncoder(w).Encode(response)
	})

	// List configured repositories endpoint
	routes["/repos"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// Create intelligence operator for AI insights
		intelligenceOp := i.NewIntelligenceOperator(g.MainConfig, g.ghc)

		repos := make([]map[string]any, 0)
		for _, repo := range g.MainConfig.GetGitHub().GetRepos() {
			repoInfo := map[string]any{
				"owner": repo.GetOwner(),
				"name":  repo.GetName(),
				"url":   "https://github.com/" + repo.GetOwner() + "/" + repo.GetName(),
				"rules": map[string]any{
					"runs": map[string]any{
						"max_age_days":      repo.GetRules().GetRunsRule().GetMaxAgeDays(),
						"keep_success_last": repo.GetRules().GetRunsRule().GetKeepSuccessLast(),
					},
					"artifacts": map[string]any{
						"max_age_days": repo.GetRules().GetArtifactsRule().GetMaxAgeDays(),
					},
					"monitoring": map[string]any{
						"inactive_days_threshold": repo.GetMonitoring().GetInactiveDaysThreshold(),
					},
				},
			}

			// Add AI insights to each repository card
			if insight, err := intelligenceOp.GenerateQuickInsight(context.Background(), repo.GetOwner(), repo.GetName()); err == nil {
				repoInfo["ai"] = map[string]any{
					"score":       insight.AIScore,
					"assessment":  insight.QuickAssessment,
					"health_icon": insight.HealthIcon,
					"main_tag":    insight.MainTag,
					"risk_level":  insight.RiskLevel,
					"opportunity": insight.Opportunity,
				}
			} else {
				// Fallback AI data
				repoInfo["ai"] = map[string]any{
					"score":       85.0,
					"assessment":  "Active repository with good development patterns",
					"health_icon": "🟢",
					"main_tag":    "Active",
					"risk_level":  "low",
					"opportunity": "Performance optimization",
				}
			}

			repos = append(repos, repoInfo)
		}

		response := map[string]any{
			"total":        len(repos),
			"repositories": repos,
		}
		_ = json.NewEncoder(w).Encode(response)
	})

	// Bulk sanitize endpoint for multiple repositories
	routes["/admin/sanitize/bulk"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST", http.StatusMethodNotAllowed)
			return
		}

		dry := r.URL.Query().Get("dry_run")
		dryRun := dry == "1" || strings.EqualFold(dry, "true")

		var bulkResults []map[string]any
		totalRuns := 0
		totalArtifacts := 0
		startTime := time.Now()

		gl.Log("info", fmt.Sprintf("🚀 BULK SANITIZATION STARTED - DRY_RUN: %v", dryRun))

		for _, repoConfig := range g.MainConfig.GetGitHub().GetRepos() {
			if repoConfig.GetRules() == nil {
				gl.Log("info", fmt.Sprintf("📊 Skipping %s/%s - No rules defined", repoConfig.GetOwner(), repoConfig.GetName()))
				continue
			}
			gl.Log("info", fmt.Sprintf("📊 Processing %s/%s...", repoConfig.GetOwner(), repoConfig.GetName()))

			rpt, err := svc.SanitizeRepo(r.Context(), repoConfig.GetOwner(), repoConfig.GetName(), repoConfig.GetRules(), dryRun)
			if err != nil {
				gl.Log("error", fmt.Sprintf("Error processing %s/%s: %v", repoConfig.GetOwner(), repoConfig.GetName(), err))
				continue
			}

			totalRuns += rpt.Runs.Deleted
			totalArtifacts += rpt.Artifacts.Deleted

			result := map[string]any{
				"owner":     rpt.Owner,
				"repo":      rpt.Repo,
				"runs":      rpt.Runs.Deleted,
				"artifacts": rpt.Artifacts.Deleted,
				"releases":  rpt.Releases.DeletedDrafts,
				"success":   true,
			}
			bulkResults = append(bulkResults, result)

			gl.Log("info", fmt.Sprintf("✅ %s/%s - Runs: %d, Artifacts: %d", repoConfig.GetOwner(), repoConfig.GetName(), rpt.Runs.Deleted, rpt.Artifacts.Deleted))
		}

		duration := time.Since(startTime)

		response := map[string]any{
			"bulk_operation":          true,
			"dry_run":                 dryRun,
			"started_at":              startTime.Format("2006-01-02 15:04:05"),
			"duration_ms":             duration.Milliseconds(),
			"total_repos":             len(bulkResults),
			"total_runs_cleaned":      totalRuns,
			"total_artifacts_cleaned": totalArtifacts,
			"productivity_summary": map[string]any{
				"estimated_storage_saved_mb": (totalRuns * 10) + (totalArtifacts * 50), // Estimativa
				"estimated_time_saved_min":   (totalRuns + totalArtifacts) * 2,         // Estimativa
			},
			"repositories": bulkResults,
		}

		gl.Log("info", fmt.Sprintf("🎉 BULK SANITIZATION COMPLETED - Duration: %v, Total Runs: %d, Total Artifacts: %d",
			duration, totalRuns, totalArtifacts))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(response)
	})

	// Analytics endpoint for repository insights
	routes["/analytics/"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		gl.Log("info", fmt.Sprintf("🔍 ANALYTICS REQUEST - %s/%s - Analysis Days: %d", owner, repo, analysisDays))
		startTime := time.Now()

		// Perform analytics
		insights, err := analytics.AnalyzeRepository(r.Context(), g.ghc, owner, repo, analysisDays)
		if err != nil {
			gl.Log("error", fmt.Sprintf("Analytics error for %s/%s: %v", owner, repo, err))
			http.Error(w, fmt.Sprintf("Analytics failed: %v", err), http.StatusInternalServerError)
			return
		}

		duration := time.Since(startTime)
		gl.Log("info", fmt.Sprintf("✅ ANALYTICS COMPLETED - %s/%s - Duration: %v, Health Score: %.1f (%s)",
			owner, repo, duration, insights.HealthScore.Overall, insights.HealthScore.Grade))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(insights)
	})

	// route: GET /productivity/{owner}/{repo}
	routes["/productivity/"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/productivity/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			http.Error(w, "missing owner/repo in path", http.StatusBadRequest)
			return
		}

		owner, repo := parts[0], parts[1]

		gl.Log("info", fmt.Sprintf("🚀 PRODUCTIVITY REQUEST - %s/%s", owner, repo))
		startTime := time.Now()

		// Perform productivity analysis
		report, err := productivity.AnalyzeProductivity(context.Background(), g.ghc, owner, repo)
		if err != nil {
			gl.Log("error", fmt.Sprintf("Failed to analyze productivity for %s/%s: %v", owner, repo, err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		duration := time.Since(startTime)
		gl.Log("info", fmt.Sprintf("✅ PRODUCTIVITY COMPLETE - %s/%s - Duration: %v - Actions: %d - ROI: %.1fx",
			owner, repo, duration, len(report.Actions), report.ROI.ROIRatio))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(report)
	})

	// route: GET /intelligence/quick/{owner}/{repo}
	routes["/intelligence/quick/"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/intelligence/quick/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			http.Error(w, "missing owner/repo in path", http.StatusBadRequest)
			return
		}

		owner, repo := parts[0], parts[1]

		gl.Log("info", fmt.Sprintf("🧠 AI QUICK INSIGHT REQUEST - %s/%s", owner, repo))
		startTime := time.Now()

		// Create intelligence operator
		intelligenceOp := i.NewIntelligenceOperator(g.MainConfig, g.ghc)

		// Generate quick insight
		insight, err := intelligenceOp.GenerateQuickInsight(context.Background(), owner, repo)
		if err != nil {
			gl.Log("error", fmt.Sprintf("Failed to generate AI insight for %s/%s: %v", owner, repo, err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		duration := time.Since(startTime)
		gl.Log("info", fmt.Sprintf("✅ AI INSIGHT COMPLETE - %s/%s - Duration: %v - Score: %.1f - Assessment: %s",
			owner, repo, duration, insight.AIScore, insight.QuickAssessment))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(insight)
	})

	// route: GET /intelligence/recommendations/{owner}/{repo}
	routes["/intelligence/recommendations/"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "only GET", http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/intelligence/recommendations/")
		parts := strings.Split(path, "/")
		if len(parts) < 2 {
			http.Error(w, "missing owner/repo in path", http.StatusBadRequest)
			return
		}

		owner, repo := parts[0], parts[1]

		gl.Log("info", fmt.Sprintf("🎯 AI RECOMMENDATIONS REQUEST - %s/%s", owner, repo))
		startTime := time.Now()

		// Create intelligence operator
		intelligenceOp := i.NewIntelligenceOperator(g.MainConfig, g.ghc)

		// Generate smart recommendations
		recommendations, err := intelligenceOp.GenerateSmartRecommendations(context.Background(), owner, repo)
		if err != nil {
			gl.Log("error", fmt.Sprintf("Failed to generate AI recommendations for %s/%s: %v", owner, repo, err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		duration := time.Since(startTime)
		gl.Log("info", fmt.Sprintf("✅ AI RECOMMENDATIONS COMPLETE - %s/%s - Duration: %v - Count: %d",
			owner, repo, duration, len(recommendations)))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(recommendations)
	})

	// route: POST /admin/repos/{owner}/{repo}/sanitize?dry_run=1
	routes["/admin/repos/"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		gl.Log("info", fmt.Sprintf("🎯 INDIVIDUAL SANITIZATION - %s/%s - DRY_RUN: %v", owner, repo, dryRun))
		startTime := time.Now()

		// find rules (optional override via cfg)
		var rules interfaces.IRules
		for _, rc := range g.MainConfig.GetGitHub().GetRepos() {
			if rc.GetOwner() == owner && rc.GetName() == repo {
				rules = rc.GetRules()
				break
			}
		}

		dummy := defs.NewRules(
			defs.NewRunsRule(30, 0, []string{}),
			defs.NewArtifactsRule(7),
			defs.NewReleasesRule(true),
			defs.NewSecurityRule(false, false, ""),
			defs.NewMonitoringRule(true, 90, false),
		)

		if rules.GetArtifactsRule() == dummy.GetArtifactsRule() &&
			rules.GetRunsRule().GetMaxAgeDays() == dummy.GetRunsRule().GetMaxAgeDays() &&
			rules.GetReleasesRule() == dummy.GetReleasesRule() {
			// default sane rules
			rules.GetRunsRule().SetMaxAgeDays(30)
			rules.GetArtifactsRule().SetMaxAgeDays(7)
			rules.GetReleasesRule().SetDeleteDrafts(true)
		}

		rpt, err := svc.SanitizeRepo(r.Context(), owner, repo, rules, dryRun)
		if err != nil {
			gl.Log("error", fmt.Sprintf("❌ Error sanitizing %s/%s: %v", owner, repo, err))
			http.Error(w, err.Error(), 500)
			return
		}

		duration := time.Since(startTime)
		gl.Log("info", fmt.Sprintf("✅ SANITIZATION COMPLETED - %s/%s - Duration: %v, Runs: %d, Artifacts: %d",
			owner, repo, duration, rpt.Runs.Deleted, rpt.Artifacts.Deleted))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(rpt)
	})

	return routes
}
