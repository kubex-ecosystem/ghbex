// Package analytics provides advanced repository intelligence and insights.
package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"
)

// AnalyzeRepository performs comprehensive repository analysis
func AnalyzeRepository(ctx context.Context, client *github.Client, owner, repo string, analysisDays int) (*InsightsReport, error) {
	if analysisDays <= 0 {
		analysisDays = 90 // Default to 90 days
	}

	report := &InsightsReport{
		Owner:        owner,
		Repo:         repo,
		GeneratedAt:  time.Now(),
		AnalysisDays: analysisDays,
	}

	// Calculate analysis period
	since := time.Now().AddDate(0, 0, -analysisDays)

	// Analyze development patterns
	devPatterns, err := analyzeDevelopmentPatterns(ctx, client, owner, repo, since)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze development patterns: %w", err)
	}
	report.DevPatterns = devPatterns

	// Analyze code intelligence
	codeIntel, err := analyzeCodeIntelligence(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze code intelligence: %w", err)
	}
	report.CodeIntel = codeIntel

	// Calculate health score
	healthScore := calculateHealthScore(devPatterns, codeIntel)
	report.HealthScore = healthScore

	// Analyze community insights
	community, err := analyzeCommunityInsights(ctx, client, owner, repo, since)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze community insights: %w", err)
	}
	report.Community = community

	// Calculate productivity metrics
	productivity := calculateProductivityMetrics(devPatterns, community)
	report.Productivity = productivity

	// Generate recommendations
	recommendations := generateRecommendations(report)
	report.Recommendations = recommendations

	return report, nil
}

func GetRepositoryInsights(ctx context.Context, client *github.Client, owner, repo string, days int) (*InsightsReport, error) {
	var repositoryReport *InsightsReport

	// Implementation for gathering repository insights
	repositoryReport, err := AnalyzeRepository(ctx, client, owner, repo, days)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze repository: %w", err)
	}

	return &InsightsReport{
		AnalysisDays: repositoryReport.AnalysisDays,
		CodeIntel:    repositoryReport.CodeIntel,
		Productivity: repositoryReport.Productivity,
	}, nil
}
