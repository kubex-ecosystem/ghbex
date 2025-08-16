// Package artifacts provides functions to manage GitHub artifacts.
package artifacts

import (
	"context"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/defs"
	"github.com/rafa-mori/ghbex/internal/utils"
)

type IArtifacts interface {
}

func CleanArtifacts(ctx context.Context, cli *github.Client, owner, repo string, r defs.ArtifactsRule, dry bool) (deleted int, ids []int64, err error) {
	cut := utils.Cutoff(r.MaxAgeDays)
	opt := &github.ListOptions{PerPage: 100}

	for {
		arts, resp, e := cli.Actions.ListArtifacts(ctx, owner, repo, opt)
		if e != nil {
			err = e
			return
		}

		for _, a := range arts.Artifacts {
			ids = append(ids, a.GetID())
			// Delete if older than cutoff date
			if !cut.IsZero() && a.GetCreatedAt().Time.Before(cut) {
				if dry {
					deleted++
					continue
				}
				if e := deleteArtifact(ctx, cli, owner, repo, a.GetID()); e == nil {
					deleted++
				}
			}
		}

		// Check for next page
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return
}

func deleteArtifact(ctx context.Context, cli *github.Client, owner, repo string, id int64) error {
	_, err := cli.Actions.DeleteArtifact(ctx, owner, repo, id)
	return err
}
