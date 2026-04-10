// Command auth-login demonstrates EOS authentication using the Developer Auth Tool.
//
// Prerequisites:
//   - EOS Developer Auth Tool running (download from Epic Developer Portal)
//   - Valid EOS credentials set as environment variables
//
// Environment variables:
//
//	EOS_PRODUCT_ID, EOS_SANDBOX_ID, EOS_DEPLOYMENT_ID
//	EOS_CLIENT_ID, EOS_CLIENT_SECRET
//	EOS_DEV_AUTH_HOST      - DevAuth Tool address (e.g., "localhost:6547")
//	EOS_DEV_AUTH_CREDENTIAL - DevAuth credential name
//
// Usage:
//
//	go run ./examples/auth-login
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	"github.com/mydev/go-eos/eos/auth"
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
		ProductName:    "go-eos-example",
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

	// Auth login via Developer Auth Tool
	devHost := requireEnv("EOS_DEV_AUTH_HOST")
	devCred := requireEnv("EOS_DEV_AUTH_CREDENTIAL")

	fmt.Println("Logging in via Developer Auth Tool...")
	loginResult, err := p.Auth().Login(ctx, auth.LoginOptions{
		CredentialType: types.LoginCredentialDeveloper,
		ID:             devHost,
		Token:          devCred,
	})
	if err != nil {
		return fmt.Errorf("auth login: %w", err)
	}
	fmt.Printf("Auth login successful!\n")
	fmt.Printf("  EpicAccountId: %s\n", loginResult.LocalUserId)

	// Register for auth status changes
	removeAuthNotify := p.Auth().AddNotifyLoginStatusChanged(func(info auth.LoginStatusChangedInfo) {
		fmt.Printf("Auth status changed: %d -> %d (user: %s)\n",
			info.PrevStatus, info.CurrentStatus, info.LocalUserId)
	})
	defer removeAuthNotify()

	// Copy ID token for Connect login
	idToken, err := p.Auth().CopyIdToken(loginResult.LocalUserId)
	if err != nil {
		return fmt.Errorf("copy id token: %w", err)
	}

	// Connect login using the ID token
	fmt.Println("Logging in to Connect...")
	connectResult, err := p.Connect().Login(ctx, connect.LoginOptions{
		CredentialType: types.ExternalCredentialEpicIDToken,
		Token:          idToken,
	})
	if err != nil {
		// If user doesn't exist, create one
		var eosErr *types.Result
		if connectResult != nil && connectResult.ContinuanceToken != 0 {
			fmt.Println("Product user not found, creating...")
			userId, createErr := p.Connect().CreateUser(ctx, connectResult.ContinuanceToken)
			if createErr != nil {
				return fmt.Errorf("create user: %w", createErr)
			}
			fmt.Printf("  ProductUserId: %s\n", *userId)
		} else if err != nil {
			_ = eosErr
			return fmt.Errorf("connect login: %w", err)
		}
	} else {
		fmt.Printf("Connect login successful!\n")
		fmt.Printf("  ProductUserId: %s\n", connectResult.LocalUserId)
	}

	// Register for connect notifications
	removeConnectNotify := p.Connect().AddNotifyLoginStatusChanged(func(info connect.LoginStatusChangedInfo) {
		fmt.Printf("Connect status changed: %d -> %d (user: %s)\n",
			info.PreviousStatus, info.CurrentStatus, info.LocalUserId)
	})
	defer removeConnectNotify()

	fmt.Println("\nLogged in. Press Ctrl+C to quit.")
	<-ctx.Done()
	fmt.Println("\nShutting down...")
	return nil
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}
