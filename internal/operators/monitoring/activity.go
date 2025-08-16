// Package monitoring provides functions to monitor GitHub repository activity.
package monitoring

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v61/github"
)

// ActivityReport represents repository activity analysis
type ActivityReport struct {
	Owner        string            `json:"owner"`
	Repo         string            `json:"repo"`
	LastActivity time.Time         `json:"last_activity"`
	IsInactive   bool              `json:"is_inactive"`
	DaysInactive int               `json:"days_inactive"`
	PRStats      *PullRequestStats `json:"pr_stats"`
	IssueStats   *IssueStats       `json:"issue_stats"`
	CommitStats  *CommitStats      `json:"commit_stats"`
}

// PullRequestStats represents PR statistics
type PullRequestStats struct {
	Open     int       `json:"open"`
	Closed   int       `json:"closed"`
	Merged   int       `json:"merged"`
	LastPR   time.Time `json:"last_pr"`
	OldestPR time.Time `json:"oldest_pr"`
}

// IssueStats represents issue statistics
type IssueStats struct {
	Open        int       `json:"open"`
	Closed      int       `json:"closed"`
	LastIssue   time.Time `json:"last_issue"`
	OldestIssue time.Time `json:"oldest_issue"`
}

// CommitStats represents commit statistics
type CommitStats struct {
	LastCommit    time.Time `json:"last_commit"`
	CommitsLast30 int       `json:"commits_last_30_days"`
	CommitsLast7  int       `json:"commits_last_7_days"`
}

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

// analyzePullRequests analyzes pull request statistics
func analyzePullRequests(ctx context.Context, cli *github.Client, owner, repo string, stats *PullRequestStats) error {
	// List all PRs (open and closed)
	allOpt := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		prs, resp, err := cli.PullRequests.List(ctx, owner, repo, allOpt)
		if err != nil {
			return err
		}

		for _, pr := range prs {
			switch pr.GetState() {
			case "open":
				stats.Open++
				if stats.OldestPR.IsZero() || pr.GetCreatedAt().Time.Before(stats.OldestPR) {
					stats.OldestPR = pr.GetCreatedAt().Time
				}
			case "closed":
				if pr.MergedAt != nil {
					stats.Merged++
				} else {
					stats.Closed++
				}
			}

			// Track latest PR activity
			prTime := pr.GetUpdatedAt().Time
			if prTime.After(stats.LastPR) {
				stats.LastPR = prTime
			}
		}

		if resp.NextPage == 0 {
			break
		}
		allOpt.Page = resp.NextPage
	}

	return nil
}

// analyzeIssues analyzes issue statistics
func analyzeIssues(ctx context.Context, cli *github.Client, owner, repo string, stats *IssueStats) error {
	// List all issues (open and closed)
	allOpt := &github.IssueListByRepoOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		issues, resp, err := cli.Issues.ListByRepo(ctx, owner, repo, allOpt)
		if err != nil {
			return err
		}

		for _, issue := range issues {
			// Skip pull requests (they appear in issues API)
			if issue.IsPullRequest() {
				continue
			}

			switch issue.GetState() {
			case "open":
				stats.Open++
				if stats.OldestIssue.IsZero() || issue.GetCreatedAt().Time.Before(stats.OldestIssue) {
					stats.OldestIssue = issue.GetCreatedAt().Time
				}
			case "closed":
				stats.Closed++
			}

			// Track latest issue activity
			issueTime := issue.GetUpdatedAt().Time
			if issueTime.After(stats.LastIssue) {
				stats.LastIssue = issueTime
			}
		}

		if resp.NextPage == 0 {
			break
		}
		allOpt.Page = resp.NextPage
	}

	return nil
}

// analyzeCommits analyzes commit statistics
func analyzeCommits(ctx context.Context, cli *github.Client, owner, repo string, stats *CommitStats) error {
	// Get commits from the last 30 days
	since := time.Now().AddDate(0, 0, -30)

	opt := &github.CommitsListOptions{
		Since:       since,
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		commits, resp, err := cli.Repositories.ListCommits(ctx, owner, repo, opt)
		if err != nil {
			return err
		}

		for _, commit := range commits {
			commitTime := commit.GetCommit().GetCommitter().GetDate().Time

			// Track latest commit
			if commitTime.After(stats.LastCommit) {
				stats.LastCommit = commitTime
			}

			// Count commits in different periods
			daysSince := time.Since(commitTime).Hours() / 24
			if daysSince <= 30 {
				stats.CommitsLast30++
			}
			if daysSince <= 7 {
				stats.CommitsLast7++
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return nil
}

// getLatestActivity returns the most recent activity timestamp
func getLatestActivity(times ...time.Time) time.Time {
	var latest time.Time
	for _, t := range times {
		if t.After(latest) {
			latest = t
		}
	}
	return latest
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
