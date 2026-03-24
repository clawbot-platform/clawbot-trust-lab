package scenario

import (
	"context"
	"testing"

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
