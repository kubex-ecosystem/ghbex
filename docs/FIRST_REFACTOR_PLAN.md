# Roteiro objetivo (em commits pequenos)

1. **Tipos centrais + contratos** (diretório `internal/defs`): mover/normalizar tipos e erros.
2. **Operadores com responsabilidades claras**: `collector → normalizer → scorer → recommender → actions`.
3. **Middlewares (wrappers/observers)** universais: `WithMeter`, `WithBudget`, `WithRetry`, `WithCache`, `WithTimeout`.
4. **Paginação & clients**: helper genérico e DI do `GitHubAPI` (e futuramente LLM).
5. **Mapa de mocks/placeholders/TODOs** com taxonomia e CI que falha se passarem do limite.
6. **Relatório de custo/latência** (CLI/endpoint) para comprovar ganho imediato.

---

## 1) Tipos centrais (consolidação)

**`internal/defs/types.go`** (exemplo enxuto e reaproveitável):

```go
package defs

type RepoRef struct {
    Owner string
    Name  string
    Head  string // commit SHA
}

type FileRef struct {
    Path string
    SHA  string
    Lang string
    Size int64
}

type Metric struct {
    Name   string
    Value  float64
    Unit   string
    Labels map[string]string
}

type Insight struct {
    Key     string
    Summary string
    Details map[string]any
    Score   float64 // 0..1
}

type OpInput struct {
    Repo   RepoRef
    Files  []FileRef
    Params map[string]any
}

type OpOutput struct {
    Metrics   []Metric
    Insights  []Insight
    Artifacts map[string][]byte // e.g. JSON reports
}
```

**Erros sentinela** (evitar if/else mágicos):

```go
var (
    ErrBudgetExceeded = errors.New("budget exceeded")
    ErrRateLimited    = errors.New("rate limited")
)
```

---

## 2) Operadores: separar responsabilidades

Estruture por **fase**, não por “tema”:

```
operators/
  collector/     // só I/O (GitHub, LLM, FS)
  normalizer/    // modela dados crus → modelos internos
  scorer/        // cálculos e pesos (sem rede)
  recommend/     // gera recomendações textuais a partir de métricas/insights
  actions/       // side-effects (rótulos, limpar artifacts, abrir issue...)
  common/        // utilitários compartilhados (paginador, filtros, etc.)
```

Cada operador implementa **um** contrato simples:

```go
type Operator interface {
    Name() string
    Run(ctx context.Context, in *defs.OpInput) (*defs.OpOutput, error)
}
```

---

## 3) Middlewares (wrappers/observers/chains)

Componha operadores como no padrão middleware. Você ganha **métrica, retry, orçamento e cache** sem poluir regra de negócio.

```go
type Middleware func(Operator) Operator

func Chain(op Operator, mws ...Middleware) Operator {
    for i := len(mws) - 1; i >= 0; i-- { op = mws[i](op) }
    return op
}
```

### WithMeter (latência, custo, tokens, cache-hit)

```go
func WithMeter(recorder func(map[string]any)) Middleware {
    return func(next Operator) Operator {
        return operatorFunc(func(ctx context.Context, in *defs.OpInput) (*defs.OpOutput, error) {
            start := time.Now()
            out, err := next.Run(ctx, in)
            recorder(map[string]any{
                "op": next.Name(), "dur_ms": time.Since(start).Milliseconds(),
                "repo": in.Repo, "err": fmt.Sprintf("%v", err),
            })
            return out, err
        })
    }
}
```

### WithBudget (hard-stop por consumo)

```go
type Budget interface{
    Allow(cost float64) bool   // thread-safe
    Charge(cost float64)       // acumula custo
    Left() float64
}

func WithBudget(b Budget) Middleware {
    return func(next Operator) Operator {
        return operatorFunc(func(ctx context.Context, in *defs.OpInput) (*defs.OpOutput, error) {
            if !b.Allow(0) { return nil, ErrBudgetExceeded }
            out, err := next.Run(ctx, in)
            // dica: some com custo real vindo do client/LLM via context
            return out, err
        })
    }
}
```

### WithRetry + WithTimeout + WithCache

* **Retry**: exponencial com jitter; não repetir em 4xx definitivos.
* **Timeout**: `context.WithTimeout` por operador.
* **Cache**: chave **determinística** (veja promptKey no canvas) por `opName|stepVersion|repo.Head|file.SHA|params`.

---

## 4) Clients + paginação (universal)

**Contrato fino** para mock e VCR:

```go
type GitHubAPI interface {
    ListIssues(ctx context.Context, owner, repo string, opts *github.IssueListByRepoOptions) ([]*github.Issue, *github.Response, error)
    ListPRs(ctx context.Context, owner, repo string, opts *github.PullRequestListOptions) ([]*github.PullRequest, *github.Response, error)
    ListWorkflowRuns(ctx context.Context, owner, repo string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error)
    // ...somente o que você usa
}
```

**Paginador genérico** (aplique em Issues/PRs/Workflows):

```go
func Paginate[T any](ctx context.Context, f func(page int) ([]T, *github.Response, error)) ([]T, error) {
    var all []T; page := 1
    for {
        items, resp, err := f(page)
        if rl, ok := err.(*github.RateLimitError); ok {
            time.Sleep(time.Until(rl.Rate.Reset.Time)); continue
        }
        if err != nil { return all, err }
        all = append(all, items...)
        if resp.GetNextPage() == 0 { break }
        page = resp.GetNextPage()
    }
    return all, nil
}
```

---

## 5) Mocks, pseudo-mocks & TODOs: taxonomia e guard-rail

**Tags padronizadas** (sempre com “issue-link” quando real):

* `TODO(TECHDEBT):` débito técnico aceito temporariamente.
* `MOCK(VCR):` usa gravação/replay real com TTL.
* `PLACEHOLDER:` valor/escore arbitrário até implementar métrica real.
* `FIXME:` bug conhecido com prazo.

**Script CI** (falha build se exceder limites):

```bash
#!/usr/bin/env bash
set -euo pipefail
MAX_PLACEHOLDERS=20
COUNT=$(git grep -E "PLACEHOLDER:|TODO\(TECHDEBT\)|MOCK\(VCR\)|FIXME:" -- '*.go' | wc -l)
if [ "$COUNT" -gt "$MAX_PLACEHOLDERS" ]; then
  echo "Too many placeholders/TODOs: $COUNT (max $MAX_PLACEHOLDERS)"; exit 1
fi
```

**VCR Real** (sem autoengano):

* Grave respostas **reais** (GitHub/LLM) com redaction.
* Metadata por gravação: `model`, `params`, `prompt_hash`, `rate_limit`, `latência`, `custo`.
* TTL curto (ex.: 7 dias) para atualizar a realidade.

---

## 6) Métrica real de custo/latência (minimamente viável)

* **LLM**: `prompt_tokens`, `completion_tokens`, `latency_ms`, `cost_usd`, `cache_hit`.

  * Preços por modelo em tabela local; se o provider retornar `usage`, preferir fonte primária.
* **GitHub**: `http_status`, `latency_ms`, `rate_remaining`, `retry_after`, `requests_total`.
* **Agregações**: por run/model/provider/dia + “top arquivos caros” e “cache-hit-rate”.
* **Budget**: `--budget 0.50` no CLI aborta graciosamente ao ultrapassar.

---

## 7) Exemplo de operador real “limpo” (collector → normalizer)

```go
// collector/issues.go
type IssuesCollector struct{ GH GitHubAPI }
func (c IssuesCollector) Name() string { return "collector.issues" }
func (c IssuesCollector) Run(ctx context.Context, in *defs.OpInput) (*defs.OpOutput, error) {
    items, err := Paginate(ctx, func(p int) ([]*github.Issue, *github.Response, error) {
        return c.GH.ListIssues(ctx, in.Repo.Owner, in.Repo.Name, &github.IssueListByRepoOptions{
            State: "all", ListOptions: github.ListOptions{PerPage: 100, Page: p},
        })
    })
    if err != nil { return nil, err }
    // normalize: mapeia para modelos internos, sem scoring aqui
    var metrics []defs.Metric
    metrics = append(metrics, defs.Metric{Name:"issues_total", Value: float64(len(items)), Unit:"count"})
    return &defs.OpOutput{Metrics: metrics}, nil
}
```

Encadeando com observabilidade e orçamento:

```go
op := Chain(IssuesCollector{GH: ghClient},
    WithTimeout(15*time.Second),
    WithRetry(3),
    WithBudget(budget),
    WithMeter(logzRecorder),
)
out, err := op.Run(ctx, &defs.OpInput{Repo: repo})
```

---

## 8) Aceite (Definition of Done) deste refactor

* **Operadores** não possuem chamadas diretas de rede; tudo passa por **clients** injetados.
* **Paginação** presente em Issues/PRs/Workflows.
* **Middlewares** aplicáveis a qualquer operador (meter/budget/retry/cache/timeout).
* **Tipos centrais únicos** em `defs/`.
* **Mapa de mocks/TODOs** visível e **CI barrando** excesso.
* **Relatório** `usage/cost` por run; cache-hit ≥ **70%** em re-execução do mesmo commit.

---

## 9) Quick wins (atacáveis já)

* Trocar `GetRepositoryInsights(ctx, nil, …)` por DI de client (evita crash).
* Introduzir `Paginate` e aplicar nos três pontos mais caros (Issues/PRs/Workflows).
* Implementar `WithMeter` gravando no seu **logz** (JSON estruturado) e já gerar um `report cost`.
* Criar `Budget` simples (memória) e wirear no CLI (`--budget`).
