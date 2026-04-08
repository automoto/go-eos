#!/usr/bin/env bash
#
# cgo-clang-format.sh — run clang-format on C wrapper files in eos/internal/cbinding.
#
# Usage:
#   ./scripts/cgo-clang-format.sh --check   # exit nonzero if formatting needed (CI)
#   ./scripts/cgo-clang-format.sh --fix     # auto-format C code in place
#
set -euo pipefail

if [ $# -lt 1 ] || { [ "$1" != "--check" ] && [ "$1" != "--fix" ]; }; then
    echo "Usage: $0 --check|--fix"
    exit 1
fi

MODE="$1"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

if ! command -v clang-format &>/dev/null; then
    echo "Error: clang-format not found."
    echo "  macOS:  brew install clang-format"
    echo "  Linux:  apt-get install clang-format"
    exit 1
fi

C_FILES=("$ROOT_DIR"/eos/internal/cbinding/*.c "$ROOT_DIR"/eos/internal/cbinding/*.h)

if [ "$MODE" = "--check" ]; then
    clang-format --style=file:"$ROOT_DIR/.clang-format" --dry-run -Werror "${C_FILES[@]}"
    echo "C formatting OK."
elif [ "$MODE" = "--fix" ]; then
    clang-format --style=file:"$ROOT_DIR/.clang-format" -i "${C_FILES[@]}"
    echo "C files formatted."
fi
