# 🚀 GHbex - Intelligent GitHub Repository Management Platform

**GHbex** é uma plataforma avançada de gerenciamento de repositórios GitHub com capacidades de IA, análise inteligente e automação. Projetado para sanitização, otimização e insights avançados de repositórios.

## ✨ Funcionalidades Principais

### 🧠 **Intelligence Operator (AI-Powered)**

- **Quick Insights**: Análise rápida de repositórios com IA
- **Smart Recommendations**: Recomendações inteligentes baseadas em padrões
- **Multi-Provider Support**: Gemini, OpenAI, Claude, DeepSeek, Ollama
- **Health Checking**: Verificação concorrente de disponibilidade de providers

### 🧹 **Repository Sanitization**

- **Workflow Runs**: Limpeza automática de execuções antigas
- **Artifacts**: Remoção de artefatos obsoletos
- **Draft Releases**: Eliminação de releases em rascunho
- **Bulk Operations**: Sanitização em massa de múltiplos repositórios

### 📊 **Analytics & Insights**

- **Repository Health**: Análise completa de saúde do repositório
- **Dependency Analysis**: Avaliação de dependências e vulnerabilidades
- **Activity Patterns**: Padrões de atividade e engajamento
- **Performance Metrics**: Métricas de performance e qualidade

### 🚀 **Productivity Optimization**

- **Workflow Analysis**: Análise de workflows do GitHub Actions
- **Auto-merge Rules**: Recomendações de regras de auto-merge
- **Notification Optimization**: Otimização de notificações
- **ROI Calculation**: Cálculo de retorno sobre investimento

### 🤖 **Automation Engine**

- **Pattern Recognition**: Reconhecimento de padrões de automação
- **Recommendation Engine**: Motor de recomendações automáticas
- **Confidence Scoring**: Sistema de pontuação de confiança
- **Integration Suggestions**: Sugestões de integração

## 🛠️ Tecnologias

- **Backend**: Go 1.25+ com arquitetura modular
- **GitHub Integration**: API v4 com suporte a PAT e GitHub Apps
- **AI Providers**: Suporte multi-provider com health checking
- **Authentication**: PAT (Personal Access Token) ou GitHub App (JWT + Installation)
- **Notificações**: Discord webhook integration
- **Frontend**: Dashboard web integrado

## ⚡ Quick Start

### 1. **Instalação**

```bash
# Clone o repositório
git clone https://github.com/rafa-mori/ghbex.git
cd ghbex

# Instale dependências
go mod tidy

# Build do projeto
make build-dev
```

### 2. **Configuração**

```bash
# Configure GitHub authentication
export GITHUB_TOKEN="ghp_your_personal_access_token"

# Configure AI providers (opcional)
export GEMINI_API_KEY="your_gemini_api_key"
export OPENAI_API_KEY="your_openai_api_key"

# Configure Discord notifications (opcional)
export DISCORD_WEBHOOK_URL="your_discord_webhook_url"
```

### 3. **Execução**

```bash
# Inicie o servidor
./dist/ghbex start --owner rafa-mori --port 8088 --repos 'owner/repo1,owner/repo2'

# Ou usando variáveis de ambiente
export REPO_LIST='owner/repo1,owner/repo2'
./dist/ghbex start --port 8088
```

### 4. **Acesso**

- **Dashboard**: <http://localhost:8088>
- **Health Check**: <http://localhost:8088/health>
- **API Documentation**: Veja [docs/endpoints.md](docs/endpoints.md)

## 🔧 Comandos CLI

```bash
# Iniciar servidor
ghbex start --owner <owner> --port <port> --repos '<repo1,repo2>'

# Verificar status
ghbex status

# Parar servidor
ghbex stop

# Verificar configuração
ghbex config

# Mostrar versão
ghbex version
```

## 🏗️ Arquitetura

```plaintext
ghbex/
├── cmd/                    # Pontos de entrada CLI
│   ├── main.go            # Entrypoint principal
│   └── cli/               # Comandos CLI
├── internal/              # Código interno
│   ├── operators/         # Operadores especializados
│   │   ├── intelligence/  # IA e insights
│   │   ├── sanitize/      # Limpeza de repositórios
│   │   ├── analytics/     # Análise e métricas
│   │   ├── productivity/  # Otimização de produtividade
│   │   └── automation/    # Motor de automação
│   ├── server/           # Servidor HTTP
│   ├── client/           # Cliente GitHub
│   └── config/           # Configuração
├── docs/                 # Documentação
└── support/              # Scripts de suporte
```

## 🛡️ Segurança

- ✅ **Sanitização de entrada**: Validação rigorosa de parâmetros
- ✅ **Rate limiting**: Respeito aos limites da API GitHub
- ✅ **Dry-run mode**: Modo seguro para testes
- ✅ **Repository scoping**: Apenas repositórios explicitamente configurados
- ✅ **Error recovery**: Tratamento robusto de erros e panic recovery

## 📈 Performance

- **Health Checking**: Sistema concorrente com cache para providers AI
- **Timeout Management**: Timeouts agressivos (3s) para verificações
- **Concurrent Operations**: Verificações paralelas para múltiplos providers
- **Intelligent Caching**: Cache thread-safe para evitar verificações repetitivas

## 🤝 Contribuição

1. **Fork** o projeto
2. **Clone** seu fork
3. **Branch** para sua feature (`git checkout -b feature/amazing-feature`)
4. **Commit** suas mudanças (`git commit -m 'Add amazing feature'`)
5. **Push** para a branch (`git push origin feature/amazing-feature`)
6. **Abra** um Pull Request

## 📋 Roadmap

- [ ] **Swagger Documentation**: Documentação automática de API
- [ ] **Webhook Support**: Suporte a webhooks GitHub
- [ ] **Advanced Analytics**: Métricas avançadas e dashboards
- [ ] **Team Management**: Gerenciamento de equipes e permissões
- [ ] **Scheduled Operations**: Operações agendadas e recorrentes

## 📄 Licença

Este projeto está licenciado sob a **MIT License** - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 👨‍💻 Autor

**Rafael Mori** - [@rafa-mori](https://github.com/rafa-mori)

---

🔗 **Links Úteis:**

- [Documentação de Endpoints](docs/endpoints.md)
- [Configuração Avançada](docs/config/)
- [Issues & Bug Reports](https://github.com/rafa-mori/ghbex/issues)
- [Discussões](https://github.com/rafa-mori/ghbex/discussions)
