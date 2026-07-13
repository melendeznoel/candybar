# candybar

A Go HTTP API deployed to Google App Engine (flexible environment). It bundles a few unrelated feature areas behind one router:

- **Food** — recipe/ingredient CRUD backed by Postgres
- **Social** — Twitter home/user timeline proxy
- **World Health Org** — child nutrition data proxy
- **Image** — image comparison endpoint
- **Carp** — URL crawler (work in progress)

## Setup

This repo predates Go modules — there is no `go.mod`. The repo root itself is the `GOPATH`, and the app lives at `src/main`, importing itself as `main/...`. Third-party dependencies are vendored under `src/github.com`, `src/golang.org`, `src/cloud.google.com` (gitignored), so fetch them into place before building:

```bash
export GOPATH="$(pwd)"
go get github.com/gorilla/mux github.com/lib/pq github.com/dghubble/go-twitter/twitter cloud.google.com/go/datastore
```

### Configuration

- `config.json` — Twitter API credentials (`ConsumerKey`, `ConsumerSecret`, `AccessToken`, etc.). Fill in the `<ADD>` placeholders.
- `src/main/app.yaml` — App Engine flex config. Set `POSTGRES_CONNECTION` and `cloud_sql_instances` for your environment.

For local Postgres access via Cloud SQL, run the proxy (see `.vscode/tasks.json`):

```bash
cloud_sql_proxy -instances=candybar-208003:us-central1:food=tcp:5432
```

then point `POSTGRES_CONNECTION` at `localhost` with `sslmode=disable`.

## Build & run

```bash
go build main
go run src/main/api.go
```

The server listens on `:8080`.

## Debugging (macOS)

Install Delve by following the instructions here: https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md

VS Code is preconfigured (`.vscode/launch.json`) to launch `src/main/api.go` under `dlv` on port 2345.
