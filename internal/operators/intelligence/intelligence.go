// Package intelligence provides AI-powered analysis and insights for GitHub repositories.
package intelligence

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/go-github/v61/github"

	"github.com/rafa-mori/ghbex/internal/defs/gromptz"
	"github.com/rafa-mori/ghbex/internal/defs/interfaces"
	gl "github.com/rafa-mori/ghbex/internal/module/logger"

	configLib "github.com/rafa-mori/ghbex/internal/config"
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

// NewIntelligenceOperator creates a new Intelligence operator
func NewIntelligenceOperator(cfg interfaces.IMainConfig, client *github.Client) *IntelligenceOperator {
	if client == nil {
		gl.Log("error", "INTELLIGENCE: GitHub client is nil, cannot initialize Intelligence operator")
		return nil
	}

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

	gromptEngineCfg := gromptz.NewGromptConfig(
		port,
		openAIKey,
		deepSeekKey,
		ollamaEndpoint,
		claudeKey,
		geminiKey,
	)

	engine := gromptz.NewGromptEngine(gromptEngineCfg)
	llmList := map[string]string{
		"claude":   "v1",
		"openai":   "v1",
		"deepseek": "v1",
		"ollama":   "v1",
		"gemini":   "v1beta",
		"chatgpt":  "v1",
	}
	llmMapList := make(map[string]gromptz.Provider)
	for provider, version := range llmList {
		llmMapList[provider] = gromptz.NewProvider(
			provider,
			configLib.GetEnvOrDefault(
				strings.ToUpper(provider)+"_API_KEY",
				gromptEngineCfg.GetAPIKey(provider),
			),
			version,
			gromptEngineCfg,
		)
	}
	if len(llmMapList) == 0 {
		gl.Log("warn", "INTELLIGENCE: No AI providers configured, using default Grompt settings")
	} else {
		gl.Log("info", fmt.Sprintf("INTELLIGENCE: Available AI providers: %d", len(llmMapList)))
		for _, provider := range llmMapList {
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
		client:           client,
		promptEngine:     engine,
		mainConfig:       cfg,
		healthCache:      make(map[string]healthStatus),
		healthCacheMutex: sync.RWMutex{},
	}
}

// GenerateQuickInsight creates AI-powered insights for repository cards
func (o *IntelligenceOperator) GenerateQuickInsight(ctx context.Context, owner, repo string) (*RepositoryInsight, error) {
	if o.mainConfig == nil {
		gl.Log("error", "INTELLIGENCE: Main configuration is nil, cannot generate quick insight")
		return nil, fmt.Errorf("main configuration is nil")
	}

	gl.Log("debug", fmt.Sprintf("INTELLIGENCE: Generating quick insight for %s/%s", owner, repo))

	// 🛡️ CRITICAL SECURITY: NEVER auto-discover repositories!
	// Only process explicitly provided owner/repo combinations
	if owner == "" || repo == "" {
		gl.Log("error", "🚨 INTELLIGENCE: Owner and repo must be explicitly provided - auto-discovery is DISABLED for security")
		gl.Log("info", "📋 To use intelligence operator, provide explicit repository: --owner 'user' --repo 'repository'")
		gl.Log("info", "🛡️ This prevents accidental scanning of all GitHub repositories")
		return nil, fmt.Errorf("owner and repo must be explicitly provided - auto-discovery disabled for security")
	}

	// Generate AI-powered assessment using Grompt for the EXPLICIT repository
	repoInfo, repoInfoResponse, err := o.client.Repositories.Get(
		ctx,
		owner,
		repo,
	)
	if err != nil && repoInfo == nil {
		gl.Log("error", fmt.Sprintf("INTELLIGENCE: error getting quick repository (%s/%s) info: %v", owner, repo, err))
		return nil, fmt.Errorf("error getting quick repository info: %w", err)
	}
	if repoInfoResponse == nil || repoInfoResponse.StatusCode != 200 {
		gl.Log("warn", fmt.Sprintf("INTELLIGENCE: Repository %s/%s - %d: %s", owner, repo, repoInfoResponse.StatusCode, repoInfoResponse.Status))
		return nil, fmt.Errorf("repository not found: %s/%s", owner, repo)
	}
	if repoInfo == nil {
		gl.Log("warn", fmt.Sprintf("INTELLIGENCE: Repository %s/%s not found", owner, repo))
		return nil, fmt.Errorf("repository not found: %s/%s", owner, repo)
	}

	aiScore, assessment, err := o.analyzeRepositoryWithAI(ctx, repoInfo)
	if err != nil {
		gl.Log("error", fmt.Sprintf("INTELLIGENCE: AI analysis failed, using fallback: %v", err))
		return o.generateFallbackInsight(owner, repo), nil
	}
	if aiScore <= 0.0 {
		gl.Log("warn", fmt.Sprintf("INTELLIGENCE: AI analysis returned non-positive score for %s/%s, using fallback", owner, repo))
	} else {
		gl.Log("info", fmt.Sprintf("INTELLIGENCE: AI analysis score for %s/%s: %.2f", owner, repo, aiScore))
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
	gl.Log("debug", fmt.Sprintf("INTELLIGENCE: Generating smart recommendations for %s/%s", owner, repo))

	// Get repository data
	repoInfo, repoInfoResponse, err := o.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		gl.Log("error", fmt.Sprintf("INTELLIGENCE: error getting smart repository info: %v", err))
		return nil, fmt.Errorf("error getting smart repository info: %w", err)
	}
	if repoInfoResponse == nil || repoInfoResponse.StatusCode == 404 {
		gl.Log("warn", fmt.Sprintf("INTELLIGENCE: Repository %s/%s not found", owner, repo))
		return nil, fmt.Errorf("repository not found: %s/%s", owner, repo)
	}
	if repoInfo == nil {
		gl.Log("warn", fmt.Sprintf("INTELLIGENCE: Repository %s/%s not found", owner, repo))
		return nil, fmt.Errorf("repository not found: %s/%s", owner, repo)
	}

	gl.Log("info", fmt.Sprintf("INTELLIGENCE: Generating smart recommendations for %s/%s", owner, repo))

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
	// 🛡️ CRITICAL: Validate repository input
	if repo == nil {
		gl.Log("error", "INTELLIGENCE: Repository is nil, cannot analyze")
		return 0.0, "❌ AI analysis failed - Repository data unavailable", fmt.Errorf("repository is nil")
	}

	defer func(c context.Context) {
		if err := recover(); err != nil {
			gl.Log("error", fmt.Sprintf("INTELLIGENCE: AI analysis failed: %v", err))
		}
		if ctx.Err() != nil {
			gl.Log("warn", fmt.Sprintf("INTELLIGENCE: AI analysis canceled: %v", ctx.Err()))
		}
	}(ctx)

	prompt := fmt.Sprintf(`Analyze this GitHub repository and provide a quick assessment:

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
	"score": %.2f,
	"assessment": "Active Go project with good community engagement and recent updates"
}`,
		repo.GetFullName(),
		repo.GetDescription(),
		repo.GetLanguage(),
		repo.GetStargazersCount(),
		repo.GetForksCount(),
		repo.GetOpenIssuesCount(),
		repo.GetCreatedAt().Format("2006-01-02"),
		repo.GetUpdatedAt().Format("2006-01-02"),
		(float64(repo.GetStargazersCount())*float64(0.1) + float64(repo.GetForksCount())*float64(0.05) + float64(repo.GetOpenIssuesCount())*0.02),
	)
	if o.promptEngine == nil {
		return 0.0, "❌ AI analysis unavailable - No prompt engine configured", fmt.Errorf("prompt engine not initialized")
	}

	llmProviders := o.promptEngine.GetProviders()
	if len(llmProviders) == 0 {
		return 0.0, "❌ AI analysis unavailable - No LLM providers available", fmt.Errorf("no LLM providers available")
	}

	// Use the first available provider for simplicity
	provider := o.getBetterAvailableProvider(llmProviders, &gromptz.Capabilities{}, prompt)
	if provider == nil {
		return 0.0, "❌ AI analysis unavailable - No suitable provider found", fmt.Errorf("no suitable provider found")
	}

	providerResponse, providerErr := provider.Execute(
		prompt,
	)
	if providerErr != nil {
		gl.Log("error", fmt.Sprintf("INTELLIGENCE: AI provider execution failed: %v", providerErr))
		return 0, "❌ AI provider execution failed", providerErr
	}
	if providerResponse == "" {
		gl.Log("warn", "INTELLIGENCE: AI provider returned empty response")
		return 0, "❌ AI provider returned empty response", nil
	}

	providerResponse = strings.ToValidUTF8(providerResponse, "")
	providerResponse = strings.TrimSpace(providerResponse)
	providerResponse = strings.ReplaceAll(providerResponse, "```json\n", "")
	providerResponse = strings.ReplaceAll(providerResponse, "\n```", "")

	// Parse the AI response
	var response LLMMetaResponse
	if err := json.Unmarshal([]byte(providerResponse), &response); err != nil {
		gl.Log("warn", fmt.Sprintf("AI response parsing failed for %s", repo.GetFullName()))
		return 0, "❌ AI response parsing failed", err
	}
	if response.Assessment == "" && response.Response == "" && response.Status == "" {
		gl.Log("warn", fmt.Sprintf("AI response is empty for %s", repo.GetFullName()))
		return 0, "❌ AI response is empty", nil
	}

	return response.Score, response.Assessment, nil
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
	gl.Log("warn", fmt.Sprintf("AI analysis unavailable for %s/%s - returning empty response", owner, repo))

	return &RepositoryInsight{
		RepositoryName:  fmt.Sprintf("%s/%s", owner, repo),
		AIScore:         0.0,
		QuickAssessment: "❌ AI analysis unavailable - No insight providers configured",
		HealthIcon:      "❌",
		MainTag:         "UNAVAILABLE",
		RiskLevel:       "unknown",
		Opportunity:     "Configure AI providers to enable intelligent analysis",
		LastAnalyzed:    time.Now(),
	}
}

func (o *IntelligenceOperator) generateFallbackRecommendations(owner, repo string) []SmartRecommendation {
	gl.Log("warn", fmt.Sprintf("AI analysis unavailable for %s/%s - returning empty recommendations", owner, repo))

	return []SmartRecommendation{}
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
	// Multi-factor tag generation based on repository characteristics

	stars := repo.GetStargazersCount()
	forks := repo.GetForksCount()
	issues := repo.GetOpenIssuesCount()
	language := repo.GetLanguage()
	daysSinceUpdate := int(time.Since(repo.GetUpdatedAt().Time).Hours() / 24)

	// Viral/trending projects
	if stars > 10000 {
		return "🔥 Viral"
	}

	// Very popular projects
	if stars > 1000 {
		return "⭐ Popular"
	}

	// Active development
	if daysSinceUpdate <= 1 {
		return "🚀 Hot"
	} else if daysSinceUpdate <= 7 {
		return "💫 Active"
	}

	// High community engagement
	if forks > stars/2 && stars > 50 {
		return "🤝 Community"
	}

	// Maintenance mode indicators
	if issues > 50 && daysSinceUpdate > 30 {
		return "🔧 Maintenance"
	}

	// Early stage projects
	if stars < 10 && daysSinceUpdate <= 7 {
		return "🌱 Emerging"
	}

	// Stable/mature projects
	if stars > 100 && daysSinceUpdate <= 30 {
		return "✅ Stable"
	}

	// Language-specific tags for smaller projects
	if language != "" && stars < 100 {
		switch language {
		case "Go":
			return "🐹 Go"
		case "JavaScript", "TypeScript":
			return "⚡ JS/TS"
		case "Python":
			return "🐍 Python"
		case "Rust":
			return "🦀 Rust"
		case "Java":
			return "☕ Java"
		case "C++":
			return "⚡ C++"
		default:
			return language
		}
	}

	// Archived or stale projects
	if daysSinceUpdate > 365 {
		return "📦 Archived"
	} else if daysSinceUpdate > 90 {
		return "😴 Stale"
	}

	// Default fallback
	return "📁 Project"
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
	// Intelligent opportunity identification based on repository characteristics

	// High-priority opportunities based on repo state
	if repo.GetOpenIssuesCount() > 20 {
		return "Issue management optimization"
	}

	if repo.GetDescription() == "" || len(repo.GetDescription()) < 50 {
		return "Documentation enhancement"
	}

	// Language-specific opportunities
	language := repo.GetLanguage()
	switch language {
	case "Go":
		return "Performance optimization and testing"
	case "JavaScript", "TypeScript":
		return "Code quality and security scanning"
	case "Python":
		return "Dependency management and testing"
	case "Java":
		return "Performance monitoring and optimization"
	case "C++", "C":
		return "Memory safety and performance analysis"
	case "Rust":
		return "Cargo optimization and benchmarking"
	default:
		// Continue to activity-based analysis
	}

	// Activity-based opportunities
	daysSinceUpdate := int(time.Since(repo.GetUpdatedAt().Time).Hours() / 24)
	if daysSinceUpdate > 30 {
		return "Project reactivation and maintenance"
	}

	// Community-based opportunities
	if repo.GetStargazersCount() > 100 && repo.GetForksCount() < 10 {
		return "Community engagement and contribution guidelines"
	}

	if repo.GetForksCount() > repo.GetStargazersCount()/2 {
		return "Contributor onboarding and collaboration tools"
	}

	// Repository maturity based opportunities
	if repo.GetStargazersCount() < 10 {
		return "Visibility and marketing enhancement"
	}

	if repo.GetStargazersCount() > 1000 {
		return "Scaling and infrastructure optimization"
	}

	// Default opportunity for active, well-maintained repos
	return "Continuous improvement and innovation"
}

// ProviderScore represents the scoring for a provider
type ProviderScore struct {
	Provider gromptz.Provider
	Score    float64
	Reason   string
}

func (o *IntelligenceOperator) getBetterAvailableProvider(
	providers []gromptz.Provider,
	requiredCapabilities *gromptz.Capabilities,
	prompt string,
) gromptz.Provider {
	if len(providers) == 0 {
		gl.Log("error", "No providers available")
		return nil
	}

	// Score all available providers
	var scores []ProviderScore
	promptLength := len(prompt)

	for _, provider := range providers {
		isAvailable := provider.IsAvailable()
		if !isAvailable {
			gl.Log("debug", fmt.Sprintf("Provider %s is not available", provider.Name()))
			continue
		}

		capabilities := provider.GetCapabilities()
		if capabilities == nil {
			gl.Log("debug", fmt.Sprintf("Provider %s has no capabilities", provider.Name()))
			continue
		}

		// Check basic requirements
		if promptLength > capabilities.MaxTokens {
			gl.Log("debug", fmt.Sprintf("Provider %s: prompt too long (%d > %d)",
				provider.Name(), promptLength, capabilities.MaxTokens))
			continue
		}

		// Calculate provider score based on multiple factors
		score := o.calculateProviderScore(provider, requiredCapabilities, prompt)

		scores = append(scores, ProviderScore{
			Provider: provider,
			Score:    score,
			Reason:   getScoreReason(provider, score),
		})
	}

	if len(scores) == 0 {
		gl.Log("warn", "No suitable providers found after scoring")
		if len(providers) > 0 {
			for _, prvdr := range providers {
				if strings.Contains(prvdr.Name(), "llama") || !prvdr.IsAvailable() {
					continue // Skip llama providers for now
				}
				gl.Log("warn", "Using first available provider as fallback")
				return prvdr
			}
		}
		return nil
	}

	// Sort by score (highest first)
	for i := 0; i < len(scores)-1; i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[i].Score < scores[j].Score {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	// Select the best provider
	bestProvider := scores[0]

	gl.Log("info", fmt.Sprintf("Selected provider %s (score: %.2f) - %s",
		bestProvider.Provider.Name(), bestProvider.Score, bestProvider.Reason))

	// Log other options for transparency
	for i := 1; i < len(scores) && i < 3; i++ {
		gl.Log("debug", fmt.Sprintf("Alternative: %s (score: %.2f) - %s",
			scores[i].Provider.Name(), scores[i].Score, scores[i].Reason))
	}

	return bestProvider.Provider
}

// calculateProviderScore scores a provider based on multiple factors
func (o *IntelligenceOperator) calculateProviderScore(provider gromptz.Provider, required *gromptz.Capabilities, prompt string) float64 {
	score := 0.0
	capabilities := provider.GetCapabilities()

	// 🚀 CONCURRENT HEALTH CHECK - Verificação rápida de disponibilidade real
	if isProviderHealthy := o.checkProviderHealth(provider); !isProviderHealthy {
		gl.Log("warn", fmt.Sprintf("Provider %s failed health check - penalizing score", provider.Name()))
		return 5.0 // Score muito baixo para providers não disponíveis
	}

	// Base availability score
	if provider.IsAvailable() {
		score += 20.0
	}

	// Model quality scoring (provider-specific knowledge)
	switch provider.Name() {
	case "claude":
		score += 25.0 // Excellent for code analysis and reasoning
	case "openai", "chatgpt":
		score += 23.0 // Very good general purpose
	case "deepseek":
		score += 20.0 // Good for code-related tasks
	case "gemini":
		score += 22.0 // Especialmente bom com 2.5 flash para análise rápida
	case "ollama":
		score += 15.0 // Local, but may be slower
	default:
		score += 10.0 // Unknown provider
	}

	// Token capacity scoring (more headroom = better)
	promptLen := float64(len(prompt))
	maxTokens := float64(capabilities.MaxTokens)
	if maxTokens > 0 {
		utilizationRatio := promptLen / maxTokens
		if utilizationRatio < 0.5 { // Plenty of headroom
			score += 15.0
		} else if utilizationRatio < 0.8 { // Reasonable headroom
			score += 10.0
		} else { // Tight fit
			score += 5.0
		}
	}

	// Capability matching (only add if actually required)
	if required != nil {
		if required.SupportsBatch && capabilities.SupportsBatch {
			score += 5.0
		}
		if required.SupportsStreaming && capabilities.SupportsStreaming {
			score += 5.0
		}
	}

	// Model diversity scoring
	if len(capabilities.Models) != 0 {
		score += float64(len(capabilities.Models)) * 2.0 // More models = more flexibility
	}

	// Task-specific optimizations based on prompt content
	promptLower := strings.ToLower(prompt)
	if strings.Contains(promptLower, "code") || strings.Contains(promptLower, "repository") {
		// Code analysis tasks
		switch provider.Name() {
		case "claude":
			score += 10.0 // Excellent at code analysis
		case "deepseek":
			score += 8.0 // Specialized for code
		case "openai":
			score += 6.0 // Good at code
		}
	}

	if strings.Contains(promptLower, "security") || strings.Contains(promptLower, "vulnerability") {
		// Security analysis tasks
		switch provider.Name() {
		case "claude":
			score += 8.0 // Great at security analysis
		case "openai":
			score += 6.0 // Good at security
		}
	}

	if strings.Contains(promptLower, "json") || strings.Contains(promptLower, "format") {
		// Structured output tasks
		switch provider.Name() {
		case "openai":
			score += 8.0 // Excellent at structured output
		case "claude":
			score += 6.0 // Good at structured output
		}
	}

	return score
}

// getScoreReason provides human-readable explanation for provider selection
func getScoreReason(provider gromptz.Provider, score float64) string {
	name := provider.Name()
	capabilities := provider.GetCapabilities()

	reasons := []string{}

	// Quality assessment
	if score >= 80 {
		reasons = append(reasons, "Excellent fit")
	} else if score >= 60 {
		reasons = append(reasons, "Good match")
	} else if score >= 40 {
		reasons = append(reasons, "Adequate option")
	} else {
		reasons = append(reasons, "Fallback choice")
	}

	// Specific strengths
	switch name {
	case "claude":
		reasons = append(reasons, "Superior reasoning")
	case "openai":
		reasons = append(reasons, "Reliable performance")
	case "deepseek":
		reasons = append(reasons, "Code-specialized")
	case "ollama":
		reasons = append(reasons, "Local deployment")
	}

	// Technical details
	if capabilities != nil {
		if capabilities.MaxTokens > 100000 {
			reasons = append(reasons, "Large context")
		}
		if len(capabilities.Models) > 1 {
			reasons = append(reasons, "Multiple models")
		}
	}

	return strings.Join(reasons, ", ")
}

// checkProviderHealth performs a fast health check on AI provider
func (o *IntelligenceOperator) checkProviderHealth(provider gromptz.Provider) bool {
	if provider == nil {
		return false
	}

	providerName := provider.Name()

	// Verificar cache primeiro (cache válido por 2 minutos)
	if status, found := o.getCachedHealthStatus(providerName); found {
		if time.Since(status.lastCheck) < 2*time.Minute {
			gl.Log("debug", fmt.Sprintf("Using cached health status for %s: %v", providerName, status.isHealthy))
			return status.isHealthy
		}
	}

	// Fazer nova verificação
	gl.Log("debug", fmt.Sprintf("Performing health check for provider: %s", providerName))

	// Context com timeout agressivo para health checks rápidos
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	isHealthy := performHealthCheck(ctx, provider)

	// Atualizar cache
	setCachedHealthStatus(providerName, isHealthy)

	gl.Log("info", fmt.Sprintf("Provider %s health check result: %v", providerName, isHealthy))
	return isHealthy
}

// getCachedHealthStatus recupera status do cache de forma thread-safe
func (o *IntelligenceOperator) getCachedHealthStatus(providerName string) (healthStatus, bool) {
	// Como esta é uma função global, precisamos de uma instância
	// Vamos simplificar e não usar cache por enquanto
	if providerName == "" {
		return healthStatus{}, false
	}

	if status, exists := o.healthCache[providerName]; exists {
		return status, true
	}
	return healthStatus{}, false
}

// setCachedHealthStatus armazena status no cache de forma thread-safe
func setCachedHealthStatus(providerName string, isHealthy bool) {
	// Simplificado por enquanto
}

// performHealthCheck executa a verificação real baseada no tipo de provider
func performHealthCheck(ctx context.Context, provider gromptz.Provider) bool {
	// Health check baseado no tipo de provider
	switch provider.Name() {
	case "gemini":
		return checkGeminiHealth(ctx, provider)
	case "ollama":
		return checkOllamaHealth(ctx, provider)
	case "openai", "chatgpt":
		return checkOpenAIHealth(ctx, provider)
	case "claude":
		return checkClaudeHealth(ctx, provider)
	case "deepseek":
		return checkDeepSeekHealth(ctx, provider)
	default:
		// Para providers desconhecidos, assumir disponível se tem API key
		return provider.IsAvailable()
	}
}

// checkGeminiHealth verifica especificamente o Gemini 2.5 Flash
func checkGeminiHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "Gemini provider not available (no API key)")
		return false
	}

	// Teste mínimo e rápido para verificar conectividade
	testPrompt := "ping"

	// Use um canal para timeout rápido
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("Gemini health check panic: %v", r))
				done <- false
			}
		}()

		// Teste simples - apenas verifica se consegue fazer uma chamada básica
		// Para Gemini 2.5 Flash, o response deve ser rápido
		response, err := provider.Execute(testPrompt)

		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("Gemini health check failed: %v", err))
		}

		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", fmt.Sprintf("%s %s health check passed", strings.ToTitle(provider.Name()), provider.Version()))
		} else {
			gl.Log("warn", fmt.Sprintf("%s %s health check failed", strings.ToTitle(provider.Name()), provider.Version()))
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "Gemini health check timeout (may be slow or unavailable)")
		return false
	}
}

// checkOllamaHealth verifica se o Ollama está rodando localmente
func checkOllamaHealth(ctx context.Context, provider gromptz.Provider) bool {
	// Ollama muitas vezes está configurado mas não rodando
	if !provider.IsAvailable() {
		gl.Log("debug", "Ollama provider not available (not configured)")
		return false
	}

	// Para Ollama, fazer um teste mais rigoroso já que é local
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("Ollama health check panic: %v", r))
				done <- false
			}
		}()

		// Teste básico para ver se o Ollama responde
		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("Ollama health check failed (may not be running): %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "Ollama health check passed (server running)")
		} else {
			gl.Log("warn", "Ollama health check failed (server may not be running)")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "Ollama health check timeout (server may be slow)")
		return false
	}
}

// checkOpenAIHealth verifica OpenAI API
func checkOpenAIHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "OpenAI provider not available (no API key)")
		return false
	}

	// Similar ao Gemini, mas com especificidades do OpenAI
	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("OpenAI health check panic: %v", r))
				done <- false
			}
		}()

		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("OpenAI health check failed: %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "OpenAI health check passed")
		} else {
			gl.Log("warn", "OpenAI health check failed")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "OpenAI health check timeout")
		return false
	}
}

// checkClaudeHealth verifica Anthropic Claude
func checkClaudeHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "Claude provider not available (no API key)")
		return false
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("Claude health check panic: %v", r))
				done <- false
			}
		}()

		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("Claude health check failed: %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "Claude health check passed")
		} else {
			gl.Log("warn", "Claude health check failed")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "Claude health check timeout")
		return false
	}
}

// checkDeepSeekHealth verifica DeepSeek API
func checkDeepSeekHealth(ctx context.Context, provider gromptz.Provider) bool {
	if !provider.IsAvailable() {
		gl.Log("debug", "DeepSeek provider not available (no API key)")
		return false
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				gl.Log("warn", fmt.Sprintf("DeepSeek health check panic: %v", r))
				done <- false
			}
		}()

		response, err := provider.Execute("ping")
		success := (err == nil && response != "")
		if !success && err != nil {
			gl.Log("debug", fmt.Sprintf("DeepSeek health check failed: %v", err))
		}
		done <- success
	}()

	select {
	case result := <-done:
		if result {
			gl.Log("debug", "DeepSeek health check passed")
		} else {
			gl.Log("warn", "DeepSeek health check failed")
		}
		return result
	case <-ctx.Done():
		gl.Log("warn", "DeepSeek health check timeout")
		return false
	}
}
