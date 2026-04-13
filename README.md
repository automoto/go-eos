# go-eos

[![Go Reference](https://pkg.go.dev/badge/github.com/mydev/go-eos.svg)](https://pkg.go.dev/github.com/mydev/go-eos)
[![CI](https://github.com/mydev/go-eos/actions/workflows/ci.yml/badge.svg)](https://github.com/mydev/go-eos/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mydev/go-eos)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Unofficial Go SDK for [Epic Online Services](https://dev.epicgames.com/docs/epic-online-services).

## Features

- Auth: Epic Account login (Developer Auth, Exchange Code, Persistent Auth, Account Portal)
- Connect: Product user identity (Epic ID Token, Device ID, Steam, Discord, and more)
- Lobbies: Create, search, join, and manage lobbies with attribute-based matchmaking
- Sessions: Game session creation, search, and member management
- P2P: Peer-to-peer packet send/receive with NAT traversal
- Web API: Pure-Go HTTP client for leaderboards, fetching stats, auth verification (`CGO_ENABLED=0`)

## Project Structure

| Path | Description |
|------|-------------|
| `eos/` | Native Cgo bindings to the EOS C SDK — full feature access including P2P |
| `webapi/` | Pure-Go HTTP client for the EOS REST API — no native dependencies |
| `examples/` | Working examples: auth, lobbies, device-id, p2p, web API queries |

## Prerequisites

- Go 1.26+
- C compiler (for native bindings): GCC/Clang on Linux/macOS, MinGW on Windows
- [EOS C SDK](https://dev.epicgames.com/docs/epic-online-services/eos-get-started/get-started-resources) v1.19.0.3 (download from Epic Developer Portal)

The `webapi/` package has no native dependencies and builds with `CGO_ENABLED=0`.

## Installation

```bash
go get github.com/mydev/go-eos
```

## Quick Start

Initialize the platform, log in, and create a lobby:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "runtime"

    "github.com/mydev/go-eos/eos/auth"
    "github.com/mydev/go-eos/eos/connect"
    "github.com/mydev/go-eos/eos/lobby"
    "github.com/mydev/go-eos/eos/platform"
    "github.com/mydev/go-eos/eos/types"
)

func init() { runtime.LockOSThread() }

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    cfg := platform.PlatformConfig{
        ProductName:  "my-game",
        ProductVersion: "1.0.0",
        ProductId:    os.Getenv("EOS_PRODUCT_ID"),
        SandboxId:    os.Getenv("EOS_SANDBOX_ID"),
        DeploymentId: os.Getenv("EOS_DEPLOYMENT_ID"),
        ClientId:     os.Getenv("EOS_CLIENT_ID"),
        ClientSecret: os.Getenv("EOS_CLIENT_SECRET"),
    }

    if err := platform.RunOnMainThread(ctx, cfg, func(p *platform.Platform) error {
        ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
        defer stop()

        // Auth login via Developer Auth Tool
        result, err := p.Auth().Login(ctx, auth.LoginOptions{
            CredentialType: types.LoginCredentialDeveloper,
            ID:             os.Getenv("EOS_DEV_AUTH_HOST"),
            Token:          os.Getenv("EOS_DEV_AUTH_CREDENTIAL"),
        })
        if err != nil {
            return fmt.Errorf("auth login: %w", err)
        }

        // Connect login with the ID token
        idToken, _ := p.Auth().CopyIdToken(result.LocalUserId)
        cr, err := p.Connect().Login(ctx, connect.LoginOptions{
            CredentialType: types.ExternalCredentialEpicIDToken,
            Token:          idToken,
        })
        if err != nil {
            return fmt.Errorf("connect login: %w", err)
        }

        // Create a lobby
        lobbyId, err := p.Lobby().CreateLobby(ctx, lobby.CreateLobbyOptions{
            LocalUserId: cr.LocalUserId,
            MaxMembers:  4,
            BucketId:    "game:mode=default",
        })
        if err != nil {
            return fmt.Errorf("create lobby: %w", err)
        }
        fmt.Printf("Lobby created: %s\n", lobbyId)

        <-ctx.Done()
        return nil
    }); err != nil {
        log.Fatal(err)
    }
}
```

See the [Getting Started Guide](docs/getting-started.md) for a full walkthrough, or browse the [examples](examples/).

## Documentation

- [Getting Started](docs/getting-started.md) end-to-end setup and first lobby
- [Architecture](docs/architecture.md) how the internals work (thread worker, callbacks, Cgo)
- [Web API Guide](docs/web-api-guide.md) pure-Go server-side usage
- [Ebitengine Guide](docs/ebitengine-guide.md) integrating go-eos with Ebitengine
- [Distribution](docs/distribution.md) building and shipping per-platform binaries
- [Migration from C++](docs/migration-from-cpp.md) mapping EOS C/C++ API to go-eos
- [EOS SDK Versions](docs/eos-sdk-versions.md) SDK version management

## Supported Platforms

- Windows x64
- Linux x64
- macOS (amd64, arm64)

## License

[MIT](LICENSE)
