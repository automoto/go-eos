#ifndef AUTH_WRAPPER_H
#define AUTH_WRAPPER_H

#include <stdint.h>

void eos_auth_login(uintptr_t handle, int credentialType, const char* id, const char* token,
					uint64_t scopeFlags, int externalType, uintptr_t clientData);
void eos_auth_logout(uintptr_t handle, uintptr_t localUserId, uintptr_t clientData);
void eos_auth_delete_persistent_auth(uintptr_t handle, const char* refreshToken,
									 uintptr_t clientData);
int eos_auth_get_logged_in_accounts_count(uintptr_t handle);
uintptr_t eos_auth_get_logged_in_account_by_index(uintptr_t handle, int index);
int eos_auth_copy_user_auth_token(uintptr_t handle, uintptr_t localUserId, uintptr_t* outToken);

const char* eos_auth_token_get_app(uintptr_t t);
const char* eos_auth_token_get_client_id(uintptr_t t);
uintptr_t eos_auth_token_get_account_id(uintptr_t t);
const char* eos_auth_token_get_access_token(uintptr_t t);
double eos_auth_token_get_expires_in(uintptr_t t);
const char* eos_auth_token_get_expires_at(uintptr_t t);
int eos_auth_token_get_auth_type(uintptr_t t);
const char* eos_auth_token_get_refresh_token(uintptr_t t);
double eos_auth_token_get_refresh_expires_in(uintptr_t t);
const char* eos_auth_token_get_refresh_expires_at(uintptr_t t);
void eos_auth_token_release(uintptr_t t);

int eos_auth_copy_id_token(uintptr_t handle, uintptr_t accountId, uintptr_t* outToken);
const char* eos_auth_id_token_get_jwt(uintptr_t t);
uintptr_t eos_auth_id_token_get_account_id(uintptr_t t);
void eos_auth_id_token_release(uintptr_t t);

uint64_t eos_auth_add_notify_login_status_changed(uintptr_t handle, uintptr_t clientData);
void eos_auth_remove_notify_login_status_changed(uintptr_t handle, uint64_t id);

#endif
