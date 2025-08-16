# GHBEX - Sanitization Engine Implementation Plan

## 🎯 OBJECTIVE

Transform GHBEX into a complete GitHub repository sanitization and monitoring engine, integrated with GoBE via MCP.

## 📋 CURRENT STATUS (What's ALREADY working!)

### ✅ IMPLEMENTED

- [x] GitHub Authentication (PAT + GitHub App)
- [x] Workflow Runs Cleanup (complete with pagination)
- [x] Basic Artifacts & Releases structure
- [x] Report System (JSON + Markdown)
- [x] Discord Notifications
- [x] HTTP Server with sanitize endpoint
- [x] Configuration System (YAML)
- [x] Dry-run support

### 🔄 IN PROGRESS (Current pilot branch)

- [ ] Complete Artifacts cleanup implementation
- [ ] Complete Releases cleanup implementation
- [ ] SSH Key rotation system
- [ ] Repository monitoring (PRs, Issues, Activity)
- [ ] MCP integration bridge

## 🚀 IMPLEMENTATION ROADMAP

### PHASE 1: Complete Core Sanitization (Current Sprint)

**Target: Get basic sanitization 100% functional**

1. **Complete Artifacts Cleanup** (1-2 hours)
   - Implement pagination in artifacts operator
   - Add age filtering and deletion logic
   - Test with real repositories

2. **Complete Releases Cleanup** (1-2 hours)
   - Implement draft deletion
   - Add release management logic
   - Test draft cleanup functionality

3. **Testing & Validation** (1 hour)
   - Create test repository scenarios
   - Validate all cleanup operations
   - Ensure reports are accurate

### PHASE 2: Security & SSH Management (Next Sprint)

**Target: Add SSH key rotation and security features**

1. **SSH Key Rotation System** (3-4 hours)
   - Create `internal/operators/security/` package
   - Implement SSH key generation and rotation
   - GitHub Deploy Keys API integration
   - Repository-specific key management

2. **Security Auditing** (2 hours)
   - Repository security settings check
   - Vulnerability scanning integration
   - Security policy compliance

### PHASE 3: Monitoring & Analytics (Following Sprint)

**Target: Add proactive monitoring capabilities**

1. **Repository Activity Monitoring** (3 hours)
   - PR and Issue tracking
   - Commit activity analysis
   - Contributor activity patterns
   - Inactivity detection algorithms

2. **Alerting System** (2 hours)
   - Smart notification logic
   - Priority-based alerts
   - Historical trend analysis

### PHASE 4: MCP Integration (Final Sprint)

**Target: Bridge GHBEX with GoBE MCP system**

1. **MCP Bridge Creation** (4 hours)
   - Create MCP controller in GoBE
   - GHBEX service integration
   - API communication layer
   - Task scheduling via GoBE

2. **Orchestration Layer** (3 hours)
   - Multi-repository operations
   - Batch processing capabilities
   - Progress tracking and reporting

## 🏗️ TECHNICAL ARCHITECTURE

### Current Structure (Working)

```
ghbex/
├── internal/
│   ├── client/           # ✅ GitHub API clients
│   ├── operators/        # ✅ Core operations
│   │   ├── workflows/    # ✅ Workflow runs cleanup
│   │   ├── artifacts/    # 🔄 Basic structure
│   │   ├── releases/     # 🔄 Basic structure
│   │   └── sanitize/     # ✅ Report generation
│   ├── notifiers/        # ✅ Discord integration
│   └── server/           # ✅ HTTP server
```

### Planned Extensions

```
ghbex/
├── internal/
│   ├── operators/
│   │   ├── security/     # 🆕 SSH & security
│   │   └── monitoring/   # 🆕 Activity tracking
│   ├── analytics/        # 🆕 Data analysis
│   └── bridge/           # 🆕 MCP integration
```

## 🎯 SUCCESS METRICS

### Phase 1 Success Criteria

- [ ] Successfully clean 100+ workflow runs in test repo
- [ ] Delete 50+ old artifacts without issues
- [ ] Remove 10+ draft releases safely
- [ ] Generate comprehensive reports
- [ ] Send Discord notifications with attachments

### Phase 2 Success Criteria

- [ ] Rotate SSH keys on 5+ repositories
- [ ] Generate security audit reports
- [ ] Detect and report security vulnerabilities

### Phase 3 Success Criteria

- [ ] Monitor 10+ repositories for activity
- [ ] Detect inactive repositories (30+ days)
- [ ] Track PR/Issue trends over time
- [ ] Send proactive alerts for important changes

### Phase 4 Success Criteria

- [ ] Execute sanitization tasks via GoBE MCP
- [ ] Schedule periodic maintenance jobs
- [ ] Orchestrate multi-repo operations
- [ ] Provide unified dashboard in GoBE

## 💡 IMPLEMENTATION PRINCIPLES

### No Mocks Philosophy

- Use real GitHub repositories for testing
- Integration tests with actual API calls
- Contract-based testing between components
- Behavioral testing for complex workflows

### Interface-First Design

- Define contracts before implementations
- Keep internal packages unexposed until necessary
- Use dependency injection for testability
- Composition over inheritance patterns

### Incremental Development

- Each phase delivers working functionality
- Continuous integration and testing
- Regular progress validation
- User feedback incorporation

## 🔧 DEVELOPMENT WORKFLOW

### Current Sprint Tasks (Phase 1)

1. **TODAY**: Complete artifacts cleanup implementation
2. **TODAY**: Complete releases cleanup implementation
3. **TODAY**: End-to-end testing with real repositories
4. **TOMORROW**: Performance optimization and error handling

### Next Steps

- Move to Phase 2 after Phase 1 validation
- Regular progress reviews and adjustments
- Continuous integration with GoBE development
- Documentation updates for each completed phase

---

**READY TO START PHASE 1? Let's make it happen! 🚀**
