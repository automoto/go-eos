package platform

import (
	"context"
	"fmt"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/internal/threadworker"
)

type Platform struct {
	handle         cbinding.EOS_HPlatform
	worker         *threadworker.Worker
	notify         *callback.NotificationRegistry
	authHandle     cbinding.EOS_HAuth
	connectHandle  cbinding.EOS_HConnect
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

	p := &Platform{
		handle:         handle,
		notify:         callback.NewNotificationRegistry(),
		authHandle:     cbinding.EOS_Platform_GetAuthInterface(handle),
		connectHandle:  cbinding.EOS_Platform_GetConnectInterface(handle),
		lobbyHandle:    cbinding.EOS_Platform_GetLobbyInterface(handle),
		sessionsHandle: cbinding.EOS_Platform_GetSessionsInterface(handle),
		p2pHandle:      cbinding.EOS_Platform_GetP2PInterface(handle),
	}

	p.worker = threadworker.New(
		func() { cbinding.EOS_Platform_Tick(handle) },
		threadworker.WithTickInterval(cfg.tickInterval()),
	)
	p.worker.Start(context.Background())

	return p, nil
}

func (p *Platform) Shutdown() error {
	p.worker.Stop()
	cbinding.EOS_Platform_Release(p.handle)
	cbinding.EOS_Shutdown()
	return nil
}

func (p *Platform) Auth() cbinding.EOS_HAuth         { return p.authHandle }
func (p *Platform) Connect() cbinding.EOS_HConnect    { return p.connectHandle }
func (p *Platform) Lobby() cbinding.EOS_HLobby        { return p.lobbyHandle }
func (p *Platform) Sessions() cbinding.EOS_HSessions   { return p.sessionsHandle }
func (p *Platform) P2P() cbinding.EOS_HP2P             { return p.p2pHandle }
func (p *Platform) Worker() *threadworker.Worker       { return p.worker }
func (p *Platform) Notifications() *callback.NotificationRegistry { return p.notify }

func Run(ctx context.Context, cfg PlatformConfig, fn func(p *Platform) error) error {
	p, err := Initialize(cfg)
	if err != nil {
		return err
	}
	defer p.Shutdown()

	return fn(p)
}
