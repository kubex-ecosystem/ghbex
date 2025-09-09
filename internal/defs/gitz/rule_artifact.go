package gitz

import "github.com/kubex-ecosystem/gemx/ghbex/internal/defs/interfaces"

type ArtifactsRule struct {
	MaxAgeDays int `yaml:"max_age_days" json:"max_age_days"`
}

func NewArtifactsRuleType(maxAgeDays int) *ArtifactsRule {
	return &ArtifactsRule{
		MaxAgeDays: maxAgeDays,
	}
}

func NewArtifactsRule(maxAgeDays int) interfaces.IArtifactsRule {
	return NewArtifactsRuleType(maxAgeDays)
}

func (r *ArtifactsRule) GetMaxAgeDays() int      { return r.MaxAgeDays }
func (r *ArtifactsRule) SetMaxAgeDays(days int)  { r.MaxAgeDays = days }
func (r *ArtifactsRule) GetRuleName() string     { return "artifacts" }
func (r *ArtifactsRule) SetRuleName(name string) { /* // No-op for artifacts rule */ }
