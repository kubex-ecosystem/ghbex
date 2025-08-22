package automation

import "time"

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
