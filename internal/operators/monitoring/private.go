package monitoring

import (
	"context"
	"time"

	"github.com/google/go-github/v61/github"
)

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
