// Package defs contains type definitions used across the application.
package defs

import (
	"context"
	"time"
)

type INotifiers interface {
	// Send sends a notification with the given title and text.
	Send(ctx context.Context, title, text string, files ...Attachment) error
}

type Notifiers []struct {
	Type    string `yaml:"type"`
	Webhook string `yaml:"webhook"`
}

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

type GitHub struct {
	Auth  *GitHubAuth `yaml:"auth"`
	Repos []RepoCfg   `yaml:"repos"`
}

type RepoCfg struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
	Rules Rules  `yaml:"rules"`
}

type Config struct {
	Runtime `yaml:"runtime"`
	GitHub  `yaml:"github"`
}

type Report struct {
	Owner  string    `json:"owner"`
	Repo   string    `json:"repo"`
	When   time.Time `json:"when"`
	DryRun bool      `json:"dry_run"`

	Runs struct {
		Deleted int     `json:"deleted"`
		Kept    int     `json:"kept"`
		IDs     []int64 `json:"ids"`
	} `json:"runs"`

	Artifacts struct {
		Deleted int     `json:"deleted"`
		IDs     []int64 `json:"ids"`
	} `json:"artifacts"`

	Releases struct {
		DeletedDrafts int      `json:"deleted_drafts"`
		Tags          []string `json:"tags"`
	} `json:"releases"`

	Security struct {
		SSHKeysRotated int   `json:"ssh_keys_rotated"`
		OldKeysRemoved int   `json:"old_keys_removed"`
		NewKeyID       int64 `json:"new_key_id,omitempty"`
	} `json:"security"`

	Monitoring struct {
		IsInactive    bool `json:"is_inactive"`
		DaysInactive  int  `json:"days_inactive"`
		OpenPRs       int  `json:"open_prs"`
		OpenIssues    int  `json:"open_issues"`
		CommitsLast30 int  `json:"commits_last_30"`
	} `json:"monitoring"`

	Notes []string `json:"notes"`
}

type Attachment struct {
	Name string
	Body []byte
}

type RunsRule struct {
	MaxAgeDays      int      `yaml:"max_age_days"`
	KeepSuccessLast int      `yaml:"keep_success_last"`
	OnlyWorkflows   []string `yaml:"only_workflows"`
}

type ArtifactsRule struct {
	MaxAgeDays int `yaml:"max_age_days"`
}

type ReleasesRule struct {
	DeleteDrafts bool `yaml:"delete_drafts"`
}

type SecurityRule struct {
	RotateSSHKeys bool   `yaml:"rotate_ssh_keys"`
	RemoveOldKeys bool   `yaml:"remove_old_keys"`
	KeyPattern    string `yaml:"key_pattern"`
}

type MonitoringRule struct {
	CheckInactivity       bool `yaml:"check_inactivity"`
	InactiveDaysThreshold int  `yaml:"inactive_days_threshold"`
	MonitorPRs            bool `yaml:"monitor_prs"`
	MonitorIssues         bool `yaml:"monitor_issues"`
}

type Rules struct {
	Runs       RunsRule       `yaml:"runs"`
	Artifacts  ArtifactsRule  `yaml:"artifacts"`
	Releases   ReleasesRule   `yaml:"releases"`
	Security   SecurityRule   `yaml:"security"`
	Monitoring MonitoringRule `yaml:"monitoring"`
}
