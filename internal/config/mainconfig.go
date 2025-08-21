// Package config provides the main configuration for the GHBEX application.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rafa-mori/ghbex/internal/defs/common"
	"github.com/rafa-mori/ghbex/internal/defs/core"
	"github.com/rafa-mori/ghbex/internal/defs/gitz"
	"github.com/rafa-mori/ghbex/internal/defs/gromptz"
	"github.com/rafa-mori/ghbex/internal/defs/interfaces"

	"gopkg.in/yaml.v3"

	gl "github.com/rafa-mori/ghbex/internal/module/logger"
)

// MainConfig holds the main configuration for the GHBEX application.
type MainConfig struct {
	ConfigFilePath    string `yaml:"-" json:"-"`
	*core.Runtime     `yaml:"runtime" json:"runtime"`
	*core.Server      `yaml:"server" json:"server"`
	*gitz.GitHub      `yaml:"github" json:"github"`
	*common.Notifiers `mapstructure:",squash"`
	Grompt            gromptz.Grompt `yaml:"-" json:"-"`
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

	gromptEngineCfg := gromptz.NewGromptConfig(
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
		Runtime:        core.NewRuntimeType(debug, disableDryRun, reportDir, background),
		Server:         core.NewServerType(bindAddr, port),
		GitHub: gitz.NewGitHubType(
			gitz.NewGitHubAuthType(
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
		Notifiers: common.NewNotifiersType(
			common.NewNotifierType("slack", GetEnvOrDefault("SLACK_WEBHOOK_URL", "")),
			common.NewNotifierType("discord", GetEnvOrDefault("DISCORD_WEBHOOK_URL", "")),
			common.NewNotifierType("email", GetEnvOrDefault("EMAIL_SMTP_SERVER", "")),
		),
		Grompt: gromptz.NewPromptEngine(gromptEngineCfg),
	}

	// 🛡️ CRITICAL SECURITY: NEVER auto-discover all repositories!
	// Only use explicitly provided repositories to prevent accidental universe scanning

	// Priority order: CLI args -> ENV var -> empty (NEVER auto-discover)
	var explicitRepos []string
	if len(repositories) > 0 {
		explicitRepos = repositories
		gl.Log("info", "🎯 Using %d repositories from CLI arguments", len(explicitRepos))
	} else {
		repoListEnv := GetEnvOrDefault("REPO_LIST", "")
		if repoListEnv != "" {
			explicitRepos = strings.Split(repoListEnv, ",")
			gl.Log("info", "🎯 Using %d repositories from REPO_LIST env var", len(explicitRepos))
		} else {
			gl.Log("warning", "🚨 NO REPOSITORIES CONFIGURED - Using EMPTY list for safety")
			gl.Log("info", "📋 To configure repositories, use:")
			gl.Log("info", "   • CLI flag: --repos 'owner/repo1,owner/repo2'")
			gl.Log("info", "   • ENV var: REPO_LIST='owner/repo1,owner/repo2'")
			gl.Log("info", "🛡️ This prevents accidental scanning of all GitHub repositories")
			explicitRepos = []string{}
		}
	}

	// Process only explicitly configured repositories
	if len(explicitRepos) > 0 {
		gl.Log("info", "✅ Processing explicitly configured repositories:")
		repos := make([]*gitz.RepoCfg, 0, len(explicitRepos))
		for _, repoSpec := range explicitRepos {
			// Clean and validate repo specification
			repoSpec = strings.TrimSpace(repoSpec)
			if repoSpec == "" {
				continue
			}

			parts := strings.Split(repoSpec, "/")
			if len(parts) != 2 {
				gl.Log("warning", "⚠️ Invalid repository format '%s' - expected 'owner/repo'", repoSpec)
				continue
			}

			repoOwner := strings.TrimSpace(parts[0])
			repoName := strings.TrimSpace(parts[1])

			if repoOwner == "" || repoName == "" {
				gl.Log("warning", "⚠️ Invalid repository format '%s' - owner and repo cannot be empty", repoSpec)
				continue
			}

			gl.Log("info", "   📦 %s/%s", repoOwner, repoName)
			repos = append(repos, gitz.NewRepoCfgType(
				repoOwner,
				repoName,
				gitz.NewRules(
					gitz.NewRunsRuleType(
						GetEnvOrDefault("GITHUB_REPO_RUNS_MAX_AGE_DAYS", 30),
						GetEnvOrDefault("GITHUB_REPO_RUNS_MAX_PARALLEL", 5),
						[]string{"completed", "failure", "cancelled", "timed_out", "action_required"},
					),
					gitz.NewArtifactsRuleType(
						GetEnvOrDefault("GITHUB_REPO_ARTIFACTS_MAX_AGE_DAYS", 30),
					),
					gitz.NewReleasesRuleType(
						GetEnvOrDefault("GITHUB_REPO_RELEASES_DELETE_DRAFTS", false),
					),
					gitz.NewSecurityRuleType(
						GetEnvOrDefault("GITHUB_REPO_ROTATE_SSH_KEYS", false),
						GetEnvOrDefault("GITHUB_REPO_REMOVE_OLD_KEYS", false),
						GetEnvOrDefault("GITHUB_REPO_KEY_PATTERN", ""),
					),
					gitz.NewMonitoringRuleType(
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
		c.Runtime = core.NewRuntimeType(
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
		c.Server = core.NewServerType(
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
		c.GitHub = gitz.NewGitHubType(
			gitz.NewGitHubAuthType(
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
		notifiers := []*common.Notifier{}
		notifiers = append(notifiers, common.NewNotifierType("slack", GetEnvOrDefault("SLACK_WEBHOOK_URL", "")))
		notifiers = append(notifiers, common.NewNotifierType("discord", GetEnvOrDefault("DISCORD_WEBHOOK_URL", "")))
		notifiers = append(notifiers, common.NewNotifierType("email", GetEnvOrDefault("EMAIL_SMTP_SERVER", "")))
		c.Notifiers = common.NewNotifiersType(
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

func (c *MainConfig) GetGrompt() gromptz.PromptEngine {
	if c == nil {
		return nil
	}
	if c.Grompt == nil {
		c.Grompt = gromptz.NewGromptEngine(
			gromptz.NewGromptConfig(
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

func (c *MainConfig) SetGrompt(grompt gromptz.PromptEngine) {
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
	var obj any = &MainConfig{
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
