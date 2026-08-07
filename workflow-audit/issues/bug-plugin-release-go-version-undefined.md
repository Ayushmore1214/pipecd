# Bug: plugin_release.yaml references undefined GO_VERSION and uses non-existent checkout@v5

**File:** `.github/workflows/plugin_release.yaml`

## Problem

Two defects that will cause this workflow to fail on every plugin release:

1. `setup-go` references `${{ env.GO_VERSION }}` but `GO_VERSION` is never declared in this workflow's `env:` block. It expands to an empty string, causing `setup-go` to pick an arbitrary default Go version rather than the project-required `1.26.2`.

2. `actions/checkout@v5` is used on line 21. v5 does not exist — all other workflows in this repo use the SHA-pinned `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2`.

## Fix

- Add `GO_VERSION: 1.26.2` to the workflow `env:` block.
- Change `actions/checkout@v5` to `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2`.

## Impact

Every manual plugin release silently uses the wrong Go version or fails checkout entirely.
