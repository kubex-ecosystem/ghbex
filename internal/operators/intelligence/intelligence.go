// Package intelligence provides AI-powered analysis and insights for GitHub repositories.
package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/operators/analytics"
	"github.com/rafa-mori/ghbex/internal/operators/productivity"
	"github.com/rafa-mori/grompt"
)

// Operator provides AI-powered analysis using Grompt engine
type Operator struct {
	client       *github.Client
	promptEngine grompt.PromptEngine
	analytics    *analytics.LOCAnalysis
	productivity *productivity.ProductivityReport
}

// HumanizedReport represents an AI-processed, human-readable analysis
type HumanizedReport struct {
	RepositoryName    string                 `json:"repository_name"`
	OverallAssessment OverallAssessment      `json:"overall_assessment"`
	KeyInsights       []KeyInsight           `json:"key_insights"`
	Recommendations   []Recommendation       `json:"recommendations"`
	HealthSummary     HealthSummary          `json:"health_summary"`
	ProductivityTips  []ProductivityTip      `json:"productivity_tips"`
	RiskFactors       []RiskFactor           `json:"risk_factors"`
	NextSteps         []NextStep             `json:"next_steps"`
	GeneratedAt       time.Time              `json:"generated_at"`
	Metadata          map[string]interface{} `json:"metadata"`
}

// OverallAssessment provides executive summary
type OverallAssessment struct {
	Grade         string   `json:"grade"`
	Score         float64  `json:"score"`
	Summary       string   `json:"summary"`
	KeyStrengths  []string `json:"key_strengths"`
	KeyWeaknesses []string `json:"key_weaknesses"`
	Trend         string   `json:"trend"` // "improving", "stable", "declining"
}

// KeyInsight represents important findings
type KeyInsight struct {
	Category    string `json:"category"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Impact      string `json:"impact"` // "high", "medium", "low"
	Evidence    string `json:"evidence"`
}

// Recommendation provides actionable advice
type Recommendation struct {
	Priority    string   `json:"priority"` // "critical", "high", "medium", "low"
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Effort      string   `json:"effort"`   // "low", "medium", "high"
	Impact      string   `json:"impact"`   // "low", "medium", "high"
	Timeline    string   `json:"timeline"` // "immediate", "short", "medium", "long"
	Resources   []string `json:"resources"`
}

// HealthSummary provides health breakdown
type HealthSummary struct {
	OverallHealth    float64            `json:"overall_health"`
	CategoryScores   map[string]float64 `json:"category_scores"`
	HealthTrend      string             `json:"health_trend"`
	CriticalIssues   []string           `json:"critical_issues"`
	ImprovementAreas []string           `json:"improvement_areas"`
}

// ProductivityTip provides actionable productivity advice
type ProductivityTip struct {
	Area       string `json:"area"`
	Tip        string `json:"tip"`
	Benefit    string `json:"benefit"`
	Difficulty string `json:"difficulty"`
	ROI        string `json:"roi"`
}

// RiskFactor identifies potential risks
type RiskFactor struct {
	Type        string `json:"type"`
	Level       string `json:"level"` // "critical", "high", "medium", "low"
	Description string `json:"description"`
	Mitigation  string `json:"mitigation"`
	Probability string `json:"probability"`
}

// NextStep provides concrete actions
type NextStep struct {
	Order        int      `json:"order"`
	Action       string   `json:"action"`
	Owner        string   `json:"owner"`
	Timeline     string   `json:"timeline"`
	Dependencies []string `json:"dependencies"`
}

// NewOperator creates a new Intelligence operator
func NewOperator(client *github.Client, analyticsOp *analytics.CodeIntelligence, productivityOp *productivity.ProductivityReport) *Operator {
	// Initialize Grompt with basic config
	config := grompt.DefaultConfig()

	llmMapList := map[string]grompt.APIConfig{
		"openai":   config.GetAPIConfig("openai"),
		"claude":   config.GetAPIConfig("claude"),
		"gemini":   config.GetAPIConfig("gemini"),
		"chatgpt":  config.GetAPIConfig("chatgpt"),
		"deepseek": config.GetAPIConfig("deepseek"),
		"ollama":   config.GetAPIConfig("ollama"),
	}

	for name, apiConfig := range llmMapList {
		if apiConfig != nil {
			if apiConfig.IsAvailable() {
				log.Printf("INTELLIGENCE: Using %s API for AI processing", name)
				if err := config.SetAPIKey(name, config.GetAPIKey(name)); err != nil {
					log.Printf("INTELLIGENCE: Failed to set API key for %s: %v", name, err)
				}
			} else {
				log.Printf("INTELLIGENCE: %s API is not available, skipping", name)
			}
		}
	}

	return &Operator{
		client:       client,
		promptEngine: grompt.NewPromptEngine(config),
		analytics:    analyticsOp.LinesOfCode,
		productivity: productivityOp,
	}
}

// GenerateReport creates a comprehensive, humanized analysis report
func (o *Operator) GenerateReport(ctx context.Context, owner, repo string, days int) (*HumanizedReport, error) {
	start := time.Now()
	log.Printf("INTELLIGENCE: Starting humanized report generation for %s/%s", owner, repo)

	// Gather raw data from other operators
	analyticsData, err := o.analytics.GetInsights(ctx, owner, repo, days)
	if err != nil {
		log.Printf("INTELLIGENCE: Failed to get analytics data: %v", err)
		return nil, fmt.Errorf("failed to get analytics data: %w", err)
	}

	productivityData, err := o.productivity.GetRepositoryActions(ctx, owner, repo)
	if err != nil {
		log.Printf("INTELLIGENCE: Failed to get productivity data: %v", err)
		return nil, fmt.Errorf("failed to get productivity data: %w", err)
	}

	// Create comprehensive data structure for AI processing
	reportData := map[string]interface{}{
		"repository":    fmt.Sprintf("%s/%s", owner, repo),
		"analytics":     analyticsData,
		"productivity":  productivityData,
		"analysis_date": time.Now().Format("2006-01-02"),
		"days_analyzed": days,
	}

	// Convert to JSON for AI processing
	dataJSON, err := json.MarshalIndent(reportData, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create AI prompt for comprehensive analysis
	prompt := o.createAnalysisPrompt(string(dataJSON), owner, repo)

	// If AI is not available, generate a structured report from raw data
	// This ensures ZERO MOCKS - we always return real data
	report := o.generateFallbackReport(analyticsData, productivityData, owner, repo)

	// Try to enhance with AI if available
	if o.promptEngine != nil {
		enhancedReport, err := o.enhanceWithAI(prompt, report)
		if err != nil {
			log.Printf("INTELLIGENCE: AI enhancement failed, using fallback: %v", err)
		} else {
			report = enhancedReport
		}
	}

	log.Printf("INTELLIGENCE: Report generation completed in %v", time.Since(start))
	return report, nil
}

// createAnalysisPrompt creates a comprehensive prompt for AI analysis
func (o *Operator) createAnalysisPrompt(dataJSON, owner, repo string) string {
	return fmt.Sprintf(`
You are a senior DevOps consultant and GitHub repository expert. Analyze the following repository data and provide insights.

Repository: %s/%s

Raw Data:
%s

Please provide a comprehensive analysis including:

1. OVERALL ASSESSMENT:
- Executive summary in 2-3 sentences
- Letter grade (A-F) with justification
- Key strengths and weaknesses
- Trend analysis (improving/stable/declining)

2. KEY INSIGHTS:
- 3-5 most important findings
- Evidence-based observations
- Impact assessment for each finding

3. RECOMMENDATIONS:
- Prioritized list of actionable improvements
- Effort vs impact analysis
- Timeline suggestions
- Resource requirements

4. HEALTH BREAKDOWN:
- Detailed health category analysis
- Critical issues identification
- Improvement opportunities

5. PRODUCTIVITY OPTIMIZATION:
- Specific tips for team productivity
- ROI-focused suggestions
- Quick wins vs long-term investments

6. RISK ASSESSMENT:
- Potential risks and their likelihood
- Mitigation strategies
- Dependencies and bottlenecks

7. NEXT STEPS:
- Concrete action plan
- Priority order
- Owner assignments
- Timeline expectations

Focus on practical, actionable advice that development teams can implement immediately.
Use clear, professional language suitable for both technical and non-technical stakeholders.
`, owner, repo, dataJSON)
}

// generateFallbackReport creates a structured report without AI enhancement
func (o *Operator) generateFallbackReport(analyticsData *analytics.InsightsReport, productivityData *productivity.ProductivityAction, owner, repo string) *HumanizedReport {
	// Calculate overall grade based on health score
	grade := "F"
	if analyticsData.HealthScore.Overall >= 90 {
		grade = "A"
	} else if analyticsData.HealthScore.Overall >= 70 {
		grade = "B"
	} else if analyticsData.HealthScore.Overall >= 60 {
		grade = "C"
	} else if analyticsData.HealthScore.Overall >= 50 {
		grade = "D"
	}

	// Generate insights based on real data
	insights := o.generateInsightsFromData(analyticsData, productivityData)
	recommendations := o.generateRecommendationsFromData(analyticsData, productivityData)
	productivityTips := o.generateProductivityTipsFromData(productivityData)
	riskFactors := o.generateRiskFactorsFromData(analyticsData)

	productivityReport, err := productivity.AnalyzeProductivity(
		contexto.Background(),
		github.Client,
		owner,
		"",
	)
	if err != nil {
		log.Printf("Failed to analyze productivity data: %v", err)
	}

	return &HumanizedReport{
		RepositoryName: fmt.Sprintf("%s/%s", owner, repo),
		OverallAssessment: OverallAssessment{
			Grade:         grade,
			Score:         analyticsData.HealthScore.Overall,
			Summary:       o.generateSummaryFromData(analyticsData, productivityReport, grade),
			KeyWeaknesses: o.identifyWeaknesses(analyticsData.CodeIntel, productivityData),
			Trend:         o.assessTrend(analyticsData.CodeIntel),
		},
		KeyInsights:     insights,
		Recommendations: recommendations,
		HealthSummary: HealthSummary{
			OverallHealth:    analyticsData.HealthScore.Overall,
			CategoryScores:   o.extractCategoryScores(analyticsData.CodeIntel),
			HealthTrend:      o.assessTrend(analyticsData.CodeIntel),
			CriticalIssues:   o.identifyCriticalIssues(analyticsData),
			ImprovementAreas: o.identifyImprovementAreas(analyticsData, productivityData),
		},
		ProductivityTips: productivityTips,
		RiskFactors:      riskFactors,
		NextSteps:        o.generateNextSteps(recommendations),
		GeneratedAt:      time.Now(),
		Metadata: map[string]interface{}{
			"source":        "ghbex-intelligence",
			"ai_enhanced":   false,
			"data_sources":  []string{"analytics", "productivity"},
			"analysis_type": "comprehensive",
		},
	}
}

// enhanceWithAI attempts to improve the report using AI processing
func (o *Operator) enhanceWithAI(prompt string, baseReport *HumanizedReport) (*HumanizedReport, error) {
	// Try to process with AI
	result, err := o.promptEngine.ProcessPrompt(prompt, map[string]interface{}{
		"base_report": baseReport,
	})
	if err != nil {
		return nil, fmt.Errorf("AI processing failed: %w", err)
	}

	// For now, return the base report with AI metadata
	// In a full implementation, we would parse the AI response
	// and enhance the report structure
	baseReport.Metadata["ai_enhanced"] = true
	baseReport.Metadata["ai_response"] = result.Response
	baseReport.Metadata["ai_provider"] = result.Provider

	return baseReport, nil
}

// Helper methods for data analysis
func (o *Operator) generateSummaryFromData(analytics *analytics.InsightsReport, productivity *productivity.ProductivityReport, grade string) string {
	repoInfo := analytics
	health := analytics.HealthScore.Overall
	actions := len(productivity.Actions)

	return fmt.Sprintf("Repository %s shows a health score of %.1f/100 (Grade %s) with %d recommended productivity improvements. Primary language is %s with %d contributors and %d open issues.",
		repoInfo.Repo,
		health,
		grade,
		actions,
		repoInfo.CodeIntel.PrimaryLanguage,
		repoInfo.AnalysisDays,
		repoInfo.OpenIssuesCount,
	)
}

func (o *Operator) identifyStrengths(
	analytics *analytics.CodeIntelligence,
	productivity *productivity.ProductivityAction,
) {
	strengths := []string{}

	if analytics.Complexity.MaintainabilityIndex > 70 {
		strengths = append(strengths, "Active maintenance and regular updates")
	}
	if analytics.Complexity.MaintainabilityIndex < 3.0 {
		strengths = append(strengths, "Maintainable code complexity")
	}
	highestKPI := productivity.Impact

	if highestKPI > 0 {
		strengths = append(strengths, "High-impact productivity opportunities")
	}
}

func (o *Operator) identifyWeaknesses(
	insReport *analytics.CodeIntelligence,
	actReport *productivity.ProductivityAction,
) []string {
	weaknesses := []string{}

	if insReport.Complexity.MaintainabilityIndex < 50 {
		weaknesses = append(weaknesses, "Below-average maintainability index")
	}
	if insReport.Complexity.CyclomaticComplexity > 10 {
		weaknesses = append(weaknesses, "High cyclomatic complexity")
	}
	if insReport.Complexity.CodeDuplication > 10 {
		weaknesses = append(weaknesses, "High code duplication")
	}
	if len(actReport.Steps) > 5 {
		weaknesses = append(weaknesses, "Multiple productivity optimization opportunities")
	}

	return weaknesses
}

func (o *Operator) assessTrend(analytics *analytics.CodeIntelligence) string {
	// Simple trend assessment based on available data
	if analytics.Complexity.MaintainabilityIndex > 70 {
		return "stable"
	} else if analytics.Complexity.CyclomaticComplexity > 50 {
		return "stable"
	}
	return "needs attention"
}

func (o *Operator) extractCategoryScores(analytics *analytics.CodeIntelligence) map[string]float64 {
	return map[string]float64{
		"overall":     analytics.Complexity.CodeDuplication,
		"maintenance": analytics.Complexity.MaintainabilityIndex,
		// "community":   analytics.Complexity.CommunityIndex,
		"dependency": float64(analytics.Dependencies.TotalDependencies),
		"devex":      analytics.Complexity.CyclomaticComplexity,
	}
}

func (o *Operator) identifyCriticalIssues(analytics *analytics.InsightsReport) []string {
	issues := []string{}

	if analytics.HealthScore.Overall < 30 {
		issues = append(issues, "Critically low health score requires immediate attention")
	}
	if analytics.Productivity.DevexScore < 30 {
		issues = append(issues, "Dependency vulnerabilities or outdated packages")
	}
	if analytics.Productivity.BugFixRate > 50 {
		issues = append(issues, "Excessive backlog of open issues")
	}

	return issues
}

func (o *Operator) identifyImprovementAreas(analytics *analytics.InsightsReport, productivity *productivity.ProductivityAction) []string {
	areas := []string{}

	if analytics.Community.Growth.Forks.Growth < 50 {
		areas = append(areas, "Community engagement and contributor attraction")
	}
	if analytics.CodeIntel.Complexity.MaintainabilityIndex < 10 {
		areas = append(areas, "Regular maintenance and update practices")
	}
	if len(analytics.Recommendations) > 0 {
		areas = append(areas, "Workflow automation and productivity optimization")
	}

	return areas
}

func (o *Operator) generateInsightsFromData(analytics *analytics.InsightsReport, productivity *productivity.ProductivityAction) []KeyInsight {
	insights := []KeyInsight{}

	// Health insight
	impact := "medium"
	if analytics.HealthMetrics.OverallScore < 50 {
		impact = "high"
	}

	insights = append(insights, KeyInsight{
		Category:    "Repository Health",
		Title:       fmt.Sprintf("Health Score: %.1f/100", analytics.HealthMetrics.OverallScore),
		Description: "Overall repository health based on maintenance, community, and code quality metrics",
		Impact:      impact,
		Evidence: fmt.Sprintf("Score breakdown: Maintenance %.1f, Community %.1f, DevEx %.1f",
			analytics.HealthMetrics.MaintenanceScore, analytics.HealthMetrics.CommunityScore, analytics.HealthMetrics.DevExScore),
	})

	// Productivity insight
	if len(productivity.RecommendedActions) > 0 {
		insights = append(insights, KeyInsight{
			Category:    "Productivity",
			Title:       fmt.Sprintf("%d Optimization Opportunities", len(productivity.RecommendedActions)),
			Description: "Identified opportunities to improve team productivity and workflow efficiency",
			Impact:      "medium",
			Evidence:    fmt.Sprintf("Potential savings: $%.0f, ROI: %.1fx", productivity.Summary.TotalSavings, productivity.Summary.ROI),
		})
	}

	// Code complexity insight
	if analytics.CodeIntelligence.Complexity > 0 {
		complexityImpact := "low"
		if analytics.CodeIntelligence.Complexity > 5 {
			complexityImpact = "high"
		} else if analytics.CodeIntelligence.Complexity > 3 {
			complexityImpact = "medium"
		}

		insights = append(insights, KeyInsight{
			Category:    "Code Quality",
			Title:       fmt.Sprintf("Code Complexity: %.2f", analytics.CodeIntelligence.Complexity),
			Description: "Cyclomatic complexity indicating code maintainability",
			Impact:      complexityImpact,
			Evidence: fmt.Sprintf("Total lines: %d, Code lines: %d, Comments: %d",
				analytics.CodeIntelligence.TotalLines, analytics.CodeIntelligence.CodeLines, analytics.CodeIntelligence.CommentLines),
		})
	}

	return insights
}

func (o *Operator) generateRecommendationsFromData(analytics *analytics.InsightsResponse, productivity *productivity.ActionsResponse) []Recommendation {
	recommendations := []Recommendation{}

	// Convert productivity actions to recommendations
	for _, action := range productivity.RecommendedActions {
		priority := "medium"
		if action.Priority == "high" {
			priority = "high"
		} else if action.Priority == "low" {
			priority = "low"
		}

		recommendations = append(recommendations, Recommendation{
			Priority:    priority,
			Title:       action.Title,
			Description: action.Description,
			Effort:      action.Effort,
			Impact:      action.Impact,
			Timeline:    action.Timeline,
			Resources:   []string{action.Type},
		})
	}

	// Add health-based recommendations
	if analytics.HealthMetrics.OverallScore < 50 {
		recommendations = append(recommendations, Recommendation{
			Priority:    "high",
			Title:       "Improve Repository Health",
			Description: "Focus on maintenance, community engagement, and code quality improvements",
			Effort:      "medium",
			Impact:      "high",
			Timeline:    "medium",
			Resources:   []string{"maintenance", "community", "documentation"},
		})
	}

	return recommendations
}

func (o *Operator) generateProductivityTipsFromData(productivity *productivity.ActionsResponse) []ProductivityTip {
	tips := []ProductivityTip{}

	for _, action := range productivity.RecommendedActions {
		tip := ProductivityTip{
			Area:       action.Type,
			Tip:        action.Title,
			Benefit:    action.Description,
			Difficulty: action.Effort,
			ROI:        fmt.Sprintf("%.1fx", action.ROI),
		}
		tips = append(tips, tip)
	}

	return tips
}

func (o *Operator) generateRiskFactorsFromData(analytics *analytics.InsightsResponse) []RiskFactor {
	risks := []RiskFactor{}

	if analytics.HealthMetrics.DependencyScore < 50 {
		risks = append(risks, RiskFactor{
			Type:        "Dependencies",
			Level:       "medium",
			Description: "Potential security or compatibility issues with dependencies",
			Mitigation:  "Regular dependency updates and security audits",
			Probability: "medium",
		})
	}

	if analytics.RepositoryInfo.OpenIssues > 30 {
		risks = append(risks, RiskFactor{
			Type:        "Maintenance",
			Level:       "medium",
			Description: "Growing backlog may indicate maintenance challenges",
			Mitigation:  "Implement issue triage and resolution workflows",
			Probability: "high",
		})
	}

	return risks
}

func (o *Operator) generateNextSteps(recommendations []Recommendation) []NextStep {
	steps := []NextStep{}

	// Convert high-priority recommendations to next steps
	order := 1
	for _, rec := range recommendations {
		if rec.Priority == "high" || rec.Priority == "critical" {
			steps = append(steps, NextStep{
				Order:        order,
				Action:       rec.Title,
				Owner:        "Development Team",
				Timeline:     rec.Timeline,
				Dependencies: []string{},
			})
			order++
		}
	}

	// Add general next steps if none from recommendations
	if len(steps) == 0 {
		steps = append(steps, NextStep{
			Order:        1,
			Action:       "Review repository health metrics",
			Owner:        "Tech Lead",
			Timeline:     "immediate",
			Dependencies: []string{},
		})
	}

	return steps
}

// GetRepositoryActions
