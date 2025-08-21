package analytics

import "time"

// InsightsReport represents comprehensive repository analytics
type InsightsReport struct {
	Owner        string    `json:"owner"`
	Repo         string    `json:"repo"`
	GeneratedAt  time.Time `json:"generated_at"`
	AnalysisDays int       `json:"analysis_days"`

	// Development Patterns
	DevPatterns *DevelopmentPatterns `json:"development_patterns"`

	// Code Intelligence
	CodeIntel *CodeIntelligence `json:"code_intelligence"`

	// Health Score (0-100)
	HealthScore *HealthScore `json:"health_score"`

	// Community Insights
	Community *CommunityInsights `json:"community_insights"`

	// Productivity Metrics
	Productivity *ProductivityMetrics `json:"productivity_metrics"`

	// Recommendations
	Recommendations []string `json:"recommendations"`
}

// DevelopmentPatterns analyzes commit patterns and development habits
type DevelopmentPatterns struct {
	CommitFrequency   *CommitFrequency  `json:"commit_frequency"`
	TimeDistribution  *TimeDistribution `json:"time_distribution"`
	CommitTypes       map[string]int    `json:"commit_types"`
	AverageCommitSize float64           `json:"average_commit_size"`
	BranchingStrategy string            `json:"branching_strategy"`
}

// CommitFrequency represents commit frequency analysis
type CommitFrequency struct {
	Daily        float64        `json:"daily_average"`
	Weekly       float64        `json:"weekly_average"`
	Monthly      float64        `json:"monthly_average"`
	PeakDays     []string       `json:"peak_days"`
	Trend        string         `json:"trend"` // "increasing", "decreasing", "stable"
	LastCommit   time.Time      `json:"last_commit"`
	CommitsByDay map[string]int `json:"commits_by_day"`
}

// TimeDistribution shows when development happens
type TimeDistribution struct {
	HourlyDistribution map[int]int    `json:"hourly_distribution"`
	WeekdayVsWeekend   map[string]int `json:"weekday_vs_weekend"`
	TimezonePattern    string         `json:"timezone_pattern"`
	WorkingHours       []int          `json:"working_hours"`
}

// CodeIntelligence provides code analysis insights
type CodeIntelligence struct {
	Languages       map[string]float64  `json:"languages"`
	PrimaryLanguage string              `json:"primary_language"`
	Complexity      *ComplexityMetrics  `json:"complexity"`
	Dependencies    *DependencyAnalysis `json:"dependencies"`
	FileTypes       map[string]int      `json:"file_types"`
	LinesOfCode     *LOCAnalysis        `json:"lines_of_code"`
}

// ComplexityMetrics analyzes code complexity
type ComplexityMetrics struct {
	CyclomaticComplexity float64 `json:"cyclomatic_complexity"`
	CodeDuplication      float64 `json:"code_duplication"`
	TechnicalDebt        string  `json:"technical_debt"`
	MaintainabilityIndex float64 `json:"maintainability_index"`
}

// DependencyAnalysis checks dependencies health
type DependencyAnalysis struct {
	TotalDependencies int      `json:"total_dependencies"`
	OutdatedCount     int      `json:"outdated_count"`
	VulnerableCount   int      `json:"vulnerable_count"`
	LicenseIssues     int      `json:"license_issues"`
	DependencyHealth  float64  `json:"dependency_health"`
	CriticalUpdates   []string `json:"critical_updates"`
}

// LOCAnalysis provides lines of code insights
type LOCAnalysis struct {
	Total      int            `json:"total"`
	Code       int            `json:"code"`
	Comments   int            `json:"comments"`
	Blank      int            `json:"blank"`
	ByLanguage map[string]int `json:"by_language"`
	GrowthRate float64        `json:"growth_rate"`
}

// HealthScore calculates overall repository health (0-100)
type HealthScore struct {
	Overall   float64            `json:"overall"`
	Breakdown map[string]float64 `json:"breakdown"`
	Grade     string             `json:"grade"` // "A+", "A", "B+", "B", "C+", "C", "D", "F"
	Factors   []string           `json:"factors"`
}

// CommunityInsights analyzes community engagement
type CommunityInsights struct {
	Contributors  *ContributorAnalysis  `json:"contributors"`
	Collaboration *CollaborationMetrics `json:"collaboration"`
	Growth        *GrowthMetrics        `json:"growth"`
	Diversity     *DiversityMetrics     `json:"diversity"`
}

// ContributorAnalysis analyzes contributor patterns
type ContributorAnalysis struct {
	Total           int               `json:"total"`
	Active          int               `json:"active"`
	TopContributors []ContributorInfo `json:"top_contributors"`
	NewContributors int               `json:"new_contributors"`
	RetentionRate   float64           `json:"retention_rate"`
}

// ContributorInfo represents individual contributor data
type ContributorInfo struct {
	Login       string    `json:"login"`
	Commits     int       `json:"commits"`
	Additions   int       `json:"additions"`
	Deletions   int       `json:"deletions"`
	FirstCommit time.Time `json:"first_commit"`
	LastCommit  time.Time `json:"last_commit"`
}

// CollaborationMetrics measures teamwork effectiveness
type CollaborationMetrics struct {
	PRReviewRate      float64 `json:"pr_review_rate"`
	AverageReviewTime float64 `json:"average_review_time"`
	IssueResponseTime float64 `json:"issue_response_time"`
	CrossTeamCommits  int     `json:"cross_team_commits"`
}

// GrowthMetrics tracks repository growth
type GrowthMetrics struct {
	Stars        *GrowthTrend `json:"stars"`
	Forks        *GrowthTrend `json:"forks"`
	Watchers     *GrowthTrend `json:"watchers"`
	Contributors *GrowthTrend `json:"contributors"`
}

// GrowthTrend represents growth analysis
type GrowthTrend struct {
	Current  int     `json:"current"`
	Previous int     `json:"previous"`
	Growth   float64 `json:"growth"`
	Trend    string  `json:"trend"`
}

// DiversityMetrics analyzes contribution diversity
type DiversityMetrics struct {
	TimezoneSpread      int     `json:"timezone_spread"`
	GeographicDiversity float64 `json:"geographic_diversity"`
	ContributionBalance float64 `json:"contribution_balance"`
}

// ProductivityMetrics measures development efficiency
type ProductivityMetrics struct {
	CodeChurn       float64 `json:"code_churn"`
	BugFixRate      float64 `json:"bug_fix_rate"`
	FeatureVelocity float64 `json:"feature_velocity"`
	DeploymentFreq  float64 `json:"deployment_frequency"`
	LeadTime        float64 `json:"lead_time"`
	MTTRecover      float64 `json:"mtt_recover"`
	ChangeFailRate  float64 `json:"change_fail_rate"`
	DevexScore      float64 `json:"devex_score"`
}
