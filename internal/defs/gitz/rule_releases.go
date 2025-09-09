package gitz

import "github.com/kubex-ecosystem/gemx/ghbex/internal/defs/interfaces"

type ReleasesRule struct {
	DeleteDrafts bool `yaml:"delete_drafts" json:"delete_drafts"`
}

func NewReleasesRuleType(deleteDrafts bool) *ReleasesRule {
	return &ReleasesRule{
		DeleteDrafts: deleteDrafts,
	}
}

func NewReleasesRule(deleteDrafts bool) interfaces.IReleasesRule {
	return NewReleasesRuleType(deleteDrafts)
}

func (r *ReleasesRule) GetDeleteDrafts() bool             { return r.DeleteDrafts }
func (r *ReleasesRule) SetDeleteDrafts(deleteDrafts bool) { r.DeleteDrafts = deleteDrafts }
func (r *ReleasesRule) GetRuleName() string               { return "releases" }
func (r *ReleasesRule) SetRuleName(name string)           { /* // No-op for releases rule */ }
