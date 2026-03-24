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
	servicecommerce "clawbot-trust-lab/internal/services/commerce"
	serviceevents "clawbot-trust-lab/internal/services/events"
	servicescenario "clawbot-trust-lab/internal/services/scenario"
	servicetrust "clawbot-trust-lab/internal/services/trust"
)

type Dependencies struct {
	ControlPlane controlplane.Client
	Memory       memory.Client
	Scenarios    *scenario.Service
	Trust        *trust.Service
	Replay       *replay.Service
	Benchmark    *benchmark.Service
	Commerce     *servicecommerce.Service
	Events       *serviceevents.Service
	TrustFlow    *servicetrust.Service
	Execution    *servicescenario.Service
}

func Build(cfg config.Config) (Dependencies, error) {
	controlPlaneClient := controlplane.New(cfg.ControlPlaneURL, cfg.ControlPlaneTimeout)
	memoryClient := memory.New(cfg.ClawMemBaseURL, cfg.ClawMemTimeout)
	scenarioLoader := loader.New(cfg.ScenarioPacksDir)
	scenarioService, err := scenario.NewService(scenarioLoader)
	if err != nil {
		return Dependencies{}, err
	}
	replayStore, err := store.NewFileReplayStore(cfg.ReplayArchiveDir)
	if err != nil {
		return Dependencies{}, err
	}
	worldStore := store.NewCommerceWorldStore()
	commerceService := servicecommerce.NewService(worldStore)
	eventService := serviceevents.NewService(worldStore)
	trustFlowService := servicetrust.NewService(worldStore)
	trustService := trust.NewService(scenarioService, store.NewInMemoryTrustArtifactStore(), memoryClient)
	replayService := replay.NewService(replayStore, memoryClient)
	executionService := servicescenario.NewService(scenarioService, commerceService, eventService, trustFlowService, trustService, replayService)

	return Dependencies{
		ControlPlane: controlPlaneClient,
		Memory:       memoryClient,
		Scenarios:    scenarioService,
		Trust:        trustService,
		Replay:       replayService,
		Benchmark:    benchmark.NewService(controlPlaneClient),
		Commerce:     commerceService,
		Events:       eventService,
		TrustFlow:    trustFlowService,
		Execution:    executionService,
	}, nil
}

func Ready(ctx context.Context, deps Dependencies) error {
	if err := deps.ControlPlane.Health(ctx); err != nil {
		return fmt.Errorf("control-plane health check failed: %w", err)
	}
	if err := deps.Memory.Health(ctx); err != nil {
		return fmt.Errorf("clawmem health check failed: %w", err)
	}
	return nil
}
