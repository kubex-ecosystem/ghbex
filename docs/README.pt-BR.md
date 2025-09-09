# ![GHbex Banner](docs/assets/nm_banner_md.png)

<p align="center"><b>GHbex</b> — Plataforma inteligente de gestão de repositórios GitHub, com automação, análise avançada e integração multi-AI.</p>
<p align="center"><em>Automatize, otimize e monitore seus repositórios GitHub com inteligência e segurança.</em></p>

[![Kubex Go Dist CI](https://github.com/kubex-ecosystem/ghbex/actions/workflows/kubex_go_release.yml/badge.svg)](https://github.com/kubex-ecosystem/ghbex/actions/workflows/kubex_go_release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E=1.21-blue)](go.mod)
[![Releases](https://img.shields.io/github/v/release/rafa-mori/ghbex?include_prereleases)](https://github.com/kubex-ecosystem/ghbex/releases)

---

<!--
<p align="center">
  <img src="docs/assets/ghbex_demo.gif" alt="Demonstração animada do GHbex" width="80%"/>
  <br><em>GIF: Demonstração animada do GHbex (adicione aqui quando disponível)</em>
</p>
-->

## � Sumário

- [Sobre o Projeto](#-sobre-o-projeto)
- [Principais Funcionalidades](#-principais-funcionalidades)
- [Instalação](#-instalação)
- [Configuração](#️-configuração)
- [Uso Rápido](#-uso-rápido)
- [Exemplos de Uso](#-exemplos-de-uso)
- [CLI](#️-cli)
- [Arquitetura](#️-arquitetura)
- [Segurança](#-segurança)
- [Performance](#-performance)
- [Contribuição](#-contribuição)
- [Roadmap](#�️-roadmap)
- [Licença](#-licença)
- [Autor](#-autor)
- [Links Úteis](#-links-úteis)

---

## 🧩 Sobre o Projeto

O **GHbex** é uma plataforma avançada para gestão de repositórios GitHub, com recursos de inteligência artificial, automação, análise e otimização. Permite desde a sanitização e limpeza de repositórios até recomendações inteligentes, análise de dependências, automação de workflows e integração com múltiplos provedores de IA (Gemini, OpenAI, Claude, DeepSeek, Ollama).

Ideal para times DevOps, engenheiros de software e mantenedores que buscam automação, governança e insights sobre seus repositórios.

---

## ✨ Principais Funcionalidades

- **Operador de Inteligência (AI-Powered):**
  - Análise rápida e recomendações inteligentes
  - Suporte multi-provedor (Gemini, OpenAI, Claude, DeepSeek, Ollama)
  - Health check concorrente dos provedores
- **Sanitização de Repositórios:**
  - Limpeza automática de workflows antigos, artefatos e releases em draft
  - Operações em massa para múltiplos repositórios
- **Analytics & Insights:**
  - Análise de saúde, dependências, vulnerabilidades e padrões de atividade
  - Métricas de performance e engajamento
- **Otimização de Produtividade:**
  - Análise de workflows, sugestões de auto-merge, otimização de notificações
  - Cálculo de ROI
- **Motor de Automação:**
  - Reconhecimento de padrões, recomendações automáticas, scoring de confiança
  - Sugestões de integrações

---

## ⚡ Instalação

Requisitos: Go >= 1.21

```bash
# Clone o repositório
git clone https://github.com/kubex-ecosystem/ghbex.git
cd ghbex

# Instale as dependências
go mod tidy

# Build do projeto
make build-dev
```

---

## ⚙️ Configuração

```bash
# Autenticação GitHub
export GITHUB_TOKEN="ghp_seu_token_pessoal"

# Provedores de IA (opcional)
export GEMINI_API_KEY="sua_gemini_api_key"
export OPENAI_API_KEY="sua_openai_api_key"

# Notificações Discord (opcional)
export DISCORD_WEBHOOK_URL="sua_discord_webhook_url"
```

---

## 🚀 Uso Rápido

```bash
# Inicie o servidor
./dist/ghbex start --owner rafa-mori --port 8088 --repos 'owner/repo1,owner/repo2'

# Ou usando variáveis de ambiente
export REPO_LIST='owner/repo1,owner/repo2'
./dist/ghbex start --port 8088
```

### Acesso

- **Dashboard**: <http://localhost:8088>
- **Health Check**: <http://localhost:8088/health>
- **API Docs**: [docs/endpoints.md](docs/endpoints.md)

---

## � Exemplos de Uso

### Análise Inteligente de Repositório

```bash
ghbex intelligence --repo rafa-mori/ghbex
```

### Sanitização em Massa

```bash
ghbex sanitize --repos 'rafa-mori/ghbex,rafa-mori/logz'
```

### Recomendações de Automação

```bash
ghbex automation --repo rafa-mori/ghbex
```

---

## 🖥️ CLI

```bash
# Iniciar servidor
ghbex start --owner <owner> --port <port> --repos '<repo1,repo2>'

# Verificar status
ghbex status

# Parar servidor
ghbex stop

# Verificar configuração
ghbex config

# Exibir versão
ghbex version
```

---

## 🏗️ Arquitetura

```plaintext
ghbex/
├── cmd/                    # Entrypoints CLI
│   ├── main.go            # Entrypoint principal
│   └── cli/               # Comandos CLI
├── internal/              # Código interno
│   ├── operators/         # Operadores especializados
│   │   ├── intelligence/  # IA e insights
│   │   ├── sanitize/      # Limpeza
│   │   ├── analytics/     # Métricas
│   │   ├── productivity/  # Otimização
│   │   └── automation/    # Automação
│   ├── server/           # HTTP server
│   ├── client/           # GitHub client
│   └── config/           # Configuração
├── docs/                 # Documentação
└── support/              # Scripts de apoio
```

---

## � Segurança

- **Sanitização de entrada**: Validação rigorosa de parâmetros
- **Rate limiting**: Respeita limites da API do GitHub
- **Modo dry-run**: Execução segura para testes
- **Escopo restrito**: Só repositórios explicitamente configurados
- **Recuperação de erros**: Tratamento robusto de erros e panics

---

## � Performance

- **Health check concorrente** para provedores de IA
- **Timeouts agressivos** (3s) para checagens
- **Operações paralelas** para múltiplos provedores
- **Cache inteligente**: thread-safe para evitar repetições

---

## 🤝 Contribuição

1. Fork o projeto
2. Clone seu fork
3. Crie um branch (`git checkout -b feature/nova-feature`)
4. Commit suas alterações (`git commit -m 'feat: nova feature'`)
5. Push para o branch (`git push origin feature/nova-feature`)
6. Abra um Pull Request

---

## �️ Roadmap

- [ ] Swagger Docs: documentação automática da API
- [ ] Webhook Support: integração GitHub webhooks
- [ ] Analytics avançado: dashboards e métricas
- [ ] Gestão de times: permissões e times
- [ ] Operações agendadas: execuções recorrentes

---

## 📄 Licença

Projeto sob licença **MIT** — veja o arquivo [LICENSE](LICENSE).

---

## 👨‍💻 Autor

**Rafael Mori** — [@rafa-mori](https://github.com/kubex-ecosystem)

---

## 🔗 Links Úteis

- [Documentação de Endpoints](docs/endpoints.md)
- [Configuração Avançada](docs/config/)
- [Issues & Bug Reports](https://github.com/kubex-ecosystem/ghbex/issues)
- [Discussions](https://github.com/kubex-ecosystem/ghbex/discussions)
