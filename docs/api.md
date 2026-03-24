# API

Phase 5 keeps the earlier trust/replay APIs and adds deterministic commerce-world execution plus inspection endpoints.

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

## Startup order

1. `clawbot-server`
2. `clawmem`
3. `clawbot-trust-lab`
