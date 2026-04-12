// Package p2p wraps the EOS P2P Networking interface.
//
// # Polling-based receive
//
// Unlike Auth/Connect/Lobby/Sessions, where async operations complete via
// callbacks fired during Tick, P2P's receive path is poll-based: the caller's
// game loop must call ReceivePacket periodically. EOS_Platform_Tick does NOT
// synthesize a callback for queued incoming packets — they sit in the SDK's
// internal queue until something pulls them out. The polling frequency
// directly determines the receive latency.
//
// # Connection-request flow
//
// The typical pattern for accepting incoming connections is:
//
//  1. Register AddNotifyPeerConnectionRequest.
//  2. When a request arrives, the callback decides whether to accept.
//  3. If accepting, call P2P.AcceptConnection.
//  4. Once accepted, AddNotifyPeerConnectionEstablished fires.
package p2p

import (
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
)

// P2P is the public handle for the EOS P2P Networking interface. Construct
// one via Platform.P2P() — direct construction is reserved for internal use.
type P2P struct {
	handle cbinding.EOS_HP2P
	worker *threadworker.Worker
}

// New constructs a P2P from a raw cbinding handle and the platform's worker.
// Game code should use Platform.P2P() instead of calling this directly.
func New(handle cbinding.EOS_HP2P, worker *threadworker.Worker) *P2P {
	return &P2P{handle: handle, worker: worker}
}
