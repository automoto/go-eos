//go:build eosstub

package p2p_test

import (
	"context"
	"testing"

	"github.com/mydev/go-eos/eos/p2p"
)

// BenchmarkSendPacket_Stub measures the wrapper overhead for SendPacket
// (validation + worker dispatch + struct conversion + stub callback).
// PRD success criterion: total wrapper overhead < 5ms; the stub-only path
// should be well under that — current target is < 100µs/op.
func BenchmarkSendPacket_Stub(b *testing.B) {
	p, cleanup := setupP2P(b)
	defer cleanup()

	payload := make([]byte, 256)
	opts := p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "bench"},
		Data:         payload,
		Reliability:  p2p.ReliableOrdered,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := p.SendPacket(context.Background(), opts); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkReceivePacket_Stub measures the receive-path wrapper overhead.
// Each iteration sends one packet first so the queue has data to drain.
func BenchmarkReceivePacket_Stub(b *testing.B) {
	p, cleanup := setupP2P(b)
	defer cleanup()

	payload := make([]byte, 256)
	opts := p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "bench"},
		Data:         payload,
		Reliability:  p2p.ReliableOrdered,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		if err := p.SendPacket(context.Background(), opts); err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		if _, err := p.ReceivePacket(testUserId); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSendReceiveLoopback measures the full round-trip through the
// stub layer in a single timed step. Useful to compare against
// BenchmarkSendPacket_Stub + BenchmarkReceivePacket_Stub for sanity.
func BenchmarkSendReceiveLoopback(b *testing.B) {
	p, cleanup := setupP2P(b)
	defer cleanup()

	payload := make([]byte, 256)
	opts := p2p.SendOptions{
		LocalUserId:  testUserId,
		RemoteUserId: testUserId,
		Socket:       p2p.SocketId{Name: "bench"},
		Data:         payload,
		Reliability:  p2p.ReliableOrdered,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := p.SendPacket(context.Background(), opts); err != nil {
			b.Fatal(err)
		}
		if _, err := p.ReceivePacket(testUserId); err != nil {
			b.Fatal(err)
		}
	}
}
