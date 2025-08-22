package gitz

import "github.com/rafa-mori/ghbex/internal/defs/interfaces"

type RepoCfg struct {
	Owner  string `yaml:"owner" json:"owner"`
	Name   string `yaml:"name" json:"name"`
	*Rules `yaml:"rules" json:"rules"`
}

func NewRepoCfgType(owner, name string, rules interfaces.IRules) *RepoCfg {
	rs := &Rules{}
	if rules != nil {
		rs = rules.(*Rules)
	}
	return &RepoCfg{
		Owner: owner,
		Name:  name,
		Rules: rs,
	}
}

func NewRepoCfg(owner, name string, rules interfaces.IRules) interfaces.IRepoCfg {
	return NewRepoCfgType(owner, name, rules)
}

func (r *RepoCfg) SetRules(rules interfaces.IRules) {
	if rules == nil {
		r.Rules = nil
	}
	rs, ok := rules.(*Rules)
	if !ok {
		return
	}
	r.Rules = rs
}
func (r *RepoCfg) GetMonitoring() interfaces.IMonitoringRule {
	if r.Rules != nil {
		return r.GetMonitoringRule()
	}
	return nil
}
func (r *RepoCfg) SetMonitoring(monitoring interfaces.IMonitoringRule) {
	if r.Rules != nil {
		r.SetMonitoringRule(monitoring)
	}
}

func (r *RepoCfg) GetOwner() string            { return r.Owner }
func (r *RepoCfg) SetOwner(owner string)       { r.Owner = owner }
func (r *RepoCfg) GetName() string             { return r.Name }
func (r *RepoCfg) SetName(name string)         { r.Name = name }
func (r *RepoCfg) GetRules() interfaces.IRules { return r.Rules }
