//go:build eosstub

package cbinding

import (
	"runtime/cgo"
	"sync/atomic"

	"github.com/mydev/go-eos/eos/internal/callback"
)

var stubAuthNotifyCounter uint64

func EOS_Auth_Login(handle EOS_HAuth, opts *EOS_Auth_LoginOptions, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Auth_LoginCallbackInfo{
				ResultCode:        EOS_EResult_Success,
				LocalUserId:       EOS_EpicAccountId(1),
				SelectedAccountId: EOS_EpicAccountId(1),
			},
		})
	}()
}

func EOS_Auth_Logout(handle EOS_HAuth, opts *EOS_Auth_LogoutOptions, clientData uintptr) {
	localUserId := opts.LocalUserId
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Auth_LogoutCallbackInfo{
				ResultCode:  EOS_EResult_Success,
				LocalUserId: localUserId,
			},
		})
	}()
}

func EOS_Auth_DeletePersistentAuth(handle EOS_HAuth, opts *EOS_Auth_DeletePersistentAuthOptions, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Auth_DeletePersistentAuthCallbackInfo{
				ResultCode: EOS_EResult_Success,
			},
		})
	}()
}

func EOS_Auth_GetLoggedInAccountsCount(handle EOS_HAuth) int32 {
	return 0
}

func EOS_Auth_GetLoggedInAccountByIndex(handle EOS_HAuth, index int32) EOS_EpicAccountId {
	return 0
}

func EOS_Auth_CopyUserAuthToken(handle EOS_HAuth, localUserId EOS_EpicAccountId) (*EOS_Auth_Token, EOS_EResult) {
	return nil, EOS_EResult_NotFound
}

func EOS_Auth_CopyIdToken(handle EOS_HAuth, accountId EOS_EpicAccountId) (string, EOS_EResult) {
	return "stub-jwt-token", EOS_EResult_Success
}

func EOS_Auth_AddNotifyLoginStatusChanged(handle EOS_HAuth, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubAuthNotifyCounter, 1))
}

func EOS_Auth_RemoveNotifyLoginStatusChanged(handle EOS_HAuth, id EOS_NotificationId) {}
