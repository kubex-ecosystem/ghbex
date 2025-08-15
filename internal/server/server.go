package ghserver

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
	githubx "github.com/rafa-mori/ghbex/internal/client"
	config "github.com/rafa-mori/ghbex/internal/config"
	"github.com/rafa-mori/ghbex/internal/manager"
	notify "github.com/rafa-mori/ghbex/internal/notifiers"
	sanitize "github.com/rafa-mori/ghbex/internal/operators"
	"github.com/rafa-mori/ghbex/internal/state"
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
	switch strings.ToLower(gh.Auth.Kind) {
	case "pat":
		ghc, err = githubx.NewPAT(ctx, githubx.PATConfig{
			Token:     os.ExpandEnv(gh.Auth.Token),
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
	var notifierz []sanitize.Notifier
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

		// find rules (optional override via cfg)
		var rules state.Rules
		for _, rc := range g.MainConfig.GetGitHub().Repos {
			if rc.Owner == owner && rc.Name == repo {
				rules = rc.Rules
				break
			}
		}

		var dummy state.Rules
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
			http.Error(w, err.Error(), 500)
			return
		}
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
