package auth

import (
	"context"
	"fmt"
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

// Auth wraps the EOS Auth interface for Epic Account authentication.
type Auth struct {
	handle cbinding.EOS_HAuth
	worker *threadworker.Worker
}

// New creates a new Auth instance from the given platform handle and worker.
func New(handle cbinding.EOS_HAuth, worker *threadworker.Worker) *Auth {
	return &Auth{handle: handle, worker: worker}
}

// LoginOptions configures an Auth login request.
type LoginOptions struct {
	CredentialType types.LoginCredentialType
	ID             string
	Token          string
	ScopeFlags     types.AuthScopeFlags
	ExternalType   types.ExternalCredentialType
}

// LoginResult holds the outcome of a successful Auth login.
type LoginResult struct {
	LocalUserId       types.EpicAccountId
	SelectedAccountId types.EpicAccountId
	PinGrantInfo      *PinGrantInfo
}

// PinGrantInfo contains device-code flow information for out-of-band user verification.
type PinGrantInfo struct {
	UserCode                string
	VerificationURI         string
	ExpiresIn               int32
	VerificationURIComplete string
}

// Token represents an EOS Auth user access token and its metadata.
type Token struct {
	App              string
	ClientId         string
	AccountId        types.EpicAccountId
	AccessToken      string
	ExpiresIn        float64
	ExpiresAt        string
	AuthType         int32
	RefreshToken     string
	RefreshExpiresIn float64
	RefreshExpiresAt string
}

// LoginStatusChangedInfo is delivered when an Epic Account's login status changes.
type LoginStatusChangedInfo struct {
	LocalUserId   types.EpicAccountId
	PrevStatus    types.LoginStatus
	CurrentStatus types.LoginStatus
}

// Login authenticates via the EOS Auth interface. See EOS_Auth_Login.
func (a *Auth) Login(ctx context.Context, opts LoginOptions) (*LoginResult, error) {
	oneshot := callback.NewOneShot()

	if err := a.worker.Submit(func() {
		cbinding.EOS_Auth_Login(a.handle, &cbinding.EOS_Auth_LoginOptions{
			CredentialType: cbinding.EOS_ELoginCredentialType(opts.CredentialType),
			ID:             opts.ID,
			Token:          opts.Token,
			ScopeFlags:     cbinding.EOS_EAuthScopeFlags(opts.ScopeFlags),
			ExternalType:   cbinding.EOS_EExternalCredentialType(opts.ExternalType),
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return nil, err
	}
	result, err := oneshot.Wait(ctx)
	if err != nil {
		return nil, err
	}

	info := result.Data.(*cbinding.EOS_Auth_LoginCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(info.ResultCode))
	}

	var localStr, selectedStr string
	if err := a.worker.Submit(func() {
		localStr = string(cbinding.EOS_EpicAccountId_ToString(info.LocalUserId))
		selectedStr = string(cbinding.EOS_EpicAccountId_ToString(info.SelectedAccountId))
	}); err != nil {
		return nil, fmt.Errorf("id conversion: %w", err)
	}
	lr := &LoginResult{
		LocalUserId:       types.EpicAccountId(localStr),
		SelectedAccountId: types.EpicAccountId(selectedStr),
	}
	if info.PinGrantInfo != nil {
		lr.PinGrantInfo = &PinGrantInfo{
			UserCode:                info.PinGrantInfo.UserCode,
			VerificationURI:         info.PinGrantInfo.VerificationURI,
			ExpiresIn:               info.PinGrantInfo.ExpiresIn,
			VerificationURIComplete: info.PinGrantInfo.VerificationURIComplete,
		}
	}
	return lr, nil
}

// Logout signs out the given Epic Account. See EOS_Auth_Logout.
func (a *Auth) Logout(ctx context.Context, localUserId types.EpicAccountId) error {
	oneshot := callback.NewOneShot()
	cId := cbinding.EOS_EpicAccountId_FromString(string(localUserId))

	if err := a.worker.Submit(func() {
		cbinding.EOS_Auth_Logout(a.handle, &cbinding.EOS_Auth_LogoutOptions{
			LocalUserId: cId,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}

	info := result.Data.(*cbinding.EOS_Auth_LogoutCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

// DeletePersistentAuth removes locally stored persistent auth credentials. See EOS_Auth_DeletePersistentAuth.
func (a *Auth) DeletePersistentAuth(ctx context.Context) error {
	oneshot := callback.NewOneShot()

	if err := a.worker.Submit(func() {
		cbinding.EOS_Auth_DeletePersistentAuth(a.handle, &cbinding.EOS_Auth_DeletePersistentAuthOptions{}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}

	info := result.Data.(*cbinding.EOS_Auth_DeletePersistentAuthCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

// GetLoggedInAccountsCount returns the number of currently logged-in accounts. See EOS_Auth_GetLoggedInAccountsCount.
func (a *Auth) GetLoggedInAccountsCount() int {
	var count int32
	if err := a.worker.Submit(func() {
		count = cbinding.EOS_Auth_GetLoggedInAccountsCount(a.handle)
	}); err != nil {
		return 0
	}
	return int(count)
}

// GetLoggedInAccountByIndex returns the EpicAccountId at the given index. See EOS_Auth_GetLoggedInAccountByIndex.
func (a *Auth) GetLoggedInAccountByIndex(index int) types.EpicAccountId {
	var result string
	if err := a.worker.Submit(func() {
		id := cbinding.EOS_Auth_GetLoggedInAccountByIndex(a.handle, int32(index))
		result = string(cbinding.EOS_EpicAccountId_ToString(id))
	}); err != nil {
		return ""
	}
	return types.EpicAccountId(result)
}

// CopyIdToken returns the JWT ID token for the given account. See EOS_Auth_CopyIdToken.
func (a *Auth) CopyIdToken(accountId types.EpicAccountId) (string, error) {
	cId := cbinding.EOS_EpicAccountId_FromString(string(accountId))
	var jwt string
	var result cbinding.EOS_EResult

	if err := a.worker.Submit(func() {
		jwt, result = cbinding.EOS_Auth_CopyIdToken(a.handle, cId)
	}); err != nil {
		return "", err
	}
	if result != cbinding.EOS_EResult_Success {
		return "", types.NewResult(int(result))
	}
	return jwt, nil
}

// CopyUserAuthToken returns the access token for the given account. See EOS_Auth_CopyUserAuthToken.
func (a *Auth) CopyUserAuthToken(localUserId types.EpicAccountId) (*Token, error) {
	cId := cbinding.EOS_EpicAccountId_FromString(string(localUserId))
	var token *cbinding.EOS_Auth_Token
	var result cbinding.EOS_EResult

	var accountIdStr string
	if err := a.worker.Submit(func() {
		token, result = cbinding.EOS_Auth_CopyUserAuthToken(a.handle, cId)
		if result == cbinding.EOS_EResult_Success && token != nil {
			accountIdStr = string(cbinding.EOS_EpicAccountId_ToString(token.AccountId))
		}
	}); err != nil {
		return nil, err
	}

	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &Token{
		App:              token.App,
		ClientId:         token.ClientId,
		AccountId:        types.EpicAccountId(accountIdStr),
		AccessToken:      token.AccessToken,
		ExpiresIn:        token.ExpiresIn,
		ExpiresAt:        token.ExpiresAt,
		AuthType:         token.AuthType,
		RefreshToken:     token.RefreshToken,
		RefreshExpiresIn: token.RefreshExpiresIn,
		RefreshExpiresAt: token.RefreshExpiresAt,
	}, nil
}

// AddNotifyLoginStatusChanged registers a callback for login status changes. See EOS_Auth_AddNotifyLoginStatusChanged.
func (a *Auth) AddNotifyLoginStatusChanged(fn func(LoginStatusChangedInfo)) callback.RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_Auth_LoginStatusChangedCallbackInfo)
		fn(LoginStatusChangedInfo{
			LocalUserId:   types.EpicAccountId(cbinding.EOS_EpicAccountId_ToString(info.LocalUserId)),
			PrevStatus:    types.LoginStatus(info.PrevStatus),
			CurrentStatus: types.LoginStatus(info.CurrentStatus),
		})
	})
	handle := cgo.NewHandle(notifyFn)

	var notifyId cbinding.EOS_NotificationId
	if err := a.worker.Submit(func() {
		notifyId = cbinding.EOS_Auth_AddNotifyLoginStatusChanged(a.handle, uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}

	return func() {
		_ = a.worker.Submit(func() {
			cbinding.EOS_Auth_RemoveNotifyLoginStatusChanged(a.handle, notifyId)
		})
		handle.Delete()
	}
}
