// Command ebitengine-lobby demonstrates integrating go-eos with Ebitengine
// for lobby-based matchmaking with a graphical display.
//
// This example shows the recommended integration pattern:
//   - EOS platform runs on a background worker via platform.Initialize()
//   - Ebitengine runs on the main thread via ebiten.RunGame()
//   - EOS operations are called from the game's Update() method
//
// Usage:
//
//	cd examples/ebitengine-lobby
//	go run . --create --bucket-id "game:mode=deathmatch"
//	go run . --join --bucket-id "game:mode=deathmatch"
//
// Environment variables:
//
//	EOS_PRODUCT_ID, EOS_SANDBOX_ID, EOS_DEPLOYMENT_ID
//	EOS_CLIENT_ID, EOS_CLIENT_SECRET
//	EOS_DEV_AUTH_HOST, EOS_DEV_AUTH_CREDENTIAL
package main

import (
	"context"
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/mydev/go-eos/eos/auth"
	"github.com/mydev/go-eos/eos/connect"
	"github.com/mydev/go-eos/eos/lobby"
	"github.com/mydev/go-eos/eos/platform"
	"github.com/mydev/go-eos/eos/types"
)

var (
	flagCreate   = flag.Bool("create", false, "Create a new lobby")
	flagJoin     = flag.Bool("join", false, "Search for and join a lobby")
	flagBucketId = flag.String("bucket-id", "game:mode=default", "Lobby bucket ID")
)

func main() {
	flag.Parse()

	if !*flagCreate && !*flagJoin {
		fmt.Println("Usage: ebitengine-lobby --create|--join --bucket-id <bucket>")
		os.Exit(1)
	}

	cfg := platform.PlatformConfig{
		ProductName:    "ebitengine-lobby-demo",
		ProductVersion: "1.0.0",
		ProductId:      requireEnv("EOS_PRODUCT_ID"),
		SandboxId:      requireEnv("EOS_SANDBOX_ID"),
		DeploymentId:   requireEnv("EOS_DEPLOYMENT_ID"),
		ClientId:       requireEnv("EOS_CLIENT_ID"),
		ClientSecret:   requireEnv("EOS_CLIENT_SECRET"),
	}

	// Initialize EOS on a background worker thread.
	// Ebitengine owns the main thread for its window/event loop.
	p, err := platform.Initialize(cfg)
	if err != nil {
		log.Fatalf("platform init: %v", err)
	}

	game := &Game{
		platform: p,
		status:   "Initializing...",
	}

	// Start the EOS login and lobby flow in the background.
	go game.eosFlow()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("go-eos + Ebitengine — Lobby Demo")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	p.Shutdown()
}

// Game implements ebiten.Game and holds the EOS state.
type Game struct {
	platform      *platform.Platform
	productUserId types.ProductUserId
	lobbyId       string
	members       []string
	status        string
	events        []string

	mu sync.Mutex
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 40, A: 255})

	g.mu.Lock()
	status := g.status
	lobbyId := g.lobbyId
	members := append([]string{}, g.members...)
	events := append([]string{}, g.events...)
	g.mu.Unlock()

	var lines []string
	lines = append(lines, fmt.Sprintf("Status: %s", status))
	lines = append(lines, "")

	if lobbyId != "" {
		lines = append(lines, fmt.Sprintf("Lobby: %s", lobbyId))
		lines = append(lines, fmt.Sprintf("Members (%d):", len(members)))
		for _, m := range members {
			lines = append(lines, fmt.Sprintf("  - %s", m))
		}
	}

	if len(events) > 0 {
		lines = append(lines, "")
		lines = append(lines, "Events:")
		// Show last 10 events
		start := 0
		if len(events) > 10 {
			start = len(events) - 10
		}
		for _, e := range events[start:] {
			lines = append(lines, fmt.Sprintf("  %s", e))
		}
	}

	ebitenutil.DebugPrint(screen, strings.Join(lines, "\n"))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func (g *Game) setStatus(s string) {
	g.mu.Lock()
	g.status = s
	g.mu.Unlock()
}

func (g *Game) addEvent(msg string) {
	g.mu.Lock()
	g.events = append(g.events, msg)
	g.mu.Unlock()
}

func (g *Game) setLobby(id string, members []string) {
	g.mu.Lock()
	g.lobbyId = id
	g.members = members
	g.mu.Unlock()
}

// eosFlow runs the EOS login and lobby operations on a background goroutine.
func (g *Game) eosFlow() {
	ctx := context.Background()
	p := g.platform

	// Auth login
	g.setStatus("Logging in (Auth)...")
	devHost := requireEnv("EOS_DEV_AUTH_HOST")
	devCred := requireEnv("EOS_DEV_AUTH_CREDENTIAL")

	loginResult, err := p.Auth().Login(ctx, auth.LoginOptions{
		CredentialType: types.LoginCredentialDeveloper,
		ID:             devHost,
		Token:          devCred,
	})
	if err != nil {
		g.setStatus(fmt.Sprintf("Auth login failed: %v", err))
		return
	}
	g.addEvent(fmt.Sprintf("Auth OK: %s", loginResult.LocalUserId))

	// Connect login
	g.setStatus("Logging in (Connect)...")
	idToken, err := p.Auth().CopyIdToken(loginResult.LocalUserId)
	if err != nil {
		g.setStatus(fmt.Sprintf("CopyIdToken failed: %v", err))
		return
	}

	connectResult, err := p.Connect().Login(ctx, connect.LoginOptions{
		CredentialType: types.ExternalCredentialEpicIDToken,
		Token:          idToken,
	})
	if err != nil {
		if connectResult != nil && connectResult.ContinuanceToken != 0 {
			g.addEvent("Creating product user...")
			userId, createErr := p.Connect().CreateUser(ctx, connectResult.ContinuanceToken)
			if createErr != nil {
				g.setStatus(fmt.Sprintf("CreateUser failed: %v", createErr))
				return
			}
			g.productUserId = *userId
		} else {
			g.setStatus(fmt.Sprintf("Connect login failed: %v", err))
			return
		}
	} else {
		g.productUserId = connectResult.LocalUserId
	}
	g.addEvent(fmt.Sprintf("Connect OK: %s", g.productUserId))

	// Register for lobby events
	p.Lobby().AddNotifyLobbyMemberStatusReceived(func(info lobby.MemberStatusInfo) {
		g.addEvent(fmt.Sprintf("Member %s status: %d", info.TargetUserId, info.CurrentStatus))
	})
	p.Lobby().AddNotifyLobbyUpdateReceived(func(info lobby.LobbyUpdateInfo) {
		g.addEvent(fmt.Sprintf("Lobby updated: %s", info.LobbyId))
	})

	if *flagCreate {
		g.createLobby(ctx)
	} else {
		g.joinLobby(ctx)
	}
}

func (g *Game) createLobby(ctx context.Context) {
	g.setStatus("Creating lobby...")

	lobbyId, err := g.platform.Lobby().CreateLobby(ctx, lobby.CreateLobbyOptions{
		LocalUserId:     g.productUserId,
		MaxMembers:      4,
		PermissionLevel: lobby.PermissionPublicAdvertised,
		AllowInvites:    true,
		BucketId:        *flagBucketId,
	})
	if err != nil {
		g.setStatus(fmt.Sprintf("Create lobby failed: %v", err))
		return
	}

	g.addEvent(fmt.Sprintf("Lobby created: %s", lobbyId))
	g.setLobby(lobbyId, []string{string(g.productUserId)})
	g.setStatus("In lobby (owner) — waiting for players")
}

func (g *Game) joinLobby(ctx context.Context) {
	g.setStatus(fmt.Sprintf("Searching: %s", *flagBucketId))

	search, err := g.platform.Lobby().CreateLobbySearch(10)
	if err != nil {
		g.setStatus(fmt.Sprintf("Create search failed: %v", err))
		return
	}
	defer search.Release()

	if err := search.SetParameter("bucket", *flagBucketId, lobby.ComparisonEqual); err != nil {
		g.setStatus(fmt.Sprintf("Set parameter failed: %v", err))
		return
	}

	results, err := search.Find(ctx, g.productUserId)
	if err != nil {
		g.setStatus(fmt.Sprintf("Search failed: %v", err))
		return
	}

	if len(results) == 0 {
		g.setStatus("No lobbies found")
		return
	}

	g.addEvent(fmt.Sprintf("Found %d lobbies", len(results)))
	details := results[0]
	info, infoErr := details.Info()

	if err := g.platform.Lobby().JoinLobby(ctx, g.productUserId, details); err != nil {
		g.setStatus(fmt.Sprintf("Join failed: %v", err))
		for _, d := range results {
			d.Release()
		}
		return
	}

	var lobbyId string
	var memberStrs []string
	if infoErr == nil {
		lobbyId = info.LobbyId
		memberStrs = []string{string(info.LobbyOwnerUserId), string(g.productUserId)}
	}

	for _, d := range results {
		d.Release()
	}

	g.addEvent("Joined lobby!")
	g.setLobby(lobbyId, memberStrs)
	g.setStatus("In lobby — connected")
}

func requireEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return val
}
