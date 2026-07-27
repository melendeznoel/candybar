# candybar

A Go CLI application that bundles a few unrelated feature areas behind command handlers:

- **Food** — recipe/ingredient domain logic backed by Postgres
- **World Health Org** — child nutrition data fetcher
- **Image** — image comparison service
- **Carp** — URL crawler (work in progress)

## Setup

```bash
go mod tidy
```

### Configuration

- `app.yaml` — App Engine flex config. Set `POSTGRES_CONNECTION` and `cloud_sql_instances` for your environment.

For local Postgres access via Cloud SQL, run the proxy (see `.vscode/tasks.json`):

```bash
cloud_sql_proxy -instances=candybar-208003:us-central1:food=tcp:5432
```

then point `POSTGRES_CONNECTION` at `localhost` with `sslmode=disable`.

## Build & run

```bash
go build ./...
go run ./cmd/candybar
```

Example CLI command:

```bash
go run ./cmd/candybar who-infant-nutrition USA
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
