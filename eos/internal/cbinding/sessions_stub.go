//go:build eosstub

package cbinding

import (
	"runtime/cgo"
	"sync/atomic"

	"github.com/mydev/go-eos/eos/internal/callback"
)

var stubSessionsNotifyCounter uint64

// Core lifecycle

func EOS_Sessions_UpdateSession(handle EOS_HSessions, modHandle EOS_HSessionModification, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Sessions_UpdateSessionCallbackInfo{
				ResultCode:  EOS_EResult_Success,
				SessionName: "stub-session",
				SessionId:   "stub-session-id",
			},
		})
	}()
}

func EOS_Sessions_DestroySession(handle EOS_HSessions, sessionName string, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_Sessions_DestroySessionCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

func EOS_Sessions_JoinSession(handle EOS_HSessions, sessionName string, sessionDetails EOS_HSessionDetails, localUserId EOS_ProductUserId, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_Sessions_JoinSessionCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

func EOS_Sessions_StartSession(handle EOS_HSessions, sessionName string, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_Sessions_StartSessionCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

func EOS_Sessions_EndSession(handle EOS_HSessions, sessionName string, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_Sessions_EndSessionCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

func EOS_Sessions_RegisterPlayers(handle EOS_HSessions, sessionName string, playerIds []EOS_ProductUserId, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_Sessions_RegisterPlayersCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

func EOS_Sessions_UnregisterPlayers(handle EOS_HSessions, sessionName string, playerIds []EOS_ProductUserId, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_Sessions_UnregisterPlayersCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

// Modification handle

func EOS_Sessions_CreateSessionModification(handle EOS_HSessions, opts *EOS_Sessions_CreateSessionModificationOptions) (EOS_HSessionModification, EOS_EResult) {
	return EOS_HSessionModification(1), EOS_EResult_Success
}

func EOS_Sessions_UpdateSessionModification(handle EOS_HSessions, sessionName string) (EOS_HSessionModification, EOS_EResult) {
	return EOS_HSessionModification(1), EOS_EResult_Success
}

func EOS_SessionModification_SetBucketId(mod EOS_HSessionModification, bucketId string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_SetPermissionLevel(mod EOS_HSessionModification, level EOS_EOnlineSessionPermissionLevel) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_SetMaxPlayers(mod EOS_HSessionModification, max uint32) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_SetJoinInProgressAllowed(mod EOS_HSessionModification, allowed bool) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_SetInvitesAllowed(mod EOS_HSessionModification, allowed bool) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_SetHostAddress(mod EOS_HSessionModification, addr string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_AddAttributeInt64(mod EOS_HSessionModification, key string, val int64, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_AddAttributeDouble(mod EOS_HSessionModification, key string, val float64, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_AddAttributeBool(mod EOS_HSessionModification, key string, val bool, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_AddAttributeString(mod EOS_HSessionModification, key string, val string, advType EOS_ESessionAttributeAdvertisementType) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_RemoveAttribute(mod EOS_HSessionModification, key string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionModification_Release(mod EOS_HSessionModification) {}

// Search

func EOS_Sessions_CreateSessionSearch(handle EOS_HSessions, maxResults uint32) (EOS_HSessionSearch, EOS_EResult) {
	return EOS_HSessionSearch(1), EOS_EResult_Success
}

func EOS_SessionSearch_Find(search EOS_HSessionSearch, localUserId EOS_ProductUserId, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data:       &EOS_SessionSearch_FindCallbackInfo{ResultCode: EOS_EResult_Success},
		})
	}()
}

func EOS_SessionSearch_SetParameterInt64(search EOS_HSessionSearch, key string, val int64, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionSearch_SetParameterDouble(search EOS_HSessionSearch, key string, val float64, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionSearch_SetParameterBool(search EOS_HSessionSearch, key string, val bool, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionSearch_SetParameterString(search EOS_HSessionSearch, key string, val string, op EOS_EComparisonOp) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionSearch_SetSessionId(search EOS_HSessionSearch, sessionId string) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionSearch_SetTargetUserId(search EOS_HSessionSearch, targetUserId EOS_ProductUserId) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionSearch_SetMaxResults(search EOS_HSessionSearch, maxResults uint32) EOS_EResult {
	return EOS_EResult_Success
}
func EOS_SessionSearch_GetSearchResultCount(search EOS_HSessionSearch) uint32 { return 0 }
func EOS_SessionSearch_CopySearchResultByIndex(search EOS_HSessionSearch, index uint32) (EOS_HSessionDetails, EOS_EResult) {
	return 0, EOS_EResult_NotFound
}
func EOS_SessionSearch_Release(search EOS_HSessionSearch) {}

// Session details

func EOS_SessionDetails_CopyInfo(details EOS_HSessionDetails) (*EOS_SessionDetails_Info, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_SessionDetails_GetSessionAttributeCount(details EOS_HSessionDetails) uint32 { return 0 }
func EOS_SessionDetails_CopySessionAttributeByIndex(details EOS_HSessionDetails, index uint32) (*EOS_Sessions_Attribute, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_SessionDetails_CopySessionAttributeByKey(details EOS_HSessionDetails, key string) (*EOS_Sessions_Attribute, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}
func EOS_SessionDetails_Release(details EOS_HSessionDetails) {}

// Active session

func EOS_Sessions_CopyActiveSessionHandle(handle EOS_HSessions, sessionName string) (EOS_HActiveSession, EOS_EResult) {
	return 0, EOS_EResult_NotFound
}
func EOS_ActiveSession_GetRegisteredPlayerCount(active EOS_HActiveSession) uint32 { return 0 }
func EOS_ActiveSession_GetRegisteredPlayerByIndex(active EOS_HActiveSession, index uint32) EOS_ProductUserId {
	return 0
}
func EOS_ActiveSession_Release(active EOS_HActiveSession) {}

// Notifications

func EOS_Sessions_AddNotifySessionInviteReceived(handle EOS_HSessions, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubSessionsNotifyCounter, 1))
}
func EOS_Sessions_RemoveNotifySessionInviteReceived(handle EOS_HSessions, id EOS_NotificationId) {
}
