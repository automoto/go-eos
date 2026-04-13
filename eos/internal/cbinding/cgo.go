//go:build !eosstub

package cbinding

/*
#cgo CFLAGS: -DEOS_CGO -I${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Include
#cgo darwin LDFLAGS: -L${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin -lEOSSDK-Mac-Shipping -Wl,-rpath,${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin -framework CoreFoundation
#cgo linux LDFLAGS: -L${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin -lEOSSDK-Linux-Shipping -Wl,-rpath,${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin
#cgo windows LDFLAGS: -L${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin -lEOSSDK-Win64-Shipping
*/
import "C"
