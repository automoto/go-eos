//go:build windows && amd64

package cbinding

// When real EOS SDK bindings are generated, uncomment these directives:
//
// #cgo CFLAGS: -I${EOS_SDK_PATH}/Include
// #cgo LDFLAGS: -L${EOS_SDK_PATH}/Bin -lEOSSDK-Win64-Shipping
//
// TODO: These directives will be activated when c-for-go generates
// real Cgo bindings from the EOS SDK headers.
