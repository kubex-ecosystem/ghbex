package sanitize

import (
	"time"

	"github.com/google/go-github/v61/github"
)

// IntelligentSanitizer provides AI-powered repository cleanup and optimization
type IntelligentSanitizer struct {
	client *github.Client
}

// SanitizationReport contains intelligent cleanup analysis and actions
type SanitizationReport struct {
	Repository       string                `json:"repository"`
	Timestamp        time.Time             `json:"timestamp"`
	DryRun           bool                  `json:"dry_run"`
	OverallHealth    float64               `json:"overall_health"`
	ActionsPerformed []SanitizationAction  `json:"actions_performed"`
	Recommendations  []string              `json:"recommendations"`
	Savings          *ResourceSavings      `json:"savings"`
	SecurityImpacts  []SecurityImprovement `json:"security_impacts"`
	QualityImpacts   []QualityImprovement  `json:"quality_impacts"`
}

// SanitizationAction represents a cleanup action taken
type SanitizationAction struct {
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"`
	ItemsCount  int       `json:"items_count"`
	Savings     string    `json:"savings"`
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
}

// ResourceSavings quantifies cleanup benefits
type ResourceSavings struct {
	StorageMB        float64 `json:"storage_mb"`
	ComputeMinutes   int     `json:"compute_minutes"`
	SecurityRisk     string  `json:"security_risk_reduction"`
	MaintenanceHours float64 `json:"maintenance_hours_saved"`
}

// SecurityImprovement tracks security enhancements
type SecurityImprovement struct {
	Area        string `json:"area"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

// QualityImprovement tracks code quality enhancements
type QualityImprovement struct {
	Metric      string  `json:"metric"`
	Before      float64 `json:"before"`
	After       float64 `json:"after"`
	Improvement float64 `json:"improvement"`
}
