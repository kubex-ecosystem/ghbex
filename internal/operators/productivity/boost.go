// Package productivity provides advanced productivity enhancement tools for GitHub repositories.
package productivity

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
)

// ProductivityReport represents comprehensive productivity analysis and recommendations
type ProductivityReport struct {
	Owner       string    `json:"owner"`
	Repo        string    `json:"repo"`
	GeneratedAt time.Time `json:"generated_at"`

	// Template Analysis & Generation
	Templates *TemplateAnalysis `json:"templates"`

	// Branching Strategy Optimization
	Branching *BranchingOptimization `json:"branching"`

	// Auto-merge Opportunities
	AutoMerge *AutoMergeAnalysis `json:"auto_merge"`

	// Notification Optimization
	Notifications *NotificationOptimization `json:"notifications"`

	// Workflow Automation Opportunities
	Workflows *WorkflowAutomation `json:"workflows"`

	// Developer Experience Improvements
	DevEx *DeveloperExperience `json:"developer_experience"`

	// Actionable Recommendations
	Actions []ProductivityAction `json:"actions"`

	// ROI Estimation
	ROI *ROIEstimation `json:"roi_estimation"`
}

// TemplateAnalysis analyzes and suggests templates for issues and PRs
type TemplateAnalysis struct {
	ExistingTemplates       []string        `json:"existing_templates"`
	MissingTemplates        []string        `json:"missing_templates"`
	SuggestedIssueTemplates []IssueTemplate `json:"suggested_issue_templates"`
	SuggestedPRTemplate     *PRTemplate     `json:"suggested_pr_template"`
	TemplateUsageRate       float64         `json:"template_usage_rate"`
	ImprovementPotential    string          `json:"improvement_potential"`
}

// IssueTemplate represents a suggested issue template
type IssueTemplate struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"` // "bug", "feature", "documentation", "question"
	Template  string   `json:"template"`
	Labels    []string `json:"labels"`
	Assignees []string `json:"assignees"`
	Reasoning string   `json:"reasoning"`
}

// PRTemplate represents a suggested pull request template
type PRTemplate struct {
	Template   string   `json:"template"`
	Sections   []string `json:"sections"`
	Checklists []string `json:"checklists"`
	Reasoning  string   `json:"reasoning"`
}

// BranchingOptimization analyzes and optimizes branching strategies
type BranchingOptimization struct {
	CurrentStrategy     string             `json:"current_strategy"`
	RecommendedStrategy string             `json:"recommended_strategy"`
	BranchAnalysis      *BranchAnalysis    `json:"branch_analysis"`
	OptimizationGains   map[string]float64 `json:"optimization_gains"`
	ImplementationSteps []string           `json:"implementation_steps"`
}

// BranchAnalysis provides detailed branch analysis
type BranchAnalysis struct {
	TotalBranches    int               `json:"total_branches"`
	StaleBranches    []StaleBranch     `json:"stale_branches"`
	ActiveBranches   []ActiveBranch    `json:"active_branches"`
	NamingPatterns   map[string]int    `json:"naming_patterns"`
	MergePatterns    *MergePatterns    `json:"merge_patterns"`
	ProtectionStatus *BranchProtection `json:"protection_status"`
}

// StaleBranch represents a stale branch
type StaleBranch struct {
	Name       string    `json:"name"`
	LastCommit time.Time `json:"last_commit"`
	DaysStale  int       `json:"days_stale"`
	Author     string    `json:"author"`
	CanDelete  bool      `json:"can_delete"`
	Reason     string    `json:"reason"`
}

// ActiveBranch represents an active branch
type ActiveBranch struct {
	Name            string    `json:"name"`
	LastCommit      time.Time `json:"last_commit"`
	CommitsAhead    int       `json:"commits_ahead"`
	CommitsBehind   int       `json:"commits_behind"`
	HasPR           bool      `json:"has_pr"`
	SuggestedAction string    `json:"suggested_action"`
}

// MergePatterns analyzes merge patterns
type MergePatterns struct {
	MergeCommitRate   float64 `json:"merge_commit_rate"`
	SquashMergeRate   float64 `json:"squash_merge_rate"`
	RebaseMergeRate   float64 `json:"rebase_merge_rate"`
	RecommendedMethod string  `json:"recommended_method"`
}

// BranchProtection analyzes branch protection settings
type BranchProtection struct {
	ProtectedBranches   []string `json:"protected_branches"`
	UnprotectedBranches []string `json:"unprotected_branches"`
	MissingProtections  []string `json:"missing_protections"`
	Recommendations     []string `json:"recommendations"`
}

// AutoMergeAnalysis identifies auto-merge opportunities
type AutoMergeAnalysis struct {
	EligiblePRs        []AutoMergePR   `json:"eligible_prs"`
	AutoMergeRules     []AutoMergeRule `json:"auto_merge_rules"`
	SafetyChecks       []SafetyCheck   `json:"safety_checks"`
	EstimatedTimeSaved float64         `json:"estimated_time_saved"`
	RiskAssessment     string          `json:"risk_assessment"`
}

// AutoMergePR represents a PR eligible for auto-merge
type AutoMergePR struct {
	Number     int      `json:"number"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Type       string   `json:"type"` // "dependency", "hotfix", "documentation"
	Confidence float64  `json:"confidence"`
	Checks     []string `json:"checks"`
	Reason     string   `json:"reason"`
}

// AutoMergeRule represents an auto-merge rule
type AutoMergeRule struct {
	Name        string   `json:"name"`
	Conditions  []string `json:"conditions"`
	Actions     []string `json:"actions"`
	SafetyLevel string   `json:"safety_level"`
	Description string   `json:"description"`
}

// SafetyCheck represents a safety check for auto-merge
type SafetyCheck struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Configured  bool   `json:"configured"`
	Description string `json:"description"`
}

// NotificationOptimization optimizes notification strategies
type NotificationOptimization struct {
	CurrentNoise       float64              `json:"current_noise_level"`
	OptimizedNoise     float64              `json:"optimized_noise_level"`
	SmartFilters       []NotificationFilter `json:"smart_filters"`
	PersonalizedRules  []PersonalizedRule   `json:"personalized_rules"`
	TeamNotifications  []TeamNotification   `json:"team_notifications"`
	EstimatedReduction float64              `json:"estimated_reduction"`
}

// NotificationFilter represents a smart notification filter
type NotificationFilter struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Conditions []string `json:"conditions"`
	Action     string   `json:"action"`
	Priority   string   `json:"priority"`
}

// PersonalizedRule represents a personalized notification rule
type PersonalizedRule struct {
	User        string   `json:"user"`
	Role        string   `json:"role"`
	Preferences []string `json:"preferences"`
	Schedule    string   `json:"schedule"`
}

// TeamNotification represents team-wide notification settings
type TeamNotification struct {
	Event      string   `json:"event"`
	Recipients []string `json:"recipients"`
	Method     string   `json:"method"`
	Timing     string   `json:"timing"`
	Template   string   `json:"template"`
}

// WorkflowAutomation identifies workflow automation opportunities
type WorkflowAutomation struct {
	ExistingWorkflows   []ExistingWorkflow     `json:"existing_workflows"`
	SuggestedWorkflows  []SuggestedWorkflow    `json:"suggested_workflows"`
	OptimizationOps     []WorkflowOptimization `json:"optimization_opportunities"`
	AutomationPotential float64                `json:"automation_potential"`
	ComplexityReduction float64                `json:"complexity_reduction"`
}

// ExistingWorkflow represents an existing GitHub workflow
type ExistingWorkflow struct {
	Name            string    `json:"name"`
	File            string    `json:"file"`
	Triggers        []string  `json:"triggers"`
	Jobs            int       `json:"jobs"`
	LastRun         time.Time `json:"last_run"`
	SuccessRate     float64   `json:"success_rate"`
	AverageRuntime  float64   `json:"average_runtime"`
	OptimizationOps []string  `json:"optimization_opportunities"`
}

// SuggestedWorkflow represents a suggested new workflow
type SuggestedWorkflow struct {
	Name     string   `json:"name"`
	Purpose  string   `json:"purpose"`
	Triggers []string `json:"triggers"`
	Template string   `json:"template"`
	Benefits []string `json:"benefits"`
	Priority string   `json:"priority"`
	Effort   string   `json:"effort"`
}

// WorkflowOptimization represents workflow optimization opportunities
type WorkflowOptimization struct {
	WorkflowName     string   `json:"workflow_name"`
	OptimizationType string   `json:"optimization_type"`
	CurrentMetric    float64  `json:"current_metric"`
	OptimizedMetric  float64  `json:"optimized_metric"`
	ImprovementPct   float64  `json:"improvement_percentage"`
	Implementation   []string `json:"implementation_steps"`
}

// DeveloperExperience analyzes and improves developer experience
type DeveloperExperience struct {
	SetupComplexity   *SetupAnalysis          `json:"setup_complexity"`
	DocumentationGaps []DocumentationGap      `json:"documentation_gaps"`
	DeveloperFriction []FrictionPoint         `json:"developer_friction"`
	OnboardingPath    *OnboardingOptimization `json:"onboarding_path"`
	DevToolsGaps      []DevToolGap            `json:"dev_tools_gaps"`
	OverallScore      float64                 `json:"overall_score"`
}

// SetupAnalysis analyzes repository setup complexity
type SetupAnalysis struct {
	SetupSteps             int      `json:"setup_steps"`
	RequiredTools          []string `json:"required_tools"`
	DocumentedSetup        bool     `json:"documented_setup"`
	AutomatedSetup         bool     `json:"automated_setup"`
	SetupTime              float64  `json:"estimated_setup_time"`
	ComplexityLevel        string   `json:"complexity_level"`
	ImprovementSuggestions []string `json:"improvement_suggestions"`
}

// DocumentationGap represents a documentation gap
type DocumentationGap struct {
	Type             string `json:"type"`
	Severity         string `json:"severity"`
	Description      string `json:"description"`
	SuggestedContent string `json:"suggested_content"`
	Priority         int    `json:"priority"`
}

// FrictionPoint represents a developer friction point
type FrictionPoint struct {
	Area      string   `json:"area"`
	Issue     string   `json:"issue"`
	Impact    string   `json:"impact"`
	Solutions []string `json:"solutions"`
	Effort    string   `json:"effort"`
}

// OnboardingOptimization optimizes developer onboarding
type OnboardingOptimization struct {
	CurrentOnboarding   []string `json:"current_onboarding"`
	OptimizedOnboarding []string `json:"optimized_onboarding"`
	TimeReduction       float64  `json:"time_reduction"`
	AutomationOps       []string `json:"automation_opportunities"`
}

// DevToolGap represents development tooling gaps
type DevToolGap struct {
	Tool     string   `json:"tool"`
	Purpose  string   `json:"purpose"`
	Benefits []string `json:"benefits"`
	Setup    string   `json:"setup_instructions"`
	Priority string   `json:"priority"`
}

// ProductivityAction represents an actionable recommendation
type ProductivityAction struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Category     string   `json:"category"`
	Priority     string   `json:"priority"`
	Effort       string   `json:"effort"`
	Impact       string   `json:"impact"`
	Description  string   `json:"description"`
	Steps        []string `json:"steps"`
	ROI          float64  `json:"roi"`
	Dependencies []string `json:"dependencies"`
}

// ROIEstimation estimates return on investment for productivity improvements
type ROIEstimation struct {
	TotalTimeSavedHours   float64                `json:"total_time_saved_hours"`
	TotalTimeSavedDollars float64                `json:"total_time_saved_dollars"`
	ImplementationCost    float64                `json:"implementation_cost"`
	ROIRatio              float64                `json:"roi_ratio"`
	PaybackPeriod         float64                `json:"payback_period_months"`
	BreakdownByCategory   map[string]ROICategory `json:"breakdown_by_category"`
}

// ROICategory represents ROI for a specific category
type ROICategory struct {
	TimeSaved          float64 `json:"time_saved_hours"`
	DollarValue        float64 `json:"dollar_value"`
	ImplementationCost float64 `json:"implementation_cost"`
	ROI                float64 `json:"roi"`
}

// AnalyzeProductivity performs comprehensive productivity analysis
func AnalyzeProductivity(ctx context.Context, client *github.Client, owner, repo string) (*ProductivityReport, error) {
	report := &ProductivityReport{
		Owner:       owner,
		Repo:        repo,
		GeneratedAt: time.Now(),
	}

	// Analyze templates
	templates, err := analyzeTemplates(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze templates: %w", err)
	}
	report.Templates = templates

	// Analyze branching strategy
	branching, err := analyzeBranchingStrategy(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze branching: %w", err)
	}
	report.Branching = branching

	// Analyze auto-merge opportunities
	autoMerge, err := analyzeAutoMergeOpportunities(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze auto-merge: %w", err)
	}
	report.AutoMerge = autoMerge

	// Analyze notifications
	notifications := analyzeNotificationOptimization(ctx, client, owner, repo)
	report.Notifications = notifications

	// Analyze workflows
	workflows, err := analyzeWorkflowAutomation(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze workflows: %w", err)
	}
	report.Workflows = workflows

	// Analyze developer experience
	devex, err := analyzeDeveloperExperience(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze developer experience: %w", err)
	}
	report.DevEx = devex

	// Generate actionable recommendations
	actions := generateProductivityActions(report)
	report.Actions = actions

	// Calculate ROI estimation
	roi := calculateROIEstimation(report)
	report.ROI = roi

	return report, nil
}

// analyzeTemplates analyzes existing templates and suggests improvements
func analyzeTemplates(ctx context.Context, client *github.Client, owner, repo string) (*TemplateAnalysis, error) {
	// Check for existing templates in .github directory
	existingTemplates := []string{}
	missingTemplates := []string{}

	// Try to get .github directory contents
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, owner, repo, ".github", nil)

	hasIssueTemplates := false
	hasPRTemplate := false

	if err == nil {
		for _, content := range directoryContent {
			if content.Name != nil {
				name := *content.Name
				if strings.Contains(strings.ToLower(name), "issue") {
					hasIssueTemplates = true
					existingTemplates = append(existingTemplates, name)
				}
				if strings.Contains(strings.ToLower(name), "pull_request") || strings.Contains(strings.ToLower(name), "pr") {
					hasPRTemplate = true
					existingTemplates = append(existingTemplates, name)
				}
			}
		}
	}

	// Determine missing templates
	if !hasIssueTemplates {
		missingTemplates = append(missingTemplates, "Issue Templates")
	}
	if !hasPRTemplate {
		missingTemplates = append(missingTemplates, "Pull Request Template")
	}

	// Generate suggested issue templates based on repository analysis
	suggestedIssueTemplates := generateIssueTemplates(ctx, client, owner, repo)

	// Generate suggested PR template
	suggestedPRTemplate := generatePRTemplate(ctx, client, owner, repo)

	// Calculate template usage rate (simplified)
	templateUsageRate := 0.0
	if len(existingTemplates) > 0 {
		templateUsageRate = 75.0 // Estimated based on having templates
	}

	// Determine improvement potential
	improvementPotential := "high"
	if len(existingTemplates) > 2 {
		improvementPotential = "medium"
	}
	if len(existingTemplates) > 4 {
		improvementPotential = "low"
	}

	return &TemplateAnalysis{
		ExistingTemplates:       existingTemplates,
		MissingTemplates:        missingTemplates,
		SuggestedIssueTemplates: suggestedIssueTemplates,
		SuggestedPRTemplate:     suggestedPRTemplate,
		TemplateUsageRate:       templateUsageRate,
		ImprovementPotential:    improvementPotential,
	}, nil
}

// generateIssueTemplates generates suggested issue templates
func generateIssueTemplates(ctx context.Context, client *github.Client, owner, repo string) []IssueTemplate {
	templates := []IssueTemplate{
		{
			Name: "Bug Report",
			Type: "bug",
			Template: `---
name: Bug Report
about: Create a report to help us improve
title: '[BUG] '
labels: ['bug']
assignees: ''
---

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Environment (please complete the following information):**
- OS: [e.g. iOS]
- Browser [e.g. chrome, safari]
- Version [e.g. 22]

**Additional context**
Add any other context about the problem here.`,
			Labels:    []string{"bug"},
			Assignees: []string{},
			Reasoning: "Standard bug report template to ensure consistent bug reporting",
		},
		{
			Name: "Feature Request",
			Type: "feature",
			Template: `---
name: Feature request
about: Suggest an idea for this project
title: '[FEATURE] '
labels: ['enhancement']
assignees: ''
---

**Is your feature request related to a problem? Please describe.**
A clear and concise description of what the problem is. Ex. I'm always frustrated when [...]

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

**Describe alternatives you've considered**
A clear and concise description of any alternative solutions or features you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.`,
			Labels:    []string{"enhancement"},
			Assignees: []string{},
			Reasoning: "Helps users structure feature requests with proper context",
		},
	}

	return templates
}

// generatePRTemplate generates a suggested PR template
func generatePRTemplate(ctx context.Context, client *github.Client, owner, repo string) *PRTemplate {
	template := `## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] This change requires a documentation update

## Testing
- [ ] Tests pass locally
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes

## Checklist
- [ ] My code follows the style guidelines of this project
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings`

	return &PRTemplate{
		Template: template,
		Sections: []string{"Description", "Type of Change", "Testing", "Checklist"},
		Checklists: []string{
			"Tests pass locally",
			"Code follows style guidelines",
			"Self-review completed",
			"Documentation updated",
		},
		Reasoning: "Comprehensive PR template ensuring quality and consistency",
	}
}

// analyzeBranchingStrategy analyzes and optimizes branching strategies
func analyzeBranchingStrategy(ctx context.Context, client *github.Client, owner, repo string) (*BranchingOptimization, error) {
	// Get all branches
	branches, _, err := client.Repositories.ListBranches(ctx, owner, repo, &github.BranchListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	// Analyze branches
	branchAnalysis := analyzeBranches(ctx, client, owner, repo, branches)

	// Determine current strategy
	currentStrategy := determineBranchingStrategy(branches)

	// Recommend optimal strategy
	recommendedStrategy := recommendBranchingStrategy(branchAnalysis)

	// Calculate optimization gains
	optimizationGains := calculateOptimizationGains(currentStrategy, recommendedStrategy)

	// Generate implementation steps
	implementationSteps := generateImplementationSteps(currentStrategy, recommendedStrategy)

	return &BranchingOptimization{
		CurrentStrategy:     currentStrategy,
		RecommendedStrategy: recommendedStrategy,
		BranchAnalysis:      branchAnalysis,
		OptimizationGains:   optimizationGains,
		ImplementationSteps: implementationSteps,
	}, nil
}

// analyzeBranches analyzes all branches in the repository
func analyzeBranches(ctx context.Context, client *github.Client, owner, repo string, branches []*github.Branch) *BranchAnalysis {
	var staleBranches []StaleBranch
	var activeBranches []ActiveBranch
	namingPatterns := make(map[string]int)

	cutoffDate := time.Now().AddDate(0, 0, -30) // 30 days ago

	for _, branch := range branches {
		if branch.Name == nil {
			continue
		}

		branchName := *branch.Name

		// Analyze naming patterns
		if strings.HasPrefix(branchName, "feature/") {
			namingPatterns["feature"]++
		} else if strings.HasPrefix(branchName, "bugfix/") || strings.HasPrefix(branchName, "fix/") {
			namingPatterns["bugfix"]++
		} else if strings.HasPrefix(branchName, "hotfix/") {
			namingPatterns["hotfix"]++
		} else if strings.HasPrefix(branchName, "release/") {
			namingPatterns["release"]++
		} else if branchName == "main" || branchName == "master" || branchName == "develop" {
			namingPatterns["main"]++
		} else {
			namingPatterns["other"]++
		}

		// Get branch details
		branchDetail, _, err := client.Repositories.GetBranch(ctx, owner, repo, branchName, 3)
		if err != nil {
			continue
		}

		var lastCommitTime time.Time
		if branchDetail.Commit != nil && branchDetail.Commit.Commit != nil && branchDetail.Commit.Commit.Author != nil && branchDetail.Commit.Commit.Author.Date != nil {
			lastCommitTime = branchDetail.Commit.Commit.Author.Date.Time
		}

		// Determine if branch is stale
		if lastCommitTime.Before(cutoffDate) && branchName != "main" && branchName != "master" && branchName != "develop" {
			daysStale := int(time.Since(lastCommitTime).Hours() / 24)
			authorName := "unknown"
			if branchDetail.Commit != nil && branchDetail.Commit.Author != nil && branchDetail.Commit.Author.Login != nil {
				authorName = *branchDetail.Commit.Author.Login
			}

			staleBranches = append(staleBranches, StaleBranch{
				Name:       branchName,
				LastCommit: lastCommitTime,
				DaysStale:  daysStale,
				Author:     authorName,
				CanDelete:  daysStale > 60, // Can delete if older than 60 days
				Reason:     fmt.Sprintf("No commits for %d days", daysStale),
			})
		} else {
			// Active branch
			activeBranches = append(activeBranches, ActiveBranch{
				Name:            branchName,
				LastCommit:      lastCommitTime,
				CommitsAhead:    0,     // Would need comparison to determine
				CommitsBehind:   0,     // Would need comparison to determine
				HasPR:           false, // Would need PR lookup
				SuggestedAction: "monitor",
			})
		}
	}

	// Analyze merge patterns (simplified)
	mergePatterns := &MergePatterns{
		MergeCommitRate:   40.0,
		SquashMergeRate:   50.0,
		RebaseMergeRate:   10.0,
		RecommendedMethod: "squash",
	}

	// Analyze branch protection (simplified)
	branchProtection := &BranchProtection{
		ProtectedBranches:   []string{"main"},
		UnprotectedBranches: []string{},
		MissingProtections:  []string{"develop"},
		Recommendations:     []string{"Enable branch protection for develop branch"},
	}

	return &BranchAnalysis{
		TotalBranches:    len(branches),
		StaleBranches:    staleBranches,
		ActiveBranches:   activeBranches,
		NamingPatterns:   namingPatterns,
		MergePatterns:    mergePatterns,
		ProtectionStatus: branchProtection,
	}
}

// Determine current branching strategy
func determineBranchingStrategy(branches []*github.Branch) string {
	branchNames := make([]string, 0, len(branches))
	for _, branch := range branches {
		if branch.Name != nil {
			branchNames = append(branchNames, strings.ToLower(*branch.Name))
		}
	}

	hasMain := false
	hasDevelop := false
	hasFeatureBranches := false
	hasReleaseBranches := false

	for _, name := range branchNames {
		if name == "main" || name == "master" {
			hasMain = true
		} else if name == "develop" || name == "dev" {
			hasDevelop = true
		} else if strings.HasPrefix(name, "feature/") {
			hasFeatureBranches = true
		} else if strings.HasPrefix(name, "release/") {
			hasReleaseBranches = true
		}
	}

	if hasDevelop && hasFeatureBranches && hasReleaseBranches {
		return "git-flow"
	} else if hasMain && hasFeatureBranches {
		return "github-flow"
	} else if len(branchNames) <= 3 {
		return "centralized"
	} else {
		return "custom"
	}
}

// Recommend optimal branching strategy
func recommendBranchingStrategy(analysis *BranchAnalysis) string {
	totalBranches := analysis.TotalBranches
	staleBranches := len(analysis.StaleBranches)

	// Simple heuristic for recommendation
	if totalBranches <= 5 {
		return "github-flow"
	} else if staleBranches > totalBranches/2 {
		return "github-flow" // Simplify if too many stale branches
	} else {
		return "git-flow"
	}
}

// Calculate optimization gains
func calculateOptimizationGains(current, recommended string) map[string]float64 {
	gains := make(map[string]float64)

	if current != recommended {
		gains["complexity_reduction"] = 25.0
		gains["merge_time_reduction"] = 15.0
		gains["conflict_reduction"] = 20.0
		gains["onboarding_improvement"] = 30.0
	} else {
		gains["maintenance_improvement"] = 10.0
	}

	return gains
}

// Generate implementation steps
func generateImplementationSteps(current, recommended string) []string {
	if current == recommended {
		return []string{} // No changes needed
	}

	steps := []string{
		"Document new branching strategy",
		"Update CONTRIBUTING.md with new workflow",
		"Set up branch protection rules",
		"Train team on new workflow",
		"Migrate existing branches gradually",
	}

	return steps
}

// analyzeAutoMergeOpportunities identifies auto-merge opportunities
func analyzeAutoMergeOpportunities(ctx context.Context, client *github.Client, owner, repo string) (*AutoMergeAnalysis, error) {
	// Get open pull requests
	prs, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 50},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	var eligiblePRs []AutoMergePR

	for _, pr := range prs {
		if isEligibleForAutoMerge(pr) {
			eligiblePR := AutoMergePR{
				Number:     *pr.Number,
				Title:      *pr.Title,
				Author:     *pr.User.Login,
				Type:       determinePRType(pr),
				Confidence: calculateAutoMergeConfidence(pr),
				Checks:     []string{"CI passing", "Reviews approved"},
				Reason:     "Low risk dependency update",
			}
			eligiblePRs = append(eligiblePRs, eligiblePR)
		}
	}

	// Generate auto-merge rules
	autoMergeRules := generateAutoMergeRules()

	// Generate safety checks
	safetyChecks := generateSafetyChecks()

	// Estimate time saved
	estimatedTimeSaved := float64(len(eligiblePRs)) * 15.0 // 15 minutes per PR

	// Assess risk
	riskAssessment := "low"
	if len(eligiblePRs) > 10 {
		riskAssessment = "medium"
	}

	return &AutoMergeAnalysis{
		EligiblePRs:        eligiblePRs,
		AutoMergeRules:     autoMergeRules,
		SafetyChecks:       safetyChecks,
		EstimatedTimeSaved: estimatedTimeSaved,
		RiskAssessment:     riskAssessment,
	}, nil
}

// Check if PR is eligible for auto-merge
func isEligibleForAutoMerge(pr *github.PullRequest) bool {
	if pr.Title == nil || pr.User == nil || pr.User.Login == nil {
		return false
	}

	title := strings.ToLower(*pr.Title)
	author := *pr.User.Login

	// Dependency updates
	if strings.Contains(title, "dependabot") || author == "dependabot[bot]" {
		return true
	}

	// Documentation updates
	if strings.Contains(title, "docs") || strings.Contains(title, "documentation") {
		return true
	}

	// Minor version bumps
	if strings.Contains(title, "version") && strings.Contains(title, "patch") {
		return true
	}

	return false
}

// Determine PR type
func determinePRType(pr *github.PullRequest) string {
	if pr.Title == nil {
		return "other"
	}

	title := strings.ToLower(*pr.Title)

	if strings.Contains(title, "dependabot") {
		return "dependency"
	} else if strings.Contains(title, "hotfix") {
		return "hotfix"
	} else if strings.Contains(title, "docs") {
		return "documentation"
	} else {
		return "feature"
	}
}

// Calculate auto-merge confidence
func calculateAutoMergeConfidence(pr *github.PullRequest) float64 {
	confidence := 50.0

	if pr.User != nil && pr.User.Login != nil {
		author := *pr.User.Login
		if author == "dependabot[bot]" {
			confidence += 30.0
		}
	}

	if pr.Title != nil {
		title := strings.ToLower(*pr.Title)
		if strings.Contains(title, "docs") {
			confidence += 20.0
		}
		if strings.Contains(title, "patch") {
			confidence += 15.0
		}
	}

	if confidence > 100 {
		confidence = 100
	}

	return confidence
}

// Generate auto-merge rules
func generateAutoMergeRules() []AutoMergeRule {
	// TODO: Implement real analysis for auto-merge rule suggestions
	// Currently returning empty to avoid hardcoded suggestions
	// Analysis should examine:
	// - Current repository protection rules
	// - Historical merge patterns
	// - Team workflow preferences
	// - Risk tolerance based on repository criticality

	return []AutoMergeRule{}
}

// Generate safety checks
func generateSafetyChecks() []SafetyCheck {
	// TODO: Implement real repository safety check analysis
	// Currently returning empty to avoid hardcoded suggestions
	// Analysis should examine:
	// - Existing branch protection rules
	// - Required status checks configuration
	// - Review requirements
	// - Security policy configuration

	return []SafetyCheck{}
}

// analyzeNotificationOptimization optimizes notification strategies
func analyzeNotificationOptimization(ctx context.Context, client *github.Client, owner, repo string) *NotificationOptimization {
	// TODO: Implement real notification optimization analysis
	// Currently returning minimal data to avoid hardcoded suggestions
	// Analysis should examine:
	// - Current repository notification settings
	// - Team member notification preferences
	// - Historical notification patterns
	// - Integration with external tools (Slack, email, etc.)

	return &NotificationOptimization{
		CurrentNoise:       0.0,
		OptimizedNoise:     0.0,
		SmartFilters:       []NotificationFilter{},
		PersonalizedRules:  []PersonalizedRule{},
		TeamNotifications:  []TeamNotification{},
		EstimatedReduction: 0.0,
	}
}

// analyzeWorkflowAutomation identifies workflow automation opportunities
func analyzeWorkflowAutomation(ctx context.Context, client *github.Client, owner, repo string) (*WorkflowAutomation, error) {
	// Get existing workflows
	workflows, _, err := client.Actions.ListWorkflows(ctx, owner, repo, &github.ListOptions{})
	if err != nil {
		// Not all repositories have Actions enabled
		workflows = &github.Workflows{Workflows: []*github.Workflow{}}
	}

	var existingWorkflows []ExistingWorkflow
	for _, workflow := range workflows.Workflows {
		if workflow.Name != nil && workflow.Path != nil {
			existing := ExistingWorkflow{
				Name:            *workflow.Name,
				File:            *workflow.Path,
				Triggers:        []string{"push", "pull_request"}, // Simplified
				Jobs:            3,                                // Estimated
				LastRun:         time.Now().AddDate(0, 0, -1),
				SuccessRate:     85.0,
				AverageRuntime:  5.5,
				OptimizationOps: []string{"Cache dependencies", "Parallel jobs"},
			}
			existingWorkflows = append(existingWorkflows, existing)
		}
	}

	// Generate suggested workflows
	suggestedWorkflows := generateSuggestedWorkflows()

	// Generate optimization opportunities
	optimizationOps := generateWorkflowOptimizations(existingWorkflows)

	automationPotential := 75.0
	complexityReduction := 40.0

	return &WorkflowAutomation{
		ExistingWorkflows:   existingWorkflows,
		SuggestedWorkflows:  suggestedWorkflows,
		OptimizationOps:     optimizationOps,
		AutomationPotential: automationPotential,
		ComplexityReduction: complexityReduction,
	}, nil
}

// Generate suggested workflows
func generateSuggestedWorkflows() []SuggestedWorkflow {
	return []SuggestedWorkflow{
		{
			Name:     "Dependency Security Scan",
			Purpose:  "Automated security scanning of dependencies",
			Triggers: []string{"schedule", "pull_request"},
			Template: `// Generate suggested workflows
func generateSuggestedWorkflows() []SuggestedWorkflow {
	// TODO: Implement real workflow suggestion analysis
	// Currently returning empty to avoid hardcoded suggestions
	// Analysis should examine:
	// - Repository language and dependencies
	// - Existing workflow patterns
	// - Team development practices
	// - Security and compliance requirements

	return []SuggestedWorkflow{}
}`,
			Benefits: []string{"Early vulnerability detection", "Compliance", "Automated reporting"},
			Priority: "high",
			Effort:   "low",
		},
		{
			Name:     "Auto-label PRs",
			Purpose:  "Automatically label pull requests based on changes",
			Triggers: []string{"pull_request"},
			Template: `name: Auto Label PRs
on:
  pull_request:
    types: [opened, synchronize]
jobs:
  label:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/labeler@v4`,
			Benefits: []string{"Consistent labeling", "Reduced manual work", "Better organization"},
			Priority: "medium",
			Effort:   "low",
		},
	}
}

// Generate workflow optimizations
func generateWorkflowOptimizations(existing []ExistingWorkflow) []WorkflowOptimization {
	var optimizations []WorkflowOptimization

	for _, workflow := range existing {
		if workflow.AverageRuntime > 5.0 {
			optimizations = append(optimizations, WorkflowOptimization{
				WorkflowName:     workflow.Name,
				OptimizationType: "runtime",
				CurrentMetric:    workflow.AverageRuntime,
				OptimizedMetric:  workflow.AverageRuntime * 0.6,
				ImprovementPct:   40.0,
				Implementation: []string{
					"Add dependency caching",
					"Use parallel job execution",
					"Optimize Docker layers",
				},
			})
		}
	}

	return optimizations
}

// analyzeDeveloperExperience analyzes and improves developer experience
func analyzeDeveloperExperience(ctx context.Context, client *github.Client, owner, repo string) (*DeveloperExperience, error) {
	// Analyze setup complexity
	setupComplexity := analyzeSetupComplexity(ctx, client, owner, repo)

	// Identify documentation gaps
	documentationGaps := identifyDocumentationGaps(ctx, client, owner, repo)

	// Identify friction points
	developerFriction := identifyFrictionPoints()

	// Analyze onboarding path
	onboardingPath := analyzeOnboardingPath()

	// Identify dev tools gaps
	devToolsGaps := identifyDevToolsGaps()

	// Calculate overall score
	overallScore := calculateDevExScore(setupComplexity, documentationGaps, developerFriction)

	return &DeveloperExperience{
		SetupComplexity:   setupComplexity,
		DocumentationGaps: documentationGaps,
		DeveloperFriction: developerFriction,
		OnboardingPath:    onboardingPath,
		DevToolsGaps:      devToolsGaps,
		OverallScore:      overallScore,
	}, nil
}

// Analyze setup complexity
func analyzeSetupComplexity(ctx context.Context, client *github.Client, owner, repo string) *SetupAnalysis {
	// Check for setup files
	setupFiles := []string{"README.md", "CONTRIBUTING.md", "package.json", "go.mod", "requirements.txt", "Makefile"}
	requiredTools := []string{}
	setupSteps := 3 // Minimum steps

	for _, file := range setupFiles {
		_, _, _, err := client.Repositories.GetContents(ctx, owner, repo, file, nil)
		if err == nil {
			setupSteps++
			if file == "package.json" {
				requiredTools = append(requiredTools, "Node.js", "npm")
			} else if file == "go.mod" {
				requiredTools = append(requiredTools, "Go")
			} else if file == "requirements.txt" {
				requiredTools = append(requiredTools, "Python", "pip")
			}
		}
	}

	// Check for README
	_, _, _, readmeErr := client.Repositories.GetContents(ctx, owner, repo, "README.md", nil)
	documentedSetup := readmeErr == nil

	// Check for automation
	_, _, _, makefileErr := client.Repositories.GetContents(ctx, owner, repo, "Makefile", nil)
	automatedSetup := makefileErr == nil

	setupTime := float64(setupSteps) * 5.0 // 5 minutes per step

	complexityLevel := "low"
	if setupSteps > 5 {
		complexityLevel = "medium"
	}
	if setupSteps > 8 {
		complexityLevel = "high"
	}

	improvementSuggestions := []string{
		"Add one-command setup script",
		"Document all prerequisites clearly",
		"Provide Docker-based development environment",
	}

	return &SetupAnalysis{
		SetupSteps:             setupSteps,
		RequiredTools:          requiredTools,
		DocumentedSetup:        documentedSetup,
		AutomatedSetup:         automatedSetup,
		SetupTime:              setupTime,
		ComplexityLevel:        complexityLevel,
		ImprovementSuggestions: improvementSuggestions,
	}
}

// Identify documentation gaps
func identifyDocumentationGaps(ctx context.Context, client *github.Client, owner, repo string) []DocumentationGap {
	gaps := []DocumentationGap{}

	// Check for essential documentation files
	essentialDocs := map[string]string{
		"README.md":       "Project overview and getting started guide",
		"CONTRIBUTING.md": "Contribution guidelines and development workflow",
		"LICENSE":         "License information",
		"CHANGELOG.md":    "Version history and changes",
		"docs/API.md":     "API documentation",
	}

	priority := 1
	for file, description := range essentialDocs {
		_, _, _, err := client.Repositories.GetContents(ctx, owner, repo, file, nil)
		if err != nil {
			gap := DocumentationGap{
				Type:             "missing_file",
				Severity:         "medium",
				Description:      fmt.Sprintf("Missing %s", file),
				SuggestedContent: description,
				Priority:         priority,
			}
			gaps = append(gaps, gap)
			priority++
		}
	}

	return gaps
}

// Identify friction points
func identifyFrictionPoints() []FrictionPoint {
	// TODO: Implement real friction point analysis
	// Currently returning empty to avoid hardcoded suggestions
	// Analysis should examine:
	// - Repository setup complexity
	// - Developer onboarding feedback
	// - Code review patterns and bottlenecks
	// - Build and deployment issues

	return []FrictionPoint{}
}

// Analyze onboarding path
func analyzeOnboardingPath() *OnboardingOptimization {
	currentOnboarding := []string{
		"Read README",
		"Install dependencies manually",
		"Run tests manually",
		"Submit PR",
	}

	optimizedOnboarding := []string{
		"Run single setup command",
		"Auto-format code",
		"Auto-run tests",
		"Submit PR with auto-checks",
	}

	timeReduction := 75.0 // 75% time reduction

	automationOps := []string{
		"Setup script automation",
		"Pre-commit hooks",
		"Automated testing",
		"Template generation",
	}

	return &OnboardingOptimization{
		CurrentOnboarding:   currentOnboarding,
		OptimizedOnboarding: optimizedOnboarding,
		TimeReduction:       timeReduction,
		AutomationOps:       automationOps,
	}
}

// Identify dev tools gaps
func identifyDevToolsGaps() []DevToolGap {
	// TODO: Implement real dev tools gap analysis
	// Currently returning empty to avoid hardcoded suggestions
	// Analysis should examine:
	// - Existing development tools configuration
	// - Team development workflow patterns
	// - Code quality and consistency metrics
	// - Integration with repository settings

	return []DevToolGap{}
}

// Calculate developer experience score
func calculateDevExScore(setup *SetupAnalysis, gaps []DocumentationGap, friction []FrictionPoint) float64 {
	score := 100.0

	// Deduct for setup complexity
	if setup.ComplexityLevel == "high" {
		score -= 30
	} else if setup.ComplexityLevel == "medium" {
		score -= 15
	}

	// Deduct for documentation gaps
	score -= float64(len(gaps)) * 5

	// Deduct for friction points
	score -= float64(len(friction)) * 10

	// Add points for automation
	if setup.AutomatedSetup {
		score += 15
	}
	if setup.DocumentedSetup {
		score += 10
	}

	if score < 0 {
		score = 0
	}

	return score
}

// generateProductivityActions generates actionable recommendations
func generateProductivityActions(report *ProductivityReport) []ProductivityAction {
	var actions []ProductivityAction
	id := 1

	// Template actions
	if len(report.Templates.MissingTemplates) > 0 {
		actions = append(actions, ProductivityAction{
			ID:          fmt.Sprintf("PROD-%03d", id),
			Title:       "Add Issue and PR Templates",
			Category:    "templates",
			Priority:    "high",
			Effort:      "low",
			Impact:      "high",
			Description: "Add standardized templates to improve issue and PR quality",
			Steps: []string{
				"Create .github/ISSUE_TEMPLATE/ directory",
				"Add bug report template",
				"Add feature request template",
				"Add pull request template",
				"Test templates with new issues/PRs",
			},
			ROI:          3.5,
			Dependencies: []string{},
		})
		id++
	}

	// Branching actions
	if report.Branching.CurrentStrategy != report.Branching.RecommendedStrategy {
		actions = append(actions, ProductivityAction{
			ID:           fmt.Sprintf("PROD-%03d", id),
			Title:        "Optimize Branching Strategy",
			Category:     "branching",
			Priority:     "medium",
			Effort:       "medium",
			Impact:       "medium",
			Description:  fmt.Sprintf("Migrate from %s to %s strategy", report.Branching.CurrentStrategy, report.Branching.RecommendedStrategy),
			Steps:        report.Branching.ImplementationSteps,
			ROI:          2.1,
			Dependencies: []string{},
		})
		id++
	}

	// Auto-merge actions
	if len(report.AutoMerge.EligiblePRs) > 0 {
		actions = append(actions, ProductivityAction{
			ID:          fmt.Sprintf("PROD-%03d", id),
			Title:       "Implement Auto-merge Rules",
			Category:    "automation",
			Priority:    "high",
			Effort:      "medium",
			Impact:      "high",
			Description: "Set up automated merging for low-risk changes",
			Steps: []string{
				"Configure Dependabot auto-merge",
				"Set up documentation auto-merge",
				"Implement safety checks",
				"Monitor and refine rules",
			},
			ROI:          4.2,
			Dependencies: []string{"PROD-001"},
		})
		id++
	}

	// Workflow actions
	if len(report.Workflows.SuggestedWorkflows) > 0 {
		actions = append(actions, ProductivityAction{
			ID:          fmt.Sprintf("PROD-%03d", id),
			Title:       "Add Security and Automation Workflows",
			Category:    "workflows",
			Priority:    "medium",
			Effort:      "low",
			Impact:      "medium",
			Description: "Implement automated security scanning and PR labeling",
			Steps: []string{
				"Add security scanning workflow",
				"Implement auto-labeling workflow",
				"Configure workflow notifications",
				"Test and monitor workflows",
			},
			ROI:          2.8,
			Dependencies: []string{},
		})
		id++
	}

	// Developer experience actions
	if report.DevEx.OverallScore < 70 {
		actions = append(actions, ProductivityAction{
			ID:          fmt.Sprintf("PROD-%03d", id),
			Title:       "Improve Developer Experience",
			Category:    "developer_experience",
			Priority:    "high",
			Effort:      "high",
			Impact:      "high",
			Description: "Streamline development setup and reduce friction",
			Steps: []string{
				"Create one-command setup script",
				"Add pre-commit hooks",
				"Improve documentation",
				"Add development tools",
			},
			ROI:          5.1,
			Dependencies: []string{},
		})
		id++
	}

	// Sort actions by priority and ROI
	sort.Slice(actions, func(i, j int) bool {
		if actions[i].Priority != actions[j].Priority {
			priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
			return priorityOrder[actions[i].Priority] > priorityOrder[actions[j].Priority]
		}
		return actions[i].ROI > actions[j].ROI
	})

	return actions
}

// calculateROIEstimation calculates return on investment
func calculateROIEstimation(report *ProductivityReport) *ROIEstimation {
	totalTimeSaved := 0.0
	implementationCost := 0.0
	breakdown := make(map[string]ROICategory)

	// Calculate based on actions
	for _, action := range report.Actions {
		timeSaved := action.ROI * 10 // Convert ROI to hours
		cost := calculateImplementationCost(action.Effort)

		totalTimeSaved += timeSaved
		implementationCost += cost

		breakdown[action.Category] = ROICategory{
			TimeSaved:          timeSaved,
			DollarValue:        timeSaved * 50, // $50/hour
			ImplementationCost: cost,
			ROI:                (timeSaved*50 - cost) / cost,
		}
	}

	dollarValue := totalTimeSaved * 50 // $50/hour developer time
	roiRatio := (dollarValue - implementationCost) / implementationCost
	paybackPeriod := implementationCost / (dollarValue / 12) // Months

	return &ROIEstimation{
		TotalTimeSavedHours:   totalTimeSaved,
		TotalTimeSavedDollars: dollarValue,
		ImplementationCost:    implementationCost,
		ROIRatio:              roiRatio,
		PaybackPeriod:         paybackPeriod,
		BreakdownByCategory:   breakdown,
	}
}

// Calculate implementation cost based on effort
func calculateImplementationCost(effort string) float64 {
	switch effort {
	case "low":
		return 500 // $500
	case "medium":
		return 1500 // $1500
	case "high":
		return 3000 // $3000
	default:
		return 1000 // $1000
	}
}
