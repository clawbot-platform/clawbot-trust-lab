# Architecture

`clawbot-trust-lab` is the first vertical repository built on top of the `clawbot-platform` shared foundation.

## Upstream dependencies

- `clawbot-server` provides the shared foundation and control-plane APIs
- `clawmem` provides persistent memory-oriented APIs for trust and replay summaries

## Phase 2 through Phase 9 topology

The runtime stays intentionally small:

- one Go HTTP service
- typed trust-lab domain packages
- a control-plane client abstraction
- a real `clawmem` HTTP client
- a scenario-pack loader
- a deterministic scenario execution layer
- a deterministic benchmark-round runner
- a bounded local benchmark scheduler for multi-day demo runs
- a lightweight reporting layer
- a recommendation layer that stays in shadow-mode / recommendation-only posture
- a small in-memory commerce world store
- a file-backed replay archive store
- local trust artifact storage for the trust-lab API surface
- no heavy simulation or execution engine yet

## Multi-repo flow

1. A scenario pack is loaded from disk.
2. `POST /api/v1/scenarios/execute` runs one deterministic commerce scenario.
3. Scenario execution updates a local world state of buyers, merchants, products, orders, payments, refunds, trust decisions, and events.
4. Scenario execution creates a trust artifact and replay case, which in turn write memory records to `clawmem`.
5. Phase 6 detection APIs evaluate the resulting world state.
6. Phase 7 benchmark APIs run stable scenarios, challenger variants, and replay regression cases, then emit report artifacts under `reports/<round-id>/`.
7. Phase 9 recommendations and long-run summaries aggregate round history into operator-facing production-bridge outputs.

Write failures to `clawmem` are treated as request failures. Status enrichment is best-effort and reports degraded memory context rather than failing the whole status endpoint.

## Phase 5 through Phase 7 role in the roadmap

These phases together create the baseline loop that later phases will use for:

- a detection baseline
- benchmark rounds against stable and suspicious flows
- replay-backed investigation
- Red Queen style mutation and adaptation

Phase 9 makes the loop easier to adopt in practice by framing it as a sidecar:

- `evaluation_mode = shadow`
- `blocking_mode = recommendation_only`
- ordinary Tier A and Tier B commerce data are sufficient
- Tier C agentic overlays improve differentiation but remain optional

## Boundary decisions

- `clawbot-server` remains the owner of shared control-plane logic
- `clawmem` remains the owner of memory internals
- `clawbot-trust-lab` owns trust-lab domain evolution, trust artifact workflows, replay archival, and orchestration-facing app code
