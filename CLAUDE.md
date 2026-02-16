# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Rules

- NEVER read, fetch, or display real secrets/keys from 1Password or .env files. When testing, use dummy values only.

## Commands

```bash
make build          # Build binary to bin/export-key (CGO_ENABLED=0, injects version via ldflags)
make test           # Run all tests: go test -v ./...
go test -v -run TestParse ./internal/keyitem  # Run a single test
make lint           # golangci-lint run ./...
make fmt            # go fmt ./...
```

Test with the dotenv backend (no 1Password auth needed):
```bash
EK_BACKEND=dotenv EK_DOTENV_PATHS=/path/to/.env ./bin/export-key list
EK_BACKEND=dotenv EK_DOTENV_PATHS=/path/to/.env ./bin/export-key select 1
```

## Commit Messages

Use conventional commit format with imperative mood:

```
type: short description
```

Types: `feat`, `fix`, `docs`, `refactor`, `chore`, `test`, `ci`

Examples:
- `feat: add multi-select to TUI picker`
- `fix: preserve dotenv insertion order`
- `docs: update README install instructions`
- `refactor: extract export logic from select command`
- `test: add dotenv edge case coverage`

## Architecture

See [README.md](README.md) for user-facing docs (install, usage, key naming, backends, configuration).

### Stdout/stderr contract

This is the most important invariant: `select.go` prints `export VAR="value"` to stdout and status messages to stderr. The shell wrapper in `internal/shell/` evals stdout. Breaking this (e.g., printing a log to stdout) will inject garbage into the user's shell.

### Data flow

`select.go` orchestrates everything: loads config → creates backend → lists items → parses into KeyItems → routes to number/query/project/TUI picker → fetches secret → prints export.

### Backend interface

`backend.Backend` has two methods: `ListItems() ([]string, error)` and `GetSecret(title string) (string, error)`. Implementations: `OnePassword` (shells out to `op` CLI) and `Dotenv` (reads .env files). `registry.go` maps config strings to implementations.

### Key naming parsing

`keyitem.Parse()` is the source of truth for splitting item titles on the first `-` into env var name and project.

### Version injection

`Version`, `Commit`, `BuildDate` are package-level vars in `internal/cmd/version.go`, injected via `-ldflags -X` at build time. The Makefile, `.goreleaser.yaml`, and GitHub Actions all use this pattern.
