# GHbex Repo Sanitizer (Go + google/go-github)

MVP to sanitize GitHub repositories:

- **Authentication**: PAT *or* GitHub App (JWT + Installation Token) — no external libraries.
- **Actions**: delete old **workflow runs**, outdated **artifacts**, and **draft** **releases**.
- **HTTP**: `POST /admin/repos/{owner}/{repo}/sanitize?dry_run=true`
- **Reports**: saves `.json` and `.md` in `_reports/YYYY-MM-DD/`.
- **Discord**: sends summary via webhook.

> Designed to fit into your current BE. If you already have an API Gateway/Router, just “wire” the `Service` and `Client`.

## Quick Start

```bash
# deps
go mod tidy

# configure
cp config/sanitize.yaml.example config/sanitize.yaml
# set DISCORD_WEBHOOK_URL if you want notifications

# run
go run ./cmd/server

# test (dry-run)
curl -X POST "http://localhost:8088/admin/repos/rafa-mori/grompt/sanitize?dry_run=true" | jq
```

## Structure

- `internal/githubx`: GitHub client (PAT/App), no `ghinstallation`.
- `internal/sanitize`: rules and execution (runs/artifacts/releases).
- `internal/notify`: simple Discord notifier.
- `cmd/server`: minimal HTTP to plug into your gateway later.

## Notes

- For **GitHub App**, place the `.pem` in the `config` path.
- For **GHES**, configure `base_url`/`upload_url` (e.g.: `https://ghe.example/api/v3/`).
- **Dry-run** is enabled by default in YAML. Only disable it when you’re comfortable.

---
