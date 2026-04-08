//go:build !eosstub

package cbinding

/*
#cgo CFLAGS: -I${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Include
#include "eos_sdk.h"
#include "eos_logging.h"
*/
import "C"

//export goEOSLogCallback
func goEOSLogCallback(msg *C.EOS_LogMessage) {
	if logCallback == nil {
		return
	}
	logCallback(&EOS_LogMessage{
		Category: C.GoString(msg.Category),
		Level:    EOS_ELogLevel(msg.Level),
		Message:  C.GoString(msg.Message),
	})
}
