# Maintenance: build.yaml and gen.yaml are missing concurrency controls

**Files:** `.github/workflows/build.yaml`, `.github/workflows/gen.yaml`

## Problem

`lint.yaml` and `test.yaml` both use the concurrency pattern:

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: ${{ github.event_name == 'pull_request' }}
```

`build.yaml` and `gen.yaml` do not. On active PRs with multiple pushes, stale build and codegen-validation jobs queue up and waste CI minutes. On push to master the cancel-in-progress is false so those runs complete safely.

## Fix

Add the same `concurrency:` block to `build.yaml` and `gen.yaml`.
