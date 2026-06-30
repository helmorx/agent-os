# Changelog

All notable changes to HELMOR Agent are documented here.

## Unreleased

- Renamed public product branding to HELMOR Agent.
- Renamed the scoped npm alias package metadata to `@helmoragent/agent`.

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
