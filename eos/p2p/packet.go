package p2p

import (
	"context"
	"fmt"

	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/types"
)

// Reliability mirrors EOS_EPacketReliability.
type Reliability int

const (
	UnreliableUnordered Reliability = 0
	ReliableUnordered   Reliability = 1
	ReliableOrdered     Reliability = 2
)

// SendOptions describes a single outgoing packet.
type SendOptions struct {
	LocalUserId  types.ProductUserId
	RemoteUserId types.ProductUserId
	Socket       SocketId
	Channel      uint8
	Data         []byte
	Reliability  Reliability

	// AllowDelayedDelivery: if false and there is no open connection to the
	// peer, the packet is dropped. If true, it is queued until the connection
	// is established.
	AllowDelayedDelivery bool

	// DisableAutoAcceptConnection: if true, SendPacket will not implicitly
	// open a connection — the caller must call AcceptConnection first.
	DisableAutoAcceptConnection bool
}

// IncomingPacket is the result of ReceivePacket. Data is a fresh
// Go-allocated buffer; see MEM-5 in docs/prd.md and the package GoDoc for
// the single-copy reasoning.
type IncomingPacket struct {
	Sender  types.ProductUserId
	Socket  SocketId
	Channel uint8
	Data    []byte
}

// SendPacket queues a packet for delivery to RemoteUserId. Returns nil on
// success or a typed error (ErrEmptyPacket, ErrPacketTooLarge,
// ErrInvalidSocketId) for Go-side validation failures, or an *eos/types.Result
// for SDK-reported errors.
func (p *P2P) SendPacket(ctx context.Context, opts SendOptions) error {
	if len(opts.Data) == 0 {
		return ErrEmptyPacket
	}
	if len(opts.Data) > cbinding.EOS_P2P_MAX_PACKET_SIZE {
		return ErrPacketTooLarge
	}
	if err := opts.Socket.Validate(); err != nil {
		return err
	}

	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_SendPacket(p.handle, &cbinding.EOS_P2P_SendPacketOptions{
			LocalUserId:                 cbinding.EOS_ProductUserId_FromString(string(opts.LocalUserId)),
			RemoteUserId:                cbinding.EOS_ProductUserId_FromString(string(opts.RemoteUserId)),
			SocketId:                    &cbinding.EOS_P2P_SocketId{Name: opts.Socket.Name},
			Channel:                     opts.Channel,
			Data:                        opts.Data,
			AllowDelayedDelivery:        opts.AllowDelayedDelivery,
			Reliability:                 cbinding.EOS_EPacketReliability(opts.Reliability),
			DisableAutoAcceptConnection: opts.DisableAutoAcceptConnection,
		})
	}); err != nil {
		return fmt.Errorf("p2p send: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// ReceivePacket pulls the next queued packet for localUserId from any
// channel. Returns ErrNoPacket if the queue is empty.
//
// The returned IncomingPacket.Data is a fresh Go-allocated buffer; the SDK
// memory is not borrowed (see MEM-5 reasoning in package GoDoc).
//
// This method is poll-based: callers must call it on a tight cadence in
// their game loop. EOS_Platform_Tick does not deliver packets — they sit in
// the SDK's internal queue until ReceivePacket pulls them out.
func (p *P2P) ReceivePacket(localUserId types.ProductUserId) (*IncomingPacket, error) {
	return p.receive(localUserId, nil)
}

// ReceivePacketOnChannel is like ReceivePacket but only returns packets on
// the specified channel. See ReceivePacket for the polling discussion.
func (p *P2P) ReceivePacketOnChannel(localUserId types.ProductUserId, channel uint8) (*IncomingPacket, error) {
	return p.receive(localUserId, &channel)
}

func (p *P2P) receive(localUserId types.ProductUserId, channel *uint8) (*IncomingPacket, error) {
	cUserId := cbinding.EOS_ProductUserId_FromString(string(localUserId))

	var size uint32
	var sizeResult cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		size, sizeResult = cbinding.EOS_P2P_GetNextReceivedPacketSize(p.handle, cUserId, channel)
	}); err != nil {
		return nil, fmt.Errorf("p2p get next size: %w", err)
	}
	if sizeResult == cbinding.EOS_EResult_NotFound || size == 0 {
		return nil, ErrNoPacket
	}
	if sizeResult != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(sizeResult))
	}

	var result *cbinding.EOS_P2P_ReceivePacketResult
	var recvResult cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result, recvResult = cbinding.EOS_P2P_ReceivePacket(p.handle, cUserId, size, channel)
	}); err != nil {
		return nil, fmt.Errorf("p2p receive: %w", err)
	}
	if recvResult == cbinding.EOS_EResult_NotFound {
		return nil, ErrNoPacket
	}
	if recvResult != cbinding.EOS_EResult_Success {
		return nil, types.NewResult(int(recvResult))
	}

	return &IncomingPacket{
		Sender:  types.ProductUserId(cbinding.EOS_ProductUserId_ToString(result.PeerId)),
		Socket:  SocketId{Name: result.SocketId.Name},
		Channel: result.Channel,
		Data:    result.Data,
	}, nil
}
