package bootstrap

import (
	"context"
	"fmt"

	"clawbot-trust-lab/internal/clients/controlplane"
	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/config"
	"clawbot-trust-lab/internal/platform/store"
)

type Dependencies struct {
	ControlPlane controlplane.Client
	Memory       memory.Client
	Scenarios    store.ScenarioCatalog
}

func Build(cfg config.Config) Dependencies {
	return Dependencies{
		ControlPlane: controlplane.New(cfg.ControlPlaneURL, cfg.ControlPlaneTimeout),
		Memory:       memory.NewStub(cfg.MemoryURL),
		Scenarios:    store.NewInMemoryScenarioCatalog(),
	}
}

func Ready(ctx context.Context, deps Dependencies) error {
	if err := deps.ControlPlane.Health(ctx); err != nil {
		return fmt.Errorf("control-plane health check failed: %w", err)
	}
	return nil
}
