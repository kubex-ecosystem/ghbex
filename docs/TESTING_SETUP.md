# GHBEX Real-World Testing Setup

## 🔧 Quick Test Configuration

### 1. GitHub Personal Access Token (Easiest for testing)

```bash
# Create your config file
cp /srv/apps/LIFE/KUBEX/ghbex/docs/config/sanitize.yaml /srv/apps/LIFE/KUBEX/ghbex/config/sanitize.yaml

# Export your GitHub token (get from: https://github.com/settings/tokens)
export GITHUB_TOKEN="your_token_here"

# Optional: Discord webhook for notifications
export DISCORD_WEBHOOK_URL="your_discord_webhook_here"
```

### 2. Test Repository Setup

We'll need a **test repository** to safely run operations:

#### Option A: Create Test Repo

```bash
# Use GitHub CLI to create a test repo
gh repo create ghbex-test --private --description "Test repository for GHBEX sanitization"
```

#### Option B: Use Existing Repo (Safer)

- Choose a repository with old workflow runs/artifacts
- Make sure it's not critical (we'll use dry_run first!)

### 3. Configuration File

Edit `/srv/apps/LIFE/KUBEX/ghbex/config/sanitize.yaml`:

```yaml
runtime:
  dry_run: true  # ALWAYS start with dry_run!
  report_dir: ./_reports

github:
  auth:
    kind: "pat"
    token: "${GITHUB_TOKEN}"
  repos:
    - owner: "YOUR_USERNAME"
      name: "TEST_REPO_NAME"
      rules:
        runs:
          max_age_days: 30
          keep_success_last: 5
        artifacts:
          max_age_days: 7
        releases:
          delete_drafts: true
        security:
          rotate_ssh_keys: false  # Start with false for testing
          remove_old_keys: false
        monitoring:
          check_inactivity: true
          inactive_days_threshold: 30
          monitor_prs: true
          monitor_issues: true

notifiers:
  - type: "stdout"
  # - type: "discord"
  #   webhook: "${DISCORD_WEBHOOK_URL}"
```

## 🚀 Ready to Test?

Once you have the config set up, just tell me:

1. **Your GitHub username** (for the config)
2. **Test repository name** (or if you want to create one)
3. **GitHub token ready** (we'll test dry_run first)

And we'll make this baby **SING!** 🎵

### What We'll Test First

1. ✅ **Connection** - GitHub API authentication
2. ✅ **Dry Run** - See what would be cleaned without doing it
3. ✅ **Reports** - Generate beautiful markdown reports
4. ✅ **Monitoring** - Check repository activity
5. ✅ **Real Operations** - When you're ready!

**READY TO ROCK AND ROLL?** 🤘
