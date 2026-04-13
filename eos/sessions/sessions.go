package sessions

import (
	"context"
	"runtime/cgo"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/types"
)

// Sessions wraps the EOS Sessions interface for creating, searching, and managing game sessions.
type Sessions struct {
	handle cbinding.EOS_HSessions
	worker *threadworker.Worker
}

// New constructs a Sessions from a raw cbinding handle and the platform's worker.
func New(handle cbinding.EOS_HSessions, worker *threadworker.Worker) *Sessions {
	return &Sessions{handle: handle, worker: worker}
}

// SessionPermissionLevel mirrors EOS_EOnlineSessionPermissionLevel.
type SessionPermissionLevel int

const (
	// PermissionPublicAdvertised allows anyone to find and join the session.
	PermissionPublicAdvertised SessionPermissionLevel = 0
	// PermissionJoinViaPresence restricts joining to presence-based discovery.
	PermissionJoinViaPresence SessionPermissionLevel = 1
	// PermissionInviteOnly restricts joining to invited players only.
	PermissionInviteOnly SessionPermissionLevel = 2
)

// AdvertisementType mirrors EOS_ESessionAttributeAdvertisementType.
type AdvertisementType int

const (
	// AdvertisementDontAdvertise keeps the attribute private to the session.
	AdvertisementDontAdvertise AdvertisementType = 0
	// AdvertisementAdvertise makes the attribute visible in search results.
	AdvertisementAdvertise AdvertisementType = 1
)

// SessionState mirrors EOS_EOnlineSessionState.
type SessionState int

const (
	// StateNoSession indicates no session exists.
	StateNoSession SessionState = 0
	// StateCreating indicates the session is being created.
	StateCreating SessionState = 1
	// StatePending indicates the session is pending.
	StatePending SessionState = 2
	// StateStarting indicates the session is starting.
	StateStarting SessionState = 3
	// StateInProgress indicates the session is actively running.
	StateInProgress SessionState = 4
	// StateEnding indicates the session is ending.
	StateEnding SessionState = 5
	// StateEnded indicates the session has ended.
	StateEnded SessionState = 6
	// StateDestroying indicates the session is being destroyed.
	StateDestroying SessionState = 7
)

// ComparisonOp mirrors EOS_EComparisonOp for session search parameter filtering.
type ComparisonOp = int

// CreateSessionOptions holds the parameters for creating a new session.
type CreateSessionOptions struct {
	SessionName string
	BucketId    string
	MaxPlayers  uint32
	LocalUserId types.ProductUserId
}

// SessionInfo contains the metadata for a session, as returned by EOS_SessionDetails_CopyInfo.
type SessionInfo struct {
	SessionId                string
	HostAddress              string
	NumOpenPublicConnections uint32
	OwnerUserId              types.ProductUserId
	BucketId                 string
	NumPublicConnections     uint32
	AllowJoinInProgress      bool
	PermissionLevel          SessionPermissionLevel
	InvitesAllowed           bool
}

// Attribute represents a key-value pair attached to a session.
type Attribute struct {
	Key               string
	Value             any
	AdvertisementType AdvertisementType
}

// InviteReceivedInfo is the payload delivered by AddNotifySessionInviteReceived.
type InviteReceivedInfo struct {
	LocalUserId  types.ProductUserId
	TargetUserId types.ProductUserId
	InviteId     string
}

// CreateSessionModification creates a new session modification handle. See EOS_Sessions_CreateSessionModification.
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

// UpdateSession applies a session modification and returns the updated session info. See EOS_Sessions_UpdateSession.
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

// DestroySession destroys the named session. See EOS_Sessions_DestroySession.
func (s *Sessions) DestroySession(ctx context.Context, sessionName string) error {
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_DestroySession(s.handle, sessionName, clientData)
	})
}

// JoinSession joins an existing session using the provided session details. See EOS_Sessions_JoinSession.
func (s *Sessions) JoinSession(ctx context.Context, sessionName string, details cbinding.EOS_HSessionDetails, localUserId types.ProductUserId) error {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_JoinSession(s.handle, sessionName, details, cUserId, clientData)
	})
}

// StartSession transitions the named session to the in-progress state. See EOS_Sessions_StartSession.
func (s *Sessions) StartSession(ctx context.Context, sessionName string) error {
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_StartSession(s.handle, sessionName, clientData)
	})
}

// EndSession transitions the named session to the ended state. See EOS_Sessions_EndSession.
func (s *Sessions) EndSession(ctx context.Context, sessionName string) error {
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_EndSession(s.handle, sessionName, clientData)
	})
}

// RegisterPlayers registers the given players with the named session. See EOS_Sessions_RegisterPlayers.
func (s *Sessions) RegisterPlayers(ctx context.Context, sessionName string, players []types.ProductUserId) error {
	cIds := make([]cbinding.EOS_ProductUserId, len(players))
	for i, p := range players {
		cIds[i] = cbinding.EOS_ProductUserId_FromString(string(p))
	}
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_RegisterPlayers(s.handle, sessionName, cIds, clientData)
	})
}

// UnregisterPlayers removes the given players from the named session. See EOS_Sessions_UnregisterPlayers.
func (s *Sessions) UnregisterPlayers(ctx context.Context, sessionName string, players []types.ProductUserId) error {
	cIds := make([]cbinding.EOS_ProductUserId, len(players))
	for i, p := range players {
		cIds[i] = cbinding.EOS_ProductUserId_FromString(string(p))
	}
	return s.simpleAsync(ctx, func(clientData uintptr) {
		cbinding.EOS_Sessions_UnregisterPlayers(s.handle, sessionName, cIds, clientData)
	})
}

// CreateSessionSearch creates a session search handle with the given max result count. See EOS_Sessions_CreateSessionSearch.
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

// AddNotifySessionInviteReceived registers a callback for session invite notifications. See EOS_Sessions_AddNotifySessionInviteReceived.
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
