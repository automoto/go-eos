# Ebitengine Integration Guide

go-eos has no dependency on Ebitengine. This guide explains how to structure a
game that uses both libraries together.

## The Challenge

On macOS, both Ebitengine and the EOS SDK want the main OS thread:

- Ebitengine drives its window event loop on the main thread.
- EOS SDK uses Apple networking, which dispatches HTTP completions through
  the main thread's run loop.

go-eos solves this with two initialization paths:

| Function | Main thread usage | Best for |
|---|---|---|
| `platform.RunOnMainThread` | Blocks the main thread for EOS tick loop | CLI tools, headless servers |
| `platform.Initialize` | Starts EOS tick on a background goroutine | GUI apps like Ebitengine |

For Ebitengine games, use `platform.Initialize`. It spins up a background
worker goroutine (with `runtime.LockOSThread`) that handles
`EOS_Platform_Tick` at ~16ms intervals, leaving the main thread free for
Ebitengine.

## Integration Pattern

```
main thread (Ebitengine)         background goroutine (go-eos worker)
  |                                |
  |  ebiten.RunGame(game)         |  runtime.LockOSThread()
  |    -> Game.Update()           |  EOS_Platform_Tick every ~16ms
  |    -> Game.Draw()             |  dispatches SDK callbacks
  |    -> Game.Layout()           |
  |                                |
  |  game.platform.Shutdown()     |  stops
```

1. Call `platform.Initialize(cfg)` -- starts the background tick worker.
2. Call `ebiten.RunGame(game)` on the main thread.
3. In `Game.Update()`, interact with EOS (lobby, auth, etc.).
4. When the game exits, call `platform.Shutdown()`.

## Minimal Code Outline

```go
package main

import (
    "log"
    "os"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/mydev/go-eos/eos/platform"
)

type Game struct {
    plat *platform.Platform
}

func (g *Game) Update() error {
    // Interact with EOS here, e.g.:
    //   g.plat.Lobby().CreateLobby(ctx, opts)
    //   g.plat.Connect().Login(ctx, opts)
    return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Render your game
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 640, 480
}

func main() {
    cfg := platform.PlatformConfig{
        ProductName:    "my-game",
        ProductVersion: "1.0.0",
        ProductId:      os.Getenv("EOS_PRODUCT_ID"),
        SandboxId:      os.Getenv("EOS_SANDBOX_ID"),
        DeploymentId:   os.Getenv("EOS_DEPLOYMENT_ID"),
        ClientId:       os.Getenv("EOS_CLIENT_ID"),
        ClientSecret:   os.Getenv("EOS_CLIENT_SECRET"),
    }

    plat, err := platform.Initialize(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer plat.Shutdown()

    if err := ebiten.RunGame(&Game{plat: plat}); err != nil {
        log.Fatal(err)
    }
}
```

Note: do **not** call `runtime.LockOSThread()` in `init()` when using
Ebitengine -- Ebitengine manages the main thread itself.

## Thread Safety

All go-eos SDK calls go through the internal thread worker via `Submit()`.
This means calls like `Lobby().CreateLobby()` and `Connect().Login()` are
safe to call from `Game.Update()` -- they serialize onto the worker's locked
OS thread automatically.

Async operations return results via channels or callbacks. Store results in
your `Game` struct and read them during `Update()`.

```go
func (g *Game) Update() error {
    select {
    case result := <-g.lobbyResultCh:
        // handle lobby creation result
    default:
    }
    return nil
}
```

Avoid calling EOS methods from `Draw()` -- keep SDK interaction in `Update()`.

## Complete Example

See `examples/ebitengine-lobby` for a working example with lobby creation,
member list rendering, and chat.

## Note on Dependencies

go-eos has no Ebitengine dependency. The `examples/ebitengine-lobby` directory
has its own `go.mod` to keep Ebitengine out of the root module's dependency
tree.
