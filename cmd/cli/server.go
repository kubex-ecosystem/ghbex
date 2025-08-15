package cli

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
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	githubx "github.com/rafa-mori/ghbex/internal/ghclient"
	"github.com/rafa-mori/ghbex/internal/notify"
	"github.com/rafa-mori/ghbex/internal/sanitize"
)

type cfgRoot struct {
	Runtime struct {
		DryRun    bool   `yaml:"dry_run"`
		ReportDir string `yaml:"report_dir"`
	} `yaml:"runtime"`
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	GitHub struct {
		Auth struct {
			Kind           string `yaml:"kind"` // pat|app
			Token          string `yaml:"token"`
			AppID          int64  `yaml:"app_id"`
			InstallationID int64  `yaml:"installation_id"`
			PrivateKeyPath string `yaml:"private_key_path"`
			BaseURL        string `yaml:"base_url"`
			UploadURL      string `yaml:"upload_url"`
		} `yaml:"auth"`
		Repos []sanitize.RepoCfg `yaml:"repos"`
	} `yaml:"github"`
	Notifiers []struct {
		Type    string `yaml:"type"`
		Webhook string `yaml:"webhook"`
	} `yaml:"notifiers"`
}

func ServerCmdList() []*cobra.Command {
	var cmds []*cobra.Command

	// Define your server commands here
	cmds = append(cmds, startServer())
	cmds = append(cmds, stopServer())
	cmds = append(cmds, statusServer())
	return cmds
}

func stopServer() *cobra.Command {
	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the server",
		Annotations: GetDescriptions([]string{
			"This command stops server.",
			"This command stops the Grompt server and releases any resources held by it.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Stop server logic
		},
	}
	return stopCmd
}

func statusServer() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get server status",
		Annotations: GetDescriptions([]string{
			"This command gets the status of the server.",
			"This command checks if the server is running and returns its status.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Get server status logic
		},
	}
	return statusCmd
}

func startServer() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		Annotations: GetDescriptions([]string{
			"This command starts the server.",
			"This command initializes the server and starts waiting for help to build prompts.",
		}, false),
		Run: func(cmd *cobra.Command, args []string) {
			// Start server logic
		},
	}
	return startCmd
}

func main() {
	// load config
	b, err := os.ReadFile("config/sanitize.yaml")
	if err != nil {
		log.Printf("using example config: %v", err)
		b, _ = os.ReadFile("config/sanitize.yaml.example")
	}

	var cfg cfgRoot
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatal(err)
	}

	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8088"
	}

	// build github client
	ctx := context.Background()
	var ghc *github.Client
	switch strings.ToLower(cfg.GitHub.Auth.Kind) {
	case "pat":
		ghc, err = githubx.NewPAT(ctx, githubx.PATConfig{
			Token:     os.ExpandEnv(cfg.GitHub.Auth.Token),
			BaseURL:   cfg.GitHub.Auth.BaseURL,
			UploadURL: cfg.GitHub.Auth.UploadURL,
		})
	case "app":
		ghc, err = githubx.NewApp(ctx, githubx.AppConfig{
			AppID:          cfg.GitHub.Auth.AppID,
			InstallationID: cfg.GitHub.Auth.InstallationID,
			PrivateKeyPath: cfg.GitHub.Auth.PrivateKeyPath,
			BaseURL:        cfg.GitHub.Auth.BaseURL,
			UploadURL:      cfg.GitHub.Auth.UploadURL,
		})
	default:
		err = errors.New("github.auth.kind must be pat|app")
	}

	if err != nil {
		log.Fatal(err)
	}

	// notifiers
	var notifiers []sanitize.Notifier
	for _, n := range cfg.Notifiers {
		switch n.Type {
		case "discord":
			notifiers = append(notifiers, notify.Discord{Webhook: os.ExpandEnv(n.Webhook)})
		case "stdout":
			notifiers = append(notifiers, notify.Stdout{})
		}
	}

	// service
	svc := sanitize.New(ghc, sanitize.Config{
		Runtime: struct {
			DryRun    bool   "yaml:\"dry_run\""
			ReportDir string "yaml:\"report_dir\""
		}{DryRun: cfg.Runtime.DryRun, ReportDir: cfg.Runtime.ReportDir},
		GitHub: struct {
			Repos []sanitize.RepoCfg "yaml:\"repos\""
		}{Repos: cfg.GitHub.Repos},
	}, notifiers...)

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
		var rules sanitize.Rules
		for _, rc := range cfg.GitHub.Repos {
			if rc.Owner == owner && rc.Name == repo {
				rules = rc.Rules
				break
			}
		}

		var dummy sanitize.Rules
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
		Addr:              cfg.Server.Addr,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", cfg.Server.Addr)
	log.Fatal(srv.ListenAndServe())
}
