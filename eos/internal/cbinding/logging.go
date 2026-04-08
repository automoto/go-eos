//go:build !eosstub

package cbinding

/*
#include "logging_wrapper.h"
#include "eos_sdk.h"
#include "eos_logging.h"
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
