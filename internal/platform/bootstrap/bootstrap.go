package bootstrap

import (
	"context"
	"fmt"
	"log/slog"

	"clawbot-trust-lab/internal/clients/controlplane"
	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/config"
	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/platform/loader"
	"clawbot-trust-lab/internal/platform/store"
	servicebenchmark "clawbot-trust-lab/internal/services/benchmark"
	servicecommerce "clawbot-trust-lab/internal/services/commerce"
	servicedetection "clawbot-trust-lab/internal/services/detection"
	serviceevents "clawbot-trust-lab/internal/services/events"
	serviceoperator "clawbot-trust-lab/internal/services/operator"
	servicereporting "clawbot-trust-lab/internal/services/reporting"
	servicescenario "clawbot-trust-lab/internal/services/scenario"
	servicetrust "clawbot-trust-lab/internal/services/trust"
)

type Dependencies struct {
	ControlPlane controlplane.Client
	Memory       memory.Client
	Scenarios    *scenario.Service
	Trust        *trust.Service
	Replay       *replay.Service
	Benchmark    *servicebenchmark.Service
	Commerce     *servicecommerce.Service
	Events       *serviceevents.Service
	TrustFlow    *servicetrust.Service
	Execution    *servicescenario.Service
	Detection    *servicedetection.Service
	Operator     *serviceoperator.Service
}

func Build(cfg config.Config, logger *slog.Logger) (Dependencies, error) {
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
	detectionStore := store.NewDetectionStore()
	benchmarkStore := store.NewBenchmarkStore()
	operatorStore := store.NewOperatorStore()
	historicalState := LoadHistoricalState(cfg.ReportsDir, logger)
	for _, round := range historicalState.Rounds {
		benchmarkStore.PutHistorical(round)
	}
	for _, item := range historicalState.DetectionResults {
		detectionStore.Put(item)
	}
	if logger != nil && len(historicalState.Rounds) > 0 {
		logger.Info("bootstrapped historical benchmark rounds from reports", "round_count", len(historicalState.Rounds), "reports_dir", cfg.ReportsDir)
	}
	commerceService := servicecommerce.NewService(worldStore)
	eventService := serviceevents.NewService(worldStore)
	trustFlowService := servicetrust.NewService(worldStore)
	trustService := trust.NewService(scenarioService, store.NewInMemoryTrustArtifactStore(), memoryClient)
	replayService := replay.NewService(replayStore, memoryClient)
	executionService := servicescenario.NewService(scenarioService, commerceService, eventService, trustFlowService, trustService, replayService)
	detectionService := servicedetection.NewService(worldStore, executionService, replayService, memoryClient, detectionStore)
	reportingService := servicereporting.NewService(cfg.ReportsDir)
	benchmarkRegistrationService := benchmark.NewService(controlPlaneClient)
	benchmarkRoundService := servicebenchmark.NewService(benchmarkRegistrationService, executionService, detectionService, replayService, benchmarkStore, reportingService)
	operatorService := serviceoperator.NewService(benchmarkRoundService, detectionService, operatorStore)

	return Dependencies{
		ControlPlane: controlPlaneClient,
		Memory:       memoryClient,
		Scenarios:    scenarioService,
		Trust:        trustService,
		Replay:       replayService,
		Benchmark:    benchmarkRoundService,
		Commerce:     commerceService,
		Events:       eventService,
		TrustFlow:    trustFlowService,
		Execution:    executionService,
		Detection:    detectionService,
		Operator:     operatorService,
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
