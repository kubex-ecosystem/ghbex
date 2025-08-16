// Package sanitize provides functionalities for operating on GitHub resources.
package sanitize

import (
	"fmt"
	"time"

	"github.com/rafa-mori/ghbex/internal/defs"
)

func ToMarkdown(r *defs.Report) string {
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

## security
- ssh keys rotated: %d
- old keys removed: %d
- new key id: %d

## monitoring
- is inactive: %v
- days inactive: %d
- open PRs: %d
- open issues: %d
- commits last 30 days: %d

notes:
%v
`,
		r.Owner, r.Repo, r.When.Format(time.RFC3339), r.DryRun,
		r.Runs.Deleted, r.Runs.Kept,
		r.Artifacts.Deleted,
		r.Releases.DeletedDrafts, r.Releases.Tags,
		r.Security.SSHKeysRotated, r.Security.OldKeysRemoved, r.Security.NewKeyID,
		r.Monitoring.IsInactive, r.Monitoring.DaysInactive, r.Monitoring.OpenPRs, r.Monitoring.OpenIssues, r.Monitoring.CommitsLast30,
		r.Notes,
	)
}
