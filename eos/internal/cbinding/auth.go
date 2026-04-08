//go:build !eosstub

package cbinding

/*
#cgo CFLAGS: -I${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Include

#include "eos_sdk.h"
#include "eos_auth.h"
#include <stdlib.h>
#include <stdint.h>

static const char* safe_str(const char* s) { return s ? s : ""; }

// Forward declarations for Go export functions (implemented in auth_callback.go)
extern void goAuthLoginCallback(int resultCode, uintptr_t clientData,
	uintptr_t localUserId, uintptr_t selectedAccountId,
	uintptr_t continuanceToken, int hasPinGrant,
	const char* pinUserCode, const char* pinVerificationURI,
	int pinExpiresIn, const char* pinVerificationURIComplete);
extern void goAuthLogoutCallback(int resultCode, uintptr_t clientData,
	uintptr_t localUserId);
extern void goAuthDeletePersistentAuthCallback(int resultCode, uintptr_t clientData);
extern void goAuthLoginStatusChangedCallback(uintptr_t clientData,
	uintptr_t localUserId, int prevStatus, int currentStatus);

// Trampolines — called by EOS SDK, forward to Go exports with primitive types
static void authLoginTrampoline(const EOS_Auth_LoginCallbackInfo* data) {
	int hasPinGrant = data->PinGrantInfo != NULL;
	goAuthLoginCallback(
		(int)data->ResultCode,
		(uintptr_t)data->ClientData,
		(uintptr_t)data->LocalUserId,
		(uintptr_t)data->SelectedAccountId,
		(uintptr_t)data->ContinuanceToken,
		hasPinGrant,
		hasPinGrant ? safe_str(data->PinGrantInfo->UserCode) : "",
		hasPinGrant ? safe_str(data->PinGrantInfo->VerificationURI) : "",
		hasPinGrant ? data->PinGrantInfo->ExpiresIn : 0,
		hasPinGrant ? safe_str(data->PinGrantInfo->VerificationURIComplete) : ""
	);
}

static void authLogoutTrampoline(const EOS_Auth_LogoutCallbackInfo* data) {
	goAuthLogoutCallback(
		(int)data->ResultCode,
		(uintptr_t)data->ClientData,
		(uintptr_t)data->LocalUserId
	);
}

static void authDeletePersistentAuthTrampoline(const EOS_Auth_DeletePersistentAuthCallbackInfo* data) {
	goAuthDeletePersistentAuthCallback(
		(int)data->ResultCode,
		(uintptr_t)data->ClientData
	);
}

static void authLoginStatusChangedTrampoline(const EOS_Auth_LoginStatusChangedCallbackInfo* data) {
	goAuthLoginStatusChangedCallback(
		(uintptr_t)data->ClientData,
		(uintptr_t)data->LocalUserId,
		(int)data->PrevStatus,
		(int)data->CurrentStatus
	);
}

// C wrapper functions — convert between Go-friendly uintptr and C opaque types
static void eos_auth_login(uintptr_t handle, int credentialType, const char* id,
	const char* token, uint64_t scopeFlags, int externalType, uintptr_t clientData) {
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

static void eos_auth_logout(uintptr_t handle, uintptr_t localUserId, uintptr_t clientData) {
	EOS_Auth_LogoutOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_LOGOUT_API_LATEST;
	opts.LocalUserId = (EOS_EpicAccountId)localUserId;
	EOS_Auth_Logout((EOS_HAuth)handle, &opts, (void*)clientData, &authLogoutTrampoline);
}

static void eos_auth_delete_persistent_auth(uintptr_t handle, const char* refreshToken,
	uintptr_t clientData) {
	EOS_Auth_DeletePersistentAuthOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_DELETEPERSISTENTAUTH_API_LATEST;
	opts.RefreshToken = refreshToken;
	EOS_Auth_DeletePersistentAuth((EOS_HAuth)handle, &opts, (void*)clientData,
		&authDeletePersistentAuthTrampoline);
}

static int eos_auth_get_logged_in_accounts_count(uintptr_t handle) {
	return EOS_Auth_GetLoggedInAccountsCount((EOS_HAuth)handle);
}

static uintptr_t eos_auth_get_logged_in_account_by_index(uintptr_t handle, int index) {
	return (uintptr_t)EOS_Auth_GetLoggedInAccountByIndex((EOS_HAuth)handle, index);
}

// CopyUserAuthToken — copies token, extracts fields, releases the C token
static int eos_auth_copy_user_auth_token(uintptr_t handle, uintptr_t localUserId,
	uintptr_t* outToken) {
	EOS_Auth_CopyUserAuthTokenOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_COPYUSERAUTHTOKEN_API_LATEST;
	EOS_Auth_Token* token = NULL;
	EOS_EResult result = EOS_Auth_CopyUserAuthToken(
		(EOS_HAuth)handle, &opts, (EOS_EpicAccountId)localUserId, &token);
	*outToken = (uintptr_t)token;
	return (int)result;
}

static const char* eos_auth_token_get_app(uintptr_t t) { return safe_str(((EOS_Auth_Token*)t)->App); }
static const char* eos_auth_token_get_client_id(uintptr_t t) { return safe_str(((EOS_Auth_Token*)t)->ClientId); }
static uintptr_t eos_auth_token_get_account_id(uintptr_t t) { return (uintptr_t)((EOS_Auth_Token*)t)->AccountId; }
static const char* eos_auth_token_get_access_token(uintptr_t t) { return safe_str(((EOS_Auth_Token*)t)->AccessToken); }
static double eos_auth_token_get_expires_in(uintptr_t t) { return ((EOS_Auth_Token*)t)->ExpiresIn; }
static const char* eos_auth_token_get_expires_at(uintptr_t t) { return safe_str(((EOS_Auth_Token*)t)->ExpiresAt); }
static int eos_auth_token_get_auth_type(uintptr_t t) { return (int)((EOS_Auth_Token*)t)->AuthType; }
static const char* eos_auth_token_get_refresh_token(uintptr_t t) { return safe_str(((EOS_Auth_Token*)t)->RefreshToken); }
static double eos_auth_token_get_refresh_expires_in(uintptr_t t) { return ((EOS_Auth_Token*)t)->RefreshExpiresIn; }
static const char* eos_auth_token_get_refresh_expires_at(uintptr_t t) { return safe_str(((EOS_Auth_Token*)t)->RefreshExpiresAt); }
static void eos_auth_token_release(uintptr_t t) { EOS_Auth_Token_Release((EOS_Auth_Token*)t); }

// Notification helpers
static uint64_t eos_auth_add_notify_login_status_changed(uintptr_t handle, uintptr_t clientData) {
	EOS_Auth_AddNotifyLoginStatusChangedOptions opts = {0};
	opts.ApiVersion = EOS_AUTH_ADDNOTIFYLOGINSTATUSCHANGED_API_LATEST;
	return (uint64_t)EOS_Auth_AddNotifyLoginStatusChanged(
		(EOS_HAuth)handle, &opts, (void*)clientData, &authLoginStatusChangedTrampoline);
}

static void eos_auth_remove_notify_login_status_changed(uintptr_t handle, uint64_t id) {
	EOS_Auth_RemoveNotifyLoginStatusChanged((EOS_HAuth)handle, (EOS_NotificationId)id);
}
*/
import "C"
import "unsafe"

func EOS_Auth_Login(handle EOS_HAuth, opts *EOS_Auth_LoginOptions, clientData uintptr) {
	var cID, cToken *C.char
	if opts.ID != "" {
		cID = C.CString(opts.ID)
		defer C.free(unsafe.Pointer(cID))
	}
	if opts.Token != "" {
		cToken = C.CString(opts.Token)
		defer C.free(unsafe.Pointer(cToken))
	}
	C.eos_auth_login(C.uintptr_t(handle),
		C.int(opts.CredentialType), cID, cToken,
		C.uint64_t(opts.ScopeFlags), C.int(opts.ExternalType),
		C.uintptr_t(clientData))
}

func EOS_Auth_Logout(handle EOS_HAuth, opts *EOS_Auth_LogoutOptions, clientData uintptr) {
	C.eos_auth_logout(C.uintptr_t(handle), C.uintptr_t(opts.LocalUserId), C.uintptr_t(clientData))
}

func EOS_Auth_DeletePersistentAuth(handle EOS_HAuth, opts *EOS_Auth_DeletePersistentAuthOptions, clientData uintptr) {
	var cRefreshToken *C.char
	if opts.RefreshToken != "" {
		cRefreshToken = C.CString(opts.RefreshToken)
		defer C.free(unsafe.Pointer(cRefreshToken))
	}
	C.eos_auth_delete_persistent_auth(C.uintptr_t(handle), cRefreshToken, C.uintptr_t(clientData))
}

func EOS_Auth_GetLoggedInAccountsCount(handle EOS_HAuth) int32 {
	return int32(C.eos_auth_get_logged_in_accounts_count(C.uintptr_t(handle)))
}

func EOS_Auth_GetLoggedInAccountByIndex(handle EOS_HAuth, index int32) EOS_EpicAccountId {
	return EOS_EpicAccountId(C.eos_auth_get_logged_in_account_by_index(C.uintptr_t(handle), C.int(index)))
}

func EOS_Auth_CopyUserAuthToken(handle EOS_HAuth, localUserId EOS_EpicAccountId) (*EOS_Auth_Token, EOS_EResult) {
	var tokenPtr C.uintptr_t
	result := EOS_EResult(C.eos_auth_copy_user_auth_token(C.uintptr_t(handle), C.uintptr_t(localUserId), &tokenPtr))
	if result != EOS_EResult_Success || tokenPtr == 0 {
		return nil, result
	}
	defer C.eos_auth_token_release(tokenPtr)
	return &EOS_Auth_Token{
		App:              C.GoString(C.eos_auth_token_get_app(tokenPtr)),
		ClientId:         C.GoString(C.eos_auth_token_get_client_id(tokenPtr)),
		AccountId:        EOS_EpicAccountId(C.eos_auth_token_get_account_id(tokenPtr)),
		AccessToken:      C.GoString(C.eos_auth_token_get_access_token(tokenPtr)),
		ExpiresIn:        float64(C.eos_auth_token_get_expires_in(tokenPtr)),
		ExpiresAt:        C.GoString(C.eos_auth_token_get_expires_at(tokenPtr)),
		AuthType:         int32(C.eos_auth_token_get_auth_type(tokenPtr)),
		RefreshToken:     C.GoString(C.eos_auth_token_get_refresh_token(tokenPtr)),
		RefreshExpiresIn: float64(C.eos_auth_token_get_refresh_expires_in(tokenPtr)),
		RefreshExpiresAt: C.GoString(C.eos_auth_token_get_refresh_expires_at(tokenPtr)),
	}, EOS_EResult_Success
}

func EOS_Auth_AddNotifyLoginStatusChanged(handle EOS_HAuth, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_auth_add_notify_login_status_changed(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Auth_RemoveNotifyLoginStatusChanged(handle EOS_HAuth, id EOS_NotificationId) {
	C.eos_auth_remove_notify_login_status_changed(C.uintptr_t(handle), C.uint64_t(id))
}
