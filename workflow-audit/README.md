# GitHub Workflow Audit

This folder documents a full audit of every GitHub Actions workflow in `.github/workflows/`. Each change has its own branch and pull request in the fork. The `issues/` folder contains detailed write-ups for every problem found.

---

## What Was Audited

All 20 workflow files were reviewed for:
- Bugs that cause workflows to silently fail
- Security risks (permissions, secrets, third-party actions)
- Missing safety controls (timeouts, concurrency)
- Consistency problems across the codebase
- Wasted CI resources

---

## Changes Made

### 1. Fix undefined `GO_VERSION` in `plugin_release.yaml`

**Branch:** `fix/plugin-release-go-version-and-checkout`
**Issue file:** `issues/bug-plugin-release-go-version-undefined.md`

#### What was wrong

Every time a maintainer manually triggered a plugin release, the workflow tried to set up Go using a variable called `GO_VERSION` — but that variable was never actually defined in the file. Go would pick a random default version instead of the required `1.26.2`. The file also used `actions/checkout@v5`, which does not exist (the latest is v4).

#### What was changed

- Added `GO_VERSION: 1.26.2` to the workflow's environment variables section, matching every other workflow in the project.
- Changed `actions/checkout@v5` to the same SHA-pinned v4.2.2 used everywhere else.

#### Why it matters

Every plugin release was being built with an incorrect or unpredictable Go version. The binary output could be different from what was tested. This is a silent correctness bug — the workflow would succeed but the result would be wrong.

---

### 2. Remove `code-butler.yaml` (security risk)

**Branch:** `fix/remove-code-butler`
**Issue file:** `issues/security-code-butler-unpinned-abandoned.md`

#### What was wrong

This workflow used an action called `ca-dp/code-butler@v1` to provide AI code review on pull requests using `/review` and `/chat` commands. Two serious problems:

1. **Not pinned to a specific version.** The `@v1` tag is mutable — the owner of that repository can change what code runs under it at any time. If their account or repo were ever compromised, arbitrary code would run inside PipeCD's CI with access to the `OPENAI_API_KEY` secret.

2. **Already broken.** The action uses `gpt-4-1106-preview`, a model OpenAI deprecated in early 2024. The `/review` and `/chat` commands have likely been non-functional for over a year.

#### What was changed

The entire `code-butler.yaml` file was deleted.

#### Why it matters

Passing a secret API key to an unmaintained third-party action with a mutable tag is one of the highest-risk things a workflow can do. Even if the risk of compromise is low, the impact — leaking the API key and executing arbitrary CI code — is high. The reward (AI review via a broken action) is zero. Removing it costs nothing.

#### What to do if AI review is wanted again

Instead of using a third-party action, write the logic directly in the workflow using `actions/github-script`. That keeps all the code inside the repo, under version control, with no external trust required.

---

### 3. Add permissions, concurrency controls, and timeouts across 11 workflows

**Branch:** `fix/workflow-permissions-concurrency-timeouts`
**Issue files:** `issues/security-missing-permissions-declarations.md`, `issues/maintenance-missing-concurrency-controls.md`, `issues/maintenance-missing-timeouts.md`

This is the largest single change. It touches 11 workflow files but every individual edit is small and safe — no logic changes, just adding missing configuration.

#### 3a. Explicit Permissions

**Files changed:** `build.yaml`, `test.yaml`, `gen.yaml`, `codeql-analysis.yaml`, `stale.yaml`, `cherry_pick.yaml`, `release.yaml`, `prerelease.yaml`

**What was wrong:**

GitHub Actions workflows have a `permissions:` setting that controls what the workflow is allowed to do with the repository. When it is not set, the workflow inherits whatever the repository's default is — which can be overly broad (for example, full `contents: write` when only `contents: read` is needed).

The most critical case was `codeql-analysis.yaml`. CodeQL is a security scanning tool that runs weekly to find vulnerabilities in Go and JavaScript code. Uploading its results to GitHub's Security tab requires a permission called `security-events: write`. That permission was missing. This means the scan was running every Monday but the findings were potentially being discarded silently — nobody could see them in the Security tab.

**What was changed:**

Every workflow that was missing permissions now has an explicit `permissions:` block with only the minimum permissions it actually needs:

| Workflow | Permissions added |
|---|---|
| `codeql-analysis.yaml` | `contents: read`, `security-events: write`, `actions: read` |
| `stale.yaml` | `issues: write`, `pull-requests: write` |
| `cherry_pick.yaml` | `contents: write`, `pull-requests: write` |
| `release.yaml` | `contents: write`, `pull-requests: write` |
| `prerelease.yaml` | `contents: write` |
| `build.yaml`, `test.yaml`, `gen.yaml` | `contents: read` |

**Why it matters:**

Least-privilege permissions are a basic security principle. A workflow that only reads code should never have write access — if something goes wrong (a bug, a supply chain attack, a confused action), narrower permissions limit the blast radius. The CodeQL fix is also a correctness fix: security findings are useless if they are never uploaded.

---

#### 3b. Concurrency Controls

**Files changed:** `build.yaml`, `gen.yaml`

**What was wrong:**

When you push multiple commits quickly to a pull request, GitHub starts a new CI run for each push. Without concurrency controls, all those runs queue up and execute in parallel — wasting CI minutes running builds and checks for code that has already been superseded by a newer commit.

`lint.yaml` and `test.yaml` already had the correct pattern:

```yaml
concurrency:
  group: ${{ github.workflow }}-${{ github.event.pull_request.number || github.ref }}
  cancel-in-progress: ${{ github.event_name == 'pull_request' }}
```

This means: "for a given PR, cancel any old run when a new one starts." On pushes to `master` (not PRs), it does not cancel — those always run to completion.

`build.yaml` and `gen.yaml` were missing this entirely.

**What was changed:**

Added the same concurrency block to `build.yaml` and `gen.yaml`.

**Why it matters:**

Faster feedback for contributors (they are not waiting for old runs to finish), lower CI costs (no wasted compute), and a cleaner GitHub UI (no pile of cancelled runs to scroll through).

---

#### 3c. Timeouts

**Files changed:** `build.yaml`, `test.yaml`, `codeql-analysis.yaml`, `gen.yaml`, `stale.yaml`, `cherry_pick.yaml`, `release.yaml`, `prerelease.yaml`, `first-time-contributor.yaml`, `thank-you.yaml`

**What was wrong:**

GitHub Actions jobs that have no `timeout-minutes` setting can run forever if something hangs. For regular GitHub-hosted runners this is annoying (it holds up the queue). For self-hosted runners like `oracle-vm-8cpu-32gb-x86-64` — which PipeCD uses for heavy Go tests and Docker builds — a hung job locks that runner and blocks every other job that needs it. There is no automatic recovery.

`build_tool.yaml` had `timeout-minutes: 15` and `publish_tool.yaml` had `timeout-minutes: 30` as good examples, but most other workflows did not follow the pattern.

**What was changed:**

Every job that was missing a timeout now has one appropriate to what it does:

| Job | Timeout |
|---|---|
| Build jobs (Go, web, chart) | 30 minutes |
| Integration tests | 60 minutes |
| CodeQL analysis | 120 minutes |
| Codegen validation | 15 minutes |
| Release/cherry-pick/stale/prerelease | 10–15 minutes |
| First-time contributor welcome, thank-you comment | 5 minutes |

**Why it matters:**

A real example of the risk: a Docker build that gets stuck network-waiting will sit forever on a self-hosted runner, blocking every release and merge until someone notices and manually cancels it hours later. Timeouts are a safety net.

---

#### 3d. Runner and Checkout Consistency

**Files changed:** `first-time-contributor.yaml`, `thank-you.yaml`

**What was wrong:**

Both workflows used `ubuntu-latest` as the runner instead of the explicit `ubuntu-24.04` that every other workflow in the repo uses. `ubuntu-latest` is a moving target — GitHub can change what version it points to at any time, potentially breaking workflows.

`thank-you.yaml` also used `actions/checkout@v4` (a mutable tag) while every other workflow in the repo uses the SHA-pinned version `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683`.

**What was changed:**

- Both workflows now use `ubuntu-24.04` explicitly.
- `thank-you.yaml` now uses the SHA-pinned checkout action.

**Why it matters:**

Consistency makes the codebase easier to maintain and audit. When Dependabot opens a PR to update the checkout action SHA, you want every workflow to be covered — not have some workflows silently stuck on a floating tag that never gets updated.

---

### 4. Replace QEMU emulation with native ARM runners in `publish_tool.yaml`

**Branch:** `fix/publish-tool-native-arm-runners`
**Issue file:** `issues/maintenance-publish-tool-qemu-vs-native-arm.md`

#### What was wrong

PipeCD's tool images (things like `actions-gh-release`, `codegen`, `piped-base`) need to run on both x86 (amd64) and ARM (arm64) machines. There are two ways to build multi-architecture Docker images:

1. **QEMU emulation** — Run everything on one x86 machine, and have it pretend to be ARM when building the ARM version. It is simple to set up but slow, and it can produce incorrect binaries for architecture-specific code because the CPU is being simulated, not real.

2. **Native runners** — Run the ARM build on an actual ARM machine. GitHub provides `ubuntu-24.04-arm` runners for exactly this. Fast and correct.

The PR check workflow (`build_tool.yaml`) already used native runners — the right way. But the publish workflow (`publish_tool.yaml`) — which creates the actual images that users download — was using QEMU. This meant that what was being validated in CI was architecturally different from what was being published.

The experimental pipedv1 publish workflow (`publish_pipedv1_exp.yaml`) already had the correct native-runner pattern with a digest merge step. `publish_tool.yaml` just needed to match it.

#### What was changed

Rewrote `publish_tool.yaml` to use the same two-stage approach:

- **`build` job:** Runs in a matrix of platform × image (2 platforms × 6 images = 12 parallel jobs). Each job uses a native runner (`ubuntu-24.04` for amd64, `ubuntu-24.04-arm` for arm64) and builds one image for one architecture. Instead of pushing a full tagged image, it pushes a digest (a content-addressed blob without a tag) and uploads the digest ID as a GitHub Actions artifact.

- **`merge` job:** Runs once per image (6 jobs) after `build` finishes. Downloads the two digests (one amd64, one arm64) and uses `docker buildx imagetools create` to combine them into a single multi-architecture manifest with the proper tag.

#### Why it matters

- **Correctness:** The ARM image is now built on real ARM hardware, not emulated x86.
- **Speed:** Native ARM builds are significantly faster than QEMU emulation.
- **Consistency:** What CI validates during the PR is now structurally identical to what gets published on merge. This eliminates an entire class of "works in CI, broken in production on ARM" bugs.

---

### 5. Extract Helm chart push from matrix in `publish_image_chart.yaml`

**Branch:** `fix/extract-helm-push-from-matrix`
**Issue file:** `issues/maintenance-helm-push-runs-in-every-matrix-job.md`

#### What was wrong

`publish_image_chart.yaml` builds and pushes 7 Docker images to 2 registries in a parallel matrix, giving up to 14 simultaneous jobs. Each of those 14 jobs also contained the complete Helm chart build and push block. So the same three Helm charts (`pipecd`, `piped`, `helloworld`) were being built and pushed 14 times per workflow run.

OCI registries (where Helm charts are stored) reject duplicate pushes for the same tag, so nothing was actually corrupted. But it meant:

- 13 out of 14 Helm pushes were doing wasted work
- CI logs had 14 identical Helm publish sections, making it hard to see what actually happened
- Any Helm push failure would appear 14 times, not once

`trigger-event-watcher` and `release-quickstart-manifests` — both downstream jobs that need the charts to be available — had a similar problem: they depended on `artifacts` (the image matrix), which meant they could start running as soon as one image was pushed, before the charts were necessarily done.

#### What was changed

- Removed the Helm build and push steps from the `artifacts` matrix job entirely.
- Added a new `publish-helm-charts` job that runs once, after `artifacts` completes, on a standard `ubuntu-24.04` runner.
- Updated `trigger-event-watcher` and `release-quickstart-manifests` to depend on `publish-helm-charts` instead of `artifacts`, so they correctly wait for charts to be published before running.

The flow is now: `artifacts` (14 parallel image builds) → `publish-helm-charts` (1 chart publish) → `trigger-event-watcher` and `release-quickstart-manifests`.

#### Why it matters

- **Cleaner CI logs:** Helm publish appears exactly once, not 14 times.
- **Correct job ordering:** Downstream jobs now wait for charts to be ready, which they actually need.
- **Less registry noise:** 13 fewer redundant OCI pushes per release.
- **Easier debugging:** If a Helm push fails, you see one failure, not 14.

---

## Issues That Need Maintainer Discussion (No Code Change Yet)

### Unimplemented `gh_release` job in `publish_binary.yaml`

**Issue file:** `issues/bug-publish-binary-stub-gh-release.md`

The `publish_binary.yaml` workflow has a `gh_release` job that contains only `echo "not implemented"`. The real binary publishing job (`binary`) is set to wait for this stub before running. This means on every tag push, a no-op job runs first and then binaries are published — but no GitHub Release is created by this workflow.

It is unclear whether this stub was a placeholder that was never finished, or whether the GitHub Release is intended to be created by a different mechanism (`prerelease.yaml` + `release.yaml`). This needs a conversation with the release-process owners (Warashi, t-kikuc) before any code change is made.

---

## Files in This Folder

```
workflow-audit/
├── README.md                          ← you are here
└── issues/
    ├── bug-plugin-release-go-version-undefined.md
    ├── bug-publish-binary-stub-gh-release.md
    ├── security-code-butler-unpinned-abandoned.md
    ├── security-missing-permissions-declarations.md
    ├── maintenance-missing-concurrency-controls.md
    ├── maintenance-missing-timeouts.md
    ├── maintenance-action-pinning-inconsistency.md
    ├── maintenance-publish-tool-qemu-vs-native-arm.md
    └── maintenance-helm-push-runs-in-every-matrix-job.md
```

Each issue file is a standalone write-up of one problem: what it is, what the impact is, and what the recommended fix is.
