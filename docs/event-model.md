# Event Model

Phase 5 treats events as first-class outputs.

## Event categories

- transaction events
- trust events

## Transaction event types

- `order_created`
- `payment_authorized`
- `refund_requested`
- `refund_decision_recorded`

## Trust event types

- `order_submitted_by_agent`
- `mandate_checked`
- `provenance_attached`
- `trust_decision_recorded`
- `approval_recorded`

## Shared event shape

All events are append-only records with:

- `id`
- `category`
- `event_type`
- `entity_type`
- `entity_id`
- `scenario_id`
- `actor_id`
- `occurred_at`
- `metadata`

This shape is small enough to stay explainable and rich enough to power later replay, benchmark, and detection work.
