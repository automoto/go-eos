# Distributing a Go Game That Uses go-eos

This guide covers what a Go game developer needs to build and ship a game that uses `go-eos`. It is opinionated for the common cases and links to authoritative sources for the rest.

If you only need server-side EOS access (leaderboards, auth verification), use the `webapi/` package instead — it is pure Go, builds with `CGO_ENABLED=0`, and has none of the issues described here.

---

## TL;DR — what ships with your game

A go-eos game is **not** a single self-contained binary. You ship at least two files per platform:

| Platform | Game binary | EOS runtime library | Where the library goes |
|---|---|---|---|
| Windows x64 | `mygame.exe` | `EOSSDK-Win64-Shipping.dll` | Same directory as the `.exe` |
| Linux x64 | `mygame` | `libEOSSDK-Linux-Shipping.so` | Same directory, with `RPATH=$ORIGIN` baked in (or set `LD_LIBRARY_PATH`) |
| macOS (universal or per-arch) | `mygame.app/Contents/MacOS/mygame` | `libEOSSDK-Mac-Shipping.dylib` | `mygame.app/Contents/Frameworks/`, with `@rpath/@executable_path` set |

Plus your EOS credentials (ProductId, SandboxId, DeploymentId, ClientId, ClientSecret) — typically baked into the binary at build time via `-ldflags -X` or shipped as a config file.

---

## 1. Build prerequisites

### Per developer machine

| Tool | Why | Notes |
|---|---|---|
| Go 1.26+ | Required by `go-eos` | See `go.mod` for the floor |
| Native C compiler | Cgo needs it | macOS: Xcode CLT (`xcode-select --install`). Linux: `gcc` or `clang`. Windows: MinGW-w64 (`gcc.exe` on PATH) — MSVC is not supported by Cgo |
| EOS C SDK 1.19.x | Headers + per-platform shared library | Download from [Epic Developer Portal](https://dev.epicgames.com/portal) — requires accepting Epic's license. Cannot be redistributed |
| Epic Developer Portal account | To get credentials | Free signup. You will get: ProductId, SandboxId, DeploymentId, ClientId, ClientSecret |

### One-time portal setup

The first time you set up an EOS product, you must:

1. Create a Product in the portal
2. Create a Sandbox under that product (the portal creates a default "Live" sandbox)
3. Create a Deployment within the sandbox
4. Create a Client with the permissions your game needs (at minimum: Auth, Connect, Lobby — Sessions and P2P if you use them)
5. Create an Application and link it to the Client (this is what enables `Connect.Login` with `EpicIDToken`)
6. (Optional, anonymous-only games) Enable Device ID under the Application's Identity Providers — without this, `Connect.CreateDeviceId` will fail

The portal UI changes regularly; consult [Epic's getting-started docs](https://dev.epicgames.com/docs/) for the current click path.

---

## 2. Where the runtime library lives

The EOS shared library has to be findable by the OS dynamic loader at startup. Each platform has its own search rules.

### Windows

The Win32 loader searches the directory containing the `.exe` first. **The simplest distribution layout is to drop `EOSSDK-Win64-Shipping.dll` next to your `mygame.exe`.** Players also need the **Microsoft Visual C++ Redistributable** (2015 or later). Bundle it in your installer or instruct players to install it.

```
MyGame/
  mygame.exe
  EOSSDK-Win64-Shipping.dll
  vcredist_x64.exe         (optional, in installer)
```

Windows builds require MinGW-w64 (`gcc.exe` on PATH). The `#cgo windows LDFLAGS` directive in `eos/internal/cbinding/cgo.go` links against `EOSSDK-Win64-Shipping.dll`'s import library. No rpath configuration is needed — Windows finds DLLs next to the `.exe` by default.

### Linux

The Linux loader searches `RPATH`, then `LD_LIBRARY_PATH`, then standard system paths. You almost always want to bake `$ORIGIN` (the directory containing the binary) into the RPATH so the game finds the SDK without env vars:

```bash
go build -ldflags="-r '\$ORIGIN'" -o mygame ./cmd/mygame
```

The `-r` ldflag tells the Go linker to set DT_RUNPATH on the resulting ELF binary. Layout:

```
mygame/
  mygame
  libEOSSDK-Linux-Shipping.so
```

Verify with `readelf -d mygame | grep RUNPATH` after the build.

### macOS

macOS uses install names on each library plus `@rpath` resolution from the loading binary. Two complications:

1. `@loader_path` and similar tokens are not allowed in Cgo `LDFLAGS` (a Cgo flag-safety restriction). The current go-eos build sets the rpath to the absolute `static/.../Bin` directory, which works for development but not for distribution.
2. macOS apps need to be signed and notarized for Gatekeeper to allow unsigned downloads to launch on a fresh machine.

For distribution, the canonical layout is an `.app` bundle:

```
MyGame.app/
  Contents/
    Info.plist
    MacOS/
      mygame                              (your binary)
    Frameworks/
      libEOSSDK-Mac-Shipping.dylib
```

After building, you need to:

1. Strip the development rpath and add a runtime one. `install_name_tool -delete_rpath <devpath> mygame` then `install_name_tool -add_rpath @executable_path/../Frameworks mygame`.
2. Verify the SDK library's install name with `otool -D libEOSSDK-Mac-Shipping.dylib`. If it is not `@rpath/libEOSSDK-Mac-Shipping.dylib`, fix it: `install_name_tool -id @rpath/libEOSSDK-Mac-Shipping.dylib libEOSSDK-Mac-Shipping.dylib`.
3. Sign both the dylib and the binary with your Developer ID, with hardened runtime enabled: `codesign --force --options runtime --sign "Developer ID Application: ..." Contents/Frameworks/libEOSSDK-Mac-Shipping.dylib` and the same for `Contents/MacOS/mygame`.
4. Notarize the bundled `.app` (or a `.dmg`/`.zip` containing it) via `xcrun notarytool submit ... --wait` and staple the ticket.

There is no shortcut for any of this. Apple's own [code signing and notarization docs](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution) are the source of truth.

**Packaging script:** `scripts/package-macos.sh` automates the bundle creation, rpath fixup, and optional code signing:

```bash
go build -o mygame ./cmd/mygame
./scripts/package-macos.sh mygame /path/to/libEOSSDK-Mac-Shipping.dylib
./scripts/package-macos.sh mygame /path/to/libEOSSDK-Mac-Shipping.dylib --sign "Developer ID Application: My Company"
```

### Development vs distribution rpath

During development, `cgo.go` sets the rpath to the `static/.../SDK/Bin` directory via `${SRCDIR}`. This works for `go run` and `go build` during development but is an absolute path that won't exist on other machines.

For distribution, each platform needs a different strategy:

| Platform | Development rpath | Distribution strategy |
|----------|------------------|----------------------|
| Windows | N/A (no rpath) | DLL next to `.exe` — works by default |
| Linux | `${SRCDIR}/.../SDK/Bin` | Build with `-ldflags="-r '\$ORIGIN'"` to set RUNPATH to the binary's directory |
| macOS | `${SRCDIR}/.../SDK/Bin` | Use `scripts/package-macos.sh` or manually run `install_name_tool` to set `@executable_path/../Frameworks` |

---

## 3. Cross-compilation: don't

Cgo + cross-compilation is famously painful. To produce a Windows build from a macOS host, you need a Windows-targeting C compiler (MinGW-w64), the Windows EOS SDK, and matching Win32 link flags. The Go ecosystem's usual `GOOS=windows go build` does not work because the C side can't follow.

**Realistic options, in order of recommendation:**

1. Build natively on each target. Three machines or three CI runners (one per OS). This is what Epic's own samples and most shipping Cgo projects do.
2. Use Docker for the Linux build from any host. A `golang:1.26` container with `gcc` and the Linux EOS SDK mounted is straightforward.
3. `zig cc` as the Cgo compiler. Zig can target Windows and Linux from any host and is becoming a popular Cgo cross-compilation backend, but it is not currently exercised by go-eos. If you go this route, expect to debug linker issues.

Per the docs: go-eos's own CI runs only on Linux to keep runner costs down. macOS and Windows are validated manually.

---

## 4. Storefront-specific notes

EOS is store-agnostic, but each storefront has a preferred login path that gives you a much better player experience than DevAuth or Account Portal.

### Epic Games Store

The Epic launcher invokes your game with command-line arguments including an **exchange code**. Skip DevAuth entirely: parse the exchange code from `os.Args` and call `Auth.Login` with `LoginCredentialExchangeCode`. Then continue to `Connect.Login` with the resulting Epic ID token. The player never sees a login prompt.

### Steam

Use the Steam → Connect path. From your game, obtain a Steam session ticket via the Steamworks SDK (or a community Go binding to it), then call `Connect.Login` directly with `ExternalCredentialSteamSessionTicket`. **You skip the Auth interface entirely** — the player never has to know EOS exists. Your `Connect.LoginOptions.DisplayName` should be the player's Steam persona name.

### itch.io / direct download

Two reasonable paths:

- Anonymous (Device ID). See `examples/connect-deviceid` and `Connect.CreateDeviceId` + `Connect.Login` with `ExternalCredentialDeviceIDAccessToken`. Zero player friction. Trade-off: each device is a separate identity unless the player later links Steam/Discord/etc.
- Your own auth backend. Issue your players a JWT signed by your backend, then call `Connect.Login` with the appropriate external credential type. Requires you to run an auth server.

### GOG Galaxy / Discord / etc.

The same pattern as Steam: get a session ticket from the platform's SDK, pass it to `Connect.Login` with the matching `ExternalCredential*` type. See `eos/types/enums.go` for the full list of supported credential types.

---

## 5. Credentials in the binary

Your `ClientId` and `ClientSecret` are not "secret" in the OAuth-confidential sense — they ship inside every game client. Treat them like API keys with restricted scope, not like passwords. EOS knows this; the rate-limiting and abuse protections live on Epic's side.

The standard pattern is to bake them in at build time:

```bash
go build -ldflags "
  -X main.eosProductId=$EOS_PRODUCT_ID
  -X main.eosSandboxId=$EOS_SANDBOX_ID
  -X main.eosDeploymentId=$EOS_DEPLOYMENT_ID
  -X main.eosClientId=$EOS_CLIENT_ID
  -X main.eosClientSecret=$EOS_CLIENT_SECRET
" -o mygame ./cmd/mygame
```

Don't commit them to your repo. Use your CI's secrets store.

---

## 6. License and legal

| Topic | Status |
|---|---|
| EOS pricing | Free. No per-CCU fee, no revenue share, no MTX cut. This is the actual reason it exists |
| EOS SDK redistribution | The runtime shared library may ship inside your game. The SDK headers and source archives may **not** be redistributed. This is why `static/` is gitignored in this repo |
| Epic Developer Agreement | You must accept it at portal signup. Includes trust-and-safety obligations standard to any platform service |
| Player data / privacy | EOS stores ProductUserIds and any data you push to it (lobby attributes, leaderboard scores). Disclose in your privacy policy |
| go-eos license | See `LICENSE` in this repository — go-eos itself is open source and unaffiliated with Epic |

---

## 7. Troubleshooting

### "EOSSDK-Win64-Shipping.dll was not found"

The DLL is not next to the `.exe`, or the player is missing the MSVC++ Redistributable. Verify `where.exe EOSSDK-Win64-Shipping.dll` from the game directory. If the redistributable is missing, the error often comes back as a vague Windows error code instead of a clear "missing DLL" — bundle it in your installer.

### "image not found: @rpath/libEOSSDK-Mac-Shipping.dylib"

The dylib's install name or the binary's rpath is wrong. Diagnose with:

```bash
otool -L mygame                                 # what does the binary expect?
otool -D Contents/Frameworks/libEOSSDK-Mac-Shipping.dylib  # what does the dylib say it is?
otool -l mygame | grep -A2 LC_RPATH             # what rpaths are baked in?
```

The expectation, the dylib's self-reported install name, and one of the binary's rpaths must agree. Fix mismatches with `install_name_tool`.

### "error while loading shared libraries: libEOSSDK-Linux-Shipping.so: cannot open shared object file"

The .so is not findable. Either the binary's RPATH/RUNPATH does not include the directory containing the library, or the player's `LD_LIBRARY_PATH` is not set. Verify with:

```bash
readelf -d mygame | grep -E 'PATH|NEEDED'
ldd mygame
```

### "killed: 9" on macOS first launch

Gatekeeper blocked an unsigned/unnotarized binary downloaded from the internet. Either sign and notarize properly, or instruct technical players to run `xattr -cr MyGame.app` to remove the quarantine attribute. The latter is a development-only workaround; do not ship that as a customer instruction for a paid game.

### `EOS_NotConfigured` errors at runtime

Almost always a portal configuration issue, not a code issue. Common causes:

- The Client is missing a permission (e.g., `LobbyCreate`).
- No Application is linked to the Client (causes `Connect.Login` with `EpicIDToken` to fail with this error).
- For Device ID: Device ID is not enabled as an Identity Provider on the Application.

The portal's audit log is the fastest way to see what permission was denied.

### Hangs on first SDK call on macOS

The SDK is being called from a non-main OS thread. Use `platform.RunOnMainThread` and put `runtime.LockOSThread()` in `init()`. See `docs/go-dev-docs.md` for the full explanation.

---

## 8. Going further

- `examples/auth-login` — Epic Account login (most users)
- `examples/lobby-chat` — lobby create/search/join with two processes
- `examples/connect-deviceid` — fully anonymous login, no Epic account
- `docs/go-dev-docs.md` — the macOS main-thread requirement and other developer footguns
