# Detection Model

Phase 6 introduces the first deterministic detection baseline for `clawbot-trust-lab`.

The detector is intentionally:

- rule-based
- local and deterministic
- explainable without hidden heuristics
- small enough to benchmark and discuss in an interview

## DetectionResult

Each evaluation stores a `DetectionResult` with:

- `id`
- `scenario_id`
- `order_id`
- `refund_id`
- `trust_decision_refs`
- `replay_case_refs`
- `status`
- `score`
- `grade`
- `triggered_rules`
- `reason_codes`
- `recommendation`
- `evaluated_at`
- `metadata`

## Detection status

- `clean`
- `suspicious`
- `step_up_required`
- `blocked`

## Risk grades

- `low`
- `moderate`
- `high`
- `critical`

## Recommendations

- `allow`
- `review`
- `step_up`
- `block`

## Feature extraction

The baseline detector derives explicit features from:

- commerce state
- transaction and trust events
- trust decisions
- replay history
- `clawmem` scenario context

Current features include:

- `delegated_actor_present`
- `fully_delegated_action`
- `mandate_present`
- `mandate_missing`
- `mandate_expired`
- `provenance_present`
- `provenance_missing`
- `approval_present`
- `approval_missing`
- `refund_requested`
- `refund_requested_by_agent`
- `refund_without_authorization`
- `order_submitted_by_agent`
- `trust_decision_step_up`
- `replay_history_present`
- `memory_context_present`

The stored detection context also includes counts for:

- total related events
- trust events
- trust-decision reason codes
- replay history
- memory status

## Baseline rules

The active rules are:

1. `missing_mandate_delegated_action`
2. `missing_provenance_sensitive_action`
3. `refund_weak_authorization`
4. `agent_refund_without_approval`
5. `prior_step_up_decision`
6. `repeat_suspicious_context`

Each rule emits a `RuleHit` with:

- `rule_id`
- `title`
- `severity`
- `reason`
- `metadata`

## Scoring

The detector uses a bounded additive score:

- each rule contributes a fixed severity value
- the total score drives status, grade, and recommendation

Current thresholds:

- `0-14` => `clean`, `low`, `allow`
- `15-39` => `suspicious`, `moderate`, `review`
- `40-79` => `step_up_required`, `high`, `step_up`
- `80+` => `blocked`, `critical`, `block`

The thresholds are intentionally easy to inspect and change.
