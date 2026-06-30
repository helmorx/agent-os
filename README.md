<p align="center">
  <img src="assets/helmor-hero.svg" alt="HELMOR DevSuite banner" width="100%">
</p>

<p align="center">
  <a href="https://github.com/helmorx/devsuite/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/helmorx/devsuite/ci.yml?branch=main&label=ci&style=for-the-badge" alt="CI"></a>
  <a href="https://github.com/helmorx/devsuite/blob/main/LICENSE"><img src="https://img.shields.io/github/license/helmorx/devsuite?style=for-the-badge" alt="Apache-2.0 license"></a>
  <a href="https://github.com/helmorx/devsuite/releases"><img src="https://img.shields.io/github/v/release/helmorx/devsuite?style=for-the-badge&include_prereleases" alt="Latest release"></a>
  <a href="https://github.com/helmorx/devsuite/stargazers"><img src="https://img.shields.io/github/stars/helmorx/devsuite?style=for-the-badge" alt="GitHub stars"></a>
</p>

<p align="center">
  <b>A local-first development engine for AI-assisted coding.</b><br>
  Reduce wasted tokens, project drift, hallucinated changes, and unsafe agent actions.
</p>

---

## Why HELMOR

AI coding agents are powerful, but they waste context, forget project rules, invent files/APIs, run noisy commands, and drift away from your product intent. HELMOR gives every project a local operating layer:

| Problem | HELMOR response |
|---|---|
| Agents re-read the same repo context | Compact context cards, session state, handoffs |
| Shell/test output burns tokens | `rtk` first, `sqz` fallback, concise checks |
| Agents use the wrong package manager | package-runner detection and guard rails |
| Risky actions happen too early | observe, guard, and strict project modes |
| UI starts looking generic | deterministic design detectors inspired by Impeccable-style checks |
| New sessions lose decisions | `.helmor/project.json`, task state, handoff memory |

## Install

macOS:

```bash
brew tap helmorx/devsuite https://github.com/helmorx/devsuite
brew install helmor
```

Windows:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/helmorx/devsuite/main/install/install.ps1 | iex"
```

Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/helmorx/devsuite/main/install/install.sh | sh
```

## Quick Start

```bash
helmor install
helmor status
helmor doctor
helmor dashboard
```

Existing projects start in `observe` mode. HELMOR warns, routes, and summarizes without blocking most actions. Move to enforcement when you are ready:

```bash
helmor init --mode guard --force
helmor init --mode strict --force
```

<p align="center">
  <img src="assets/terminal-preview.svg" alt="HELMOR terminal dashboard preview" width="88%">
</p>

## What You Get

**Token Reduction Engine**

- Prefer `rtk` for compressed shell, git, and test output.
- Use `sqz` as fallback for compression and handoff support.
- Detect `nub`, npm, pnpm, yarn, bun, Go, Rust, and Python runners.
- Push agents toward graph/code discovery before broad repo scans.

**Anti-Drift Engine**

- Detect truth files like `PRODUCT.md`, `DESIGN.md`, `ARCHITECTURE.md`, `AGENTS.md`, `CLAUDE.md`, `README.md`, `PRD.md`, and `TRD.md`.
- Preserve `.helmor/context-card.md` and `.helmor/handoff.md`.
- Track task stage, touched files, pending checks, and risky closeout state.

**Agent Adapters**

- Codex and Claude Code: hook-compatible command entrypoints.
- Cursor and Windsurf: generated project rules.
- Other agents: use the same `helmor hook --event <EventName>` interface.

**Detector Packs**

- secret-shaped filenames
- destructive git commands
- package-runner bypass
- unsafe production/mainnet/provider commands
- stale or missing truth files
- token-tool gaps
- UI/design drift patterns

## Modes

| Mode | Best for | Behavior |
|---|---|---|
| `observe` | existing projects, onboarding | warns and routes, does not surprise-block |
| `guard` | active development | blocks secrets, destructive git, runner bypass, unsafe deploy actions |
| `strict` | launch/security-sensitive work | enforces closeout, checks, handoffs, and security review |

## Commands

```bash
helmor init
helmor install
helmor uninstall
helmor status
helmor doctor
helmor dashboard
helmor task start "feature work"
helmor task finish
helmor checks
helmor handoff
helmor reduce-tokens
helmor verify
helmor security
helmor design init
helmor design audit
helmor design polish
```

Hook-compatible commands:

```bash
helmor hook --event SessionStart
helmor hook --event UserPromptSubmit
helmor hook --event PreToolUse
helmor hook --event PostToolUse
helmor hook --event Stop
helmor hook --event PreCompact
helmor hook --event SessionEnd
```

## Project Profile

HELMOR stores project rules locally:

```text
.helmor/project.json
```

The profile tracks stack, runner, truth files, checks, policies, enabled adapters, and built-in skill modules. It is plain JSON and can be reviewed or edited by the developer.

## Built For

- solo developers shipping with AI agents
- teams that want safer AI coding workflows
- vibe coders who need less hallucination and more structure
- high-risk apps that need launch/security discipline
- frontend teams that want deterministic UI polish checks

## Release Notes

V1 is local-first and does not upload source code, require an account, or send telemetry.

Homebrew SHA placeholders in `Formula/helmor.rb` are replaced after the first GitHub Release is published.

## License

Apache-2.0. See [LICENSE](LICENSE).
