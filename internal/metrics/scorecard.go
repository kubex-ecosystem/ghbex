package metrics

import (
	"time"
)

// Enhanced Scorecard types that integrate with current GHbex structure

// EnhancedScorecard representa um scorecard completo seguindo o schema v1
type EnhancedScorecard struct {
	SchemaVersion string    `json:"schema_version" yaml:"schema_version"`
	Owner         string    `json:"owner" yaml:"owner"`
	Repo          string    `json:"repo" yaml:"repo"`
	GeneratedAt   time.Time `json:"generated_at" yaml:"generated_at"`
	PeriodDays    int       `json:"period_days" yaml:"period_days"`

	// Windows for trend calculations
	Windows Windows `json:"windows" yaml:"windows"`

	// Confidence levels for each metric group
	Confidence Confidence `json:"confidence" yaml:"confidence"`

	// Main metric groups
	Health    HealthMetrics    `json:"health" yaml:"health"`
	Dora      DoraMetrics      `json:"dora" yaml:"dora"`
	Code      CodeMetrics      `json:"code" yaml:"code"`
	Deps      DepsMetrics      `json:"deps" yaml:"deps"`
	Community CommunityMetrics `json:"community" yaml:"community"`

	// Provenance for auditability
	Provenance *Provenance `json:"provenance,omitempty" yaml:"provenance,omitempty"`
}

type Windows struct {
	TrendDays int `json:"trend_days" yaml:"trend_days"`
}

type Confidence struct {
	Dora      float64 `json:"dora" yaml:"dora"`
	Code      float64 `json:"code" yaml:"code"`
	Deps      float64 `json:"deps" yaml:"deps"`
	Community float64 `json:"community" yaml:"community"`
}

type HealthMetrics struct {
	CHI   float64 `json:"chi" yaml:"chi"`
	Grade string  `json:"grade" yaml:"grade"`
}

type CodeMetrics struct {
	PrimaryLanguage string             `json:"primary_language" yaml:"primary_language"`
	MI              float64            `json:"mi" yaml:"mi"`
	DuplicationPct  float64            `json:"duplication_pct" yaml:"duplication_pct"`
	CyclomaticAvg   float64            `json:"cyclomatic_avg" yaml:"cyclomatic_avg"`
	LocTotal        int                `json:"loc_total,omitempty" yaml:"loc_total,omitempty"`
	LocCode         int                `json:"loc_code,omitempty" yaml:"loc_code,omitempty"`
	LanguagesPct    map[string]float64 `json:"languages_pct,omitempty" yaml:"languages_pct,omitempty"`
	Trend           []float64          `json:"trend,omitempty" yaml:"trend,omitempty"`
}

type DepsMetrics struct {
	Count      int     `json:"count" yaml:"count"`
	Vulnerable int     `json:"vulnerable" yaml:"vulnerable"`
	Outdated   int     `json:"outdated" yaml:"outdated"`
	Health     float64 `json:"health" yaml:"health"`
}

type CommunityMetrics struct {
	Contributors      int       `json:"contributors" yaml:"contributors"`
	BusFactor         int       `json:"bus_factor" yaml:"bus_factor"`
	FirstReviewP50    float64   `json:"first_review_p50,omitempty" yaml:"first_review_p50,omitempty"`
	FirstReviewP90    float64   `json:"first_review_p90,omitempty" yaml:"first_review_p90,omitempty"`
	ReviewCoveragePct float64   `json:"review_coverage_pct,omitempty" yaml:"review_coverage_pct,omitempty"`
	OnboardingDaysP50 float64   `json:"onboarding_days_p50,omitempty" yaml:"onboarding_days_p50,omitempty"`
	TrendLeadTime     []float64 `json:"trend_lead_time,omitempty" yaml:"trend_lead_time,omitempty"`
}

type Provenance struct {
	Sources []DataSource `json:"sources" yaml:"sources"`
	Notes   string       `json:"notes,omitempty" yaml:"notes,omitempty"`
}

type DataSource struct {
	Type      string    `json:"type" yaml:"type"`         // "github", "llm", "static_analysis"
	Provider  string    `json:"provider" yaml:"provider"` // "github.com", "openai", "cloc"
	Version   string    `json:"version,omitempty" yaml:"version,omitempty"`
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
}

// ConvertFromGHbexModel converte do modelo atual do GHbex para EnhancedScorecard
func ConvertFromGHbexModel(data map[string]interface{}, owner, repo string) *EnhancedScorecard {
	scorecard := &EnhancedScorecard{
		SchemaVersion: "repo_scorecard@1.0.0",
		Owner:         owner,
		Repo:          repo,
		GeneratedAt:   time.Now(),
		PeriodDays:    60, // padrão
		Windows: Windows{
			TrendDays: 60,
		},
		Confidence: Confidence{
			Dora:      0.7,
			Code:      0.6,
			Deps:      0.9,
			Community: 0.5,
		},
	}

	// Extract health score
	if health, ok := data["health_score"].(map[string]interface{}); ok {
		if overall, ok := health["overall"].(float64); ok {
			scorecard.Health.CHI = overall
			scorecard.Health.Grade = GradeFromCHI(overall)
		}
	}

	// Extract code intelligence
	if codeIntel, ok := data["code_intelligence"].(map[string]interface{}); ok {
		if primaryLang, ok := codeIntel["primary_language"].(string); ok {
			scorecard.Code.PrimaryLanguage = primaryLang
		}

		if complexity, ok := codeIntel["complexity"].(map[string]interface{}); ok {
			if mi, ok := complexity["maintainability_index"].(float64); ok {
				scorecard.Code.MI = mi
			}
			if dup, ok := complexity["code_duplication"].(float64); ok {
				scorecard.Code.DuplicationPct = dup
			}
			if cyclo, ok := complexity["cyclomatic_complexity"].(float64); ok {
				scorecard.Code.CyclomaticAvg = cyclo
			}
		}

		if deps, ok := codeIntel["dependencies"].(map[string]interface{}); ok {
			if total, ok := deps["total_dependencies"].(float64); ok {
				scorecard.Deps.Count = int(total)
			}
			if vuln, ok := deps["vulnerable_count"].(float64); ok {
				scorecard.Deps.Vulnerable = int(vuln)
			}
			if outdated, ok := deps["outdated_count"].(float64); ok {
				scorecard.Deps.Outdated = int(outdated)
			}
			if health, ok := deps["dependency_health"].(float64); ok {
				scorecard.Deps.Health = health
			}
		}

		if languages, ok := codeIntel["languages"].(map[string]interface{}); ok {
			scorecard.Code.LanguagesPct = make(map[string]float64)
			for lang, pct := range languages {
				if pctFloat, ok := pct.(float64); ok {
					scorecard.Code.LanguagesPct[lang] = pctFloat
				}
			}
		}
	}

	// Extract community insights
	if community, ok := data["community_insights"].(map[string]interface{}); ok {
		if contributors, ok := community["contributors"].(map[string]interface{}); ok {
			if total, ok := contributors["total"].(float64); ok {
				scorecard.Community.Contributors = int(total)
			}
		}

		if collaboration, ok := community["collaboration"].(map[string]interface{}); ok {
			if reviewTime, ok := collaboration["average_review_time"].(float64); ok {
				scorecard.Community.FirstReviewP50 = reviewTime
			}
		}
	}

	// Extract productivity metrics for DORA
	if productivity, ok := data["productivity_metrics"].(map[string]interface{}); ok {
		if deployFreq, ok := productivity["deployment_frequency"].(float64); ok {
			scorecard.Dora.DeploymentFrequency = deployFreq
		}
		if leadTime, ok := productivity["lead_time"].(float64); ok {
			scorecard.Dora.LeadTimeP95 = leadTime
		}
		if changeFailRate, ok := productivity["change_fail_rate"].(float64); ok {
			scorecard.Dora.ChangeFailRate = changeFailRate
		}
		if mttr, ok := productivity["mtt_recover"].(float64); ok {
			scorecard.Dora.MTTR = mttr
		}
		scorecard.Dora.PeriodUnit = "week"
	}

	// Add provenance
	scorecard.Provenance = &Provenance{
		Sources: []DataSource{
			{
				Type:      "github",
				Provider:  "github.com",
				Timestamp: time.Now(),
			},
		},
		Notes: "Converted from GHbex analysis model",
	}

	return scorecard
}

// CalculateBusFactor estima o bus factor baseado na distribuição de commits
func CalculateBusFactor(contributors []map[string]interface{}) int {
	if len(contributors) == 0 {
		return 1
	}

	// Simplified bus factor calculation
	// In reality, you'd analyze commit distribution more carefully
	switch {
	case len(contributors) >= 5:
		return min(len(contributors)/2, 5)
	case len(contributors) >= 3:
		return 2
	default:
		return 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
