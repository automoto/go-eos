package p2p

import (
	"fmt"

	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/types"
)

// AcceptConnection accepts (or proactively requests) a P2P connection from
// remoteUserId on the given socket. Idempotent — calling on an
// already-accepted connection returns nil.
func (p *P2P) AcceptConnection(localUserId, remoteUserId types.ProductUserId, socket SocketId) error {
	if err := socket.Validate(); err != nil {
		return err
	}
	cSocket := &cbinding.EOS_P2P_SocketId{Name: socket.Name}
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_AcceptConnection(p.handle, &cbinding.EOS_P2P_AcceptConnectionOptions{
			LocalUserId:  cbinding.EOS_ProductUserId_FromString(string(localUserId)),
			RemoteUserId: cbinding.EOS_ProductUserId_FromString(string(remoteUserId)),
			SocketId:     cSocket,
		})
	}); err != nil {
		return fmt.Errorf("p2p accept: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// CloseConnection closes a single peer connection on the given socket. Any
// queued packets (including reliable ones) are flushed.
func (p *P2P) CloseConnection(localUserId, remoteUserId types.ProductUserId, socket SocketId) error {
	if err := socket.Validate(); err != nil {
		return err
	}
	cSocket := &cbinding.EOS_P2P_SocketId{Name: socket.Name}
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_CloseConnection(p.handle, &cbinding.EOS_P2P_CloseConnectionOptions{
			LocalUserId:  cbinding.EOS_ProductUserId_FromString(string(localUserId)),
			RemoteUserId: cbinding.EOS_ProductUserId_FromString(string(remoteUserId)),
			SocketId:     cSocket,
		})
	}); err != nil {
		return fmt.Errorf("p2p close: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// CloseConnections closes all connections for localUserId on the given
// socket. Useful in shutdown paths to flush any queued reliable packets
// before the platform releases.
func (p *P2P) CloseConnections(localUserId types.ProductUserId, socket SocketId) error {
	if err := socket.Validate(); err != nil {
		return err
	}
	cSocket := &cbinding.EOS_P2P_SocketId{Name: socket.Name}
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_CloseConnections(p.handle, &cbinding.EOS_P2P_CloseConnectionsOptions{
			LocalUserId: cbinding.EOS_ProductUserId_FromString(string(localUserId)),
			SocketId:    cSocket,
		})
	}); err != nil {
		return fmt.Errorf("p2p close all: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}
