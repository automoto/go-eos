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

//export goAuthLoginCallback
func goAuthLoginCallback(resultCode C.int, clientData C.uintptr_t,
	localUserId C.uintptr_t, selectedAccountId C.uintptr_t,
	continuanceToken C.uintptr_t, hasPinGrant C.int,
	pinUserCode *C.char, pinVerificationURI *C.char,
	pinExpiresIn C.int, pinVerificationURIComplete *C.char) {

	info := &EOS_Auth_LoginCallbackInfo{
		ResultCode:        EOS_EResult(resultCode),
		LocalUserId:       EOS_EpicAccountId(localUserId),
		SelectedAccountId: EOS_EpicAccountId(selectedAccountId),
		ContinuanceToken:  EOS_ContinuanceToken(continuanceToken),
	}

	if hasPinGrant != 0 {
		info.PinGrantInfo = &EOS_Auth_PinGrantInfo{
			UserCode:                C.GoString(pinUserCode),
			VerificationURI:         C.GoString(pinVerificationURI),
			ExpiresIn:               int32(pinExpiresIn),
			VerificationURIComplete: C.GoString(pinVerificationURIComplete),
		}
	}

	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data:       info,
	})
}

//export goAuthLogoutCallback
func goAuthLogoutCallback(resultCode C.int, clientData C.uintptr_t, localUserId C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Auth_LogoutCallbackInfo{
			ResultCode:  EOS_EResult(resultCode),
			LocalUserId: EOS_EpicAccountId(localUserId),
		},
	})
}

//export goAuthDeletePersistentAuthCallback
func goAuthDeletePersistentAuthCallback(resultCode C.int, clientData C.uintptr_t) {
	h := cgo.Handle(clientData)
	callback.CompleteByHandle(h, callback.OneShotResult{
		ResultCode: int(resultCode),
		Data: &EOS_Auth_DeletePersistentAuthCallbackInfo{
			ResultCode: EOS_EResult(resultCode),
		},
	})
}

//export goAuthLoginStatusChangedCallback
func goAuthLoginStatusChangedCallback(clientData C.uintptr_t,
	localUserId C.uintptr_t, prevStatus C.int, currentStatus C.int) {

	h := cgo.Handle(clientData)
	fn := h.Value().(callback.NotifyFunc)
	fn(&EOS_Auth_LoginStatusChangedCallbackInfo{
		LocalUserId:   EOS_EpicAccountId(localUserId),
		PrevStatus:    EOS_ELoginStatus(prevStatus),
		CurrentStatus: EOS_ELoginStatus(currentStatus),
	})
}
