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

//export goSessionsUpdateCallback
func goSessionsUpdateCallback(resultCode C.int, clientData C.uintptr_t, sessionName *C.char, sessionId *C.char) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Sessions_UpdateSessionCallbackInfo{
			ResultCode:  EOS_EResult(resultCode),
			SessionName: C.GoString(sessionName),
			SessionId:   C.GoString(sessionId),
		},
	})
}

//export goSessionsDestroyCallback
func goSessionsDestroyCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       &EOS_Sessions_DestroySessionCallbackInfo{ResultCode: EOS_EResult(resultCode)},
	})
}

//export goSessionsJoinCallback
func goSessionsJoinCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       &EOS_Sessions_JoinSessionCallbackInfo{ResultCode: EOS_EResult(resultCode)},
	})
}

//export goSessionsStartCallback
func goSessionsStartCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       &EOS_Sessions_StartSessionCallbackInfo{ResultCode: EOS_EResult(resultCode)},
	})
}

//export goSessionsEndCallback
func goSessionsEndCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       &EOS_Sessions_EndSessionCallbackInfo{ResultCode: EOS_EResult(resultCode)},
	})
}

//export goSessionsRegisterPlayersCallback
func goSessionsRegisterPlayersCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       &EOS_Sessions_RegisterPlayersCallbackInfo{ResultCode: EOS_EResult(resultCode)},
	})
}

//export goSessionsUnregisterPlayersCallback
func goSessionsUnregisterPlayersCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       &EOS_Sessions_UnregisterPlayersCallbackInfo{ResultCode: EOS_EResult(resultCode)},
	})
}

//export goSessionSearchFindCallback
func goSessionSearchFindCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       &EOS_SessionSearch_FindCallbackInfo{ResultCode: EOS_EResult(resultCode)},
	})
}

//export goSessionsInviteReceivedCallback
func goSessionsInviteReceivedCallback(clientData C.uintptr_t, localUserId C.uintptr_t, targetUserId C.uintptr_t, inviteId *C.char) {
	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Sessions_SessionInviteReceivedCallbackInfo{
		LocalUserId:  EOS_ProductUserId(localUserId),
		TargetUserId: EOS_ProductUserId(targetUserId),
		InviteId:     C.GoString(inviteId),
	})
}
