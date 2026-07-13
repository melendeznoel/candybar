# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Go HTTP API ("candybar") deployed to Google App Engine flexible environment. It bundles several unrelated feature areas behind one router: a Twitter timeline proxy, a food/recipe CRUD API backed by Postgres, a WHO child-nutrition data proxy, an image comparison endpoint, and a web crawler ("Carp"). Treat these as independent modules that happen to share a process — changes in one package rarely need to touch another.

## Build system: GOPATH, not Go modules

There is no `go.mod`. The repo root itself is the `GOPATH`, and the app lives at `src/main`, importing itself as `main/...` (e.g. `main/endpoints`, `main/food`). Third-party deps (`github.com/gorilla/mux`, `github.com/lib/pq`, `github.com/dghubble/go-twitter`, `cloud.google.com/go/datastore`, etc.) are expected to be vendored under `src/github.com`, `src/golang.org`, `src/cloud.google.com` — these are gitignored, so a fresh checkout needs them fetched into place before anything builds.

To work on this repo, set `GOPATH` to the repo root first:

```bash
export GOPATH="$(pwd)"          # from repo root
go get github.com/gorilla/mux github.com/lib/pq github.com/dghubble/go-twitter/twitter cloud.google.com/go/datastore ...
go build main
go run src/main/api.go
```

The entry point is `src/main/api.go`; the VS Code debug config (`.vscode/launch.json`) launches it directly with `dlv` on port 2345 (see `src/main/documents/delve-settings.txt`).

There are no test files, no linter config, and no Makefile in the repo currently — there is nothing to run beyond `go build`/`go vet` until tests are added.

## Configuration

- `src/main/app.yaml` — App Engine flex config. `POSTGRES_CONNECTION` env var and `cloud_sql_instances` are placeholders (`<ADD>`) that must be filled in per environment; for local dev against Cloud SQL, use the `cloud sql proxy` task in `.vscode/tasks.json` (`cloud_sql_proxy -instances=candybar-208003:us-central1:food=tcp:5432`) and point `POSTGRES_CONNECTION` at `localhost` with `sslmode=disable`.
- `config.json` (repo root) — Twitter API credentials (`ConsumerKey`, `ConsumerSecret`, `AccessToken`, etc.), currently placeholders.
- The App Engine bootstrap (`appengine.Main()`, App Engine context in `api.go`) is commented out — the app currently runs as a plain `net/http` server on `:8080` via `http.ListenAndServe`, not through the App Engine runtime.

## Architecture

**Routing**: `api.go`'s `BuildRoutes()` returns a flat `helper.Routes` slice mapping method+path to a handler function pulled from each feature package's Controller. `helper.BuildRouter` (`src/main/helper/HelperService.go`) turns that into a `gorilla/mux` router, wrapping every handler in a logging middleware (`helper.Logger`). To add an endpoint: write the handler in the relevant feature package's `*Controller.go`, then register it in `BuildRoutes()`.

**Per-feature package layout** (see `food/` as the fullest example): `Models.go` (request/response and domain types) → `*Controller.go` (HTTP handlers: decode request, call service, encode response) → `*Service.go` (business logic) → `Repository.go` / `Database.go` (SQL access via `database/sql` + `lib/pq`, raw connection built from `POSTGRES_CONNECTION`). Not every package has all layers (e.g. `carp` and `worldHealthOrg` skip the repository layer since they call external HTTP/scraping sources instead of a DB).

**Two persistence backends coexist**: `repositories/PostgresRepo.go` (raw Postgres via `lib/pq`) and `repositories/DatastoreRepo.go` (GCP Datastore, via `cloud.google.com/go/datastore`) are both present as thin, mostly-unused helpers — most feature code (e.g. `food/Database.go`) opens its own `*sql.DB` directly rather than going through `repositories/`.

**Social (Twitter)** (`src/main/social/`): wraps `github.com/dghubble/go-twitter` to fetch home/user timelines and remaps the SDK's tweet structs into this app's own leaner `Tweet`/`User`/`Hashtag`/`TweetURL` models (`Models.go`), extracting media/video URLs, hashtags, and @-mentions by hand from tweet entities.

**Carp** (`src/main/carp/`): a URL crawler/scraper (`CarpCrawlService.go`, `CarpAnchorService.go`) — currently stubbed (`Crawl` handler has a `// TODO: REPLACE` and always calls into `crawl(urls)`).

**World Health Org** (`src/main/worldHealthOrg/`): proxies WHO child-nutrition data through `ChildHealthService.go`; used by both its own endpoint and pulled into `social.RunTasks()`.

**Helper package** (`src/main/helper/`): cross-cutting utilities — routing (`HelperService.go`), generic HTTP utility calls (`HttpUtilityService.go`), background task scheduling (`TaskService.go`), and misc utilities (`UtilityService.go`). Shared `Route`/`Routes` types live in `helper/Models.go`.
