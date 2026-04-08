// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_logging.h"

extern void goEOSLogCallback(EOS_LogMessage* Message);

static void eosLogCallbackTrampoline(const EOS_LogMessage* Message) {
	goEOSLogCallback((EOS_LogMessage*)Message);
}

int eos_logging_set_callback(void) {
	return (int)EOS_Logging_SetCallback((EOS_LogMessageFunc)&eosLogCallbackTrampoline);
}

#endif /* EOS_CGO */
