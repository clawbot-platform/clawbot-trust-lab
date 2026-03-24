# API

Phase 3 exposes real trust-lab behavior in addition to the Phase 2 shell endpoints.

## System

- `GET /healthz`
- `GET /readyz`
- `GET /version`

## Scenarios

- `GET /api/v1/scenarios/types`
- `GET /api/v1/scenarios/packs`
- `GET /api/v1/scenarios/packs/{id}`

## Trust

- `POST /api/v1/trust/artifacts`
- `GET /api/v1/trust/artifacts`
- `GET /api/v1/trust/status`

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
