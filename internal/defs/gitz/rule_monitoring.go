package gitz

import "github.com/rafa-mori/ghbex/internal/defs/interfaces"

type MonitoringRule struct {
	CheckInactivity       bool `yaml:"check_inactivity" json:"check_inactivity"`
	InactiveDaysThreshold int  `yaml:"inactive_days_threshold" json:"inactive_days_threshold"`
	MonitorPRs            bool `yaml:"monitor_prs" json:"monitor_prs"`
	MonitorIssues         bool `yaml:"monitor_issues" json:"monitor_issues"`
}

func NewMonitoringRuleType(checkInactivity bool, inactiveDaysThreshold int, monitorPRs bool) *MonitoringRule {
	return &MonitoringRule{
		CheckInactivity:       checkInactivity,
		InactiveDaysThreshold: inactiveDaysThreshold,
		MonitorPRs:            monitorPRs,
	}
}

func NewMonitoringRule(checkInactivity bool, inactiveDaysThreshold int, monitorPRs bool) interfaces.IMonitoringRule {
	return NewMonitoringRuleType(checkInactivity, inactiveDaysThreshold, monitorPRs)
}

func (r *MonitoringRule) GetCheckInactivity() bool          { return r.CheckInactivity }
func (r *MonitoringRule) SetCheckInactivity(check bool)     { r.CheckInactivity = check }
func (r *MonitoringRule) GetInactiveDaysThreshold() int     { return r.InactiveDaysThreshold }
func (r *MonitoringRule) SetInactiveDaysThreshold(days int) { r.InactiveDaysThreshold = days }
func (r *MonitoringRule) GetMonitorPRs() bool               { return r.MonitorPRs }
func (r *MonitoringRule) SetMonitorPRs(monitor bool)        { r.MonitorPRs = monitor }
func (r *MonitoringRule) GetRuleName() string               { return "monitoring" }
func (r *MonitoringRule) SetRuleName(name string)           { /* // No-op for monitoring rule */ }

func (r *MonitoringRule) GetArtifacts() interfaces.IArtifactsRule {
	if r == nil {
		return nil
	}

	return r.GetArtifacts()
}
