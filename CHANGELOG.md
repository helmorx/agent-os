# Changelog

All notable changes to HELMOR Agent are documented here.

## 0.1.6 - 2026-07-02

- Updated all repository URL references from `helmorx/agent-os` to `helmorx/helmoragent` (badges, install scripts, Homebrew formula, npm package metadata, docs) to match the renamed GitHub repository. No logic or functionality changed. The Go module path is intentionally left as `github.com/helmorx/agent-os`, which continues to resolve via GitHub's automatic redirect.

## 0.1.5 - 2026-07-02

- Re-release of 0.1.4: no functional changes. The 0.1.4 tag was moved to amend commit message metadata after the version had already published, which correctly blocked a same-version republish on npm. 0.1.5 ships that identical content cleanly.

## 0.1.4 - 2026-07-02

- Renamed public product branding to HELMOR Agent.
- Renamed the scoped npm alias package metadata to `@helmoragent/agent`.
- Fixed PreToolUse/PostToolUse/Stop hook output shapes that Codex's stricter hook schema was rejecting (missing `hookEventName` tag, unsupported `suppressOutput` on those two events, unsupported `permissionDecision:allow` and `decision:approve` values) — every allowed tool call and clean session close was failing under Codex.
- Fixed `helmor install`/`helmor doctor` merging and checking Claude Code's global hooks at `~/.claude/settings.local.json`, a path Claude Code never reads at the user level; the correct target is `~/.claude/settings.json`. The global hook merge was silently a no-op for every Claude Code user until now.
- Fixed legacy-hook migration detection using a hardcoded machine-specific path instead of a home-relative one.
- Reduced false-positive secret-path detector findings: stricter word-boundary matching, `.gitignore`-aware file scanning, and no more flagging scripts/tests that merely mention secret-related terms in their names.

## 0.1.3 - 2026-06-30

- Switched npm distribution to `@helmoragent/helmor` because npm blocks the unscoped `helmor` package name as too similar to an existing package.
- Kept the installed CLI command as `helmor`.
- Added scoped npm alias package metadata.

## 0.1.2 - 2026-06-30

- Added npm distribution under the short package name `helmor`.
- Added npm-family one-off commands for `npx`, `pnpm dlx`, `yarn dlx`, and `bunx`.
- Added scoped npm alias package metadata.
- Added npm package tests and release checks for Go/npm version alignment.

## 0.1.1 - 2026-06-30

- Reframed public repository positioning around the local agent watcher and end-to-end product development lifecycle.
- Added dedicated docs for commands, skills, quickstart, security, and contributing.
- Kept the README marketing-first with deeper references moved into `docs/`.
- Prepared source-built Homebrew packaging under the `helmoragent` formula name.

## 0.1.0 - 2026-06-30

- Initial public release of HELMOR Agent.
- Added single-binary CLI for AI-assisted development workflows.
- Added `observe`, `guard`, and `strict` project modes.
- Added Codex, Claude Code, Cursor, and Windsurf adapters.
- Added detector packs for secrets, unsafe commands, package-runner drift, and design drift.
- Added macOS Homebrew, Windows PowerShell, and Linux shell installers.
