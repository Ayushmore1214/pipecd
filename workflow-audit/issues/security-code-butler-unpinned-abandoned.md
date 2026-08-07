# Security: code-butler.yaml uses unpinned third-party action with API key access

**File:** `.github/workflows/code-butler.yaml`

## Problem

`ca-dp/code-butler@v1` is:
- **Not SHA-pinned**: the `v1` mutable tag means any push to the upstream repo moves what this action executes. Anyone who compromises that repo gets code execution in PipeCD CI with access to `OPENAI_API_KEY`.
- **Uses a deprecated model**: `gpt-4-1106-preview` was deprecated by OpenAI in early 2024. The action may already be broken.
- **Third-party, low-visibility**: `ca-dp` is not a well-known GitHub Actions publisher. The action has minimal community scrutiny.

## Impact

- `OPENAI_API_KEY` secret is passed to an unpinned, third-party action.
- `/review` and `/chat` commands on PRs may already be non-functional.

## Recommended Fix

Option A (preferred): Remove the workflow and revoke/rotate `OPENAI_API_KEY`. Check git history for last successful use — if it has not worked since the model was deprecated, it is already dead.

Option B: If AI review is actively used, replace with a custom `actions/github-script` step that calls the API directly (code lives in-repo, no third-party action trust required), and use a current model.
