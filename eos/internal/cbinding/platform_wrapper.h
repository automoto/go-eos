#ifndef PLATFORM_WRAPPER_H
#define PLATFORM_WRAPPER_H

#include <stdint.h>

int eos_initialize(const char* productName, const char* productVersion);
uintptr_t eos_platform_create(const char* productId, const char* sandboxId,
							  const char* deploymentId, const char* clientId,
							  const char* clientSecret);
void eos_platform_tick(uintptr_t handle);
void eos_platform_release(uintptr_t handle);
uintptr_t eos_platform_get_auth(uintptr_t handle);
uintptr_t eos_platform_get_connect(uintptr_t handle);
uintptr_t eos_platform_get_lobby(uintptr_t handle);
uintptr_t eos_platform_get_sessions(uintptr_t handle);
uintptr_t eos_platform_get_p2p(uintptr_t handle);

#endif
