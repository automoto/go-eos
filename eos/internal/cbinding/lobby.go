//go:build !eosstub

package cbinding

/*
#include "lobby_wrapper.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Core lifecycle

func EOS_Lobby_CreateLobby(handle EOS_HLobby, opts *EOS_Lobby_CreateLobbyOptions, clientData uintptr) {
	cBucket := C.CString(opts.BucketId)
	defer C.free(unsafe.Pointer(cBucket))
	allowInvites := 0
	if opts.AllowInvites {
		allowInvites = 1
	}
	C.eos_lobby_create(C.uintptr_t(handle), C.uintptr_t(opts.LocalUserId),
		C.uint32_t(opts.MaxLobbyMembers), C.int(opts.PermissionLevel),
		C.int(allowInvites), cBucket, C.uintptr_t(clientData))
}

func EOS_Lobby_DestroyLobby(handle EOS_HLobby, opts *EOS_Lobby_DestroyLobbyOptions, clientData uintptr) {
	cLobbyId := C.CString(opts.LobbyId)
	defer C.free(unsafe.Pointer(cLobbyId))
	C.eos_lobby_destroy(C.uintptr_t(handle), C.uintptr_t(opts.LocalUserId),
		cLobbyId, C.uintptr_t(clientData))
}

func EOS_Lobby_JoinLobby(handle EOS_HLobby, opts *EOS_Lobby_JoinLobbyOptions, clientData uintptr) {
	C.eos_lobby_join(C.uintptr_t(handle), C.uintptr_t(opts.LobbyDetailsHandle),
		C.uintptr_t(opts.LocalUserId), C.uintptr_t(clientData))
}

func EOS_Lobby_LeaveLobby(handle EOS_HLobby, opts *EOS_Lobby_LeaveLobbyOptions, clientData uintptr) {
	cLobbyId := C.CString(opts.LobbyId)
	defer C.free(unsafe.Pointer(cLobbyId))
	C.eos_lobby_leave(C.uintptr_t(handle), C.uintptr_t(opts.LocalUserId),
		cLobbyId, C.uintptr_t(clientData))
}

func EOS_Lobby_UpdateLobby(handle EOS_HLobby, opts *EOS_Lobby_UpdateLobbyOptions, clientData uintptr) {
	C.eos_lobby_update(C.uintptr_t(handle), C.uintptr_t(opts.LobbyModificationHandle),
		C.uintptr_t(clientData))
}

func EOS_Lobby_SendInvite(handle EOS_HLobby, opts *EOS_Lobby_SendInviteOptions, clientData uintptr) {
	cLobbyId := C.CString(opts.LobbyId)
	defer C.free(unsafe.Pointer(cLobbyId))
	C.eos_lobby_send_invite(C.uintptr_t(handle), cLobbyId,
		C.uintptr_t(opts.LocalUserId), C.uintptr_t(opts.TargetUserId),
		C.uintptr_t(clientData))
}

func EOS_Lobby_QueryInvites(handle EOS_HLobby, opts *EOS_Lobby_QueryInvitesOptions, clientData uintptr) {
	C.eos_lobby_query_invites(C.uintptr_t(handle), C.uintptr_t(opts.LocalUserId),
		C.uintptr_t(clientData))
}

// Modification handle

func EOS_Lobby_UpdateLobbyModification(handle EOS_HLobby, localUserId EOS_ProductUserId, lobbyId string) (EOS_HLobbyModification, EOS_EResult) {
	cLobbyId := C.CString(lobbyId)
	defer C.free(unsafe.Pointer(cLobbyId))
	var mod C.uintptr_t
	result := EOS_EResult(C.eos_lobby_update_lobby_modification(C.uintptr_t(handle),
		C.uintptr_t(localUserId), cLobbyId, &mod))
	return EOS_HLobbyModification(mod), result
}

func EOS_LobbyModification_SetBucketId(mod EOS_HLobbyModification, bucketId string) EOS_EResult {
	cBucket := C.CString(bucketId)
	defer C.free(unsafe.Pointer(cBucket))
	return EOS_EResult(C.eos_lobby_mod_set_bucket_id(C.uintptr_t(mod), cBucket))
}

func EOS_LobbyModification_SetPermissionLevel(mod EOS_HLobbyModification, level EOS_ELobbyPermissionLevel) EOS_EResult {
	return EOS_EResult(C.eos_lobby_mod_set_permission_level(C.uintptr_t(mod), C.int(level)))
}

func EOS_LobbyModification_SetMaxMembers(mod EOS_HLobbyModification, max uint32) EOS_EResult {
	return EOS_EResult(C.eos_lobby_mod_set_max_members(C.uintptr_t(mod), C.uint32_t(max)))
}

func EOS_LobbyModification_SetInvitesAllowed(mod EOS_HLobbyModification, allowed bool) EOS_EResult {
	v := 0
	if allowed {
		v = 1
	}
	return EOS_EResult(C.eos_lobby_mod_set_invites_allowed(C.uintptr_t(mod), C.int(v)))
}

func EOS_LobbyModification_AddAttributeInt64(mod EOS_HLobbyModification, key string, val int64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_mod_add_attr_int64(C.uintptr_t(mod), cKey, C.int64_t(val), C.int(vis)))
}

func EOS_LobbyModification_AddAttributeDouble(mod EOS_HLobbyModification, key string, val float64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_mod_add_attr_double(C.uintptr_t(mod), cKey, C.double(val), C.int(vis)))
}

func EOS_LobbyModification_AddAttributeBool(mod EOS_HLobbyModification, key string, val bool, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	v := 0
	if val {
		v = 1
	}
	return EOS_EResult(C.eos_lobby_mod_add_attr_bool(C.uintptr_t(mod), cKey, C.int(v), C.int(vis)))
}

func EOS_LobbyModification_AddAttributeString(mod EOS_HLobbyModification, key string, val string, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	return EOS_EResult(C.eos_lobby_mod_add_attr_string(C.uintptr_t(mod), cKey, cVal, C.int(vis)))
}

func EOS_LobbyModification_RemoveAttribute(mod EOS_HLobbyModification, key string) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_mod_remove_attr(C.uintptr_t(mod), cKey))
}

func EOS_LobbyModification_AddMemberAttributeInt64(mod EOS_HLobbyModification, key string, val int64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_mod_add_member_attr_int64(C.uintptr_t(mod), cKey, C.int64_t(val), C.int(vis)))
}

func EOS_LobbyModification_AddMemberAttributeDouble(mod EOS_HLobbyModification, key string, val float64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_mod_add_member_attr_double(C.uintptr_t(mod), cKey, C.double(val), C.int(vis)))
}

func EOS_LobbyModification_AddMemberAttributeBool(mod EOS_HLobbyModification, key string, val bool, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	v := 0
	if val {
		v = 1
	}
	return EOS_EResult(C.eos_lobby_mod_add_member_attr_bool(C.uintptr_t(mod), cKey, C.int(v), C.int(vis)))
}

func EOS_LobbyModification_AddMemberAttributeString(mod EOS_HLobbyModification, key string, val string, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	return EOS_EResult(C.eos_lobby_mod_add_member_attr_string(C.uintptr_t(mod), cKey, cVal, C.int(vis)))
}

func EOS_LobbyModification_RemoveMemberAttribute(mod EOS_HLobbyModification, key string) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_mod_remove_member_attr(C.uintptr_t(mod), cKey))
}

func EOS_LobbyModification_Release(mod EOS_HLobbyModification) {
	C.eos_lobby_mod_release(C.uintptr_t(mod))
}

// Search

func EOS_Lobby_CreateLobbySearch(handle EOS_HLobby, maxResults uint32) (EOS_HLobbySearch, EOS_EResult) {
	var search C.uintptr_t
	result := EOS_EResult(C.eos_lobby_create_search(C.uintptr_t(handle), C.uint32_t(maxResults), &search))
	return EOS_HLobbySearch(search), result
}

func EOS_LobbySearch_Find(search EOS_HLobbySearch, localUserId EOS_ProductUserId, clientData uintptr) {
	C.eos_lobby_search_find(C.uintptr_t(search), C.uintptr_t(localUserId), C.uintptr_t(clientData))
}

func EOS_LobbySearch_SetParameterInt64(search EOS_HLobbySearch, key string, val int64, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_search_set_param_int64(C.uintptr_t(search), cKey, C.int64_t(val), C.int(op)))
}

func EOS_LobbySearch_SetParameterDouble(search EOS_HLobbySearch, key string, val float64, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_lobby_search_set_param_double(C.uintptr_t(search), cKey, C.double(val), C.int(op)))
}

func EOS_LobbySearch_SetParameterBool(search EOS_HLobbySearch, key string, val bool, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	v := 0
	if val {
		v = 1
	}
	return EOS_EResult(C.eos_lobby_search_set_param_bool(C.uintptr_t(search), cKey, C.int(v), C.int(op)))
}

func EOS_LobbySearch_SetParameterString(search EOS_HLobbySearch, key string, val string, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	return EOS_EResult(C.eos_lobby_search_set_param_string(C.uintptr_t(search), cKey, cVal, C.int(op)))
}

func EOS_LobbySearch_SetLobbyId(search EOS_HLobbySearch, lobbyId string) EOS_EResult {
	cId := C.CString(lobbyId)
	defer C.free(unsafe.Pointer(cId))
	return EOS_EResult(C.eos_lobby_search_set_lobby_id(C.uintptr_t(search), cId))
}

func EOS_LobbySearch_SetTargetUserId(search EOS_HLobbySearch, targetUserId EOS_ProductUserId) EOS_EResult {
	return EOS_EResult(C.eos_lobby_search_set_target_user_id(C.uintptr_t(search), C.uintptr_t(targetUserId)))
}

func EOS_LobbySearch_SetMaxResults(search EOS_HLobbySearch, maxResults uint32) EOS_EResult {
	return EOS_EResult(C.eos_lobby_search_set_max_results(C.uintptr_t(search), C.uint32_t(maxResults)))
}

func EOS_LobbySearch_GetSearchResultCount(search EOS_HLobbySearch) uint32 {
	return uint32(C.eos_lobby_search_get_search_result_count(C.uintptr_t(search)))
}

func EOS_LobbySearch_CopySearchResultByIndex(search EOS_HLobbySearch, index uint32) (EOS_HLobbyDetails, EOS_EResult) {
	var details C.uintptr_t
	result := EOS_EResult(C.eos_lobby_search_copy_search_result_by_index(C.uintptr_t(search), C.uint32_t(index), &details))
	return EOS_HLobbyDetails(details), result
}

func EOS_LobbySearch_Release(search EOS_HLobbySearch) {
	C.eos_lobby_search_release(C.uintptr_t(search))
}

// Details

func EOS_Lobby_CopyLobbyDetailsHandle(handle EOS_HLobby, localUserId EOS_ProductUserId, lobbyId string) (EOS_HLobbyDetails, EOS_EResult) {
	cId := C.CString(lobbyId)
	defer C.free(unsafe.Pointer(cId))
	var details C.uintptr_t
	result := EOS_EResult(C.eos_lobby_copy_details_handle(C.uintptr_t(handle), C.uintptr_t(localUserId), cId, &details))
	return EOS_HLobbyDetails(details), result
}

func EOS_LobbyDetails_GetLobbyOwner(details EOS_HLobbyDetails) EOS_ProductUserId {
	return EOS_ProductUserId(C.eos_lobby_details_get_owner(C.uintptr_t(details)))
}

func EOS_LobbyDetails_GetMemberCount(details EOS_HLobbyDetails) uint32 {
	return uint32(C.eos_lobby_details_get_member_count(C.uintptr_t(details)))
}

func EOS_LobbyDetails_GetMemberByIndex(details EOS_HLobbyDetails, index uint32) EOS_ProductUserId {
	return EOS_ProductUserId(C.eos_lobby_details_get_member_by_index(C.uintptr_t(details), C.uint32_t(index)))
}

func EOS_LobbyDetails_GetAttributeCount(details EOS_HLobbyDetails) uint32 {
	return uint32(C.eos_lobby_details_get_attribute_count(C.uintptr_t(details)))
}

func EOS_LobbyDetails_GetMemberAttributeCount(details EOS_HLobbyDetails, targetUserId EOS_ProductUserId) uint32 {
	return uint32(C.eos_lobby_details_get_member_attribute_count(C.uintptr_t(details), C.uintptr_t(targetUserId)))
}

func EOS_LobbyDetails_CopyInfo(details EOS_HLobbyDetails) (*EOS_LobbyDetails_Info, EOS_EResult) {
	var infoPtr C.uintptr_t
	result := EOS_EResult(C.eos_lobby_details_copy_info(C.uintptr_t(details), &infoPtr))
	if result != EOS_EResult_Success || infoPtr == 0 {
		return nil, result
	}
	defer C.eos_lobby_details_info_release(infoPtr)
	return &EOS_LobbyDetails_Info{
		LobbyId:          C.GoString(C.eos_lobby_info_get_lobby_id(infoPtr)),
		LobbyOwnerUserId: EOS_ProductUserId(C.eos_lobby_info_get_owner(infoPtr)),
		PermissionLevel:  EOS_ELobbyPermissionLevel(C.eos_lobby_info_get_permission_level(infoPtr)),
		AvailableSlots:   uint32(C.eos_lobby_info_get_available_slots(infoPtr)),
		MaxMembers:       uint32(C.eos_lobby_info_get_max_members(infoPtr)),
		AllowInvites:     C.eos_lobby_info_get_allow_invites(infoPtr) != 0,
		BucketId:         C.GoString(C.eos_lobby_info_get_bucket_id(infoPtr)),
	}, EOS_EResult_Success
}

func EOS_LobbyDetails_CopyAttributeByIndex(details EOS_HLobbyDetails, index uint32) (*EOS_Lobby_Attribute, EOS_EResult) {
	var attrPtr C.uintptr_t
	result := EOS_EResult(C.eos_lobby_details_copy_attr_by_index(C.uintptr_t(details), C.uint32_t(index), &attrPtr))
	if result != EOS_EResult_Success || attrPtr == 0 {
		return nil, result
	}
	defer C.eos_lobby_attr_release(attrPtr)
	return readLobbyAttribute(attrPtr), EOS_EResult_Success
}

func EOS_LobbyDetails_CopyAttributeByKey(details EOS_HLobbyDetails, key string) (*EOS_Lobby_Attribute, EOS_EResult) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var attrPtr C.uintptr_t
	result := EOS_EResult(C.eos_lobby_details_copy_attr_by_key(C.uintptr_t(details), cKey, &attrPtr))
	if result != EOS_EResult_Success || attrPtr == 0 {
		return nil, result
	}
	defer C.eos_lobby_attr_release(attrPtr)
	return readLobbyAttribute(attrPtr), EOS_EResult_Success
}

func EOS_LobbyDetails_CopyMemberAttributeByIndex(details EOS_HLobbyDetails, targetUserId EOS_ProductUserId, index uint32) (*EOS_Lobby_Attribute, EOS_EResult) {
	var attrPtr C.uintptr_t
	result := EOS_EResult(C.eos_lobby_details_copy_member_attr_by_index(C.uintptr_t(details), C.uintptr_t(targetUserId), C.uint32_t(index), &attrPtr))
	if result != EOS_EResult_Success || attrPtr == 0 {
		return nil, result
	}
	defer C.eos_lobby_attr_release(attrPtr)
	return readLobbyAttribute(attrPtr), EOS_EResult_Success
}

func EOS_LobbyDetails_CopyMemberAttributeByKey(details EOS_HLobbyDetails, targetUserId EOS_ProductUserId, key string) (*EOS_Lobby_Attribute, EOS_EResult) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var attrPtr C.uintptr_t
	result := EOS_EResult(C.eos_lobby_details_copy_member_attr_by_key(C.uintptr_t(details), C.uintptr_t(targetUserId), cKey, &attrPtr))
	if result != EOS_EResult_Success || attrPtr == 0 {
		return nil, result
	}
	defer C.eos_lobby_attr_release(attrPtr)
	return readLobbyAttribute(attrPtr), EOS_EResult_Success
}

func EOS_LobbyDetails_Release(details EOS_HLobbyDetails) {
	C.eos_lobby_details_release(C.uintptr_t(details))
}

// Invite helpers

func EOS_Lobby_GetInviteCount(handle EOS_HLobby, localUserId EOS_ProductUserId) uint32 {
	return uint32(C.eos_lobby_get_invite_count(C.uintptr_t(handle), C.uintptr_t(localUserId)))
}

func EOS_Lobby_GetInviteIdByIndex(handle EOS_HLobby, localUserId EOS_ProductUserId, index uint32) (string, EOS_EResult) {
	var buf [256]C.char
	bufLen := C.int32_t(256)
	result := EOS_EResult(C.eos_lobby_get_invite_id_by_index(C.uintptr_t(handle), C.uintptr_t(localUserId), C.uint32_t(index), &buf[0], &bufLen))
	if result != EOS_EResult_Success {
		return "", result
	}
	return C.GoString(&buf[0]), EOS_EResult_Success
}

func EOS_Lobby_CopyLobbyDetailsHandleByInviteId(handle EOS_HLobby, inviteId string) (EOS_HLobbyDetails, EOS_EResult) {
	cId := C.CString(inviteId)
	defer C.free(unsafe.Pointer(cId))
	var details C.uintptr_t
	result := EOS_EResult(C.eos_lobby_copy_details_handle_by_invite_id(C.uintptr_t(handle), cId, &details))
	return EOS_HLobbyDetails(details), result
}

// Notifications

func EOS_Lobby_AddNotifyLobbyUpdateReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_lobby_add_notify_update_received(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Lobby_RemoveNotifyLobbyUpdateReceived(handle EOS_HLobby, id EOS_NotificationId) {
	C.eos_lobby_remove_notify_update_received(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_Lobby_AddNotifyLobbyMemberUpdateReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_lobby_add_notify_member_update_received(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Lobby_RemoveNotifyLobbyMemberUpdateReceived(handle EOS_HLobby, id EOS_NotificationId) {
	C.eos_lobby_remove_notify_member_update_received(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_Lobby_AddNotifyLobbyMemberStatusReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_lobby_add_notify_member_status_received(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Lobby_RemoveNotifyLobbyMemberStatusReceived(handle EOS_HLobby, id EOS_NotificationId) {
	C.eos_lobby_remove_notify_member_status_received(C.uintptr_t(handle), C.uint64_t(id))
}

func EOS_Lobby_AddNotifyLobbyInviteReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_lobby_add_notify_invite_received(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Lobby_RemoveNotifyLobbyInviteReceived(handle EOS_HLobby, id EOS_NotificationId) {
	C.eos_lobby_remove_notify_invite_received(C.uintptr_t(handle), C.uint64_t(id))
}

// Helper to read a C lobby attribute into a Go struct
func readLobbyAttribute(attrPtr C.uintptr_t) *EOS_Lobby_Attribute {
	attr := &EOS_Lobby_Attribute{
		Key:        C.GoString(C.eos_lobby_attr_get_key(attrPtr)),
		ValueType:  EOS_EAttributeType(C.eos_lobby_attr_get_type(attrPtr)),
		Visibility: EOS_ELobbyAttributeVisibility(C.eos_lobby_attr_get_visibility(attrPtr)),
	}
	switch attr.ValueType {
	case EOS_AT_Int64:
		attr.AsInt64 = int64(C.eos_lobby_attr_get_int64(attrPtr))
	case EOS_AT_Double:
		attr.AsDouble = float64(C.eos_lobby_attr_get_double(attrPtr))
	case EOS_AT_Boolean:
		attr.AsBool = C.eos_lobby_attr_get_bool(attrPtr) != 0
	case EOS_AT_String:
		attr.AsString = C.GoString(C.eos_lobby_attr_get_string(attrPtr))
	}
	return attr
}
