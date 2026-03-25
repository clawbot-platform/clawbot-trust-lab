package benchmark

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"clawbot-trust-lab/internal/clients/memory"
	domainbenchmark "clawbot-trust-lab/internal/domain/benchmark"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	domainreplay "clawbot-trust-lab/internal/domain/replay"
	domainscenario "clawbot-trust-lab/internal/domain/scenario"
	domaintrust "clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/platform/store"
	servicecommerce "clawbot-trust-lab/internal/services/commerce"
	servicedetection "clawbot-trust-lab/internal/services/detection"
	serviceevents "clawbot-trust-lab/internal/services/events"
	servicereporting "clawbot-trust-lab/internal/services/reporting"
	servicescenario "clawbot-trust-lab/internal/services/scenario"
	servicetrust "clawbot-trust-lab/internal/services/trust"
)

type registrarStub struct{}

func (registrarStub) RegisterRound(context.Context, domainbenchmark.RegistrationRequest) (domainbenchmark.RegistrationResult, error) {
	return domainbenchmark.RegistrationResult{RegistrationID: "bench-reg-1", Status: "accepted"}, nil
}

func (registrarStub) Status() map[string]any {
	return map[string]any{"registrations": 0, "last_status": "idle"}
}

type scenarioCatalogStub struct {
	items map[string]domainscenario.Scenario
}

func (s scenarioCatalogStub) ListScenarios() []domainscenario.Scenario {
	out := make([]domainscenario.Scenario, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	return out
}

func (s scenarioCatalogStub) GetScenario(id string) (domainscenario.Scenario, error) {
	return s.items[id], nil
}

type memoryClientStub struct{}

func (memoryClientStub) Health(context.Context) error { return nil }
func (memoryClientStub) StoreReplayCase(context.Context, memory.StoreReplayCaseRequest) error {
	return nil
}
func (memoryClientStub) FetchSimilarCases(context.Context, memory.FetchSimilarCasesRequest) (memory.FetchSimilarCasesResponse, error) {
	return memory.FetchSimilarCasesResponse{}, nil
}
func (memoryClientStub) StoreTrustArtifact(context.Context, memory.StoreTrustArtifactRequest) error {
	return nil
}
func (memoryClientStub) LoadScenarioContext(_ context.Context, request memory.LoadScenarioContextRequest) (memory.LoadScenarioContextResponse, error) {
	count := 1
	if request.ScenarioID == "commerce-challenger-weakened-provenance-purchase" {
		count = 2
	}
	return memory.LoadScenarioContextResponse{
		ScenarioID: request.ScenarioID,
		Context:    map[string]any{"record_count": count},
	}, nil
}

func newRoundService(t *testing.T) *Service {
	t.Helper()

	world := store.NewCommerceWorldStore()
	replayStore, err := store.NewFileReplayStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewFileReplayStore() error = %v", err)
	}
	scenarioService := servicescenario.NewService(
		scenarioCatalogStub{items: map[string]domainscenario.Scenario{
			"commerce-clean-agent-assisted-purchase":           {ID: "commerce-clean-agent-assisted-purchase", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-suspicious-refund-attempt":               {ID: "commerce-suspicious-refund-attempt", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-challenger-weakened-provenance-purchase": {ID: "commerce-challenger-weakened-provenance-purchase", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-challenger-expired-mandate-purchase":     {ID: "commerce-challenger-expired-mandate-purchase", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-challenger-approval-removed-refund":      {ID: "commerce-challenger-approval-removed-refund", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
		}},
		servicecommerce.NewService(world),
		serviceevents.NewService(world),
		servicetrust.NewService(world),
		domaintrust.NewService(scenarioCatalogForArtifacts{}, store.NewInMemoryTrustArtifactStore(), memoryClientStub{}),
		domainreplay.NewService(replayStore, memoryClientStub{}),
	)

	replayService := domainreplay.NewService(replayStore, memoryClientStub{})
	detectionService := servicedetection.NewService(world, scenarioService, replayService, memoryClientStub{}, store.NewDetectionStore())
	reportingService := servicereporting.NewService(t.TempDir())
	service := NewService(registrarStub{}, scenarioService, detectionService, replayService, store.NewBenchmarkStore(), reportingService)
	service.now = func() time.Time { return time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC) }
	return service
}

type scenarioCatalogForArtifacts struct{}

func (scenarioCatalogForArtifacts) ListScenarios() []domainscenario.Scenario { return nil }
func (scenarioCatalogForArtifacts) GetScenario(id string) (domainscenario.Scenario, error) {
	return domainscenario.Scenario{ID: id, Type: domainscenario.ScenarioTypeCommercePurchase}, nil
}

func TestRunRoundCreatesReportsAndPromotion(t *testing.T) {
	service := newRoundService(t)

	round, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() error = %v", err)
	}

	if round.RoundStatus != domainbenchmark.RoundStatusCompleted {
		t.Fatalf("expected completed round, got %s", round.RoundStatus)
	}
	if len(round.StableScenarioRefs) != 2 {
		t.Fatalf("expected 2 stable scenarios, got %d", len(round.StableScenarioRefs))
	}
	if len(round.ChallengerVariantRefs) != 3 {
		t.Fatalf("expected 3 challenger variants, got %d", len(round.ChallengerVariantRefs))
	}
	if len(round.PromotionResults) == 0 {
		t.Fatal("expected at least one promotion result")
	}
	if round.Summary.RobustnessOutcome != domainbenchmark.RobustnessOutcomeNewBlindSpotDiscovered {
		t.Fatalf("unexpected robustness outcome: %s", round.Summary.RobustnessOutcome)
	}
	if len(round.Reports.Artifacts) != 5 {
		t.Fatalf("expected 5 report artifacts, got %d", len(round.Reports.Artifacts))
	}
	for _, artifact := range round.Reports.Artifacts {
		if _, err := os.Stat(artifact.Path); err != nil {
			t.Fatalf("expected artifact %s to exist: %v", artifact.Path, err)
		}
	}
}

func TestRunRoundRetestsPriorPromotions(t *testing.T) {
	service := newRoundService(t)

	firstRound, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() first error = %v", err)
	}
	if len(firstRound.PromotionResults) == 0 {
		t.Fatal("expected initial round to produce a promotion")
	}

	service.now = func() time.Time { return time.Date(2026, 3, 25, 12, 10, 0, 0, time.UTC) }
	secondRound, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() second error = %v", err)
	}

	if secondRound.Summary.ReplayRetestCount == 0 {
		t.Fatal("expected replay retests in second round")
	}
	if secondRound.Summary.ReplayPassRate >= 1 {
		t.Fatalf("expected replay pass rate below 1, got %.2f", secondRound.Summary.ReplayPassRate)
	}
	if len(secondRound.Delta) == 0 {
		t.Fatal("expected detection delta in second round")
	}
}

func TestRunRoundReportsEndpointDataIsStored(t *testing.T) {
	service := newRoundService(t)

	round, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() error = %v", err)
	}

	stored, err := service.GetRoundReports(round.ID)
	if err != nil {
		t.Fatalf("GetRoundReports() error = %v", err)
	}
	if stored.Directory == "" {
		t.Fatal("expected report directory")
	}
	if filepath.Base(stored.Directory) != round.ID {
		t.Fatalf("unexpected report directory %s", stored.Directory)
	}
}

func TestWeakenedProvenanceVariantPromotesOnDetectorMiss(t *testing.T) {
	service := newRoundService(t)

	round, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() error = %v", err)
	}

	found := false
	for _, item := range round.PromotionResults {
		if item.ChallengerVariantID == "variant-weakened-provenance" && item.PromotionReason == domainbenchmark.PromotionReasonDetectorMiss {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected weakened provenance promotion in %#v", round.PromotionResults)
	}
}

func TestStatusRankOrdering(t *testing.T) {
	if !meetsMinimumStatus(detectionmodel.DetectionStatusStepUpRequired, detectionmodel.DetectionStatusSuspicious) {
		t.Fatal("expected step-up to satisfy suspicious floor")
	}
	if meetsMinimumStatus(detectionmodel.DetectionStatusClean, detectionmodel.DetectionStatusSuspicious) {
		t.Fatal("did not expect clean to satisfy suspicious floor")
	}
}
