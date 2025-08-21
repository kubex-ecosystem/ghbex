package productivity

import "time"

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
