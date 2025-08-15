// Package sanitize defines the configuration types used for sanitizing GitHub repositories.
package sanitize

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

type Rules struct {
	Runs      RunsRule      `yaml:"runs"`
	Artifacts ArtifactsRule `yaml:"artifacts"`
	Releases  ReleasesRule  `yaml:"releases"`
}

type RepoCfg struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
	Rules Rules  `yaml:"rules"`
}

type Config struct {
	Runtime struct {
		DryRun    bool   `yaml:"dry_run"`
		ReportDir string `yaml:"report_dir"`
	} `yaml:"runtime"`
	GitHub struct {
		Repos []RepoCfg `yaml:"repos"`
	} `yaml:"github"`
}
