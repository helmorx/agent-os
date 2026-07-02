# HELMOR Agent

Npm installer for the native `helmor` CLI.

```bash
npm i -g @helmoragent/helmor
helmor install
```

One-off usage without a global install:

```bash
npx @helmoragent/helmor@latest install
pnpm dlx @helmoragent/helmor install
yarn dlx @helmoragent/helmor install
bunx @helmoragent/helmor install
```

This package downloads the matching HELMOR release binary from GitHub, verifies it against the published SHA-256 checksum file, and exposes the `helmor` command.

`helmor install` initializes project files and merges Codex/Claude global hook entries. Use `helmor install --project-only` to write only project-local files and generated adapters.

Set `HELMOR_VERSION`, `HELMOR_REPO`, or `HELMOR_NPM_QUIET=1` to override the release tag, source repository, or install logging.
