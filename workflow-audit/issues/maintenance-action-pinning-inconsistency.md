# Maintenance: Inconsistent action pinning across workflows

## Problem

The repo pins security-sensitive Docker actions to SHA (good) but leaves many other actions on mutable tags. Dependabot tracks `github-actions` weekly for version bumps but does not enforce SHA pinning.

## Unpinned Actions (not SHA-pinned)

| Action | Used in |
|---|---|
| `actions/setup-go@v3` | build.yaml, test.yaml, publish_binary.yaml, plugin_release.yaml, publish_pipedv1_exp.yaml |
| `actions/setup-node@v3` | build.yaml, test.yaml, publish_site.yaml |
| `azure/setup-helm@v4` | build.yaml, lint.yaml, publish_image_chart.yaml, publish_site.yaml, publish_pipedv1_exp.yaml |
| `actions/labeler@v4` | labeler.yaml |
| `actions/stale@v8` | stale.yaml |
| `actions/github-script@v7` | first-time-contributor.yaml, thank-you.yaml |
| `codecov/codecov-action@v3` | test.yaml |
| `pipe-cd/actions-gh-release@v2.6.0` | prerelease.yaml |
| `github/codeql-action/*@v3` | codeql-analysis.yaml |
| `peter-evans/create-pull-request@v6` | release.yaml, publish_image_chart.yaml, publish_pipedv1_exp.yaml |
| `actions/checkout@v4` | thank-you.yaml (all others are SHA-pinned) |

## Priority

Highest risk: `ca-dp/code-butler@v1` (see separate security issue), `codecov/codecov-action`, `peter-evans/create-pull-request`.
Lower risk: GitHub-owned actions (`setup-go`, `setup-node`, `labeler`, `stale`) where Dependabot provides version tracking.
