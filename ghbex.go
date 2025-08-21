// Package ghbex provides a set of utilities for working with GitHub repositories.
package ghbex

import (
	"context"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/config"
	"github.com/rafa-mori/ghbex/internal/defs/gitz"
	"github.com/rafa-mori/ghbex/internal/defs/interfaces"
	"github.com/rafa-mori/ghbex/internal/defs/notifiers"
	"github.com/rafa-mori/ghbex/internal/operators/analytics"
	"github.com/rafa-mori/ghbex/internal/operators/automation"
	ghserver "github.com/rafa-mori/ghbex/internal/server"
)

type MainConfig = interfaces.IMainConfig
type GHServerEngine = ghserver.GHServerEngine
type GitHub = gitz.GitHub

func NewMainConfigObj() (MainConfig, error) {
	return config.NewMainConfig(
		"",
		"",
		"",
		"",
		[]string{},
		false,
		true,
		false,
	)
}

/* OPERATORS - API EXPOSE (ABSTRACT) */

type OperatorStatus struct {
	Status   string
	Error    error
	Metadata map[string]any
}

type OperatorsRegistry interface {
	GetOperators() ([]string, error)
	AddOperator(name string, fn func(...any)) error
	RemoveOperator(name string) error
}

type OperatorsManager interface {
	DispatchOperator(name string, args ...any) (any, error)
	CancelOperator(name string) error
	IsOperatorRunning(name string) bool
	GetOperatorStatus(name string) (*OperatorStatus, error)
	MonitorOperator(name string, callback func(status *OperatorStatus)) error
}

/* OPERATORS - API EXPOSE (OLD VERSIONS) */

// ANALYTICS

type InsightsReport = analytics.InsightsReport

func AnalyzeRepository(ctx context.Context, client *github.Client, owner, repo string, analysisDays int) (*InsightsReport, error) {
	return analytics.AnalyzeRepository(ctx, client, owner, repo, analysisDays)
}
func GetRepositoryInsights(ctx context.Context, owner, repo string, days int) (*InsightsReport, error) {
	return analytics.GetRepositoryInsights(ctx, owner, repo, days)
}

// AUTOMATION

type AutomationReport = automation.AutomationReport

func AnalyzeAutomation(ctx context.Context, client *github.Client, owner, repo string, analysisDays int) (*AutomationReport, error) {
	return automation.AnalyzeAutomation(ctx, client, owner, repo, analysisDays)
}

type Service = automation.Service

func NewService(cli *github.Client, cfg interfaces.IMainConfig, ntf ...*notifiers.Discord) *Service {
	return automation.New(cli, cfg, ntf...)
}
