# Maintenance: Helm chart is pushed once per matrix job in publish_image_chart.yaml

**File:** `.github/workflows/publish_image_chart.yaml`

## Problem

The `artifacts` matrix has 7 images × 2 registries = up to 14 parallel jobs. Every job runs the full Helm chart publish block (install helm, login, build chart, helm push). The same charts are pushed 14 times, producing redundant OCI pushes and confusing logs.

The push is idempotent (OCI rejects duplicate tags), so no corruption occurs, but it wastes time and creates unnecessary registry noise.

## Fix

Extract the Helm chart build and push into a separate job that `needs: artifacts` and runs once. This is already the pattern used for `trigger-event-watcher` and `release-quickstart-manifests` in the same file.
