# candybar

A Go CLI application that bundles a few unrelated feature areas behind command handlers:

- **Food** — openFDA food-recall lookup
- **World Health Org** — child nutrition and infant mortality data fetcher
- **Image** — image comparison service

`internal/carp` (a URL crawler) and the Postgres-backed recipe/ingredient logic in `internal/food` are present in the codebase but not yet wired up to a CLI command.

## Setup

```bash
go mod tidy
```

## Build & Run

```bash
go build ./...
go run ./cmd/candybar help
```

Available commands:

```bash
go run ./cmd/candybar who-infant-nutrition USA
go run ./cmd/candybar who-infant-deaths USA
go run ./cmd/candybar food-recalls peanut butter
go run ./cmd/candybar compare-images figures.json
cat figures.json | go run ./cmd/candybar compare-images
```

## Install & Invoke

```bash
go install ./cmd/candybar
```

```bash
candybar food-recalls peanuts
```

## Quality checks

```bash
gofmt -w .
go test ./...
go vet ./...
```

CI also runs `staticcheck` and `govulncheck` on every push and pull request.

## Debugging (macOS)

Install Delve by following the instructions here: https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md

VS Code is preconfigured (`.vscode/launch.json`) to launch `cmd/candybar/main.go` under `dlv` on port 2345.
