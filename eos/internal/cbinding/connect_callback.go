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

//export goConnectLoginCallback
func goConnectLoginCallback(resultCode C.int, clientData C.uintptr_t,
	localUserId C.uintptr_t, continuanceToken C.uintptr_t) {

	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Connect_LoginCallbackInfo{
			ResultCode:       EOS_EResult(resultCode),
			LocalUserId:      EOS_ProductUserId(localUserId),
			ContinuanceToken: EOS_ContinuanceToken(continuanceToken),
		},
	})
}

//export goConnectCreateUserCallback
func goConnectCreateUserCallback(resultCode C.int, clientData C.uintptr_t,
	localUserId C.uintptr_t) {

	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Connect_CreateUserCallbackInfo{
			ResultCode:  EOS_EResult(resultCode),
			LocalUserId: EOS_ProductUserId(localUserId),
		},
	})
}

//export goConnectLinkAccountCallback
func goConnectLinkAccountCallback(resultCode C.int, clientData C.uintptr_t,
	localUserId C.uintptr_t) {

	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Connect_LinkAccountCallbackInfo{
			ResultCode:  EOS_EResult(resultCode),
			LocalUserId: EOS_ProductUserId(localUserId),
		},
	})
}

//export goConnectAuthExpirationCallback
func goConnectAuthExpirationCallback(clientData C.uintptr_t, localUserId C.uintptr_t) {
	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Connect_AuthExpirationCallbackInfo{
		LocalUserId: EOS_ProductUserId(localUserId),
	})
}

//export goConnectLoginStatusChangedCallback
func goConnectLoginStatusChangedCallback(clientData C.uintptr_t,
	localUserId C.uintptr_t, previousStatus C.int, currentStatus C.int) {

	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Connect_LoginStatusChangedCallbackInfo{
		LocalUserId:    EOS_ProductUserId(localUserId),
		PreviousStatus: EOS_ELoginStatus(previousStatus),
		CurrentStatus:  EOS_ELoginStatus(currentStatus),
	})
}
