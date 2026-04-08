#!/usr/bin/env bash
set -euo pipefail

echo "=== EOS SDK Setup Helper ==="
echo ""
echo "The EOS C SDK must be downloaded from the Epic Developer Portal."
echo "It cannot be redistributed and is not included in this repository."
echo ""
echo "Steps:"
echo "  1. Go to https://dev.epicgames.com/portal"
echo "  2. Navigate to your product > SDK Downloads"
echo "  3. Download the C SDK for your platform(s)"
echo "  4. Extract to a directory and set EOS_SDK_PATH"
echo ""

if [ -z "${EOS_SDK_PATH:-}" ]; then
    echo "WARNING: EOS_SDK_PATH is not set."
    echo ""
    echo "Set it to the root of the extracted EOS SDK:"
    echo "  export EOS_SDK_PATH=/path/to/eos-sdk"
    echo ""
    exit 1
fi

echo "EOS_SDK_PATH is set to: $EOS_SDK_PATH"
echo ""

echo "Checking expected directory structure..."
errors=0

if [ ! -d "$EOS_SDK_PATH/Include" ]; then
    echo "  MISSING: $EOS_SDK_PATH/Include/"
    errors=$((errors + 1))
else
    echo "  OK: $EOS_SDK_PATH/Include/"
fi

if [ ! -d "$EOS_SDK_PATH/Bin" ]; then
    echo "  MISSING: $EOS_SDK_PATH/Bin/"
    errors=$((errors + 1))
else
    echo "  OK: $EOS_SDK_PATH/Bin/"
fi

echo ""
echo "Expected platform libraries:"
echo "  Windows: \$EOS_SDK_PATH/Bin/EOSSDK-Win64-Shipping.dll"
echo "  Linux:   \$EOS_SDK_PATH/Bin/libEOSSDK-Linux-Shipping.so"
echo "  macOS:   \$EOS_SDK_PATH/Bin/libEOSSDK-Mac-Shipping.dylib"
echo ""

if [ "$errors" -gt 0 ]; then
    echo "Found $errors issue(s). Please verify your EOS SDK installation."
    exit 1
fi

echo "SDK directory structure looks correct."
