# eos-bindgen

Generates Go bindings from EOS C SDK headers for `eos/internal/cbinding`.

## Status

Not yet implemented. Currently the `cbinding` package contains hand-written stubs.

## Planned Approach

Uses [c-for-go](https://github.com/nicholasgasior/c-for-go) with a YAML manifest to parse EOS C headers and emit Go types, constants, and function wrappers. See PRD requirements GEN-1 through GEN-7.

## Expected SDK Layout

```
$EOS_SDK_PATH/
  Include/          # C headers (eos_sdk.h, eos_auth.h, etc.)
  Bin/
    Win64/          # EOSSDK-Win64-Shipping.dll
    Linux/          # libEOSSDK-Linux-Shipping.so
    Mac/            # libEOSSDK-Mac-Shipping.dylib
```

## Usage (once implemented)

```bash
export EOS_SDK_PATH=/path/to/eos-sdk
go generate ./eos/internal/cbinding/...
```
