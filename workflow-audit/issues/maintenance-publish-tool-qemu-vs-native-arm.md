# Maintenance: publish_tool.yaml uses QEMU emulation while build_tool.yaml uses native ARM runners

**Files:** `.github/workflows/build_tool.yaml`, `.github/workflows/publish_tool.yaml`

## Problem

`build_tool.yaml` (PR check) builds tool images using native `ubuntu-24.04-arm` runners — fast and architecturally correct.

`publish_tool.yaml` (merge/tag publish) builds the same images using QEMU emulation on a single x86 runner — slow and can produce subtly incorrect binaries for architecture-specific code.

`publish_pipedv1_exp.yaml` uses the correct pattern: native runners per arch + digest merge into a multi-arch manifest.

## Fix

Migrate `publish_tool.yaml` to the same native runner + digest merge pattern used in `publish_pipedv1_exp.yaml`. This ensures that what is validated in CI is exactly what gets published.
