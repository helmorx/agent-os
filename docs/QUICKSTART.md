# Quickstart

HELMOR Agent is designed for existing repositories and new product builds. It starts in `observe` mode so teams can adopt it without surprise blocking.

## 1. Install HELMOR

Node.js:

```bash
npm i -g @helmoragent/helmor
```

One-off:

```bash
npx @helmoragent/helmor@latest install
pnpm dlx @helmoragent/helmor install
yarn dlx @helmoragent/helmor install
bunx @helmoragent/helmor install
```

macOS Homebrew:

```bash
brew install helmorx/tap/helmoragent
```

Windows:

```powershell
irm https://raw.githubusercontent.com/helmorx/agent-os/main/install/install.ps1 | iex
```

Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/helmorx/agent-os/main/install/install.sh | sh
```

## 2. Enable A Project

Run this inside the project root:

```bash
helmor install
```

HELMOR creates `.helmor/project.json`, generated agent adapters, and local runtime files for context, state, and handoff.

## 3. Check The Project

```bash
helmor status
helmor dashboard
helmor doctor
```

## 4. Choose Enforcement

| Mode | Use it when |
|---|---|
| `observe` | You are adopting HELMOR in an existing project. |
| `guard` | You want HELMOR to block common unsafe actions. |
| `strict` | You are working on release, launch, or security-sensitive tasks. |

```bash
helmor init --mode guard --force
helmor init --mode strict --force
```

## 5. Work With Agents

HELMOR supports Codex, Claude Code, Cursor, and Windsurf in v1. Unsupported agents can still call the hook-compatible command interface documented in [Commands](COMMANDS.md).
