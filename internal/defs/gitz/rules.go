package gitz

import "github.com/kubex-ecosystem/ghbex/internal/defs/interfaces"

type Rules struct {
	*RunsRule       `yaml:"runs" json:"runs"`
	*ArtifactsRule  `yaml:"artifacts" json:"artifacts"`
	*ReleasesRule   `yaml:"releases" json:"releases"`
	*SecurityRule   `yaml:"security" json:"security"`
	*MonitoringRule `yaml:"monitoring" json:"monitoring"`
}

func NewRulesType(
	runs interfaces.IRunsRule,
	artifacts interfaces.IArtifactsRule,
	releases interfaces.IReleasesRule,
	security interfaces.ISecurityRule,
	monitoring interfaces.IMonitoringRule,
) *Rules {
	if runs == nil {
		runs = &RunsRule{}
	}
	if artifacts == nil {
		artifacts = &ArtifactsRule{}
	}
	if releases == nil {
		releases = &ReleasesRule{}
	}
	if security == nil {
		security = &SecurityRule{}
	}
	if monitoring == nil {
		monitoring = &MonitoringRule{}
	}

	return &Rules{
		RunsRule:       runs.(*RunsRule),
		ArtifactsRule:  artifacts.(*ArtifactsRule),
		ReleasesRule:   releases.(*ReleasesRule),
		SecurityRule:   security.(*SecurityRule),
		MonitoringRule: monitoring.(*MonitoringRule),
	}
}

func NewRules(
	runs interfaces.IRunsRule,
	artifacts interfaces.IArtifactsRule,
	releases interfaces.IReleasesRule,
	security interfaces.ISecurityRule,
	monitoring interfaces.IMonitoringRule,
) interfaces.IRules {
	return NewRulesType(
		runs,
		artifacts,
		releases,
		security,
		monitoring,
	)
}

func (r *Rules) GetRuns() interfaces.IRule                     { return r.RunsRule }
func (r *Rules) GetRuleName() string                           { return "rules" }
func (r *Rules) SetRuleName(name string)                       { /* // No-op for rules */ }
func (r *Rules) GetArtifacts() interfaces.IRule                { return r.ArtifactsRule }
func (r *Rules) GetReleases() interfaces.IRule                 { return r.ReleasesRule }
func (r *Rules) GetSecurity() interfaces.IRule                 { return r.SecurityRule }
func (r *Rules) GetMonitoring() interfaces.IRule               { return r.MonitoringRule }
func (r *Rules) GetSecurityRule() interfaces.ISecurityRule     { return r.SecurityRule }
func (r *Rules) GetMonitoringRule() interfaces.IMonitoringRule { return r.MonitoringRule }
func (r *Rules) GetReleasesRule() interfaces.IReleasesRule     { return r.ReleasesRule }
func (r *Rules) GetArtifactsRule() interfaces.IArtifactsRule   { return r.ArtifactsRule }
func (r *Rules) GetRunsRule() interfaces.IRunsRule             { return r.RunsRule }

func (r *Rules) SetRuns(rule interfaces.IRunsRule) {
	if rule == nil {
		r.RunsRule = nil
	} else {
		if _, ok := rule.(*RunsRule); !ok {
			r.RunsRule = rule.(*RunsRule)
		}
	}
}
func (r *Rules) SetArtifacts(rule interfaces.IArtifactsRule) {
	if rule == nil {
		r.ArtifactsRule = nil
	} else {
		if _, ok := rule.(*ArtifactsRule); !ok {
			r.ArtifactsRule = rule.(*ArtifactsRule)
		}
	}
}
func (r *Rules) SetReleases(rule interfaces.IReleasesRule) {
	if rule == nil {
		r.ReleasesRule = nil
	} else {
		if _, ok := rule.(*ReleasesRule); !ok {
			r.ReleasesRule = rule.(*ReleasesRule)
		}
	}
}
func (r *Rules) SetSecurity(rule interfaces.ISecurityRule) {
	if rule == nil {
		r.SecurityRule = nil
	} else {
		if _, ok := rule.(*SecurityRule); !ok {
			r.SecurityRule = rule.(*SecurityRule)
		}
	}
}
func (r *Rules) SetMonitoring(rule interfaces.IMonitoringRule) {
	if rule == nil {
		r.MonitoringRule = nil
	} else {
		if _, ok := rule.(*MonitoringRule); !ok {
			r.MonitoringRule = rule.(*MonitoringRule)
		}
	}
}
func (r *Rules) SetSecurityRule(rule interfaces.ISecurityRule) {
	if rule == nil {
		r.SecurityRule = nil
	} else {
		if _, ok := rule.(*SecurityRule); !ok {
			r.SecurityRule = rule.(*SecurityRule)
		}
	}
}
func (r *Rules) SetMonitoringRule(rule interfaces.IMonitoringRule) {
	if rule == nil {
		r.MonitoringRule = nil
	} else {
		if _, ok := rule.(*MonitoringRule); !ok {
			r.MonitoringRule = rule.(*MonitoringRule)
		}
	}
}
func (r *Rules) SetReleasesRule(rule interfaces.IReleasesRule) {
	if rule == nil {
		r.ReleasesRule = nil
	} else {
		if _, ok := rule.(*ReleasesRule); !ok {
			r.ReleasesRule = rule.(*ReleasesRule)
		}
	}
}
func (r *Rules) SetArtifactsRule(rule interfaces.IArtifactsRule) {
	if rule == nil {
		r.ArtifactsRule = nil
	} else {
		if _, ok := rule.(*ArtifactsRule); !ok {
			r.ArtifactsRule = rule.(*ArtifactsRule)
		}
	}
}
func (r *Rules) SetRunsRule(rule interfaces.IRunsRule) {
	if rule == nil {
		r.RunsRule = nil
	} else {
		if _, ok := rule.(*RunsRule); !ok {
			r.RunsRule = rule.(*RunsRule)
		}
	}
}
