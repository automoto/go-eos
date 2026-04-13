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
	"github.com/mydev/go-eos/eos/lobby"
	"github.com/mydev/go-eos/eos/p2p"
	"github.com/mydev/go-eos/eos/sessions"
)

// Platform wraps EOS_HPlatform and owns the SDK lifecycle, tick loop, and interface accessors.
type Platform struct {
	handle   cbinding.EOS_HPlatform
	worker   *threadworker.Worker
	notify   *callback.NotificationRegistry
	auth     *auth.Auth
	connect  *connect.Connect
	lobby    *lobby.Lobby
	sessions *sessions.Sessions
	p2p      *p2p.P2P
}

// Initialize creates and initializes the EOS platform. See EOS_Platform_Create.
func Initialize(cfg PlatformConfig) (*Platform, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Start the worker first so all SDK calls run on its LockOSThread goroutine.
	// The tick function no-ops until handle is set.
	p := &Platform{
		notify: callback.NewNotificationRegistry(),
	}

	worker := threadworker.New(
		func() {
			if p.handle != 0 {
				cbinding.EOS_Platform_Tick(p.handle)
			}
		},
		threadworker.WithTickInterval(cfg.tickInterval()),
	)
	worker.Start(context.Background())
	p.worker = worker

	// All EOS SDK calls must execute on the worker's locked OS thread (THR-1).
	var initErr error
	if err := worker.Submit(func() {
		result := cbinding.EOS_Initialize(&cbinding.EOS_InitializeOptions{
			ProductName:    cfg.ProductName,
			ProductVersion: cfg.ProductVersion,
		})
		if result != cbinding.EOS_EResult_Success {
			initErr = fmt.Errorf("EOS_Initialize failed: %d", result)
			return
		}

		initLogging()

		handle := cbinding.EOS_Platform_Create(&cbinding.EOS_Platform_Options{
			ProductId:    cfg.ProductId,
			SandboxId:    cfg.SandboxId,
			DeploymentId: cfg.DeploymentId,
			ClientId:     cfg.ClientId,
			ClientSecret: cfg.ClientSecret,
		})
		if handle == 0 {
			cbinding.EOS_Shutdown()
			initErr = fmt.Errorf("EOS_Platform_Create returned null handle")
			return
		}

		p.handle = handle
		p.auth = auth.New(cbinding.EOS_Platform_GetAuthInterface(handle), worker)
		p.connect = connect.New(cbinding.EOS_Platform_GetConnectInterface(handle), worker)
		p.lobby = lobby.New(cbinding.EOS_Platform_GetLobbyInterface(handle), worker)
		p.sessions = sessions.New(cbinding.EOS_Platform_GetSessionsInterface(handle), worker)
		p.p2p = p2p.New(cbinding.EOS_Platform_GetP2PInterface(handle), worker)
	}); err != nil {
		worker.Stop()
		return nil, fmt.Errorf("worker submit failed: %w", err)
	}
	if initErr != nil {
		worker.Stop()
		return nil, initErr
	}

	return p, nil
}

// Shutdown releases the EOS platform and stops the tick loop. See EOS_Platform_Release.
func (p *Platform) Shutdown() error {
	// Release on the worker's locked OS thread, then zero the handle
	// so the tick function no-ops before the worker fully stops.
	_ = p.worker.Submit(func() {
		cbinding.EOS_Platform_Release(p.handle)
		p.handle = 0
		cbinding.EOS_Shutdown()
	})
	p.worker.Stop()
	return nil
}

// Auth returns the Auth interface wrapper.
func (p *Platform) Auth() *auth.Auth { return p.auth }

// Connect returns the Connect interface wrapper.
func (p *Platform) Connect() *connect.Connect { return p.connect }

// Lobby returns the Lobby interface wrapper.
func (p *Platform) Lobby() *lobby.Lobby { return p.lobby }

// Sessions returns the Sessions interface wrapper.
func (p *Platform) Sessions() *sessions.Sessions { return p.sessions }

// P2P returns the P2P interface wrapper.
func (p *Platform) P2P() *p2p.P2P { return p.p2p }

// Worker returns the platform's thread worker for scheduling SDK calls.
func (p *Platform) Worker() *threadworker.Worker { return p.worker }

// Notifications returns the platform's notification registry.
func (p *Platform) Notifications() *callback.NotificationRegistry { return p.notify }

// Run initializes the platform, calls fn, and shuts down when fn returns.
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

// RunOnMainThread initializes the SDK and drives the tick loop on the
// calling goroutine. On macOS the EOS SDK's HTTP layer dispatches
// completions through the main thread's run loop, so all SDK work must
// happen there. The caller MUST have called runtime.LockOSThread() in
// an init() function to guarantee it owns the main OS thread.
//
// fn runs on a separate goroutine; when it returns the platform shuts down.
func RunOnMainThread(ctx context.Context, cfg PlatformConfig, fn func(p *Platform) error) error {
	if err := cfg.validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// The caller holds the main OS thread (via init() + runtime.LockOSThread).
	// All SDK calls below run directly on this goroutine — no Submit needed.
	result := cbinding.EOS_Initialize(&cbinding.EOS_InitializeOptions{
		ProductName:    cfg.ProductName,
		ProductVersion: cfg.ProductVersion,
	})
	if result != cbinding.EOS_EResult_Success {
		return fmt.Errorf("EOS_Initialize failed: %d", result)
	}

	initLogging()

	handle := cbinding.EOS_Platform_Create(&cbinding.EOS_Platform_Options{
		ProductId:    cfg.ProductId,
		SandboxId:    cfg.SandboxId,
		DeploymentId: cfg.DeploymentId,
		ClientId:     cfg.ClientId,
		ClientSecret: cfg.ClientSecret,
	})
	if handle == 0 {
		cbinding.EOS_Shutdown()
		return fmt.Errorf("EOS_Platform_Create returned null handle")
	}

	p := &Platform{
		notify: callback.NewNotificationRegistry(),
	}

	worker := threadworker.New(
		func() {
			if p.handle != 0 {
				cbinding.EOS_Platform_Tick(p.handle)
			}
		},
		threadworker.WithTickInterval(cfg.tickInterval()),
	)
	p.worker = worker
	p.handle = handle
	p.auth = auth.New(cbinding.EOS_Platform_GetAuthInterface(handle), worker)
	p.connect = connect.New(cbinding.EOS_Platform_GetConnectInterface(handle), worker)
	p.lobby = lobby.New(cbinding.EOS_Platform_GetLobbyInterface(handle), worker)
	p.sessions = sessions.New(cbinding.EOS_Platform_GetSessionsInterface(handle), worker)
	p.p2p = p2p.New(cbinding.EOS_Platform_GetP2PInterface(handle), worker)

	// Run fn on a separate goroutine. When it finishes, release the SDK
	// (via Submit so it runs on this same main thread) then stop the loop.
	// Use Background() so the worker keeps ticking while fn does cleanup
	// after ctx cancellation (e.g. LeaveLobby on Ctrl+C).
	workerCtx, workerCancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)
	go func() {
		fnErr := fn(p)
		_ = worker.Submit(func() {
			cbinding.EOS_Platform_Release(p.handle)
			p.handle = 0
			cbinding.EOS_Shutdown()
		})
		workerCancel()
		errCh <- fnErr
	}()

	// Drive the tick loop on the calling (main) goroutine — this blocks
	// until workerCancel() is called above. Init, tick, and shutdown all
	// happen on this same OS thread, satisfying THR-1.
	worker.StartBlocking(workerCtx)

	return <-errCh
}
