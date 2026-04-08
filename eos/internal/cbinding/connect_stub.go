//go:build eosstub

package cbinding

import (
	"runtime/cgo"
	"sync/atomic"

	"github.com/mydev/go-eos/eos/internal/callback"
)

var stubConnectNotifyCounter uint64

func EOS_Connect_Login(handle EOS_HConnect, opts *EOS_Connect_LoginOptions, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Connect_LoginCallbackInfo{
				ResultCode:  EOS_EResult_Success,
				LocalUserId: EOS_ProductUserId(1),
			},
		})
	}()
}

func EOS_Connect_CreateUser(handle EOS_HConnect, opts *EOS_Connect_CreateUserOptions, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Connect_CreateUserCallbackInfo{
				ResultCode:  EOS_EResult_Success,
				LocalUserId: EOS_ProductUserId(1),
			},
		})
	}()
}

func EOS_Connect_LinkAccount(handle EOS_HConnect, opts *EOS_Connect_LinkAccountOptions, clientData uintptr) {
	localUserId := opts.LocalUserId
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Connect_LinkAccountCallbackInfo{
				ResultCode:  EOS_EResult_Success,
				LocalUserId: localUserId,
			},
		})
	}()
}

func EOS_Connect_GetLoggedInUsersCount(handle EOS_HConnect) int32 {
	return 0
}

func EOS_Connect_GetLoggedInUserByIndex(handle EOS_HConnect, index int32) EOS_ProductUserId {
	return 0
}

func EOS_Connect_AddNotifyAuthExpiration(handle EOS_HConnect, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubConnectNotifyCounter, 1))
}

func EOS_Connect_RemoveNotifyAuthExpiration(handle EOS_HConnect, id EOS_NotificationId) {}

func EOS_Connect_AddNotifyLoginStatusChanged(handle EOS_HConnect, clientData uintptr) EOS_NotificationId {
	return EOS_NotificationId(atomic.AddUint64(&stubConnectNotifyCounter, 1))
}

func EOS_Connect_RemoveNotifyLoginStatusChanged(handle EOS_HConnect, id EOS_NotificationId) {}
