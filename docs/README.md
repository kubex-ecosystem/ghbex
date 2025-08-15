# MCP Repo Sanitizer (Go + google/go-github)

MVP para higienizar repositórios GitHub:

- **Autenticação**: PAT *ou* GitHub App (JWT + Installation Token) — sem libs externas.
- **Ações**: apagar **workflow runs** antigos, **artifacts** velhos, **releases** em **draft**.
- **HTTP**: `POST /admin/repos/{owner}/{repo}/sanitize?dry_run=true`
- **Reports**: salva `.json` e `.md` em `_reports/YYYY-MM-DD/`.
- **Discord**: envia resumo via webhook.

> Pensado para encaixar no seu BE atual. Se já tiver API Gateway/Router, só “wirear” o `Service` e o `Client`.

## Uso rápido

```bash
# deps
go mod tidy

# configure
cp config/sanitize.yaml.example config/sanitize.yaml
# set DISCORD_WEBHOOK_URL se quiser notificação

# run
go run ./cmd/server

# testar (dry-run)
curl -X POST "http://localhost:8088/admin/repos/rafa-mori/grompt/sanitize?dry_run=true" | jq
````

## Estrutura

- `internal/githubx`: client GitHub (PAT/App), sem `ghinstallation`.
- `internal/sanitize`: regras e execução (runs/artifacts/releases).
- `internal/notify`: Discord notifier simples.
- `cmd/server`: HTTP mínimo pra você colar no seu gateway depois.

## Notas

- Para **GitHub App**, coloque o `.pem` no caminho do `config`.
- Para **GHES**, configure `base_url`/`upload_url` (ex.: `https://ghe.example/api/v3/`).
- **Dry-run** por padrão no YAML. Só desligue quando estiver confortável.

---
