package intelligence

import (
	"sync"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/kubex-ecosystem/ghbex/internal/defs/gromptz"
	"github.com/kubex-ecosystem/ghbex/internal/defs/interfaces"
)

type LLMMetaResponse struct {
	AIProvider string  `json:"ai_provider,omitempty"`
	AIModel    string  `json:"ai_model,omitempty"`
	AIEngine   string  `json:"ai_engine,omitempty"`
	AIType     string  `json:"ai_type,omitempty"`
	Attachment []byte  `json:"attachment,omitempty"`
	Response   string  `json:"response,omitempty"`
	Score      float64 `json:"score,omitempty"`
	Assessment string  `json:"assessment,omitempty"`
	Summary    string  `json:"summary,omitempty"`
	Status     string  `json:"status,omitempty"`
	Severity   string  `json:"severity,omitempty"`
	Suggestion string  `json:"suggestion,omitempty"`
	StatusCode int     `json:"status_code,omitempty"`
}

// IntelligenceOperator provides AI-powered analysis using Grompt engine
type IntelligenceOperator struct {
	client       *github.Client
	promptEngine gromptz.PromptEngine
	mainConfig   interfaces.IMainConfig

	// Health check cache para evitar verificações repetitivas
	healthCache      map[string]healthStatus
	healthCacheMutex sync.RWMutex
}

// healthStatus armazena o status de saúde de um provider com timestamp
type healthStatus struct {
	isHealthy bool
	lastCheck time.Time
}

// RepositoryInsight provides quick AI insights for repository cards
type RepositoryInsight struct {
	RepositoryName  string    `json:"repository_name" yaml:"repository_name"`
	AIScore         float64   `json:"ai_score" yaml:"ai_score"`
	QuickAssessment string    `json:"quick_assessment" yaml:"quick_assessment"`
	HealthIcon      string    `json:"health_icon" yaml:"health_icon"`
	MainTag         string    `json:"main_tag" yaml:"main_tag"`
	RiskLevel       string    `json:"risk_level" yaml:"risk_level"`
	Opportunity     string    `json:"opportunity" yaml:"opportunity"`
	LastAnalyzed    time.Time `json:"last_analyzed" yaml:"last_analyzed"`
}

// SmartRecommendation provides contextual recommendations
type SmartRecommendation struct {
	ID          string    `json:"id" yaml:"id"`
	Type        string    `json:"type" yaml:"type"` // "security", "performance", "maintenance", "enhancement"
	Title       string    `json:"title" yaml:"title"`
	Description string    `json:"description" yaml:"description"`
	Impact      string    `json:"impact" yaml:"impact"`
	Effort      string    `json:"effort" yaml:"effort"`
	Urgency     string    `json:"urgency" yaml:"urgency"`
	GeneratedAt time.Time `json:"generated_at" yaml:"generated_at"`
}

// HumanizedReport represents a comprehensive AI analysis
type HumanizedReport struct {
	RepositoryName    string                `json:"repository_name" yaml:"repository_name"`
	OverallAssessment OverallAssessment     `json:"overall_assessment" yaml:"overall_assessment"`
	KeyInsights       []KeyInsight          `json:"key_insights" yaml:"key_insights"`
	Recommendations   []SmartRecommendation `json:"recommendations" yaml:"recommendations"`
	ProductivityTips  []ProductivityTip     `json:"productivity_tips" yaml:"productivity_tips"`
	RiskFactors       []RiskFactor          `json:"risk_factors" yaml:"risk_factors"`
	NextSteps         []NextStep            `json:"next_steps" yaml:"next_steps"`
	GeneratedAt       time.Time             `json:"generated_at" yaml:"generated_at"`
	Metadata          map[string]any        `json:"metadata" yaml:"metadata"`
}

// OverallAssessment provides executive summary
type OverallAssessment struct {
	Grade         string   `json:"grade" yaml:"grade"`
	Score         float64  `json:"score" yaml:"score"`
	Summary       string   `json:"summary" yaml:"summary"`
	KeyStrengths  []string `json:"key_strengths" yaml:"key_strengths"`
	KeyWeaknesses []string `json:"key_weaknesses" yaml:"key_weaknesses"`
	Trend         string   `json:"trend" yaml:"trend"` // "improving", "stable", "declining"
}

// KeyInsight represents important findings
type KeyInsight struct {
	Category    string `json:"category" yaml:"category"`
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Impact      string `json:"impact" yaml:"impact"` // "high", "medium", "low"
	Evidence    string `json:"evidence" yaml:"evidence"`
}

// ProductivityTip provides actionable productivity advice
type ProductivityTip struct {
	Area       string `json:"area" yaml:"area"`
	Tip        string `json:"tip" yaml:"tip"`
	Benefit    string `json:"benefit" yaml:"benefit"`
	Difficulty string `json:"difficulty" yaml:"difficulty"`
	ROI        string `json:"roi" yaml:"roi"`
}

// RiskFactor identifies potential risks
type RiskFactor struct {
	Type        string `json:"type" yaml:"type"`
	Level       string `json:"level" yaml:"level"` // "critical", "high", "medium", "low"
	Description string `json:"description" yaml:"description"`
	Mitigation  string `json:"mitigation" yaml:"mitigation"`
	Probability string `json:"probability" yaml:"probability"`
}

// NextStep provides concrete actions
type NextStep struct {
	Order        int      `json:"order" yaml:"order"`
	Action       string   `json:"action" yaml:"action"`
	Owner        string   `json:"owner" yaml:"owner"`
	Timeline     string   `json:"timeline" yaml:"timeline"`
	Dependencies []string `json:"dependencies" yaml:"dependencies"`
}
