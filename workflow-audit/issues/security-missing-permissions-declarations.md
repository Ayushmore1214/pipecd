# Security: Multiple workflows missing explicit permissions declarations

**Files:** `cherry_pick.yaml`, `codeql-analysis.yaml`, `stale.yaml`, `release.yaml`, `build.yaml`, `test.yaml`

## Problem

Workflows without explicit `permissions:` blocks inherit the repository-level default, which may be `contents: write` in some configurations. This violates the principle of least privilege and makes the security posture dependent on a repo setting rather than workflow intent.

Most critically: `codeql-analysis.yaml` is missing `security-events: write`, which is **required** for CodeQL to upload SARIF results to the Security tab. Without it, the scan runs every Monday but findings may be silently discarded.

## Required Permissions per Workflow

| Workflow | Minimum permissions needed |
|---|---|
| `codeql-analysis.yaml` | `contents: read`, `security-events: write`, `actions: read` |
| `stale.yaml` | `issues: write`, `pull-requests: write` |
| `cherry_pick.yaml` | `contents: write`, `pull-requests: write` |
| `release.yaml` | `contents: write`, `pull-requests: write` |
| `build.yaml` | `contents: read` |
| `test.yaml` | `contents: read` |
