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
func (benchmarkServiceStub) Status() map[string]any {
	return map[string]any{"registrations": 1, "last_status": "accepted_stub"}
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

func TestSystemHandlerHealth(t *testing.T) {
	handler := NewSystemHandler(func(context.Context) error { return nil }, version.Current())
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	handler.Health(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
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
