# GitHub Actions Workflow Improvement Proposals for PipeCD

This document outlines 7 high-impact, low-risk workflow improvements designed to reduce maintainer effort, improve contributor experience, and align with PipeCD's GitOps, safety-first philosophy.

---

## 1. DCO Sign-off Check Workflow

### Problem
Contributors frequently forget to sign off their commits with `git commit -s`, leading to DCO check failures. Maintainers must manually ask contributors to fix this, causing back-and-forth and delays.

### Why This Matters for PipeCD
- DCO compliance is mandatory per CONTRIBUTING.md
- Currently, the `hack/ensure-dco.sh` script runs locally but fails silently in CI until the `check` target is run
- Early detection prevents wasted CI cycles and reviewer time

### Workflow Logic
```yaml
trigger: pull_request (opened, synchronize, reopened)
steps:
  1. Check out repository
  2. Run DCO check on all commits in the PR
  3. If missing sign-off:
     - Fail the check with clear error message
     - Post a helpful comment explaining how to fix (git rebase --signoff)
outputs: Pass/Fail status with actionable guidance
```

### Impact
- Reduces maintainer load by eliminating manual DCO reminders
- Helps new contributors understand requirements faster
- Catches issues before other CI runs waste resources

---

## 2. PR Description Validation Workflow

### Problem
Contributors often submit PRs with incomplete descriptions (empty "What this PR does", missing issue links, no user-facing change info). Maintainers spend time asking for this information.

### Why This Matters for PipeCD
- The PR template in `.github/PULL_REQUEST_TEMPLATE.md` has required sections
- Complete PR descriptions are essential for release notes and changelog generation
- Breaking changes need explicit documentation per CONTRIBUTING.md

### Workflow Logic
```yaml
trigger: pull_request (opened, edited)
steps:
  1. Parse PR body
  2. Check for required sections:
     - "What this PR does" is not empty
     - Issue link exists (Fixes #XXX)
     - User-facing change section is addressed
  3. If incomplete:
     - Add a comment with missing items
     - Set status to pending (not failing, to allow WIP PRs)
outputs: Status check with specific missing items listed
```

### Impact
- Saves maintainers from asking for missing information
- Ensures consistent PR quality for release notes
- Transparent and predictable expectations for contributors

---

## 3. PipeCD Application Config Validation Workflow

### Problem
The `examples/` directory contains PipeCD application configurations (`app.pipecd.yaml`). These may have syntax errors, invalid field names, or outdated API versions that aren't caught until deployment.

### Why This Matters for PipeCD
- Examples are used in documentation and by new users
- Broken examples damage trust and increase support burden
- As a CD system, PipeCD should validate its own configurations

### Workflow Logic
```yaml
trigger: pull_request (paths: examples/**)
steps:
  1. Find all app.pipecd.yaml files in examples/
  2. Validate YAML syntax
  3. Check apiVersion matches supported versions (pipecd.dev/v1beta1)
  4. Validate required fields (kind, spec.name)
  5. Report validation errors with file paths
outputs: Pass/Fail with specific validation errors per file
```

### Impact
- Catches example configuration errors before merge
- Ensures examples stay in sync with API changes
- Reduces user confusion and support requests

---

## 4. Documentation Link Checker Workflow

### Problem
Documentation often contains broken internal links (moved/renamed pages) and external links (deprecated URLs). These are discovered by users, not CI.

### Why This Matters for PipeCD
- PipeCD's documentation at pipecd.dev is a primary resource
- Broken links damage credibility and user experience
- Hugo-based docs can have subtle link issues with relative paths

### Workflow Logic
```yaml
trigger: 
  - pull_request (paths: docs/**)
  - schedule (weekly)
steps:
  1. Build Hugo site in docs/
  2. Run link checker (lychee or similar)
  3. Report broken internal links (fail)
  4. Report broken external links (warn, don't fail)
  5. Cache results to avoid rate limiting
outputs: List of broken links with file locations
```

### Impact
- Proactively catches documentation issues
- Reduces user-reported broken link issues
- Maintains documentation quality over time

---

## 5. First-time Contributor Welcome Workflow

### Problem
New contributors don't know about DCO requirements, code generation checks, or contribution guidelines until their PR fails.

### Why This Matters for PipeCD
- PipeCD wants to build a welcoming community
- The codebase has specific requirements (codegen, DCO, license headers)
- CONTRIBUTING.md is comprehensive but not always read

### Workflow Logic
```yaml
trigger: pull_request_target (opened)
conditions: actor is first-time contributor
steps:
  1. Check if user has previous merged PRs
  2. If first-time:
     - Post welcome comment with:
       - Link to CONTRIBUTING.md
       - Reminder about DCO sign-off
       - Link to development setup
       - Mention of `make check` command
     - Add "first-contribution" label
outputs: Welcome comment and label
```

### Impact
- Sets contributors up for success from the start
- Reduces friction and failed CI iterations
- Builds community by showing appreciation

---

## 6. PR Changed Files Area Summary Workflow

### Problem
Large PRs touch multiple areas (go, web, docs, examples). Maintainers must manually scan files to understand scope and assign appropriate reviewers.

### Why This Matters for PipeCD
- Different areas have different code owners and expertise
- The labeler workflow adds labels but doesn't summarize impact
- Multi-area PRs need special coordination

### Workflow Logic
```yaml
trigger: pull_request (opened, synchronize)
steps:
  1. Get list of changed files
  2. Categorize by area:
     - Go code (cmd/, pkg/)
     - Web (web/)
     - Docs (docs/)
     - Examples (examples/)
     - Manifests (manifests/)
     - CI/Build (.github/, hack/, Makefile)
  3. Generate summary comment with:
     - Affected areas
     - Count of files changed per area
     - Suggested reviewers based on CODEOWNERS
outputs: Comment with structured change summary
```

### Impact
- Helps maintainers quickly understand PR scope
- Facilitates better reviewer assignment
- Highlights cross-cutting changes that need extra attention

---

## 7. Kubernetes/Helm Manifest Linting Workflow

### Problem
Helm charts and Kubernetes manifests in `manifests/` can have subtle issues (deprecated APIs, missing required fields) that aren't caught until deployment.

### Why This Matters for PipeCD
- PipeCD ships Helm charts for pipecd, piped, and site
- Users install these charts in production
- Manifest issues can cause deployment failures

### Workflow Logic
```yaml
trigger: pull_request (paths: manifests/**)
steps:
  1. Run helm lint on all charts (already exists, enhance it)
  2. Run kubeconform/kubeval for Kubernetes API validation
  3. Check for deprecated API versions
  4. Validate values.yaml against chart requirements
outputs: Pass/Fail with specific lint/validation errors
```

### Impact
- Catches manifest issues before release
- Ensures charts work with supported Kubernetes versions
- Reduces deployment failures for users

---

## Implementation Priority

| Proposal | Effort | Impact | Priority |
|----------|--------|--------|----------|
| 1. DCO Check | Low | High | P0 |
| 5. First-time Welcome | Low | High | P0 |
| 3. Examples Validation | Medium | High | P1 |
| 2. PR Description Validation | Medium | Medium | P1 |
| 6. Changed Files Summary | Medium | Medium | P2 |
| 4. Docs Link Checker | Medium | Medium | P2 |
| 7. Manifest Linting | Low | Medium | P2 |

---

## Design Principles Applied

1. **Reuse existing tooling**: Uses existing scripts like `hack/ensure-dco.sh`
2. **Safety first**: Workflows that add comments don't block PRs unnecessarily
3. **Transparency**: Clear feedback with actionable next steps
4. **Composable**: Each workflow is independent and can be enabled/disabled
5. **Maintainer empathy**: Reduces repetitive tasks without adding complexity
