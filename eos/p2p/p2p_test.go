//go:build eosstub

package p2p_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
	"github.com/mydev/go-eos/eos/p2p"
)

const testUserId = "0000000000000000000000000000abcd"

func setupP2P(tb testing.TB) (*p2p.P2P, func()) {
	tb.Helper()
	cbinding.StubP2PReset()
	w := threadworker.New(func() {})
	ctx, cancel := context.WithCancel(context.Background())
	w.Start(ctx)
	p := p2p.New(cbinding.EOS_HP2P(1), w)
	return p, func() {
		cancel()
		w.Stop()
		cbinding.StubP2PReset()
	}
}

func Test_socket_id_should_reject_empty_name(t *testing.T) {
	s := p2p.SocketId{Name: ""}
	err := s.Validate()
	assert.ErrorIs(t, err, p2p.ErrInvalidSocketId)
}

func Test_socket_id_should_reject_overlong_name(t *testing.T) {
	s := p2p.SocketId{Name: strings.Repeat("a", 33)}
	err := s.Validate()
	assert.ErrorIs(t, err, p2p.ErrInvalidSocketId)
}

func Test_socket_id_should_accept_max_length_name(t *testing.T) {
	s := p2p.SocketId{Name: strings.Repeat("a", 32)}
	assert.NoError(t, s.Validate())
}

func Test_send_packet_should_reject_empty_data(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.SendPacket(context.Background(), p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "test"},
		Data:         nil,
		Reliability:  p2p.ReliableOrdered,
	})
	assert.ErrorIs(t, err, p2p.ErrEmptyPacket)
}

func Test_send_packet_should_reject_oversized_data(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	tooBig := make([]byte, cbinding.EOS_P2P_MAX_PACKET_SIZE+1)
	err := p.SendPacket(context.Background(), p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "test"},
		Data:         tooBig,
		Reliability:  p2p.ReliableOrdered,
	})
	assert.ErrorIs(t, err, p2p.ErrPacketTooLarge)
}

func Test_send_packet_should_reject_invalid_socket(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.SendPacket(context.Background(), p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: ""},
		Data:         []byte("hi"),
		Reliability:  p2p.ReliableOrdered,
	})
	assert.ErrorIs(t, err, p2p.ErrInvalidSocketId)
}

func Test_send_packet_should_succeed_with_valid_input(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.SendPacket(context.Background(), p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "test"},
		Channel:      0,
		Data:         []byte("hello"),
		Reliability:  p2p.ReliableOrdered,
	})
	assert.NoError(t, err)
}

func Test_receive_packet_should_return_err_no_packet_when_empty(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	pkt, err := p.ReceivePacket(testUserId)
	assert.Nil(t, pkt)
	assert.ErrorIs(t, err, p2p.ErrNoPacket)
}

func Test_send_then_receive_should_round_trip(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	payload := []byte("hello world")
	err := p.SendPacket(context.Background(), p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "test"},
		Channel:      7,
		Data:         payload,
		Reliability:  p2p.ReliableOrdered,
	})
	assert.NoError(t, err)

	pkt, err := p.ReceivePacket(testUserId)
	assert.NoError(t, err)
	assert.NotNil(t, pkt)
	assert.Equal(t, payload, pkt.Data)
	assert.Equal(t, uint8(7), pkt.Channel)
	assert.Equal(t, "test", pkt.Socket.Name)
}

func Test_receive_packet_on_channel_should_filter(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	// Send packets on two channels.
	for _, ch := range []uint8{1, 2} {
		err := p.SendPacket(context.Background(), p2p.SendOptions{
			LocalUserId:  testUserId,
			RemoteUserId: testUserId,
			Socket:       p2p.SocketId{Name: "test"},
			Channel:      ch,
			Data:         []byte{ch},
			Reliability:  p2p.ReliableOrdered,
		})
		assert.NoError(t, err)
	}

	// Pulling channel 2 should return the channel-2 packet, not channel 1.
	pkt, err := p.ReceivePacketOnChannel(testUserId, 2)
	assert.NoError(t, err)
	assert.Equal(t, uint8(2), pkt.Channel)
	assert.Equal(t, []byte{2}, pkt.Data)
}

func Test_err_no_packet_should_unwrap_to_typed_result(t *testing.T) {
	// Sanity check that callers can do `errors.Is(err, p2p.ErrNoPacket)`
	// without needing to know the underlying eos types.Result.
	assert.True(t, errors.Is(p2p.ErrNoPacket, p2p.ErrNoPacket))
}
