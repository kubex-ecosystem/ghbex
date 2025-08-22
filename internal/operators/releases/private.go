package releases

import (
	"context"

	"github.com/google/go-github/v61/github"
)

func deleteRelease(ctx context.Context, cli *github.Client, owner, repo string, id int64) error {
	_, err := cli.Repositories.DeleteRelease(ctx, owner, repo, id)
	return err
}
