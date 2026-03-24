package detection

import (
	"context"
	"testing"
	"time"

	"clawbot-trust-lab/internal/clients/memory"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	domainreplay "clawbot-trust-lab/internal/domain/replay"
	domainscenario "clawbot-trust-lab/internal/domain/scenario"
	domaintrust "clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/platform/store"
	commerceSvc "clawbot-trust-lab/internal/services/commerce"
	eventsvc "clawbot-trust-lab/internal/services/events"
	scenariosvc "clawbot-trust-lab/internal/services/scenario"
	trustsvc "clawbot-trust-lab/internal/services/trust"
)

type memoryClientStub struct {
	contextResponse memory.LoadScenarioContextResponse
	contextErr      error
}

func (s memoryClientStub) Health(context.Context) error { return nil }

func (s memoryClientStub) StoreReplayCase(context.Context, memory.StoreReplayCaseRequest) error {
	return nil
}

func (s memoryClientStub) FetchSimilarCases(context.Context, memory.FetchSimilarCasesRequest) (memory.FetchSimilarCasesResponse, error) {
	return memory.FetchSimilarCasesResponse{}, nil
}

func (s memoryClientStub) StoreTrustArtifact(context.Context, memory.StoreTrustArtifactRequest) error {
	return nil
}

func (s memoryClientStub) LoadScenarioContext(context.Context, memory.LoadScenarioContextRequest) (memory.LoadScenarioContextResponse, error) {
	if s.contextErr != nil {
		return memory.LoadScenarioContextResponse{}, s.contextErr
	}
	return s.contextResponse, nil
}

type replayReaderStub struct {
	items []domainreplay.ReplayCase
}

func (s replayReaderStub) ListCases() []domainreplay.ReplayCase {
	return append([]domainreplay.ReplayCase(nil), s.items...)
}

type trustArtifactWriterStub struct{}

func (trustArtifactWriterStub) CreateArtifact(context.Context, domaintrust.CreateArtifactInput) (domaintrust.TrustArtifact, error) {
	return domaintrust.TrustArtifact{ID: "ta-test"}, nil
}

type replayWriterStub struct{}

func (replayWriterStub) CreateCase(context.Context, domainreplay.CreateCaseInput) (domainreplay.ReplayCase, error) {
	return domainreplay.ReplayCase{ID: "rc-test"}, nil
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

func newExecutionService() (*scenariosvc.Service, *store.CommerceWorldStore) {
	world := store.NewCommerceWorldStore()
	execution := scenariosvc.NewService(
		scenarioCatalogStub{items: map[string]domainscenario.Scenario{
			"commerce-clean-agent-assisted-purchase": {
				ID:   "commerce-clean-agent-assisted-purchase",
				Name: "Clean Agent-Assisted Purchase",
				Type: domainscenario.ScenarioTypeCommercePurchase,
			},
			"commerce-suspicious-refund-attempt": {
				ID:   "commerce-suspicious-refund-attempt",
				Name: "Suspicious Refund Attempt",
				Type: domainscenario.ScenarioTypeCommerceRefundReview,
			},
		}},
		commerceSvc.NewService(world),
		eventsvc.NewService(world),
		trustsvc.NewService(world),
		trustArtifactWriterStub{},
		replayWriterStub{},
	)
	return execution, world
}

func TestEvaluateCleanScenarioReturnsClean(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-clean-agent-assisted-purchase"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{items: []domainreplay.ReplayCase{
			{ID: "rc-commerce-clean-agent-assisted-purchase", ScenarioID: "commerce-clean-agent-assisted-purchase"},
		}},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-clean-agent-assisted-purchase",
			Context:    map[string]any{"record_count": 2},
		}},
		store.NewDetectionStore(),
	)
	service.now = func() time.Time { return time.Date(2026, 3, 24, 12, 0, 0, 0, time.UTC) }

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-clean-agent-assisted-purchase"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Status != detectionmodel.DetectionStatusClean {
		t.Fatalf("expected clean status, got %s", result.Status)
	}
	if result.Recommendation != detectionmodel.RecommendationAllow {
		t.Fatalf("expected allow recommendation, got %s", result.Recommendation)
	}
	if len(result.TriggeredRules) != 0 {
		t.Fatalf("expected no triggered rules, got %#v", result.TriggeredRules)
	}
}

func TestEvaluateSuspiciousScenarioTriggersRuleHits(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-suspicious-refund-attempt"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{items: []domainreplay.ReplayCase{
			{ID: "rc-commerce-suspicious-refund-attempt", ScenarioID: "commerce-suspicious-refund-attempt"},
		}},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-suspicious-refund-attempt",
			Context:    map[string]any{"record_count": 1},
		}},
		store.NewDetectionStore(),
	)
	service.now = func() time.Time { return time.Date(2026, 3, 24, 12, 5, 0, 0, time.UTC) }

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-suspicious-refund-attempt"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Status != detectionmodel.DetectionStatusStepUpRequired {
		t.Fatalf("expected step_up_required status, got %s", result.Status)
	}
	if result.Recommendation != detectionmodel.RecommendationStepUp {
		t.Fatalf("expected step_up recommendation, got %s", result.Recommendation)
	}
	if result.Score != 75 {
		t.Fatalf("expected score 75, got %d", result.Score)
	}
	assertReasonCode(t, result.ReasonCodes, "agent_refund_without_approval")
	assertReasonCode(t, result.ReasonCodes, "missing_mandate_delegated_action")
	assertReasonCode(t, result.ReasonCodes, "prior_step_up_decision")
	assertReasonCode(t, result.ReasonCodes, "refund_weak_authorization")
}

func TestEvaluateSuspiciousScenarioIncludesRepeatContextRule(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-suspicious-refund-attempt"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{items: []domainreplay.ReplayCase{
			{ID: "rc-commerce-suspicious-refund-attempt-1", ScenarioID: "commerce-suspicious-refund-attempt"},
			{ID: "rc-commerce-suspicious-refund-attempt-2", ScenarioID: "commerce-suspicious-refund-attempt"},
		}},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-suspicious-refund-attempt",
			Context:    map[string]any{"record_count": 3},
		}},
		store.NewDetectionStore(),
	)

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-suspicious-refund-attempt"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	assertReasonCode(t, result.ReasonCodes, "repeat_suspicious_context")
}

func TestEvaluateByOrderIDResolvesExecution(t *testing.T) {
	execution, world := newExecutionService()
	run, err := execution.Execute(context.Background(), "commerce-clean-agent-assisted-purchase")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-clean-agent-assisted-purchase",
			Context:    map[string]any{"record_count": 0},
		}},
		store.NewDetectionStore(),
	)

	result, err := service.Evaluate(context.Background(), EvaluateInput{OrderID: run.Entities.OrderRefs[0]})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.OrderID != run.Entities.OrderRefs[0] {
		t.Fatalf("expected order id %s, got %s", run.Entities.OrderRefs[0], result.OrderID)
	}
}

func assertReasonCode(t *testing.T, codes []string, want string) {
	t.Helper()
	for _, code := range codes {
		if code == want {
			return
		}
	}
	t.Fatalf("expected reason code %s in %#v", want, codes)
}
