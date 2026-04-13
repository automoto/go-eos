#!/usr/bin/env bash
# package-macos.sh — create a macOS .app bundle for a go-eos game
#
# Usage:
#   ./scripts/package-macos.sh <binary> <dylib> [--sign <identity>]
#
# Example:
#   go build -o mygame ./cmd/mygame
#   ./scripts/package-macos.sh mygame /path/to/libEOSSDK-Mac-Shipping.dylib
#   ./scripts/package-macos.sh mygame /path/to/libEOSSDK-Mac-Shipping.dylib --sign "Developer ID Application: My Company"

set -euo pipefail

usage() {
    echo "Usage: $0 <binary> <dylib> [--sign <identity>]"
    echo ""
    echo "Arguments:"
    echo "  binary    Path to the compiled game binary"
    echo "  dylib     Path to libEOSSDK-Mac-Shipping.dylib"
    echo "  --sign    (Optional) Code signing identity for codesign"
    exit 1
}

if [ $# -lt 2 ]; then
    usage
fi

BINARY="$1"
DYLIB="$2"
SIGN_IDENTITY=""

shift 2
while [ $# -gt 0 ]; do
    case "$1" in
        --sign)
            SIGN_IDENTITY="$2"
            shift 2
            ;;
        *)
            usage
            ;;
    esac
done

if [ ! -f "$BINARY" ]; then
    echo "Error: binary not found: $BINARY"
    exit 1
fi

if [ ! -f "$DYLIB" ]; then
    echo "Error: dylib not found: $DYLIB"
    exit 1
fi

APP_NAME="$(basename "$BINARY")"
BUNDLE="${APP_NAME}.app"

echo "Creating ${BUNDLE}..."

# Create bundle directory structure
rm -rf "$BUNDLE"
mkdir -p "${BUNDLE}/Contents/MacOS"
mkdir -p "${BUNDLE}/Contents/Frameworks"

# Copy binary and dylib
cp "$BINARY" "${BUNDLE}/Contents/MacOS/${APP_NAME}"
cp "$DYLIB" "${BUNDLE}/Contents/Frameworks/"

DYLIB_NAME="$(basename "$DYLIB")"

# Create minimal Info.plist
cat > "${BUNDLE}/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>${APP_NAME}</string>
    <key>CFBundleName</key>
    <string>${APP_NAME}</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.${APP_NAME}</string>
    <key>CFBundleVersion</key>
    <string>1.0.0</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.15</string>
</dict>
</plist>
PLIST

# Fix dylib install name to use @rpath
echo "Fixing dylib install name..."
install_name_tool -id "@rpath/${DYLIB_NAME}" "${BUNDLE}/Contents/Frameworks/${DYLIB_NAME}"

# Remove development rpath(s) from binary and add the bundle-relative one
echo "Fixing binary rpath..."
# Remove all existing rpaths (ignore errors if none exist)
otool -l "${BUNDLE}/Contents/MacOS/${APP_NAME}" | grep -A1 LC_RPATH | grep path | awk '{print $2}' | while read -r rpath; do
    install_name_tool -delete_rpath "$rpath" "${BUNDLE}/Contents/MacOS/${APP_NAME}" 2>/dev/null || true
done
install_name_tool -add_rpath "@executable_path/../Frameworks" "${BUNDLE}/Contents/MacOS/${APP_NAME}"

# Code sign if identity provided
if [ -n "$SIGN_IDENTITY" ]; then
    echo "Signing with identity: ${SIGN_IDENTITY}"
    codesign --force --options runtime --sign "$SIGN_IDENTITY" "${BUNDLE}/Contents/Frameworks/${DYLIB_NAME}"
    codesign --force --options runtime --sign "$SIGN_IDENTITY" "${BUNDLE}/Contents/MacOS/${APP_NAME}"
    echo "Signed. To notarize, run:"
    echo "  xcrun notarytool submit ${BUNDLE}.zip --apple-id <email> --team-id <team> --password <app-password> --wait"
else
    echo "Skipping code signing (no --sign identity provided)"
fi

echo ""
echo "Bundle created: ${BUNDLE}"
echo "  Binary:    ${BUNDLE}/Contents/MacOS/${APP_NAME}"
echo "  Framework: ${BUNDLE}/Contents/Frameworks/${DYLIB_NAME}"
echo ""
echo "Verify with:"
echo "  otool -L ${BUNDLE}/Contents/MacOS/${APP_NAME}"
echo "  otool -l ${BUNDLE}/Contents/MacOS/${APP_NAME} | grep -A2 LC_RPATH"
