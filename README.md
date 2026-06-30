<p align="center">
  <img src="assets/helmor-hero.svg" alt="HELMOR Agent OS banner" width="100%">
</p>

<p align="center">
  <a href="https://github.com/helmorx/agent-os/actions/workflows/ci.yml"><img src="https://img.shields.io/github/actions/workflow/status/helmorx/agent-os/ci.yml?branch=main&label=ci&style=flat-square" alt="CI"></a>
  <a href="https://github.com/helmorx/agent-os/releases"><img src="https://img.shields.io/github/v/release/helmorx/agent-os?style=flat-square&include_prereleases" alt="Latest release"></a>
  <a href="https://github.com/helmorx/agent-os/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-Apache--2.0-0E1724?style=flat-square" alt="Apache-2.0 license"></a>
  <a href="#install"><img src="https://img.shields.io/badge/install-brew%20%7C%20powershell%20%7C%20curl-19C37D?style=flat-square" alt="Install"></a>
</p>

<h1 align="center">The local operating layer for AI coding agents</h1>

<p align="center">
  HELMOR helps Codex, Claude Code, Cursor, and Windsurf spend fewer tokens, remember repo context, avoid drift, and stop unsafe actions before they ship.
</p>

<p align="center">
  <a href="#install"><b>Install</b></a>
  ·
  <a href="https://helmor.io"><b>Website</b></a>
  ·
  <a href="https://x.com/helmorlabs"><b>X</b></a>
  ·
  <a href="#why-helmor"><b>Why HELMOR</b></a>
  ·
  <a href="#agent-support"><b>Agent support</b></a>
  ·
  <a href="#commands"><b>Commands</b></a>
</p>

---

## Why HELMOR

AI agents are fast, but they often waste tokens rediscovering the project, invent missing APIs, forget earlier decisions, run the wrong commands, or drift away from product truth. HELMOR gives each repository a local operating layer for safer AI-assisted development.

<table>
  <tr>
    <td width="33%">
      <h3>Reduce wasted tokens</h3>
      <p>Prefer compact shell/test output, context cards, handoffs, and graph-first discovery instead of repeated repo scans.</p>
    </td>
    <td width="33%">
      <h3>Stop project drift</h3>
      <p>Keep agents aligned to truth files, package runners, checks, policies, and task state in <code>.helmor/project.json</code>.</p>
    </td>
    <td width="33%">
      <h3>Guard risky actions</h3>
      <p>Detect secrets, destructive git, package-runner bypass, unsafe deploys, and launch/security closeout gaps.</p>
    </td>
  </tr>
</table>

## Install

```bash
brew tap helmorx/agent-os https://github.com/helmorx/agent-os
brew install helmor
```

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/helmorx/agent-os/main/install/install.ps1 | iex"
```

```bash
curl -fsSL https://raw.githubusercontent.com/helmorx/agent-os/main/install/install.sh | sh
```

## First Run

```bash
helmor install
helmor status
helmor doctor
helmor dashboard
```

Existing projects start in `observe` mode, so HELMOR warns and routes without surprise-blocking your workflow.

```bash
helmor init --mode guard --force
helmor init --mode strict --force
```

<p align="center">
  <img src="assets/terminal-preview.svg" alt="HELMOR terminal dashboard preview" width="88%">
</p>

## What It Adds To A Project

```text
.helmor/
  project.json          repo profile, checks, policies, tools, adapters
  context-card.md       compact context for new sessions
  handoff.md            closeout summary for the next agent
  state.json            local runtime state, ignored by git
```

HELMOR is local-first. It does not require an account, upload your source, or send telemetry.

## Agent Support

| Agent | V1 support | Integration style |
|---|---:|---|
| Codex | yes | hook-compatible command entrypoints |
| Claude Code | yes | hook-compatible command entrypoints |
| Cursor | yes | generated project rules |
| Windsurf | yes | generated project rules |
| Other agents | compatible | use `helmor hook --event <EventName>` |

## Modes

| Mode | Use it when | Behavior |
|---|---|---|
| `observe` | adopting HELMOR in an existing repo | warn, route, summarize |
| `guard` | active development with agents | block secrets, destructive git, wrong runner, unsafe deploys |
| `strict` | release, launch, security-sensitive work | enforce checks, handoffs, closeout, security review |

## Detector Packs

<table>
  <tr>
    <td><b>Secrets</b><br>secret-shaped filenames and unsafe paths</td>
    <td><b>Shell/Git</b><br>destructive git and unsafe deploy commands</td>
    <td><b>Runner Drift</b><br>wrong package manager and noisy retries</td>
  </tr>
  <tr>
    <td><b>Truth Files</b><br>missing project authority docs</td>
    <td><b>Token Waste</b><br>missing or unused compression/discovery tools</td>
    <td><b>Design Drift</b><br>generic AI UI patterns inspired by modern design audits</td>
  </tr>
</table>

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

Hook-compatible entrypoints:

```bash
helmor hook --event SessionStart
helmor hook --event UserPromptSubmit
helmor hook --event PreToolUse
helmor hook --event PostToolUse
helmor hook --event Stop
helmor hook --event PreCompact
helmor hook --event SessionEnd
```

## Built For

- developers shipping real projects with AI agents
- teams that want repeatable AI coding workflows
- vibe coders who need less hallucination and more structure
- frontend teams that want deterministic UI polish checks
- high-risk apps that need launch and security discipline

## License

Apache-2.0. See [LICENSE](LICENSE).
