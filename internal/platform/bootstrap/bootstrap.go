package bootstrap

import (
	"context"
	"fmt"

	"clawbot-trust-lab/internal/clients/controlplane"
	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/config"
	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/platform/loader"
	"clawbot-trust-lab/internal/platform/store"
)

type Dependencies struct {
	ControlPlane controlplane.Client
	Memory       memory.Client
	Scenarios    *scenario.Service
	Trust        *trust.Service
	Replay       *replay.Service
	Benchmark    *benchmark.Service
}

func Build(cfg config.Config) (Dependencies, error) {
	controlPlaneClient := controlplane.New(cfg.ControlPlaneURL, cfg.ControlPlaneTimeout)
	memoryClient := memory.NewStub(cfg.MemoryURL)
	scenarioLoader := loader.New(cfg.ScenarioPacksDir)
	scenarioService, err := scenario.NewService(scenarioLoader)
	if err != nil {
		return Dependencies{}, err
	}
	replayStore, err := store.NewFileReplayStore(cfg.ReplayArchiveDir)
	if err != nil {
		return Dependencies{}, err
	}

	return Dependencies{
		ControlPlane: controlPlaneClient,
		Memory:       memoryClient,
		Scenarios:    scenarioService,
		Trust:        trust.NewService(scenarioService, store.NewInMemoryTrustArtifactStore(), memoryClient),
		Replay:       replay.NewService(replayStore),
		Benchmark:    benchmark.NewService(controlPlaneClient),
	}, nil
}

func Ready(ctx context.Context, deps Dependencies) error {
	if err := deps.ControlPlane.Health(ctx); err != nil {
		return fmt.Errorf("control-plane health check failed: %w", err)
	}
	return nil
}
