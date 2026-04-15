# zotero-go-client

A Go client library for the [Zotero Web API v3](https://www.zotero.org/support/dev/web_api/v3/basics).

## Features

- Full coverage of the Zotero Web API v3 (items, collections, searches, tags, groups, full-text, schema, keys, deleted content)
- Automatic pagination with Go 1.24 iterators (`iter.Seq2`)
- WebSocket streaming via the [Zotero Streaming API](https://www.zotero.org/support/dev/web_api/v3/streaming_api)
- Functional options for request parameters and client configuration
- Typed error helpers (`IsNotFound`, `IsConflict`, `IsRateLimited`, `IsPreconditionFailed`)
- Respects `Backoff` and `Retry-After` headers

## Requirements

- Go 1.24 or later

## Installation

```sh
go get github.com/abtris/zotero-go-client/zotero
```

## Getting your API key and user ID

1. **Create a Zotero account** at [zotero.org/user/register](https://www.zotero.org/user/register) if you don't have one.
2. Go to [zotero.org/settings/keys](https://www.zotero.org/settings/keys).
3. Your **user ID** is displayed at the top of the page (a numeric value — this is different from your username).
4. Click **Create new private key** and configure permissions:
   - **Personal Library** — check *Allow library access* (and *Allow write access* if your application needs to create or modify data).
   - **Default Group Permissions** — choose the level of access for group libraries.
   - Give the key a descriptive name (e.g. "My Go app") so you can identify it later.
5. Click **Save Key**. Copy the generated key — it will only be shown once.

Use the key and user ID in your code:

```go
client := zotero.NewClient("P9NiFoyLeZu2bZNvvuQPDWsd")          // explicit API key
lib    := zotero.UserLibrary("12345")                             // user ID
```

Or set the `ZOTERO_API_KEY` environment variable and pass an empty string — the client picks it up automatically:

```sh
export ZOTERO_API_KEY="P9NiFoyLeZu2bZNvvuQPDWsd"
```

```go
client := zotero.NewClient("")   // reads ZOTERO_API_KEY from the environment
```

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/abtris/zotero-go-client/zotero"
)

func main() {
    client := zotero.NewClient("your-api-key")
    lib := zotero.UserLibrary("your-user-id")
    ctx := context.Background()

    // List the 5 most recently modified items
    items, resp, err := client.Items.ListTop(ctx, lib,
        zotero.WithLimit(5),
        zotero.WithSort("dateModified"),
        zotero.WithDirection("desc"),
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Total items: %d (library version %d)\n", resp.TotalResults, resp.LastVersion)
    for _, item := range items {
        fmt.Printf("  %s — %s\n", item.Data.Title, item.Data.ItemType)
    }
}
```

See [`example/basic/main.go`](example/basic/main.go) for a more complete runnable example:

```sh
export ZOTERO_API_KEY="your-api-key"
export ZOTERO_USER_ID="your-user-id"
go run ./example/basic
```

## Usage

### Client initialization

```go
// Explicit API key
client := zotero.NewClient("api-key")

// Read API key from ZOTERO_API_KEY environment variable
client := zotero.NewClient("")

// With options
client := zotero.NewClient("",
    zotero.WithBaseURL("https://custom-endpoint.example.com"),
    zotero.WithHTTPClient(myHTTPClient),
    zotero.WithUserAgent("my-app/1.0"),
)
```

### Libraries

Every data operation requires a `LibraryID` to scope the request to a user or group library:

```go
userLib  := zotero.UserLibrary("12345")
groupLib := zotero.GroupLibrary("67890")
```

### Services

| Service              | Description                              |
| -------------------- | ---------------------------------------- |
| `client.Items`       | CRUD for items, top-level, trash, children, publications |
| `client.Collections` | CRUD for collections and subcollections  |
| `client.Searches`    | Saved search CRUD                        |
| `client.Tags`        | Tag listing, per-item/collection tags, bulk delete |
| `client.Schema`      | Item types, fields, creator types, templates |
| `client.Keys`        | API key info and deletion                |
| `client.Groups`      | List user groups                         |
| `client.FullText`    | Full-text content indexing               |
| `client.Deleted`     | Deleted content since a version          |

### Pagination

All `List` methods accept functional options for pagination:

```go
items, resp, err := client.Items.List(ctx, lib,
    zotero.WithLimit(50),
    zotero.WithStart(100),
    zotero.WithSort("title"),
)
// resp.TotalResults — total matching items
// resp.Links.Next   — URL of the next page (empty if last page)
```

Use `ListAll` to iterate over every result automatically:

```go
for item, err := range client.Items.ListAll(ctx, lib) {
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(item.Data.Title)
}
```

### Error handling

```go
_, _, err := client.Items.Get(ctx, lib, "BADKEY")
if zotero.IsNotFound(err) {
    // 404
} else if zotero.IsRateLimited(err) {
    // 429 — check err.(*zotero.APIError).RetryAfter
} else if zotero.IsPreconditionFailed(err) {
    // 412 — version conflict
}
```

### Streaming

```go
stream, err := zotero.NewStreamClient(ctx, "api-key")
if err != nil {
    log.Fatal(err)
}
defer stream.Close()

err = stream.Subscribe(ctx, []zotero.StreamSubscription{
    {APIKey: "api-key", Topics: []string{"/users/12345"}},
})

for evt := range stream.Events() {
    fmt.Printf("event: %s topic: %s\n", evt.Event, evt.Topic)
}
```

## Project structure

```
├── go.mod
├── go.sum
├── openapi.yaml              # OpenAPI 3.1 spec for the Zotero API v3
├── example/
│   └── basic/
│       └── main.go           # Runnable usage example
└── zotero/                   # Library package
    ├── zotero.go             # Package doc, constants
    ├── client.go             # Client, HTTP transport, pagination
    ├── errors.go             # APIError, WriteResponse, error helpers
    ├── types.go              # Domain types (Item, Collection, Search, …)
    ├── options.go            # LibraryID, RequestOption, ClientOption
    ├── items.go              # ItemsService
    ├── collections.go        # CollectionsService
    ├── searches.go           # SearchesService
    ├── tags.go               # TagsService
    ├── schema.go             # SchemaService
    ├── keys.go               # KeysService
    ├── groups.go             # GroupsService
    ├── fulltext.go           # FullTextService
    ├── deleted.go            # DeletedService
    ├── streaming.go          # StreamClient (WebSocket)
    └── *_test.go             # Tests for each service
```

## Contributing

Contributions are welcome! Here's how to get started.

### Prerequisites

- [Go 1.24+](https://go.dev/dl/)
- A [Zotero account](https://www.zotero.org/user/register) and [API key](https://www.zotero.org/settings/keys) (only needed for manual testing against the live API)

### Setup

```sh
git clone https://github.com/abtris/zotero-go-client.git
cd zotero-go-client
go mod download
```

### Running tests

All tests use `httptest.Server` — no network access or API key required:

```sh
go test ./...           # run all tests
go test -v ./zotero/    # verbose output
go test -race ./...     # with race detector
go test -cover ./...    # with coverage
```

### Code quality

```sh
go vet ./...            # static analysis
go build ./...          # verify compilation
```

### Adding a new service

1. Create `zotero/<service>.go` with a `<Service>Service` struct and methods following the existing pattern (accept `context.Context`, `LibraryID`, return `(*Type, *Response, error)`).
2. Register the service in `client.go` — add a field to `Client` and initialize it in `NewClient`.
3. Add types to `zotero/types.go` if needed.
4. Create `zotero/<service>_test.go` — use the `testServer` helper from `client_test.go` to set up a fake HTTP server.
5. Add endpoints to `openapi.yaml`.

### Writing tests

Tests use the shared `testServer` helper that creates an `httptest.Server` and a `Client` pointing to it:

```go
func TestMyEndpoint(t *testing.T) {
    client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
        // Assert request path, method, headers, body
        // Write response
        json.NewEncoder(w).Encode(expectedResponse)
    })

    ctx := context.Background()
    result, _, err := client.MyService.MyMethod(ctx, UserLibrary("123"))
    if err != nil {
        t.Fatal(err)
    }
    // Assert result
}
```

### Commit guidelines

- Keep commits focused — one logical change per commit.
- Run `go test ./...` and `go vet ./...` before pushing.
- Add tests for new functionality.

## API reference

See the [Zotero Web API v3 documentation](https://www.zotero.org/support/dev/web_api/v3/basics) for details on request parameters, response formats, and rate limiting.

The included [`openapi.yaml`](openapi.yaml) file provides a machine-readable OpenAPI 3.1 specification of the API.

## License

This project is available under the [MIT License](LICENSE).