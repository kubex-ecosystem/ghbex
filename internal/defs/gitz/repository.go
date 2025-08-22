package gitz

import "github.com/google/go-github/v61/github"

type Repository struct {
	ID                        int64                `json:"id,omitempty"`
	NodeID                    string               `json:"node_id,omitempty"`
	Owner                     string               `json:"owner,omitempty"`
	Name                      string               `json:"name,omitempty"`
	FullName                  string               `json:"full_name,omitempty"`
	Description               string               `json:"description,omitempty"`
	Homepage                  string               `json:"homepage,omitempty"`
	CodeOfConduct             github.CodeOfConduct `json:"code_of_conduct,omitempty"`
	DefaultBranch             string               `json:"default_branch,omitempty"`
	MasterBranch              string               `json:"master_branch,omitempty"`
	CreatedAt                 github.Timestamp     `json:"created_at,omitempty"`
	PushedAt                  github.Timestamp     `json:"pushed_at,omitempty"`
	UpdatedAt                 github.Timestamp     `json:"updated_at,omitempty"`
	HTMLURL                   string               `json:"html_url,omitempty"`
	CloneURL                  string               `json:"clone_url,omitempty"`
	GitURL                    string               `json:"git_url,omitempty"`
	MirrorURL                 string               `json:"mirror_url,omitempty"`
	SSHURL                    string               `json:"ssh_url,omitempty"`
	SVNURL                    string               `json:"svn_url,omitempty"`
	Language                  string               `json:"language,omitempty"`
	Fork                      bool                 `json:"fork,omitempty"`
	ForksCount                int                  `json:"forks_count,omitempty"`
	NetworkCount              int                  `json:"network_count,omitempty"`
	OpenIssuesCount           int                  `json:"open_issues_count,omitempty"`
	OpenIssues                int                  `json:"open_issues,omitempty"` // Deprecated: Replaced by OpenIssuesCount. For backward compatibility OpenIssues is still populated.
	StargazersCount           int                  `json:"stargazers_count,omitempty"`
	SubscribersCount          int                  `json:"subscribers_count,omitempty"`
	WatchersCount             int                  `json:"watchers_count,omitempty"` // Deprecated: Replaced by StargazersCount. For backward compatibility WatchersCount is still populated.
	Watchers                  int                  `json:"watchers,omitempty"`       // Deprecated: Replaced by StargazersCount. For backward compatibility Watchers is still populated.
	Size                      int                  `json:"size,omitempty"`
	AutoInit                  bool                 `json:"auto_init,omitempty"`
	Parent                    *Repository          `json:"parent,omitempty"`
	Source                    *Repository          `json:"source,omitempty"`
	TemplateRepository        *Repository          `json:"template_repository,omitempty"`
	Organization              *github.Organization `json:"organization,omitempty"`
	Permissions               map[string]bool      `json:"permissions,omitempty"`
	AllowRebaseMerge          bool                 `json:"allow_rebase_merge,omitempty"`
	AllowUpdateBranch         bool                 `json:"allow_update_branch,omitempty"`
	AllowSquashMerge          bool                 `json:"allow_squash_merge,omitempty"`
	AllowMergeCommit          bool                 `json:"allow_merge_commit,omitempty"`
	AllowAutoMerge            bool                 `json:"allow_auto_merge,omitempty"`
	AllowForking              bool                 `json:"allow_forking,omitempty"`
	WebCommitSignoffRequired  bool                 `json:"web_commit_signoff_required,omitempty"`
	DeleteBranchOnMerge       bool                 `json:"delete_branch_on_merge,omitempty"`
	UseSquashPRTitleAsDefault bool                 `json:"use_squash_pr_title_as_default,omitempty"`
	SquashMergeCommitTitle    string               `json:"squash_merge_commit_title,omitempty"`   // Can be one of: "PR_TITLE", "COMMIT_OR_PR_TITLE"
	SquashMergeCommitMessage  string               `json:"squash_merge_commit_message,omitempty"` // Can be one of: "PR_BODY", "COMMIT_MESSAGES", "BLANK"
	MergeCommitTitle          string               `json:"merge_commit_title,omitempty"`          // Can be one of: "PR_TITLE", "MERGE_MESSAGE"
	MergeCommitMessage        string               `json:"merge_commit_message,omitempty"`        // Can be one of: "PR_BODY", "PR_TITLE", "BLANK"
	Topics                    []string             `json:"topics,omitempty"`
	CustomProperties          map[string]string    `json:"custom_properties,omitempty"`
	Archived                  bool                 `json:"archived,omitempty"`
	Disabled                  bool                 `json:"disabled,omitempty"`

	// Only provided when using RepositoriesService.Get while in preview
	License *github.License `json:"license,omitempty"`

	// Additional mutable fields when creating and editing a repository
	Private           bool   `json:"private,omitempty"`
	HasIssues         bool   `json:"has_issues,omitempty"`
	HasWiki           bool   `json:"has_wiki,omitempty"`
	HasPages          bool   `json:"has_pages,omitempty"`
	HasProjects       bool   `json:"has_projects,omitempty"`
	HasDownloads      bool   `json:"has_downloads,omitempty"`
	HasDiscussions    bool   `json:"has_discussions,omitempty"`
	IsTemplate        bool   `json:"is_template,omitempty"`
	LicenseTemplate   string `json:"license_template,omitempty"`
	GitignoreTemplate string `json:"gitignore_template,omitempty"`

	// Options for configuring Advanced Security and Secret Scanning
	SecurityAndAnalysis *github.SecurityAndAnalysis `json:"security_and_analysis,omitempty"`

	// Creating an organization repository. Required for non-owners.
	TeamID *int64 `json:"team_id,omitempty"`

	// API URLs
	URL              string              `json:"url,omitempty"`
	ArchiveURL       string              `json:"archive_url,omitempty"`
	AssigneesURL     string              `json:"assignees_url,omitempty"`
	BlobsURL         string              `json:"blobs_url,omitempty"`
	BranchesURL      string              `json:"branches_url,omitempty"`
	CollaboratorsURL string              `json:"collaborators_url,omitempty"`
	CommentsURL      string              `json:"comments_url,omitempty"`
	CommitsURL       string              `json:"commits_url,omitempty"`
	CompareURL       string              `json:"compare_url,omitempty"`
	ContentsURL      string              `json:"contents_url,omitempty"`
	ContributorsURL  string              `json:"contributors_url,omitempty"`
	DeploymentsURL   string              `json:"deployments_url,omitempty"`
	DownloadsURL     string              `json:"downloads_url,omitempty"`
	EventsURL        string              `json:"events_url,omitempty"`
	ForksURL         string              `json:"forks_url,omitempty"`
	GitCommitsURL    string              `json:"git_commits_url,omitempty"`
	GitRefsURL       string              `json:"git_refs_url,omitempty"`
	GitTagsURL       string              `json:"git_tags_url,omitempty"`
	HooksURL         string              `json:"hooks_url,omitempty"`
	IssueCommentURL  string              `json:"issue_comment_url,omitempty"`
	IssueEventsURL   string              `json:"issue_events_url,omitempty"`
	IssuesURL        string              `json:"issues_url,omitempty"`
	KeysURL          string              `json:"keys_url,omitempty"`
	LabelsURL        string              `json:"labels_url,omitempty"`
	LanguagesURL     string              `json:"languages_url,omitempty"`
	MergesURL        string              `json:"merges_url,omitempty"`
	MilestonesURL    string              `json:"milestones_url,omitempty"`
	NotificationsURL string              `json:"notifications_url,omitempty"`
	PullsURL         string              `json:"pulls_url,omitempty"`
	ReleasesURL      string              `json:"releases_url,omitempty"`
	StargazersURL    string              `json:"stargazers_url,omitempty"`
	StatusesURL      string              `json:"statuses_url,omitempty"`
	SubscribersURL   string              `json:"subscribers_url,omitempty"`
	SubscriptionURL  string              `json:"subscription_url,omitempty"`
	TagsURL          string              `json:"tags_url,omitempty"`
	TreesURL         string              `json:"trees_url,omitempty"`
	TeamsURL         string              `json:"teams_url,omitempty"`
	TextMatches      []*github.TextMatch `json:"text_matches,omitempty"`
	Visibility       string              `json:"visibility,omitempty"`
	RoleName         string              `json:"role_name,omitempty"`
}
