# Contributing

Thanks for helping improve HELMOR Agent OS.

## Development Setup

Requirements:

- Go 1.24 or newer
- Git

Run checks before opening a pull request:

```bash
go test ./...
go vet ./...
```

## Pull Request Guidelines

- Keep changes focused and easy to review.
- Use professional commit messages.
- Update docs when behavior, commands, policies, or public positioning changes.
- Do not commit generated binaries, release archives, secrets, or local `.helmor` state.
- Include tests for behavior changes.

## Documentation Changes

The README is marketing-first. Put deeper implementation details in `docs/` and link to them from the README.

## Security

Do not open public issues for vulnerabilities, leaked credentials, bypasses, or exploit details. Follow [SECURITY.md](SECURITY.md).
