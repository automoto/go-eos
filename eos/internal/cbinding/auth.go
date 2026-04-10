//go:build !eosstub

package cbinding

/*
#include "auth_wrapper.h"
#include <stdlib.h>
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

func EOS_Auth_CopyIdToken(handle EOS_HAuth, accountId EOS_EpicAccountId) (string, EOS_EResult) {
	var tokenPtr C.uintptr_t
	result := EOS_EResult(C.eos_auth_copy_id_token(C.uintptr_t(handle), C.uintptr_t(accountId), &tokenPtr))
	if result != EOS_EResult_Success || tokenPtr == 0 {
		return "", result
	}
	jwt := C.GoString(C.eos_auth_id_token_get_jwt(tokenPtr))
	C.eos_auth_id_token_release(tokenPtr)
	return jwt, EOS_EResult_Success
}

func EOS_Auth_AddNotifyLoginStatusChanged(handle EOS_HAuth, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_auth_add_notify_login_status_changed(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Auth_RemoveNotifyLoginStatusChanged(handle EOS_HAuth, id EOS_NotificationId) {
	C.eos_auth_remove_notify_login_status_changed(C.uintptr_t(handle), C.uint64_t(id))
}
