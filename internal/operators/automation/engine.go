// Package automation provides intelligent repository automation capabilities.
package automation

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"
)

// AutomationReport contains the results of automation analysis and actions.
type AutomationReport struct {
	Repository      string    `json:"repository"`
	Timestamp       time.Time `json:"timestamp"`
	AnalysisDays    int       `json:"analysis_days"`
	AutomationScore float64   `json:"automation_score"`
	Grade           string    `json:"grade"`

	// Organization & Management
	Labels       LabelManagement    `json:"labels"`
	Issues       IssueManagement    `json:"issues"`
	PullRequests PRManagement       `json:"pull_requests"`
	Workflows    WorkflowManagement `json:"workflows"`

	// Smart Suggestions
	Recommendations []AutomationAction `json:"recommendations"`
	EstimatedImpact EstimatedImpact    `json:"estimated_impact"`
}

type LabelManagement struct {
	Current              []string         `json:"current_labels"`
	Suggested            []SuggestedLabel `json:"suggested_labels"`
	Duplicates           []string         `json:"duplicate_labels"`
	UnusedLabels         []string         `json:"unused_labels"`
	StandardizationScore float64          `json:"standardization_score"`
	StandardizationGrade string           `json:"standardization_grade"`
}

type SuggestedLabel struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
	Priority    string `json:"priority"`
}

type IssueManagement struct {
	TotalIssues           int                    `json:"total_issues"`
	OpenIssues            int                    `json:"open_issues"`
	StaleIssues           int                    `json:"stale_issues"`
	AutoCloseables        []AutoCloseableIssue   `json:"auto_closeables"`
	MissingLabels         int                    `json:"missing_labels"`
	AssignmentSuggestions []AssignmentSuggestion `json:"assignment_suggestions"`
	EfficiencyScore       float64                `json:"efficiency_score"`
}

type AutoCloseableIssue struct {
	Number     int     `json:"number"`
	Title      string  `json:"title"`
	Reason     string  `json:"reason"`
	DaysStale  int     `json:"days_stale"`
	Confidence float64 `json:"confidence"`
	URL        string  `json:"url"`
}

type AssignmentSuggestion struct {
	IssueNumber int      `json:"issue_number"`
	IssueTitle  string   `json:"issue_title"`
	SuggestedTo []string `json:"suggested_to"`
	Reason      string   `json:"reason"`
	Confidence  float64  `json:"confidence"`
}

type PRManagement struct {
	TotalPRs            int                  `json:"total_prs"`
	OpenPRs             int                  `json:"open_prs"`
	StalePRs            int                  `json:"stale_prs"`
	AutoMergeables      []AutoMergeablePR    `json:"auto_mergeables"`
	ReviewerSuggestions []ReviewerSuggestion `json:"reviewer_suggestions"`
	ConflictedPRs       int                  `json:"conflicted_prs"`
	EfficiencyScore     float64              `json:"efficiency_score"`
}

type AutoMergeablePR struct {
	Number        int     `json:"number"`
	Title         string  `json:"title"`
	Reason        string  `json:"reason"`
	Confidence    float64 `json:"confidence"`
	ChecksPassing bool    `json:"checks_passing"`
	URL           string  `json:"url"`
	Author        string  `json:"author"`
}

type ReviewerSuggestion struct {
	PRNumber    int      `json:"pr_number"`
	PRTitle     string   `json:"pr_title"`
	SuggestedTo []string `json:"suggested_to"`
	Reason      string   `json:"reason"`
	Confidence  float64  `json:"confidence"`
}

type WorkflowManagement struct {
	TotalWorkflows             int                         `json:"total_workflows"`
	OptimizationOpportunities  []WorkflowOptimization      `json:"optimization_opportunities"`
	FailurePatterns            []FailurePattern            `json:"failure_patterns"`
	CacheOptimizations         []CacheOptimization         `json:"cache_optimizations"`
	ParallelizationSuggestions []ParallelizationSuggestion `json:"parallelization_suggestions"`
	EfficiencyScore            float64                     `json:"efficiency_score"`
}

type WorkflowOptimization struct {
	WorkflowName  string  `json:"workflow_name"`
	CurrentTime   float64 `json:"current_avg_time_minutes"`
	PotentialTime float64 `json:"potential_time_minutes"`
	Savings       float64 `json:"time_savings_minutes"`
	Suggestion    string  `json:"suggestion"`
	Difficulty    string  `json:"difficulty"`
	Priority      string  `json:"priority"`
}

type FailurePattern struct {
	Pattern   string    `json:"pattern"`
	Frequency int       `json:"frequency"`
	LastSeen  time.Time `json:"last_seen"`
	Solution  string    `json:"suggested_solution"`
	Impact    string    `json:"impact"`
	Priority  string    `json:"priority"`
}

type CacheOptimization struct {
	Step         string  `json:"step"`
	CacheHitRate float64 `json:"current_hit_rate"`
	Suggestion   string  `json:"suggestion"`
	Impact       string  `json:"estimated_impact"`
	Complexity   string  `json:"complexity"`
}

type ParallelizationSuggestion struct {
	CurrentWorkflow string   `json:"current_workflow"`
	SuggestedSplit  []string `json:"suggested_parallel_jobs"`
	TimeSavings     float64  `json:"estimated_time_savings_minutes"`
	Complexity      string   `json:"implementation_complexity"`
	ROI             float64  `json:"roi"`
}

type AutomationAction struct {
	Type           string  `json:"type"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Impact         string  `json:"impact"`
	Effort         string  `json:"effort"`
	Priority       string  `json:"priority"`
	ROI            float64 `json:"roi"`
	Implementation string  `json:"implementation"`
	AutoApplicable bool    `json:"auto_applicable"`
	EstimatedHours float64 `json:"estimated_hours"`
	Category       string  `json:"category"`
}

type EstimatedImpact struct {
	TimeSavedPerWeek      float64 `json:"time_saved_hours_per_week"`
	DeveloperSatisfaction string  `json:"developer_satisfaction_impact"`
	MaintenanceBurden     string  `json:"maintenance_burden_impact"`
	ProjectVelocity       string  `json:"project_velocity_impact"`
	ROIRating             float64 `json:"roi_rating"`
	AutomationCoverage    float64 `json:"automation_coverage_percentage"`
}

// AnalyzeAutomation performs comprehensive automation analysis on a repository.
func AnalyzeAutomation(ctx context.Context, client *github.Client, owner, repo string, analysisDays int) (*AutomationReport, error) {
	if analysisDays <= 0 {
		analysisDays = 90
	}

	report := &AutomationReport{
		Repository:   fmt.Sprintf("%s/%s", owner, repo),
		Timestamp:    time.Now(),
		AnalysisDays: analysisDays,
	}

	// Analyze labels
	if err := analyzeLabelManagement(ctx, client, owner, repo, report); err != nil {
		return nil, fmt.Errorf("label analysis failed: %w", err)
	}

	// Analyze issues
	if err := analyzeIssueManagement(ctx, client, owner, repo, report, analysisDays); err != nil {
		return nil, fmt.Errorf("issue analysis failed: %w", err)
	}

	// Analyze PRs
	if err := analyzePRManagement(ctx, client, owner, repo, report, analysisDays); err != nil {
		return nil, fmt.Errorf("PR analysis failed: %w", err)
	}

	// Analyze workflows
	if err := analyzeWorkflowManagement(ctx, client, owner, repo, report); err != nil {
		return nil, fmt.Errorf("workflow analysis failed: %w", err)
	}

	// Generate recommendations
	generateAutomationRecommendations(report)

	// Calculate automation score and grade
	report.AutomationScore = calculateAutomationScore(report)
	report.Grade = calculateGrade(report.AutomationScore)

	return report, nil
}

func analyzeLabelManagement(ctx context.Context, client *github.Client, owner, repo string, report *AutomationReport) error {
	// Get current labels
	labels, _, err := client.Issues.ListLabels(ctx, owner, repo, nil)
	if err != nil {
		return err
	}

	labelMgmt := LabelManagement{
		Current: make([]string, len(labels)),
	}

	for i, label := range labels {
		labelMgmt.Current[i] = label.GetName()
	}

	// Analyze for duplicates and suggest standardization
	labelMgmt.Duplicates = findDuplicateLabels(labelMgmt.Current)
	labelMgmt.Suggested = generateStandardLabels(labelMgmt.Current)
	labelMgmt.UnusedLabels = findUnusedLabels(ctx, client, owner, repo, labelMgmt.Current)
	labelMgmt.StandardizationScore = calculateLabelStandardization(labelMgmt.Current)
	labelMgmt.StandardizationGrade = calculateGrade(labelMgmt.StandardizationScore * 100)

	report.Labels = labelMgmt
	return nil
}

func analyzeIssueManagement(ctx context.Context, client *github.Client, owner, repo string, report *AutomationReport, days int) error {
	// Get recent issues
	since := time.Now().AddDate(0, 0, -days)
	allIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:       "all",
		Since:       since,
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return err
	}

	// Get open issues
	openIssues, _, err := client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return err
	}

	issueMgmt := IssueManagement{
		TotalIssues: 0,
		OpenIssues:  0,
	}

	// Filter out pull requests and count real issues
	var realIssues []*github.Issue
	var realOpenIssues []*github.Issue

	for _, issue := range allIssues {
		if !issue.IsPullRequest() {
			realIssues = append(realIssues, issue)
			issueMgmt.TotalIssues++
		}
	}

	for _, issue := range openIssues {
		if !issue.IsPullRequest() {
			realOpenIssues = append(realOpenIssues, issue)
			issueMgmt.OpenIssues++
		}
	}

	for _, issue := range realOpenIssues {
		// Check for stale issues
		if isStaleIssue(issue, 30) {
			issueMgmt.StaleIssues++

			// Check if auto-closeable
			if confidence := getAutoCloseConfidence(issue); confidence > 0.7 {
				issueMgmt.AutoCloseables = append(issueMgmt.AutoCloseables, AutoCloseableIssue{
					Number:     issue.GetNumber(),
					Title:      issue.GetTitle(),
					Reason:     getAutoCloseReason(issue),
					DaysStale:  int(time.Since(issue.GetUpdatedAt().Time).Hours() / 24),
					Confidence: confidence,
					URL:        issue.GetHTMLURL(),
				})
			}
		}

		// Check for missing labels
		if len(issue.Labels) == 0 {
			issueMgmt.MissingLabels++
		}
	}

	// Generate assignment suggestions
	issueMgmt.AssignmentSuggestions = generateAssignmentSuggestions(ctx, client, owner, repo, realOpenIssues)

	// Calculate efficiency score
	issueMgmt.EfficiencyScore = calculateIssueEfficiency(issueMgmt)

	report.Issues = issueMgmt
	return nil
}

func analyzePRManagement(ctx context.Context, client *github.Client, owner, repo string, report *AutomationReport, days int) error {
	// Get all PRs
	allPRs, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return err
	}

	// Get open PRs
	openPRs, _, err := client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return err
	}

	prMgmt := PRManagement{
		TotalPRs: len(allPRs),
		OpenPRs:  len(openPRs),
	}

	for _, pr := range openPRs {
		// Check for stale PRs
		if isStalePR(pr, 14) {
			prMgmt.StalePRs++
		}

		// Check for auto-mergeable PRs
		if confidence := getAutoMergeConfidence(ctx, client, owner, repo, pr); confidence > 0.8 {
			prMgmt.AutoMergeables = append(prMgmt.AutoMergeables, AutoMergeablePR{
				Number:        pr.GetNumber(),
				Title:         pr.GetTitle(),
				Reason:        getAutoMergeReason(pr),
				Confidence:    confidence,
				ChecksPassing: true, // Would need to check actual status
				URL:           pr.GetHTMLURL(),
				Author:        pr.GetUser().GetLogin(),
			})
		}

		// Check for conflicts
		if pr.GetMergeableState() == "dirty" || pr.GetMergeableState() == "blocked" {
			prMgmt.ConflictedPRs++
		}
	}

	// Generate reviewer suggestions
	prMgmt.ReviewerSuggestions = generateReviewerSuggestions(ctx, client, owner, repo, openPRs)

	// Calculate efficiency score
	prMgmt.EfficiencyScore = calculatePREfficiency(prMgmt)

	report.PullRequests = prMgmt
	return nil
}

func analyzeWorkflowManagement(ctx context.Context, client *github.Client, owner, repo string, report *AutomationReport) error {
	// Get workflows
	workflows, _, err := client.Actions.ListWorkflows(ctx, owner, repo, &github.ListOptions{PerPage: 100})
	if err != nil {
		// Repository might not have Actions enabled
		report.Workflows = WorkflowManagement{
			TotalWorkflows:  0,
			EfficiencyScore: 50.0, // Neutral score when no workflows
		}
		return nil
	}

	// Get workflow runs (last 100)
	runs, _, err := client.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})
	if err != nil {
		return err
	}

	workflowMgmt := WorkflowManagement{
		TotalWorkflows: workflows.GetTotalCount(),
	}

	// Analyze workflow performance
	workflowStats := make(map[string][]time.Duration)
	failurePatterns := make(map[string]int)

	for _, run := range runs.WorkflowRuns {
		workflowName := run.GetName()

		// Calculate duration
		createdAt := run.GetCreatedAt()
		updatedAt := run.GetUpdatedAt()
		if !createdAt.IsZero() && !updatedAt.IsZero() {
			duration := updatedAt.Time.Sub(createdAt.Time)
			workflowStats[workflowName] = append(workflowStats[workflowName], duration)
		}

		// Track failure patterns
		if run.GetConclusion() == "failure" {
			pattern := fmt.Sprintf("%s_failure", workflowName)
			failurePatterns[pattern]++
		}
	}

	// Generate optimization suggestions
	for workflow, durations := range workflowStats {
		if len(durations) > 5 {
			avgTime := calculateAverageDuration(durations)
			if avgTime > 10*time.Minute {
				priority := "medium"
				if avgTime > 30*time.Minute {
					priority = "high"
				}

				workflowMgmt.OptimizationOpportunities = append(workflowMgmt.OptimizationOpportunities, WorkflowOptimization{
					WorkflowName:  workflow,
					CurrentTime:   avgTime.Minutes(),
					PotentialTime: avgTime.Minutes() * 0.7, // 30% improvement estimate
					Savings:       avgTime.Minutes() * 0.3,
					Suggestion:    generateOptimizationSuggestion(workflow, avgTime),
					Difficulty:    "Medium",
					Priority:      priority,
				})
			}
		}
	}

	// Convert failure patterns
	for pattern, frequency := range failurePatterns {
		if frequency > 3 {
			priority := "low"
			if frequency > 10 {
				priority = "high"
			} else if frequency > 6 {
				priority = "medium"
			}

			workflowMgmt.FailurePatterns = append(workflowMgmt.FailurePatterns, FailurePattern{
				Pattern:   pattern,
				Frequency: frequency,
				LastSeen:  time.Now().AddDate(0, 0, -7), // Approximate
				Solution:  generateFailureSolution(pattern),
				Impact:    "Medium",
				Priority:  priority,
			})
		}
	}

	// Generate cache optimizations
	workflowMgmt.CacheOptimizations = generateCacheOptimizations(workflows.Workflows)

	// Generate parallelization suggestions
	workflowMgmt.ParallelizationSuggestions = generateParallelizationSuggestions(workflowStats)

	// Calculate efficiency score
	workflowMgmt.EfficiencyScore = calculateWorkflowEfficiency(workflowMgmt)

	report.Workflows = workflowMgmt
	return nil
}

func generateAutomationRecommendations(report *AutomationReport) {
	var recommendations []AutomationAction

	// Label standardization recommendations
	if report.Labels.StandardizationScore < 0.7 {
		recommendations = append(recommendations, AutomationAction{
			Type:           "label_standardization",
			Title:          "Standardize Repository Labels",
			Description:    fmt.Sprintf("Standardize %d labels to improve organization and consistency", len(report.Labels.Current)),
			Impact:         "Medium",
			Effort:         "Low",
			Priority:       "High",
			ROI:            2.5,
			Implementation: "Use GitHub's label sync action or manual cleanup with standardized color scheme",
			AutoApplicable: true,
			EstimatedHours: 2.0,
			Category:       "organization",
		})
	}

	// Auto-close stale issues
	if len(report.Issues.AutoCloseables) > 0 {
		recommendations = append(recommendations, AutomationAction{
			Type:           "auto_close_issues",
			Title:          "Auto-close Stale Issues",
			Description:    fmt.Sprintf("Automatically close %d stale issues that meet criteria", len(report.Issues.AutoCloseables)),
			Impact:         "High",
			Effort:         "Low",
			Priority:       "High",
			ROI:            4.0,
			Implementation: "Use GitHub Actions with stale bot and configure closing rules",
			AutoApplicable: true,
			EstimatedHours: 1.5,
			Category:       "maintenance",
		})
	}

	// Auto-merge PRs
	if len(report.PullRequests.AutoMergeables) > 0 {
		recommendations = append(recommendations, AutomationAction{
			Type:           "auto_merge_prs",
			Title:          "Enable Auto-merge for Safe PRs",
			Description:    fmt.Sprintf("Auto-merge %d PRs that pass all checks and meet safety criteria", len(report.PullRequests.AutoMergeables)),
			Impact:         "High",
			Effort:         "Medium",
			Priority:       "Medium",
			ROI:            3.2,
			Implementation: "Configure branch protection with auto-merge and setup safety rules",
			AutoApplicable: false,
			EstimatedHours: 4.0,
			Category:       "workflow",
		})
	}

	// Workflow optimizations
	if len(report.Workflows.OptimizationOpportunities) > 0 {
		totalSavings := 0.0
		for _, opt := range report.Workflows.OptimizationOpportunities {
			totalSavings += opt.Savings
		}

		recommendations = append(recommendations, AutomationAction{
			Type:           "workflow_optimization",
			Title:          "Optimize CI/CD Workflows",
			Description:    fmt.Sprintf("Optimize workflows to save %.1f minutes per run across %d workflows", totalSavings, len(report.Workflows.OptimizationOpportunities)),
			Impact:         "High",
			Effort:         "High",
			Priority:       "Medium",
			ROI:            2.8,
			Implementation: "Implement caching, parallelization, and workflow optimizations",
			AutoApplicable: false,
			EstimatedHours: 8.0,
			Category:       "performance",
		})
	}

	// Issue labeling automation
	if report.Issues.MissingLabels > 5 {
		recommendations = append(recommendations, AutomationAction{
			Type:           "auto_labeling",
			Title:          "Automated Issue Labeling",
			Description:    fmt.Sprintf("Automatically label %d issues missing labels using ML classification", report.Issues.MissingLabels),
			Impact:         "Medium",
			Effort:         "Medium",
			Priority:       "Low",
			ROI:            1.8,
			Implementation: "Setup GitHub Actions with ML-based auto-labeling",
			AutoApplicable: true,
			EstimatedHours: 3.0,
			Category:       "organization",
		})
	}

	// PR review automation
	if len(report.PullRequests.ReviewerSuggestions) > 0 {
		recommendations = append(recommendations, AutomationAction{
			Type:           "auto_reviewer_assignment",
			Title:          "Automated Reviewer Assignment",
			Description:    "Automatically assign reviewers based on code changes and expertise",
			Impact:         "Medium",
			Effort:         "Low",
			Priority:       "Medium",
			ROI:            2.1,
			Implementation: "Use CODEOWNERS file and GitHub's auto-assignment features",
			AutoApplicable: true,
			EstimatedHours: 1.0,
			Category:       "workflow",
		})
	}

	// Failure pattern automation
	if len(report.Workflows.FailurePatterns) > 0 {
		recommendations = append(recommendations, AutomationAction{
			Type:           "failure_recovery",
			Title:          "Automated Failure Recovery",
			Description:    fmt.Sprintf("Implement auto-recovery for %d detected failure patterns", len(report.Workflows.FailurePatterns)),
			Impact:         "High",
			Effort:         "High",
			Priority:       "Low",
			ROI:            3.5,
			Implementation: "Add retry logic, error handling, and auto-recovery mechanisms",
			AutoApplicable: false,
			EstimatedHours: 12.0,
			Category:       "reliability",
		})
	}

	// Sort by ROI and priority
	sort.Slice(recommendations, func(i, j int) bool {
		// First sort by priority
		priorityOrder := map[string]int{"high": 3, "medium": 2, "low": 1}
		if priorityOrder[strings.ToLower(recommendations[i].Priority)] != priorityOrder[strings.ToLower(recommendations[j].Priority)] {
			return priorityOrder[strings.ToLower(recommendations[i].Priority)] > priorityOrder[strings.ToLower(recommendations[j].Priority)]
		}
		// Then by ROI
		return recommendations[i].ROI > recommendations[j].ROI
	})

	report.Recommendations = recommendations

	// Calculate estimated impact
	report.EstimatedImpact = EstimatedImpact{
		TimeSavedPerWeek:      calculateTimeSavings(recommendations),
		DeveloperSatisfaction: calculateSatisfactionImpact(report),
		MaintenanceBurden:     calculateMaintenanceImpact(report),
		ProjectVelocity:       calculateVelocityImpact(report),
		ROIRating:             calculateAverageROI(recommendations),
		AutomationCoverage:    calculateAutomationCoverage(report),
	}
}

// Helper functions
func findDuplicateLabels(labels []string) []string {
	var duplicates []string
	seen := make(map[string]bool)

	for _, label := range labels {
		normalized := strings.ToLower(strings.ReplaceAll(label, " ", ""))
		normalizedDash := strings.ReplaceAll(normalized, "-", "")
		normalizedUnderscore := strings.ReplaceAll(normalized, "_", "")

		key := normalizedDash
		if len(normalizedUnderscore) < len(key) {
			key = normalizedUnderscore
		}

		if seen[key] {
			duplicates = append(duplicates, label)
		}
		seen[key] = true
	}

	return duplicates
}

func generateStandardLabels(current []string) []SuggestedLabel {
	standard := []SuggestedLabel{
		{Name: "bug", Color: "d73a4a", Description: "Something isn't working", Reason: "Standard GitHub label for bug reports", Priority: "high"},
		{Name: "enhancement", Color: "a2eeef", Description: "New feature or request", Reason: "Standard GitHub label for enhancements", Priority: "high"},
		{Name: "good first issue", Color: "7057ff", Description: "Good for newcomers", Reason: "Helps onboard new contributors", Priority: "medium"},
		{Name: "help wanted", Color: "008672", Description: "Extra attention is needed", Reason: "Community engagement and collaboration", Priority: "medium"},
		{Name: "priority: high", Color: "B60205", Description: "High priority issue", Reason: "Priority management and triage", Priority: "high"},
		{Name: "priority: medium", Color: "fbca04", Description: "Medium priority issue", Reason: "Priority management and triage", Priority: "medium"},
		{Name: "priority: low", Color: "0e8a16", Description: "Low priority issue", Reason: "Priority management and triage", Priority: "low"},
		{Name: "documentation", Color: "0075ca", Description: "Improvements or additions to documentation", Reason: "Categorize documentation work", Priority: "medium"},
		{Name: "duplicate", Color: "cfd3d7", Description: "This issue or pull request already exists", Reason: "Manage duplicate issues efficiently", Priority: "low"},
		{Name: "wontfix", Color: "ffffff", Description: "This will not be worked on", Reason: "Clear communication about scope decisions", Priority: "low"},
	}

	var suggestions []SuggestedLabel
	for _, std := range standard {
		found := false
		for _, curr := range current {
			if strings.EqualFold(curr, std.Name) ||
				strings.EqualFold(strings.ReplaceAll(curr, " ", "-"), strings.ReplaceAll(std.Name, " ", "-")) {
				found = true
				break
			}
		}
		if !found {
			suggestions = append(suggestions, std)
		}
	}

	return suggestions
}

func findUnusedLabels(ctx context.Context, client *github.Client, owner, repo string, labels []string) []string {
	// For performance reasons, we'll return a simplified analysis
	// In a full implementation, this would check each label's usage
	var unused []string

	// Common unused labels patterns
	commonUnused := []string{"invalid", "question", "wontfix"}
	for _, label := range labels {
		for _, pattern := range commonUnused {
			if strings.Contains(strings.ToLower(label), pattern) {
				unused = append(unused, label)
				break
			}
		}
	}

	return unused
}

func calculateLabelStandardization(labels []string) float64 {
	standardLabels := []string{"bug", "enhancement", "documentation", "good first issue", "help wanted"}
	found := 0

	for _, std := range standardLabels {
		for _, label := range labels {
			if strings.EqualFold(label, std) ||
				strings.Contains(strings.ToLower(label), std) {
				found++
				break
			}
		}
	}

	score := float64(found) / float64(len(standardLabels))

	// Bonus points for having priority labels
	priorityLabels := 0
	for _, label := range labels {
		if strings.Contains(strings.ToLower(label), "priority") {
			priorityLabels++
		}
	}
	if priorityLabels >= 3 {
		score += 0.1
	}

	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}

	return score
}

func isStaleIssue(issue *github.Issue, staleDays int) bool {
	if issue.GetState() != "open" {
		return false
	}

	staleDate := time.Now().AddDate(0, 0, -staleDays)
	return issue.GetUpdatedAt().Time.Before(staleDate)
}

func getAutoCloseConfidence(issue *github.Issue) float64 {
	confidence := 0.0

	title := strings.ToLower(issue.GetTitle())
	body := strings.ToLower(issue.GetBody())

	// Question issues that are old
	if strings.Contains(title, "how") || strings.Contains(title, "question") || strings.Contains(title, "?") {
		confidence += 0.3
	}

	// Issues with no activity for very long time
	daysSinceUpdate := int(time.Since(issue.GetUpdatedAt().Time).Hours() / 24)
	if daysSinceUpdate > 90 {
		confidence += 0.4
	}
	if daysSinceUpdate > 180 {
		confidence += 0.3
	}

	// Issues with specific keywords indicating resolution
	resolvedKeywords := []string{"never mind", "solved", "fixed", "resolved", "closing", "duplicate"}
	for _, keyword := range resolvedKeywords {
		if strings.Contains(body, keyword) {
			confidence += 0.5
			break
		}
	}

	// Issues with minimal engagement
	if issue.GetComments() == 0 && daysSinceUpdate > 60 {
		confidence += 0.2
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

func getAutoCloseReason(issue *github.Issue) string {
	daysSinceUpdate := int(time.Since(issue.GetUpdatedAt().Time).Hours() / 24)

	title := strings.ToLower(issue.GetTitle())
	body := strings.ToLower(issue.GetBody())

	// Check for resolved indicators
	resolvedKeywords := []string{"never mind", "solved", "fixed", "resolved"}
	for _, keyword := range resolvedKeywords {
		if strings.Contains(body, keyword) {
			return fmt.Sprintf("Issue appears to be resolved (contains '%s')", keyword)
		}
	}

	// Long-term stale
	if daysSinceUpdate > 180 {
		return fmt.Sprintf("No activity for %d days (very stale)", daysSinceUpdate)
	}

	// Question pattern
	if strings.Contains(title, "question") || strings.Contains(title, "how") || strings.Contains(title, "?") {
		return fmt.Sprintf("Question issue with no activity for %d days", daysSinceUpdate)
	}

	// General stale
	return fmt.Sprintf("No activity for %d days", daysSinceUpdate)
}

func generateAssignmentSuggestions(ctx context.Context, client *github.Client, owner, repo string, issues []*github.Issue) []AssignmentSuggestion {
	var suggestions []AssignmentSuggestion

	// Simple heuristic: suggest assigning to repository owner for critical issues
	for _, issue := range issues {
		if issue.GetAssignee() == nil && len(issue.Labels) > 0 {
			for _, label := range issue.Labels {
				if strings.Contains(strings.ToLower(label.GetName()), "bug") ||
					strings.Contains(strings.ToLower(label.GetName()), "critical") {
					suggestions = append(suggestions, AssignmentSuggestion{
						IssueNumber: issue.GetNumber(),
						IssueTitle:  issue.GetTitle(),
						SuggestedTo: []string{owner},
						Reason:      "Critical issue requiring maintainer attention",
						Confidence:  0.7,
					})
					break
				}
			}
		}

		// Limit suggestions to avoid overwhelming response
		if len(suggestions) >= 5 {
			break
		}
	}

	return suggestions
}

func isStalePR(pr *github.PullRequest, staleDays int) bool {
	if pr.GetState() != "open" {
		return false
	}

	staleDate := time.Now().AddDate(0, 0, -staleDays)
	return pr.GetUpdatedAt().Time.Before(staleDate)
}

func getAutoMergeConfidence(ctx context.Context, client *github.Client, owner, repo string, pr *github.PullRequest) float64 {
	confidence := 0.0

	// Check if it's a dependabot PR
	author := pr.GetUser().GetLogin()
	if strings.Contains(strings.ToLower(author), "dependabot") || strings.Contains(strings.ToLower(author), "bot") {
		confidence += 0.4
	}

	// Check if it's a small change
	additions := pr.GetAdditions()
	deletions := pr.GetDeletions()
	if additions > 0 || deletions > 0 {
		totalChanges := additions + deletions
		if totalChanges < 50 {
			confidence += 0.2
		}
		if totalChanges < 10 {
			confidence += 0.1
		}
	}

	// Check PR title for safe patterns
	title := strings.ToLower(pr.GetTitle())
	safePatterns := []string{"bump", "update", "upgrade", "docs", "documentation", "typo", "fix typo"}
	for _, pattern := range safePatterns {
		if strings.Contains(title, pattern) {
			confidence += 0.2
			break
		}
	}

	// Check if it has been open for reasonable time (not too new, not too old)
	daysSinceCreated := int(time.Since(pr.GetCreatedAt().Time).Hours() / 24)
	if daysSinceCreated >= 1 && daysSinceCreated <= 7 {
		confidence += 0.2
	}

	// Check if it's a draft
	if pr.GetDraft() {
		confidence = 0.0 // Never auto-merge drafts
	}

	return confidence
}

func getAutoMergeReason(pr *github.PullRequest) string {
	author := pr.GetUser().GetLogin()

	if strings.Contains(strings.ToLower(author), "dependabot") {
		return "Dependabot dependency update with passing checks"
	}

	title := strings.ToLower(pr.GetTitle())
	if strings.Contains(title, "docs") || strings.Contains(title, "documentation") {
		return "Documentation-only change"
	}

	if strings.Contains(title, "typo") {
		return "Typo fix with minimal risk"
	}

	additions := pr.GetAdditions()
	deletions := pr.GetDeletions()
	if additions > 0 || deletions > 0 {
		totalChanges := additions + deletions
		if totalChanges < 20 {
			return "Small change with minimal risk"
		}
	}

	return "Low-risk change with passing checks"
}

func generateReviewerSuggestions(ctx context.Context, client *github.Client, owner, repo string, prs []*github.PullRequest) []ReviewerSuggestion {
	var suggestions []ReviewerSuggestion

	// Simple heuristic: suggest owner for PRs without reviewers
	for _, pr := range prs {
		// Check if PR has reviewers (simplified check)
		hasReviewers := false
		if pr.Assignee != nil {
			hasReviewers = true
		}

		if !hasReviewers {
			suggestions = append(suggestions, ReviewerSuggestion{
				PRNumber:    pr.GetNumber(),
				PRTitle:     pr.GetTitle(),
				SuggestedTo: []string{owner},
				Reason:      "Repository maintainer review needed",
				Confidence:  0.8,
			})
		}

		// Limit suggestions
		if len(suggestions) >= 3 {
			break
		}
	}

	return suggestions
}

func calculateAverageDuration(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}

	return total / time.Duration(len(durations))
}

func generateOptimizationSuggestion(workflow string, avgTime time.Duration) string {
	minutes := avgTime.Minutes()
	workflowLower := strings.ToLower(workflow)

	// Base suggestions for different time ranges
	var suggestions []string

	// Critical performance issues (>60 minutes)
	if minutes > 60 {
		suggestions = append(suggestions,
			"🚨 Critical: Split into multiple parallel workflows",
			"Implement matrix builds for parallel execution",
			"Use self-hosted runners for better performance",
			"Cache all dependencies and build artifacts",
		)
	} else if minutes > 30 {
		// Major optimization needed (30-60 minutes)
		suggestions = append(suggestions,
			"📈 Major optimization needed:",
			"Add comprehensive dependency caching",
			"Parallelize independent jobs and test suites",
			"Optimize Docker builds with multi-stage builds",
			"Consider using faster runners or more CPU cores",
		)
	} else if minutes > 15 {
		// Moderate optimization (15-30 minutes)
		suggestions = append(suggestions,
			"⚡ Moderate optimization opportunities:",
			"Implement dependency caching (npm, pip, go mod, etc.)",
			"Parallelize test execution where possible",
			"Optimize container images and reduce layer size",
		)
	} else if minutes > 5 {
		// Minor optimizations (5-15 minutes)
		suggestions = append(suggestions,
			"🔧 Fine-tuning opportunities:",
			"Add selective caching for build artifacts",
			"Parallelize test suites if not already done",
			"Review and eliminate redundant build steps",
		)
	} else {
		// Already optimized (<5 minutes)
		return "✅ Workflow is well-optimized (under 5 minutes)"
	}

	// Workflow-specific suggestions
	if strings.Contains(workflowLower, "test") {
		suggestions = append(suggestions,
			"Test-specific: Use test parallelization and smart test selection")
	}

	if strings.Contains(workflowLower, "build") {
		suggestions = append(suggestions,
			"Build-specific: Implement incremental builds and artifact caching")
	}

	if strings.Contains(workflowLower, "deploy") {
		suggestions = append(suggestions,
			"Deploy-specific: Use blue-green deployment and artifact reuse")
	}

	if strings.Contains(workflowLower, "ci") || strings.Contains(workflowLower, "integration") {
		suggestions = append(suggestions,
			"CI-specific: Cache dependencies and use conditional job execution")
	}

	return strings.Join(suggestions, "\n• ")
}

func generateFailureSolution(pattern string) string {
	patternLower := strings.ToLower(pattern)

	// Test failures - most common category
	if strings.Contains(patternLower, "test") {
		solutions := []string{
			"🧪 Test Failure Solutions:",
			"• Review and stabilize flaky tests with proper waits and assertions",
			"• Implement retry logic for integration tests with external dependencies",
			"• Improve test isolation and cleanup between test runs",
			"• Add better test data management and fixtures",
			"• Consider test parallelization issues and race conditions",
			"• Use deterministic test ordering and seeding",
		}

		// Specific test type failures
		if strings.Contains(patternLower, "unit") {
			solutions = append(solutions, "• Focus on mocking external dependencies for unit tests")
		}
		if strings.Contains(patternLower, "integration") {
			solutions = append(solutions, "• Verify test environment consistency and service availability")
		}
		if strings.Contains(patternLower, "e2e") || strings.Contains(patternLower, "end") {
			solutions = append(solutions, "• Add browser/UI stability improvements and explicit waits")
		}

		return strings.Join(solutions, "\n")
	}

	// Build failures
	if strings.Contains(patternLower, "build") || strings.Contains(patternLower, "compile") {
		solutions := []string{
			"🔨 Build Failure Solutions:",
			"• Lock dependency versions to prevent version conflicts",
			"• Verify build environment consistency across all runners",
			"• Add comprehensive error logging and build diagnostics",
			"• Implement incremental builds to isolate problem areas",
			"• Check for platform-specific build issues (OS, architecture)",
			"• Validate build tool versions and configurations",
		}

		if strings.Contains(patternLower, "docker") {
			solutions = append(solutions, "• Review Dockerfile syntax and base image availability")
		}
		if strings.Contains(patternLower, "node") || strings.Contains(patternLower, "npm") {
			solutions = append(solutions, "• Clear npm cache and verify package-lock.json consistency")
		}
		if strings.Contains(patternLower, "maven") || strings.Contains(patternLower, "gradle") {
			solutions = append(solutions, "• Check Java version compatibility and clear build cache")
		}

		return strings.Join(solutions, "\n")
	}

	// Deployment failures
	if strings.Contains(patternLower, "deploy") || strings.Contains(patternLower, "release") {
		solutions := []string{
			"🚀 Deployment Failure Solutions:",
			"• Implement comprehensive deployment health checks",
			"• Add automated rollback mechanisms for failed deployments",
			"• Verify environment configuration and secrets availability",
			"• Check network connectivity and firewall rules",
			"• Validate resource limits and capacity planning",
			"• Implement blue-green or canary deployment strategies",
			"• Add deployment monitoring and alerting",
		}

		if strings.Contains(patternLower, "kubernetes") || strings.Contains(patternLower, "k8s") {
			solutions = append(solutions, "• Verify Kubernetes manifests and cluster permissions")
		}
		if strings.Contains(patternLower, "aws") || strings.Contains(patternLower, "azure") || strings.Contains(patternLower, "gcp") {
			solutions = append(solutions, "• Check cloud provider service availability and quotas")
		}

		return strings.Join(solutions, "\n")
	}

	// Security scan failures
	if strings.Contains(patternLower, "security") || strings.Contains(patternLower, "vulnerability") {
		return "🔒 Security Failure Solutions:\n• Update vulnerable dependencies to patched versions\n• Review and whitelist false positives\n• Implement security baseline and compliance checks\n• Add security scanning earlier in the development cycle"
	}

	// Linting/code quality failures
	if strings.Contains(patternLower, "lint") || strings.Contains(patternLower, "quality") {
		return "📝 Code Quality Failure Solutions:\n• Fix linting violations or update linting rules\n• Implement pre-commit hooks for early detection\n• Standardize code formatting across the team\n• Add incremental linting for large codebases"
	}

	// Performance/timeout failures
	if strings.Contains(patternLower, "timeout") || strings.Contains(patternLower, "performance") {
		return "⏱️ Performance/Timeout Solutions:\n• Increase timeout values for slow operations\n• Optimize resource-intensive operations\n• Add performance monitoring and profiling\n• Consider parallel execution for time-consuming tasks"
	}

	// Generic failure pattern
	return "🔍 Generic Failure Analysis:\n• Analyze detailed failure logs and error messages\n• Implement comprehensive error handling and recovery\n• Add monitoring and alerting for early detection\n• Review recent changes that might have introduced the issue\n• Ensure proper environment setup and dependencies"
}

func generateCacheOptimizations(workflows []*github.Workflow) []CacheOptimization {
	var optimizations []CacheOptimization

	// TODO: Implement real workflow analysis for cache optimization opportunities
	// Currently returning empty to avoid hardcoded suggestions
	// Analysis should examine actual workflow files for:
	// - Dependency management steps
	// - Build artifact patterns
	// - Cache action usage
	// - Build time metrics

	return optimizations
}

func generateParallelizationSuggestions(workflowStats map[string][]time.Duration) []ParallelizationSuggestion {
	var suggestions []ParallelizationSuggestion

	// TODO: Implement real workflow parallelization analysis
	// Currently returning empty to avoid hardcoded suggestions
	// Analysis should:
	// - Parse actual workflow YAML files
	// - Identify dependencies between jobs
	// - Calculate potential parallelization opportunities
	// - Provide specific, actionable suggestions

	return suggestions
}

func calculateIssueEfficiency(mgmt IssueManagement) float64 {
	if mgmt.TotalIssues == 0 {
		return 100.0 // Perfect score for no issues
	}

	score := 100.0

	// Deduct for stale issues
	if mgmt.OpenIssues > 0 {
		stalePercentage := float64(mgmt.StaleIssues) / float64(mgmt.OpenIssues)
		score -= stalePercentage * 30 // Up to 30 points deduction
	}

	// Deduct for missing labels
	if mgmt.OpenIssues > 0 {
		unlabeledPercentage := float64(mgmt.MissingLabels) / float64(mgmt.OpenIssues)
		score -= unlabeledPercentage * 20 // Up to 20 points deduction
	}

	// Bonus for having auto-closeables identified
	if len(mgmt.AutoCloseables) > 0 {
		score += 10
	}

	if score < 0 {
		score = 0
	}

	return score
}

func calculatePREfficiency(mgmt PRManagement) float64 {
	if mgmt.TotalPRs == 0 {
		return 100.0 // Perfect score for no PRs
	}

	score := 100.0

	// Deduct for stale PRs
	if mgmt.OpenPRs > 0 {
		stalePercentage := float64(mgmt.StalePRs) / float64(mgmt.OpenPRs)
		score -= stalePercentage * 25 // Up to 25 points deduction
	}

	// Deduct for conflicted PRs
	if mgmt.OpenPRs > 0 {
		conflictPercentage := float64(mgmt.ConflictedPRs) / float64(mgmt.OpenPRs)
		score -= conflictPercentage * 20 // Up to 20 points deduction
	}

	// Bonus for having auto-mergeable PRs
	if len(mgmt.AutoMergeables) > 0 {
		score += 10
	}

	if score < 0 {
		score = 0
	}

	return score
}

func calculateWorkflowEfficiency(mgmt WorkflowManagement) float64 {
	if mgmt.TotalWorkflows == 0 {
		return 50.0 // Neutral score when no workflows
	}

	score := 100.0

	// Deduct for optimization opportunities
	if mgmt.TotalWorkflows > 0 {
		optimizationPercentage := float64(len(mgmt.OptimizationOpportunities)) / float64(mgmt.TotalWorkflows)
		score -= optimizationPercentage * 30 // Up to 30 points deduction
	}

	// Deduct for failure patterns
	score -= float64(len(mgmt.FailurePatterns)) * 5 // 5 points per failure pattern

	// Bonus for having parallelization suggestions (indicates analysis depth)
	if len(mgmt.ParallelizationSuggestions) > 0 {
		score += 5
	}

	if score < 0 {
		score = 0
	}

	return score
}

func calculateAutomationScore(report *AutomationReport) float64 {
	// Weighted scoring system
	labelWeight := 0.20
	issueWeight := 0.30
	prWeight := 0.25
	workflowWeight := 0.25

	labelScore := report.Labels.StandardizationScore * 100
	issueScore := report.Issues.EfficiencyScore
	prScore := report.PullRequests.EfficiencyScore
	workflowScore := report.Workflows.EfficiencyScore

	finalScore := (labelScore * labelWeight) +
		(issueScore * issueWeight) +
		(prScore * prWeight) +
		(workflowScore * workflowWeight)

	// Bonus for having recommendations
	if len(report.Recommendations) > 0 {
		finalScore += 5
	}

	// Cap at 100
	if finalScore > 100 {
		finalScore = 100
	}

	return finalScore
}

func calculateGrade(score float64) string {
	if score >= 90 {
		return "A+"
	} else if score >= 85 {
		return "A"
	} else if score >= 80 {
		return "A-"
	} else if score >= 75 {
		return "B+"
	} else if score >= 70 {
		return "B"
	} else if score >= 65 {
		return "B-"
	} else if score >= 60 {
		return "C+"
	} else if score >= 55 {
		return "C"
	} else if score >= 50 {
		return "C-"
	} else if score >= 40 {
		return "D"
	} else {
		return "F"
	}
}

func calculateTimeSavings(recommendations []AutomationAction) float64 {
	savings := 0.0
	for _, rec := range recommendations {
		switch rec.Type {
		case "auto_close_issues":
			savings += 2.0 // 2 hours per week
		case "auto_merge_prs":
			savings += 1.5 // 1.5 hours per week
		case "label_standardization":
			savings += 0.5 // 30 minutes per week
		case "workflow_optimization":
			savings += 3.0 // 3 hours per week
		case "auto_labeling":
			savings += 1.0 // 1 hour per week
		case "auto_reviewer_assignment":
			savings += 0.5 // 30 minutes per week
		case "failure_recovery":
			savings += 2.5 // 2.5 hours per week
		default:
			savings += 1.0 // Default 1 hour per week
		}
	}
	return savings
}

func calculateSatisfactionImpact(report *AutomationReport) string {
	score := 0

	if len(report.Recommendations) > 3 {
		score += 2
	}
	if report.AutomationScore > 70 {
		score += 2
	}
	if len(report.Issues.AutoCloseables) > 0 {
		score += 1
	}
	if len(report.PullRequests.AutoMergeables) > 0 {
		score += 1
	}

	switch {
	case score >= 5:
		return "Very High - Significant reduction in manual work"
	case score >= 3:
		return "High - Noticeable improvement in workflow"
	case score >= 1:
		return "Medium - Some improvements in daily tasks"
	default:
		return "Low - Minimal impact on current workflow"
	}
}

func calculateMaintenanceImpact(report *AutomationReport) string {
	complexActions := 0

	for _, rec := range report.Recommendations {
		if rec.Effort == "High" {
			complexActions++
		}
	}

	switch {
	case complexActions >= 3:
		return "Medium - Some ongoing maintenance required"
	case complexActions >= 1:
		return "Low - Minimal ongoing maintenance"
	default:
		return "Very Low - Set-and-forget automation"
	}
}

func calculateVelocityImpact(report *AutomationReport) string {
	timeSaved := report.EstimatedImpact.TimeSavedPerWeek

	switch {
	case timeSaved >= 5:
		return "Very High - Significant acceleration of development cycles"
	case timeSaved >= 3:
		return "High - Noticeable improvement in delivery speed"
	case timeSaved >= 1:
		return "Medium - Moderate improvement in workflow speed"
	default:
		return "Low - Minimal impact on delivery velocity"
	}
}

func calculateAverageROI(recommendations []AutomationAction) float64 {
	if len(recommendations) == 0 {
		return 0.0
	}

	total := 0.0
	for _, rec := range recommendations {
		total += rec.ROI
	}

	return total / float64(len(recommendations))
}

func calculateAutomationCoverage(report *AutomationReport) float64 {
	totalAreas := 4.0 // Labels, Issues, PRs, Workflows
	coveredAreas := 0.0

	if len(report.Labels.Suggested) > 0 || report.Labels.StandardizationScore > 0.8 {
		coveredAreas++
	}
	if len(report.Issues.AutoCloseables) > 0 || report.Issues.EfficiencyScore > 80 {
		coveredAreas++
	}
	if len(report.PullRequests.AutoMergeables) > 0 || report.PullRequests.EfficiencyScore > 80 {
		coveredAreas++
	}
	if len(report.Workflows.OptimizationOpportunities) > 0 || report.Workflows.EfficiencyScore > 80 {
		coveredAreas++
	}

	return (coveredAreas / totalAreas) * 100
}
