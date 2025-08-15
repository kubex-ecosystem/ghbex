// Package app provides utilities for cleaning up GitHub resources such as workflow runs, artifacts, and releases.
package app

import (
	"context"
	"slices"
	"time"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/state"
)

func CleanRuns(ctx context.Context, cli *github.Client, owner, repo string, r state.RunsRule, dry bool) (deleted, kept int, ids []int64, err error) {
	opt := &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{PerPage: 100}}
	cut := cutoff(r.MaxAgeDays)
	allow := func(name string) bool {
		if len(r.OnlyWorkflows) == 0 {
			return true
		}
		return slices.Contains(r.OnlyWorkflows, name)
	}

	for {
		rs, resp, e := cli.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opt)
		if e != nil {
			err = e
			return
		}
		for _, run := range rs.WorkflowRuns {
			if !allow(run.GetName()) {
				continue
			}
			ids = append(ids, run.GetID())

			// keep N latest successful
			if run.GetStatus() == "completed" && run.GetConclusion() == "success" && kept < r.KeepSuccessLast {
				kept++
				continue
			}
			// age filter
			if !cut.IsZero() && run.GetCreatedAt().Time.After(cut) {
				continue
			}
			if dry {
				deleted++
				continue
			}
			if e := deleteRun(ctx, cli, owner, repo, run.GetID()); e == nil {
				deleted++
			}
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}

		return
	}
}

func deleteRun(ctx context.Context, cli *github.Client, owner, repo string, id int64) error {
	_, err := cli.Actions.DeleteWorkflowRun(ctx, owner, repo, id)
	return err
}

func CleanArtifacts(ctx context.Context, cli *github.Client, owner, repo string, r state.ArtifactsRule, dry bool) (deleted int, ids []int64, err error) {
	cut := cutoff(r.MaxAgeDays)
	opt := &github.ListOptions{PerPage: 100}
	for {
		arts, resp, e := cli.Actions.ListArtifacts(ctx, owner, repo, opt)
		if e != nil {
			err = e
			return
		}
		for _, a := range arts.Artifacts {
			ids = append(ids, a.GetID())
			if cut.IsZero() || a.GetCreatedAt().Time.Before(cut) {
				if dry {
					deleted++
					continue
				}
				if e := deleteArtifact(ctx, cli, owner, repo, a.GetID()); e == nil {
					deleted++
				}
			}
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}
}

func deleteArtifact(ctx context.Context, cli *github.Client, owner, repo string, id int64) error {
	_, err := cli.Actions.DeleteArtifact(ctx, owner, repo, id)
	return err
}

func CleanReleases(ctx context.Context, cli *github.Client, owner, repo string, r state.ReleasesRule, dry bool) (deletedDrafts int, tags []string, err error) {
	opt := &github.ListOptions{PerPage: 100}
	for {
		rel, resp, e := cli.Repositories.ListReleases(ctx, owner, repo, opt)
		if e != nil {
			err = e
			return
		}
		for _, rr := range rel {
			if rr.GetDraft() && r.DeleteDrafts {
				if dry {
					deletedDrafts++
					continue
				}
				if e := deleteRelease(ctx, cli, owner, repo, rr.GetID()); e == nil {
					deletedDrafts++
				}
			}
			if rr.TagName != nil {
				tags = append(tags, rr.GetTagName())
			}
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}
}

func deleteRelease(ctx context.Context, cli *github.Client, owner, repo string, id int64) error {
	_, err := cli.Repositories.DeleteRelease(ctx, owner, repo, id)
	return err
}

func cutoff(days int) time.Time {
	if days <= 0 {
		return time.Time{}
	}
	return time.Now().AddDate(0, 0, -days)
}
