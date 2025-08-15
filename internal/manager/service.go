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
	"github.com/rafa-mori/ghbex/internal/app"
	ghclient "github.com/rafa-mori/ghbex/internal/client"
	"github.com/rafa-mori/ghbex/internal/config"
	ghcfg "github.com/rafa-mori/ghbex/internal/config"
	ghntf "github.com/rafa-mori/ghbex/internal/notifiers"
	"github.com/rafa-mori/ghbex/internal/state"
)

type Service struct {
	cli       *github.Client
	cfg       ghcfg.MainConfig
	notifiers []*ghntf.Discord
}

func New(cli *github.Client, cfg ghcfg.MainConfig, ntf ...*ghntf.Discord) *Service {
	return &Service{cli: cli, cfg: cfg, notifiers: ntf}
}

func (s *Service) SanitizeRepo(ctx context.Context, owner, repo string, rules state.Rules, dryRun bool) (*ghclient.Report, error) {
	rpt := &ghclient.Report{Owner: owner, Repo: repo, When: time.Now(), DryRun: dryRun}

	d1, k1, ids1, err := app.CleanRuns(ctx, s.cli, owner, repo, rules.Runs, dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "runs: "+err.Error())
	}
	rpt.Runs.Deleted, rpt.Runs.Kept, rpt.Runs.IDs = d1, k1, ids1

	d2, ids2, err := app.CleanArtifacts(ctx, s.cli, owner, repo, rules.Artifacts, dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "artifacts: "+err.Error())
	}
	rpt.Artifacts.Deleted, rpt.Artifacts.IDs = d2, ids2

	d3, tags, err := app.CleanReleases(ctx, s.cli, owner, repo, rules.Releases, dryRun)
	if err != nil {
		rpt.Notes = append(rpt.Notes, "releases: "+err.Error())
	}
	rpt.Releases.DeletedDrafts, rpt.Releases.Tags = d3, tags

	// persist report
	dir := filepath.Join(s.cfg.GetRuntime().ReportDir, time.Now().Format("2006-01-02"))
	_ = os.MkdirAll(dir, 0o755)

	jb, _ := json.MarshalIndent(rpt, "", "  ")
	md := ghclient.ToMarkdown(rpt)
	_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s_%s.json", owner, repo)), jb, 0o644)
	_ = os.WriteFile(filepath.Join(dir, fmt.Sprintf("%s_%s.md", owner, repo)), []byte(md), 0o644)

	// notify
	title := fmt.Sprintf("Repo sanitize: %s/%s (dry_run=%v)", owner, repo, dryRun)
	for _, n := range s.notifiers {
		_ = n.Send(ctx, title, md,
			config.Attachment{Name: "report.json", Body: jb},
			config.Attachment{Name: "report.md", Body: []byte(md)},
		)
	}
	return rpt, nil

}
