# Contributing Guide

Thanks for your interest in contributing.

## Development Setup

1. Install Go 1.22+.
2. Fork and clone the repository.
3. Run:

```bash
go mod tidy
go test ./...
go vet ./...
```

## Coding Rules

- Keep changes focused and small.
- Add or update tests for behavior changes.
- Prefer explicit error handling over panic.
- Keep public API changes documented in `README.md` and `CHANGELOG.md`.

## Pull Request Checklist

- [ ] Tests pass locally: `go test ./...`
- [ ] Static checks pass locally: `go vet ./...`
- [ ] Documentation updated for behavior/API changes
- [ ] Changelog updated when needed

## Commit Message Convention

Recommended prefixes:
- `feat:` new functionality
- `fix:` bug fixes
- `refactor:` non-behavioral refactor
- `test:` test updates
- `docs:` documentation changes
- `ci:` workflow and automation updates
