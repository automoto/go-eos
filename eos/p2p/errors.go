package p2p

import (
	"errors"

	"github.com/mydev/go-eos/eos/types"
)

// ErrNoPacket indicates that ReceivePacket was called but no packets are
// queued for the requested user (and channel, if specified). It wraps the
// EOS_NotFound result code.
var ErrNoPacket = types.NewResult(types.CodeNotFound)

// ErrEmptyPacket is returned by SendPacket when called with no payload.
// EOS_P2P_SendPacket technically allows zero-length data but it is almost
// always a bug, so the wrapper rejects it Go-side.
var ErrEmptyPacket = errors.New("p2p: empty packet")

// ErrPacketTooLarge is returned by SendPacket when the payload exceeds
// EOS_P2P_MAX_PACKET_SIZE (1170 bytes).
var ErrPacketTooLarge = errors.New("p2p: packet exceeds EOS_P2P_MAX_PACKET_SIZE")

// ErrInvalidSocketId is returned when a SocketId fails Validate (empty name
// or longer than 32 chars).
var ErrInvalidSocketId = errors.New("p2p: socket id name must be 1-32 chars")
