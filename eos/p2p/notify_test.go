//go:build eosstub

package p2p_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/p2p"
)

func Test_add_notify_peer_connection_request_should_return_remove_func(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	remove := p.AddNotifyPeerConnectionRequest(testUserId, nil, func(p2p.IncomingConnectionRequest) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}

func Test_add_notify_peer_connection_request_with_socket_filter(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	socket := p2p.SocketId{Name: "filtered"}
	remove := p.AddNotifyPeerConnectionRequest(testUserId, &socket, func(p2p.IncomingConnectionRequest) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}

func Test_add_notify_peer_connection_established_should_return_remove_func(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	remove := p.AddNotifyPeerConnectionEstablished(testUserId, nil, func(p2p.PeerConnectionEstablished) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}

func Test_add_notify_peer_connection_closed_should_return_remove_func(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	remove := p.AddNotifyPeerConnectionClosed(testUserId, nil, func(p2p.PeerConnectionClosed) {})
	assert.NotNil(t, remove)
	assert.NotPanics(t, func() { remove() })
}

func Test_remove_notify_should_be_idempotent(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	remove := p.AddNotifyPeerConnectionEstablished(testUserId, nil, func(p2p.PeerConnectionEstablished) {})
	assert.NotPanics(t, func() {
		remove()
		remove()
	})
}

func Test_remove_notify_should_be_safe_under_concurrent_calls(t *testing.T) {
	p, cleanup := setupP2P(t)
	defer cleanup()

	remove := p.AddNotifyPeerConnectionEstablished(testUserId, nil, func(p2p.PeerConnectionEstablished) {})

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NotPanics(t, func() { remove() })
		}()
	}
	wg.Wait()
}
