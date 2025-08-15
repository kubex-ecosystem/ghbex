package operators

import (
	"context"
	"slices"

	"github.com/google/go-github/v61/github"
	"github.com/rafa-mori/ghbex/internal/defs"
	"github.com/rafa-mori/ghbex/internal/utils"
)

func CleanRuns(ctx context.Context, cli *github.Client, owner, repo string, r defs.RunsRule, dry bool) (deleted, kept int, ids []int64, err error) {
	opt := &github.ListWorkflowRunsOptions{ListOptions: github.ListOptions{PerPage: 100}}
	cut := utils.Cutoff(r.MaxAgeDays)
	allow := func(name string) bool {
		if len(r.OnlyWorkflows) == 0 {
			return true
		}
		return slices.Contains(r.OnlyWorkflows, name)
	}

	// If the rule is to delete runs, we need to paginate through all workflow runs
	for {
		// List workflow runs for the repository
		rs, resp, e := cli.Actions.ListRepositoryWorkflowRuns(ctx, owner, repo, opt)
		if e != nil {
			err = e
			return
		}

		// If there are more pages, we need to fetch them
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
		} // END OF INNER for

		// If there are more pages, we need to fetch them
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	} // END OF EXTERNAL for

	return
}

func deleteRun(ctx context.Context, cli *github.Client, owner, repo string, id int64) error {
	_, err := cli.Actions.DeleteWorkflowRun(ctx, owner, repo, id)
	return err
}
