package detection

import (
	"context"
	"errors"
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
			"commerce-h2-human-refund-valid-history": {
				ID:   "commerce-h2-human-refund-valid-history",
				Name: "Human Refund with Valid History",
				Type: domainscenario.ScenarioTypeCommerceRefundReview,
				FeatureModel: domainscenario.FeatureTierModel{
					TierA: []string{"refund_indicator", "order_age", "amount"},
					TierB: []string{"historical_refund_rate", "approval_history"},
					TierC: []string{},
				},
			},
			"commerce-a3-agent-assisted-refund-approval-evidence": {
				ID:   "commerce-a3-agent-assisted-refund-approval-evidence",
				Name: "Agent-Assisted Refund with Approval Evidence",
				Type: domainscenario.ScenarioTypeCommerceRefundReview,
				FeatureModel: domainscenario.FeatureTierModel{
					TierA: []string{"refund_indicator", "amount", "delegated_indicator"},
					TierB: []string{"historical_refund_rate", "approval_history", "recent_attempt_count"},
					TierC: []string{"mandate_status", "provenance_confidence", "delegation_mode"},
				},
			},
			"commerce-s4-repeated-agent-refund-attempts": {
				ID:   "commerce-s4-repeated-agent-refund-attempts",
				Name: "Repeated Agent Refund Attempts",
				Type: domainscenario.ScenarioTypeCommerceRefundReview,
				FeatureModel: domainscenario.FeatureTierModel{
					TierA: []string{"refund_indicator", "amount", "delegated_indicator"},
					TierB: []string{"repeat_attempt_count", "historical_refund_rate", "prior_review_outcomes"},
					TierC: []string{"mandate_status", "approval_evidence", "delegation_mode"},
				},
			},
			"commerce-clean-agent-assisted-purchase": {
				ID:   "commerce-clean-agent-assisted-purchase",
				Name: "Clean Agent-Assisted Purchase",
				Type: domainscenario.ScenarioTypeCommercePurchase,
				FeatureModel: domainscenario.FeatureTierModel{
					TierA: []string{"amount", "merchant_category"},
					TierB: []string{"buyer_history", "recent_attempt_count"},
					TierC: []string{"mandate_status", "provenance_confidence"},
				},
			},
			"commerce-suspicious-refund-attempt": {
				ID:   "commerce-suspicious-refund-attempt",
				Name: "Suspicious Refund Attempt",
				Type: domainscenario.ScenarioTypeCommerceRefundReview,
				FeatureModel: domainscenario.FeatureTierModel{
					TierA: []string{"refund_indicator", "amount"},
					TierB: []string{"historical_refund_rate", "repeat_attempt_count"},
					TierC: []string{"mandate_status", "approval_evidence"},
				},
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

func TestEvaluateHumanRefundWithValidHistoryReturnsClean(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-h2-human-refund-valid-history"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-h2-human-refund-valid-history",
			Context:    map[string]any{"record_count": 0},
		}},
		store.NewDetectionStore(),
	)

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-h2-human-refund-valid-history"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Status != detectionmodel.DetectionStatusClean {
		t.Fatalf("expected clean status, got %s", result.Status)
	}
	if len(result.TriggeredRules) != 0 {
		t.Fatalf("expected no triggered rules, got %#v", result.TriggeredRules)
	}
}

func TestEvaluateAgentRefundWithApprovalEvidenceReturnsClean(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-a3-agent-assisted-refund-approval-evidence"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-a3-agent-assisted-refund-approval-evidence",
			Context:    map[string]any{"record_count": 1},
		}},
		store.NewDetectionStore(),
	)

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-a3-agent-assisted-refund-approval-evidence"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Status != detectionmodel.DetectionStatusClean {
		t.Fatalf("expected clean status, got %s", result.Status)
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
	if result.Score < 75 {
		t.Fatalf("expected score at least 75, got %d", result.Score)
	}
	assertReasonCode(t, result.ReasonCodes, "agent_refund_without_approval")
	assertReasonCode(t, result.ReasonCodes, "missing_mandate_delegated_action")
	assertReasonCode(t, result.ReasonCodes, "prior_step_up_decision")
	assertReasonCode(t, result.ReasonCodes, "refund_weak_authorization")
}

func TestEvaluateRepeatedAgentRefundAttemptsNoLongerReturnsClean(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-s4-repeated-agent-refund-attempts"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-s4-repeated-agent-refund-attempts",
			Context:    map[string]any{"record_count": 0},
		}},
		store.NewDetectionStore(),
	)

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-s4-repeated-agent-refund-attempts"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if result.Status == detectionmodel.DetectionStatusClean {
		t.Fatalf("expected repeated agent refund attempts to be escalated, got %#v", result)
	}
	assertReasonCode(t, result.ReasonCodes, "repeat_suspicious_context")
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

func TestEvaluateTierCRemainsOptionalForHumanBaseline(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-h2-human-refund-valid-history"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-h2-human-refund-valid-history",
			Context:    map[string]any{"record_count": 0},
		}},
		store.NewDetectionStore(),
	)

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-h2-human-refund-valid-history"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	contextData, ok := result.Metadata["context"].(detectionmodel.DetectionContext)
	if !ok {
		t.Fatalf("expected typed context metadata, got %#v", result.Metadata["context"])
	}
	if !contextData.TierProfile.TierAAvailable || !contextData.TierProfile.TierBAvailable {
		t.Fatalf("expected tier A and B to be available, got %#v", contextData.TierProfile)
	}
	if contextData.TierProfile.TierCUsed {
		t.Fatalf("expected tier C to remain optional for this baseline, got %#v", contextData.TierProfile)
	}
}

func TestBuildContextRecognizesFloatMemoryCountAndTierCUse(t *testing.T) {
	execution, world := newExecutionService()
	run, err := execution.Execute(context.Background(), "commerce-clean-agent-assisted-purchase")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{items: []domainreplay.ReplayCase{{ID: "rc-1", ScenarioID: run.Scenario.ID}}},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: run.Scenario.ID,
			Context:    map[string]any{"record_count": 2.0},
		}},
		store.NewDetectionStore(),
	)

	contextData := service.buildContext(context.Background(), run)
	if !contextData.MemoryContextPresent || contextData.MemoryStatus != "ok" {
		t.Fatalf("expected float record_count to enable memory context, got %#v", contextData)
	}
	if !contextData.TierProfile.TierCAvailable || !contextData.TierProfile.TierCUsed {
		t.Fatalf("expected delegated purchase to use tier C when available, got %#v", contextData.TierProfile)
	}
}

func TestBuildContextDegradesWhenMemoryUnavailable(t *testing.T) {
	execution, world := newExecutionService()
	run, err := execution.Execute(context.Background(), "commerce-h2-human-refund-valid-history")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{},
		memoryClientStub{contextErr: errors.New("unavailable")},
		store.NewDetectionStore(),
	)

	contextData := service.buildContext(context.Background(), run)
	if contextData.MemoryContextPresent || contextData.MemoryStatus != "degraded" {
		t.Fatalf("expected degraded memory context, got %#v", contextData)
	}
}

func TestEvaluateRulesCoversScopeDriftHighValueAndActorSwitch(t *testing.T) {
	service := NewService(nil, nil, nil, nil, store.NewDetectionStore())

	hits := service.evaluateRules(detectionmodel.DetectionContext{
		ScenarioID: "commerce-v7-high-value-delegated-purchase",
		Features: map[string]bool{
			"delegated_actor_present":       true,
			"merchant_scope_drift":          true,
			"high_value_delegated_purchase": true,
			"actor_switch_to_agent":         true,
		},
		Signals: map[string]any{"buyer_spend_baseline": int64(4200)},
	})

	if len(hits) != 3 {
		t.Fatalf("expected 3 rule hits, got %#v", hits)
	}
	assertReasonCode(t, []string{hits[0].RuleID, hits[1].RuleID, hits[2].RuleID}, ruleActorSwitchSensitiveAction)
	assertReasonCode(t, []string{hits[0].RuleID, hits[1].RuleID, hits[2].RuleID}, ruleHighValueDelegatedPurchase)
	assertReasonCode(t, []string{hits[0].RuleID, hits[1].RuleID, hits[2].RuleID}, ruleMerchantScopeDriftDelegated)
}

func TestResolveExecutionRequiresInput(t *testing.T) {
	service := NewService(nil, nil, nil, nil, store.NewDetectionStore())
	if _, err := service.resolveExecution(EvaluateInput{}); err == nil {
		t.Fatal("expected resolveExecution with no identifiers to fail")
	}
}

func TestRulesAndDeriveOutcomeRemainStructured(t *testing.T) {
	service := NewService(nil, nil, nil, nil, store.NewDetectionStore())
	rules := service.Rules()
	if len(rules) != len(ruleCatalog) {
		t.Fatalf("expected %d rules, got %d", len(ruleCatalog), len(rules))
	}

	score, status, grade, recommendation := deriveOutcome([]detectionmodel.RuleHit{
		{Severity: 25},
		{Severity: 20},
	})
	if score != 45 || status != detectionmodel.DetectionStatusStepUpRequired || grade != detectionmodel.RiskGradeHigh || recommendation != detectionmodel.RecommendationStepUp {
		t.Fatalf("unexpected derived outcome: score=%d status=%s grade=%s recommendation=%s", score, status, grade, recommendation)
	}
}

func TestListGetAndSummaryReflectStoredResults(t *testing.T) {
	execution, world := newExecutionService()
	if _, err := execution.Execute(context.Background(), "commerce-suspicious-refund-attempt"); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	service := NewService(
		world,
		execution,
		replayReaderStub{},
		memoryClientStub{contextResponse: memory.LoadScenarioContextResponse{
			ScenarioID: "commerce-suspicious-refund-attempt",
			Context:    map[string]any{"record_count": 1},
		}},
		store.NewDetectionStore(),
	)

	result, err := service.Evaluate(context.Background(), EvaluateInput{ScenarioID: "commerce-suspicious-refund-attempt"})
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}

	list := service.ListResults()
	if len(list) != 1 || list[0].ID != result.ID {
		t.Fatalf("expected stored result in list, got %#v", list)
	}

	stored, err := service.GetResult(result.ID)
	if err != nil {
		t.Fatalf("GetResult() error = %v", err)
	}
	if stored.ID != result.ID {
		t.Fatalf("expected stored id %s, got %#v", result.ID, stored)
	}

	summary := service.Summary()
	if summary.Total != 1 || summary.LastResultID != result.ID {
		t.Fatalf("expected detection summary to reflect stored result, got %#v", summary)
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
