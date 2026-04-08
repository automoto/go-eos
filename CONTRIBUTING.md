# Contributing to go-eos

## Prerequisites

- Go 1.26.2+
- C compiler (Xcode Clang on macOS, GCC/Clang on Linux)
- EOS C SDK v1.19.0.3 placed in `static/` (gitignored, download from [Epic Developer Portal](https://dev.epicgames.com/portal))
- `golangci-lint` v2+ (`brew install golangci-lint` or see [install docs](https://golangci-lint.run/docs/install/))
- `clang-format` (`brew install clang-format` on macOS, `apt-get install clang-format` on Linux)

## Build and Test

```bash
# Build with real Cgo (requires EOS SDK in static/)
go build ./...

# Run all unit tests (uses pure Go stubs, no SDK needed)
go test -tags eosstub -race ./...

# Run non-platform unit tests without stubs
go test -race ./eos/types/... ./eos/internal/threadworker/... ./eos/internal/callback/...

# Lint Go code
make lint

# Lint C code (checks formatting)
make lint-c

# Auto-format C code
make format-c

# Check for known vulnerabilities
make vulncheck
```

## Build Tags

| Tag | Purpose |
|-----|---------|
| `eosstub` | Uses pure Go stubs instead of real Cgo. Enables unit testing without the EOS SDK. |
| (default) | Builds with real Cgo, links against the EOS SDK shared library. |

Unit tests should always be run with `-tags eosstub`. The real Cgo build is verified by `go build ./...`.

## Project Structure

```
eos/
  types/              # Shared Go types (Result, IDs, enums)
  platform/           # Platform init, tick loop, shutdown
  auth/               # Auth interface wrapper
  connect/            # Connect interface wrapper
  internal/
    cbinding/          # Cgo bindings to EOS C SDK
      *.go             # Go wrapper functions + stubs
      *_wrapper.c/.h   # C implementations (trampolines, SDK wrappers)
    callback/          # Oneshot + notification callback infrastructure
    threadworker/      # OS-thread-locked worker goroutine
```

## C Code Architecture

All C code lives in `eos/internal/cbinding/` as separate `.c` and `.h` files. The Go preambles only contain `#include` directives.

**Design principle: thin C wrappers, comprehensive Go tests via stubs.**

The C layer is intentionally kept as thin as possible -- it handles only what must be done in C:
- **Trampolines**: receive SDK callbacks and forward to Go exports with primitive types (`uintptr_t`, `int`, `const char*`) to avoid `go vet` warnings
- **SDK wrappers**: build options structs, set `ApiVersion`, call SDK functions
- **Type casting**: convert between opaque C pointers and `uintptr_t` at the C boundary

All business logic, type dispatch, and error handling lives in Go where it is tested using the `eosstub` stubs. This avoids the need for a C testing framework while keeping the C code trivially correct by construction.

The `.c` files use `#ifdef EOS_CGO` guards and `//go:build !eosstub` constraints so they compile to nothing when building with stubs.

## Testing Conventions

- **TDD**: write failing tests before implementation
- **AAA pattern**: Arrange-Act-Assert
- **testify assert**: use `github.com/stretchr/testify/assert`
- **Table-based tests**: use where appropriate to keep tests concise
- **Test names**: describe behavior, e.g. `Test_login_should_succeed_with_developer_credentials`

## Code Quality

- Run `make lint` before submitting (must pass with zero issues)
- Run `make lint-c` to check C formatting
- Use early returns to reduce nesting
- Proper error handling; avoid panics
- Keep code simple; avoid unnecessary abstractions
