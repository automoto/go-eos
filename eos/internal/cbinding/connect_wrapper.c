// clang-format off
//go:build !eosstub
// clang-format on

#ifdef EOS_CGO

#include "eos_sdk.h"
#include "eos_connect.h"
#include <stdint.h>

/* Forward declarations for Go export functions (implemented in connect_callback.go) */
extern void goConnectLoginCallback(int resultCode, uintptr_t clientData, uintptr_t localUserId,
								   uintptr_t continuanceToken);
extern void goConnectCreateUserCallback(int resultCode, uintptr_t clientData,
										uintptr_t localUserId);
extern void goConnectLinkAccountCallback(int resultCode, uintptr_t clientData,
										 uintptr_t localUserId);
extern void goConnectAuthExpirationCallback(uintptr_t clientData, uintptr_t localUserId);
extern void goConnectLoginStatusChangedCallback(uintptr_t clientData, uintptr_t localUserId,
												int previousStatus, int currentStatus);

/* Trampolines */

static void connectLoginTrampoline(const EOS_Connect_LoginCallbackInfo* data) {
	goConnectLoginCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
						   (uintptr_t)data->LocalUserId, (uintptr_t)data->ContinuanceToken);
}

static void connectCreateUserTrampoline(const EOS_Connect_CreateUserCallbackInfo* data) {
	goConnectCreateUserCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
								(uintptr_t)data->LocalUserId);
}

static void connectLinkAccountTrampoline(const EOS_Connect_LinkAccountCallbackInfo* data) {
	goConnectLinkAccountCallback((int)data->ResultCode, (uintptr_t)data->ClientData,
								 (uintptr_t)data->LocalUserId);
}

static void connectAuthExpirationTrampoline(const EOS_Connect_AuthExpirationCallbackInfo* data) {
	goConnectAuthExpirationCallback((uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId);
}

static void
connectLoginStatusChangedTrampoline(const EOS_Connect_LoginStatusChangedCallbackInfo* data) {
	goConnectLoginStatusChangedCallback((uintptr_t)data->ClientData, (uintptr_t)data->LocalUserId,
										(int)data->PreviousStatus, (int)data->CurrentStatus);
}

/* Wrapper functions */

void eos_connect_login(uintptr_t handle, int credentialType, const char* token,
					   const char* displayName, uintptr_t clientData) {
	EOS_Connect_Credentials creds = {0};
	creds.ApiVersion = EOS_CONNECT_CREDENTIALS_API_LATEST;
	creds.Type = (EOS_EExternalCredentialType)credentialType;
	creds.Token = token;

	EOS_Connect_LoginOptions opts = {0};
	opts.ApiVersion = EOS_CONNECT_LOGIN_API_LATEST;
	opts.Credentials = &creds;

	if (displayName != NULL) {
		EOS_Connect_UserLoginInfo userInfo = {0};
		userInfo.ApiVersion = EOS_CONNECT_USERLOGININFO_API_LATEST;
		userInfo.DisplayName = displayName;
		opts.UserLoginInfo = &userInfo;
	}

	EOS_Connect_Login((EOS_HConnect)handle, &opts, (void*)clientData, &connectLoginTrampoline);
}

void eos_connect_create_user(uintptr_t handle, uintptr_t continuanceToken, uintptr_t clientData) {
	EOS_Connect_CreateUserOptions opts = {0};
	opts.ApiVersion = EOS_CONNECT_CREATEUSER_API_LATEST;
	opts.ContinuanceToken = (EOS_ContinuanceToken)continuanceToken;
	EOS_Connect_CreateUser((EOS_HConnect)handle, &opts, (void*)clientData,
						   &connectCreateUserTrampoline);
}

void eos_connect_link_account(uintptr_t handle, uintptr_t localUserId, uintptr_t continuanceToken,
							  uintptr_t clientData) {
	EOS_Connect_LinkAccountOptions opts = {0};
	opts.ApiVersion = EOS_CONNECT_LINKACCOUNT_API_LATEST;
	opts.LocalUserId = (EOS_ProductUserId)localUserId;
	opts.ContinuanceToken = (EOS_ContinuanceToken)continuanceToken;
	EOS_Connect_LinkAccount((EOS_HConnect)handle, &opts, (void*)clientData,
							&connectLinkAccountTrampoline);
}

int eos_connect_get_logged_in_users_count(uintptr_t handle) {
	return EOS_Connect_GetLoggedInUsersCount((EOS_HConnect)handle);
}

uintptr_t eos_connect_get_logged_in_user_by_index(uintptr_t handle, int index) {
	return (uintptr_t)EOS_Connect_GetLoggedInUserByIndex((EOS_HConnect)handle, index);
}

uint64_t eos_connect_add_notify_auth_expiration(uintptr_t handle, uintptr_t clientData) {
	EOS_Connect_AddNotifyAuthExpirationOptions opts = {0};
	opts.ApiVersion = EOS_CONNECT_ADDNOTIFYAUTHEXPIRATION_API_LATEST;
	return (uint64_t)EOS_Connect_AddNotifyAuthExpiration(
		(EOS_HConnect)handle, &opts, (void*)clientData, &connectAuthExpirationTrampoline);
}

void eos_connect_remove_notify_auth_expiration(uintptr_t handle, uint64_t id) {
	EOS_Connect_RemoveNotifyAuthExpiration((EOS_HConnect)handle, (EOS_NotificationId)id);
}

uint64_t eos_connect_add_notify_login_status_changed(uintptr_t handle, uintptr_t clientData) {
	EOS_Connect_AddNotifyLoginStatusChangedOptions opts = {0};
	opts.ApiVersion = EOS_CONNECT_ADDNOTIFYLOGINSTATUSCHANGED_API_LATEST;
	return (uint64_t)EOS_Connect_AddNotifyLoginStatusChanged(
		(EOS_HConnect)handle, &opts, (void*)clientData, &connectLoginStatusChangedTrampoline);
}

void eos_connect_remove_notify_login_status_changed(uintptr_t handle, uint64_t id) {
	EOS_Connect_RemoveNotifyLoginStatusChanged((EOS_HConnect)handle, (EOS_NotificationId)id);
}

#endif /* EOS_CGO */
