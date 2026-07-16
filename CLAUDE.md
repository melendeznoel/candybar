# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Go CLI application ("candybar") that bundles several unrelated feature areas in one binary: Twitter integration logic, food/recipe domain logic backed by Postgres, a WHO child-nutrition data fetcher, image comparison services, and a web crawler ("Carp"). Treat these as independent modules that happen to share a process — changes in one package rarely need to touch another.

## Build system: Go modules

Go modules are enabled with the module root at `src/main` so existing `main/...` imports continue to work.

To work on this repo:

```bash
cd src/main
go mod tidy
go build ./...
go run . help
```

The entry point is `src/main/main.go`; the VS Code debug config (`.vscode/launch.json`) launches it directly with `dlv` on port 2345 (see `src/main/documents/delve-settings.txt`).

There is a starter test suite in `src/main/worldHealthOrg` and CI runs `gofmt`, `go test`, `go vet`, `staticcheck`, and `govulncheck`.

## Configuration

- `src/main/app.yaml` — App Engine flex config. `POSTGRES_CONNECTION` env var and `cloud_sql_instances` are placeholders (`<ADD>`) that must be filled in per environment; for local dev against Cloud SQL, use the `cloud sql proxy` task in `.vscode/tasks.json` (`cloud_sql_proxy -instances=candybar-208003:us-central1:food=tcp:5432`) and point `POSTGRES_CONNECTION` at `localhost` with `sslmode=disable`.
- `config.json` (repo root) — Twitter API credentials (`ConsumerKey`, `ConsumerSecret`, `AccessToken`, etc.), currently placeholders.
- The App Engine bootstrap (`appengine.Main()`, App Engine context in `main.go`) is commented out — the app currently runs as a CLI, not through the App Engine runtime.

## Architecture

**CLI entrypoint**: `src/main/main.go` is the command dispatcher. Add a command by extending the switch in `main()` and calling the relevant service-layer function.

**Per-feature package layout** (see `food/` as the fullest example): `Models.go` (request/response and domain types) → `*Controller.go` (HTTP handlers: decode request, call service, encode response) → `*Service.go` (business logic) → `Repository.go` / `Database.go` (SQL access via `database/sql` + `lib/pq`, raw connection built from `POSTGRES_CONNECTION`). Not every package has all layers (e.g. `carp` and `worldHealthOrg` skip the repository layer since they call external HTTP/scraping sources instead of a DB).

**Persistence helpers**: `repositories/PostgresRepo.go` (raw Postgres via `lib/pq`) is present as a thin, mostly-unused helper — most feature code (e.g. `food/Database.go`) opens its own `*sql.DB` directly rather than going through `repositories/`.

**Social (Twitter)** (`src/main/social/`): wraps `github.com/dghubble/go-twitter` to fetch home/user timelines and remaps the SDK's tweet structs into this app's own leaner `Tweet`/`User`/`Hashtag`/`TweetURL` models (`Models.go`), extracting media/video URLs, hashtags, and @-mentions by hand from tweet entities.

**Carp** (`src/main/carp/`): a URL crawler/scraper (`CarpCrawlService.go`, `CarpAnchorService.go`) — currently stubbed (`Crawl` handler has a `// TODO: REPLACE` and always calls into `crawl(urls)`).

**World Health Org** (`src/main/worldHealthOrg/`): proxies WHO child-nutrition data through `ChildHealthService.go`; used by both its own endpoint and pulled into `social.RunTasks()`.

**Helper package** (`src/main/helper/`): cross-cutting utilities — routing (`HelperService.go`), generic HTTP utility calls (`HttpUtilityService.go`), background task scheduling (`TaskService.go`), and misc utilities (`UtilityService.go`). Shared `Route`/`Routes` types live in `helper/Models.go`.
