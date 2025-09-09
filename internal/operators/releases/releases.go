// Package releases provides functions to manage GitHub releases.
package releases

import (
	"context"

	"github.com/google/go-github/v61/github"
	"github.com/kubex-ecosystem/gemx/ghbex/internal/defs/interfaces"
)

func CleanReleases(ctx context.Context, cli *github.Client, owner, repo string, r interfaces.IReleasesRule, dry bool) (deletedDrafts int, tags []string, err error) {
	opt := &github.ListOptions{PerPage: 100}

	// If the rule is to delete drafts, we need to paginate through all releases
	for {
		rel, resp, e := cli.Repositories.ListReleases(ctx, owner, repo, opt)
		if e != nil {
			err = e
			return
		}

		// If there are more pages, we need to fetch them
		for _, rr := range rel {
			if rr.GetDraft() && r.GetDeleteDrafts() {
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
		} // END OF INNER for

		// If there are more pages, we need to fetch them
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	} // END OF EXTERNAL for

	return
}
