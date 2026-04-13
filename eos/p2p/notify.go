package p2p

import (
	"runtime/cgo"
	"sync"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/types"
)

// ConnectionEstablishedType reports whether a Peer Connection Established
// notification is for a brand-new connection or a reconnect / network-type
// transition (direct ↔ relayed).
type ConnectionEstablishedType int

const (
	// NewConnection indicates a brand-new P2P connection.
	NewConnection ConnectionEstablishedType = 0
	// ReconnectionEstablished indicates a reconnect or network-type transition.
	ReconnectionEstablished ConnectionEstablishedType = 1
)

// NetworkConnectionType reports the underlying transport for an established
// connection.
type NetworkConnectionType int

const (
	// NoNetworkConnection indicates no underlying network connection.
	NoNetworkConnection NetworkConnectionType = 0
	// DirectConnection indicates a direct peer-to-peer connection.
	DirectConnection NetworkConnectionType = 1
	// RelayedConnection indicates the connection is routed through an Epic relay server.
	RelayedConnection NetworkConnectionType = 2
)

// ConnectionClosedReason mirrors EOS_EConnectionClosedReason.
type ConnectionClosedReason int

const (
	// ClosedUnknown indicates the connection closed for an unknown reason.
	ClosedUnknown ConnectionClosedReason = 0
	// ClosedByLocalUser indicates the local user closed the connection.
	ClosedByLocalUser ConnectionClosedReason = 1
	// ClosedByPeer indicates the remote peer closed the connection.
	ClosedByPeer ConnectionClosedReason = 2
	// ClosedTimedOut indicates the connection timed out.
	ClosedTimedOut ConnectionClosedReason = 3
	// ClosedTooManyConnections indicates the connection was dropped due to too many connections.
	ClosedTooManyConnections ConnectionClosedReason = 4
	// ClosedInvalidMessage indicates the connection was closed due to an invalid message.
	ClosedInvalidMessage ConnectionClosedReason = 5
	// ClosedInvalidData indicates the connection was closed due to invalid data.
	ClosedInvalidData ConnectionClosedReason = 6
	// ClosedConnectionFailed indicates the connection attempt failed.
	ClosedConnectionFailed ConnectionClosedReason = 7
	// ClosedConnectionClosed indicates the connection was already closed.
	ClosedConnectionClosed ConnectionClosedReason = 8
	// ClosedNegotiationFailed indicates the connection negotiation failed.
	ClosedNegotiationFailed ConnectionClosedReason = 9
	// ClosedUnexpectedError indicates an unexpected error closed the connection.
	ClosedUnexpectedError ConnectionClosedReason = 10
	// ClosedConnectionIgnored indicates the incoming connection request was ignored.
	ClosedConnectionIgnored ConnectionClosedReason = 11
)

// IncomingConnectionRequest is the payload for AddNotifyPeerConnectionRequest.
type IncomingConnectionRequest struct {
	LocalUserId  types.ProductUserId
	RemoteUserId types.ProductUserId
	Socket       SocketId
}

// PeerConnectionEstablished is the payload for AddNotifyPeerConnectionEstablished.
type PeerConnectionEstablished struct {
	LocalUserId    types.ProductUserId
	RemoteUserId   types.ProductUserId
	Socket         SocketId
	ConnectionType ConnectionEstablishedType
	NetworkType    NetworkConnectionType
}

// PeerConnectionClosed is the payload for AddNotifyPeerConnectionClosed.
type PeerConnectionClosed struct {
	LocalUserId  types.ProductUserId
	RemoteUserId types.ProductUserId
	Socket       SocketId
	Reason       ConnectionClosedReason
}

// RemoveNotifyFunc tears down a previously registered notification. Safe
// to call multiple times — the second call is a no-op.
type RemoveNotifyFunc = callback.RemoveNotifyFunc

// AddNotifyPeerConnectionRequest fires when a remote peer wants to open a
// connection. The handler decides whether to call AcceptConnection. socket
// may be nil to receive requests for all sockets.
//
// The handler runs on the platform's worker goroutine during Tick. SDK
// calls (AcceptConnection, SendPacket, etc.) may be invoked directly from
// inside the handler — the worker detects re-entrance and executes them
// inline on the same locked OS thread. Keep the handler short to avoid
// starving the tick loop; offload long-running work to a goroutine.
func (p *P2P) AddNotifyPeerConnectionRequest(localUserId types.ProductUserId, socket *SocketId, fn func(IncomingConnectionRequest)) RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_P2P_OnIncomingConnectionRequestInfo)
		fn(IncomingConnectionRequest{
			LocalUserId:  types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId)),
			RemoteUserId: types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.RemoteUserId)),
			Socket:       SocketId{Name: info.SocketId.Name},
		})
	})
	handle := cgo.NewHandle(notifyFn)

	var notifyId cbinding.EOS_NotificationId
	if err := p.worker.Submit(func() {
		notifyId = cbinding.EOS_P2P_AddNotifyPeerConnectionRequest(p.handle,
			cbinding.EOS_ProductUserId_FromString(string(localUserId)),
			toCSocket(socket), uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}
	return makeRemove(handle, func() {
		_ = p.worker.Submit(func() {
			cbinding.EOS_P2P_RemoveNotifyPeerConnectionRequest(p.handle, notifyId)
		})
	})
}

// AddNotifyPeerConnectionEstablished fires when a connection is opened or
// reopened (e.g., after a network type transition). socket may be nil for
// all sockets.
//
// The handler runs on the platform's worker goroutine — see
// AddNotifyPeerConnectionRequest for the threading note.
func (p *P2P) AddNotifyPeerConnectionEstablished(localUserId types.ProductUserId, socket *SocketId, fn func(PeerConnectionEstablished)) RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_P2P_OnPeerConnectionEstablishedInfo)
		fn(PeerConnectionEstablished{
			LocalUserId:    types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId)),
			RemoteUserId:   types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.RemoteUserId)),
			Socket:         SocketId{Name: info.SocketId.Name},
			ConnectionType: ConnectionEstablishedType(info.ConnectionType),
			NetworkType:    NetworkConnectionType(info.NetworkType),
		})
	})
	handle := cgo.NewHandle(notifyFn)

	var notifyId cbinding.EOS_NotificationId
	if err := p.worker.Submit(func() {
		notifyId = cbinding.EOS_P2P_AddNotifyPeerConnectionEstablished(p.handle,
			cbinding.EOS_ProductUserId_FromString(string(localUserId)),
			toCSocket(socket), uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}
	return makeRemove(handle, func() {
		_ = p.worker.Submit(func() {
			cbinding.EOS_P2P_RemoveNotifyPeerConnectionEstablished(p.handle, notifyId)
		})
	})
}

// AddNotifyPeerConnectionClosed fires when a previously open or pending
// connection closes. socket may be nil for all sockets.
//
// The handler runs on the platform's worker goroutine — see
// AddNotifyPeerConnectionRequest for the threading note.
func (p *P2P) AddNotifyPeerConnectionClosed(localUserId types.ProductUserId, socket *SocketId, fn func(PeerConnectionClosed)) RemoveNotifyFunc {
	notifyFn := callback.NotifyFunc(func(data any) {
		info := data.(*cbinding.EOS_P2P_OnRemoteConnectionClosedInfo)
		fn(PeerConnectionClosed{
			LocalUserId:  types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.LocalUserId)),
			RemoteUserId: types.ProductUserId(cbinding.EOS_ProductUserId_ToString(info.RemoteUserId)),
			Socket:       SocketId{Name: info.SocketId.Name},
			Reason:       ConnectionClosedReason(info.Reason),
		})
	})
	handle := cgo.NewHandle(notifyFn)

	var notifyId cbinding.EOS_NotificationId
	if err := p.worker.Submit(func() {
		notifyId = cbinding.EOS_P2P_AddNotifyPeerConnectionClosed(p.handle,
			cbinding.EOS_ProductUserId_FromString(string(localUserId)),
			toCSocket(socket), uintptr(handle))
	}); err != nil {
		handle.Delete()
		return func() {}
	}
	return makeRemove(handle, func() {
		_ = p.worker.Submit(func() {
			cbinding.EOS_P2P_RemoveNotifyPeerConnectionClosed(p.handle, notifyId)
		})
	})
}

func toCSocket(s *SocketId) *cbinding.EOS_P2P_SocketId {
	if s == nil {
		return nil
	}
	return &cbinding.EOS_P2P_SocketId{Name: s.Name}
}

// makeRemove returns an idempotent RemoveNotifyFunc that runs the SDK-side
// remove and then deletes the cgo.Handle. Safe to call multiple times from
// any goroutine — only the first call has any effect.
func makeRemove(handle cgo.Handle, sdkRemove func()) RemoveNotifyFunc {
	var once sync.Once
	return func() {
		once.Do(func() {
			sdkRemove()
			handle.Delete()
		})
	}
}
