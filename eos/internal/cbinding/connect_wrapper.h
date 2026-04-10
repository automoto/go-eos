#ifndef CONNECT_WRAPPER_H
#define CONNECT_WRAPPER_H

#include <stdint.h>

void eos_connect_login(uintptr_t handle, int credentialType, const char* token,
					   const char* displayName, uintptr_t clientData);
void eos_connect_create_user(uintptr_t handle, uintptr_t continuanceToken, uintptr_t clientData);
void eos_connect_link_account(uintptr_t handle, uintptr_t localUserId, uintptr_t continuanceToken,
							  uintptr_t clientData);
int eos_connect_get_logged_in_users_count(uintptr_t handle);
uintptr_t eos_connect_get_logged_in_user_by_index(uintptr_t handle, int index);

uint64_t eos_connect_add_notify_auth_expiration(uintptr_t handle, uintptr_t clientData);
void eos_connect_remove_notify_auth_expiration(uintptr_t handle, uint64_t id);
uint64_t eos_connect_add_notify_login_status_changed(uintptr_t handle, uintptr_t clientData);
void eos_connect_remove_notify_login_status_changed(uintptr_t handle, uint64_t id);

void eos_connect_create_device_id(uintptr_t handle, const char* deviceModel, uintptr_t clientData);
void eos_connect_delete_device_id(uintptr_t handle, uintptr_t clientData);

#endif
