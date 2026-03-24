package events

import (
	"testing"
	"time"

	domainevents "clawbot-trust-lab/internal/domain/events"
	"clawbot-trust-lab/internal/platform/store"
)

func TestRecordEvents(t *testing.T) {
	world := store.NewCommerceWorldStore()
	service := NewService(world)

	service.RecordTransaction("evt-1", domainevents.TransactionEventOrderCreated, "order", "order-1", "scenario-1", "actor-1", time.Date(2026, 3, 24, 9, 0, 0, 0, time.UTC), nil)
	service.RecordTrust("evt-2", domainevents.TrustEventTrustDecisionRecorded, "trust_decision", "decision-1", "scenario-1", "actor-1", time.Date(2026, 3, 24, 9, 1, 0, 0, time.UTC), nil)

	if len(service.ListEvents()) != 2 {
		t.Fatalf("expected 2 events, got %d", len(service.ListEvents()))
	}
}
