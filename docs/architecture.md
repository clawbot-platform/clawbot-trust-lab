# Architecture

`clawbot-trust-lab` is the first vertical repository built on top of the `clawbot-platform` shared foundation.

## Upstream dependencies

- `clawbot-server` provides the shared foundation and control-plane APIs
- `clawmem` provides persistent memory-oriented APIs for trust and replay summaries

## Phase 2 through Phase 5 topology

The runtime stays intentionally small:

- one Go HTTP service
- typed trust-lab domain packages
- a control-plane client abstraction
- a real `clawmem` HTTP client
- a scenario-pack loader
- a deterministic scenario execution layer
- a small in-memory commerce world store
- a file-backed replay archive store
- local trust artifact storage for the trust-lab API surface
- no heavy simulation or execution engine yet

## Multi-repo flow

1. A scenario pack is loaded from disk.
2. `POST /api/v1/scenarios/execute` runs one deterministic commerce scenario.
3. Scenario execution updates a local world state of buyers, merchants, products, orders, payments, refunds, trust decisions, and events.
4. Scenario execution creates a trust artifact and replay case, which in turn write memory records to `clawmem`.
5. Inspection APIs expose orders, events, and trust decisions.

Write failures to `clawmem` are treated as request failures. Status enrichment is best-effort and reports degraded memory context rather than failing the whole status endpoint.

## Phase 5 role in the roadmap

Phase 5 is the event-first world model that later phases will use for:

- a detection baseline
- benchmark rounds against stable and suspicious flows
- replay-backed investigation
- Red Queen style mutation and adaptation

## Boundary decisions

- `clawbot-server` remains the owner of shared control-plane logic
- `clawmem` remains the owner of memory internals
- `clawbot-trust-lab` owns trust-lab domain evolution, trust artifact workflows, replay archival, and orchestration-facing app code
