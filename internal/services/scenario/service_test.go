package scenario

import (
	"context"
	"testing"

	"clawbot-trust-lab/internal/domain/actors"
	domainreplay "clawbot-trust-lab/internal/domain/replay"
	domainscenario "clawbot-trust-lab/internal/domain/scenario"
	domaintrust "clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/platform/store"
	commerceSvc "clawbot-trust-lab/internal/services/commerce"
	eventsvc "clawbot-trust-lab/internal/services/events"
	trustsvc "clawbot-trust-lab/internal/services/trust"
)

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

type trustArtifactWriterStub struct{}

func (trustArtifactWriterStub) CreateArtifact(context.Context, domaintrust.CreateArtifactInput) (domaintrust.TrustArtifact, error) {
	return domaintrust.TrustArtifact{ID: "ta-test"}, nil
}

type replayWriterStub struct{}

func (replayWriterStub) CreateCase(context.Context, domainreplay.CreateCaseInput) (domainreplay.ReplayCase, error) {
	return domainreplay.ReplayCase{ID: "rc-test"}, nil
}

func TestExecuteCleanPurchase(t *testing.T) {
	world := store.NewCommerceWorldStore()
	service := NewService(
		scenarioCatalogStub{items: map[string]domainscenario.Scenario{
			"commerce-clean-agent-assisted-purchase": {
				ID:   "commerce-clean-agent-assisted-purchase",
				Name: "Clean Agent-Assisted Purchase",
				Type: domainscenario.ScenarioTypeCommercePurchase,
			},
		}},
		commerceSvc.NewService(world),
		eventsvc.NewService(world),
		trustsvc.NewService(world),
		trustArtifactWriterStub{},
		replayWriterStub{},
	)

	result, err := service.Execute(context.Background(), "commerce-clean-agent-assisted-purchase")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Entities.OrderRefs) != 1 {
		t.Fatalf("expected one order ref, got %#v", result.Entities.OrderRefs)
	}
	if len(result.TrustDecisions) != 1 {
		t.Fatalf("expected one trust decision, got %d", len(result.TrustDecisions))
	}
	if len(result.EventRefs) == 0 {
		t.Fatal("expected event refs")
	}
}

func TestExecuteSuspiciousRefund(t *testing.T) {
	world := store.NewCommerceWorldStore()
	service := NewService(
		scenarioCatalogStub{items: map[string]domainscenario.Scenario{
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

	result, err := service.Execute(context.Background(), "commerce-suspicious-refund-attempt")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if len(result.Entities.RefundRefs) != 1 {
		t.Fatalf("expected one refund ref, got %#v", result.Entities.RefundRefs)
	}
	if !result.TrustDecisions[0].StepUpRequired {
		t.Fatal("expected step up required decision")
	}
}

func TestListScenariosReturnsCatalogItems(t *testing.T) {
	service := NewService(
		scenarioCatalogStub{items: map[string]domainscenario.Scenario{
			"scenario-a": {ID: "scenario-a"},
			"scenario-b": {ID: "scenario-b"},
		}},
		commerceSvc.NewService(store.NewCommerceWorldStore()),
		eventsvc.NewService(store.NewCommerceWorldStore()),
		trustsvc.NewService(store.NewCommerceWorldStore()),
		trustArtifactWriterStub{},
		replayWriterStub{},
	)

	scenarios := service.ListScenarios()
	if len(scenarios) != 2 {
		t.Fatalf("expected 2 scenarios, got %d", len(scenarios))
	}
}

func TestGetExecutionResultAndOrderLookup(t *testing.T) {
	world := store.NewCommerceWorldStore()
	service := NewService(
		scenarioCatalogStub{items: map[string]domainscenario.Scenario{
			"commerce-clean-agent-assisted-purchase": {
				ID:   "commerce-clean-agent-assisted-purchase",
				Name: "Clean Agent-Assisted Purchase",
				Type: domainscenario.ScenarioTypeCommercePurchase,
			},
		}},
		commerceSvc.NewService(world),
		eventsvc.NewService(world),
		trustsvc.NewService(world),
		trustArtifactWriterStub{},
		replayWriterStub{},
	)

	if _, err := service.GetExecutionResult("commerce-clean-agent-assisted-purchase"); err == nil {
		t.Fatal("expected GetExecutionResult before execution to fail")
	}
	if _, err := service.GetExecutionResultByOrderID("order-missing"); err == nil {
		t.Fatal("expected GetExecutionResultByOrderID for unknown order to fail")
	}

	result, err := service.Execute(context.Background(), "commerce-clean-agent-assisted-purchase")
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	stored, err := service.GetExecutionResult("commerce-clean-agent-assisted-purchase")
	if err != nil {
		t.Fatalf("GetExecutionResult() error = %v", err)
	}
	if stored.Scenario.ID != result.Scenario.ID {
		t.Fatalf("expected stored result for %s, got %#v", result.Scenario.ID, stored)
	}

	byOrder, err := service.GetExecutionResultByOrderID(result.Entities.OrderRefs[0])
	if err != nil {
		t.Fatalf("GetExecutionResultByOrderID() error = %v", err)
	}
	if byOrder.Scenario.ID != result.Scenario.ID {
		t.Fatalf("expected order lookup to resolve %s, got %#v", result.Scenario.ID, byOrder)
	}
}

func TestBlueprintForScenarioSupportsAliasesAndRejectsUnknownScenario(t *testing.T) {
	item := domainscenario.Scenario{
		ID: "commerce-v1-weakened-provenance",
	}

	blueprint, err := blueprintForScenario(item)
	if err != nil {
		t.Fatalf("blueprintForScenario() error = %v", err)
	}
	if blueprint.FlowKind != flowKindPurchase {
		t.Fatalf("expected purchase flow for alias scenario, got %#v", blueprint)
	}
	if weak, ok := blueprint.SignalContext["provenance_low_confidence"].(bool); !ok || !weak {
		t.Fatalf("expected weak provenance alias to carry low-confidence signal, got %#v", blueprint.SignalContext)
	}

	if _, err := blueprintForScenario(domainscenario.Scenario{ID: "commerce-unknown"}); err == nil {
		t.Fatal("expected unknown scenario blueprint lookup to fail")
	}
}

func TestSignalContextHandlesTierAvailability(t *testing.T) {
	humanScenario := domainscenario.Scenario{
		ID: "commerce-h2-human-refund-valid-history",
		FeatureModel: domainscenario.FeatureTierModel{
			TierA: []string{"amount"},
			TierB: []string{"approval_history"},
		},
	}
	humanPlan, err := blueprintForScenario(humanScenario)
	if err != nil {
		t.Fatalf("blueprintForScenario() human error = %v", err)
	}
	humanContext := signalContext(humanScenario, seedWorld(), humanPlan)
	if got, _ := humanContext["tier_c_available"].(bool); got {
		t.Fatalf("expected human baseline to work without tier C, got %#v", humanContext)
	}
	if got, _ := humanContext["tier_c_used_in_scenario"].(bool); got {
		t.Fatalf("expected human baseline not to mark tier C used, got %#v", humanContext)
	}

	agentScenario := domainscenario.Scenario{
		ID: "commerce-a1-agent-assisted-purchase-valid-controls",
		FeatureModel: domainscenario.FeatureTierModel{
			TierA: []string{"amount"},
			TierB: []string{"buyer_history"},
			TierC: []string{"mandate_status", "provenance_confidence"},
		},
	}
	agentPlan, err := blueprintForScenario(agentScenario)
	if err != nil {
		t.Fatalf("blueprintForScenario() agent error = %v", err)
	}
	agentContext := signalContext(agentScenario, seedWorld(), agentPlan)
	if got, _ := agentContext["tier_c_available"].(bool); !got {
		t.Fatalf("expected agent-assisted scenario to expose tier C availability, got %#v", agentContext)
	}
	if got, _ := agentContext["evaluation_mode"].(string); got != "shadow" {
		t.Fatalf("expected production-bridge shadow mode signal, got %#v", agentContext)
	}
}

func TestBlueprintForScenarioCoversMultiplePhase9Families(t *testing.T) {
	cases := []struct {
		id       string
		flowKind string
		check    func(t *testing.T, blueprint scenarioBlueprint)
	}{
		{
			id:       "commerce-h1-direct-human-purchase",
			flowKind: flowKindPurchase,
			check: func(t *testing.T, blueprint scenarioBlueprint) {
				t.Helper()
				if blueprint.DelegationMode != actors.DelegationModeDirectHuman {
					t.Fatalf("expected direct human delegation, got %#v", blueprint)
				}
			},
		},
		{
			id:       "commerce-a2-fully-delegated-replenishment-purchase",
			flowKind: flowKindPurchase,
			check: func(t *testing.T, blueprint scenarioBlueprint) {
				t.Helper()
				if blueprint.DelegationMode != actors.DelegationModeFullyDelegated {
					t.Fatalf("expected fully delegated purchase, got %#v", blueprint)
				}
			},
		},
		{
			id:       "commerce-s5-merchant-scope-drift-delegated-action",
			flowKind: flowKindPurchase,
			check: func(t *testing.T, blueprint scenarioBlueprint) {
				t.Helper()
				if matched, _ := blueprint.SignalContext["merchant_scope_match"].(bool); matched {
					t.Fatalf("expected merchant scope drift signal, got %#v", blueprint.SignalContext)
				}
			},
		},
		{
			id:       "commerce-v4-actor-switch-human-to-agent",
			flowKind: flowKindRefund,
			check: func(t *testing.T, blueprint scenarioBlueprint) {
				t.Helper()
				if switched, _ := blueprint.SignalContext["actor_switch_to_agent"].(bool); !switched {
					t.Fatalf("expected actor-switch challenger signal, got %#v", blueprint.SignalContext)
				}
			},
		},
		{
			id:       "commerce-v7-high-value-delegated-purchase",
			flowKind: flowKindPurchase,
			check: func(t *testing.T, blueprint scenarioBlueprint) {
				t.Helper()
				if threshold, _ := blueprint.SignalContext["high_value_threshold"].(int); threshold == 0 {
					t.Fatalf("expected high-value threshold signal, got %#v", blueprint.SignalContext)
				}
			},
		},
	}

	for _, tc := range cases {
		blueprint, err := blueprintForScenario(domainscenario.Scenario{ID: tc.id})
		if err != nil {
			t.Fatalf("blueprintForScenario(%s) error = %v", tc.id, err)
		}
		if blueprint.FlowKind != tc.flowKind {
			t.Fatalf("expected %s flow for %s, got %#v", tc.flowKind, tc.id, blueprint)
		}
		tc.check(t, blueprint)
	}
}
