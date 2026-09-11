# candybar

A Go CLI application that bundles a few unrelated feature areas behind command handlers:

- **Food** — openFDA food-recall lookup
- **Drug** — openFDA drug-recall lookup
- **World Health Org** — child nutrition and infant mortality data fetcher
- **Image** — image comparison service

`internal/communication` (a URL crawler/scraper), the Postgres-backed recipe/ingredient logic in `internal/food`, and the Qdrant vector store client in `internal/repositories` are present in the codebase but not yet wired up to a CLI command.

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
go run ./cmd/candybar drug-recalls ibuprofen
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

## Vector search (Qdrant)

`internal/repositories/qdrant.go` wraps the official [`github.com/qdrant/go-client`](https://github.com/qdrant/go-client) gRPC client so feature packages can store and query embeddings without building protobuf requests by hand. It is a library helper — there is no CLI command for it yet.

Point it at a Qdrant instance (note the **gRPC** port `6334`, not the REST port `6333`):

```bash
# Local, e.g. docker run -p 6333:6333 -p 6334:6334 qdrant/qdrant
export QDRANT_HOST="localhost"
export QDRANT_PORT="6334"

# Or Qdrant Cloud — an https URL turns TLS on automatically
export QDRANT_URL="https://[CLUSTER_ID].[REGION].cloud.qdrant.io:6334"
export QDRANT_API_KEY="[API_KEY]"
```

```go
package main

import (
	"context"
	"log"

	"candybar/internal/repositories"

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	ctx := context.Background()

	// Reads QDRANT_URL (or QDRANT_HOST/QDRANT_PORT) and QDRANT_API_KEY.
	store, err := repositories.QdrantFromEnv()
	if err != nil {
		log.Fatalf("qdrant: %v", err)
	}

	defer store.Close()

	// Create the collection once; this is a no-op if it already exists.
	if err := store.CreateCollection(ctx, "recipes", 1536, qdrant.Distance_Cosine); err != nil {
		log.Fatalf("create collection: %v", err)
	}

	// embedding and query come from whatever embedding model you use.
	var embedding, query []float32

	points := []repositories.QdrantPoint{
		{
			ID:      uint64(1),
			Vector:  embedding,
			Payload: map[string]interface{}{"title": "Pad Thai", "cuisine": "thai"},
		},
	}

	if err := store.Upsert(ctx, "recipes", points); err != nil {
		log.Fatalf("upsert: %v", err)
	}

	hits, err := store.Search(ctx, "recipes", repositories.QdrantSearchRequest{
		Vector:      query,
		Limit:       5,
		WithPayload: true,
		Filter: &qdrant.Filter{
			Must: []*qdrant.Condition{qdrant.NewMatchKeyword("cuisine", "thai")},
		},
	})
	if err != nil {
		log.Fatalf("search: %v", err)
	}

	for _, hit := range hits {
		log.Printf("%v scored %.3f: %v", hit.ID, hit.Score, hit.Payload["title"])
	}
}
```

Other operations on the same value:

```go
store.Retrieve(ctx, "recipes", []interface{}{uint64(1)})    // fetch points by ID
store.Count(ctx, "recipes", nil)                            // exact count, nil filter counts all
store.Scroll(ctx, "recipes", 100, nil, nil)                 // page through; returns the next offset
store.Delete(ctx, "recipes", []interface{}{uint64(1)})      // delete by ID
store.DeleteByFilter(ctx, "recipes", filter)                // delete by filter (nil is rejected)
store.DeleteCollection(ctx, "recipes")
```

Notes:

- Point IDs must be an unsigned integer or a UUID string — the only two forms Qdrant accepts.
- Filters are the client library's `*qdrant.Filter`, built with the `qdrant.NewMatch*`, `qdrant.NewRange` and related helpers.
- Requests get a 30s deadline unless the context already carries one; change it with `store.WithTimeout(d)`.
- `store.Client()` exposes the underlying `*qdrant.Client` for anything the wrapper does not cover (snapshots, aliases, grouped or batch queries).

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
