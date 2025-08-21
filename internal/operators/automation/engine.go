// Package automation provides intelligent repository automation capabilities.
package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"
)

// AnalyzeAutomation performs comprehensive automation analysis on a repository.
func AnalyzeAutomation(ctx context.Context, client *github.Client, owner, repo string, analysisDays int) (*AutomationReport, error) {
	if analysisDays <= 0 {
		analysisDays = 90
	}

	report := &AutomationReport{
		Repository:   fmt.Sprintf("%s/%s", owner, repo),
		Timestamp:    time.Now(),
		AnalysisDays: analysisDays,
	}

	// Analyze labels
	if err := analyzeLabelManagement(ctx, client, owner, repo, report); err != nil {
		return nil, fmt.Errorf("label analysis failed: %w", err)
	}

	// Analyze issues
	if err := analyzeIssueManagement(ctx, client, owner, repo, report, analysisDays); err != nil {
		return nil, fmt.Errorf("issue analysis failed: %w", err)
	}

	// Analyze PRs
	if err := analyzePRManagement(ctx, client, owner, repo, report, analysisDays); err != nil {
		return nil, fmt.Errorf("PR analysis failed: %w", err)
	}

	// Analyze workflows
	if err := analyzeWorkflowManagement(ctx, client, owner, repo, report); err != nil {
		return nil, fmt.Errorf("workflow analysis failed: %w", err)
	}

	// Generate recommendations
	generateAutomationRecommendations(report)

	// Calculate automation score and grade
	report.AutomationScore = calculateAutomationScore(report)
	report.Grade = calculateGrade(report.AutomationScore)

	return report, nil
}
