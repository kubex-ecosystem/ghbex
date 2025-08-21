# 📋 GHbex API Endpoints Documentation

Esta documentação contém todos os endpoints disponíveis na API do GHbex com exemplos práticos usando curl.

## 🚀 Base URL

```bash
export GHBEX_URL="http://localhost:8088"
```

---

## 📊 **Health & Status**

### GET /health

Verifica o status de saúde do servidor e configurações.

**Exemplo:**

```bash
curl -X GET "$GHBEX_URL/health" | jq
```

**Resposta:**

```json
{
  "status": "ok",
  "version": "0.0.1",
  "github_auth": true,
  "config_repos": 2
}
```

---

## 📦 **Repository Management**

### GET /repos

Lista todos os repositórios configurados com insights de IA.

**Exemplo:**

```bash
curl -X GET "$GHBEX_URL/repos" | jq
```

**Resposta:**

```json
{
  "total": 2,
  "repositories": [
    {
      "owner": "rafa-mori",
      "name": "ghbex",
      "url": "https://github.com/rafa-mori/ghbex",
      "rules": {
        "runs": {
          "max_age_days": 30,
          "keep_success_last": 5
        },
        "artifacts": {
          "max_age_days": 7
        },
        "monitoring": {
          "inactive_days_threshold": 90
        }
      },
      "ai": {
        "score": 87.5,
        "assessment": "Active Go project with good community engagement",
        "health_icon": "🟢",
        "main_tag": "Active",
        "risk_level": "low",
        "opportunity": "Performance optimization"
      }
    }
  ]
}
```

---

## 🧹 **Sanitization Operations**

### POST /admin/repos/{owner}/{repo}/sanitize

Sanitiza um repositório específico (individual).

**Parâmetros:**

- `dry_run`: `true|false` (padrão: `false`)

**Exemplo (Dry Run):**

```bash
curl -X POST "$GHBEX_URL/admin/repos/rafa-mori/ghbex/sanitize?dry_run=true" | jq
```

**Exemplo (Execução Real):**

```bash
curl -X POST "$GHBEX_URL/admin/repos/rafa-mori/ghbex/sanitize?dry_run=false" | jq
```

**Resposta:**

```json
{
  "owner": "rafa-mori",
  "repo": "ghbex",
  "dry_run": true,
  "timestamp": "2025-08-21T15:30:45Z",
  "runs": {
    "total_found": 45,
    "deleted": 12,
    "kept": 33,
    "saved_space_mb": 120
  },
  "artifacts": {
    "total_found": 23,
    "deleted": 8,
    "kept": 15,
    "saved_space_mb": 400
  },
  "releases": {
    "total_drafts": 3,
    "deleted_drafts": 3
  },
  "duration_ms": 2450
}
```

### POST /admin/sanitize/bulk

Sanitiza múltiplos repositórios configurados em lote.

**Parâmetros:**

- `dry_run`: `true|false` (padrão: `false`)

**Exemplo (Dry Run):**

```bash
curl -X POST "$GHBEX_URL/admin/sanitize/bulk?dry_run=true" | jq
```

**Exemplo (Execução Real):**

```bash
curl -X POST "$GHBEX_URL/admin/sanitize/bulk?dry_run=false" | jq
```

**Resposta:**

```json
{
  "bulk_operation": true,
  "dry_run": true,
  "started_at": "2025-08-21 15:30:45",
  "duration_ms": 5230,
  "total_repos": 3,
  "total_runs_cleaned": 45,
  "total_artifacts_cleaned": 32,
  "productivity_summary": {
    "estimated_storage_saved_mb": 2050,
    "estimated_time_saved_min": 154
  },
  "repositories": [
    {
      "owner": "rafa-mori",
      "repo": "ghbex",
      "runs": 15,
      "artifacts": 12,
      "releases": 2,
      "success": true
    }
  ]
}
```

---

## 🧠 **Intelligence Operations (AI)**

### GET /intelligence/quick/{owner}/{repo}

Gera insights rápidos de IA para um repositório.

**Exemplo:**

```bash
curl -X GET "$GHBEX_URL/intelligence/quick/rafa-mori/ghbex" | jq
```

**Resposta:**

```json
{
  "ai_score": 87.5,
  "quick_assessment": "Active Go project with excellent code quality and community engagement",
  "health_icon": "🟢",
  "main_tag": "High Quality",
  "risk_level": "low",
  "opportunity": "Consider implementing automated testing workflows",
  "provider_used": "gemini-2.5-flash",
  "analysis_duration_ms": 1200,
  "confidence": 0.92
}
```

### GET /intelligence/recommendations/{owner}/{repo}

Gera recomendações inteligentes detalhadas.

**Exemplo:**

```bash
curl -X GET "$GHBEX_URL/intelligence/recommendations/rafa-mori/ghbex" | jq
```

**Resposta:**

```json
[
  {
    "category": "Security",
    "title": "Enable Dependabot Alerts",
    "description": "Configure automated dependency vulnerability scanning",
    "priority": "high",
    "effort": "low",
    "impact": "high",
    "confidence": 0.95,
    "implementation": {
      "steps": [
        "Go to repository Settings",
        "Navigate to Security & analysis",
        "Enable Dependabot alerts"
      ],
      "estimated_time_minutes": 5
    }
  }
]
```

---

## 📊 **Analytics Operations**

### GET /analytics/{owner}/{repo}

Realiza análise completa de métricas do repositório.

**Parâmetros:**

- `days`: Número de dias para análise (padrão: 90)

**Exemplo:**

```bash
curl -X GET "$GHBEX_URL/analytics/rafa-mori/ghbex?days=30" | jq
```

**Resposta:**

```json
{
  "health_score": {
    "overall": 87.5,
    "grade": "A",
    "components": {
      "activity": 92.0,
      "quality": 85.0,
      "security": 88.0,
      "community": 86.0
    }
  },
  "activity_analysis": {
    "commits_count": 156,
    "contributors_count": 3,
    "issues_opened": 8,
    "issues_closed": 6,
    "pull_requests": 12
  },
  "dependency_health": {
    "total_dependencies": 45,
    "outdated_count": 3,
    "vulnerable_count": 0,
    "health_percentage": 93.3
  },
  "analysis_period_days": 30,
  "generated_at": "2025-08-21T15:30:45Z"
}
```

---

## 🚀 **Productivity Operations**

### GET /productivity/{owner}/{repo}

Analisa produtividade e sugere otimizações.

**Exemplo:**

```bash
curl -X GET "$GHBEX_URL/productivity/rafa-mori/ghbex" | jq
```

**Resposta:**

```json
{
  "productivity_score": 82.5,
  "grade": "B+",
  "workflow_analysis": {
    "total_workflows": 4,
    "success_rate": 94.2,
    "average_duration_minutes": 3.5,
    "bottlenecks": ["test execution", "dependency installation"]
  },
  "actions": [
    {
      "type": "optimization",
      "title": "Enable workflow caching",
      "description": "Cache dependencies to reduce build time by ~40%",
      "priority": "high",
      "estimated_savings_minutes": 2.1,
      "implementation_effort": "low"
    }
  ],
  "roi": {
    "current_time_cost_hours": 12.5,
    "potential_savings_hours": 5.2,
    "roi_ratio": 2.4
  }
}
```

---

## 🤖 **Automation Operations**

### GET /automation/{owner}/{repo}

Analisa oportunidades de automação.

**Parâmetros:**

- `days`: Período de análise em dias (padrão: 30)

**Exemplo:**

```bash
curl -X GET "$GHBEX_URL/automation/rafa-mori/ghbex?days=60" | jq
```

**Resposta:**

```json
{
  "automation_score": 75.0,
  "grade": "B",
  "analysis_period_days": 60,
  "recommendations": [
    {
      "category": "CI/CD",
      "title": "Auto-merge for dependency updates",
      "description": "Automatically merge Dependabot PRs when tests pass",
      "confidence": 0.88,
      "effort": "medium",
      "impact": "high",
      "implementation": {
        "type": "github_action",
        "template": "auto-merge-dependabot.yml"
      }
    }
  ],
  "current_automation": {
    "ci_cd_coverage": 85.0,
    "testing_automation": 90.0,
    "deployment_automation": 60.0
  }
}
```

---

## 🏠 **Frontend Dashboard**

### GET /

Acessa o dashboard web integrado.

**Exemplo:**

```bash
# Abrir no navegador
open "$GHBEX_URL"

# Ou usando curl para verificar se está respondendo
curl -X GET "$GHBEX_URL" | head -20
```

---

## 🔧 **Exemplos de Uso Avançados**

### Análise Completa de Repositório

```bash
#!/bin/bash
OWNER="rafa-mori"
REPO="ghbex"

echo "🔍 Análise Completa de $OWNER/$REPO"
echo "=================================="

echo "📊 Analytics:"
curl -s "$GHBEX_URL/analytics/$OWNER/$REPO" | jq '.health_score'

echo "🧠 AI Insights:"
curl -s "$GHBEX_URL/intelligence/quick/$OWNER/$REPO" | jq '.quick_assessment'

echo "🚀 Productivity:"
curl -s "$GHBEX_URL/productivity/$OWNER/$REPO" | jq '.productivity_score'

echo "🤖 Automation:"
curl -s "$GHBEX_URL/automation/$OWNER/$REPO" | jq '.automation_score'
```

### Sanitização Segura com Verificação

```bash
#!/bin/bash
OWNER="rafa-mori"
REPO="ghbex"

echo "🧹 Sanitização Segura de $OWNER/$REPO"
echo "====================================="

echo "1. Executando dry-run..."
DRY_RESULT=$(curl -s "$GHBEX_URL/admin/repos/$OWNER/$REPO/sanitize?dry_run=true")
RUNS_TO_DELETE=$(echo $DRY_RESULT | jq '.runs.deleted')

echo "   Runs a serem deletados: $RUNS_TO_DELETE"

if [ "$RUNS_TO_DELETE" -gt 0 ]; then
    echo "2. Confirma execução? (y/N)"
    read -r confirmation
    if [ "$confirmation" = "y" ]; then
        echo "3. Executando sanitização real..."
        curl -s "$GHBEX_URL/admin/repos/$OWNER/$REPO/sanitize?dry_run=false" | jq
    else
        echo "3. Operação cancelada."
    fi
else
    echo "2. Nenhuma limpeza necessária."
fi
```

### Monitoramento de Saúde

```bash
#!/bin/bash
echo "🏥 Health Check do GHbex"
echo "======================="

HEALTH=$(curl -s "$GHBEX_URL/health")
STATUS=$(echo $HEALTH | jq -r '.status')
AUTH=$(echo $HEALTH | jq -r '.github_auth')
REPOS=$(echo $HEALTH | jq -r '.config_repos')

echo "Status: $STATUS"
echo "GitHub Auth: $AUTH"
echo "Repositórios Configurados: $REPOS"

if [ "$STATUS" = "ok" ] && [ "$AUTH" = "true" ]; then
    echo "✅ Sistema funcionando normalmente"
else
    echo "❌ Problemas detectados"
fi
```

---

## 🛠️ **Status Codes**

| Status Code | Significado |
|-------------|-------------|
| `200` | Sucesso |
| `400` | Bad Request (parâmetros inválidos) |
| `404` | Endpoint não encontrado |
| `405` | Método não permitido |
| `500` | Erro interno do servidor |

## 🔑 **Authentication**

Todos os endpoints usam a autenticação configurada no servidor (PAT ou GitHub App). Não é necessário passar tokens nos requests - a autenticação é gerenciada internamente.

## ⚡ **Rate Limiting**

O GHbex respeita automaticamente os rate limits da API GitHub e implementa retry logic inteligente para garantir operações confiáveis.
