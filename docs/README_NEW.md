# GHbex - GitHub Intelligence Platform

**GHBEX** is a comprehensive GitHub repository intelligence and automation platform built with Go. Originally started as a sanitization tool, it has evolved into a powerful multi-operator system for GitHub analytics, productivity optimization, and automation.

## 🚀 **Core Operators**

### ✅ **SANITIZATION** - Repository Cleanup

- Delete old workflow runs and artifacts
- Remove draft releases and outdated content
- Configurable rules and dry-run support
- JSON/Markdown reporting

### ✅ **ANALYTICS** - Deep Repository Intelligence

- Repository health scoring (0-100)
- Community and contributor analysis
- Code intelligence and language distribution
- Issue/PR insights and maintenance metrics
- DevEx (Developer Experience) scoring

### ✅ **PRODUCTIVITY** - Development Optimization

- Template analysis and recommendations
- Branching strategy optimization
- Auto-merge opportunity detection
- Developer experience improvements
- ROI calculations and implementation guides

### 🔄 **AUTOMATION** - Smart Workflows *(In Development)*

- Organizational automation workflows
- Repository governance automation
- Communication and notification automation
- Policy enforcement and compliance

## 🛠 **API Endpoints**

```bash
# Health & Status
GET  /health
GET  /repos

# Analytics Intelligence
GET  /analytics/{owner}/{repo}?days=90

# Productivity Analysis
GET  /productivity/{owner}/{repo}

# Sanitization Operations
POST /admin/repos/{owner}/{repo}/sanitize?dry_run=1
POST /admin/sanitize/bulk

# Authentication
# Set GITHUB_PAT_TOKEN environment variable
```

## ⚡ **Quick Start**

```bash
# Install dependencies
go mod tidy

# Configure (optional - has sensible defaults)
cp config/sanitize.yaml.example config/sanitize.yaml

# Build
make build-dev

# Start server
GITHUB_PAT_TOKEN=<your_token> ./dist/ghbex start

# Test analytics
curl "http://localhost:8088/analytics/owner/repo" | jq

# Test productivity
curl "http://localhost:8088/productivity/owner/repo" | jq

# Test sanitization (dry-run)
curl -X POST "http://localhost:8088/admin/repos/owner/repo/sanitize?dry_run=1" | jq
```

## 📊 **Real-World Results**

**Analytics Example:**

```bash
# Repository health analysis
✅ Health Score: 75.8/100 (Grade: B)
✅ DevEx Score: 85/100
✅ Community Activity: High
✅ Code Quality: Excellent
```

**Productivity Example:**

```bash
# Productivity optimization
✅ 5 actionable recommendations
✅ ROI: 3.2x return on investment
✅ $12,500 estimated time savings
✅ 65% setup complexity reduction
```

## 🏗 **Architecture**

- **Go 1.21+** with google/go-github/v61
- **Interface-based design** for extensibility
- **Real-time performance logging**
- **No mocks - real GitHub API testing**
- **Modular operator system**

## 📁 **Project Structure**

```
├── cmd/                     # CLI and server entry points
├── internal/
│   ├── operators/           # Core business logic
│   │   ├── analytics/       # Repository intelligence
│   │   ├── productivity/    # Development optimization
│   │   ├── sanitize/        # Repository cleanup
│   │   └── automation/      # Smart workflows (WIP)
│   ├── server/              # HTTP server and routing
│   ├── client/              # GitHub API client
│   └── config/              # Configuration management
├── config/                  # Configuration files
└── _reports/                # Generated reports
```

## 🔧 **Configuration**

### Authentication

- **PAT Token**: `GITHUB_PAT_TOKEN` environment variable
- **GitHub Apps**: JWT + Installation Token support
- **GHES**: Custom base URL configuration

### Server

- **Port**: `:8088` (configurable)
- **Timeout**: 5s read header timeout
- **CORS**: Enabled for development

## 📈 **Performance**

- **Analytics**: 6-8 seconds average response
- **Productivity**: 6-8 seconds average response
- **Sanitization**: 15-30 seconds (depends on repository size)
- **Memory**: ~50MB base footprint
- **Concurrency**: Safe for concurrent requests

## 🎯 **Development Status**

**Current Version**: MVP+ (Pilot Phase)
**Status**: Active development on `pilot` branch
**Next Release**: Automation operator completion

## 🤝 **Contributing**

This project follows Go best practices and emphasizes:

- **No mocks** - Real-world testing only
- **Interface segregation** - Clean, testable code
- **Performance focus** - Sub-10s response times
- **Real value** - Practical productivity gains

## 📝 **Documentation**

- 📋 [Session Progress](./SESSION_PROGRESS.md) - Current development status
- 🔧 [Technical Details](./docs/) - Deep dive documentation
- 🚀 [Getting Started Guide](./docs/README.pt-BR.md) - Portuguese guide

## 🌟 **Vision**

Building a **revolutionary GitHub intelligence platform** that helps developers and organizations optimize their workflows, improve productivity, and automate repetitive tasks.

***"Something that will really help A LOT of people!"***

---

**License**: MIT | **Author**: Rafael Mori | **Organization**: Kubex
