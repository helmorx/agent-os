# HELMOR Agent OS

Npm installer for the native `helmor` CLI.

```bash
npm i -g helmor
helmor install
```

One-off usage without a global install:

```bash
npx helmor@latest install
pnpm dlx helmor install
yarn dlx helmor install
bunx helmor install
```

This package downloads the matching HELMOR release binary from GitHub, verifies it against the published SHA-256 checksum file, and exposes the `helmor` command.

Set `HELMOR_VERSION`, `HELMOR_REPO`, or `HELMOR_NPM_QUIET=1` to override the release tag, source repository, or install logging.
