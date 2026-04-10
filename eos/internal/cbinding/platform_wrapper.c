// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_init.h"
#include "eos_logging.h"
#include <stdint.h>

#ifdef __APPLE__
#include <CoreFoundation/CoreFoundation.h>
#endif

int eos_initialize(const char* productName, const char* productVersion) {
	EOS_InitializeOptions opts = {0};
	opts.ApiVersion = EOS_INITIALIZE_API_LATEST;
	opts.ProductName = productName;
	opts.ProductVersion = productVersion;
	return (int)EOS_Initialize(&opts);
}

uintptr_t eos_platform_create(const char* productId, const char* sandboxId,
							  const char* deploymentId, const char* clientId,
							  const char* clientSecret) {
	EOS_Platform_Options opts = {0};
	opts.ApiVersion = EOS_PLATFORM_OPTIONS_API_LATEST;
	opts.ProductId = productId;
	opts.SandboxId = sandboxId;
	opts.DeploymentId = deploymentId;
	opts.ClientCredentials.ClientId = clientId;
	opts.ClientCredentials.ClientSecret = clientSecret;
	opts.bIsServer = EOS_FALSE;
	return (uintptr_t)EOS_Platform_Create(&opts);
}

void eos_platform_tick(uintptr_t handle) {
	EOS_Platform_Tick((EOS_HPlatform)handle);
#ifdef __APPLE__
	// Pump the macOS run loop so SDK networking callbacks fire.
	// The EOS SDK registers CFRunLoopSources on the calling thread;
	// Go does not drive the run loop, so we must pump it here.
	while (CFRunLoopRunInMode(kCFRunLoopDefaultMode, 0, true) == kCFRunLoopRunHandledSource) {
	}
#endif
}

void eos_platform_release(uintptr_t handle) { EOS_Platform_Release((EOS_HPlatform)handle); }

uintptr_t eos_platform_get_auth(uintptr_t handle) {
	return (uintptr_t)EOS_Platform_GetAuthInterface((EOS_HPlatform)handle);
}

uintptr_t eos_platform_get_connect(uintptr_t handle) {
	return (uintptr_t)EOS_Platform_GetConnectInterface((EOS_HPlatform)handle);
}

uintptr_t eos_platform_get_lobby(uintptr_t handle) {
	return (uintptr_t)EOS_Platform_GetLobbyInterface((EOS_HPlatform)handle);
}

uintptr_t eos_platform_get_sessions(uintptr_t handle) {
	return (uintptr_t)EOS_Platform_GetSessionsInterface((EOS_HPlatform)handle);
}

uintptr_t eos_platform_get_p2p(uintptr_t handle) {
	return (uintptr_t)EOS_Platform_GetP2PInterface((EOS_HPlatform)handle);
}

#endif /* EOS_CGO */
