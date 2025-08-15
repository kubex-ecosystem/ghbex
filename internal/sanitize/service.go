package sanitize

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v61/github"
)

type Notifier interface {
	Send(ctx context.Context, title, text string, files ...Attachment) error
}
type Attachment struct {
	Name string
	Body []byte
}

type Service struct {
	cli       *github.Client
	cfg       Config
	notifiers []Notifier
}

func New(cli *github.Client, cfg Config, ntf ...Notifier) *Service {
	return &Service{cli: cli, cfg: cfg, notifiers: ntf}
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

func (s *Service) SanitizeRepo(ctx context.Context, owner, repo string, rules Rules, dryRun bool) (*Report, error) {
	rpt := &Report{Owner: owner, Repo: repo, When: time.Now(), DryRun: dryRun}

	d1, k1, ids1, err := cleanRuns(ctx, s.cli, owner, repo, rules.Runs, dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "runs: "+err.Error())
	}
	rpt.Runs.Deleted, rpt.Runs.Kept, rpt.Runs.IDs = d1, k1, ids1

	d2, ids2, err := cleanArtifacts(ctx, s.cli, owner, repo, rules.Artifacts, dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "artifacts: "+err.Error())
	}
	rpt.Artifacts.Deleted, rpt.Artifacts.IDs = d2, ids2

	d3, tags, err := cleanReleases(ctx, s.cli, owner, repo, rules.Releases, dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "releases: "+err.Error())
	}
	rpt.Releases.DeletedDrafts, rpt.Releases.Tags = d3, tags

	// persist report
	dir := filepath.Join(s.cfg.Runtime.ReportDir, time.Now().Format("2006-01-02"))
	_ = os.MkdirAll(dir, 0o755)

	jb, _ := json.MarshalIndent(rpt, "", "  ")
	md := toMarkdown(rpt)
	_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s_%s.json", owner, repo)), jb, 0o644)
	_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s_%s.md", owner, repo)), []byte(md), 0o644)

	// notify
	title := fmt.Sprintf("Repo sanitize: %s/%s (dry_run=%v)", owner, repo, dryRun)
	for _, n := range s.notifiers {
		_ = n.Send(ctx, title, md,
			Attachment{Name: "report.json", Body: jb},
			Attachment{Name: "report.md", Body: []byte(md)},
		)
	}
	return rpt, nil

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
