// Package app defines the configuration types used for sanitizing GitHub repositories.
package app

import "github.com/rafa-mori/ghbex/internal/state"

type RepoCfg struct {
	Owner string      `yaml:"owner"`
	Name  string      `yaml:"name"`
	Rules state.Rules `yaml:"rules"`
}

type Runtime struct {
	DryRun    bool   `yaml:"dry_run"`
	ReportDir string `yaml:"report_dir"`
}

type GitHub struct {
	Repos []RepoCfg `yaml:"repos"`
}

type Config struct {
	Runtime `yaml:"runtime"`
	GitHub  `yaml:"github"`
}
