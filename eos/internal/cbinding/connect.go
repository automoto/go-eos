//go:build !eosstub

package cbinding

/*
#include "connect_wrapper.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

func EOS_Connect_Login(handle EOS_HConnect, opts *EOS_Connect_LoginOptions, clientData uintptr) {
	var cToken, cDisplayName *C.char
	if opts.Token != "" {
		cToken = C.CString(opts.Token)
		defer C.free(unsafe.Pointer(cToken))
	}
	if opts.DisplayName != "" {
		cDisplayName = C.CString(opts.DisplayName)
		defer C.free(unsafe.Pointer(cDisplayName))
	}
	C.eos_connect_login(C.uintptr_t(handle), C.int(opts.CredentialType),
		cToken, cDisplayName, C.uintptr_t(clientData))
}

func EOS_Connect_CreateUser(handle EOS_HConnect, opts *EOS_Connect_CreateUserOptions, clientData uintptr) {
	C.eos_connect_create_user(C.uintptr_t(handle), C.uintptr_t(opts.ContinuanceToken),
		C.uintptr_t(clientData))
}

func EOS_Connect_LinkAccount(handle EOS_HConnect, opts *EOS_Connect_LinkAccountOptions, clientData uintptr) {
	C.eos_connect_link_account(C.uintptr_t(handle), C.uintptr_t(opts.LocalUserId),
		C.uintptr_t(opts.ContinuanceToken), C.uintptr_t(clientData))
}

func EOS_Connect_GetLoggedInUsersCount(handle EOS_HConnect) int32 {
	return int32(C.eos_connect_get_logged_in_users_count(C.uintptr_t(handle)))
}

func EOS_Connect_GetLoggedInUserByIndex(handle EOS_HConnect, index int32) EOS_ProductUserId {
	return EOS_ProductUserId(C.eos_connect_get_logged_in_user_by_index(C.uintptr_t(handle), C.int(index)))
}

func EOS_Connect_AddNotifyAuthExpiration(handle EOS_HConnect, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_connect_add_notify_auth_expiration(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Connect_RemoveNotifyAuthExpiration(handle EOS_HConnect, id EOS_NotificationId) {
	C.eos_connect_remove_notify_auth_expiration(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_Connect_AddNotifyLoginStatusChanged(handle EOS_HConnect, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_connect_add_notify_login_status_changed(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Connect_RemoveNotifyLoginStatusChanged(handle EOS_HConnect, id EOS_NotificationId) {
	C.eos_connect_remove_notify_login_status_changed(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_Connect_CreateDeviceId(handle EOS_HConnect, opts *EOS_Connect_CreateDeviceIdOptions, clientData uintptr) {
	cModel := C.CString(opts.DeviceModel)
	defer C.free(unsafe.Pointer(cModel))
	C.eos_connect_create_device_id(C.uintptr_t(handle), cModel, C.uintptr_t(clientData))
}

func EOS_Connect_DeleteDeviceId(handle EOS_HConnect, clientData uintptr) {
	C.eos_connect_delete_device_id(C.uintptr_t(handle), C.uintptr_t(clientData))
}
