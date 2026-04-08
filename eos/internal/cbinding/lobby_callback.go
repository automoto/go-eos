//go:build !eosstub

package cbinding

/*
#include <stdint.h>
*/
import "C"

import (
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
)

// Lifecycle callbacks — complete oneshot via handle

//export goLobbyCreateCallback
func goLobbyCreateCallback(resultCode C.int, clientData C.uintptr_t, lobbyId *C.char) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Lobby_CreateLobbyCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
			LobbyId:    C.GoString(lobbyId),
		},
	})
}

//export goLobbyDestroyCallback
func goLobbyDestroyCallback(resultCode C.int, clientData C.uintptr_t, lobbyId *C.char) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Lobby_DestroyLobbyCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
			LobbyId:    C.GoString(lobbyId),
		},
	})
}

//export goLobbyJoinCallback
func goLobbyJoinCallback(resultCode C.int, clientData C.uintptr_t, lobbyId *C.char) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Lobby_JoinLobbyCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
			LobbyId:    C.GoString(lobbyId),
		},
	})
}

//export goLobbyLeaveCallback
func goLobbyLeaveCallback(resultCode C.int, clientData C.uintptr_t, lobbyId *C.char) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Lobby_LeaveLobbyCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
			LobbyId:    C.GoString(lobbyId),
		},
	})
}

//export goLobbyUpdateCallback
func goLobbyUpdateCallback(resultCode C.int, clientData C.uintptr_t, lobbyId *C.char) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Lobby_UpdateLobbyCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
			LobbyId:    C.GoString(lobbyId),
		},
	})
}

//export goLobbySendInviteCallback
func goLobbySendInviteCallback(resultCode C.int, clientData C.uintptr_t, lobbyId *C.char) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Lobby_SendInviteCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
			LobbyId:    C.GoString(lobbyId),
		},
	})
}

//export goLobbyQueryInvitesCallback
func goLobbyQueryInvitesCallback(resultCode C.int, clientData C.uintptr_t, localUserId C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Lobby_QueryInvitesCallbackInfo{
			ResultCode:  EOS_EResult(resultCode),
			LocalUserId: EOS_ProductUserId(localUserId),
		},
	})
}

//export goLobbySearchFindCallback
func goLobbySearchFindCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_LobbySearch_FindCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
		},
	})
}

// Notification callbacks — dispatch via NotifyFunc in handle

//export goLobbyUpdateReceivedCallback
func goLobbyUpdateReceivedCallback(clientData C.uintptr_t, lobbyId *C.char) {
	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Lobby_LobbyUpdateReceivedCallbackInfo{
		LobbyId: C.GoString(lobbyId),
	})
}

//export goLobbyMemberUpdateReceivedCallback
func goLobbyMemberUpdateReceivedCallback(clientData C.uintptr_t, lobbyId *C.char, targetUserId C.uintptr_t) {
	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Lobby_LobbyMemberUpdateReceivedCallbackInfo{
		LobbyId:      C.GoString(lobbyId),
		TargetUserId: EOS_ProductUserId(targetUserId),
	})
}

//export goLobbyMemberStatusReceivedCallback
func goLobbyMemberStatusReceivedCallback(clientData C.uintptr_t, lobbyId *C.char, targetUserId C.uintptr_t, currentStatus C.int) {
	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Lobby_LobbyMemberStatusReceivedCallbackInfo{
		LobbyId:       C.GoString(lobbyId),
		TargetUserId:  EOS_ProductUserId(targetUserId),
		CurrentStatus: EOS_ELobbyMemberStatus(currentStatus),
	})
}

//export goLobbyInviteReceivedCallback
func goLobbyInviteReceivedCallback(clientData C.uintptr_t, inviteId *C.char, localUserId C.uintptr_t, targetUserId C.uintptr_t) {
	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Lobby_LobbyInviteReceivedCallbackInfo{
		InviteId:     C.GoString(inviteId),
		LocalUserId:  EOS_ProductUserId(localUserId),
		TargetUserId: EOS_ProductUserId(targetUserId),
	})
}
