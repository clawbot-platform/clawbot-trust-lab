package benchmark

import (
	"context"
	"fmt"
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
	item, ok := s.items[id]
	if !ok {
		return domainscenario.Scenario{}, errNotFound("scenario", id)
	}
	return item, nil
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
			"commerce-h1-direct-human-purchase":                   {ID: "commerce-h1-direct-human-purchase", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-h2-human-refund-valid-history":              {ID: "commerce-h2-human-refund-valid-history", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-a1-agent-assisted-purchase-valid-controls":  {ID: "commerce-a1-agent-assisted-purchase-valid-controls", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-a2-fully-delegated-replenishment-purchase":  {ID: "commerce-a2-fully-delegated-replenishment-purchase", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-a3-agent-assisted-refund-approval-evidence": {ID: "commerce-a3-agent-assisted-refund-approval-evidence", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-s1-refund-weak-authorization":               {ID: "commerce-s1-refund-weak-authorization", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-s4-repeated-agent-refund-attempts":          {ID: "commerce-s4-repeated-agent-refund-attempts", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-s2-delegated-purchase-weak-provenance":      {ID: "commerce-s2-delegated-purchase-weak-provenance", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-s3-approval-removed-after-authorization":    {ID: "commerce-s3-approval-removed-after-authorization", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-s5-merchant-scope-drift-delegated-action":   {ID: "commerce-s5-merchant-scope-drift-delegated-action", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-v1-weakened-provenance":                     {ID: "commerce-v1-weakened-provenance", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-v2-expired-inactive-mandate":                {ID: "commerce-v2-expired-inactive-mandate", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-v3-approval-removed":                        {ID: "commerce-v3-approval-removed", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-v4-actor-switch-human-to-agent":             {ID: "commerce-v4-actor-switch-human-to-agent", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-v5-repeat-attempt-escalation":               {ID: "commerce-v5-repeat-attempt-escalation", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-v6-merchant-scope-drift":                    {ID: "commerce-v6-merchant-scope-drift", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-v7-high-value-delegated-purchase":           {ID: "commerce-v7-high-value-delegated-purchase", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-challenger-weakened-provenance-purchase":    {ID: "commerce-challenger-weakened-provenance-purchase", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-challenger-expired-mandate-purchase":        {ID: "commerce-challenger-expired-mandate-purchase", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-challenger-approval-removed-refund":         {ID: "commerce-challenger-approval-removed-refund", PackID: "challenger-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
			"commerce-clean-agent-assisted-purchase":              {ID: "commerce-clean-agent-assisted-purchase", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommercePurchase},
			"commerce-suspicious-refund-attempt":                  {ID: "commerce-suspicious-refund-attempt", PackID: "commerce-pack", Type: domainscenario.ScenarioTypeCommerceRefundReview},
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
	if len(round.StableScenarioRefs) != 7 {
		t.Fatalf("expected 7 stable scenarios, got %d", len(round.StableScenarioRefs))
	}
	if len(round.ChallengerVariantRefs) != 10 {
		t.Fatalf("expected 10 challenger variants, got %d", len(round.ChallengerVariantRefs))
	}

	assertNoPromotionsForTunedWeakCases(t, round.PromotionResults)
	assertTunedWeakCasesMeetMinimumPosture(t, round.ScenarioResults)

	if len(round.Reports.Artifacts) != 8 {
		t.Fatalf("expected 8 report artifacts, got %d", len(round.Reports.Artifacts))
	}
	for _, artifact := range round.Reports.Artifacts {
		if _, err := os.Stat(artifact.Path); err != nil {
			t.Fatalf("expected artifact %s to exist: %v", artifact.Path, err)
		}
	}
	if round.Summary.EvaluationMode != "shadow" {
		t.Fatalf("expected evaluation mode shadow, got %#v", round.Summary)
	}
	if round.Summary.BlockingMode != "recommendation_only" {
		t.Fatalf("expected blocking mode recommendation_only, got %#v", round.Summary)
	}
	if round.Summary.ExistingControlNote == "" || round.Summary.RecommendedFollowUp == "" {
		t.Fatalf("expected production-bridge summary fields to be populated, got %#v", round.Summary)
	}
	if len(round.Recommendations) == 0 {
		t.Fatal("expected recommendation data on round")
	}
	for _, item := range round.Recommendations {
		if item.LinkedRoundID != round.ID {
			t.Fatalf("expected recommendation linked to round %s, got %#v", round.ID, item)
		}
		if len(item.LinkedScenarioIDs) == 0 {
			t.Fatalf("expected recommendation to link scenarios, got %#v", item)
		}
		if item.SuggestedAction == "" || item.Rationale == "" {
			t.Fatalf("expected recommendation to be actionable, got %#v", item)
		}
		if item.ExistingControlIntegrationNote == "" {
			t.Fatalf("expected recommendation to include sidecar guidance, got %#v", item)
		}
	}
}

func TestRunRoundCalibratesStableScenariosBelievably(t *testing.T) {
	service := newRoundService(t)

	round, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() error = %v", err)
	}

	resultsByScenario := map[string]domainbenchmark.ScenarioResult{}
	for _, item := range round.ScenarioResults {
		resultsByScenario[item.ScenarioID] = item
	}

	benignStable := []string{
		"commerce-h1-direct-human-purchase",
		"commerce-h2-human-refund-valid-history",
		"commerce-a1-agent-assisted-purchase-valid-controls",
		"commerce-a2-fully-delegated-replenishment-purchase",
		"commerce-a3-agent-assisted-refund-approval-evidence",
	}
	for _, scenarioID := range benignStable {
		item := resultsByScenario[scenarioID]
		if item.FinalDetectionStatus != detectionmodel.DetectionStatusClean {
			t.Fatalf("expected benign stable scenario %s to be clean, got %#v", scenarioID, item)
		}
		if !item.Passed {
			t.Fatalf("expected benign stable scenario %s to pass its floor, got %#v", scenarioID, item)
		}
	}

	suspiciousStable := []string{
		"commerce-s1-refund-weak-authorization",
		"commerce-s4-repeated-agent-refund-attempts",
	}
	for _, scenarioID := range suspiciousStable {
		item := resultsByScenario[scenarioID]
		if item.FinalDetectionStatus == detectionmodel.DetectionStatusClean {
			t.Fatalf("expected suspicious stable scenario %s not to be clean, got %#v", scenarioID, item)
		}
		if !item.Passed {
			t.Fatalf("expected suspicious stable scenario %s to satisfy its floor, got %#v", scenarioID, item)
		}
	}
}

func TestRunRoundRetestsPriorPromotions(t *testing.T) {
	service := newRoundService(t)

	firstRound, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() first error = %v", err)
	}
	assertNoPromotionsForTunedWeakCases(t, firstRound.PromotionResults)
	assertTunedWeakCasesMeetMinimumPosture(t, firstRound.ScenarioResults)

	service.now = func() time.Time { return time.Date(2026, 3, 25, 12, 10, 0, 0, time.UTC) }
	secondRound, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() second error = %v", err)
	}

	assertNoPromotionsForTunedWeakCases(t, secondRound.PromotionResults)
	assertTunedWeakCasesMeetMinimumPosture(t, secondRound.ScenarioResults)

	if len(secondRound.Delta) == 0 {
		t.Fatal("expected detection delta in second round")
	}
}

func TestPromotedCasesDeduplicatesReplayRegressionInputs(t *testing.T) {
	round := domainbenchmark.BenchmarkRound{
		PromotionResults: []domainbenchmark.PromotionDecision{
			{
				ID:                  "promo-1",
				ScenarioID:          "commerce-s2-delegated-purchase-weak-provenance",
				ChallengerVariantID: "variant-v1-weakened-provenance",
				PromotionReason:     domainbenchmark.PromotionReasonSuspiciousBehaviorTooLow,
				Promoted:            true,
				CreatedAt:           time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
			},
			{
				ID:                  "promo-2",
				ScenarioID:          "commerce-s2-delegated-purchase-weak-provenance",
				ChallengerVariantID: "variant-v1-weakened-provenance",
				PromotionReason:     domainbenchmark.PromotionReasonDetectorMiss,
				Promoted:            true,
				CreatedAt:           time.Date(2026, 3, 25, 12, 1, 0, 0, time.UTC),
			},
		},
	}

	items := promotedCases(round)
	if len(items) != 1 {
		t.Fatalf("expected 1 deduplicated replay promotion, got %#v", items)
	}
	if items[0].PromotionReason != domainbenchmark.PromotionReasonDetectorMiss {
		t.Fatalf("expected strongest promotion reason to win, got %#v", items[0])
	}
}

func TestDedupePromotionsKeepsRoundCountsClean(t *testing.T) {
	items := []domainbenchmark.PromotionDecision{
		{
			ID:                  "promo-1",
			ScenarioID:          "commerce-s3-approval-removed-after-authorization",
			ChallengerVariantID: "variant-v3-approval-removed",
			PromotionReason:     domainbenchmark.PromotionReasonSuspiciousBehaviorTooLow,
			Promoted:            true,
			CreatedAt:           time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
		},
		{
			ID:                  "promo-2",
			ScenarioID:          "commerce-s3-approval-removed-after-authorization",
			ChallengerVariantID: "variant-v3-approval-removed",
			PromotionReason:     domainbenchmark.PromotionReasonMeaningfulRegression,
			Promoted:            true,
			CreatedAt:           time.Date(2026, 3, 25, 12, 1, 0, 0, time.UTC),
		},
	}

	deduped := dedupePromotions(items)
	if len(deduped) != 1 {
		t.Fatalf("expected one promotion after dedupe, got %#v", deduped)
	}
	if deduped[0].PromotionReason != domainbenchmark.PromotionReasonMeaningfulRegression {
		t.Fatalf("expected strongest reason to be preserved, got %#v", deduped[0])
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

func TestRunScheduledExecutesMultipleRounds(t *testing.T) {
	service := newRoundService(t)
	counter := 0
	service.now = func() time.Time {
		base := time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)
		value := base.Add(time.Duration(counter) * time.Minute)
		counter++
		return value
	}

	items, err := service.RunScheduled(context.Background(), domainbenchmark.SchedulerControlInput{
		ScenarioFamily: "commerce",
		Interval:       "1ms",
		MaxRuns:        2,
		DryRun:         true,
	})
	if err != nil {
		t.Fatalf("RunScheduled() error = %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 scheduled rounds, got %d", len(items))
	}
	if service.SchedulerStatus().ExecutedRuns < 2 {
		t.Fatalf("expected executed runs to be tracked, got %#v", service.SchedulerStatus())
	}
}

func TestListRoundsIncludesHistoricalRoundsAfterBootstrap(t *testing.T) {
	store := store.NewBenchmarkStore()
	store.PutHistorical(domainbenchmark.BenchmarkRound{
		ID:          "round-20260324120000",
		CompletedAt: time.Date(2026, 3, 24, 12, 0, 0, 0, time.UTC),
		Summary:     domainbenchmark.RoundSummary{RoundID: "round-20260324120000", ScenarioFamily: "commerce"},
		Reports: domainbenchmark.ReportIndex{
			RoundID:   "round-20260324120000",
			Directory: filepath.Join(t.TempDir(), "round-20260324120000"),
		},
	})

	service := NewService(registrarStub{}, nil, nil, nil, store, nil)
	items := service.ListRounds()
	if len(items) != 1 {
		t.Fatalf("expected 1 historical round, got %d", len(items))
	}
	if items[0].ID != "round-20260324120000" {
		t.Fatalf("unexpected historical round %#v", items[0])
	}
}

func TestLongRunSummaryAggregatesRecommendations(t *testing.T) {
	service := newRoundService(t)
	counter := 0
	service.now = func() time.Time {
		base := time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)
		value := base.Add(time.Duration(counter) * time.Minute)
		counter++
		return value
	}

	if _, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"}); err != nil {
		t.Fatalf("RunRound() first error = %v", err)
	}
	if _, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"}); err != nil {
		t.Fatalf("RunRound() second error = %v", err)
	}

	summary := service.LongRunSummary()
	if summary.RoundsExecuted != 2 {
		t.Fatalf("expected 2 rounds executed, got %#v", summary)
	}
	if len(summary.PromotionsOverTime) != 2 {
		t.Fatalf("expected promotions over time entries, got %#v", summary.PromotionsOverTime)
	}
	if summary.RecommendationCounts[domainbenchmark.RecommendationTypeMonitorInShadowMode] == 0 {
		t.Fatalf("expected shadow-mode recommendation count, got %#v", summary.RecommendationCounts)
	}
}
func TestRunRoundTunedWeakCasesMeetMinimumPosture(t *testing.T) {
	service := newRoundService(t)

	round, err := service.RunRound(context.Background(), domainbenchmark.RunInput{ScenarioFamily: "commerce"})
	if err != nil {
		t.Fatalf("RunRound() error = %v", err)
	}

	assertTunedWeakCasesMeetMinimumPosture(t, round.ScenarioResults)
}
func TestStatusRankOrdering(t *testing.T) {
	if !meetsMinimumStatus(detectionmodel.DetectionStatusStepUpRequired, detectionmodel.DetectionStatusSuspicious) {
		t.Fatal("expected step-up to satisfy suspicious floor")
	}
	if meetsMinimumStatus(detectionmodel.DetectionStatusClean, detectionmodel.DetectionStatusSuspicious) {
		t.Fatal("did not expect clean to satisfy suspicious floor")
	}
}

func errNotFound(kind string, id string) error {
	return fmt.Errorf("%s %s not found", kind, id)
}

var tunedWeakCaseMinimums = map[string]detectionmodel.DetectionStatus{
	"commerce-v2-expired-inactive-mandate":             detectionmodel.DetectionStatusStepUpRequired,
	"commerce-v3-approval-removed":                     detectionmodel.DetectionStatusStepUpRequired,
	"commerce-s3-approval-removed-after-authorization": detectionmodel.DetectionStatusStepUpRequired,
}

func assertNoPromotionsForTunedWeakCases(t *testing.T, promotions []domainbenchmark.PromotionDecision) {
	t.Helper()

	for _, promo := range promotions {
		if _, ok := tunedWeakCaseMinimums[promo.ScenarioID]; ok {
			t.Fatalf("expected tuned round to avoid promotion for targeted scenario %s, got %#v", promo.ScenarioID, promo)
		}
	}
}

func assertTunedWeakCasesMeetMinimumPosture(t *testing.T, results []domainbenchmark.ScenarioResult) {
	t.Helper()

	resultsByScenario := make(map[string]domainbenchmark.ScenarioResult, len(results))
	for _, item := range results {
		resultsByScenario[item.ScenarioID] = item
	}

	for scenarioID, minimum := range tunedWeakCaseMinimums {
		item, ok := resultsByScenario[scenarioID]
		if !ok {
			t.Fatalf("expected tuned round to include scenario %s", scenarioID)
		}
		if !item.Passed {
			t.Fatalf("expected tuned weak case %s to meet its minimum posture, got %#v", scenarioID, item)
		}
		if item.FinalDetectionStatus != minimum {
			t.Fatalf("expected tuned weak case %s to land on %s, got %#v", scenarioID, minimum, item)
		}
	}
}
