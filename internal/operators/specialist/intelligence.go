// Package specialist provides AI-powered analysis and insights for GitHub repositories.
package specialist

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/grompt"
)

// IntelligenceOperator provides AI-powered analysis using Grompt engine
type IntelligenceOperator struct {
	client       *github.Client
	promptEngine grompt.PromptEngine
}

// RepositoryInsight provides quick AI insights for repository cards
type RepositoryInsight struct {
	RepositoryName  string    `json:"repository_name"`
	AIScore         float64   `json:"ai_score"`
	QuickAssessment string    `json:"quick_assessment"`
	HealthIcon      string    `json:"health_icon"`
	MainTag         string    `json:"main_tag"`
	RiskLevel       string    `json:"risk_level"`
	Opportunity     string    `json:"opportunity"`
	LastAnalyzed    time.Time `json:"last_analyzed"`
}

// SmartRecommendation provides contextual recommendations
type SmartRecommendation struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // "security", "performance", "maintenance", "enhancement"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Impact      string    `json:"impact"`
	Effort      string    `json:"effort"`
	Urgency     string    `json:"urgency"`
	GeneratedAt time.Time `json:"generated_at"`
}

// HumanizedReport represents a comprehensive AI analysis
type HumanizedReport struct {
	RepositoryName    string                 `json:"repository_name"`
	OverallAssessment OverallAssessment      `json:"overall_assessment"`
	KeyInsights       []KeyInsight           `json:"key_insights"`
	Recommendations   []SmartRecommendation  `json:"recommendations"`
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

// NewIntelligenceOperator creates a new Intelligence operator
func NewIntelligenceOperator(client *github.Client) *IntelligenceOperatorA {
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

	return &IntelligenceOperatorA{
		client:       client,
		promptEngine: grompt.NewPromptEngine(config),
	}
}

// GenerateQuickInsight creates AI-powered insights for repository cards
func (o *IntelligenceOperatorA) GenerateQuickInsight(ctx context.Context, owner, repo string) (*RepositoryInsight, error) {
	log.Printf("INTELLIGENCE: Generating quick insight for %s/%s", owner, repo)

	// Get basic repository info
	repoInfo, _, err := o.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return o.generateFallbackInsight(owner, repo), nil
	}

	// Generate AI-powered assessment using Grompt
	aiScore, assessment, err := o.analyzeRepositoryWithAI(ctx, repoInfo)
	if err != nil {
		log.Printf("INTELLIGENCE: AI analysis failed, using fallback: %v", err)
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
func (o *IntelligenceOperatorA) GenerateSmartRecommendations(ctx context.Context, owner, repo string) ([]SmartRecommendation, error) {
	log.Printf("INTELLIGENCE: Generating smart recommendations for %s/%s", owner, repo)

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
		log.Printf("INTELLIGENCE: Failed to get issues: %v", err)
	}

	// Generate AI recommendations
	recommendations, err := o.generateAIRecommendations(ctx, repoInfo, issues)
	if err != nil {
		log.Printf("INTELLIGENCE: AI recommendations failed, using fallback: %v", err)
		return o.generateFallbackRecommendations(owner, repo), nil
	}

	return recommendations, nil
}

// analyzeRepositoryWithAI uses Grompt to analyze repository
func (o *IntelligenceOperatorA) analyzeRepositoryWithAI(ctx context.Context, repo *github.Repository) (float64, string, error) {
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
		log.Printf("⚠️  WARNING: AI parsing failed for %s, using simulated data", repo.GetFullName())
		return 0.0, "⚠️  SIMULATED - AI analysis unavailable", nil
	}

	return result.Score, result.Assessment, nil
}

// generateAIRecommendations creates smart recommendations using AI
func (o *IntelligenceOperatorA) generateAIRecommendations(ctx context.Context, repo *github.Repository, issues []*github.Issue) ([]SmartRecommendation, error) {
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
	log.Printf("⚠️  WARNING: Using SIMULATED data for %s/%s - AI analysis not available", owner, repo)

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
	log.Printf("⚠️  WARNING: Using SIMULATED recommendations for %s/%s - AI analysis not available", owner, repo)

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
func (o *IntelligenceOperatorA) getHealthIcon(score float64) string {
	if score >= 90 {
		return "🟢"
	} else if score >= 70 {
		return "🟡"
	} else {
		return "🔴"
	}
}

func (o *IntelligenceOperatorA) generateMainTag(repo *github.Repository) string {
	if repo.GetStargazersCount() > 100 {
		return "Popular"
	} else if repo.GetUpdatedAt().After(time.Now().AddDate(0, 0, -7)) {
		return "Active"
	} else if repo.GetLanguage() != "" {
		return repo.GetLanguage()
	}
	return "Project"
}

func (o *IntelligenceOperatorA) calculateRiskLevel(repo *github.Repository, aiScore float64) string {
	if aiScore < 60 || repo.GetOpenIssuesCount() > 50 {
		return "high"
	} else if aiScore < 80 || repo.GetOpenIssuesCount() > 20 {
		return "medium"
	}
	return "low"
}

func (o *IntelligenceOperatorA) identifyOpportunity(repo *github.Repository) string {
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
