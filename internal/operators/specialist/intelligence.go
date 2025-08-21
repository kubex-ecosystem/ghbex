// Package specialist provides AI-powered analysis and insights for GitHub repositories.
package specialist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v61/github"

	gl "github.com/rafa-mori/ghbex/internal/module/logger"
)

// GromptRequest represents a request to Grompt server
type GromptRequest struct {
	Prompt   string                 `json:"prompt"`
	Context  map[string]interface{} `json:"context,omitempty"`
	Model    string                 `json:"model,omitempty"`
	Provider string                 `json:"provider,omitempty"`
}

// GromptResponse represents a response from Grompt server
type GromptResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
	Success  bool   `json:"success"`
}

// IntelligenceOperator provides AI-powered analysis using Grompt server
type IntelligenceOperator struct {
	client     *github.Client
	gromptURL  string
	httpClient *http.Client
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

// callGromptServer sends a request to Grompt server and returns the AI response
func (o *IntelligenceOperator) callGromptServer(prompt string, context map[string]interface{}) (string, error) {
	// Prepare request
	request := GromptRequest{
		Prompt:  prompt,
		Context: context,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Make HTTP request to Grompt server
	url := o.gromptURL + "/api/prompt"
	resp, err := o.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error calling Grompt server: %v", err)
	}
	defer resp.Body.Close()

	// Parse response
	var gromptResp GromptResponse
	if err := json.NewDecoder(resp.Body).Decode(&gromptResp); err != nil {
		return "", fmt.Errorf("error parsing grompt response: %v", err)
	}

	// Check for errors
	if !gromptResp.Success || gromptResp.Error != "" {
		return "", fmt.Errorf("grompt server error: %s", gromptResp.Error)
	}

	return gromptResp.Response, nil
}

// NewIntelligenceOperator creates a new instance that connects to Grompt server
func NewIntelligenceOperator() (*IntelligenceOperator, error) {
	// Create GitHub client
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		gl.Log("error", "GITHUB_TOKEN environment variable is not set")
		return nil, fmt.Errorf("GITHUB_TOKEN is required to connect to GitHub API")
	}

	var client *github.Client
	if token != "" {
		client = github.NewTokenClient(context.Background(), token)
	} else {
		client = github.NewClient(nil)
	}

	// Connect to Grompt server (running on port 8080)
	gromptURL := os.Getenv("GROMPT_SERVER_URL")
	if gromptURL == "" {
		gromptURL = "http://localhost:8080" // Default Grompt server URL
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Test connection to Grompt server
	testURL := gromptURL + "/health"
	resp, err := httpClient.Get(testURL)
	if err != nil {
		gl.Log("error", fmt.Sprintf("Cannot connect to Grompt server at %s: %v", gromptURL, err))
		gl.Log("info", "   Make sure Grompt server is running on port 8080")
	} else {
		resp.Body.Close()
		gl.Log("info", fmt.Sprintf("✅ Connected to Grompt server at %s", gromptURL))
	}

	return &IntelligenceOperator{
		client:     client,
		gromptURL:  gromptURL,
		httpClient: httpClient,
	}, nil
}

// GenerateQuickInsight creates AI-powered insights for repository cards
func (o *IntelligenceOperator) GenerateQuickInsight(ctx context.Context, owner, repo string) (*RepositoryInsight, error) {
	gl.Log("info", fmt.Sprintf("Generating quick insight for %s/%s", owner, repo))

	// Get basic repository info
	repoInfo, _, err := o.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return o.generateFallbackInsight(owner, repo), nil
	}

	// Generate AI-powered assessment using Grompt
	aiScore, assessment, err := o.analyzeRepositoryWithAI(ctx, repoInfo)
	if err != nil {
		gl.Log("error", fmt.Sprintf("AI analysis failed, using fallback: %v", err))
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
	gl.Log("info", fmt.Sprintf("Generating smart recommendations for %s/%s", owner, repo))

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
		gl.Log("error", fmt.Sprintf("Failed to get issues: %v", err))
	}

	// Generate AI recommendations
	recommendations, err := o.generateAIRecommendations(ctx, repoInfo, issues)
	if err != nil {
		gl.Log("error", fmt.Sprintf("AI recommendations failed, using fallback: %v", err))
		return o.generateFallbackRecommendations(owner, repo), nil
	}

	return recommendations, nil
}

// analyzeRepositoryWithAI uses Grompt to analyze repository
func (o *IntelligenceOperator) analyzeRepositoryWithAI(ctx context.Context, repo *github.Repository) (float64, string, error) {
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

	gl.Log("info", fmt.Sprintf("🤖 INTELLIGENCE: Sending prompt to Grompt server for %s", repo.GetFullName()))
	response, err := o.callGromptServer(prompt, map[string]interface{}{
		"repository": repo.GetFullName(),
		"language":   repo.GetLanguage(),
		"stars":      repo.GetStargazersCount(),
		"forks":      repo.GetForksCount(),
	})
	if err != nil {
		gl.Log("error", fmt.Sprintf("Grompt server failed: %v", err))
		gl.Log("info", "Using fallback data")
		return o.getFallbackQuickInsight(repo)
	}

	gl.Log("info", fmt.Sprintf("✅ INTELLIGENCE: AI response received: %s", response[:min(100, len(response))]))

	// Parse JSON response
	var result struct {
		Score      float64 `json:"score"`
		Assessment string  `json:"assessment"`
	}

	if err := json.Unmarshal([]byte(response), &result); err != nil {
		// Fallback if JSON parsing fails - CLEARLY MARKED AS SIMULATED
		gl.Log("error", fmt.Sprintf("Failed to parse AI response: %v", err))
		return 0.0, "⚠️  SIMULATED - AI analysis unavailable", nil
	}

	return result.Score, result.Assessment, nil
}

// generateAIRecommendations creates smart recommendations using AI
func (o *IntelligenceOperator) generateAIRecommendations(ctx context.Context, repo *github.Repository, issues []*github.Issue) ([]SmartRecommendation, error) {
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

	response, err := o.callGromptServer(prompt, map[string]interface{}{
		"repository": repo.GetFullName(),
		"language":   repo.GetLanguage(),
		"issues":     issuesContext,
	})
	if err != nil {
		gl.Log("error", fmt.Sprintf("Grompt server failed for recommendations: %v", err))
		gl.Log("info", "Using fallback recommendations")
		return o.generateFallbackRecommendations(repo.GetOwner().GetLogin(), repo.GetName()), nil
	}

	var recommendations []SmartRecommendation
	if err := json.Unmarshal([]byte(response), &recommendations); err != nil {
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

// getFallbackQuickInsight returns simulated data when AI is unavailable
func (o *IntelligenceOperator) getFallbackQuickInsight(repo *github.Repository) (float64, string, error) {
	// Return obviously simulated data with transparency warning
	assessment := fmt.Sprintf("⚠️ SIMULATED DATA - AI analysis unavailable\n\nRepository: %s\nLanguage: %s\nStars: %d\nThis is demonstration data only.",
		repo.GetFullName(),
		repo.GetLanguage(),
		repo.GetStargazersCount(),
	)

	return 0.0, assessment, nil
}

// Fallback methods for when AI is not available
func (o *IntelligenceOperator) generateFallbackInsight(owner, repo string) *RepositoryInsight {
	gl.Log("warning", fmt.Sprintf("Using SIMULATED data for %s/%s - AI analysis not available", owner, repo))

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
	gl.Log("warning", fmt.Sprintf("Using SIMULATED recommendations for %s/%s - AI analysis not available", owner, repo))

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
	// Intelligent opportunity identification based on repository characteristics

	stars := repo.GetStargazersCount()
	forks := repo.GetForksCount()
	language := repo.GetLanguage()
	hasIssues := repo.GetHasIssues()
	openIssues := repo.GetOpenIssuesCount()
	size := repo.GetSize()

	// High-value projects (many stars/forks) = community/performance focus
	if stars > 100 || forks > 50 {
		if stars > forks*3 {
			return "🌟 Community Engagement - High visibility project can benefit from contributor guides, better issue templates, and community events"
		}
		return "⚡ Performance Optimization - Popular project should focus on speed, scalability, and resource efficiency"
	}

	// Active development (many issues) = documentation/code quality
	if hasIssues && openIssues > 10 {
		if openIssues > 50 {
			return "📋 Issue Management - High issue volume suggests need for better triage, automation, and contributor onboarding"
		}
		return "📚 Documentation Enhancement - Active project with issues needs better docs, examples, and troubleshooting guides"
	}

	// Large codebase = architecture/testing focus
	if size > 10000 { // KB
		return "🏗️ Architecture Modernization - Large codebase benefits from refactoring, modularization, and technical debt reduction"
	}

	// Language-specific opportunities
	switch strings.ToLower(language) {
	case "javascript", "typescript":
		return "🔒 Security Hardening - JavaScript projects benefit from dependency auditing, CSP implementation, and security linting"
	case "python":
		return "🧪 Test Coverage Expansion - Python projects can leverage pytest, coverage analysis, and automated testing"
	case "go":
		return "📦 Go Module Optimization - Go projects benefit from dependency management, build optimization, and concurrent programming"
	case "java":
		return "🚀 Performance Tuning - Java projects can benefit from JVM optimization, memory profiling, and microservices architecture"
	case "rust":
		return "🛡️ Memory Safety Validation - Rust projects benefit from unsafe code review, fuzzing, and performance benchmarking"
	}

	// Small/new projects = foundation building
	if stars < 10 && forks < 5 {
		return "🌱 Project Foundation - New project needs README improvement, CI/CD setup, and development workflow establishment"
	}

	// Default comprehensive improvement
	return "🔄 Code Quality Boost - Repository would benefit from linting setup, automated testing, and development workflow improvements"
}
