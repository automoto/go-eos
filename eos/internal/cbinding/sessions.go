//go:build !eosstub

package cbinding

/*
#include "sessions_wrapper.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// Core lifecycle

func EOS_Sessions_UpdateSession(handle EOS_HSessions, modHandle EOS_HSessionModification, clientData uintptr) {
	C.eos_sessions_update_session(C.uintptr_t(handle), C.uintptr_t(modHandle), C.uintptr_t(clientData))
}

func EOS_Sessions_DestroySession(handle EOS_HSessions, sessionName string, clientData uintptr) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	C.eos_sessions_destroy_session(C.uintptr_t(handle), cName, C.uintptr_t(clientData))
}

func EOS_Sessions_JoinSession(handle EOS_HSessions, sessionName string, sessionDetails EOS_HSessionDetails, localUserId EOS_ProductUserId, clientData uintptr) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	C.eos_sessions_join_session(C.uintptr_t(handle), cName, C.uintptr_t(sessionDetails), C.uintptr_t(localUserId), C.uintptr_t(clientData))
}

func EOS_Sessions_StartSession(handle EOS_HSessions, sessionName string, clientData uintptr) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	C.eos_sessions_start_session(C.uintptr_t(handle), cName, C.uintptr_t(clientData))
}

func EOS_Sessions_EndSession(handle EOS_HSessions, sessionName string, clientData uintptr) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	C.eos_sessions_end_session(C.uintptr_t(handle), cName, C.uintptr_t(clientData))
}

func EOS_Sessions_RegisterPlayers(handle EOS_HSessions, sessionName string, playerIds []EOS_ProductUserId, clientData uintptr) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	cIds := make([]C.uintptr_t, len(playerIds))
	for i, id := range playerIds {
		cIds[i] = C.uintptr_t(id)
	}
	C.eos_sessions_register_players(C.uintptr_t(handle), cName, &cIds[0], C.uint32_t(len(playerIds)), C.uintptr_t(clientData))
}

func EOS_Sessions_UnregisterPlayers(handle EOS_HSessions, sessionName string, playerIds []EOS_ProductUserId, clientData uintptr) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	cIds := make([]C.uintptr_t, len(playerIds))
	for i, id := range playerIds {
		cIds[i] = C.uintptr_t(id)
	}
	C.eos_sessions_unregister_players(C.uintptr_t(handle), cName, &cIds[0], C.uint32_t(len(playerIds)), C.uintptr_t(clientData))
}

// Modification handle

func EOS_Sessions_CreateSessionModification(handle EOS_HSessions, opts *EOS_Sessions_CreateSessionModificationOptions) (EOS_HSessionModification, EOS_EResult) {
	cName := C.CString(opts.SessionName)
	defer C.free(unsafe.Pointer(cName))
	cBucket := C.CString(opts.BucketId)
	defer C.free(unsafe.Pointer(cBucket))
	var mod C.uintptr_t
	result := EOS_EResult(C.eos_sessions_create_session_modification(C.uintptr_t(handle), cName, cBucket, C.uint32_t(opts.MaxPlayers), C.uintptr_t(opts.LocalUserId), &mod))
	return EOS_HSessionModification(mod), result
}

func EOS_Sessions_UpdateSessionModification(handle EOS_HSessions, sessionName string) (EOS_HSessionModification, EOS_EResult) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	var mod C.uintptr_t
	result := EOS_EResult(C.eos_sessions_update_session_modification(C.uintptr_t(handle), cName, &mod))
	return EOS_HSessionModification(mod), result
}

func EOS_SessionModification_SetBucketId(mod EOS_HSessionModification, bucketId string) EOS_EResult {
	c := C.CString(bucketId)
	defer C.free(unsafe.Pointer(c))
	return EOS_EResult(C.eos_session_mod_set_bucket_id(C.uintptr_t(mod), c))
}

func EOS_SessionModification_SetPermissionLevel(mod EOS_HSessionModification, level EOS_EOnlineSessionPermissionLevel) EOS_EResult {
	return EOS_EResult(C.eos_session_mod_set_permission_level(C.uintptr_t(mod), C.int(level)))
}

func EOS_SessionModification_SetMaxPlayers(mod EOS_HSessionModification, max uint32) EOS_EResult {
	return EOS_EResult(C.eos_session_mod_set_max_players(C.uintptr_t(mod), C.uint32_t(max)))
}

func EOS_SessionModification_SetJoinInProgressAllowed(mod EOS_HSessionModification, allowed bool) EOS_EResult {
	v := 0
	if allowed {
		v = 1
	}
	return EOS_EResult(C.eos_session_mod_set_join_in_progress_allowed(C.uintptr_t(mod), C.int(v)))
}

func EOS_SessionModification_SetInvitesAllowed(mod EOS_HSessionModification, allowed bool) EOS_EResult {
	v := 0
	if allowed {
		v = 1
	}
	return EOS_EResult(C.eos_session_mod_set_invites_allowed(C.uintptr_t(mod), C.int(v)))
}

func EOS_SessionModification_SetHostAddress(mod EOS_HSessionModification, addr string) EOS_EResult {
	c := C.CString(addr)
	defer C.free(unsafe.Pointer(c))
	return EOS_EResult(C.eos_session_mod_set_host_address(C.uintptr_t(mod), c))
}

func EOS_SessionModification_AddAttributeInt64(mod EOS_HSessionModification, key string, val int64, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_session_mod_add_attr_int64(C.uintptr_t(mod), cKey, C.int64_t(val), C.int(advType)))
}

func EOS_SessionModification_AddAttributeDouble(mod EOS_HSessionModification, key string, val float64, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_session_mod_add_attr_double(C.uintptr_t(mod), cKey, C.double(val), C.int(advType)))
}

func EOS_SessionModification_AddAttributeBool(mod EOS_HSessionModification, key string, val bool, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	v := 0
	if val {
		v = 1
	}
	return EOS_EResult(C.eos_session_mod_add_attr_bool(C.uintptr_t(mod), cKey, C.int(v), C.int(advType)))
}

func EOS_SessionModification_AddAttributeString(mod EOS_HSessionModification, key string, val string, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	return EOS_EResult(C.eos_session_mod_add_attr_string(C.uintptr_t(mod), cKey, cVal, C.int(advType)))
}

func EOS_SessionModification_RemoveAttribute(mod EOS_HSessionModification, key string) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_session_mod_remove_attr(C.uintptr_t(mod), cKey))
}

func EOS_SessionModification_Release(mod EOS_HSessionModification) {
	C.eos_session_mod_release(C.uintptr_t(mod))
}

// Search

func EOS_Sessions_CreateSessionSearch(handle EOS_HSessions, maxResults uint32) (EOS_HSessionSearch, EOS_EResult) {
	var search C.uintptr_t
	result := EOS_EResult(C.eos_sessions_create_search(C.uintptr_t(handle), C.uint32_t(maxResults), &search))
	return EOS_HSessionSearch(search), result
}

func EOS_SessionSearch_Find(search EOS_HSessionSearch, localUserId EOS_ProductUserId, clientData uintptr) {
	C.eos_session_search_find(C.uintptr_t(search), C.uintptr_t(localUserId), C.uintptr_t(clientData))
}

func EOS_SessionSearch_SetParameterInt64(search EOS_HSessionSearch, key string, val int64, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_session_search_set_param_int64(C.uintptr_t(search), cKey, C.int64_t(val), C.int(op)))
}

func EOS_SessionSearch_SetParameterDouble(search EOS_HSessionSearch, key string, val float64, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	return EOS_EResult(C.eos_session_search_set_param_double(C.uintptr_t(search), cKey, C.double(val), C.int(op)))
}

func EOS_SessionSearch_SetParameterBool(search EOS_HSessionSearch, key string, val bool, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	v := 0
	if val {
		v = 1
	}
	return EOS_EResult(C.eos_session_search_set_param_bool(C.uintptr_t(search), cKey, C.int(v), C.int(op)))
}

func EOS_SessionSearch_SetParameterString(search EOS_HSessionSearch, key string, val string, op EOS_EComparisonOp) EOS_EResult {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	return EOS_EResult(C.eos_session_search_set_param_string(C.uintptr_t(search), cKey, cVal, C.int(op)))
}

func EOS_SessionSearch_SetSessionId(search EOS_HSessionSearch, sessionId string) EOS_EResult {
	c := C.CString(sessionId)
	defer C.free(unsafe.Pointer(c))
	return EOS_EResult(C.eos_session_search_set_session_id(C.uintptr_t(search), c))
}

func EOS_SessionSearch_SetTargetUserId(search EOS_HSessionSearch, targetUserId EOS_ProductUserId) EOS_EResult {
	return EOS_EResult(C.eos_session_search_set_target_user_id(C.uintptr_t(search), C.uintptr_t(targetUserId)))
}

func EOS_SessionSearch_SetMaxResults(search EOS_HSessionSearch, maxResults uint32) EOS_EResult {
	return EOS_EResult(C.eos_session_search_set_max_results(C.uintptr_t(search), C.uint32_t(maxResults)))
}

func EOS_SessionSearch_GetSearchResultCount(search EOS_HSessionSearch) uint32 {
	return uint32(C.eos_session_search_get_result_count(C.uintptr_t(search)))
}

func EOS_SessionSearch_CopySearchResultByIndex(search EOS_HSessionSearch, index uint32) (EOS_HSessionDetails, EOS_EResult) {
	var details C.uintptr_t
	result := EOS_EResult(C.eos_session_search_copy_result_by_index(C.uintptr_t(search), C.uint32_t(index), &details))
	return EOS_HSessionDetails(details), result
}

func EOS_SessionSearch_Release(search EOS_HSessionSearch) {
	C.eos_session_search_release(C.uintptr_t(search))
}

// Session details

func EOS_SessionDetails_CopyInfo(details EOS_HSessionDetails) (*EOS_SessionDetails_Info, EOS_EResult) {
	var infoPtr C.uintptr_t
	result := EOS_EResult(C.eos_session_details_copy_info(C.uintptr_t(details), &infoPtr))
	if result != EOS_EResult_Success || infoPtr == 0 {
		return nil, result
	}
	defer C.eos_session_details_info_release(infoPtr)
	return &EOS_SessionDetails_Info{
		SessionId:                C.GoString(C.eos_session_info_get_session_id(infoPtr)),
		HostAddress:              C.GoString(C.eos_session_info_get_host_address(infoPtr)),
		NumOpenPublicConnections: uint32(C.eos_session_info_get_num_open_connections(infoPtr)),
		OwnerUserId:              EOS_ProductUserId(C.eos_session_info_get_owner(infoPtr)),
		BucketId:                 C.GoString(C.eos_session_info_get_bucket_id(infoPtr)),
		NumPublicConnections:     uint32(C.eos_session_info_get_num_connections(infoPtr)),
		AllowJoinInProgress:      C.eos_session_info_get_allow_join_in_progress(infoPtr) != 0,
		PermissionLevel:          EOS_EOnlineSessionPermissionLevel(C.eos_session_info_get_permission_level(infoPtr)),
		InvitesAllowed:           C.eos_session_info_get_invites_allowed(infoPtr) != 0,
	}, EOS_EResult_Success
}

func EOS_SessionDetails_GetSessionAttributeCount(details EOS_HSessionDetails) uint32 {
	return uint32(C.eos_session_details_get_attribute_count(C.uintptr_t(details)))
}

func EOS_SessionDetails_CopySessionAttributeByIndex(details EOS_HSessionDetails, index uint32) (*EOS_Sessions_Attribute, EOS_EResult) {
	var attrPtr C.uintptr_t
	result := EOS_EResult(C.eos_session_details_copy_attr_by_index(C.uintptr_t(details), C.uint32_t(index), &attrPtr))
	if result != EOS_EResult_Success || attrPtr == 0 {
		return nil, result
	}
	defer C.eos_session_attr_release(attrPtr)
	return readSessionAttribute(attrPtr), EOS_EResult_Success
}

func EOS_SessionDetails_CopySessionAttributeByKey(details EOS_HSessionDetails, key string) (*EOS_Sessions_Attribute, EOS_EResult) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var attrPtr C.uintptr_t
	result := EOS_EResult(C.eos_session_details_copy_attr_by_key(C.uintptr_t(details), cKey, &attrPtr))
	if result != EOS_EResult_Success || attrPtr == 0 {
		return nil, result
	}
	defer C.eos_session_attr_release(attrPtr)
	return readSessionAttribute(attrPtr), EOS_EResult_Success
}

func EOS_SessionDetails_Release(details EOS_HSessionDetails) {
	C.eos_session_details_release(C.uintptr_t(details))
}

// Active session

func EOS_Sessions_CopyActiveSessionHandle(handle EOS_HSessions, sessionName string) (EOS_HActiveSession, EOS_EResult) {
	cName := C.CString(sessionName)
	defer C.free(unsafe.Pointer(cName))
	var active C.uintptr_t
	result := EOS_EResult(C.eos_sessions_copy_active_session_handle(C.uintptr_t(handle), cName, &active))
	return EOS_HActiveSession(active), result
}

func EOS_ActiveSession_GetRegisteredPlayerCount(active EOS_HActiveSession) uint32 {
	return uint32(C.eos_active_session_get_registered_player_count(C.uintptr_t(active)))
}

func EOS_ActiveSession_GetRegisteredPlayerByIndex(active EOS_HActiveSession, index uint32) EOS_ProductUserId {
	return EOS_ProductUserId(C.eos_active_session_get_registered_player_by_index(C.uintptr_t(active), C.uint32_t(index)))
}

func EOS_ActiveSession_Release(active EOS_HActiveSession) {
	C.eos_active_session_release(C.uintptr_t(active))
}

// Notifications

func EOS_Sessions_AddNotifySessionInviteReceived(handle EOS_HSessions, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(C.eos_sessions_add_notify_invite_received(C.uintptr_t(handle), C.uintptr_t(clientData)))
}

func EOS_Sessions_RemoveNotifySessionInviteReceived(handle EOS_HSessions, id EOS_NotificationId) {
	C.eos_sessions_remove_notify_invite_received(C.uintptr_t(handle), C.uint64_t(id))
}

// Helper
func readSessionAttribute(attrPtr C.uintptr_t) *EOS_Sessions_Attribute {
	attr := &EOS_Sessions_Attribute{
		Key:               C.GoString(C.eos_session_attr_get_key(attrPtr)),
		ValueType:         EOS_EAttributeType(C.eos_session_attr_get_type(attrPtr)),
		AdvertisementType: EOS_ESessionAttributeAdvertisementType(C.eos_session_attr_get_advertisement_type(attrPtr)),
	}
	switch attr.ValueType {
	case EOS_AT_Int64:
		attr.AsInt64 = int64(C.eos_session_attr_get_int64(attrPtr))
	case EOS_AT_Double:
		attr.AsDouble = float64(C.eos_session_attr_get_double(attrPtr))
	case EOS_AT_Boolean:
		attr.AsBool = C.eos_session_attr_get_bool(attrPtr) != 0
	case EOS_AT_String:
		attr.AsString = C.GoString(C.eos_session_attr_get_string(attrPtr))
	}
	return attr
}
