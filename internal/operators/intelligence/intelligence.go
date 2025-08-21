// Package intelligence provides AI-powered analysis and insights for GitHub repositories.
package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"

	"github.com/rafa-mori/ghbex/internal/defs"
	"github.com/rafa-mori/ghbex/internal/interfaces"
	gl "github.com/rafa-mori/ghbex/internal/module/logger"

	configLib "github.com/rafa-mori/ghbex/internal/config"
)

// IntelligenceOperator provides AI-powered analysis using Grompt engine
type IntelligenceOperator struct {
	client       *github.Client
	promptEngine defs.PromptEngine
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

// NewIntelligenceOperator creates a new Intelligence operator
func NewIntelligenceOperator(cfg interfaces.IMainConfig, client *github.Client) *IntelligenceOperator {
	// Initialize Grompt with basic config
	var port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey string
	port = cfg.GetServer().GetPort()
	openAIKey = configLib.GetEnvOrDefault("OPENAI_API_KEY", "")
	deepSeekKey = configLib.GetEnvOrDefault("DEEPSEEK_API_KEY", "")
	ollamaEndpoint = configLib.GetEnvOrDefault("OLLAMA_API_ENDPOINT", "")
	claudeKey = configLib.GetEnvOrDefault("CLAUDE_API_KEY", "")
	geminiKey = configLib.GetEnvOrDefault("GEMINI_API_KEY", "")

	gromptEngineCfg := defs.NewGromptConfig(
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
	)
	engine := defs.NewPromptEngine(gromptEngineCfg)

	llmMapList := map[string]defs.APIConfig{
		"openai":   gromptEngineCfg.GetAPIConfig("openai"),
		"claude":   gromptEngineCfg.GetAPIConfig("claude"),
		"gemini":   gromptEngineCfg.GetAPIConfig("gemini"),
		"chatgpt":  gromptEngineCfg.GetAPIConfig("chatgpt"),
		"deepseek": gromptEngineCfg.GetAPIConfig("deepseek"),
		"ollama":   gromptEngineCfg.GetAPIConfig("ollama"),
	}
	for key, apiKey := range llmMapList {
		if apiKey == nil || apiKey.IsDemoMode() {
			apiFromEnv := configLib.GetEnvOrDefault(
				key,
				"",
			)
			if apiFromEnv != "" {
				gl.Log("notice", fmt.Sprintf("Using API key from environment for %s", key))
				gromptEngineCfg.SetAPIKey(key, apiFromEnv)
			} else {
				gl.Log("debug", fmt.Sprintf("No API key configured for %s, using default config", key))
			}
		}
	}
	providers := engine.GetProviders()
	if len(providers) == 0 {
		gl.Log("warn", "INTELLIGENCE: No AI providers configured, using default Grompt settings")
	} else {
		gl.Log("info", fmt.Sprintf("INTELLIGENCE: Available AI providers: %d", len(providers)))
		for _, provider := range providers {
			gl.Log("info", fmt.Sprintf(" - %s: %v", provider.Name(), provider.GetCapabilities()))
		}
	}
	if len(llmMapList) == 0 {
		gl.Log("warn", "INTELLIGENCE: No LLM providers configured, using default settings")
	} else {
		gl.Log("info", fmt.Sprintf("INTELLIGENCE: Available LLM providers: %d", len(llmMapList)))
		for key, apiKey := range llmMapList {
			apiTk := gromptEngineCfg.GetAPIKey(key)
			if apiTk != "" {
				gl.Log("info", fmt.Sprintf(" - %s: %s (%v)", key, apiTk[:7], apiKey.IsAvailable()))
			}
		}
	}

	return &IntelligenceOperator{
		client:       client,
		promptEngine: engine,
	}
}

// GenerateQuickInsight creates AI-powered insights for repository cards
func (o *IntelligenceOperator) GenerateQuickInsight(ctx context.Context, owner, repo string) (*RepositoryInsight, error) {
	gl.Log("info", fmt.Sprintf("INTELLIGENCE: Generating quick insight for %s/%s", owner, repo))

	// Get basic repository info
	repoInfo, _, err := o.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return o.generateFallbackInsight(owner, repo), nil
	}

	// Generate AI-powered assessment using Grompt
	aiScore, assessment, err := o.analyzeRepositoryWithAI(ctx, repoInfo)
	if err != nil {
		gl.Log("error", fmt.Sprintf("INTELLIGENCE: AI analysis failed, using fallback: %v", err))
		return o.generateFallbackInsight(owner, repo), nil
	}

	insight := &RepositoryInsight{
		RepositoryName:  fmt.Sprintf("%s/%s", owner, repo),
		AIScore:         aiScore,
		QuickAssessment: assessment,
		HealthIcon:      o.getHealthIcon(aiScore),
		MainTag:         o.generateMainTag(repoInfo),
		RiskLevel:       o.calculateRiskLevel(repoInfo, aiScore),
		Opportunity:     o.identifyOpportunity(repoInfo),
		LastAnalyzed:    time.Now(),
	}

	return insight, nil
}

// GenerateSmartRecommendations creates contextual AI recommendations
func (o *IntelligenceOperator) GenerateSmartRecommendations(ctx context.Context, owner, repo string) ([]SmartRecommendation, error) {
	gl.Log("info", fmt.Sprintf("INTELLIGENCE: Generating smart recommendations for %s/%s", owner, repo))

	// Get repository data
	repoInfo, _, err := o.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return o.generateFallbackRecommendations(owner, repo), nil
	}

	// Get recent issues and PRs for context
	issues, _, err := o.client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 10},
	})
	if err != nil {
		gl.Log("error", fmt.Sprintf("INTELLIGENCE: Failed to get issues: %v", err))
	}

	// Generate AI recommendations
	recommendations, err := o.generateAIRecommendations(ctx, repoInfo, issues)
	if err != nil {
		gl.Log("error", fmt.Sprintf("INTELLIGENCE: AI recommendations failed, using fallback: %v", err))
		return o.generateFallbackRecommendations(owner, repo), nil
	}

	return recommendations, nil
}

// analyzeRepositoryWithAI uses Grompt to analyze repository
func (o *IntelligenceOperator) analyzeRepositoryWithAI(ctx context.Context, repo *github.Repository) (float64, string, error) {
	defer func(c context.Context) {
		if err := recover(); err != nil {
			gl.Log("error", fmt.Sprintf("INTELLIGENCE: AI analysis failed: %v", err))
		}
		if ctx.Err() != nil {
			gl.Log("warn", fmt.Sprintf("INTELLIGENCE: AI analysis canceled: %v", ctx.Err()))
		}
		gl.Log("info", "INTELLIGENCE: AI analysis completed")
	}(ctx)

	prompt := fmt.Sprintf(`
Analyze this GitHub repository and provide a quick assessment:

Repository: %s
Description: %s
Language: %s
Stars: %d
Forks: %d
Open Issues: %d
Created: %s
Last Updated: %s

Please provide:
1. A score from 0-100 based on repository health and activity
2. A brief 1-sentence assessment focusing on the most important aspect

Format your response as JSON:
{
	"score": 85.5,
	"assessment": "Active Go project with good community engagement and recent updates"
}
`,
		repo.GetFullName(),
		repo.GetDescription(),
		repo.GetLanguage(),
		repo.GetStargazersCount(),
		repo.GetForksCount(),
		repo.GetOpenIssuesCount(),
		repo.GetCreatedAt().Format("2006-01-02"),
		repo.GetUpdatedAt().Format("2006-01-02"),
	)

	response, err := o.promptEngine.ProcessPrompt(prompt, map[string]interface{}{})
	if err != nil {
		return 0, "", err
	}

	// Parse JSON response
	var result struct {
		Score      float64 `json:"score"`
		Assessment string  `json:"assessment"`
	}

	if err := json.Unmarshal([]byte(response.Response), &result); err != nil {
		// Fallback if JSON parsing fails - CLEARLY MARKED AS SIMULATED
		gl.Log("warn", fmt.Sprintf("AI parsing failed for %s, using simulated data", repo.GetFullName()))
		return 0.0, "⚠️  SIMULATED - AI analysis unavailable", nil
	}

	return result.Score, result.Assessment, nil
}

// generateAIRecommendations creates smart recommendations using AI
func (o *IntelligenceOperator) generateAIRecommendations(ctx context.Context, repo *github.Repository, issues []*github.Issue) ([]SmartRecommendation, error) {
	defer func(c context.Context) {
		if err := recover(); err != nil {
			gl.Log("error", fmt.Sprintf("INTELLIGENCE: AI recommendations failed: %v", err))
		}
		if ctx.Err() != nil {
			gl.Log("warn", fmt.Sprintf("INTELLIGENCE: AI recommendations canceled: %v", ctx.Err()))
		}
		gl.Log("info", "INTELLIGENCE: AI recommendations completed")
	}(ctx)

	issuesContext := ""
	if len(issues) > 0 {
		issuesContext = fmt.Sprintf("Recent issues: %d open, latest: '%s'",
			repo.GetOpenIssuesCount(),
			issues[0].GetTitle())
	}

	prompt := fmt.Sprintf(`
Analyze this repository and suggest 3 specific, actionable recommendations:

Repository: %s (%s)
%s

Consider:
- Security improvements
- Performance optimizations
- Maintenance tasks
- Feature enhancements

Provide recommendations as JSON array:
[
	{
		"type": "security",
		"title": "Enable Dependabot",
		"description": "Automatically scan for vulnerable dependencies",
		"impact": "high",
		"effort": "low",
		"urgency": "medium"
	}
]
`, repo.GetFullName(), repo.GetLanguage(), issuesContext)

	response, err := o.promptEngine.ProcessPrompt(prompt, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var recommendations []SmartRecommendation
	if err := json.Unmarshal([]byte(response.Response), &recommendations); err != nil {
		// Return fallback recommendations
		return o.generateFallbackRecommendations(repo.GetOwner().GetLogin(), repo.GetName()), nil
	}

	// Add metadata
	for i := range recommendations {
		recommendations[i].ID = fmt.Sprintf("%s-%d", repo.GetName(), i+1)
		recommendations[i].GeneratedAt = time.Now()
	}

	return recommendations, nil
}

// Fallback methods for when AI is not available
func (o *IntelligenceOperator) generateFallbackInsight(owner, repo string) *RepositoryInsight {
	gl.Log("info", fmt.Sprintf("Using SIMULATED insight for %s/%s - AI analysis not available", owner, repo))

	return &RepositoryInsight{
		RepositoryName:  fmt.Sprintf("%s/%s", owner, repo),
		AIScore:         0.0, // Clear indicator this is not real
		QuickAssessment: "⚠️  SIMULATED DATA - AI analysis unavailable",
		HealthIcon:      "⚠️",
		MainTag:         "DEMO",
		RiskLevel:       "unknown",
		Opportunity:     "⚠️  Enable AI analysis for real insights",
		LastAnalyzed:    time.Now(),
	}
}

func (o *IntelligenceOperator) generateFallbackRecommendations(owner, repo string) []SmartRecommendation {
	gl.Log("info", fmt.Sprintf("Using SIMULATED recommendations for %s/%s - AI analysis not available", owner, repo))

	return []SmartRecommendation{
		{
			ID:          fmt.Sprintf("DEMO-%s-1", repo),
			Type:        "warning",
			Title:       "⚠️  SIMULATED DATA - Enable AI Analysis",
			Description: "This is demonstration data. Configure AI providers for real insights.",
			Impact:      "demo",
			Effort:      "demo",
			Urgency:     "demo",
			GeneratedAt: time.Now(),
		},
		{
			ID:          fmt.Sprintf("DEMO-%s-2", repo),
			Type:        "info",
			Title:       "🔧 Configure Grompt Integration",
			Description: "Set up OpenAI, Claude, or other AI providers for real analysis.",
			Impact:      "demo",
			Effort:      "demo",
			Urgency:     "demo",
			GeneratedAt: time.Now(),
		},
		{
			ID:          fmt.Sprintf("DEMO-%s-3", repo),
			Type:        "placeholder",
			Title:       "📊 Real Insights Available Soon",
			Description: "Connect AI services to get actionable repository recommendations.",
			Impact:      "demo",
			Effort:      "demo",
			Urgency:     "demo",
			GeneratedAt: time.Now(),
		},
	}
}

// Helper methods
func (o *IntelligenceOperator) getHealthIcon(score float64) string {
	if score >= 90 {
		return "🟢"
	} else if score >= 70 {
		return "🟡"
	} else {
		return "🔴"
	}
}

func (o *IntelligenceOperator) generateMainTag(repo *github.Repository) string {
	if repo.GetStargazersCount() > 100 {
		return "Popular"
	} else if repo.GetUpdatedAt().After(time.Now().AddDate(0, 0, -7)) {
		return "Active"
	} else if repo.GetLanguage() != "" {
		return repo.GetLanguage()
	}
	return "Project"
}

func (o *IntelligenceOperator) calculateRiskLevel(repo *github.Repository, aiScore float64) string {
	if aiScore < 60 || repo.GetOpenIssuesCount() > 50 {
		return "high"
	} else if aiScore < 80 || repo.GetOpenIssuesCount() > 20 {
		return "medium"
	}
	return "low"
}

func (o *IntelligenceOperator) identifyOpportunity(repo *github.Repository) string {
	opportunities := []string{
		"Documentation enhancement",
		"Performance optimization",
		"Security improvements",
		"Community engagement",
		"Code quality boost",
		"Test coverage expansion",
	}

	// Simple deterministic selection based on repo characteristics
	index := (repo.GetStargazersCount() + repo.GetForksCount()) % len(opportunities)
	return opportunities[index]
}
