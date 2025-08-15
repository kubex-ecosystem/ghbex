package operators

import (
	"context"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/defs"
)

func CleanReleases(ctx context.Context, cli *github.Client, owner, repo string, r defs.ReleasesRule, dry bool) (deletedDrafts int, tags []string, err error) {
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
