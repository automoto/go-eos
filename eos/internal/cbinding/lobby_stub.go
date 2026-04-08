//go:build eosstub

package cbinding

import (
	"runtime/cgo"
	"sync/atomic"

	"github.com/mydev/go-eos/eos/internal/callback"
)

var stubLobbyNotifyCounter uint64

// Core lifecycle

func EOS_Lobby_CreateLobby(handle EOS_HLobby, opts *EOS_Lobby_CreateLobbyOptions, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Lobby_CreateLobbyCallbackInfo{
				ResultCode: EOS_EResult_Success,
				LobbyId:    "stub-lobby-id",
			},
		})
	}()
}

func EOS_Lobby_DestroyLobby(handle EOS_HLobby, opts *EOS_Lobby_DestroyLobbyOptions, clientData uintptr) {
	lobbyId := opts.LobbyId
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Lobby_DestroyLobbyCallbackInfo{
				ResultCode: EOS_EResult_Success,
				LobbyId:    lobbyId,
			},
		})
	}()
}

func EOS_Lobby_JoinLobby(handle EOS_HLobby, opts *EOS_Lobby_JoinLobbyOptions, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Lobby_JoinLobbyCallbackInfo{
				ResultCode: EOS_EResult_Success,
				LobbyId:    "stub-lobby-id",
			},
		})
	}()
}

func EOS_Lobby_LeaveLobby(handle EOS_HLobby, opts *EOS_Lobby_LeaveLobbyOptions, clientData uintptr) {
	lobbyId := opts.LobbyId
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Lobby_LeaveLobbyCallbackInfo{
				ResultCode: EOS_EResult_Success,
				LobbyId:    lobbyId,
			},
		})
	}()
}

func EOS_Lobby_UpdateLobby(handle EOS_HLobby, opts *EOS_Lobby_UpdateLobbyOptions, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Lobby_UpdateLobbyCallbackInfo{
				ResultCode: EOS_EResult_Success,
				LobbyId:    "stub-lobby-id",
			},
		})
	}()
}

func EOS_Lobby_SendInvite(handle EOS_HLobby, opts *EOS_Lobby_SendInviteOptions, clientData uintptr) {
	lobbyId := opts.LobbyId
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Lobby_SendInviteCallbackInfo{
				ResultCode: EOS_EResult_Success,
				LobbyId:    lobbyId,
			},
		})
	}()
}

func EOS_Lobby_QueryInvites(handle EOS_HLobby, opts *EOS_Lobby_QueryInvitesOptions, clientData uintptr) {
	localUserId := opts.LocalUserId
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Lobby_QueryInvitesCallbackInfo{
				ResultCode:  EOS_EResult_Success,
				LocalUserId: localUserId,
			},
		})
	}()
}

// Modification handle

func EOS_Lobby_UpdateLobbyModification(handle EOS_HLobby, localUserId EOS_ProductUserId, lobbyId string) (EOS_HLobbyModification, EOS_EResult) {
	return EOS_HLobbyModification(1), EOS_EResult_Success
}

func EOS_LobbyModification_SetBucketId(mod EOS_HLobbyModification, bucketId string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_SetPermissionLevel(mod EOS_HLobbyModification, level EOS_ELobbyPermissionLevel) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_SetMaxMembers(mod EOS_HLobbyModification, max uint32) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_SetInvitesAllowed(mod EOS_HLobbyModification, allowed bool) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddAttributeInt64(mod EOS_HLobbyModification, key string, val int64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddAttributeDouble(mod EOS_HLobbyModification, key string, val float64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddAttributeBool(mod EOS_HLobbyModification, key string, val bool, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddAttributeString(mod EOS_HLobbyModification, key string, val string, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_RemoveAttribute(mod EOS_HLobbyModification, key string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddMemberAttributeInt64(mod EOS_HLobbyModification, key string, val int64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddMemberAttributeDouble(mod EOS_HLobbyModification, key string, val float64, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddMemberAttributeBool(mod EOS_HLobbyModification, key string, val bool, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_AddMemberAttributeString(mod EOS_HLobbyModification, key string, val string, vis EOS_ELobbyAttributeVisibility) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_RemoveMemberAttribute(mod EOS_HLobbyModification, key string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbyModification_Release(mod EOS_HLobbyModification) {}

// Search

func EOS_Lobby_CreateLobbySearch(handle EOS_HLobby, maxResults uint32) (EOS_HLobbySearch, EOS_EResult) {
	return EOS_HLobbySearch(1), EOS_EResult_Success
}

func EOS_LobbySearch_Find(search EOS_HLobbySearch, localUserId EOS_ProductUserId, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_LobbySearch_FindCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

func EOS_LobbySearch_SetParameterInt64(search EOS_HLobbySearch, key string, val int64, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbySearch_SetParameterDouble(search EOS_HLobbySearch, key string, val float64, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbySearch_SetParameterBool(search EOS_HLobbySearch, key string, val bool, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbySearch_SetParameterString(search EOS_HLobbySearch, key string, val string, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbySearch_SetLobbyId(search EOS_HLobbySearch, lobbyId string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbySearch_SetTargetUserId(search EOS_HLobbySearch, targetUserId EOS_ProductUserId) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbySearch_SetMaxResults(search EOS_HLobbySearch, maxResults uint32) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_LobbySearch_GetSearchResultCount(search EOS_HLobbySearch) uint32         { return 0 }
func EOS_LobbySearch_CopySearchResultByIndex(search EOS_HLobbySearch, index uint32) (EOS_HLobbyDetails, EOS_EResult) {
	return 0, EOS_EResult_NotFound
}
func EOS_LobbySearch_Release(search EOS_HLobbySearch) {}

// Details

func EOS_Lobby_CopyLobbyDetailsHandle(handle EOS_HLobby, localUserId EOS_ProductUserId, lobbyId string) (EOS_HLobbyDetails, EOS_EResult) {
	return 0, EOS_EResult_NotFound
}
func EOS_LobbyDetails_GetLobbyOwner(details EOS_HLobbyDetails) EOS_ProductUserId { return 0 }
func EOS_LobbyDetails_GetMemberCount(details EOS_HLobbyDetails) uint32           { return 0 }
func EOS_LobbyDetails_GetMemberByIndex(details EOS_HLobbyDetails, index uint32) EOS_ProductUserId {
	return 0
}
func EOS_LobbyDetails_GetAttributeCount(details EOS_HLobbyDetails) uint32 { return 0 }
func EOS_LobbyDetails_GetMemberAttributeCount(details EOS_HLobbyDetails, targetUserId EOS_ProductUserId) uint32 {
	return 0
}
func EOS_LobbyDetails_CopyInfo(details EOS_HLobbyDetails) (*EOS_LobbyDetails_Info, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_LobbyDetails_CopyAttributeByIndex(details EOS_HLobbyDetails, index uint32) (*EOS_Lobby_Attribute, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_LobbyDetails_CopyAttributeByKey(details EOS_HLobbyDetails, key string) (*EOS_Lobby_Attribute, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_LobbyDetails_CopyMemberAttributeByIndex(details EOS_HLobbyDetails, targetUserId EOS_ProductUserId, index uint32) (*EOS_Lobby_Attribute, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_LobbyDetails_CopyMemberAttributeByKey(details EOS_HLobbyDetails, targetUserId EOS_ProductUserId, key string) (*EOS_Lobby_Attribute, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_LobbyDetails_Release(details EOS_HLobbyDetails) {}

// Invite helpers

func EOS_Lobby_GetInviteCount(handle EOS_HLobby, localUserId EOS_ProductUserId) uint32 {
	return 0
}
func EOS_Lobby_GetInviteIdByIndex(handle EOS_HLobby, localUserId EOS_ProductUserId, index uint32) (string, EOS_EResult) {
	return "", EOS_EResult_NotFound
}
func EOS_Lobby_CopyLobbyDetailsHandleByInviteId(handle EOS_HLobby, inviteId string) (EOS_HLobbyDetails, EOS_EResult) {
	return 0, EOS_EResult_NotFound
}

// Notifications

func EOS_Lobby_AddNotifyLobbyUpdateReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubLobbyNotifyCounter, 1))
}
func EOS_Lobby_RemoveNotifyLobbyUpdateReceived(handle EOS_HLobby, id EOS_NotificationId) {}
func EOS_Lobby_AddNotifyLobbyMemberUpdateReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubLobbyNotifyCounter, 1))
}
func EOS_Lobby_RemoveNotifyLobbyMemberUpdateReceived(handle EOS_HLobby, id EOS_NotificationId) {}
func EOS_Lobby_AddNotifyLobbyMemberStatusReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubLobbyNotifyCounter, 1))
}
func EOS_Lobby_RemoveNotifyLobbyMemberStatusReceived(handle EOS_HLobby, id EOS_NotificationId) {}
func EOS_Lobby_AddNotifyLobbyInviteReceived(handle EOS_HLobby, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubLobbyNotifyCounter, 1))
}
func EOS_Lobby_RemoveNotifyLobbyInviteReceived(handle EOS_HLobby, id EOS_NotificationId) {}
