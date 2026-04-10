// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_auth.h"
#include "cgo_helpers.h"
#include <stdint.h>

/* Forward declarations for Go export functions (implemented in auth_callback.go) */
extern void goAuthLoginCallback(int resultCode, uintptr_t clientData, uintptr_t localUserId,
								uintptr_t selectedAccountId, uintptr_t continuanceToken,
								int hasPinGrant, const char* pinUserCode,
								const char* pinVerificationURI, int pinExpiresIn,
								const char* pinVerificationURIComplete);
extern void goAuthLogoutCallback(int resultCode, uintptr_t clientData, uintptr_t localUserId);
extern void goAuthDeletePersistentAuthCallback(int resultCode, uintptr_t clientData);
extern void goAuthLoginStatusChangedCallback(uintptr_t clientData, uintptr_t localUserId,
											 int prevStatus, int currentStatus);

/* Trampolines — called by EOS SDK, forward to Go exports with primitive types */

static void authLoginTrampoline(const EOS_Auth_LoginCallbackInfo* data) {
	int hasPinGrant = data->PinGrantInfo != NULL;
	goAuthLoginCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						(uintptr_t)data->LocalUserId, (uintptr_t)data->SelectedAccountId,
						(uintptr_t)data->ContinuanceToken, hasPinGrant,
						hasPinGrant ? safe_str(data->PinGrantInfo->UserCode) : "",
						hasPinGrant ? safe_str(data->PinGrantInfo->VerificationURI) : "",
						hasPinGrant ? data->PinGrantInfo->ExpiresIn : 0,
						hasPinGrant ? safe_str(data->PinGrantInfo->VerificationURIComplete) : "");
}

static void authLogoutTrampoline(const EOS_Auth_LogoutCallbackInfo* data) {
	goAuthLogoutCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						 (uintptr_t)data->LocalUserId);
}

static void
authDeletePersistentAuthTrampoline(const EOS_Auth_DeletePersistentAuthCallbackInfo* data) {
	goAuthDeletePersistentAuthCallback((int)data->ResultCode, (uintptr_t)data->ClientData);
}

static void authLoginStatusChangedTrampoline(const EOS_Auth_LoginStatusChangedCallbackInfo* data) {
	goAuthLoginStatusChangedCallback((uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId,
									 (int)data->PrevStatus, (int)data->CurrentStatus);
}

/* Wrapper functions — convert between Go-friendly uintptr and C opaque types */

void eos_auth_login(uintptr_t handle, int credentialType, const char* id, const char* token,
					uint64_t scopeFlags, int externalType, uintptr_t clientData) {
	EOS_Auth_Credentials creds = {0};
	creds.ApiVersion = EOS_AUTH_CREDENTIALS_API_LATEST;
	creds.Type = (EOS_ELoginCredentialType)credentialType;
	creds.Id = id;
	creds.Token = token;
	creds.ExternalType = (EOS_EExternalCredentialType)externalType;

	EOS_Auth_LoginOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_LOGIN_API_LATEST;
	opts.Credentials = &creds;
	opts.ScopeFlags = (EOS_EAuthScopeFlags)scopeFlags;

	EOS_Auth_Login((EOS_HAuth)handle, &opts, (void*)clientData, &authLoginTrampoline);
}

void eos_auth_logout(uintptr_t handle, uintptr_t localUserId, uintptr_t clientData) {
	EOS_Auth_LogoutOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_LOGOUT_API_LATEST;
	opts.LocalUserId = (EOS_EpicAccountId)localUserId;
	EOS_Auth_Logout((EOS_HAuth)handle, &opts, (void*)clientData, &authLogoutTrampoline);
}

void eos_auth_delete_persistent_auth(uintptr_t handle, const char* refreshToken,
									 uintptr_t clientData) {
	EOS_Auth_DeletePersistentAuthOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_DELETEPERSISTENTAUTH_API_LATEST;
	opts.RefreshToken = refreshToken;
	EOS_Auth_DeletePersistentAuth((EOS_HAuth)handle, &opts, (void*)clientData,
								  &authDeletePersistentAuthTrampoline);
}

int eos_auth_get_logged_in_accounts_count(uintptr_t handle) {
	return EOS_Auth_GetLoggedInAccountsCount((EOS_HAuth)handle);
}

uintptr_t eos_auth_get_logged_in_account_by_index(uintptr_t handle, int index) {
	return (uintptr_t)EOS_Auth_GetLoggedInAccountByIndex((EOS_HAuth)handle, index);
}

int eos_auth_copy_user_auth_token(uintptr_t handle, uintptr_t localUserId, uintptr_t* outToken) {
	EOS_Auth_CopyUserAuthTokenOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_COPYUSERAUTHTOKEN_API_LATEST;
	EOS_Auth_Token* token = NULL;
	EOS_EResult result = EOS_Auth_CopyUserAuthToken((EOS_HAuth)handle, &opts,
													(EOS_EpicAccountId)localUserId, &token);
	*outToken = (uintptr_t)token;
	return (int)result;
}

const char* eos_auth_token_get_app(uintptr_t t) { return safe_str(((EOS_Auth_Token*)t)->App); }
const char* eos_auth_token_get_client_id(uintptr_t t) {
	return safe_str(((EOS_Auth_Token*)t)->ClientId);
}
uintptr_t eos_auth_token_get_account_id(uintptr_t t) {
	return (uintptr_t)((EOS_Auth_Token*)t)->AccountId;
}
const char* eos_auth_token_get_access_token(uintptr_t t) {
	return safe_str(((EOS_Auth_Token*)t)->AccessToken);
}
double eos_auth_token_get_expires_in(uintptr_t t) { return ((EOS_Auth_Token*)t)->ExpiresIn; }
const char* eos_auth_token_get_expires_at(uintptr_t t) {
	return safe_str(((EOS_Auth_Token*)t)->ExpiresAt);
}
int eos_auth_token_get_auth_type(uintptr_t t) { return (int)((EOS_Auth_Token*)t)->AuthType; }
const char* eos_auth_token_get_refresh_token(uintptr_t t) {
	return safe_str(((EOS_Auth_Token*)t)->RefreshToken);
}
double eos_auth_token_get_refresh_expires_in(uintptr_t t) {
	return ((EOS_Auth_Token*)t)->RefreshExpiresIn;
}
const char* eos_auth_token_get_refresh_expires_at(uintptr_t t) {
	return safe_str(((EOS_Auth_Token*)t)->RefreshExpiresAt);
}
void eos_auth_token_release(uintptr_t t) { EOS_Auth_Token_Release((EOS_Auth_Token*)t); }

int eos_auth_copy_id_token(uintptr_t handle, uintptr_t accountId, uintptr_t* outToken) {
	EOS_Auth_CopyIdTokenOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_COPYIDTOKEN_API_LATEST;
	opts.AccountId = (EOS_EpicAccountId)accountId;
	EOS_Auth_IdToken* token = NULL;
	EOS_EResult result = EOS_Auth_CopyIdToken((EOS_HAuth)handle, &opts, &token);
	*outToken = (uintptr_t)token;
	return (int)result;
}

const char* eos_auth_id_token_get_jwt(uintptr_t t) {
	return safe_str(((EOS_Auth_IdToken*)t)->JsonWebToken);
}

uintptr_t eos_auth_id_token_get_account_id(uintptr_t t) {
	return (uintptr_t)((EOS_Auth_IdToken*)t)->AccountId;
}

void eos_auth_id_token_release(uintptr_t t) { EOS_Auth_IdToken_Release((EOS_Auth_IdToken*)t); }

uint64_t eos_auth_add_notify_login_status_changed(uintptr_t handle, uintptr_t clientData) {
	EOS_Auth_AddNotifyLoginStatusChangedOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_ADDNOTIFYLOGINSTATUSCHANGED_API_LATEST;
	return (uint64_t)EOS_Auth_AddNotifyLoginStatusChanged(
		(EOS_HAuth)handle, &opts, (void*)clientData, &authLoginStatusChangedTrampoline);
}

void eos_auth_remove_notify_login_status_changed(uintptr_t handle, uint64_t id) {
	EOS_Auth_RemoveNotifyLoginStatusChanged((EOS_HAuth)handle, (EOS_NotificationId)id);
}

#endif /* EOS_CGO */
