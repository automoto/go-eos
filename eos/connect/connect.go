package connect

import (
	"context"
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

type ContinuanceToken = cbinding.EOS_ContinuanceToken

type Connect struct {
	handle cbinding.EOS_HConnect
	worker *threadworker.Worker
}

func New(handle cbinding.EOS_HConnect, worker *threadworker.Worker) *Connect {
	return &Connect{handle: handle, worker: worker}
}

type LoginOptions struct {
	CredentialType types.ExternalCredentialType
	Token          string
	DisplayName    string
}

type LoginResult struct {
	LocalUserId      types.ProductUserId
	ContinuanceToken ContinuanceToken
}

type LinkAccountOptions struct {
	LocalUserId      types.ProductUserId
	ContinuanceToken ContinuanceToken
}

type AuthExpirationInfo struct {
	LocalUserId types.ProductUserId
}

type LoginStatusChangedInfo struct {
	LocalUserId    types.ProductUserId
	PreviousStatus types.LoginStatus
	CurrentStatus  types.LoginStatus
}

func (c *Connect) Login(ctx context.Context, opts LoginOptions) (*LoginResult, error) {
	oneshot := callback.NewOneShot()

	if err := c.worker.Submit(func() {
		cbinding.EOS_Connect_Login(c.handle, &cbinding.EOS_Connect_LoginOptions{
			CredentialType: cbinding.EOS_EExternalCredentialType(opts.CredentialType),
			Token:          opts.Token,
			DisplayName:    opts.DisplayName,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return nil, err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return nil, err
	}

	info := result.Data.(*cbinding.EOS_Connect_LoginCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return &LoginResult{
			ContinuanceToken: info.ContinuanceToken,
		}, types.NewResult(int(info.ResultCode))
	}

	return &LoginResult{
		LocalUserId:      types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId)),
		ContinuanceToken: info.ContinuanceToken,
	}, nil
}

func (c *Connect) CreateUser(ctx context.Context, continuanceToken ContinuanceToken) (*types.ProductUserId, error) {
	oneshot := callback.NewOneShot()

	if err := c.worker.Submit(func() {
		cbinding.EOS_Connect_CreateUser(c.handle, &cbinding.EOS_Connect_CreateUserOptions{
			ContinuanceToken: continuanceToken,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return nil, err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return nil, err
	}

	info := result.Data.(*cbinding.EOS_Connect_CreateUserCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(info.ResultCode))
	}

	userId := types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId))
	return &userId, nil
}

func (c *Connect) LinkAccount(ctx context.Context, opts LinkAccountOptions) error {
	oneshot := callback.NewOneShot()
	cUserId := cbinding.EOS_ProductUserId_FromString(string(opts.LocalUserId))

	if err := c.worker.Submit(func() {
		cbinding.EOS_Connect_LinkAccount(c.handle, &cbinding.EOS_Connect_LinkAccountOptions{
			LocalUserId:      cUserId,
			ContinuanceToken: opts.ContinuanceToken,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}

	info := result.Data.(*cbinding.EOS_Connect_LinkAccountCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

func (c *Connect) GetLoggedInUsersCount() int {
	var count int32
	if err := c.worker.Submit(func() {
		count = cbinding.EOS_Connect_GetLoggedInUsersCount(c.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

func (c *Connect) GetLoggedInUserByIndex(index int) types.ProductUserId {
	var id cbinding.EOS_ProductUserId
	if err := c.worker.Submit(func() {
		id = cbinding.EOS_Connect_GetLoggedInUserByIndex(c.handle, int32(index))
	}); err != nil {
		return ""
	}
	return types.ProductUserId(cbinding.EOS_ProductUserId_ToString(id))
}

func (c *Connect) AddNotifyAuthExpiration(fn func(AuthExpirationInfo)) callback.RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_Connect_AuthExpirationCallbackInfo)
		fn(AuthExpirationInfo{
			LocalUserId: types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId)),
		})
	})
	handle := cgo.NewHandle(notifyFn)

	var notifyId cbinding.EOS_NotificationId
	if err := c.worker.Submit(func() {
		notifyId = cbinding.EOS_Connect_AddNotifyAuthExpiration(c.handle, uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}

	return func() {
		_ = c.worker.Submit(func() {
			cbinding.EOS_Connect_RemoveNotifyAuthExpiration(c.handle, notifyId)
		})
		handle.Delete()
	}
}

func (c *Connect) AddNotifyLoginStatusChanged(fn func(LoginStatusChangedInfo)) callback.RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_Connect_LoginStatusChangedCallbackInfo)
		fn(LoginStatusChangedInfo{
			LocalUserId:    types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId)),
			PreviousStatus: types.LoginStatus(info.PreviousStatus),
			CurrentStatus:  types.LoginStatus(info.CurrentStatus),
		})
	})
	handle := cgo.NewHandle(notifyFn)

	var notifyId cbinding.EOS_NotificationId
	if err := c.worker.Submit(func() {
		notifyId = cbinding.EOS_Connect_AddNotifyLoginStatusChanged(c.handle, uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}

	return func() {
		_ = c.worker.Submit(func() {
			cbinding.EOS_Connect_RemoveNotifyLoginStatusChanged(c.handle, notifyId)
		})
		handle.Delete()
	}
}
