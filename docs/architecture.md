# Architecture Guide

This document explains the internals of go-eos for contributors working on the
SDK itself. It is not a user guide.

## Overview

go-eos is an unofficial Go SDK for Epic Online Services (EOS). It wraps the
official EOS C SDK via Cgo bindings and exposes idiomatic Go interfaces.

The layers, from top to bottom:

```
Go application
  |
  v
eos/platform      -- Platform lifecycle, interface accessors (Auth, Lobby, P2P, ...)
  |
  v
eos/<interface>    -- Per-interface Go packages (auth/, connect/, lobby/, ...)
  |
  v
eos/internal/      -- Shared internals: threadworker, callback registry, cbinding
  |
  v
EOS C SDK          -- Linked via Cgo (static/ directory, gitignored)
```

Every EOS SDK call goes through `eos/internal/cbinding/`, which contains thin C
wrappers and Go exports. All business logic lives in Go.

## Thread Worker

The EOS SDK has strict thread-affinity requirements: every API call and
`EOS_Platform_Tick` must happen on the same OS thread. On macOS the requirement
is stronger -- the SDK's HTTP layer (Apple networking) needs the process's main
OS thread.

`eos/internal/threadworker/` solves this with a goroutine permanently locked to
an OS thread via `runtime.LockOSThread()`.

### Two startup modes

- `Start(ctx)` -- spawns a new goroutine, locks it to an arbitrary OS
  thread. Used when main-thread ownership is not required.
- `StartBlocking(ctx)` -- runs the worker loop on the *calling* goroutine.
  The caller must already own the main OS thread (via `runtime.LockOSThread()`
  in an `init()` function). `platform.RunOnMainThread` uses this path.

### Tick loop

The worker runs a `select` loop:

1. Drain the work channel (submitted via `Submit`/`SubmitWithContext`).
2. On each tick interval (default 16 ms), call `EOS_Platform_Tick`.
3. On context cancellation, drain remaining work items and return.

### Re-entrance detection

During `EOS_Platform_Tick`, the SDK may fire callbacks that ultimately call
`Submit` on the same worker. Enqueueing onto the work channel from within the
loop goroutine would deadlock. The worker detects this by comparing goroutine
IDs: if `Submit` is called from the worker's own goroutine, the function runs
inline immediately.

### Context decoupling

`RunOnMainThread` creates the worker's context from `context.Background()`, not
from the application context. This ensures the worker keeps ticking during
cleanup (e.g. `LeaveLobby` after Ctrl+C) so SDK callbacks can still land.

## Callback Infrastructure

`eos/internal/callback/` provides two mechanisms for receiving results from the
EOS SDK's asynchronous APIs.

### OneShotCallback

For request/response-style SDK calls (e.g. `EOS_Auth_Login`):

1. Go creates an `OneShotCallback`, which internally allocates a `cgo.Handle`
   pointing to itself.
2. The handle is passed as `uintptr_t` ClientData to the C wrapper.
3. The C wrapper calls the EOS SDK function, passing a static trampoline as the
   completion callback.
4. When the SDK completes the operation, it calls the trampoline on the tick
   thread.
5. The trampoline unpacks the C result struct into primitive types and calls the
   corresponding Go export (`//export goAuthLoginCallback`, etc.).
6. The Go export reconstructs a Go struct from the primitives, resolves the
   `cgo.Handle` back to the `OneShotCallback`, calls `Complete`, and deletes
   the handle.
7. The original caller receives the result via `Wait(ctx)` (channel receive).

### NotificationRegistry

For ongoing event subscriptions (e.g. login status changes, lobby member
updates):

- `Register(id, fn)` -- associates an EOS notification ID with a Go callback.
- `Dispatch(id, data)` -- called from the C trampoline's Go export to invoke
  the registered callback.
- `Unregister(id)` -- removes the callback when the notification is removed.
- Thread-safe via `sync.RWMutex`.

### The Go-to-C-to-Go flow

```
Go caller
  |  creates cgo.Handle, calls cbinding.EOS_Auth_Login(handle)
  v
cbinding Go func
  |  passes handle as uintptr_t to C wrapper
  v
auth_wrapper.c :: eos_auth_login()
  |  builds EOS_Auth_LoginOptions, calls EOS_Auth_Login with trampoline
  v
EOS SDK (async, completes during a future Tick)
  |
  v
auth_wrapper.c :: authLoginTrampoline()
  |  unpacks EOS_Auth_LoginCallbackInfo fields to primitives
  |  calls goAuthLoginCallback(resultCode, clientData, ...)
  v
auth_callback.go :: goAuthLoginCallback()
  |  reconstructs Go struct, resolves cgo.Handle, calls CompleteByHandle
  v
Original caller unblocks on Wait()
```

The key constraint: Go exports callable from C can only accept C-compatible
types. The trampolines convert opaque EOS pointers to `uintptr_t` and SDK
structs to flat primitive arguments.

## Memory Management

### EOS SDK allocations

Option structs (e.g. `EOS_Auth_LoginOptions`) are built by C wrapper functions
at function scope and passed to the SDK. The SDK reads them synchronously, so
stack-allocated structs are safe. Nested sub-structs (e.g. `EOS_Auth_Credentials`
inside login options) must also be at function scope -- placing them inside
if-blocks risks premature scope exit.

Some SDK calls return allocated objects (e.g. `EOS_Auth_CopyUserAuthToken`).
These must be released by calling the corresponding `*_Release` function
(e.g. `eos_auth_token_release`).

### cgo.Handle lifecycle

`cgo.Handle` bridges Go objects through C void pointers:

1. **Create**: `cgo.NewHandle(obj)` -- pins `obj` so the Go GC will not collect
   it.
2. **Pass**: the handle's `uintptr` value is passed as `ClientData` to C.
3. **Resolve**: in the Go export callback, `cgo.Handle(clientData).Value()`
   recovers the original Go object.
4. **Delete**: `handle.Delete()` -- unpins the object. For one-shot callbacks
   this happens in `CompleteByHandle`. The `OneShotCallback.Delete` method uses
   `atomic.Bool` to make it safe to call multiple times.

Forgetting to delete a handle leaks the pinned Go object.

### Go GC

All Go-side structs (callback info, platform config, interface wrappers) are
managed normally by the Go garbage collector. No special action is needed.

## Code Generation

The `eos/internal/cbinding/` package is partially generated from the EOS SDK's
C headers:

- `cmd/eos-bindgen` -- custom code generator that produces Go type
  definitions and C wrapper scaffolding from EOS header files.
- c-for-go -- used alongside eos-bindgen for additional binding generation.

The generated output is checked into the repository. After regenerating, files
may need manual adjustment (e.g. adding re-entrance-safe callback patterns or
fixing struct lifetime issues).

### File layout in cbinding/

Each EOS interface has a set of files:

| File | Purpose |
|------|---------|
| `auth.go` | Go functions that call into C wrappers |
| `auth_wrapper.c` | C trampolines + SDK call wrappers |
| `auth_wrapper.h` | C function declarations |
| `auth_callback.go` | Go exports (`//export`) called from C trampolines |
| `auth_types.go` | Go struct/enum definitions mirroring EOS C types |
| `auth_stub.go` | Pure Go stubs (used with `eosstub` build tag) |

## Build Tags

### `eosstub`

When building with `-tags eosstub`, all `cbinding` files guarded by
`//go:build !eosstub` are excluded. The `*_stub.go` files provide no-op or
fake implementations. This allows running tests without the EOS C SDK present.

### Default (no tag)

Real Cgo build. Requires the EOS SDK headers and libraries in the `static/`
directory. C files are additionally guarded by `#ifdef EOS_CGO` so they compile
only when the Cgo toolchain defines that macro.

### `CGO_ENABLED=0`

The `webapi/` package has no Cgo dependencies and works with
`CGO_ENABLED=0`. This is useful for server-side code that only needs the
EOS Web API (REST endpoints) without the full native SDK.
