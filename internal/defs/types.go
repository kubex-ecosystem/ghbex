// Package defs contains type definitions used across the application.
package defs

import (
	"time"
)

type Config struct {
	Runtime `yaml:"runtime" json:"runtime"`
	GitHub  `yaml:"github" json:"github"`
}

type Runs struct {
	Deleted int     `yaml:"deleted" json:"deleted"`
	Kept    int     `yaml:"kept" json:"kept"`
	IDs     []int64 `yaml:"ids" json:"ids"`
}

type Artifacts struct {
	Deleted int     `yaml:"deleted" json:"deleted"`
	IDs     []int64 `yaml:"ids" json:"ids"`
}

type Releases struct {
	DeletedDrafts int      `yaml:"deleted_drafts" json:"deleted_drafts"`
	Tags          []string `yaml:"tags" json:"tags"`
}

type Security struct {
	SSHKeysRotated int   `yaml:"ssh_keys_rotated" json:"ssh_keys_rotated"`
	OldKeysRemoved int   `yaml:"old_keys_removed" json:"old_keys_removed"`
	NewKeyID       int64 `yaml:"new_key_id,omitempty" json:"new_key_id,omitempty"`
}

type Monitoring struct {
	IsInactive    bool `yaml:"is_inactive" json:"is_inactive"`
	DaysInactive  int  `yaml:"days_inactive" json:"days_inactive"`
	OpenPRs       int  `yaml:"open_prs" json:"open_prs"`
	OpenIssues    int  `yaml:"open_issues" json:"open_issues"`
	CommitsLast30 int  `yaml:"commits_last_30" json:"commits_last_30"`
}

type Report struct {
	Owner  string    `yaml:"owner" json:"owner"`
	Repo   string    `yaml:"repo" json:"repo"`
	When   time.Time `yaml:"when" json:"when"`
	DryRun bool      `yaml:"dry_run" json:"dry_run"`

	Runs Runs `yaml:"runs" json:"runs"`

	Artifacts Artifacts `yaml:"artifacts" json:"artifacts"`

	Releases Releases `yaml:"releases" json:"releases"`

	Security Security `yaml:"security" json:"security"`

	Monitoring Monitoring `yaml:"monitoring" json:"monitoring"`

	Notes []string `yaml:"notes" json:"notes"`
}
