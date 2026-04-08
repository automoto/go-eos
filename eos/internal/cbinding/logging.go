//go:build !eosstub

package cbinding

/*
#cgo CFLAGS: -I${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Include

#include "eos_sdk.h"
#include "eos_logging.h"

// Trampoline implemented in logging_callback.go via Cgo export
extern void goEOSLogCallback(EOS_LogMessage* Message);

static void eosLogCallbackTrampoline(const EOS_LogMessage* Message) {
    goEOSLogCallback((EOS_LogMessage*)Message);
}

static EOS_EResult eos_logging_set_callback() {
    return EOS_Logging_SetCallback((EOS_LogMessageFunc)&eosLogCallbackTrampoline);
}
*/
import "C"

type EOS_LogMessage struct {
	Category string
	Level    EOS_ELogLevel
	Message  string
}

type EOS_LogMessageFunc func(msg *EOS_LogMessage)

var logCallback EOS_LogMessageFunc

func EOS_Logging_SetCallback(fn EOS_LogMessageFunc) EOS_EResult {
	logCallback = fn
	return EOS_EResult(C.eos_logging_set_callback())
}

func EOS_Logging_SetLogLevel(category EOS_ELogCategory, level EOS_ELogLevel) EOS_EResult {
	return EOS_EResult(C.EOS_Logging_SetLogLevel(C.EOS_ELogCategory(category), C.EOS_ELogLevel(level)))
}

// SimulateLogMessage is a test helper that fires the registered log callback.
func SimulateLogMessage(msg *EOS_LogMessage) {
	if logCallback != nil {
		logCallback(msg)
	}
}
