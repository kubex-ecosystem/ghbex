# ghbex — GitHub Hygienizer (MVP)

> Automação pragmática para limpar **workflow runs**, **artifacts** e **draft releases**, com analytics leves. Foco: **segurança**, **auditoria** e **produtividade**.

---

## ✨ Por que usar?

* **Economia de tempo e storage**: remove lixo operacional (runs/artifacts antigos).
* **Risco controlado**: `dry_run` por padrão + relatórios `.json/.md` (auditoria).
* **Plugável**: server HTTP simples; pensa-se em CLI/MCP/Discord bot.
* **Extensível**: flags bitwise para ligar/desligar estágios em runtime.

---

## 🚀 Quickstart

### 1) Pré‑requisitos

* Go >= 1.22
* Personal Access Token (PAT) com escopos **`repo`, `workflow`**

### 2) Executar

```bash
GITHUB_TOKEN=ghp_xxx \
DISCORD_WEBHOOK_URL=... \  # opcional
go run ./cmd/server
# server: http://localhost:8088
```

### 3) Primeira execução (dry‑run)

```bash
curl -s -X POST \
  "http://localhost:8088/admin/repos/<owner>/<repo>/sanitize?dry_run=true" | jq
```

### 4) Execução destrutiva (opcional)

> **Recomendado** usar um header explícito de confirmação (se ativado no build):

```bash
curl -s -X POST \
  -H 'X-Confirm: i-know-this-deletes' \
  "http://localhost:8088/admin/repos/<owner>/<repo>/sanitize?dry_run=false" | jq
```

### 5) Bulk

```bash
curl -s -X POST "http://localhost:8088/admin/sanitize/bulk?dry_run=true" | jq
```

---

## ⚙️ Configuração

Arquivo YAML (ex.: `config/sanitize.yaml`):

```yaml
runtime:
  dry_run: true
  report_dir: ./_reports
server:
  addr: ":8088"

github:
  auth:
    kind: "pat"              # "pat" | "app"
    token: "${GITHUB_TOKEN}" # se kind=pat
    app_id: 0                # se kind=app
    installation_id: 0       # se kind=app
    private_key_path: "./secrets/gh_app.pem"
    base_url: ""            # GHES: https://ghe.example/api/v3/
    upload_url: ""          # GHES: https://ghe.example/api/uploads/
  repos:
    - owner: "rafa-mori"
      name: "grompt"
      rules:
        runs:
          max_age_days: 30
          keep_success_last: 10
          only_workflows: []     # ["build.yml","release.yml"]
        artifacts:
          max_age_days: 7
        releases:
          delete_drafts: true

notifiers:
  - type: "stdout"
  - type: "discord"
    webhook: "${DISCORD_WEBHOOK_URL}"
```

> **Dica**: variáveis `${...}` são expandidas via `os.ExpandEnv`.

---

## 🧩 Feature Flags (bitwise)

Ligam/desligam estágios sem reiniciar o processo.

* `runs` → limpeza de workflow runs
* `artifacts` → limpeza de artifacts
* `releases` → deleção de draft releases
* `notify` → envio de notificações
* `report` → persistência de relatório `.md/.json`

### Formatos aceitos

* YAML/ENV/CLI: `"runs,artifacts,notify"`
* Interno (mask `uint64`): combinações de bits

### Exemplo (ENV → runtime)

```plaintext
STAGES="runs,artifacts,report"
```

---

## 🧪 Exemplos de requisição/resposta

### Sanitize (repo)

***Request***

```plaintext
POST /admin/repos/rafa-mori/byte_sleuth/sanitize?dry_run=false
```

***Response (resumo)***

```json
{
  "owner": "rafa-mori",
  "repo": "byte_sleuth",
  "dry_run": false,
  "runs": {"deleted": 0, "kept": 1, "ids": [15105979084]},
  "artifacts": {"deleted": 2, "ids": [3149209461, 3149207627]},
  "releases": {"deleted_drafts": 0, "tags": ["v1.0.5"]}
}
```

### Sanitize (bulk)

```json
{
  "bulk_operation": true,
  "dry_run": false,
  "duration_ms": 18293,
  "total_repos": 10,
  "total_runs_cleaned": 14,
  "repositories": [ {"owner":"rafa-mori","repo":"xtui","runs":6,"success":true}, ... ]
}
```

### Analytics

```plaintext
GET /analytics/<owner>/<repo>?days=90
```

***Response (recorte)***

```json
{
  "owner": "rafa-mori",
  "repo": "xtui",
  "analysis_days": 90,
  "code_intelligence": {
    "languages": {"Go": 81.0, "Shell": 16.7, "Makefile": 2.1},
    "complexity": {"cyclomatic_complexity": 2.59, "technical_debt": "medium"}
  },
  "health_score": {"overall": 45.76, "grade": "F"},
  "recommendations": [
    "🔧 Melhorar qualidade de código e documentação",
    "👥 Atrair mais contribuidores",
    "⚡ Acelerar deploy e reduzir lead time"
  ]
}
```

---

## 🔒 Segurança

* **Dry‑run** por padrão; habilite destrutivo conscientemente.
* Sugestão: exigir header `X-Confirm: i-know-this-deletes` quando `dry_run=false`.
* Tokens: prefira **GitHub App** (tokens de curta duração) em produção.

---

## 📈 Limites & Desempenho

* Paginação feita página‑a‑página (100 itens): evita quedas silenciosas.
* Resiliência a rate‑limit: use backoff baseado em `X-RateLimit-Remaining/Reset`.
* Bulk: serializa por recurso, paraleliza por repositório (moderado).

---

## 🧰 Troubleshooting

* **403/404**: cheque escopos do token (`repo`, `workflow`) e org policies.
* **Secondary rate‑limit**: reduza paralelismo ou habilite backoff.
* **Nada foi deletado**: verifique `keep_success_last`, janela `max_age_days` e filtros `only_workflows`.

---

## 🗺️ Roadmap curto

* Stale PR/Issues (rótulos + cutoff de inatividade)
* Backoff nativo de rate‑limit
* Upload de relatório `.md` como release asset ou Gist
* CLI com `pflag`/`cobra`

---

## 🤝 Contribuição

1. Fork → branch → PR com descrição objetiva
2. Adicione testes e atualize este README
3. Evite breaking changes fora de `v0.x`

Licença: MIT
