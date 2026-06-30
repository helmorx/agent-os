# Security Policy

HELMOR Agent is built to reduce unsafe AI-assisted development behavior. Please report security issues responsibly.

## Reporting A Vulnerability

Use a private GitHub security advisory when available. If that is not available, contact HELMOR through:

- https://helmor.io
- https://x.com/helmorlabs

Do not publish exploit details, proof-of-concept bypasses, leaked credentials, or private project information in public issues.

## What To Report

- secret handling bypasses
- destructive command bypasses
- unsafe deploy or production-action bypasses
- installer integrity issues
- release artifact checksum problems
- vulnerabilities in generated adapters or hook behavior

## Supported Versions

| Version | Supported |
|---|---:|
| `0.1.x` | yes |

## Project Safety Rules

- Never commit `.env` files, private keys, credentials, or local `.helmor` runtime state.
- Keep generated binaries and release archives out of source control.
- Verify release artifacts through published checksums.
