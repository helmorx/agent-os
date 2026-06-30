# HELMOR Commands

This is the full CLI reference for HELMOR Agent OS. The README keeps only the four primary commands so new users see the product quickly.

## Primary Commands

| Command | Purpose |
|---|---|
| `helmor install` | Initialize the project and generate agent adapters. |
| `helmor status` | Show compact project state. |
| `helmor dashboard` | Show the terminal dashboard. |
| `helmor doctor` | Run deterministic project checks. |

## General Commands

| Command | Purpose |
|---|---|
| `helmor help` | Show CLI help. |
| `helmor version` | Print the installed HELMOR version. |

## Project Setup

| Command | Purpose |
|---|---|
| `helmor init` | Create `.helmor/project.json` in the current project. |
| `helmor init --mode observe` | Initialize with warning-only behavior. |
| `helmor init --mode guard --force` | Enable stronger blocking behavior. |
| `helmor init --mode strict --force` | Enable release and security-sensitive enforcement. |
| `helmor uninstall` | Remove `.helmor` from the current project. |

## Development Workflow

| Command | Purpose |
|---|---|
| `helmor task start "feature work"` | Start a tracked task. |
| `helmor task finish` | Finish the active task and update closeout state. |
| `helmor checks` | List detected project checks. |
| `helmor handoff` | Write `.helmor/handoff.md` for the next agent. |
| `helmor reduce-tokens` | Show token-saving recommendations. |
| `helmor verify` | Alias for deterministic verification. |
| `helmor security` | Run the security detector pack. |

## Design Commands

| Command | Purpose |
|---|---|
| `helmor design init` | Add design detector guidance. |
| `helmor design audit` | Run UI/design drift checks. |
| `helmor design polish` | Show design polish recommendations. |

## Hook Entrypoints

Agents and adapters can call HELMOR through the hook interface:

```bash
helmor hook --event SessionStart
helmor hook --event UserPromptSubmit
helmor hook --event PreToolUse
helmor hook --event PostToolUse
helmor hook --event Stop
helmor hook --event PreCompact
helmor hook --event SessionEnd
```

## Modes

| Mode | Behavior |
|---|---|
| `observe` | Warn, route, summarize, and help adoption in existing projects. |
| `guard` | Block high-risk actions like secrets, destructive git, wrong runner, and unsafe deploys. |
| `strict` | Enforce checks, handoffs, closeout, and security review for release-sensitive work. |
