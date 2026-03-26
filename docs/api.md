# API

Phase 9 keeps the earlier trust/replay APIs, the deterministic commerce-world execution layer, the explainable detection baseline, the Red Queen MVP benchmark loop, and adds production-bridge recommendations, scheduled rounds, and long-run trend summaries.

## Demo walkthrough

If you want to show the system quickly, use the API in this order.

Inspect what the harness will test:

```bash
curl http://127.0.0.1:8090/api/v1/scenarios/packs/commerce-pack
curl http://127.0.0.1:8090/api/v1/scenarios/packs/challenger-pack
```

Run one benchmark round:

```bash
ROUND_ID=$(curl -s -X POST http://127.0.0.1:8090/api/v1/benchmark/rounds/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce"}' | jq -r '.data.id')
echo "$ROUND_ID"
```

Then inspect the operator-facing outputs:

```bash
curl http://127.0.0.1:8090/api/v1/operator/rounds
curl http://127.0.0.1:8090/api/v1/operator/rounds/$ROUND_ID
curl http://127.0.0.1:8090/api/v1/operator/promotions
curl http://127.0.0.1:8090/api/v1/operator/recommendations
curl http://127.0.0.1:8090/api/v1/operator/reports/$ROUND_ID
curl http://127.0.0.1:8090/api/v1/operator/trends/summary
```

Run a short scheduled loop for the “week-long demo in miniature” story:

```bash
curl -X POST http://127.0.0.1:8090/api/v1/benchmark/scheduler/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce","interval":"1s","max_runs":2,"dry_run":true}'
```

That sequence is usually enough to explain the full system:

- the catalog defines realistic stable and challenger cases
- the round runner evaluates them in shadow mode
- promotions and replay preserve meaningful failures
- recommendations translate technical findings into operator actions
- long-run trends show whether the system is improving or regressing

The rest of this document is the fuller reference.

## System

- `GET /healthz`
- `GET /readyz`
- `GET /version`

## Scenarios

- `GET /api/v1/scenarios`
- `POST /api/v1/scenarios/execute`
- `GET /api/v1/scenarios/types`
- `GET /api/v1/scenarios/packs`
- `GET /api/v1/scenarios/packs/{id}`

### Canonical IDs and legacy aliases

Phase 9 uses canonical scenario ids for the public catalog and examples. Older ids are still accepted so historical rounds, replay entries, and earlier demo scripts continue to work.

Prefer these canonical ids in new integrations:

- `commerce-a1-agent-assisted-purchase-valid-controls`
- `commerce-s1-refund-weak-authorization`
- `commerce-s2-delegated-purchase-weak-provenance`
- `commerce-v2-expired-inactive-mandate`
- `commerce-s3-approval-removed-after-authorization`

Legacy aliases still accepted by the runtime:

- `commerce-clean-agent-assisted-purchase` -> `commerce-a1-agent-assisted-purchase-valid-controls`
- `commerce-suspicious-refund-attempt` -> `commerce-s1-refund-weak-authorization`
- `commerce-challenger-weakened-provenance-purchase` -> `commerce-s2-delegated-purchase-weak-provenance`
- `commerce-challenger-expired-mandate-purchase` -> `commerce-v2-expired-inactive-mandate`
- `commerce-challenger-approval-removed-refund` -> `commerce-s3-approval-removed-after-authorization`

Execute example:

```json
{
  "scenario_id": "commerce-a1-agent-assisted-purchase-valid-controls"
}
```

Execution responses summarize:
- the scenario
- created or updated entity refs
- trust decisions
- replay case refs
- memory write outcomes

## Detection

- `POST /api/v1/detection/evaluate`
- `GET /api/v1/detection/results`
- `GET /api/v1/detection/results/{id}`
- `GET /api/v1/detection/rules`
- `GET /api/v1/detection/summary`

Evaluate example:

```json
{
  "scenario_id": "commerce-s1-refund-weak-authorization"
}
```

You can also evaluate by `order_id`:

```json
{
  "order_id": "order-a1-agent-assisted-purchase-valid-controls"
}
```

Detection responses include:
- scenario and entity refs
- status, score, grade, and recommendation
- triggered baseline rules
- reason codes
- linked trust decision and replay refs
- a serialized detection context in `metadata.context`
- `metadata.tier_profile`, showing whether Tier A, Tier B, and optional Tier C fields were available and whether Tier C was actually used

`GET /api/v1/detection/rules` lists the active baseline rules:
- `missing_mandate_delegated_action`
- `missing_provenance_sensitive_action`
- `refund_weak_authorization`
- `agent_refund_without_approval`
- `prior_step_up_decision`
- `repeat_suspicious_context`
- `merchant_scope_drift_delegated_action`
- `high_value_delegated_purchase`
- `actor_switch_sensitive_action`

`GET /api/v1/detection/summary` returns a small aggregate view across stored results, including totals by status and the last result id.

## Orders

- `GET /api/v1/orders`
- `GET /api/v1/orders/{id}`

## Events

- `GET /api/v1/events`

## Trust

- `POST /api/v1/trust/artifacts`
- `GET /api/v1/trust/artifacts`
- `GET /api/v1/trust/status`
- `GET /api/v1/trust/decisions`
- `GET /api/v1/trust/decisions/{id}`

`POST /api/v1/trust/artifacts` now requires a successful `clawmem` write. If `clawmem` is unavailable, the endpoint returns `502 Bad Gateway`.

`GET /api/v1/trust/status` accepts an optional `scenario_id` query parameter. When present, the handler attempts to load `clawmem` scenario context and includes either `memory_status=ok` plus `memory_context`, or `memory_status=degraded` plus `memory_error`.

Example:

```json
{
  "scenario_id": "starter-mandate-review"
}
```

## Replay

- `POST /api/v1/replay/cases`
- `GET /api/v1/replay/cases`
- `GET /api/v1/replay/status`

`POST /api/v1/replay/cases` now requires a successful `clawmem` write. If `clawmem` is unavailable, the endpoint returns `502 Bad Gateway`.

`GET /api/v1/replay/status` accepts an optional `scenario_id` query parameter. When present, the handler attempts to load replay memory context from `clawmem` and includes either `memory_status=ok` plus `similar_cases`, or `memory_status=degraded` plus `memory_error`.

Example:

```json
{
  "scenario_id": "starter-mandate-review",
  "trust_artifact_refs": ["ta-starter-mandate-review"],
  "benchmark_round_ref": "bench-round-1",
  "outcome_summary": "Replay matched expected artifact flow",
  "promotion_recommendation": "promote",
  "promotion_reason": "Baseline outcome is explainable"
}
```

## Benchmark

- `POST /api/v1/benchmark/rounds/register`
- `POST /api/v1/benchmark/rounds/run`
- `GET /api/v1/benchmark/rounds`
- `GET /api/v1/benchmark/rounds/{id}`
- `GET /api/v1/benchmark/rounds/{id}/summary`
- `GET /api/v1/benchmark/rounds/{id}/promotions`
- `GET /api/v1/benchmark/rounds/{id}/delta`
- `GET /api/v1/benchmark/rounds/{id}/reports`
- `GET /api/v1/benchmark/recommendations`
- `GET /api/v1/benchmark/recommendations/{id}`
- `GET /api/v1/benchmark/trends/summary`
- `GET /api/v1/benchmark/scheduler/status`
- `POST /api/v1/benchmark/scheduler/run`
- `GET /api/v1/benchmark/rounds/status`
- `GET /api/v1/benchmark/status`

### Recommendation object

Returned by:

- `GET /api/v1/benchmark/recommendations`
- `GET /api/v1/benchmark/recommendations/{id}`
- `GET /api/v1/operator/recommendations`
- `GET /api/v1/operator/recommendations/{id}`

Schema:

- `id`: stable recommendation id
- `type`: one of `add_to_replay_stable_set`, `tighten_refund_review_rule`, `require_step_up_for_delegated_refunds`, `require_provenance_for_delegated_purchase`, `investigate_repeat_refund_pattern`, `monitor_in_shadow_mode`
- `rationale`: short explanation grounded in benchmark output
- `priority`: `low`, `moderate`, or `high`
- `linked_round_id`: round that produced the recommendation
- `linked_scenario_ids`: scenario ids that contributed to the recommendation
- `linked_promotion_ids`: optional promotion ids that directly drove the recommendation
- `supporting_rule_ids`: optional rule ids that support the recommendation even when no promotion was emitted
- `suggested_action`: operator-friendly follow-up text
- `existing_control_integration_note`: short explanation of how the recommendation fits beside the incumbent fraud stack

Example response:

```json
{
  "data": [
    {
      "id": "rec-round-commerce-20260325-01-provenance",
      "type": "require_provenance_for_delegated_purchase",
      "rationale": "Delegated purchase variants with weakened provenance were under-scored relative to the expected floor.",
      "priority": "high",
      "linked_round_id": "round-commerce-20260325-01",
      "linked_scenario_ids": [
        "commerce-s2-delegated-purchase-weak-provenance",
        "commerce-v1-weakened-provenance"
      ],
      "linked_promotion_ids": [
        "promo-round-commerce-20260325-01-provenance"
      ],
      "supporting_rule_ids": [
        "missing_provenance_sensitive_action"
      ],
      "suggested_action": "Require provenance checks or step-up review before allowing delegated purchase completion.",
      "existing_control_integration_note": "Useful as a sidecar recommendation beside current purchase review thresholds and PSP checks."
    }
  ]
}
```

### TrendSummary object

Returned by:

- `GET /api/v1/benchmark/trends/summary`
- `GET /api/v1/operator/trends/summary`

Schema:

- `rounds_executed`: total historical plus live rounds included in the summary
- `promotions_over_time`: array of `{ round_id, value }` points
- `replay_pass_rate_over_time`: array of `{ round_id, value }` points where `value` is a floating-point pass rate
- `new_blind_spots_discovered`: count of new blind spots found across the summary window
- `regressions_observed`: count of replay or benchmark regressions observed
- `recommendation_counts_by_type`: map keyed by recommendation type
- `top_recurring_evasion_patterns`: ordered list of recurring blind-spot or evasion labels

Example response:

```json
{
  "data": {
    "rounds_executed": 3,
    "promotions_over_time": [
      { "round_id": "round-commerce-20260325-01", "value": 2 },
      { "round_id": "round-commerce-20260325-02", "value": 1 }
    ],
    "replay_pass_rate_over_time": [
      { "round_id": "round-commerce-20260325-01", "value": 1.0 },
      { "round_id": "round-commerce-20260325-02", "value": 0.67 }
    ],
    "new_blind_spots_discovered": 2,
    "regressions_observed": 1,
    "recommendation_counts_by_type": {
      "add_to_replay_stable_set": 2,
      "require_provenance_for_delegated_purchase": 1,
      "tighten_refund_review_rule": 1,
      "monitor_in_shadow_mode": 3
    },
    "top_recurring_evasion_patterns": [
      "weak_provenance",
      "repeat_agent_refund",
      "approval_removed"
    ]
  }
}
```

### SchedulerStatus object

Returned by:

- `GET /api/v1/benchmark/scheduler/status`
- `POST /api/v1/benchmark/scheduler/run` inside `data.status`

Schema:

- `enabled`: whether the scheduler is configured or was explicitly run
- `running`: whether a bounded scheduler loop is currently in progress
- `scenario_family`: scenario family being executed
- `interval`: parsed duration string used between scheduled rounds
- `max_runs`: configured upper bound for the current loop
- `executed_runs`: number of rounds already executed in the current or last scheduler session
- `dry_run`: whether the scheduler is running in recommendation-only demo mode
- `last_round_id`: most recent round produced by the scheduler, when available
- `last_started_at`: timestamp of the most recent scheduler-started round
- `next_run_at`: next scheduled run time when a loop is still active

Example response:

```json
{
  "data": {
    "enabled": true,
    "running": false,
    "scenario_family": "commerce",
    "interval": "24h0m0s",
    "max_runs": 7,
    "executed_runs": 7,
    "dry_run": true,
    "last_round_id": "round-commerce-20260331-07",
    "last_started_at": "2026-03-31T09:00:00Z"
  }
}
```

### SchedulerRunResponse object

Returned by:

- `POST /api/v1/benchmark/scheduler/run`

Schema:

- `rounds`: array of full `BenchmarkRound` objects created during the requested scheduler run
- `status`: `SchedulerStatus` after the run request completes
- `summary`: `TrendSummary` after those rounds are incorporated into history

Example response:

```json
{
  "data": {
    "rounds": [
      {
        "id": "round-commerce-20260325-01",
        "scenario_family": "commerce",
        "round_status": "completed"
      },
      {
        "id": "round-commerce-20260325-02",
        "scenario_family": "commerce",
        "round_status": "completed"
      }
    ],
    "status": {
      "enabled": true,
      "running": false,
      "scenario_family": "commerce",
      "interval": "1s",
      "max_runs": 2,
      "executed_runs": 2,
      "dry_run": true,
      "last_round_id": "round-commerce-20260325-02",
      "last_started_at": "2026-03-25T14:05:02Z"
    },
    "summary": {
      "rounds_executed": 2,
      "promotions_over_time": [
        { "round_id": "round-commerce-20260325-01", "value": 2 },
        { "round_id": "round-commerce-20260325-02", "value": 1 }
      ],
      "replay_pass_rate_over_time": [
        { "round_id": "round-commerce-20260325-01", "value": 1.0 },
        { "round_id": "round-commerce-20260325-02", "value": 1.0 }
      ],
      "new_blind_spots_discovered": 1,
      "regressions_observed": 0,
      "recommendation_counts_by_type": {
        "add_to_replay_stable_set": 1,
        "monitor_in_shadow_mode": 2
      },
      "top_recurring_evasion_patterns": [
        "weak_provenance"
      ]
    }
  }
}
```

Example:

```json
{
  "stable_suite": {
    "id": "stable-suite-placeholder",
    "name": "Stable Placeholder",
    "version": "v1"
  },
  "living_suite": {
    "id": "living-suite-placeholder",
    "name": "Living Placeholder",
    "mutation_policy": "phase-3-none"
  },
  "scenario_pack_id": "starter-pack",
  "scenario_pack_version": "v1",
  "replay_case_refs": ["rc-starter-mandate-review-20260101010101"],
  "notes": "Phase 3 benchmark registration slice"
}
```

Run example:

```json
{
  "scenario_family": "commerce"
}
```

Round responses include:
- stable-set results
- living-set challenger results
- replay regression results
- promotion decisions
- recommendations
- detection delta
- `evaluation_mode=shadow`
- `blocking_mode=recommendation_only`
- report directory and artifact paths

`POST /api/v1/benchmark/scheduler/run` accepts:

```json
{
  "scenario_family": "commerce",
  "interval": "1s",
  "max_runs": 2,
  "dry_run": true
}
```

It returns the executed rounds, current scheduler status, and the updated long-run summary. The scheduler is intentionally small and local-only; it is meant for bounded multi-day demo runs, not distributed orchestration.

Phase 8.1 bootstrap note:

- `GET /api/v1/benchmark/rounds` includes both live in-memory rounds and historical rounds reconstructed from `REPORTS_DIR`
- historical reconstruction requires `reports/<round-id>/round-summary.json`
- `recommendation-report.json` is generated for new rounds and backfilled for legacy rounds when missing
- bootstrap prefers `recommendation-report.json` when it exists; otherwise it reconstructs the recommendation artifact from persisted round data and writes it once
- the Phase 9 validation runner explicitly recognizes legacy rounds that are missing only `recommendation-report.json` but remain reconstructible from persisted round artifacts
- if a round exists both live and persisted, the live round wins and persisted report metadata is preserved

The report API exposes the generated artifacts under `reports/<round-id>/`, including:
- `round-summary.json`
- `round-summary.md`
- `detection-delta.json`
- `promotion-report.json`
- `recommendation-report.json`
- `executive-summary.md`

`recommendation-report.json` is a structured round artifact containing:

- `round_id`
- `evaluation_mode`
- `blocking_mode`
- `existing_control_integration_note`
- `recommended_follow_up`
- `recommendations`

## Operator

- `GET /api/v1/operator/rounds`
- `GET /api/v1/operator/rounds/{id}`
- `GET /api/v1/operator/rounds/{id}/compare?previous=<round-id>`
- `GET /api/v1/operator/promotions`
- `GET /api/v1/operator/promotions/{id}`
- `POST /api/v1/operator/promotions/{id}/review`
- `GET /api/v1/operator/detection/results/{id}`
- `GET /api/v1/operator/recommendations`
- `GET /api/v1/operator/recommendations/{id}`
- `GET /api/v1/operator/trends/summary`
- `GET /api/v1/operator/reports/{round_id}`
- `GET /api/v1/operator/reports/{round_id}/{artifact_name}`

`POST /api/v1/operator/promotions/{id}/review` accepts:

```json
{
  "status": "accepted",
  "note": "Promote this case into the replay baseline."
}
```

Allowed review statuses:
- `accepted`
- `duplicate`
- `needs_follow_up`
- `false_signal`

`GET /api/v1/operator/rounds/{id}/compare` returns deltas for:
- robustness outcome
- promotions count
- replay pass rate
- challenger count
- important findings
- detection delta count

`GET /api/v1/operator/reports/{round_id}/{artifact_name}` returns the artifact descriptor plus report content so the UI can render existing report outputs instead of regenerating them.

Phase 8.1 historical operator behavior:

- `GET /api/v1/operator/rounds` includes bootstrapped historical rounds after restart
- `GET /api/v1/operator/promotions` includes historical promotions reconstructed from `promotion-report.json` when present
- promotion review state is only returned when it has actually been stored; the API does not invent historical review values

Phase 9 operator additions:

- `GET /api/v1/operator/recommendations` exposes structured recommendation outputs generated from promotions, replay regression, and repeat patterns
- `GET /api/v1/operator/trends/summary` exposes week-long-demo style aggregates such as rounds executed, promotions over time, replay pass rate over time, blind spots, regressions, and recurring evasion patterns

## Startup order

1. `clawbot-server`
2. `clawmem`
3. `clawbot-trust-lab`
