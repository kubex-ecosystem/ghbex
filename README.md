# ![GHbex Banner](docs/assets/nm_banner_md.png)

**A Intelligent GitHub repository management platform with automation, advanced analytics, and multi-AI integration. Automate, optimize, and monitor your GitHub repositories with intelligence and security.**

---

[![Kubex Go Dist CI](https://github.com/kubex-ecosystem/gemx/ghbex/actions/workflows/kubex_go_release.yml/badge.svg)](https://github.com/kubex-ecosystem/gemx/ghbex/actions/workflows/kubex_go_release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E=1.21-blue)](go.mod)
[![Releases](https://img.shields.io/github/v/release/rafa-mori/ghbex?include_prereleases)](https://github.com/kubex-ecosystem/gemx/ghbex/releases)

---

<!--
<p align="center">
  <img src="docs/assets/ghbex_demo.gif" alt="Animated demonstration of GHbex" width="80%"/>
  <br><em>GIF: Animated demonstration of GHbex (add here when available)</em>
</p>
-->

## 📑 Table of Contents

- [About the Project](#-about-the-project)
- [Main Features](#-main-features)
- [Installation](#-installation)
- [Configuration](#️-configuration)
- [Quick Start](#-quick-start)
- [Usage Examples](#-usage-examples)
- [CLI](#️-cli)
- [Architecture](#️-architecture)
- [Security](#-security)
- [Performance](#-performance)
- [Contributing](#-contributing)
- [Roadmap](#️-roadmap)
- [License](#-license)
- [Author](#-author)
- [Useful Links](#-useful-links)

---

## 🧩 About the Project

**GHbex** is an advanced platform for GitHub repository management, featuring artificial intelligence, automation, analytics, and optimization. It enables everything from repository sanitization and cleanup to intelligent recommendations, dependency analysis, workflow automation, and integration with multiple AI providers (Gemini, OpenAI, Claude, DeepSeek, Ollama).

Ideal for DevOps teams, software engineers, and maintainers seeking automation, governance, and insights for their repositories.

---

## ✨ Main Features

- **AI-Powered Operator:**
  - Fast analysis and intelligent recommendations
  - Multi-provider support (Gemini, OpenAI, Claude, DeepSeek, Ollama)
  - Concurrent health check for providers
- **Repository Sanitization:**
  - Automatic cleanup of old workflows, artifacts, and draft releases
  - Bulk operations for multiple repositories
- **Analytics & Insights:**
  - Health, dependency, vulnerability, and activity pattern analysis
  - Performance and engagement metrics
- **Productivity Optimization:**
  - Workflow analysis, auto-merge suggestions, notification optimization
  - ROI calculation
- **Automation Engine:**
  - Pattern recognition, automatic recommendations, trust scoring
  - Integration suggestions

---

## ⚡ Installation

Requirements: Go >= 1.21

```bash
# Clone the repository
git clone https://github.com/kubex-ecosystem/gemx/ghbex.git
cd ghbex

# Install dependencies
go mod tidy

# Build the project
make build-dev
```

---

## ⚙️ Configuration

```bash
# GitHub Authentication
export GITHUB_TOKEN="ghp_your_personal_token"

# AI Providers (optional)
export GEMINI_API_KEY="your_gemini_api_key"
export OPENAI_API_KEY="your_openai_api_key"

# Discord Notifications (optional)
export DISCORD_WEBHOOK_URL="your_discord_webhook_url"
```

---

## 🚀 Quick Start

```bash
# Start the server
./dist/ghbex start --owner rafa-mori --port 8088 --repos 'owner/repo1,owner/repo2'

# Or using environment variables
export REPO_LIST='owner/repo1,owner/repo2'
./dist/ghbex start --port 8088
```

### Access

- **Dashboard**: <http://localhost:8088>
- **Health Check**: <http://localhost:8088/health>
- **API Docs**: [docs/endpoints.md](docs/endpoints.md)

---

## 🧪 Usage Examples

### Intelligent Repository Analysis

```bash
ghbex intelligence --repo rafa-mori/ghbex
```

### Bulk Sanitization

```bash
ghbex sanitize --repos 'rafa-mori/ghbex,rafa-mori/logz'
```

### Automation Recommendations

```bash
ghbex automation --repo rafa-mori/ghbex
```

---

## 🖥️ CLI

```bash
# Start server
ghbex start --owner <owner> --port <port> --repos '<repo1,repo2>'

# Check status
ghbex status

# Stop server
ghbex stop

# Check configuration
ghbex config

# Show version
ghbex version
```

---

## 🏗️ Architecture

```plaintext
ghbex/
├── cmd/                    # CLI entrypoints
│   ├── main.go            # Main entrypoint
│   └── cli/               # CLI commands
├── internal/              # Internal code
│   ├── operators/         # Specialized operators
│   │   ├── intelligence/  # AI and insights
│   │   ├── sanitize/      # Cleanup
│   │   ├── analytics/     # Metrics
│   │   ├── productivity/  # Optimization
│   │   └── automation/    # Automation
│   ├── server/           # HTTP server
│   ├── client/           # GitHub client
│   └── config/           # Configuration
├── docs/                 # Documentation
└── support/              # Support scripts
```

---

## 🔒 Security

- **Input sanitization**: Strict parameter validation
- **Rate limiting**: Respects GitHub API limits
- **Dry-run mode**: Safe execution for testing
- **Restricted scope**: Only explicitly configured repositories
- **Error recovery**: Robust error and panic handling

---

## 🚀 Performance

- **Concurrent health check** for AI providers
- **Aggressive timeouts** (3s) for checks
- **Parallel operations** for multiple providers
- **Smart cache**: thread-safe to avoid repetitions

---

## 🤝 Contributing

1. Fork the project
2. Clone your fork
3. Create a branch (`git checkout -b feature/new-feature`)
4. Commit your changes (`git commit -m 'feat: new feature'`)
5. Push to the branch (`git push origin feature/new-feature`)
6. Open a Pull Request

---

## 🗺️ Roadmap

- [ ] Swagger Docs: automatic API documentation
- [ ] Webhook Support: GitHub webhook integration
- [ ] Advanced analytics: dashboards and metrics
- [ ] Team management: permissions and teams
- [ ] Scheduled operations: recurring executions

---

## 📄 License

Project under **MIT** license — see the [LICENSE](LICENSE) file.

---

## 👨‍💻 Author

**Rafael Mori** — [@rafa-mori](https://github.com/kubex-ecosystem)

---

## 🔗 Useful Links

- [Endpoints Documentation](docs/endpoints.md)
- [Advanced Configuration](docs/config/)
- [Issues & Bug Reports](https://github.com/kubex-ecosystem/gemx/ghbex/issues)
- [Discussions](https://github.com/kubex-ecosystem/gemx/ghbex/discussions)
