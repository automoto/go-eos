// Command connect-deviceid demonstrates fully anonymous EOS authentication
// using the Device ID feature of the Connect interface.
//
// Unlike auth-login, this example does NOT use the Auth interface and does
// NOT require the player to have an Epic Games account. The Connect interface
// creates a per-device pseudo-account on first run; subsequent runs reuse it.
//
// This is the recommended path for casual games and games that ship outside
// the Epic Games Store, where forcing players to create an Epic account would
// hurt conversion. The resulting ProductUserId works with all multiplayer
// features (lobbies, sessions, P2P) exactly the same as one obtained via the
// Auth interface.
//
// Limitations:
//   - The Device ID is local to this device. The same player on a second
//     device is treated as a different user unless they explicitly link an
//     external account (Steam, Discord, etc.) via Connect.LinkAccount.
//   - If the user uninstalls the game or wipes local storage, their progress
//     is lost (no recoverable identity).
//
// Environment variables (subset of auth-login — no DevAuth needed):
//
//	EOS_PRODUCT_ID, EOS_SANDBOX_ID, EOS_DEPLOYMENT_ID
//	EOS_CLIENT_ID, EOS_CLIENT_SECRET
//
// Usage:
//
//	go run ./examples/connect-deviceid
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/mydev/go-eos/eos/connect"
	"github.com/mydev/go-eos/eos/platform"
	"github.com/mydev/go-eos/eos/types"
)

func init() {
	// Lock the main goroutine to the main OS thread. On macOS the EOS SDK's
	// HTTP layer uses Apple networking which dispatches through the main
	// thread's run loop.
	runtime.LockOSThread()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := platform.PlatformConfig{
		ProductName:    "go-eos-connect-deviceid",
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

	// Step 1: ensure a Device ID exists for this device. Idempotent —
	// returns nil if already created on a previous run.
	fmt.Println("Ensuring local Device ID exists...")
	if err := p.Connect().CreateDeviceId(ctx, deviceModel()); err != nil {
		return fmt.Errorf("create device id: %w", err)
	}

	// Step 2: log in to Connect with DeviceIDAccessToken. The Token field
	// is empty for this credential type — the SDK uses the local device's
	// stored credentials. DisplayName is required and shown in the EOS
	// Developer Portal user list.
	fmt.Println("Logging in to Connect with Device ID...")
	loginResult, err := p.Connect().Login(ctx, connect.LoginOptions{
		CredentialType: types.ExternalCredentialDeviceIDAccessToken,
		DisplayName:    "AnonymousPlayer",
	})
	if err == nil {
		fmt.Printf("Anonymous ProductUserId: %s\n", loginResult.LocalUserId)
	} else {
		// First-run path: SDK returns "user not found" with a continuance
		// token. Use it to create the product user.
		if loginResult == nil || loginResult.ContinuanceToken == 0 {
			return fmt.Errorf("connect login: %w", err)
		}
		fmt.Println("First run — creating product user...")
		userId, createErr := p.Connect().CreateUser(ctx, loginResult.ContinuanceToken)
		if createErr != nil {
			return fmt.Errorf("create user: %w", createErr)
		}
		fmt.Printf("Anonymous ProductUserId: %s\n", *userId)
	}

	fmt.Println()
	fmt.Println("Logged in anonymously — no Epic account required.")
	fmt.Println("This ProductUserId can join lobbies, send P2P packets, etc.")
	fmt.Println("Press Ctrl+C to quit.")
	<-ctx.Done()
	return nil
}

// deviceModel returns a free-form description of the host. The EOS portal
// shows this string to help operators identify devices linked to a user.
// Max 64 UTF-8 characters; longer strings are silently truncated.
func deviceModel() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}
