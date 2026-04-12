// Command p2p-echo demonstrates the EOS P2P Networking interface end to end.
//
// Two processes (server + client) authenticate via DevAuth, exchange a packet
// over EOS P2P, and the client measures round-trip time. This is the example
// referenced by the PRD success criterion "two game processes can discover
// each other via lobby search, join a lobby, and exchange P2P packets".
//
// Polling note (the most important teaching point in this example):
//
//	EOS P2P delivery latency is bounded by how often you call ReceivePacket.
//	EOS_Platform_Tick does NOT deliver packets — they sit in the SDK's
//	internal queue until something pulls them out. The receive loop below
//	polls in a tight goroutine; the polling cadence is the receive latency.
//
// Usage:
//
//	# Terminal 1 — server
//	EOS_DEV_AUTH_HOST=localhost:6547 EOS_DEV_AUTH_CREDENTIAL=player1 \
//	  go run ./examples/p2p-echo --server
//
//	# Terminal 2 — client (substitute the ProductUserId printed by the server)
//	EOS_DEV_AUTH_HOST=localhost:6547 EOS_DEV_AUTH_CREDENTIAL=player2 \
//	  go run ./examples/p2p-echo --client --remote-user=<server-puid>
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/mydev/go-eos/eos/auth"
	"github.com/mydev/go-eos/eos/connect"
	"github.com/mydev/go-eos/eos/p2p"
	"github.com/mydev/go-eos/eos/platform"
	"github.com/mydev/go-eos/eos/types"
)

const (
	socketName = "p2p-echo"
	channel    = uint8(0)
)

var (
	flagServer     = flag.Bool("server", false, "Run as echo server")
	flagClient     = flag.Bool("client", false, "Run as echo client")
	flagRemoteUser = flag.String("remote-user", "", "Remote ProductUserId (client mode only)")
	flagMessage    = flag.String("message", "hello", "Payload to send (client mode only)")
)

func init() {
	// macOS: the EOS SDK's HTTP layer dispatches through the main thread's
	// run loop. Lock the main goroutine to the main OS thread.
	runtime.LockOSThread()
}

func main() {
	flag.Parse()

	if *flagServer == *flagClient {
		fmt.Println("Usage: p2p-echo --server  OR  --client --remote-user=<puid>")
		os.Exit(1)
	}
	if *flagClient && *flagRemoteUser == "" {
		fmt.Println("--client requires --remote-user=<server-product-user-id>")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := platform.PlatformConfig{
		ProductName:    "go-eos-p2p-echo",
		ProductVersion: "1.0.0",
		ProductId:      requireEnv("EOS_PRODUCT_ID"),
		SandboxId:      requireEnv("EOS_SANDBOX_ID"),
		DeploymentId:   requireEnv("EOS_DEPLOYMENT_ID"),
		ClientId:       requireEnv("EOS_CLIENT_ID"),
		ClientSecret:   requireEnv("EOS_CLIENT_SECRET"),
	}

	if err := platform.RunOnMainThread(ctx, cfg, run); err != nil {
		log.Fatal(err)
	}
}

func run(p *platform.Platform) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	productUserId, err := loginFlow(ctx, p)
	if err != nil {
		return err
	}

	// Best-effort NAT type query so users can see their NAT class. The
	// example proceeds even if this fails.
	if nat, natErr := p.P2P().QueryNATType(ctx); natErr == nil {
		fmt.Printf("Local NAT type: %s\n", nat)
	}

	socket := p2p.SocketId{Name: socketName}

	if *flagServer {
		return runServer(ctx, p, productUserId, socket)
	}
	return runClient(ctx, p, productUserId, socket, types.ProductUserId(*flagRemoteUser))
}

func runServer(ctx context.Context, p *platform.Platform, localUserId types.ProductUserId, socket p2p.SocketId) error {
	// Auto-accept any incoming connection on our socket. The handler runs
	// on the worker goroutine — keep it short.
	removeRequest := p.P2P().AddNotifyPeerConnectionRequest(localUserId, &socket, func(req p2p.IncomingConnectionRequest) {
		fmt.Printf("[server] connection request from %s\n", req.RemoteUserId)
		if err := p.P2P().AcceptConnection(req.LocalUserId, req.RemoteUserId, req.Socket); err != nil {
			fmt.Printf("[server] accept failed: %v\n", err)
		}
	})
	defer removeRequest()

	removeEstablished := p.P2P().AddNotifyPeerConnectionEstablished(localUserId, &socket, func(e p2p.PeerConnectionEstablished) {
		fmt.Printf("[server] connection established with %s (network=%d)\n", e.RemoteUserId, e.NetworkType)
	})
	defer removeEstablished()

	fmt.Printf("\nServer ready. Tell the client:\n  --remote-user=%s\n\n", localUserId)
	fmt.Println("Waiting for packets. Press Ctrl+C to quit.")

	pollAndEcho(ctx, p, localUserId)

	// Best-effort: drop any held connections on shutdown.
	if err := p.P2P().CloseConnections(localUserId, socket); err != nil {
		fmt.Printf("[server] close connections: %v\n", err)
	}
	return nil
}

func runClient(ctx context.Context, p *platform.Platform, localUserId types.ProductUserId, socket p2p.SocketId, remoteUserId types.ProductUserId) error {
	if err := p.P2P().AcceptConnection(localUserId, remoteUserId, socket); err != nil {
		return fmt.Errorf("accept: %w", err)
	}
	defer func() {
		if err := p.P2P().CloseConnection(localUserId, remoteUserId, socket); err != nil {
			fmt.Printf("[client] close: %v\n", err)
		}
	}()

	payload := []byte(*flagMessage)
	sentAt := time.Now()
	if err := p.P2P().SendPacket(ctx, p2p.SendOptions{
		LocalUserId:          localUserId,
		RemoteUserId:         remoteUserId,
		Socket:               socket,
		Channel:              channel,
		Data:                 payload,
		Reliability:          p2p.ReliableOrdered,
		AllowDelayedDelivery: true,
	}); err != nil {
		return fmt.Errorf("send: %w", err)
	}
	fmt.Printf("[client] sent %q to %s, waiting for echo...\n", payload, remoteUserId)

	// Polling receive loop with deadline. The cadence here IS the receive
	// latency — see the package GoDoc on p2p.ReceivePacket.
	deadline := time.Now().Add(10 * time.Second)
	tick := time.NewTicker(1 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-tick.C:
			pkt, err := p.P2P().ReceivePacket(localUserId)
			if errors.Is(err, p2p.ErrNoPacket) {
				if time.Now().After(deadline) {
					return fmt.Errorf("timed out waiting for echo")
				}
				continue
			}
			if err != nil {
				return fmt.Errorf("receive: %w", err)
			}
			rtt := time.Since(sentAt)
			fmt.Printf("[client] received echo %q from %s — RTT %s\n", pkt.Data, pkt.Sender, rtt)
			return nil
		}
	}
}

// pollAndEcho is the server-side hot loop: poll ReceivePacket on a tight
// cadence, echo any payload back to its sender on the same channel.
func pollAndEcho(ctx context.Context, p *platform.Platform, localUserId types.ProductUserId) {
	tick := time.NewTicker(1 * time.Millisecond)
	defer tick.Stop()
	socket := p2p.SocketId{Name: socketName}

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			pkt, err := p.P2P().ReceivePacket(localUserId)
			if errors.Is(err, p2p.ErrNoPacket) {
				continue
			}
			if err != nil {
				fmt.Printf("[server] receive: %v\n", err)
				continue
			}
			fmt.Printf("[server] received %q from %s, echoing...\n", pkt.Data, pkt.Sender)
			sendErr := p.P2P().SendPacket(ctx, p2p.SendOptions{
				LocalUserId:          localUserId,
				RemoteUserId:         pkt.Sender,
				Socket:               socket,
				Channel:              pkt.Channel,
				Data:                 pkt.Data,
				Reliability:          p2p.ReliableOrdered,
				AllowDelayedDelivery: true,
			})
			if sendErr != nil {
				fmt.Printf("[server] echo send: %v\n", sendErr)
			}
		}
	}
}

func loginFlow(ctx context.Context, p *platform.Platform) (types.ProductUserId, error) {
	devHost := requireEnv("EOS_DEV_AUTH_HOST")
	devCred := requireEnv("EOS_DEV_AUTH_CREDENTIAL")

	fmt.Println("Logging in via Developer Auth Tool...")
	loginResult, err := p.Auth().Login(ctx, auth.LoginOptions{
		CredentialType: types.LoginCredentialDeveloper,
		ID:             devHost,
		Token:          devCred,
	})
	if err != nil {
		return "", fmt.Errorf("auth login: %w", err)
	}
	fmt.Printf("Auth: EpicAccountId=%s\n", loginResult.LocalUserId)

	idToken, err := p.Auth().CopyIdToken(loginResult.LocalUserId)
	if err != nil {
		return "", fmt.Errorf("copy id token: %w", err)
	}

	connectResult, err := p.Connect().Login(ctx, connect.LoginOptions{
		CredentialType: types.ExternalCredentialEpicIDToken,
		Token:          idToken,
	})
	if err == nil {
		fmt.Printf("Connect: ProductUserId=%s\n", connectResult.LocalUserId)
		return connectResult.LocalUserId, nil
	}
	if connectResult == nil || connectResult.ContinuanceToken == 0 {
		return "", fmt.Errorf("connect login: %w", err)
	}
	fmt.Println("First run — creating product user...")
	userId, createErr := p.Connect().CreateUser(ctx, connectResult.ContinuanceToken)
	if createErr != nil {
		return "", fmt.Errorf("create user: %w", createErr)
	}
	fmt.Printf("Connect: ProductUserId=%s\n", *userId)
	return *userId, nil
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}
