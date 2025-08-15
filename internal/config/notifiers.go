package ghconfig

import sanitize "github.com/rafa-mori/ghbex/internal/app"

type Runtime struct {
	DryRun    bool   `yaml:"dry_run"`
	ReportDir string `yaml:"report_dir"`
}

type Server struct {
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

type GitHubAuth struct {
	Kind           string `yaml:"kind"` // pat|app
	Token          string `yaml:"token"`
	AppID          int64  `yaml:"app_id"`
	InstallationID int64  `yaml:"installation_id"`
	PrivateKeyPath string `yaml:"private_key_path"`
	BaseURL        string `yaml:"base_url"`
	UploadURL      string `yaml:"upload_url"`
}

type Notifiers []struct {
	Type    string `yaml:"type"`
	Webhook string `yaml:"webhook"`
}

type GitHub struct {
	Auth  *GitHubAuth        `yaml:"auth"`
	Repos []sanitize.RepoCfg `yaml:"repos"`
}
