package lobby

import (
	"context"
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

// Lobby wraps the EOS Lobby interface for matchmaking and lobby management.
type Lobby struct {
	handle cbinding.EOS_HLobby
	worker *threadworker.Worker
}

// New creates a new Lobby instance from a platform lobby handle.
func New(handle cbinding.EOS_HLobby, worker *threadworker.Worker) *Lobby {
	return &Lobby{handle: handle, worker: worker}
}

// CreateLobbyOptions configures a new lobby creation request.
type CreateLobbyOptions struct {
	LocalUserId     types.ProductUserId
	MaxMembers      uint32
	PermissionLevel PermissionLevel
	AllowInvites    bool
	BucketId        string
}

// LobbyInfo contains metadata about a lobby instance.
type LobbyInfo struct {
	LobbyId          string
	LobbyOwnerUserId types.ProductUserId
	PermissionLevel  PermissionLevel
	AvailableSlots   uint32
	MaxMembers       uint32
	AllowInvites     bool
	BucketId         string
}

// LobbyUpdateInfo is delivered when a lobby's attributes change.
type LobbyUpdateInfo struct {
	LobbyId string
}

// MemberUpdateInfo is delivered when a lobby member's attributes change.
type MemberUpdateInfo struct {
	LobbyId      string
	TargetUserId types.ProductUserId
}

// MemberStatusInfo is delivered when a lobby member's status changes (join, leave, etc.).
type MemberStatusInfo struct {
	LobbyId       string
	TargetUserId  types.ProductUserId
	CurrentStatus MemberStatus
}

// InviteReceivedInfo is delivered when the local user receives a lobby invite.
type InviteReceivedInfo struct {
	InviteId     string
	LocalUserId  types.ProductUserId
	TargetUserId types.ProductUserId
}

// CreateLobby creates a new lobby and returns its ID. Wraps EOS_Lobby_CreateLobby.
func (l *Lobby) CreateLobby(ctx context.Context, opts CreateLobbyOptions) (string, error) {
	oneshot := callback.NewOneShot()
	cUserId := cbinding.EOS_ProductUserId_FromString(string(opts.LocalUserId))

	if err := l.worker.Submit(func() {
		cbinding.EOS_Lobby_CreateLobby(l.handle, &cbinding.EOS_Lobby_CreateLobbyOptions{
			LocalUserId:     cUserId,
			MaxLobbyMembers: opts.MaxMembers,
			PermissionLevel: cbinding.EOS_ELobbyPermissionLevel(opts.PermissionLevel),
			AllowInvites:    opts.AllowInvites,
			BucketId:        opts.BucketId,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return "", err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return "", err
	}
	info := result.Data.(*cbinding.EOS_Lobby_CreateLobbyCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return "", types.NewResult(int(info.ResultCode))
	}
	return info.LobbyId, nil
}

// DestroyLobby destroys an existing lobby. Wraps EOS_Lobby_DestroyLobby.
func (l *Lobby) DestroyLobby(ctx context.Context, localUserId types.ProductUserId, lobbyId string) error {
	oneshot := callback.NewOneShot()
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))

	if err := l.worker.Submit(func() {
		cbinding.EOS_Lobby_DestroyLobby(l.handle, &cbinding.EOS_Lobby_DestroyLobbyOptions{
			LocalUserId: cUserId,
			LobbyId:     lobbyId,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}
	info := result.Data.(*cbinding.EOS_Lobby_DestroyLobbyCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

// JoinLobby joins an existing lobby using the provided details handle. Wraps EOS_Lobby_JoinLobby.
func (l *Lobby) JoinLobby(ctx context.Context, localUserId types.ProductUserId, details *LobbyDetails) error {
	oneshot := callback.NewOneShot()
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))

	if err := l.worker.Submit(func() {
		cbinding.EOS_Lobby_JoinLobby(l.handle, &cbinding.EOS_Lobby_JoinLobbyOptions{
			LobbyDetailsHandle: details.handle,
			LocalUserId:        cUserId,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}
	info := result.Data.(*cbinding.EOS_Lobby_JoinLobbyCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

// LeaveLobby leaves a lobby the local user has joined. Wraps EOS_Lobby_LeaveLobby.
func (l *Lobby) LeaveLobby(ctx context.Context, localUserId types.ProductUserId, lobbyId string) error {
	oneshot := callback.NewOneShot()
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))

	if err := l.worker.Submit(func() {
		cbinding.EOS_Lobby_LeaveLobby(l.handle, &cbinding.EOS_Lobby_LeaveLobbyOptions{
			LocalUserId: cUserId,
			LobbyId:     lobbyId,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}
	info := result.Data.(*cbinding.EOS_Lobby_LeaveLobbyCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

// UpdateLobby applies a modification handle to an existing lobby. Wraps EOS_Lobby_UpdateLobby.
func (l *Lobby) UpdateLobby(ctx context.Context, mod *LobbyModification) error {
	oneshot := callback.NewOneShot()

	if err := l.worker.Submit(func() {
		cbinding.EOS_Lobby_UpdateLobby(l.handle, &cbinding.EOS_Lobby_UpdateLobbyOptions{
			LobbyModificationHandle: mod.handle,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}
	info := result.Data.(*cbinding.EOS_Lobby_UpdateLobbyCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

// UpdateLobbyModification creates a modification handle for a lobby. Wraps EOS_Lobby_UpdateLobbyModification.
func (l *Lobby) UpdateLobbyModification(localUserId types.ProductUserId, lobbyId string) (*LobbyModification, error) {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))
	var mod cbinding.EOS_HLobbyModification
	var result cbinding.EOS_EResult

	if err := l.worker.Submit(func() {
		mod, result = cbinding.EOS_Lobby_UpdateLobbyModification(l.handle, cUserId, lobbyId)
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &LobbyModification{handle: mod, worker: l.worker}, nil
}

// CopyLobbyDetailsHandle retrieves a details handle for a lobby. Wraps EOS_Lobby_CopyLobbyDetailsHandle.
func (l *Lobby) CopyLobbyDetailsHandle(localUserId types.ProductUserId, lobbyId string) (*LobbyDetails, error) {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))
	var details cbinding.EOS_HLobbyDetails
	var result cbinding.EOS_EResult

	if err := l.worker.Submit(func() {
		details, result = cbinding.EOS_Lobby_CopyLobbyDetailsHandle(l.handle, cUserId, lobbyId)
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &LobbyDetails{handle: details, worker: l.worker}, nil
}

// CreateLobbySearch creates a lobby search handle. Wraps EOS_Lobby_CreateLobbySearch.
func (l *Lobby) CreateLobbySearch(maxResults uint32) (*LobbySearch, error) {
	var search cbinding.EOS_HLobbySearch
	var result cbinding.EOS_EResult

	if err := l.worker.Submit(func() {
		search, result = cbinding.EOS_Lobby_CreateLobbySearch(l.handle, maxResults)
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &LobbySearch{handle: search, worker: l.worker}, nil
}

// SendInvite sends a lobby invite to another user. Wraps EOS_Lobby_SendInvite.
func (l *Lobby) SendInvite(ctx context.Context, lobbyId string, localUserId, targetUserId types.ProductUserId) error {
	oneshot := callback.NewOneShot()
	cLocal := cbinding.EOS_ProductUserId_FromString(string(localUserId))
	cTarget := cbinding.EOS_ProductUserId_FromString(string(targetUserId))

	if err := l.worker.Submit(func() {
		cbinding.EOS_Lobby_SendInvite(l.handle, &cbinding.EOS_Lobby_SendInviteOptions{
			LobbyId:      lobbyId,
			LocalUserId:  cLocal,
			TargetUserId: cTarget,
		}, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}
	info := result.Data.(*cbinding.EOS_Lobby_SendInviteCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return types.NewResult(int(info.ResultCode))
	}
	return nil
}

// GetInviteCount returns the number of pending lobby invites for the local user.
func (l *Lobby) GetInviteCount(localUserId types.ProductUserId) uint32 {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))
	var count uint32
	if err := l.worker.Submit(func() {
		count = cbinding.EOS_Lobby_GetInviteCount(l.handle, cUserId)
	}); err != nil {
		return 0
	}
	return count
}

// GetInviteIdByIndex returns the invite ID at the given index for the local user.
func (l *Lobby) GetInviteIdByIndex(localUserId types.ProductUserId, index uint32) (string, error) {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))
	var id string
	var result cbinding.EOS_EResult

	if err := l.worker.Submit(func() {
		id, result = cbinding.EOS_Lobby_GetInviteIdByIndex(l.handle, cUserId, index)
	}); err != nil {
		return "", err
	}
	if result != cbinding.EOS_EResult_Success {
		return "", types.NewResult(int(result))
	}
	return id, nil
}

// AddNotifyLobbyUpdateReceived registers a callback for lobby attribute changes.
func (l *Lobby) AddNotifyLobbyUpdateReceived(fn func(LobbyUpdateInfo)) callback.RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_Lobby_LobbyUpdateReceivedCallbackInfo)
		fn(LobbyUpdateInfo{LobbyId: info.LobbyId})
	})
	handle := cgo.NewHandle(notifyFn)
	var notifyId cbinding.EOS_NotificationId
	if err := l.worker.Submit(func() {
		notifyId = cbinding.EOS_Lobby_AddNotifyLobbyUpdateReceived(l.handle, uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}
	return func() {
		_ = l.worker.Submit(func() { cbinding.EOS_Lobby_RemoveNotifyLobbyUpdateReceived(l.handle, notifyId) })
		handle.Delete()
	}
}

// AddNotifyLobbyMemberUpdateReceived registers a callback for lobby member attribute changes.
func (l *Lobby) AddNotifyLobbyMemberUpdateReceived(fn func(MemberUpdateInfo)) callback.RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_Lobby_LobbyMemberUpdateReceivedCallbackInfo)
		fn(MemberUpdateInfo{
			LobbyId:      info.LobbyId,
			TargetUserId: types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.TargetUserId)),
		})
	})
	handle := cgo.NewHandle(notifyFn)
	var notifyId cbinding.EOS_NotificationId
	if err := l.worker.Submit(func() {
		notifyId = cbinding.EOS_Lobby_AddNotifyLobbyMemberUpdateReceived(l.handle, uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}
	return func() {
		_ = l.worker.Submit(func() { cbinding.EOS_Lobby_RemoveNotifyLobbyMemberUpdateReceived(l.handle, notifyId) })
		handle.Delete()
	}
}

// AddNotifyLobbyMemberStatusReceived registers a callback for lobby member status changes.
func (l *Lobby) AddNotifyLobbyMemberStatusReceived(fn func(MemberStatusInfo)) callback.RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_Lobby_LobbyMemberStatusReceivedCallbackInfo)
		fn(MemberStatusInfo{
			LobbyId:       info.LobbyId,
			TargetUserId:  types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.TargetUserId)),
			CurrentStatus: MemberStatus(info.CurrentStatus),
		})
	})
	handle := cgo.NewHandle(notifyFn)
	var notifyId cbinding.EOS_NotificationId
	if err := l.worker.Submit(func() {
		notifyId = cbinding.EOS_Lobby_AddNotifyLobbyMemberStatusReceived(l.handle, uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}
	return func() {
		_ = l.worker.Submit(func() { cbinding.EOS_Lobby_RemoveNotifyLobbyMemberStatusReceived(l.handle, notifyId) })
		handle.Delete()
	}
}
