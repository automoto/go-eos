# EOS SDK Version Management

## Current Target

**EOS SDK v1.19.0.3**

The SDK lives in `static/EOS-SDK-49960398-Release-v1.19.0.3/` (gitignored). Developers must download it themselves from the [Epic Developer Portal](https://dev.epicgames.com/portal).

Platform binaries:

| Platform | Library |
|----------|---------|
| macOS | `libEOSSDK-Mac-Shipping.dylib` |
| Linux | `libEOSSDK-Linux-Shipping.so` |
| Windows | `EOSSDK-Win64-Shipping.dll` |

## SDK Directory Structure

```
static/
  EOS-SDK-49960398-Release-v1.19.0.3/
    SDK/
      Include/    # C headers (referenced by CFLAGS in cgo.go)
      Bin/        # Per-platform shared libraries (referenced by LDFLAGS in cgo.go)
```

The `static/` directory is in `.gitignore` -- nothing under it is checked in.

## How to Bump the SDK Version

1. Download the new SDK release from the Epic Developer Portal.
2. Extract it into `static/` (e.g. `static/EOS-SDK-XXXXXXXX-Release-vX.YY.Z.W/`).
3. Update the SDK path in `cgo.go` -- both the `CFLAGS` include path (`SDK/Include`) and the `LDFLAGS` library path (`SDK/Bin`).
4. Re-run `eos-bindgen` (`cmd/eos-bindgen`) to regenerate the cbinding layer from the new headers via c-for-go.
5. Rebuild the project and run the full test suite.
6. Fix any compilation errors from changed or removed API symbols.

## Known Compatibility Notes

- EOS SDK is generally backwards-compatible between minor versions.
- API additions rarely break existing bindings; removals or signature changes are uncommon but possible on major bumps.
- Always review Epic's changelog before upgrading to check for breaking changes.

## Epic's SDK Changelog

Release notes are published on the Epic Developer Portal:
https://dev.epicgames.com/docs/epic-online-services/whats-new
