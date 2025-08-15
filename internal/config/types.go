// Package config provides configuration structures and interfaces for the application.
package config

import (
	"context"
	"time"

	"github.com/rafa-mori/ghbex/internal/state"
)

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
	Owner string      `yaml:"owner"`
	Name  string      `yaml:"name"`
	Rules state.Rules `yaml:"rules"`
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

	Notes []string `json:"notes"`
}

type Notifier interface {
	Send(ctx context.Context, title, text string, files ...Attachment) error
}

type Attachment struct {
	Name string
	Body []byte
}
