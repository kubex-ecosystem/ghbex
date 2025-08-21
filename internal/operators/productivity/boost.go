// Package productivity provides advanced productivity enhancement tools for GitHub repositories.
package productivity

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"
)

// AnalyzeProductivity performs comprehensive productivity analysis
func AnalyzeProductivity(ctx context.Context, client *github.Client, owner, repo string) (*ProductivityReport, error) {
	report := &ProductivityReport{
		Owner:       owner,
		Repo:        repo,
		GeneratedAt: time.Now(),
	}

	// Analyze templates
	templates, err := analyzeTemplates(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze templates: %w", err)
	}
	report.Templates = templates

	// Analyze branching strategy
	branching, err := analyzeBranchingStrategy(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze branching: %w", err)
	}
	report.Branching = branching

	// Analyze auto-merge opportunities
	autoMerge, err := analyzeAutoMergeOpportunities(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze auto-merge: %w", err)
	}
	report.AutoMerge = autoMerge

	// Analyze notifications
	notifications := analyzeNotificationOptimization(ctx, client, owner, repo)
	report.Notifications = notifications

	// Analyze workflows
	workflows, err := analyzeWorkflowAutomation(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze workflows: %w", err)
	}
	report.Workflows = workflows

	// Analyze developer experience
	devex, err := analyzeDeveloperExperience(ctx, client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze developer experience: %w", err)
	}
	report.DevEx = devex

	// Generate actionable recommendations
	actions := generateProductivityActions(report)
	report.Actions = actions

	// Calculate ROI estimation
	roi := calculateROIEstimation(report)
	report.ROI = roi

	return report, nil
}
