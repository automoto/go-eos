package main

import "fmt"

func main() {
	fmt.Println(`eos-bindgen: EOS SDK Go binding generator

This tool generates Go bindings from EOS C SDK headers using c-for-go.

Prerequisites:
  - EOS C SDK headers (download from Epic Developer Portal)
  - c-for-go (go install github.com/nicholasgasior/c-for-go@latest)

Usage:
  EOS_SDK_PATH=/path/to/eos-sdk go generate ./eos/internal/cbinding/...

Status: Not yet implemented. See docs/prd.md section 6.4.`)
}
