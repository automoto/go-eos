# go-eos

Unofficial Go SDK for Epic Online Services.

> **Status: Work in Progress** — This SDK is under active development and not yet ready for production use.

## Overview

`go-eos` provides idiomatic Go bindings for [Epic Online Services](https://dev.epicgames.com/docs/epic-online-services), enabling Go game developers to integrate authentication, matchmaking, lobbies, sessions, and peer-to-peer networking into their games.

The SDK has two components:

- **Native Bindings** (`eos/`) — Cgo bindings to the EOS C SDK for full feature access including P2P networking
- **Web API Client** (`webapi/`) — Pure-Go HTTP client for EOS REST APIs (no native dependencies)

## Prerequisites

- Go 1.26+
- C compiler (for native bindings): GCC/Clang on Linux/macOS, MinGW on Windows
- [EOS C SDK](https://dev.epicgames.com/docs/epic-online-services/eos-get-started/get-started-resources) binaries (not included; download from Epic Developer Portal)

## Installation

```
go get github.com/mydev/go-eos
```

## Quick Start

*Coming soon — see `examples/` for usage once available.*

## Supported Platforms

- Windows x64
- Linux x64
- macOS (amd64, arm64)

## License

[MIT](LICENSE)
