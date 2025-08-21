package gitz

import "github.com/rafa-mori/ghbex/internal/interfaces"

type RunsRule struct {
	MaxAgeDays      int      `yaml:"max_age_days" json:"max_age_days"`
	KeepSuccessLast int      `yaml:"keep_success_last" json:"keep_success_last"`
	OnlyWorkflows   []string `yaml:"only_workflows" json:"only_workflows"`
}

func NewRunsRuleType(maxAgeDays, keepSuccessLast int, onlyWorkflows []string) *RunsRule {
	if onlyWorkflows == nil {
		onlyWorkflows = []string{}
	}
	return &RunsRule{
		MaxAgeDays:      maxAgeDays,
		KeepSuccessLast: keepSuccessLast,
		OnlyWorkflows:   onlyWorkflows,
	}
}

func NewRunsRule(maxAgeDays, keepSuccessLast int, onlyWorkflows []string) interfaces.IRunsRule {
	return NewRunsRuleType(maxAgeDays, keepSuccessLast, onlyWorkflows)
}

func (r *RunsRule) GetMaxAgeDays() int                  { return r.MaxAgeDays }
func (r *RunsRule) SetMaxAgeDays(days int)              { r.MaxAgeDays = days }
func (r *RunsRule) GetKeepSuccessLast() int             { return r.KeepSuccessLast }
func (r *RunsRule) SetKeepSuccessLast(days int)         { r.KeepSuccessLast = days }
func (r *RunsRule) GetOnlyWorkflows() []string          { return r.OnlyWorkflows }
func (r *RunsRule) SetOnlyWorkflows(workflows []string) { r.OnlyWorkflows = workflows }
func (r *RunsRule) GetRuleName() string                 { return "runs" }
func (r *RunsRule) SetRuleName(name string)             { /* // No-op for runs rule */ }
