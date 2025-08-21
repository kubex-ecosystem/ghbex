// Package config provides the main configuration for the GHBEX application.
package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	github "github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/client"
	"github.com/rafa-mori/ghbex/internal/defs"
	"github.com/rafa-mori/ghbex/internal/interfaces"
	"github.com/rafa-mori/grompt"

	"gopkg.in/yaml.v3"

	gl "github.com/rafa-mori/ghbex/internal/module/logger"
)

// MainConfig holds the main configuration for the GHBEX application.
type MainConfig struct {
	ConfigFilePath  string `yaml:"-" json:"-"`
	*defs.Runtime   `yaml:"runtime" json:"runtime"`
	*defs.Server    `yaml:"server" json:"server"`
	*defs.GitHub    `yaml:"github" json:"github"`
	*defs.Notifiers `mapstructure:",squash"`
	Grompt          defs.Grompt `yaml:"-" json:"-"`
}

func NewMainConfigObj() (interfaces.IMainConfig, error) {
	return NewMainConfig(
		"",
		"",
		"",
		"",
		[]string{},
		false,
		true,
		false,
	)
}

func NewMainConfig(
	bindAddr, port, reportDir, owner string,
	repositories []string,
	debug, disableDryRun, background bool,
) (interfaces.IMainConfig, error) {
	return NewMainConfigType(
		bindAddr,
		port,
		reportDir,
		owner,
		repositories,
		debug,
		disableDryRun,
		background,
	)
}

func NewMainConfigType(
	bindAddr, port, reportDir, owner string,
	repositories []string,
	debug, disableDryRun, background bool,
) (*MainConfig, error) {
	LoadEnvFromCurrentDir()
	if debug {
		gl.SetDebug(debug)
	}
	if bindAddr == "" {
		bindAddr = GetEnvOrDefault("GHBEX_BIND_ADDR", "0.0.0.0")
	}
	if port == "" {
		port = GetEnvOrDefault("GHBEX_PORT", "8088")
	}
	if reportDir == "" {
		reportDir = GetEnvOrDefault("GHBEX_REPORT_DIR", "reports")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		gl.Log("error", "Failed to get home directory: %v", err)
		return nil, fmt.Errorf("failed to get home directory: %v", err)
	}
	basePath := filepath.Join(homeDir, ".kubex", "ghbex")
	if reportDir != "" && !filepath.IsAbs(reportDir) {
		reportDir = filepath.Join(basePath, reportDir)
	}
	configFilePath := GetConfigFilePath(filepath.Join(basePath, "config", "sanitize.yaml"))
	if configFilePath == "" {
		gl.Log("error", "Configuration file not found. Please create a configuration file at %s", filepath.Join(basePath, "config", "sanitize.yaml"))
		return nil, fmt.Errorf("configuration file not found")
	}
	var ollamaEndpoint, claudeAPIKey, openAIKey, deepSeekKey, geminiKey string
	ollamaEndpoint = GetEnvOrDefault("OLLAMA_API_ENDPOINT", "")
	claudeAPIKey = GetEnvOrDefault("CLAUDE_API_KEY", "")
	openAIKey = GetEnvOrDefault("OPENAI_API_KEY", "")
	deepSeekKey = GetEnvOrDefault("DEEPSEEK_API_KEY", "")
	geminiKey = GetEnvOrDefault("GEMINI_API_KEY", "")

	gromptEngineCfg := defs.NewGromptConfig(
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeAPIKey,
		geminiKey,
	)
	owner = GetEnvOrDefault("GITHUB_REPO_OWNER", owner)
	if owner == "" {
		gl.Log("error", "GITHUB_REPO_OWNER environment variable is not set. Please set it to the owner of the GitHub repository.")
		return nil, fmt.Errorf("GITHUB_REPO_OWNER environment variable is not set")
	}
	cfg := &MainConfig{
		ConfigFilePath: configFilePath,
		Runtime:        defs.NewRuntimeType(debug, disableDryRun, reportDir, background),
		Server:         defs.NewServerType(bindAddr, port),
		GitHub: defs.NewGitHubType(
			defs.NewGitHubAuthType(
				"pat",
				GetEnvOrDefault("GITHUB_PAT_TOKEN", ""),
				GetEnvOrDefault[int64]("GITHUB_APP_ID", 0),
				GetEnvOrDefault[int64]("GITHUB_INSTALLATION_ID", 0),
				GetEnvOrDefault("GITHUB_PRIVATE_KEY_PATH", ""),
				GetEnvOrDefault("GITHUB_BASE_URL", ""),
				GetEnvOrDefault("GITHUB_UPLOAD_URL", ""),
			),
			make([]interfaces.IRepoCfg, 0),
		),
		Notifiers: defs.NewNotifiersType(
			defs.NewNotifierType("slack", GetEnvOrDefault("SLACK_WEBHOOK_URL", "")),
			defs.NewNotifierType("discord", GetEnvOrDefault("DISCORD_WEBHOOK_URL", "")),
			defs.NewNotifierType("email", GetEnvOrDefault("EMAIL_SMTP_SERVER", "")),
		),
		Grompt: defs.NewPromptEngine(gromptEngineCfg),
	}

	ghc, ghcErr := client.NewPAT(
		context.Background(),
		client.PATConfig{
			Token:     cfg.Token,
			BaseURL:   cfg.BaseURL,
			UploadURL: cfg.UploadURL,
		},
	)
	if ghcErr != nil {
		gl.Log("error", "Failed to create GitHub client: %v", ghcErr)
		return nil, fmt.Errorf("failed to create GitHub client: %v", ghcErr)
	}

	repoList, repoListResp, repoListErr := ghc.Repositories.ListByAuthenticatedUser(
		context.Background(),
		&github.RepositoryListByAuthenticatedUserOptions{
			Visibility:  "all",
			Affiliation: "owner",
			//Type:        "owner",
		},
	)
	if repoListErr != nil {
		gl.Log("error", "Failed to get repository list: %v", repoListErr)
		return nil, fmt.Errorf("failed to get repository list: %v", repoListErr)
	}

	if repoListResp.NextPage == 0 {
		gl.Log("debug", "No repositories found for the owner %s", owner)
	} else {
		gl.Log("debug", "Found %d repositories for the owner %s", len(repoList), owner)
		repos := make([]*defs.RepoCfg, 0, len(repoList))
		for _, repo := range repoList {
			repos = append(repos, defs.NewRepoCfgType(
				repo.GetOwner().GetLogin(),
				repo.GetName(),
				defs.NewRules(
					defs.NewRunsRuleType(
						GetEnvOrDefault("GITHUB_REPO_RUNS_MAX_AGE_DAYS", 30),
						GetEnvOrDefault("GITHUB_REPO_RUNS_MAX_PARALLEL", 5),
						[]string{"completed", "failure", "cancelled", "timed_out", "action_required"},
					),
					defs.NewArtifactsRuleType(
						GetEnvOrDefault("GITHUB_REPO_ARTIFACTS_MAX_AGE_DAYS", 30),
					),
					defs.NewReleasesRuleType(
						GetEnvOrDefault("GITHUB_REPO_RELEASES_DELETE_DRAFTS", false),
					),
					defs.NewSecurityRuleType(
						GetEnvOrDefault("GITHUB_REPO_ROTATE_SSH_KEYS", false),
						GetEnvOrDefault("GITHUB_REPO_REMOVE_OLD_KEYS", false),
						GetEnvOrDefault("GITHUB_REPO_KEY_PATTERN", ""),
					),
					defs.NewMonitoringRuleType(
						GetEnvOrDefault("GITHUB_REPO_CHECK_INACTIVITY", false),
						GetEnvOrDefault("GITHUB_REPO_INACTIVE_DAYS_THRESHOLD", 30),
						GetEnvOrDefault("GITHUB_REPO_MONITOR_PRS", true),
					),
				),
			))
		}
		cfg.Repos = append(cfg.Repos, repos...)
	}

	return cfg, nil
}

func (c *MainConfig) GetRuntime() interfaces.IRuntime {
	if c == nil {
		return nil
	}
	if c.Runtime == nil {
		c.Runtime = defs.NewRuntimeType(
			GetEnvOrDefault("GHBEX_DEBUG", false),
			GetEnvOrDefault("GHBEX_DRY_RUN", true),
			GetEnvOrDefault("GHBEX_REPORT_DIR", "reports"),
			GetEnvOrDefault("GHBEX_BACKGROUND", false),
		)
	}
	return c.Runtime
}

func (c *MainConfig) GetServer() interfaces.IServer {
	if c == nil {
		return nil
	}
	if c.Server == nil {
		c.Server = defs.NewServerType(
			GetEnvOrDefault("SERVER_HOST", "0.0.0.0"),
			GetEnvOrDefault("SERVER_PORT", "8080"),
		)
	}
	return c.Server
}

func (c *MainConfig) GetGitHub() interfaces.IGitHub {
	if c == nil {
		return nil
	}
	if c.GitHub == nil {
		c.GitHub = defs.NewGitHubType(
			defs.NewGitHubAuthType(
				"pat",
				GetEnvOrDefault("GITHUB_PAT_TOKEN", ""),
				GetEnvOrDefault[int64]("GITHUB_APP_ID", 0),
				GetEnvOrDefault[int64]("GITHUB_INSTALLATION_ID", 0),
				GetEnvOrDefault("GITHUB_PRIVATE_KEY_PATH", ""),
				GetEnvOrDefault("GITHUB_BASE_URL", ""),
				GetEnvOrDefault("GITHUB_UPLOAD_URL", ""),
			),
			[]interfaces.IRepoCfg{},
		)
	}
	return c.GitHub
}

func (c *MainConfig) GetNotifiers() interfaces.INotifiers {
	if c == nil {
		return nil
	}
	if c.Notifiers == nil {
		notifiers := []*defs.Notifier{}
		notifiers = append(notifiers, defs.NewNotifierType("slack", GetEnvOrDefault("SLACK_WEBHOOK_URL", "")))
		notifiers = append(notifiers, defs.NewNotifierType("discord", GetEnvOrDefault("DISCORD_WEBHOOK_URL", "")))
		notifiers = append(notifiers, defs.NewNotifierType("email", GetEnvOrDefault("EMAIL_SMTP_SERVER", "")))
		c.Notifiers = defs.NewNotifiersType(
			notifiers...,
		)
	}
	return c.Notifiers
}

func (c *MainConfig) GetConfigFilePath() string {
	if c == nil {
		return ""
	}
	return c.ConfigFilePath
}

func (c *MainConfig) GetGrompt() grompt.PromptEngine {
	if c == nil {
		return nil
	}
	if c.Grompt == nil {
		c.Grompt = defs.NewGromptEngine(
			defs.NewGromptConfig(
				c.Port,
				GetEnvOrDefault("OPENAI_API_KEY", ""),
				GetEnvOrDefault("DEEPSEEK_API_KEY", ""),
				GetEnvOrDefault("OLLAMA_API_ENDPOINT", ""),
				GetEnvOrDefault("CLAUDE_API_KEY", ""),
				GetEnvOrDefault("GEMINI_API_KEY", ""),
			),
		)
	}
	return c.Grompt
}

func (c *MainConfig) SetGrompt(grompt grompt.PromptEngine) {
	if c == nil {
		return
	}
	if grompt == nil {
		c.Grompt = nil
	} else {
		c.Grompt = grompt
	}
}

func (c *MainConfig) GetConfigObject() any {
	if c == nil {
		return nil
	}
	var obj any
	obj = &MainConfig{
		Runtime:   c.Runtime,
		Server:    c.Server,
		GitHub:    c.GitHub,
		Notifiers: c.Notifiers,
		Grompt:    c.Grompt,
	}
	return obj
}

func (c *MainConfig) String() string {
	if c == nil {
		return "nil"
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(data)
}
