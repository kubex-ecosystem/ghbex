// Package operators provides functionalities for operating on GitHub resources.
package operators

import (
	"fmt"
	"time"

	"github.com/rafa-mori/ghbex/internal/config"
)

func ToMarkdown(r *config.Report) string {
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
