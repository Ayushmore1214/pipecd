# Maintenance: Several workflows and jobs are missing timeout-minutes

**Files:** `build.yaml`, `test.yaml` (integration), `codeql-analysis.yaml`, `gen.yaml`, `publish_image_chart.yaml` (artifacts job), `release.yaml`, `stale.yaml`, `prerelease.yaml`

## Problem

Jobs with no `timeout-minutes` can hang indefinitely. On self-hosted runners (`oracle-vm-8cpu-32gb-x86-64`) this is especially dangerous — a hung Docker build holds the runner and blocks all subsequent jobs.

`build_tool.yaml` uses `timeout-minutes: 15` and `publish_tool.yaml` uses `timeout-minutes: 30` — these are the right reference points.

## Recommended Timeouts

| Job | Suggested timeout |
|---|---|
| `build.yaml` go/plugin/web/chart | 30 minutes |
| `test.yaml` integration | 60 minutes |
| `codeql-analysis.yaml` analyze | 120 minutes |
| `publish_image_chart.yaml` artifacts | 60 minutes |
| `gen.yaml` code | 15 minutes |
| `release.yaml` tool | 15 minutes |
| `stale.yaml` stale | 10 minutes |
| `prerelease.yaml` gh-release | 10 minutes |
