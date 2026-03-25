# API

Phase 8 keeps the earlier trust/replay APIs, the deterministic commerce-world execution layer, the explainable detection baseline, the Red Queen MVP benchmark loop, and adds operator-facing review APIs.

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

Execute example:

```json
{
  "scenario_id": "commerce-clean-agent-assisted-purchase"
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
  "scenario_id": "commerce-suspicious-refund-attempt"
}
```

You can also evaluate by `order_id`:

```json
{
  "order_id": "order-clean-agent-assisted-purchase"
}
```

Detection responses include:
- scenario and entity refs
- status, score, grade, and recommendation
- triggered baseline rules
- reason codes
- linked trust decision and replay refs
- a serialized detection context in `metadata.context`

`GET /api/v1/detection/rules` lists the active baseline rules:
- `missing_mandate_delegated_action`
- `missing_provenance_sensitive_action`
- `refund_weak_authorization`
- `agent_refund_without_approval`
- `prior_step_up_decision`
- `repeat_suspicious_context`

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
- `GET /api/v1/benchmark/rounds/status`
- `GET /api/v1/benchmark/status`

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
- detection delta
- report directory and artifact paths

Phase 8.1 bootstrap note:

- `GET /api/v1/benchmark/rounds` includes both live in-memory rounds and historical rounds reconstructed from `REPORTS_DIR`
- historical reconstruction requires `reports/<round-id>/round-summary.json`
- if a round exists both live and persisted, the live round wins and persisted report metadata is preserved

The report API exposes the generated artifacts under `reports/<round-id>/`, including:
- `round-summary.json`
- `round-summary.md`
- `detection-delta.json`
- `promotion-report.json`
- `executive-summary.md`

## Operator

- `GET /api/v1/operator/rounds`
- `GET /api/v1/operator/rounds/{id}`
- `GET /api/v1/operator/rounds/{id}/compare?previous=<round-id>`
- `GET /api/v1/operator/promotions`
- `GET /api/v1/operator/promotions/{id}`
- `POST /api/v1/operator/promotions/{id}/review`
- `GET /api/v1/operator/detection/results/{id}`
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

## Startup order

1. `clawbot-server`
2. `clawmem`
3. `clawbot-trust-lab`
