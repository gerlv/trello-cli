# Contributing to Trello CLI

Thank you for your interest in contributing! Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before participating.

## Development Setup

**Prerequisites:** Go 1.26 or later

```bash
git clone https://github.com/gerlv/trello-cli.git
cd trello-cli
go build -o bin/trello ./cmd/trello
go test -count=1 -race ./...
go vet ./...
```

## Making Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Ensure tests pass (`go test -count=1 -race ./...`)
5. Ensure vet passes (`go vet ./...`)
6. Commit your changes (see commit message guidelines below)
7. Push to your fork and open a pull request

## Commit Messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/). Use these prefixes:

- `feat:` — new feature
- `fix:` — bug fix
- `docs:` — documentation only
- `chore:` — maintenance (deps, CI, tooling)
- `test:` — adding or updating tests
- `refactor:` — code change that neither fixes a bug nor adds a feature

Keep the subject line under 72 characters. Use the body to explain what and why, not how.

## Pull Requests

- Describe what your PR does and why
- Reference any related issues (e.g., `Fixes #42`)
- Ensure CI passes before requesting review
- Keep PRs focused — one logical change per PR

## Reporting Issues

Use [GitHub Issues](https://github.com/gerlv/trello-cli/issues) and include:

- What you expected to happen
- What actually happened
- CLI version (`trello version`)
- Operating system and architecture
- Steps to reproduce
