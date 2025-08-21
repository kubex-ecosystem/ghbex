package monitoring

import "time"

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
