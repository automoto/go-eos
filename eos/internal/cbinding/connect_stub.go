//go:build eosstub

package cbinding

import (
	"runtime/cgo"
	"sync/atomic"

	"github.com/mydev/go-eos/eos/internal/callback"
)

var stubConnectNotifyCounter uint64

// StubCreateDeviceIdResultCode lets tests force the result returned by the
// stub EOS_Connect_CreateDeviceId. Default is Success. Tests that mutate
// this MUST reset it to Success in their cleanup.
var StubCreateDeviceIdResultCode = EOS_EResult_Success

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

func EOS_Connect_CreateDeviceId(handle EOS_HConnect, opts *EOS_Connect_CreateDeviceIdOptions, clientData uintptr) {
	code := StubCreateDeviceIdResultCode
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(code),
			Data: &EOS_Connect_CreateDeviceIdCallbackInfo{
				ResultCode: code,
			},
		})
	}()
}

func EOS_Connect_DeleteDeviceId(handle EOS_HConnect, clientData uintptr) {
	go func() {
		h := cgo.Handle(clientData)
		callback.CompleteByHandle(h, callback.OneShotResult{
			ResultCode: int(EOS_EResult_Success),
			Data: &EOS_Connect_DeleteDeviceIdCallbackInfo{
				ResultCode: EOS_EResult_Success,
			},
		})
	}()
}
