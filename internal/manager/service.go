// Package manager provides functionality for managing GitHub repositories.
package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/defs/common"
	"github.com/rafa-mori/ghbex/internal/defs/gitz"
	"github.com/rafa-mori/ghbex/internal/interfaces"
	"github.com/rafa-mori/ghbex/internal/notifiers"
	artifacts "github.com/rafa-mori/ghbex/internal/operators/artifacts"
	releases "github.com/rafa-mori/ghbex/internal/operators/releases"
	sanitize "github.com/rafa-mori/ghbex/internal/operators/sanitize"
	workflows "github.com/rafa-mori/ghbex/internal/operators/workflows"
)

type Service struct {
	cli       *github.Client
	cfg       interfaces.IMainConfig
	notifiers []*notifiers.Discord
}

func New(cli *github.Client, cfg interfaces.IMainConfig, ntf ...*notifiers.Discord) *Service {
	return &Service{cli: cli, cfg: cfg, notifiers: ntf}
}

func (s *Service) SanitizeRepo(ctx context.Context, owner, repo string, rules interfaces.IRules, dryRun bool) (*gitz.Report, error) {
	rpt := &gitz.Report{Owner: owner, Repo: repo, When: time.Now(), DryRun: dryRun}

	d1, k1, ids1, err := workflows.CleanRuns(ctx, s.cli, owner, repo, rules.GetRunsRule(), dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "runs: "+err.Error())
	}
	rpt.Runs.Deleted, rpt.Runs.Kept, rpt.Runs.IDs = d1, k1, ids1

	d2, ids2, err := artifacts.CleanArtifacts(ctx, s.cli, owner, repo, rules.GetArtifactsRule(), dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "artifacts: "+err.Error())
	}
	rpt.Artifacts.Deleted, rpt.Artifacts.IDs = d2, ids2

	d3, tags, err := releases.CleanReleases(ctx, s.cli, owner, repo, rules.GetReleasesRule(), dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "releases: "+err.Error())
	}
	rpt.Releases.DeletedDrafts, rpt.Releases.Tags = d3, tags

	// persist report
	dir := filepath.Join(s.cfg.GetRuntime().GetReportDir(), time.Now().Format("2006-01-02"))
	_ = os.MkdirAll(dir, 0o755)

	jb, _ := json.MarshalIndent(rpt, "", "  ")
	md := sanitize.ToMarkdown(rpt)
	_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s_%s.json", owner, repo)), jb, 0o644)
	_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s_%s.md", owner, repo)), []byte(md), 0o644)

	// notify
	title := fmt.Sprintf("Repo sanitize: %s/%s (dry_run=%v)", owner, repo, dryRun)
	for _, n := range s.notifiers {
		_ = n.Send(ctx, title, md,
			common.NewAttachment("report.json", jb),
			common.NewAttachment("report.md", []byte(md)),
		)
	}
	return rpt, nil

}
