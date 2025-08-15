// Package core provides core functionalities for the application.-
package core

import (
	"context"
	"fmt"
	"time"
)

type Notifier interface {
	Send(ctx context.Context, title, text string, files ...Attachment) error
}

type Attachment struct {
	Name string
	Body []byte
}

type Report struct {
	Owner  string    `json:"owner"`
	Repo   string    `json:"repo"`
	When   time.Time `json:"when"`
	DryRun bool      `json:"dry_run"`

	Runs struct {
		Deleted int     `json:"deleted"`
		Kept    int     `json:"kept"`
		IDs     []int64 `json:"ids"`
	} `json:"runs"`

	Artifacts struct {
		Deleted int     `json:"deleted"`
		IDs     []int64 `json:"ids"`
	} `json:"artifacts"`

	Releases struct {
		DeletedDrafts int      `json:"deleted_drafts"`
		Tags          []string `json:"tags"`
	} `json:"releases"`

	Notes []string `json:"notes"`
}

func toMarkdown(r *Report) string {
	return fmt.Sprintf(`# Sanitize %s/%s
- when: %s
- dry_run: %v

## runs
- deleted: %d
- kept(success last): %d

## artifacts
- deleted: %d

## releases
- deleted drafts: %d
- tags: %v

notes:
%v
`,
		r.Owner, r.Repo, r.When.Format(time.RFC3339), r.DryRun,
		r.Runs.Deleted, r.Runs.Kept,
		r.Artifacts.Deleted,
		r.Releases.DeletedDrafts, r.Releases.Tags,
		r.Notes,
	)
}
