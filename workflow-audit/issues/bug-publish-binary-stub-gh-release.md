# Bug: publish_binary.yaml has an unimplemented gh_release job that all binary publishing depends on

**File:** `.github/workflows/publish_binary.yaml`

## Problem

The `gh_release` job contains only `run: echo "not implemented"` but the `binary` job declares `needs: gh_release`. This means:

- On every tag push, the stub runs (doing nothing), then binary publishing proceeds.
- No GitHub Release is created by this workflow — yet binaries are attached to a release that must have been created by some other means.
- The intent is unclear: `prerelease.yaml` creates a draft release from the RELEASE file, and `release.yaml` creates the release PR. Neither explicitly creates the final GH Release on tag push.

## Open Questions

- Was `gh_release` meant to call `pipe-cd/actions-gh-release` on tag push (like `prerelease.yaml` does on PR)?
- Should the stub be removed entirely and the `needs:` dependency dropped?
- Or should it create the GitHub Release before binaries are uploaded?

## Recommended Action

Discuss with release-process owners. If the stub is dead code, remove it and remove `needs: gh_release` from the binary job.
