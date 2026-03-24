package trust

import (
	"testing"
	"time"

	domaintrust "clawbot-trust-lab/internal/domain/trust"
	"clawbot-trust-lab/internal/platform/store"
)

func TestRecordDecision(t *testing.T) {
	world := store.NewCommerceWorldStore()
	service := NewService(world)

	decision, err := service.RecordDecision(domaintrust.TrustDecision{
		ID:            "decision-1",
		EntityType:    "order",
		EntityID:      "order-1",
		Outcome:       "accepted",
		ReasonCodes:   []string{"active_mandate"},
		MandateRef:    "mandate-1",
		ProvenanceRef: "prov-1",
		RecordedAt:    time.Date(2026, 3, 24, 9, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("RecordDecision() error = %v", err)
	}
	if decision.ID != "decision-1" {
		t.Fatalf("unexpected decision id: %s", decision.ID)
	}
}
