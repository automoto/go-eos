package platform

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mydev/go-eos/eos/auth"
	"github.com/mydev/go-eos/eos/connect"
	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
)

type Platform struct {
	handle         cbinding.EOS_HPlatform
	worker         *threadworker.Worker
	notify         *callback.NotificationRegistry
	auth           *auth.Auth
	connect        *connect.Connect
	lobbyHandle    cbinding.EOS_HLobby
	sessionsHandle cbinding.EOS_HSessions
	p2pHandle      cbinding.EOS_HP2P
}

func Initialize(cfg PlatformConfig) (*Platform, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	initLogging()

	result := cbinding.EOS_Initialize(&cbinding.EOS_InitializeOptions{
		ProductName:    cfg.ProductName,
		ProductVersion: cfg.ProductVersion,
	})
	if result != cbinding.EOS_EResult_Success {
		return nil, fmt.Errorf("EOS_Initialize failed: %d", result)
	}

	handle := cbinding.EOS_Platform_Create(&cbinding.EOS_Platform_Options{
		ProductId:    cfg.ProductId,
		SandboxId:    cfg.SandboxId,
		DeploymentId: cfg.DeploymentId,
		ClientId:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
	})
	if handle == 0 {
		return nil, fmt.Errorf("EOS_Platform_Create returned null handle")
	}

	worker := threadworker.New(
		func() { cbinding.EOS_Platform_Tick(handle) },
		threadworker.WithTickInterval(cfg.tickInterval()),
	)
	worker.Start(context.Background())

	p := &Platform{
		handle:         handle,
		worker:         worker,
		notify:         callback.NewNotificationRegistry(),
		auth:           auth.New(cbinding.EOS_Platform_GetAuthInterface(handle), worker),
		connect:        connect.New(cbinding.EOS_Platform_GetConnectInterface(handle), worker),
		lobbyHandle:    cbinding.EOS_Platform_GetLobbyInterface(handle),
		sessionsHandle: cbinding.EOS_Platform_GetSessionsInterface(handle),
		p2pHandle:      cbinding.EOS_Platform_GetP2PInterface(handle),
	}

	return p, nil
}

func (p *Platform) Shutdown() error {
	p.worker.Stop()
	cbinding.EOS_Platform_Release(p.handle)
	cbinding.EOS_Shutdown()
	return nil
}

func (p *Platform) Auth() *auth.Auth            { return p.auth }
func (p *Platform) Connect() *connect.Connect   { return p.connect }
func (p *Platform) Lobby() cbinding.EOS_HLobby  { return p.lobbyHandle }
func (p *Platform) Sessions() cbinding.EOS_HSessions { return p.sessionsHandle }
func (p *Platform) P2P() cbinding.EOS_HP2P      { return p.p2pHandle }
func (p *Platform) Worker() *threadworker.Worker       { return p.worker }
func (p *Platform) Notifications() *callback.NotificationRegistry { return p.notify }

func Run(ctx context.Context, cfg PlatformConfig, fn func(p *Platform) error) error {
	p, err := Initialize(cfg)
	if err != nil {
		return err
	}
	defer func() {
		if shutdownErr := p.Shutdown(); shutdownErr != nil {
			slog.Error("platform shutdown failed", "error", shutdownErr)
		}
	}()

	return fn(p)
}
