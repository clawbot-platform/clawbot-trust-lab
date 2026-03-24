package events

import "time"

type EventCategory string
type TransactionEventType string
type TrustEventType string

const (
	EventCategoryTransaction EventCategory = "transaction"
	EventCategoryTrust       EventCategory = "trust"
)

const (
	TransactionEventOrderCreated           TransactionEventType = "order_created"
	TransactionEventPaymentAuthorized      TransactionEventType = "payment_authorized"
	TransactionEventRefundRequested        TransactionEventType = "refund_requested"
	TransactionEventRefundDecisionRecorded TransactionEventType = "refund_decision_recorded"
)

const (
	TrustEventOrderSubmittedByAgent TrustEventType = "order_submitted_by_agent"
	TrustEventMandateChecked        TrustEventType = "mandate_checked"
	TrustEventProvenanceAttached    TrustEventType = "provenance_attached"
	TrustEventTrustDecisionRecorded TrustEventType = "trust_decision_recorded"
	TrustEventApprovalRecorded      TrustEventType = "approval_recorded"
)

type Record struct {
	ID         string         `json:"id"`
	Category   EventCategory  `json:"category"`
	EventType  string         `json:"event_type"`
	EntityType string         `json:"entity_type"`
	EntityID   string         `json:"entity_id"`
	ScenarioID string         `json:"scenario_id"`
	ActorID    string         `json:"actor_id"`
	OccurredAt time.Time      `json:"occurred_at"`
	Metadata   map[string]any `json:"metadata"`
}

type TransactionEvent struct {
	Record
}

type TrustEvent struct {
	Record
}
