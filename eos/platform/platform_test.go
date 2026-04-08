//go:build eosstub

package platform

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mydev/go-eos/eos/internal/cbinding"
)

func validConfig() PlatformConfig {
	return PlatformConfig{
		ProductName:    "test-game",
		ProductVersion: "1.0.0",
		ProductId:      "prod-123",
		SandboxId:      "sandbox-456",
		DeploymentId:   "deploy-789",
		ClientId:       "client-abc",
		ClientSecret:   "secret-xyz",
		TickInterval:   1 * time.Millisecond,
	}
}

func Test_initialize_and_shutdown_should_succeed(t *testing.T) {
	p, err := Initialize(validConfig())
	assert.NoError(t, err)

	err = p.Shutdown()
	assert.NoError(t, err)
}

func Test_initialize_should_fail_with_missing_required_fields(t *testing.T) {
	tests := []struct {
		name  string
		field string
	}{
		{"empty ProductName", "ProductName"},
		{"empty ProductId", "ProductId"},
		{"empty SandboxId", "SandboxId"},
		{"empty DeploymentId", "DeploymentId"},
		{"empty ClientId", "ClientId"},
		{"empty ClientSecret", "ClientSecret"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := validConfig()
			switch tt.field {
			case "ProductName":
				cfg.ProductName = ""
			case "ProductId":
				cfg.ProductId = ""
			case "SandboxId":
				cfg.SandboxId = ""
			case "DeploymentId":
				cfg.DeploymentId = ""
			case "ClientId":
				cfg.ClientId = ""
			case "ClientSecret":
				cfg.ClientSecret = ""
			}

			_, err := Initialize(cfg)
			assert.Error(t, err)
		})
	}
}

func Test_run_should_complete_successfully(t *testing.T) {
	err := Run(context.Background(), validConfig(), func(p *Platform) error {
		return nil
	})
	assert.NoError(t, err)
}

func Test_run_should_propagate_function_error(t *testing.T) {
	expected := errors.New("game error")
	err := Run(context.Background(), validConfig(), func(p *Platform) error {
		return expected
	})
	assert.ErrorIs(t, err, expected)
}

func Test_run_should_shutdown_on_context_cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	err := Run(ctx, validConfig(), func(p *Platform) error {
		cancel()
		time.Sleep(10 * time.Millisecond)
		return ctx.Err()
	})
	assert.Error(t, err)
}

func Test_interface_accessors_should_return_non_nil(t *testing.T) {
	p, err := Initialize(validConfig())
	assert.NoError(t, err)
	defer p.Shutdown()

	assert.NotNil(t, p.Auth())
	assert.NotNil(t, p.Connect())
	assert.NotNil(t, p.Lobby())
	assert.NotNil(t, p.Sessions())
	assert.NotEqual(t, cbinding.EOS_HP2P(0), p.P2P())
}

func Test_logging_should_forward_to_slog(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	slog.SetDefault(slog.New(handler))
	defer slog.SetDefault(slog.Default())

	p, err := Initialize(validConfig())
	assert.NoError(t, err)
	defer p.Shutdown()

	cbinding.SimulateLogMessage(&cbinding.EOS_LogMessage{
		Category: "Core",
		Level:    cbinding.EOS_LOG_Warning,
		Message:  "test warning message",
	})

	assert.Contains(t, buf.String(), "test warning message")
}
