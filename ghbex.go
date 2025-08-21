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
	"github.com/rafa-mori/ghbex/internal/operators/intelligence"
	"github.com/rafa-mori/ghbex/internal/operators/monitoring"
	"github.com/rafa-mori/ghbex/internal/operators/productivity"
	"github.com/rafa-mori/ghbex/internal/operators/releases"
	"github.com/rafa-mori/ghbex/internal/operators/sanitize"
	"github.com/rafa-mori/ghbex/internal/operators/security"
	"github.com/rafa-mori/ghbex/internal/operators/workflows"
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

/* TYPES/INTERFACES - API EXPOSE */

type Rules = interfaces.IRules
type RepoCfg = interfaces.IRepoCfg

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
type LabelManagement = automation.LabelManagement
type IssueManagement = automation.IssueManagement
type PRManagement = automation.PRManagement
type WorkflowManagement = automation.WorkflowManagement
type AutomationAction = automation.AutomationAction

func AnalyzeAutomation(ctx context.Context, client *github.Client, owner, repo string, analysisDays int) (*AutomationReport, error) {
	return automation.AnalyzeAutomation(ctx, client, owner, repo, analysisDays)
}

type Service = automation.Service

func NewService(cli *github.Client, cfg interfaces.IMainConfig, ntf ...*notifiers.Discord) *Service {
	return automation.New(cli, cfg, ntf...)
}

/* OPERATORS - API EXPOSE (INTELLIGENCE) */

type LLMMetaResponse = intelligence.LLMMetaResponse
type IntelligenceOperator = intelligence.IntelligenceOperator
type RepositoryInsight = intelligence.RepositoryInsight
type SmartRecommendation = intelligence.SmartRecommendation
type HumanizedReport = intelligence.HumanizedReport
type OverallAssessment = intelligence.OverallAssessment
type KeyInsight = intelligence.KeyInsight
type ProductivityTip = intelligence.ProductivityTip
type RiskFactor = intelligence.RiskFactor
type NextStep = intelligence.NextStep

func NewIntelligenceOperator(cfg interfaces.IMainConfig, client *github.Client) *IntelligenceOperator {
	return intelligence.NewIntelligenceOperator(cfg, client)
}

/* OPERATORS - API EXPOSE (MONITORING) */

type ActivityReport = monitoring.ActivityReport
type PullRequestStats = monitoring.PullRequestStats
type IssueStats = monitoring.IssueStats
type CommitStats = monitoring.CommitStats

func AnalyzeRepositoryActivity(ctx context.Context, cli *github.Client, owner, repo string, inactiveDaysThreshold int) (*ActivityReport, error) {
	return monitoring.AnalyzeRepositoryActivity(ctx, cli, owner, repo, inactiveDaysThreshold)
}

func CheckInactiveRepositories(ctx context.Context, cli *github.Client, repos []struct{ Owner, Name string }, inactiveDaysThreshold int) ([]*ActivityReport, error) {
	return monitoring.CheckInactiveRepositories(ctx, cli, repos, inactiveDaysThreshold)
}

/* OPERATORS - API EXPOSE (PRODUCTIVITY) */

type ProductivityReport = productivity.ProductivityReport
type TemplateAnalysis = productivity.TemplateAnalysis
type IssueTemplate = productivity.IssueTemplate
type PRTemplate = productivity.PRTemplate
type BranchingOptimization = productivity.BranchingOptimization
type BranchAnalysis = productivity.BranchAnalysis
type StaleBranch = productivity.StaleBranch
type ActiveBranch = productivity.ActiveBranch
type MergePatterns = productivity.MergePatterns
type BranchProtection = productivity.BranchProtection
type AutoMergeAnalysis = productivity.AutoMergeAnalysis
type AutoMergePR = productivity.AutoMergePR
type AutoMergeRule = productivity.AutoMergeRule
type SafetyCheck = productivity.SafetyCheck
type NotificationOptimization = productivity.NotificationOptimization
type NotificationFilter = productivity.NotificationFilter
type PersonalizedRule = productivity.PersonalizedRule
type TeamNotification = productivity.TeamNotification
type WorkflowAutomation = productivity.WorkflowAutomation
type ExistingWorkflow = productivity.ExistingWorkflow
type SuggestedWorkflow = productivity.SuggestedWorkflow
type DeveloperExperience = productivity.DeveloperExperience
type ProductivityAction = productivity.ProductivityAction
type ROIEstimation = productivity.ROIEstimation

func AnalyzeProductivity(ctx context.Context, client *github.Client, owner, repo string) (*ProductivityReport, error) {
	return productivity.AnalyzeProductivity(ctx, client, owner, repo)
}

/* OPERATORS - API EXPOSE (RELEASES) */

func CleanReleases(ctx context.Context, cli *github.Client, owner, repo string, r interfaces.IReleasesRule, dry bool) (deletedDrafts int, tags []string, err error) {
	return releases.CleanReleases(ctx, cli, owner, repo, r, dry)
}

/* OPERATORS - API EXPOSE (SANITIZE) */

type IntelligentSanitizer = sanitize.IntelligentSanitizer
type SanitizationReport = sanitize.SanitizationReport
type SanitizationAction = sanitize.SanitizationAction
type ResourceSavings = sanitize.ResourceSavings
type SecurityImprovement = sanitize.SecurityImprovement
type QualityImprovement = sanitize.QualityImprovement

func NewIntelligentSanitizer(client *github.Client) *IntelligentSanitizer {
	return sanitize.NewIntelligentSanitizer(client)
}

/* OPERATORS - API EXPOSE (SECURITY) */

type SSHKeyPair = security.SSHKeyPair

func RotateSSHKeys(ctx context.Context, cli *github.Client, owner, repo string, dry bool) (*SSHKeyPair, error) {
	return security.RotateSSHKeys(ctx, cli, owner, repo, dry)
}

func ListDeployKeys(ctx context.Context, cli *github.Client, owner, repo string) ([]*github.Key, error) {
	return security.ListDeployKeys(ctx, cli, owner, repo)
}

/* OPERATORS - API EXPOSE (WORKFLOWS) */

func CleanWorkflowRuns(ctx context.Context, cli *github.Client, owner, repo string, r interfaces.IRunsRule, dry bool) (deleted, kept int, ids []int64, err error) {
	return workflows.CleanRuns(ctx, cli, owner, repo, r, dry)
}
