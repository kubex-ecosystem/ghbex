// Package monitoring provides functions to monitor GitHub repository activity.
package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"
)

// AnalyzeRepositoryActivity analyzes repository activity and generates a report
func AnalyzeRepositoryActivity(ctx context.Context, cli *github.Client, owner, repo string, inactiveDaysThreshold int) (*ActivityReport, error) {
	report := &ActivityReport{
		Owner:       owner,
		Repo:        repo,
		PRStats:     &PullRequestStats{},
		IssueStats:  &IssueStats{},
		CommitStats: &CommitStats{},
	}

	// Analyze Pull Requests
	if err := analyzePullRequests(ctx, cli, owner, repo, report.PRStats); err != nil {
		return nil, fmt.Errorf("failed to analyze pull requests: %w", err)
	}

	// Analyze Issues
	if err := analyzeIssues(ctx, cli, owner, repo, report.IssueStats); err != nil {
		return nil, fmt.Errorf("failed to analyze issues: %w", err)
	}

	// Analyze Commits
	if err := analyzeCommits(ctx, cli, owner, repo, report.CommitStats); err != nil {
		return nil, fmt.Errorf("failed to analyze commits: %w", err)
	}

	// Determine last activity and inactivity status
	report.LastActivity = getLatestActivity(report.PRStats.LastPR, report.IssueStats.LastIssue, report.CommitStats.LastCommit)

	if !report.LastActivity.IsZero() {
		report.DaysInactive = int(time.Since(report.LastActivity).Hours() / 24)
		report.IsInactive = report.DaysInactive > inactiveDaysThreshold
	}

	return report, nil
}

// CheckInactiveRepositories checks multiple repositories for inactivity
func CheckInactiveRepositories(ctx context.Context, cli *github.Client, repos []struct{ Owner, Name string }, inactiveDaysThreshold int) ([]*ActivityReport, error) {
	var reports []*ActivityReport

	for _, repo := range repos {
		report, err := AnalyzeRepositoryActivity(ctx, cli, repo.Owner, repo.Name, inactiveDaysThreshold)
		if err != nil {
			// Log error but continue with other repositories
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}
