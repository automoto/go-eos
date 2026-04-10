// Command lobby-chat demonstrates EOS lobby-based matchmaking.
//
// Usage:
//
//	# Terminal 1 — create a lobby
//	go run ./examples/lobby-chat --create --bucket-id "game:mode=deathmatch"
//
//	# Terminal 2 — search and join
//	go run ./examples/lobby-chat --join --bucket-id "game:mode=deathmatch"
//
// Environment variables (same as auth-login):
//
//	EOS_PRODUCT_ID, EOS_SANDBOX_ID, EOS_DEPLOYMENT_ID
//	EOS_CLIENT_ID, EOS_CLIENT_SECRET
//	EOS_DEV_AUTH_HOST, EOS_DEV_AUTH_CREDENTIAL
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/mydev/go-eos/eos/auth"
	"github.com/mydev/go-eos/eos/connect"
	"github.com/mydev/go-eos/eos/lobby"
	"github.com/mydev/go-eos/eos/platform"
	"github.com/mydev/go-eos/eos/types"
)

func init() {
	runtime.LockOSThread()
}

var (
	flagCreate   = flag.Bool("create", false, "Create a new lobby")
	flagJoin     = flag.Bool("join", false, "Search for and join a lobby")
	flagBucketId = flag.String("bucket-id", "game:mode=default", "Lobby bucket ID for search/create")
)

func main() {
	flag.Parse()

	if !*flagCreate && !*flagJoin {
		fmt.Println("Usage: lobby-chat --create|--join --bucket-id <bucket>")
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := platform.PlatformConfig{
		ProductName:    "go-eos-lobby-chat",
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

	// Auth + Connect login (same as auth-login example)
	productUserId, err := loginFlow(ctx, p)
	if err != nil {
		return err
	}

	// Register for lobby member events
	removeStatus := p.Lobby().AddNotifyLobbyMemberStatusReceived(func(info lobby.MemberStatusInfo) {
		fmt.Printf("[lobby] member %s status: %d (lobby: %s)\n",
			info.TargetUserId, info.CurrentStatus, info.LobbyId)
	})
	defer removeStatus()

	removeUpdate := p.Lobby().AddNotifyLobbyUpdateReceived(func(info lobby.LobbyUpdateInfo) {
		fmt.Printf("[lobby] updated: %s\n", info.LobbyId)
	})
	defer removeUpdate()

	var activeLobbyId string

	if *flagCreate {
		lobbyId, createErr := p.Lobby().CreateLobby(ctx, lobby.CreateLobbyOptions{
			LocalUserId:     productUserId,
			MaxMembers:      4,
			PermissionLevel: lobby.PermissionPublicAdvertised,
			AllowInvites:    true,
			BucketId:        *flagBucketId,
		})
		if createErr != nil {
			return fmt.Errorf("create lobby: %w", createErr)
		}
		fmt.Printf("Lobby created: %s\n", lobbyId)
		activeLobbyId = lobbyId

		// Set a lobby attribute
		mod, modErr := p.Lobby().UpdateLobbyModification(productUserId, lobbyId)
		if modErr != nil {
			return fmt.Errorf("update lobby modification: %w", modErr)
		}
		if attrErr := mod.AddAttribute("game_mode", "deathmatch", lobby.VisibilityPublic); attrErr != nil {
			mod.Release()
			return fmt.Errorf("add attribute: %w", attrErr)
		}
		if updateErr := p.Lobby().UpdateLobby(ctx, mod); updateErr != nil {
			mod.Release()
			return fmt.Errorf("update lobby: %w", updateErr)
		}
		mod.Release()
		fmt.Println("Lobby attributes set.")
	}

	if *flagJoin {
		fmt.Printf("Searching for lobbies with bucket: %s\n", *flagBucketId)
		search, searchErr := p.Lobby().CreateLobbySearch(10)
		if searchErr != nil {
			return fmt.Errorf("create search: %w", searchErr)
		}
		defer search.Release()

		if paramErr := search.SetParameter("bucket", *flagBucketId, lobby.ComparisonEqual); paramErr != nil {
			return fmt.Errorf("set search parameter: %w", paramErr)
		}

		results, findErr := search.Find(ctx, productUserId)
		if findErr != nil {
			return fmt.Errorf("find lobbies: %w", findErr)
		}
		fmt.Printf("Found %d lobbies\n", len(results))

		if len(results) > 0 {
			details := results[0]
			info, infoErr := details.Info()
			if infoErr == nil {
				fmt.Printf("Joining lobby %s (owner: %s, %d/%d members)\n",
					info.LobbyId, info.LobbyOwnerUserId, info.MaxMembers-info.AvailableSlots, info.MaxMembers)
			}
			if joinErr := p.Lobby().JoinLobby(ctx, productUserId, details); joinErr != nil {
				return fmt.Errorf("join lobby: %w", joinErr)
			}
			fmt.Println("Joined lobby!")
			activeLobbyId = info.LobbyId
			for _, d := range results {
				d.Release()
			}
		} else {
			fmt.Println("No lobbies found.")
		}
	}

	fmt.Println("\nPress Ctrl+C to quit.")
	<-ctx.Done()
	fmt.Println("\nShutting down...")

	if activeLobbyId != "" {
		// Use a fresh context since ctx is already cancelled.
		leaveCtx, leaveCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer leaveCancel()
		if err := p.Lobby().LeaveLobby(leaveCtx, productUserId, activeLobbyId); err != nil {
			fmt.Printf("leave lobby: %v\n", err)
		}
	}
	return nil
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
	if err != nil {
		if connectResult != nil && connectResult.ContinuanceToken != 0 {
			fmt.Println("Creating product user...")
			userId, createErr := p.Connect().CreateUser(ctx, connectResult.ContinuanceToken)
			if createErr != nil {
				return "", fmt.Errorf("create user: %w", createErr)
			}
			fmt.Printf("Connect: ProductUserId=%s\n", *userId)
			return *userId, nil
		}
		return "", fmt.Errorf("connect login: %w", err)
	}
	fmt.Printf("Connect: ProductUserId=%s\n", connectResult.LocalUserId)
	return connectResult.LocalUserId, nil
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}
