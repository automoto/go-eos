# EOS Web API Client Guide

## Overview

The `webapi` package is a pure-Go HTTP client for the Epic Online Services (EOS) Web API. It has **no Cgo dependency** and builds with `CGO_ENABLED=0`, making it suitable for:

- **Game backend servers** — leaderboard queries, token verification
- **Matchmaking services** — account lookups
- **CLI tools** — scripting against EOS APIs without shipping native libraries

If your application needs real-time features like P2P networking, lobbies, or voice chat, use the native `eos/` package instead. The `webapi/` package covers the REST API surface only.

## Prerequisites

You need credentials from the [Epic Developer Portal](https://dev.epicgames.com/portal):

| Variable | Where to find it |
|----------|-----------------|
| `EOS_CLIENT_ID` | Product Settings > SDK Download & Credentials |
| `EOS_CLIENT_SECRET` | Same page |
| `EOS_DEPLOYMENT_ID` | Same page |

No EOS SDK binary, Developer Auth Tool, or native library is required.

## Installation

```bash
go get github.com/mydev/go-eos/webapi
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/mydev/go-eos/webapi"
)

func main() {
    client, err := webapi.New("your-deployment-id",
        webapi.WithClientCredentials("your-client-id", "your-client-secret"),
    )
    if err != nil {
        log.Fatal(err)
    }

    defs, err := client.GetLeaderboardDefinitions(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    for _, d := range defs {
        fmt.Printf("%s (stat: %s, aggregation: %s)\n",
            d.Spec.Name, d.Spec.RankBy.Stat, d.Spec.RankBy.Aggregation)
    }
}
```

Build and run with no native dependencies:

```bash
CGO_ENABLED=0 go run main.go
```

## Authentication

### Client Credentials (backend services)

The standard flow for server-to-server authentication. The client automatically acquires and refreshes OAuth2 tokens.

```go
client, err := webapi.New(deploymentID,
    webapi.WithClientCredentials(clientID, clientSecret),
)
```

Tokens are cached and refreshed before expiry — you never manage tokens manually.

### Exchange Code (Epic Games Launcher)

For game clients launched via the Epic Games Launcher, which provides a one-time exchange code:

```go
client, err := webapi.New(deploymentID,
    webapi.WithExchangeCode(clientID, clientSecret, exchangeCode),
)
```

The exchange code is single-use. After the initial token exchange, the client caches and refreshes using the refresh token.

## Leaderboards

### Get Definitions

Returns all leaderboard definitions for the deployment. Each definition includes its spec (name, ranking criteria, time window) and any inline player entries.

```go
defs, err := client.GetLeaderboardDefinitions(ctx)
for _, d := range defs {
    fmt.Printf("%s (stat: %s, start: %s, end: %s)\n",
        d.Spec.Name, d.Spec.RankBy.Stat, d.Spec.Start, d.Spec.End)
}
```

The `LeaderboardDefinition` struct:

```go
type LeaderboardDefinition struct {
    Spec    LeaderboardSpec    `json:"spec"`
    Players []LeaderboardEntry `json:"players"`
}

type LeaderboardSpec struct {
    Name   string            `json:"name"`
    RankBy LeaderboardRankBy `json:"rankBy"`
    Start  string            `json:"start"`
    End    string            `json:"end"`
}

type LeaderboardRankBy struct {
    Stat        string `json:"stat"`
    Aggregation string `json:"aggregation"`
}
```

### Get Rankings

```go
entries, err := client.GetLeaderboardRankings(ctx, "top-scores",
    webapi.WithOffset(0),
    webapi.WithLimit(25),
)
for _, e := range entries {
    fmt.Printf("#%d %s — %d\n", e.Rank, e.ProductUserID, e.Score)
}
```

## Auth Verification

### Verify a Token

Verify a third-party access token (e.g., one sent by a game client):

```go
info, err := client.VerifyToken(ctx, playerToken)
if err != nil {
    // token is invalid or expired
}
fmt.Printf("account=%s active=%v expires_in=%d\n",
    info.AccountID, info.Active, info.ExpiresIn)
```

### Get Account Info

Look up display names for up to 100 account IDs:

```go
accounts, err := client.GetAccounts(ctx, []string{"id1", "id2"})
for _, a := range accounts {
    fmt.Printf("%s: %s\n", a.AccountID, a.DisplayName)
}
```

## Error Handling

All methods return `*webapi.APIError` for HTTP error responses. Use `errors.Is` for broad status checks and `errors.As` for detailed inspection:

```go
_, err := client.GetLeaderboardRankings(ctx, "missing-board")
if errors.Is(err, webapi.ErrNotFound) {
    // 404 — leaderboard doesn't exist
}
if errors.Is(err, webapi.ErrRateLimited) {
    // 429 — slow down (the client already retried automatically)
}

var apiErr *webapi.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("HTTP %d: %s: %s\n", apiErr.HTTPStatus, apiErr.ErrorCode, apiErr.Message)
}
```

Available sentinels: `ErrUnauthorized` (401), `ErrForbidden` (403), `ErrNotFound` (404), `ErrRateLimited` (429), `ErrServerError` (500).

## Rate Limiting and Retries

The client includes built-in protection against rate limiting and transient failures.

**Proactive rate limiting** — a token-bucket limiter prevents bursting beyond the configured rate. Default: 10 requests/second, burst 20.

```go
client, _ := webapi.New(deploymentID,
    webapi.WithClientCredentials(clientID, clientSecret),
    webapi.WithRateLimit(50, 100), // 50 rps, burst 100
)
```

**Reactive 429 handling** — if the server returns HTTP 429, the client respects the `Retry-After` header and retries automatically.

**Exponential backoff** — transient errors (500, 502, 503, 504) and network failures trigger exponential backoff with jitter:

```go
client, _ := webapi.New(deploymentID,
    webapi.WithClientCredentials(clientID, clientSecret),
    webapi.WithRetryPolicy(webapi.RetryPolicy{
        MaxRetries:  5,
        BaseDelay:   1 * time.Second,
        MaxDelay:    60 * time.Second,
        JitterRatio: 0.5,
    }),
)
```

Default policy: 3 retries, 500ms base delay, 30s max, 0.5 jitter ratio.

## Custom HTTP Client

Inject a custom `*http.Client` for proxies, custom TLS, or timeouts:

```go
client, _ := webapi.New(deploymentID,
    webapi.WithClientCredentials(clientID, clientSecret),
    webapi.WithHTTPClient(&http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            Proxy: http.ProxyFromEnvironment,
        },
    }),
)
```

## Logging

The client logs requests via `log/slog`. Pass a custom logger:

```go
logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))
client, _ := webapi.New(deploymentID,
    webapi.WithClientCredentials(clientID, clientSecret),
    webapi.WithLogger(logger),
)
```

Log levels:
- `Debug` — request method, URL, status code, duration, attempt number
- Request/response bodies are not logged (may contain sensitive data)
- Authorization header values are never logged
