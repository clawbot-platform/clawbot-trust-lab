# Phase 5 Commerce World

Phase 5 is the first world-model phase for `clawbot-trust-lab`.

It does not implement fraud scoring yet. Instead, it creates the deterministic baseline world that later phases can score, benchmark, mutate, replay, and analyze.

## Why this phase exists

Later phases need a coherent world to observe:

- actors and delegation
- orders, payments, and refunds
- mandate and provenance surfaces
- trust decisions
- transaction and trust events

Without that baseline, a detection phase would only score placeholders.

## Phase 5 baseline flows

### Flow A: clean agent-assisted purchase

- an agent submits an order on behalf of a buyer
- a valid mandate exists
- provenance is attached
- the order is accepted
- payment is authorized
- a trust decision is recorded as acceptable

### Flow B: suspicious refund attempt

- a refund is requested by an agent
- mandate coverage is expired
- provenance is weak
- human approval is missing
- the trust decision requires step-up
- the refund path records the failure clearly in the event stream

## Why this is event-first enough

The important outputs are not just the commerce entities. They are:

- trust decisions
- transaction events
- trust events
- replay cases
- memory writes

That makes Phase 5 the correct base for Phase 6 detection and Phase 7 Red Queen style adversarial work.
