package workflows

import (
	"context"

	"github.com/google/go-github/v61/github"
)

func deleteRun(ctx context.Context, cli *github.Client, owner, repo string, id int64) error {
	_, err := cli.Actions.DeleteWorkflowRun(ctx, owner, repo, id)
	return err
}
