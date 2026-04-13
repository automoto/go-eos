package connect

import (
	"context"
	"fmt"
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

// ContinuanceToken is an opaque token used to continue an interrupted Connect login flow.
type ContinuanceToken = cbinding.EOS_ContinuanceToken

// Connect wraps the EOS Connect interface for product-user authentication.
type Connect struct {
	handle cbinding.EOS_HConnect
	worker *threadworker.Worker
}

// New creates a new Connect instance from the given platform handle and worker.
func New(handle cbinding.EOS_HConnect, worker *threadworker.Worker) *Connect {
	return &Connect{handle: handle, worker: worker}
}

// LoginOptions configures a Connect login request.
type LoginOptions struct {
	CredentialType types.ExternalCredentialType
	Token          string
	DisplayName    string
}

// LoginResult holds the outcome of a Connect login attempt.
type LoginResult struct {
	LocalUserId      types.ProductUserId
	ContinuanceToken ContinuanceToken
}

// LinkAccountOptions configures a request to link an external account to a product user.
type LinkAccountOptions struct {
	LocalUserId      types.ProductUserId
	ContinuanceToken ContinuanceToken
}

// AuthExpirationInfo is delivered when a product user's auth token is about to expire.
type AuthExpirationInfo struct {
	LocalUserId types.ProductUserId
}

// LoginStatusChangedInfo is delivered when a product user's login status changes.
type LoginStatusChangedInfo struct {
	LocalUserId    types.ProductUserId
	PreviousStatus types.LoginStatus
	CurrentStatus  types.LoginStatus
}

// Login authenticates a product user via the EOS Connect interface. See EOS_Connect_Login.
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

	var localStr string
	if err := c.worker.Submit(func() {
		localStr = string(cbinding.EOS_ProductUserId_ToString(info.LocalUserId))
	}); err != nil {
		return nil, fmt.Errorf("id conversion: %w", err)
	}
	return &LoginResult{
		LocalUserId:      types.ProductUserId(localStr),
		ContinuanceToken: info.ContinuanceToken,
	}, nil
}

// CreateUser creates a new product user from a continuance token. See EOS_Connect_CreateUser.
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

	var userStr string
	if err := c.worker.Submit(func() {
		userStr = string(cbinding.EOS_ProductUserId_ToString(info.LocalUserId))
	}); err != nil {
		return nil, fmt.Errorf("id conversion: %w", err)
	}
	userId := types.ProductUserId(userStr)
	return &userId, nil
}

// LinkAccount links an external account to an existing product user. See EOS_Connect_LinkAccount.
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

// GetLoggedInUsersCount returns the number of currently logged-in product users.
func (c *Connect) GetLoggedInUsersCount() int {
	var count int32
	if err := c.worker.Submit(func() {
		count = cbinding.EOS_Connect_GetLoggedInUsersCount(c.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

// GetLoggedInUserByIndex returns the ProductUserId at the given index.
func (c *Connect) GetLoggedInUserByIndex(index int) types.ProductUserId {
	var result string
	if err := c.worker.Submit(func() {
		id := cbinding.EOS_Connect_GetLoggedInUserByIndex(c.handle, int32(index))
		result = string(cbinding.EOS_ProductUserId_ToString(id))
	}); err != nil {
		return ""
	}
	return types.ProductUserId(result)
}

// AddNotifyAuthExpiration registers a callback for auth token expiration warnings. See EOS_Connect_AddNotifyAuthExpiration.
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

// AddNotifyLoginStatusChanged registers a callback for login status changes. See EOS_Connect_AddNotifyLoginStatusChanged.
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

// CreateDeviceId creates an anonymous per-device pseudo-account on the local
// device. After a successful (or DuplicateNotAllowed) call, the caller can
// authenticate via Connect.Login with ExternalCredentialDeviceIDAccessToken
// without requiring the player to have any external account.
//
// deviceModel is a free-form description (e.g. "PC Windows", "iPhone 15").
// Maximum 64 UTF-8 characters; longer strings are silently truncated by the SDK.
//
// Returns nil on success or if a device ID already exists (EOS_DuplicateNotAllowed).
func (c *Connect) CreateDeviceId(ctx context.Context, deviceModel string) error {
	oneshot := callback.NewOneShot()

	if err := c.worker.Submit(func() {
		cbinding.EOS_Connect_CreateDeviceId(c.handle, &cbinding.EOS_Connect_CreateDeviceIdOptions{
			DeviceModel: deviceModel,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}

	info := result.Data.(*cbinding.EOS_Connect_CreateDeviceIdCallbackInfo)
	if info.ResultCode == cbinding.EOS_EResult_Success ||
		int(info.ResultCode) == types.CodeDuplicateNotAllowed {
		return nil
	}
	return types.NewResult(int(info.ResultCode))
}

// DeleteDeviceId removes the local device ID credentials.
func (c *Connect) DeleteDeviceId(ctx context.Context) error {
	oneshot := callback.NewOneShot()

	if err := c.worker.Submit(func() {
		cbinding.EOS_Connect_DeleteDeviceId(c.handle, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}

	info := result.Data.(*cbinding.EOS_Connect_DeleteDeviceIdCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}
