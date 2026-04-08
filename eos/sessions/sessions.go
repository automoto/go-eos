package sessions

import (
	"context"
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

type Sessions struct {
	handle cbinding.EOS_HSessions
	worker *threadworker.Worker
}

func New(handle cbinding.EOS_HSessions, worker *threadworker.Worker) *Sessions {
	return &Sessions{handle: handle, worker: worker}
}

type SessionPermissionLevel int

const (
	PermissionPublicAdvertised SessionPermissionLevel = 0
	PermissionJoinViaPresence  SessionPermissionLevel = 1
	PermissionInviteOnly       SessionPermissionLevel = 2
)

type AdvertisementType int

const (
	AdvertisementDontAdvertise AdvertisementType = 0
	AdvertisementAdvertise     AdvertisementType = 1
)

type SessionState int

const (
	StateNoSession  SessionState = 0
	StateCreating   SessionState = 1
	StatePending    SessionState = 2
	StateStarting   SessionState = 3
	StateInProgress SessionState = 4
	StateEnding     SessionState = 5
	StateEnded      SessionState = 6
	StateDestroying SessionState = 7
)

type ComparisonOp = int

type CreateSessionOptions struct {
	SessionName string
	BucketId    string
	MaxPlayers  uint32
	LocalUserId types.ProductUserId
}

type SessionInfo struct {
	SessionId               string
	HostAddress             string
	NumOpenPublicConnections uint32
	OwnerUserId             types.ProductUserId
	BucketId                string
	NumPublicConnections    uint32
	AllowJoinInProgress     bool
	PermissionLevel         SessionPermissionLevel
	InvitesAllowed          bool
}

type Attribute struct {
	Key               string
	Value             any
	AdvertisementType AdvertisementType
}

type InviteReceivedInfo struct {
	LocalUserId  types.ProductUserId
	TargetUserId types.ProductUserId
	InviteId     string
}

// Core lifecycle

func (s *Sessions) CreateSessionModification(opts CreateSessionOptions) (*SessionModification, error) {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(opts.LocalUserId))
	var mod cbinding.EOS_HSessionModification
	var result cbinding.EOS_EResult

	if err := s.worker.Submit(func() {
		mod, result = cbinding.EOS_Sessions_CreateSessionModification(s.handle, &cbinding.EOS_Sessions_CreateSessionModificationOptions{
			SessionName: opts.SessionName,
			BucketId:    opts.BucketId,
			MaxPlayers:  opts.MaxPlayers,
			LocalUserId: cUserId,
		})
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &SessionModification{handle: mod, worker: s.worker}, nil
}

func (s *Sessions) UpdateSession(ctx context.Context, mod *SessionModification) (*SessionInfo, error) {
	oneshot := callback.NewOneShot()

	if err := s.worker.Submit(func() {
		cbinding.EOS_Sessions_UpdateSession(s.handle, mod.handle, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return nil, err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return nil, err
	}
	info := result.Data.(*cbinding.EOS_Sessions_UpdateSessionCallbackInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(info.ResultCode))
	}
	return &SessionInfo{SessionId: info.SessionId}, nil
}

func (s *Sessions) DestroySession(ctx context.Context, sessionName string) error {
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_DestroySession(s.handle, sessionName, clientData)
	})
}

func (s *Sessions) JoinSession(ctx context.Context, sessionName string, details cbinding.EOS_HSessionDetails, localUserId types.ProductUserId) error {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_JoinSession(s.handle, sessionName, details, cUserId, clientData)
	})
}

func (s *Sessions) StartSession(ctx context.Context, sessionName string) error {
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_StartSession(s.handle, sessionName, clientData)
	})
}

func (s *Sessions) EndSession(ctx context.Context, sessionName string) error {
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_EndSession(s.handle, sessionName, clientData)
	})
}

func (s *Sessions) RegisterPlayers(ctx context.Context, sessionName string, players []types.ProductUserId) error {
	cIds := make([]cbinding.EOS_ProductUserId, len(players))
	for i, p := range players {
		cIds[i] = cbinding.EOS_ProductUserId_FromString(string(p))
	}
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_RegisterPlayers(s.handle, sessionName, cIds, clientData)
	})
}

func (s *Sessions) UnregisterPlayers(ctx context.Context, sessionName string, players []types.ProductUserId) error {
	cIds := make([]cbinding.EOS_ProductUserId, len(players))
	for i, p := range players {
		cIds[i] = cbinding.EOS_ProductUserId_FromString(string(p))
	}
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_UnregisterPlayers(s.handle, sessionName, cIds, clientData)
	})
}

// Search

func (s *Sessions) CreateSessionSearch(maxResults uint32) (*SessionSearch, error) {
	var search cbinding.EOS_HSessionSearch
	var result cbinding.EOS_EResult

	if err := s.worker.Submit(func() {
		search, result = cbinding.EOS_Sessions_CreateSessionSearch(s.handle, maxResults)
	}); err != nil {
		return nil, err
	}
	if result != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(result))
	}
	return &SessionSearch{handle: search, worker: s.worker}, nil
}

// Notifications

func (s *Sessions) AddNotifySessionInviteReceived(fn func(InviteReceivedInfo)) callback.RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_Sessions_SessionInviteReceivedCallbackInfo)
		fn(InviteReceivedInfo{
			LocalUserId:  types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId)),
			TargetUserId: types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.TargetUserId)),
			InviteId:     info.InviteId,
		})
	})
	handle := cgo.NewHandle(notifyFn)
	var notifyId cbinding.EOS_NotificationId
	if err := s.worker.Submit(func() {
		notifyId = cbinding.EOS_Sessions_AddNotifySessionInviteReceived(s.handle, uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}
	return func() {
		_ = s.worker.Submit(func() { cbinding.EOS_Sessions_RemoveNotifySessionInviteReceived(s.handle, notifyId) })
		handle.Delete()
	}
}

// Helper for simple async operations that just return ResultCode
func (s *Sessions) simpleAsync(ctx context.Context, callFn func(clientData uintptr)) error {
	oneshot := callback.NewOneShot()

	if err := s.worker.Submit(func() {
		callFn(oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return err
	}

	if result.ResultCode != int(cbinding.EOS_EResult_Success) {
		return types.NewResult(result.ResultCode)
	}
	return nil
}
