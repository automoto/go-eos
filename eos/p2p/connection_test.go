//go:build eosstub

package p2p_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/p2p"
)

func Test_accept_connection_should_succeed(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.AcceptConnection(testUserId, testUserId, p2p.SocketId{Name: "test"})
	assert.NoError(t, err)
}

func Test_accept_connection_should_validate_socket(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.AcceptConnection(testUserId, testUserId, p2p.SocketId{Name: ""})
	assert.ErrorIs(t, err, p2p.ErrInvalidSocketId)
}

func Test_close_connection_should_succeed(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.CloseConnection(testUserId, testUserId, p2p.SocketId{Name: "test"})
	assert.NoError(t, err)
}

func Test_close_connections_should_succeed(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.CloseConnections(testUserId, p2p.SocketId{Name: "test"})
	assert.NoError(t, err)
}

func Test_query_nat_type_should_return_open_from_stub(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	nat, err := p.QueryNATType(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, p2p.NATOpen, nat)
}

func Test_query_nat_type_should_honor_stub_override(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()
	cbinding.StubP2PSetNATType(cbinding.EOS_NAT_Strict)

	nat, err := p.QueryNATType(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, p2p.NATStrict, nat)
}

func Test_get_nat_type_should_return_cached_value(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()
	cbinding.StubP2PSetNATType(cbinding.EOS_NAT_Moderate)

	nat, err := p.GetNATType()
	assert.NoError(t, err)
	assert.Equal(t, p2p.NATModerate, nat)
}

func Test_nat_type_string_should_return_human_readable(t *testing.T) {
	assert.Equal(t, "Open", p2p.NATOpen.String())
	assert.Equal(t, "Moderate", p2p.NATModerate.String())
	assert.Equal(t, "Strict", p2p.NATStrict.String())
	assert.Equal(t, "Unknown", p2p.NATUnknown.String())
}

func Test_set_then_get_relay_control_round_trips(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	assert.NoError(t, p.SetRelayControl(p2p.ForceRelays))
	rc, err := p.GetRelayControl()
	assert.NoError(t, err)
	assert.Equal(t, p2p.ForceRelays, rc)
}

func Test_set_then_get_port_range_round_trips(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	assert.NoError(t, p.SetPortRange(p2p.PortRange{Port: 5000, MaxAdditionalPortsToTry: 10}))
	r, err := p.GetPortRange()
	assert.NoError(t, err)
	assert.Equal(t, uint16(5000), r.Port)
	assert.Equal(t, uint16(10), r.MaxAdditionalPortsToTry)
}

func Test_get_packet_queue_info_should_return_empty_for_stub(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	info, err := p.GetPacketQueueInfo()
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), info.IncomingCurrentPackets)
	assert.Equal(t, uint64(0), info.OutgoingCurrentPackets)
}

func Test_set_packet_queue_size_should_succeed(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.SetPacketQueueSize(1<<20, 1<<20)
	assert.NoError(t, err)
}

func Test_clear_packet_queue_should_drain_pending(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	err := p.SendPacket(context.Background(), p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "test"},
		Data:         []byte("data"),
		Reliability:  p2p.ReliableOrdered,
	})
	assert.NoError(t, err)

	assert.NoError(t, p.ClearPacketQueue(testUserId, testUserId, p2p.SocketId{Name: "test"}))

	pkt, recvErr := p.ReceivePacket(testUserId)
	assert.Nil(t, pkt)
	assert.ErrorIs(t, recvErr, p2p.ErrNoPacket)
}
