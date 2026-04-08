package platform

import (
	"errors"
	"time"
)

const defaultTickInterval = 16 * time.Millisecond

type PlatformConfig struct {
	ProductName    string
	ProductVersion string
	ProductId      string
	SandboxId      string
	DeploymentId   string
	ClientId       string
	ClientSecret   string
	TickInterval   time.Duration
}

func (c *PlatformConfig) validate() error {
	if c.ProductName == "" {
		return errors.New("ProductName is required")
	}
	if c.ProductVersion == "" {
		return errors.New("ProductVersion is required")
	}
	if c.ProductId == "" {
		return errors.New("ProductId is required")
	}
	if c.SandboxId == "" {
		return errors.New("SandboxId is required")
	}
	if c.DeploymentId == "" {
		return errors.New("DeploymentId is required")
	}
	if c.ClientId == "" {
		return errors.New("ClientId is required")
	}
	if c.ClientSecret == "" {
		return errors.New("ClientSecret is required")
	}
	return nil
}

func (c *PlatformConfig) tickInterval() time.Duration {
	if c.TickInterval <= 0 {
		return defaultTickInterval
	}
	return c.TickInterval
}
