# Getting Started with go-eos

This guide walks you through setting up go-eos from scratch and creating a working lobby. Estimated time: 30 minutes.

go-eos is an unofficial Go SDK for [Epic Online Services](https://dev.epicgames.com/docs/epic-online-services). It provides Cgo bindings to the EOS C SDK for full feature access (auth, lobbies, P2P) and a pure-Go web API client for server-side use.

## 1. Prerequisites

- Go 1.26+ ([download](https://go.dev/dl/))
- C compiler -- GCC or Clang on Linux/macOS, MinGW on Windows
- Epic Developer Portal account -- sign up at https://dev.epicgames.com/portal

Verify your Go installation:

```bash
go version
# go version go1.26.2 ...
```

## 2. Download and Install the EOS SDK

1. Go to the [Epic Developer Portal](https://dev.epicgames.com/portal) and navigate to your product's SDK Download page.
2. Download **EOS C SDK v1.19.0.3**.
3. Extract it into the `static/` directory at the project root:

```bash
# From the go-eos project root
mkdir -p static
# Extract the SDK archive so the directory structure looks like:
#   static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/
#     Include/
#     Bin/
#     Tools/
```

The `static/` directory is gitignored -- each developer downloads their own copy. The Cgo build flags in `eos/internal/cbinding/cgo.go` reference this exact path.

## 3. Portal Setup

You need four things from the Epic Developer Portal:

### Create a Product

1. Log in to https://dev.epicgames.com/portal
2. Create a new product (or use an existing one)
3. Under your product, note the **Product ID**

### Sandbox and Deployment

1. Navigate to Product Settings > Sandboxes
2. A default sandbox exists; note its **Sandbox ID**
3. Create a deployment within the sandbox; note the **Deployment ID**

### Client Credentials

1. Go to Product Settings > SDK Download & Credentials
2. Create a new client policy (or use the default GameClient policy)
3. Note the **Client ID** and **Client Secret**

### Developer Auth Tool

The DevAuth Tool lets you authenticate locally without a browser-based OAuth flow. It ships with the SDK.

**Extract and launch:**

```bash
# macOS
cd static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Tools
unzip EOS_DevAuthTool-darwin-x64-1.2.1.zip -d DevAuthTool
./DevAuthTool/EOS_DevAuthTool.app/Contents/MacOS/EOS_DevAuthTool
```

If macOS blocks it ("unidentified developer"):
```bash
xattr -cr ./DevAuthTool/EOS_DevAuthTool.app
```

**Configure the tool:**

1. Set the **Port** to `6547` (default) and click **Login**
2. Log in with your Epic Games account in the browser window
3. Enter `player1` as the **Credential Name** and click **Create**
4. For lobby testing, repeat with a second Epic account to create a `player2` credential

Keep the DevAuth Tool running during development.

## 4. Environment Variables

Create a `.env` file (gitignored) or export these directly:

```bash
export EOS_PRODUCT_ID="your-product-id"          # From Product Settings
export EOS_SANDBOX_ID="your-sandbox-id"           # From Sandboxes page
export EOS_DEPLOYMENT_ID="your-deployment-id"     # From Deployments page
export EOS_CLIENT_ID="your-client-id"             # From SDK Credentials
export EOS_CLIENT_SECRET="your-client-secret"     # From SDK Credentials
export EOS_DEV_AUTH_HOST="localhost:6547"          # DevAuth Tool address
export EOS_DEV_AUTH_CREDENTIAL="player1"           # Credential name you created
```

Source it before running examples:

```bash
source .env
```

## 5. Initialize the Platform

Every go-eos program starts by initializing the EOS platform. On macOS, the SDK's HTTP layer requires the main OS thread, so you must lock to it and use `platform.RunOnMainThread`.

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "runtime"

    "github.com/mydev/go-eos/eos/platform"
)

func init() {
    // Required on macOS -- EOS SDK uses Apple networking on the main thread.
    runtime.LockOSThread()
}

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    cfg := platform.PlatformConfig{
        ProductName:    "my-game",
        ProductVersion: "1.0.0",
        ProductId:      os.Getenv("EOS_PRODUCT_ID"),
        SandboxId:      os.Getenv("EOS_SANDBOX_ID"),
        DeploymentId:   os.Getenv("EOS_DEPLOYMENT_ID"),
        ClientId:       os.Getenv("EOS_CLIENT_ID"),
        ClientSecret:   os.Getenv("EOS_CLIENT_SECRET"),
    }

    if err := platform.RunOnMainThread(ctx, cfg, run); err != nil {
        log.Fatal(err)
    }
}

func run(p *platform.Platform) error {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    // Your game logic goes here.

    <-ctx.Done()
    return nil
}
```

`RunOnMainThread` initializes the SDK, starts a tick loop on the main OS thread, and calls your `run` function on a separate goroutine. When `run` returns, the SDK shuts down cleanly.

## 6. Login with DevAuth

Add auth and connect imports and call `Auth().Login()` inside your `run` function:

```go
import (
    "fmt"

    "github.com/mydev/go-eos/eos/auth"
    "github.com/mydev/go-eos/eos/types"
)

// Inside run():
devHost := os.Getenv("EOS_DEV_AUTH_HOST")
devCred := os.Getenv("EOS_DEV_AUTH_CREDENTIAL")

fmt.Println("Logging in via Developer Auth Tool...")
loginResult, err := p.Auth().Login(ctx, auth.LoginOptions{
    CredentialType: types.LoginCredentialDeveloper,
    ID:             devHost,
    Token:          devCred,
})
if err != nil {
    return fmt.Errorf("auth login: %w", err)
}
fmt.Printf("Auth login successful! EpicAccountId: %s\n", loginResult.LocalUserId)
```

This authenticates your Epic Games account through the running DevAuth Tool. The `loginResult.LocalUserId` is your Epic Account ID, used in the next step.

## 7. Connect Login

EOS has two identity layers: **Auth** (Epic account) and **Connect** (product user). Multiplayer features like lobbies and P2P use the Connect identity. Bridge from Auth to Connect using an ID token:

```go
import (
    "github.com/mydev/go-eos/eos/connect"
)

// Inside run(), after Auth login:
idToken, err := p.Auth().CopyIdToken(loginResult.LocalUserId)
if err != nil {
    return fmt.Errorf("copy id token: %w", err)
}

connectResult, err := p.Connect().Login(ctx, connect.LoginOptions{
    CredentialType: types.ExternalCredentialEpicIDToken,
    Token:          idToken,
})
if err != nil {
    // First time: the user has no product user yet. Create one.
    if connectResult != nil && connectResult.ContinuanceToken != 0 {
        fmt.Println("Creating product user...")
        userId, createErr := p.Connect().CreateUser(ctx, connectResult.ContinuanceToken)
        if createErr != nil {
            return fmt.Errorf("create user: %w", createErr)
        }
        fmt.Printf("ProductUserId: %s\n", *userId)
        // Use *userId as your productUserId going forward.
    } else {
        return fmt.Errorf("connect login: %w", err)
    }
} else {
    fmt.Printf("ProductUserId: %s\n", connectResult.LocalUserId)
    // Use connectResult.LocalUserId as your productUserId going forward.
}
```

The `CreateUser` path runs once per Epic account per product. Subsequent logins return the existing product user directly.

## 8. Create and Join a Lobby

With a `productUserId` in hand, you can create and search for lobbies.

**Create a lobby:**

```go
import (
    "github.com/mydev/go-eos/eos/lobby"
)

lobbyId, err := p.Lobby().CreateLobby(ctx, lobby.CreateLobbyOptions{
    LocalUserId:     productUserId,
    MaxMembers:      4,
    BucketId:        "game:mode=default",
    PermissionLevel: lobby.PermissionPublicAdvertised,
    AllowInvites:    true,
})
if err != nil {
    return fmt.Errorf("create lobby: %w", err)
}
fmt.Printf("Lobby created: %s\n", lobbyId)
```

**Search and join from another process:**

```go
search, err := p.Lobby().CreateLobbySearch(10)
if err != nil {
    return fmt.Errorf("create search: %w", err)
}
defer search.Release()

if err := search.SetParameter("bucket", "game:mode=default", lobby.ComparisonEqual); err != nil {
    return fmt.Errorf("set search parameter: %w", err)
}

results, err := search.Find(ctx, productUserId)
if err != nil {
    return fmt.Errorf("find lobbies: %w", err)
}
fmt.Printf("Found %d lobbies\n", len(results))

if len(results) > 0 {
    if err := p.Lobby().JoinLobby(ctx, productUserId, results[0]); err != nil {
        return fmt.Errorf("join lobby: %w", err)
    }
    fmt.Println("Joined lobby!")
    // Release search results when done.
    for _, d := range results {
        d.Release()
    }
}
```

**Try it end-to-end:** open two terminals, source your `.env` in each (use `player1` in one and `player2` in the other), and run the lobby-chat example:

```bash
# Terminal 1 -- create
go run ./examples/lobby-chat --create --bucket-id "game:mode=deathmatch"

# Terminal 2 -- search and join
go run ./examples/lobby-chat --join --bucket-id "game:mode=deathmatch"
```

## 9. Next Steps

- **P2P networking** -- see `examples/p2p-echo/` for peer-to-peer packet exchange with NAT traversal
- **Game engine integration** -- see [Ebitengine Guide](ebitengine-guide.md) for using go-eos with the Ebitengine game engine
- **Server-side web API** -- see [Web API Guide](web-api-guide.md) for the pure-Go HTTP client (`CGO_ENABLED=0`, no SDK binary needed)
- **Full examples** -- browse `examples/` for auth-login, connect-deviceid, lobby-chat, p2p-echo, and webapi-query
