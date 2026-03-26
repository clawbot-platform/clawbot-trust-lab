package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/domain/commerce"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	domainevents "clawbot-trust-lab/internal/domain/events"
	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/domain/scenario"
	"clawbot-trust-lab/internal/domain/trust"
	detectionsvc "clawbot-trust-lab/internal/services/detection"
	operatorsvc "clawbot-trust-lab/internal/services/operator"
	executionsvc "clawbot-trust-lab/internal/services/scenario"
	"clawbot-trust-lab/internal/version"
)

type scenarioServiceStub struct{}

func (scenarioServiceStub) ListPacks() []scenario.ScenarioPack {
	return []scenario.ScenarioPack{{ID: "starter-pack", Name: "Starter Pack"}}
}
func (scenarioServiceStub) GetPack(string) (scenario.ScenarioPack, error) {
	return scenario.ScenarioPack{ID: "starter-pack", Name: "Starter Pack"}, nil
}

type scenarioServiceErrorStub struct {
	scenarioServiceStub
}

func (scenarioServiceErrorStub) GetPack(string) (scenario.ScenarioPack, error) {
	return scenario.ScenarioPack{}, errors.New("scenario pack not found")
}

type executionServiceStub struct {
	result executionsvc.ExecutionResult
	err    error
}

func (s executionServiceStub) ListScenarios() []scenario.Scenario {
	return []scenario.Scenario{{ID: "commerce-clean-agent-assisted-purchase", Name: "Clean Purchase"}}
}
func (s executionServiceStub) Execute(context.Context, string) (executionsvc.ExecutionResult, error) {
	if s.err != nil {
		return executionsvc.ExecutionResult{}, s.err
	}
	return s.result, nil
}

type trustServiceStub struct {
	items       []trust.TrustArtifact
	createErr   error
	contextResp memory.LoadScenarioContextResponse
	contextErr  error
}

func (s *trustServiceStub) CreateArtifact(_ context.Context, input trust.CreateArtifactInput) (trust.TrustArtifact, error) {
	if s.createErr != nil {
		return trust.TrustArtifact{}, s.createErr
	}
	item := trust.TrustArtifact{ID: "ta-" + input.ScenarioID, SourceScenarioID: input.ScenarioID}
	s.items = append(s.items, item)
	return item, nil
}
func (s *trustServiceStub) ListArtifacts() []trust.TrustArtifact { return s.items }
func (s *trustServiceStub) LoadMemoryContext(context.Context, string) (memory.LoadScenarioContextResponse, error) {
	if s.contextErr != nil {
		return memory.LoadScenarioContextResponse{}, s.contextErr
	}
	return s.contextResp, nil
}

type replayServiceStub struct {
	items       []replay.ReplayCase
	createErr   error
	similarResp memory.FetchSimilarCasesResponse
	similarErr  error
}

func (s *replayServiceStub) CreateCase(_ context.Context, input replay.CreateCaseInput) (replay.ReplayCase, error) {
	if s.createErr != nil {
		return replay.ReplayCase{}, s.createErr
	}
	item := replay.ReplayCase{ID: "rc-1", ScenarioID: input.ScenarioID}
	s.items = append(s.items, item)
	return item, nil
}
func (s *replayServiceStub) ListCases() []replay.ReplayCase { return s.items }
func (s *replayServiceStub) SimilarCases(context.Context, string) (memory.FetchSimilarCasesResponse, error) {
	if s.similarErr != nil {
		return memory.FetchSimilarCasesResponse{}, s.similarErr
	}
	return s.similarResp, nil
}

type benchmarkServiceStub struct{}

func (benchmarkServiceStub) RegisterRound(_ context.Context, _ benchmark.RegistrationRequest) (benchmark.RegistrationResult, error) {
	return benchmark.RegistrationResult{RegistrationID: "bench-1", Status: "accepted_stub", RegisteredAt: time.Now().UTC()}, nil
}
func (benchmarkServiceStub) RunRound(_ context.Context, _ benchmark.RunInput) (benchmark.BenchmarkRound, error) {
	return benchmark.BenchmarkRound{
		ID:             "round-20260325120000",
		ScenarioFamily: "commerce",
		RoundStatus:    benchmark.RoundStatusCompleted,
		StableSet:      benchmark.StableSetResult{TotalCount: 2, PassedCount: 2},
		LivingSet:      benchmark.LivingSetResult{TotalCount: 3, CaughtCount: 2, PromotionCount: 1},
		PromotionResults: []benchmark.PromotionDecision{{
			ID:                  "promo-1",
			ScenarioID:          "commerce-challenger-weakened-provenance-purchase",
			ChallengerVariantID: "variant-weakened-provenance",
			PromotionReason:     benchmark.PromotionReasonDetectorMiss,
			Rationale:           "Suspicious challenger behavior evaluated as clean.",
			Promoted:            true,
		}},
		Summary: benchmark.RoundSummary{
			RoundID:             "round-20260325120000",
			ScenarioFamily:      "commerce",
			StableScenarioCount: 2,
			ChallengerCount:     3,
			ReplayRetestCount:   0,
			PromotionCount:      1,
			ReplayPassRate:      1,
			RobustnessOutcome:   benchmark.RobustnessOutcomeNewBlindSpotDiscovered,
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
		},
		Reports: benchmark.ReportIndex{
			RoundID:   "round-20260325120000",
			Directory: "./reports/round-20260325120000",
			Artifacts: []benchmark.ReportArtifact{{Name: "round-summary.json", Path: "./reports/round-20260325120000/round-summary.json", Kind: "json"}},
		},
		Recommendations: []benchmark.Recommendation{{
			ID:                "rec-round-20260325120000-replay",
			Type:              benchmark.RecommendationTypeAddToReplayStableSet,
			LinkedRoundID:     "round-20260325120000",
			LinkedScenarioIDs: []string{"commerce-challenger-weakened-provenance-purchase"},
			SuggestedAction:   "Add the promoted case into replay.",
		}},
	}, nil
}
func (benchmarkServiceStub) RunScheduled(_ context.Context, _ benchmark.SchedulerControlInput) ([]benchmark.BenchmarkRound, error) {
	item, _ := benchmarkServiceStub{}.RunRound(context.Background(), benchmark.RunInput{ScenarioFamily: "commerce"})
	return []benchmark.BenchmarkRound{item}, nil
}
func (benchmarkServiceStub) ListRounds() []benchmark.BenchmarkRound {
	items, _ := benchmarkServiceStub{}.RunRound(context.Background(), benchmark.RunInput{ScenarioFamily: "commerce"})
	return []benchmark.BenchmarkRound{items}
}
func (benchmarkServiceStub) GetRound(id string) (benchmark.BenchmarkRound, error) {
	item, _ := benchmarkServiceStub{}.RunRound(context.Background(), benchmark.RunInput{ScenarioFamily: "commerce"})
	if item.ID != id {
		return benchmark.BenchmarkRound{}, errors.New("benchmark round not found")
	}
	return item, nil
}
func (benchmarkServiceStub) GetRoundSummary(id string) (benchmark.RoundSummary, error) {
	item, err := benchmarkServiceStub{}.GetRound(id)
	if err != nil {
		return benchmark.RoundSummary{}, err
	}
	return item.Summary, nil
}
func (benchmarkServiceStub) GetRoundPromotions(id string) ([]benchmark.PromotionDecision, error) {
	item, err := benchmarkServiceStub{}.GetRound(id)
	if err != nil {
		return nil, err
	}
	return item.PromotionResults, nil
}
func (benchmarkServiceStub) GetRoundDelta(id string) ([]benchmark.DetectionDelta, error) {
	_, err := benchmarkServiceStub{}.GetRound(id)
	if err != nil {
		return nil, err
	}
	return []benchmark.DetectionDelta{{ScenarioID: "commerce-challenger-weakened-provenance-purchase"}}, nil
}
func (benchmarkServiceStub) GetRoundReports(id string) (benchmark.ReportIndex, error) {
	item, err := benchmarkServiceStub{}.GetRound(id)
	if err != nil {
		return benchmark.ReportIndex{}, err
	}
	return item.Reports, nil
}
func (benchmarkServiceStub) ListRecommendations() []benchmark.Recommendation {
	item, _ := benchmarkServiceStub{}.RunRound(context.Background(), benchmark.RunInput{ScenarioFamily: "commerce"})
	return item.Recommendations
}
func (benchmarkServiceStub) GetRecommendation(id string) (benchmark.Recommendation, error) {
	items := benchmarkServiceStub{}.ListRecommendations()
	for _, item := range items {
		if item.ID == id {
			return item, nil
		}
	}
	return benchmark.Recommendation{}, errors.New("recommendation not found")
}
func (benchmarkServiceStub) LongRunSummary() benchmark.LongRunSummary {
	return benchmark.LongRunSummary{RoundsExecuted: 1, NewBlindSpots: 1}
}
func (benchmarkServiceStub) SchedulerStatus() benchmark.SchedulerStatus {
	return benchmark.SchedulerStatus{Enabled: true, ScenarioFamily: "commerce", Interval: "24h", MaxRuns: 7}
}
func (benchmarkServiceStub) Status() map[string]any {
	return map[string]any{"registrations": 1, "last_status": "accepted_stub", "scheduler": benchmark.SchedulerStatus{Enabled: true}}
}

type commerceServiceStub struct {
	orders map[string]commerce.Order
}

func (s commerceServiceStub) ListOrders() []commerce.Order {
	items := make([]commerce.Order, 0, len(s.orders))
	for _, order := range s.orders {
		items = append(items, order)
	}
	return items
}
func (s commerceServiceStub) GetOrder(id string) (commerce.Order, error) {
	item, ok := s.orders[id]
	if !ok {
		return commerce.Order{}, errors.New("order not found")
	}
	return item, nil
}

type eventServiceStub struct {
	items []domainevents.Record
}

func (s eventServiceStub) ListEvents() []domainevents.Record { return s.items }

type trustDecisionServiceStub struct {
	items map[string]trust.TrustDecision
}

func (s trustDecisionServiceStub) ListDecisions() []trust.TrustDecision {
	items := make([]trust.TrustDecision, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}
	return items
}
func (s trustDecisionServiceStub) GetDecision(id string) (trust.TrustDecision, error) {
	item, ok := s.items[id]
	if !ok {
		return trust.TrustDecision{}, errors.New("trust decision not found")
	}
	return item, nil
}

type detectionServiceStub struct {
	result  detectionmodel.DetectionResult
	results []detectionmodel.DetectionResult
	summary detectionmodel.DetectionRunSummary
	inputs  []detectionsvc.EvaluateInput
	err     error
}

func (s *detectionServiceStub) Evaluate(_ context.Context, input detectionsvc.EvaluateInput) (detectionmodel.DetectionResult, error) {
	s.inputs = append(s.inputs, input)
	if s.err != nil {
		return detectionmodel.DetectionResult{}, s.err
	}
	return s.result, nil
}

func (s *detectionServiceStub) ListResults() []detectionmodel.DetectionResult {
	return append([]detectionmodel.DetectionResult(nil), s.results...)
}

func (s *detectionServiceStub) GetResult(id string) (detectionmodel.DetectionResult, error) {
	for _, item := range s.results {
		if item.ID == id {
			return item, nil
		}
	}
	return detectionmodel.DetectionResult{}, errors.New("detection result not found")
}

func (s *detectionServiceStub) Rules() []detectionmodel.RuleDefinition {
	return []detectionmodel.RuleDefinition{
		{ID: "refund_weak_authorization", Title: "Refund with weak authorization", Severity: 25},
	}
}

func (s *detectionServiceStub) Summary() detectionmodel.DetectionRunSummary {
	return s.summary
}

type operatorServiceStub struct{}

func (s operatorServiceStub) ListRounds() []benchmark.BenchmarkRound {
	return []benchmark.BenchmarkRound{{
		ID:             "round-20260325120000",
		ScenarioFamily: "commerce",
		Summary: benchmark.RoundSummary{
			PromotionCount:    1,
			ReplayPassRate:    1,
			RobustnessOutcome: benchmark.RobustnessOutcomeNewBlindSpotDiscovered,
		},
	}}
}

func (s operatorServiceStub) GetRound(id string) (benchmark.BenchmarkRound, error) {
	if id != "round-20260325120000" {
		return benchmark.BenchmarkRound{}, errors.New("round not found")
	}
	return s.ListRounds()[0], nil
}

func (s operatorServiceStub) CompareRounds(current string, previous string) (benchmark.RoundComparison, error) {
	if current == "" || previous == "" {
		return benchmark.RoundComparison{}, errors.New("missing round id")
	}
	return benchmark.RoundComparison{
		CurrentRoundID:       current,
		PreviousRoundID:      previous,
		PromotionsCountDelta: 1,
	}, nil
}

func (s operatorServiceStub) ListPromotions(string) []operatorsvc.PromotionRecord {
	return []operatorsvc.PromotionRecord{{
		RoundID: "round-20260325120000",
		Promotion: benchmark.PromotionDecision{
			ID:                 "promo-1",
			ScenarioID:         "commerce-challenger-weakened-provenance-purchase",
			PromotionReason:    benchmark.PromotionReasonDetectorMiss,
			DetectionResultRef: "det-order-suspicious-refund-attempt",
		},
	}}
}

func (s operatorServiceStub) GetPromotion(id string) (operatorsvc.PromotionDetail, error) {
	if id != "promo-1" {
		return operatorsvc.PromotionDetail{}, errors.New("promotion not found")
	}
	return operatorsvc.PromotionDetail{
		RoundID: "round-20260325120000",
		Promotion: benchmark.PromotionDecision{
			ID:                 "promo-1",
			ScenarioID:         "commerce-challenger-weakened-provenance-purchase",
			PromotionReason:    benchmark.PromotionReasonDetectorMiss,
			DetectionResultRef: "det-order-suspicious-refund-attempt",
		},
		DetectionResult: detectionmodel.DetectionResult{
			ID:         "det-order-suspicious-refund-attempt",
			ScenarioID: "commerce-challenger-weakened-provenance-purchase",
			Status:     detectionmodel.DetectionStatusClean,
		},
	}, nil
}

func (s operatorServiceStub) ReviewPromotion(id string, input operatorsvc.ReviewInput) (benchmark.PromotionReview, error) {
	if id == "" || input.Status == "" {
		return benchmark.PromotionReview{}, errors.New("invalid review")
	}
	return benchmark.PromotionReview{
		PromotionID: id,
		Status:      benchmark.PromotionReviewStatus(input.Status),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

func (s operatorServiceStub) GetDetectionResult(id string) (detectionmodel.DetectionResult, error) {
	if id != "det-order-suspicious-refund-attempt" {
		return detectionmodel.DetectionResult{}, errors.New("detection result not found")
	}
	return detectionmodel.DetectionResult{
		ID:         id,
		ScenarioID: "commerce-challenger-weakened-provenance-purchase",
		Status:     detectionmodel.DetectionStatusClean,
	}, nil
}

func (s operatorServiceStub) ListRecommendations() []benchmark.Recommendation {
	return []benchmark.Recommendation{{
		ID:                "rec-round-20260325120000-replay",
		Type:              benchmark.RecommendationTypeAddToReplayStableSet,
		LinkedRoundID:     "round-20260325120000",
		LinkedScenarioIDs: []string{"commerce-challenger-weakened-provenance-purchase"},
		SuggestedAction:   "Add the promoted case into replay.",
	}}
}

func (s operatorServiceStub) GetRecommendation(id string) (benchmark.Recommendation, error) {
	for _, item := range s.ListRecommendations() {
		if item.ID == id {
			return item, nil
		}
	}
	return benchmark.Recommendation{}, errors.New("recommendation not found")
}

func (s operatorServiceStub) GetTrendSummary() benchmark.LongRunSummary {
	return benchmark.LongRunSummary{RoundsExecuted: 2, NewBlindSpots: 1}
}

func (s operatorServiceStub) GetReports(roundID string) ([]benchmark.ReportDescriptor, error) {
	if roundID != "round-20260325120000" {
		return nil, errors.New("round not found")
	}
	return []benchmark.ReportDescriptor{{
		RoundID:      roundID,
		ArtifactName: "executive-summary.md",
		Path:         "./reports/round-20260325120000/executive-summary.md",
		Kind:         "markdown",
	}}, nil
}

func (s operatorServiceStub) GetReportArtifact(roundID string, artifactName string) (operatorsvc.ReportContent, error) {
	if roundID != "round-20260325120000" || artifactName != "executive-summary.md" {
		return operatorsvc.ReportContent{}, errors.New("report artifact not found")
	}
	return operatorsvc.ReportContent{
		Descriptor: benchmark.ReportDescriptor{
			RoundID:      roundID,
			ArtifactName: artifactName,
			Path:         "./reports/round-20260325120000/executive-summary.md",
			Kind:         "markdown",
		},
		Content: "# Executive Summary",
	}, nil
}

func TestSystemHandlerHealth(t *testing.T) {
	handler := NewSystemHandler(func(context.Context) error { return nil }, version.Current())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	handler.Health(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestSystemHandlerReadyAndVersion(t *testing.T) {
	ready := NewSystemHandler(func(context.Context) error { return nil }, version.Current())
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	recorder := httptest.NewRecorder()
	ready.Ready(recorder, req)
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"ready"`)) {
		t.Fatalf("expected ready response, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	notReady := NewSystemHandler(func(context.Context) error { return errors.New("dependency down") }, version.Current())
	recorder = httptest.NewRecorder()
	notReady.Ready(recorder, req)
	if recorder.Code != http.StatusServiceUnavailable || !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"not_ready"`)) {
		t.Fatalf("expected not ready response, got %d body=%s", recorder.Code, recorder.Body.String())
	}

	versionReq := httptest.NewRequest(http.MethodGet, "/version", nil)
	versionRecorder := httptest.NewRecorder()
	ready.Version(versionRecorder, versionReq)
	if versionRecorder.Code != http.StatusOK || !bytes.Contains(versionRecorder.Body.Bytes(), []byte(`"version"`)) {
		t.Fatalf("expected version payload, got %d body=%s", versionRecorder.Code, versionRecorder.Body.String())
	}
}

func TestTrustLabHandlerListScenarios(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/scenarios", nil)
	recorder := httptest.NewRecorder()

	handler.ListScenarios(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestTrustLabHandlerExecuteScenario(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scenarios/execute", bytes.NewBufferString(`{"scenario_id":"commerce-clean-agent-assisted-purchase"}`))
	recorder := httptest.NewRecorder()

	handler.ExecuteScenario(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestTrustLabHandlerExecuteScenarioReturnsBadGatewayOnMemoryFailure(t *testing.T) {
	handler := newHandler()
	handler.execution = executionServiceStub{err: &trust.MemorySyncError{Err: errors.New("clawmem unavailable")}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scenarios/execute", bytes.NewBufferString(`{"scenario_id":"commerce-clean-agent-assisted-purchase"}`))
	recorder := httptest.NewRecorder()

	handler.ExecuteScenario(recorder, req)

	if recorder.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestTrustLabHandlerCreateArtifact(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/trust/artifacts", bytes.NewBufferString(`{"scenario_id":"starter-mandate-review"}`))
	recorder := httptest.NewRecorder()

	handler.CreateArtifact(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
}

func TestTrustLabHandlerCreateReplayCase(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/replay/cases", bytes.NewBufferString(`{"scenario_id":"starter-mandate-review","outcome_summary":"ok"}`))
	recorder := httptest.NewRecorder()

	handler.CreateReplayCase(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
}

func TestTrustLabHandlerTrustStatusIncludesMemoryContext(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trust/status?scenario_id=starter-mandate-review", nil)
	recorder := httptest.NewRecorder()

	handler.TrustStatus(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"memory_status":"ok"`)) {
		t.Fatalf("expected memory status in body: %s", recorder.Body.String())
	}
}

func TestTrustLabHandlerGetOrder(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/order-1", nil)
	req.SetPathValue("id", "order-1")
	recorder := httptest.NewRecorder()

	handler.GetOrder(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestTrustLabHandlerListTrustDecisions(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/trust/decisions", nil)
	recorder := httptest.NewRecorder()

	handler.ListTrustDecisions(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}

func TestTrustLabHandlerRegisterBenchmarkRound(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/benchmark/rounds/register", bytes.NewBufferString(`{"scenario_pack_id":"starter-pack","scenario_pack_version":"v1","replay_case_refs":["rc-1"]}`))
	recorder := httptest.NewRecorder()

	handler.RegisterBenchmarkRound(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", recorder.Code)
	}
}

func TestTrustLabHandlerRunBenchmarkRound(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/benchmark/rounds/run", bytes.NewBufferString(`{"scenario_family":"commerce"}`))
	recorder := httptest.NewRecorder()

	handler.RunBenchmarkRound(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"round_status":"completed"`)) {
		t.Fatalf("expected round status in body: %s", recorder.Body.String())
	}
}

func TestTrustLabHandlerGetBenchmarkRoundReports(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/benchmark/rounds/round-20260325120000/reports", nil)
	req.SetPathValue("id", "round-20260325120000")
	recorder := httptest.NewRecorder()

	handler.GetBenchmarkRoundReports(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"round-summary.json"`)) {
		t.Fatalf("expected report artifact in body: %s", recorder.Body.String())
	}
}

func TestTrustLabHandlerListBenchmarkRecommendations(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/benchmark/recommendations", nil)
	recorder := httptest.NewRecorder()

	handler.ListBenchmarkRecommendations(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"add_to_replay_stable_set"`)) {
		t.Fatalf("expected recommendation payload in body: %s", recorder.Body.String())
	}
}

func TestTrustLabHandlerGetBenchmarkTrendSummary(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/benchmark/trends/summary", nil)
	recorder := httptest.NewRecorder()

	handler.GetBenchmarkTrendSummary(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"rounds_executed":1`)) {
		t.Fatalf("expected trend summary in body: %s", recorder.Body.String())
	}
}

func TestTrustLabHandlerReadEndpoints(t *testing.T) {
	handler := newHandler()
	tests := []struct {
		name       string
		target     func(http.ResponseWriter, *http.Request)
		path       string
		pathValues map[string]string
		wantCode   int
		wantBody   string
	}{
		{name: "scenario types", target: handler.ScenarioTypes, path: "/api/v1/scenarios/types", wantCode: http.StatusOK, wantBody: `"types"`},
		{name: "list packs", target: handler.ListPacks, path: "/api/v1/scenarios/packs", wantCode: http.StatusOK, wantBody: `"starter-pack"`},
		{name: "get pack", target: handler.GetPack, path: "/api/v1/scenarios/packs/starter-pack", pathValues: map[string]string{"id": "starter-pack"}, wantCode: http.StatusOK, wantBody: `"Starter Pack"`},
		{name: "benchmark status", target: handler.BenchmarkStatus, path: "/api/v1/benchmark/status", wantCode: http.StatusOK, wantBody: `"registrations"`},
		{name: "list rounds", target: handler.ListBenchmarkRounds, path: "/api/v1/benchmark/rounds", wantCode: http.StatusOK, wantBody: `"round-20260325120000"`},
		{name: "get round", target: handler.GetBenchmarkRound, path: "/api/v1/benchmark/rounds/round-20260325120000", pathValues: map[string]string{"id": "round-20260325120000"}, wantCode: http.StatusOK, wantBody: `"round_status":"completed"`},
		{name: "get round summary", target: handler.GetBenchmarkRoundSummary, path: "/api/v1/benchmark/rounds/round-20260325120000/summary", pathValues: map[string]string{"id": "round-20260325120000"}, wantCode: http.StatusOK, wantBody: `"robustness_outcome"`},
		{name: "get round promotions", target: handler.GetBenchmarkRoundPromotions, path: "/api/v1/benchmark/rounds/round-20260325120000/promotions", pathValues: map[string]string{"id": "round-20260325120000"}, wantCode: http.StatusOK, wantBody: `"promo-1"`},
		{name: "get recommendation", target: handler.GetBenchmarkRecommendation, path: "/api/v1/benchmark/recommendations/rec-round-20260325120000-replay", pathValues: map[string]string{"id": "rec-round-20260325120000-replay"}, wantCode: http.StatusOK, wantBody: `"add_to_replay_stable_set"`},
		{name: "get scheduler status", target: handler.GetBenchmarkSchedulerStatus, path: "/api/v1/benchmark/scheduler/status", wantCode: http.StatusOK, wantBody: `"enabled":true`},
		{name: "list artifacts", target: handler.ListArtifacts, path: "/api/v1/trust/artifacts", wantCode: http.StatusOK, wantBody: `"data"`},
		{name: "list replay cases", target: handler.ListReplayCases, path: "/api/v1/replay/cases", wantCode: http.StatusOK, wantBody: `"data"`},
		{name: "list orders", target: handler.ListOrders, path: "/api/v1/orders", wantCode: http.StatusOK, wantBody: `"order-1"`},
		{name: "list events", target: handler.ListEvents, path: "/api/v1/events", wantCode: http.StatusOK, wantBody: `"evt-1"`},
		{name: "get trust decision", target: handler.GetTrustDecision, path: "/api/v1/trust/decisions/decision-1", pathValues: map[string]string{"id": "decision-1"}, wantCode: http.StatusOK, wantBody: `"accepted"`},
		{name: "benchmark round status", target: handler.BenchmarkRoundStatus, path: "/api/v1/benchmark/rounds/status", wantCode: http.StatusOK, wantBody: `"scheduler"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			for key, value := range tc.pathValues {
				req.SetPathValue(key, value)
			}
			recorder := httptest.NewRecorder()
			tc.target(recorder, req)
			if recorder.Code != tc.wantCode || !bytes.Contains(recorder.Body.Bytes(), []byte(tc.wantBody)) {
				t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestTrustLabHandlerAdditionalReadAndSchedulerEndpoints(t *testing.T) {
	handler := newHandler()

	t.Run("detection rules", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/detection/rules", nil)
		recorder := httptest.NewRecorder()
		handler.ListDetectionRules(recorder, req)
		if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"refund_weak_authorization"`)) {
			t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("get detection result", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/detection/results/det-order-suspicious-refund-attempt", nil)
		req.SetPathValue("id", "det-order-suspicious-refund-attempt")
		recorder := httptest.NewRecorder()
		handler.GetDetectionResult(recorder, req)
		if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"det-order-suspicious-refund-attempt"`)) {
			t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("replay status with memory context", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/replay/status?scenario_id=starter-mandate-review", nil)
		recorder := httptest.NewRecorder()
		handler.ReplayStatus(recorder, req)
		if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"memory_status":"ok"`)) {
			t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("round delta", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/benchmark/rounds/round-20260325120000/delta", nil)
		req.SetPathValue("id", "round-20260325120000")
		recorder := httptest.NewRecorder()
		handler.GetBenchmarkRoundDelta(recorder, req)
		if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"commerce-challenger-weakened-provenance-purchase"`)) {
			t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("run scheduler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/benchmark/scheduler/run", bytes.NewBufferString(`{"scenario_family":"commerce","interval":"24h","max_runs":1}`))
		recorder := httptest.NewRecorder()
		handler.RunBenchmarkScheduler(recorder, req)
		if recorder.Code != http.StatusCreated || !bytes.Contains(recorder.Body.Bytes(), []byte(`"summary"`)) {
			t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
		}
	})
}

func TestTrustLabHandlerNotFoundAndBadInputPaths(t *testing.T) {
	handler := newHandler()
	missingPackHandler := newHandler()
	missingPackHandler.scenarios = scenarioServiceErrorStub{}
	tests := []struct {
		name       string
		target     func(http.ResponseWriter, *http.Request)
		method     string
		path       string
		body       string
		pathValues map[string]string
		wantCode   int
	}{
		{name: "get detection result missing", target: handler.GetDetectionResult, method: http.MethodGet, path: "/api/v1/detection/results/missing", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "get pack missing", target: missingPackHandler.GetPack, method: http.MethodGet, path: "/api/v1/scenarios/packs/missing", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "round delta missing", target: handler.GetBenchmarkRoundDelta, method: http.MethodGet, path: "/api/v1/benchmark/rounds/missing/delta", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "get benchmark recommendation missing", target: handler.GetBenchmarkRecommendation, method: http.MethodGet, path: "/api/v1/benchmark/recommendations/missing", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "run scheduler bad json", target: handler.RunBenchmarkScheduler, method: http.MethodPost, path: "/api/v1/benchmark/scheduler/run", body: "{", wantCode: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
			for key, value := range tc.pathValues {
				req.SetPathValue(key, value)
			}
			recorder := httptest.NewRecorder()
			tc.target(recorder, req)
			if recorder.Code != tc.wantCode {
				t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestTrustLabHandlerEvaluateDetection(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/detection/evaluate", bytes.NewBufferString(`{"scenario_id":"commerce-suspicious-refund-attempt"}`))
	recorder := httptest.NewRecorder()

	handler.EvaluateDetection(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"step_up_required"`)) {
		t.Fatalf("expected detection response in body: %s", recorder.Body.String())
	}
}

func TestTrustLabHandlerListDetectionResults(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/detection/results", nil)
	recorder := httptest.NewRecorder()

	handler.ListDetectionResults(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestTrustLabHandlerDetectionSummary(t *testing.T) {
	handler := newHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/detection/summary", nil)
	recorder := httptest.NewRecorder()

	handler.DetectionSummary(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"total":1`)) {
		t.Fatalf("expected total in body: %s", recorder.Body.String())
	}
}

func TestOperatorHandlerListRounds(t *testing.T) {
	handler := NewOperatorHandler(operatorServiceStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/operator/rounds", nil)
	recorder := httptest.NewRecorder()

	handler.ListRounds(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestOperatorHandlerReadEndpoints(t *testing.T) {
	handler := NewOperatorHandler(operatorServiceStub{})

	tests := []struct {
		name       string
		target     func(http.ResponseWriter, *http.Request)
		path       string
		pathValues map[string]string
		wantCode   int
		wantBody   string
	}{
		{
			name:   "get round",
			target: handler.GetRound,
			path:   "/api/v1/operator/rounds/round-20260325120000",
			pathValues: map[string]string{
				"id": "round-20260325120000",
			},
			wantCode: http.StatusOK,
			wantBody: `"id":"round-20260325120000"`,
		},
		{
			name:     "list promotions",
			target:   handler.ListPromotions,
			path:     "/api/v1/operator/promotions",
			wantCode: http.StatusOK,
			wantBody: `"promo-1"`,
		},
		{
			name:   "get promotion",
			target: handler.GetPromotion,
			path:   "/api/v1/operator/promotions/promo-1",
			pathValues: map[string]string{
				"id": "promo-1",
			},
			wantCode: http.StatusOK,
			wantBody: `"promotion_reason":"detector_miss"`,
		},
		{
			name:   "get recommendation",
			target: handler.GetRecommendation,
			path:   "/api/v1/operator/recommendations/rec-round-20260325120000-replay",
			pathValues: map[string]string{
				"id": "rec-round-20260325120000-replay",
			},
			wantCode: http.StatusOK,
			wantBody: `"add_to_replay_stable_set"`,
		},
		{
			name:   "get reports",
			target: handler.GetReports,
			path:   "/api/v1/operator/reports/round-20260325120000",
			pathValues: map[string]string{
				"round_id": "round-20260325120000",
			},
			wantCode: http.StatusOK,
			wantBody: `"executive-summary.md"`,
		},
		{
			name:   "get detection result",
			target: handler.GetDetectionResult,
			path:   "/api/v1/operator/detection/results/det-order-suspicious-refund-attempt",
			pathValues: map[string]string{
				"id": "det-order-suspicious-refund-attempt",
			},
			wantCode: http.StatusOK,
			wantBody: `"status":"clean"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			for key, value := range tc.pathValues {
				req.SetPathValue(key, value)
			}
			recorder := httptest.NewRecorder()
			tc.target(recorder, req)
			if recorder.Code != tc.wantCode || !bytes.Contains(recorder.Body.Bytes(), []byte(tc.wantBody)) {
				t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestOperatorHandlerRecommendationsAndNotFoundPaths(t *testing.T) {
	handler := NewOperatorHandler(operatorServiceStub{})

	t.Run("list recommendations", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/operator/recommendations", nil)
		recorder := httptest.NewRecorder()
		handler.ListRecommendations(recorder, req)
		if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"add_to_replay_stable_set"`)) {
			t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
		}
	})

	tests := []struct {
		name       string
		target     func(http.ResponseWriter, *http.Request)
		path       string
		pathValues map[string]string
		query      string
		wantCode   int
	}{
		{name: "round missing", target: handler.GetRound, path: "/api/v1/operator/rounds/missing", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "compare rounds service error", target: handler.CompareRounds, path: "/api/v1/operator/rounds/compare?previous=round-1", wantCode: http.StatusBadRequest},
		{name: "promotion missing", target: handler.GetPromotion, path: "/api/v1/operator/promotions/missing", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "detection missing", target: handler.GetDetectionResult, path: "/api/v1/operator/detections/missing", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "recommendation missing", target: handler.GetRecommendation, path: "/api/v1/operator/recommendations/missing", pathValues: map[string]string{"id": "missing"}, wantCode: http.StatusNotFound},
		{name: "reports missing", target: handler.GetReports, path: "/api/v1/operator/rounds/missing/reports", pathValues: map[string]string{"round_id": "missing"}, wantCode: http.StatusNotFound},
		{name: "report artifact missing", target: handler.GetReportArtifact, path: "/api/v1/operator/rounds/round-20260325120000/reports/missing", pathValues: map[string]string{"round_id": "round-20260325120000", "artifact_name": "missing"}, wantCode: http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			for key, value := range tc.pathValues {
				req.SetPathValue(key, value)
			}
			recorder := httptest.NewRecorder()
			tc.target(recorder, req)
			if recorder.Code != tc.wantCode {
				t.Fatalf("unexpected response code=%d body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestOperatorHandlerCompareRoundsRequiresPrevious(t *testing.T) {
	handler := NewOperatorHandler(operatorServiceStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/operator/rounds/round-20260325120000/compare", nil)
	req.SetPathValue("id", "round-20260325120000")
	recorder := httptest.NewRecorder()

	handler.CompareRounds(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestOperatorHandlerReviewPromotion(t *testing.T) {
	handler := NewOperatorHandler(operatorServiceStub{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/operator/promotions/promo-1/review", bytes.NewBufferString(`{"status":"accepted","note":"Looks real."}`))
	req.SetPathValue("id", "promo-1")
	recorder := httptest.NewRecorder()

	handler.ReviewPromotion(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"status":"accepted"`)) {
		t.Fatalf("expected review status in body: %s", recorder.Body.String())
	}
}

func TestOperatorHandlerGetReportArtifact(t *testing.T) {
	handler := NewOperatorHandler(operatorServiceStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/operator/reports/round-20260325120000/executive-summary.md", nil)
	req.SetPathValue("round_id", "round-20260325120000")
	req.SetPathValue("artifact_name", "executive-summary.md")
	recorder := httptest.NewRecorder()

	handler.GetReportArtifact(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"# Executive Summary"`)) {
		t.Fatalf("expected report content in body: %s", recorder.Body.String())
	}
}

func TestOperatorHandlerGetTrendSummary(t *testing.T) {
	handler := NewOperatorHandler(operatorServiceStub{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/operator/trends/summary", nil)
	recorder := httptest.NewRecorder()

	handler.GetTrendSummary(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", recorder.Code, recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"rounds_executed":2`)) {
		t.Fatalf("expected operator trend summary in body: %s", recorder.Body.String())
	}
}

func newHandler() *TrustLabHandler {
	detectionStub := &detectionServiceStub{
		result: detectionmodel.DetectionResult{
			ID:             "det-order-suspicious-refund-attempt",
			ScenarioID:     "commerce-suspicious-refund-attempt",
			OrderID:        "order-suspicious-refund-attempt",
			Status:         detectionmodel.DetectionStatusStepUpRequired,
			Score:          55,
			Grade:          detectionmodel.RiskGradeHigh,
			ReasonCodes:    []string{"refund_weak_authorization", "agent_refund_without_approval"},
			Recommendation: detectionmodel.RecommendationStepUp,
		},
		results: []detectionmodel.DetectionResult{{
			ID:             "det-order-suspicious-refund-attempt",
			ScenarioID:     "commerce-suspicious-refund-attempt",
			OrderID:        "order-suspicious-refund-attempt",
			Status:         detectionmodel.DetectionStatusStepUpRequired,
			Score:          55,
			Grade:          detectionmodel.RiskGradeHigh,
			Recommendation: detectionmodel.RecommendationStepUp,
		}},
		summary: detectionmodel.DetectionRunSummary{
			TotalByStatus: map[detectionmodel.DetectionStatus]int{
				detectionmodel.DetectionStatusStepUpRequired: 1,
			},
			Total:        1,
			LastResultID: "det-order-suspicious-refund-attempt",
		},
	}

	return NewTrustLabHandler(
		scenarioServiceStub{},
		executionServiceStub{result: executionsvc.ExecutionResult{
			Scenario:       scenario.Scenario{ID: "commerce-clean-agent-assisted-purchase"},
			ReplayCaseRefs: []string{"rc-commerce-clean-agent-assisted-purchase"},
		}},
		&trustServiceStub{contextResp: memory.LoadScenarioContextResponse{ScenarioID: "starter-mandate-review", Context: map[string]any{"record_count": 2}}},
		&replayServiceStub{},
		benchmarkServiceStub{},
		commerceServiceStub{orders: map[string]commerce.Order{
			"order-1": {ID: "order-1", BuyerID: "buyer-1"},
		}},
		eventServiceStub{items: []domainevents.Record{{ID: "evt-1"}}},
		trustDecisionServiceStub{items: map[string]trust.TrustDecision{
			"decision-1": {ID: "decision-1", Outcome: "accepted"},
		}},
		detectionStub,
		TrustLabState{ClawMemBaseURL: "http://127.0.0.1:8088"},
	)
}
