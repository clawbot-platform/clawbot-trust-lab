# Architecture

`clawbot-trust-lab` is the first vertical repository built on top of the `clawbot-platform` shared foundation.

## Upstream dependencies

- `clawbot-server` provides the shared foundation and control-plane APIs
- `clawmem` provides persistent memory-oriented APIs for trust and replay summaries

## Phase 2 through Phase 4.1 topology

The runtime stays intentionally small:

- one Go HTTP service
- typed trust-lab domain packages
- a control-plane client abstraction
- a real `clawmem` HTTP client
- a scenario-pack loader
- a file-backed replay archive store
- local trust artifact storage for the trust-lab API surface
- no heavy simulation or execution engine yet

## Multi-repo flow

1. A scenario pack is loaded from disk.
2. `POST /api/v1/trust/artifacts` creates a trust artifact and writes a corresponding trust memory record to `clawmem`.
3. `POST /api/v1/replay/cases` creates a replay case, writes a replay memory record to `clawmem`, and archives the replay case locally.
4. Status endpoints can optionally enrich responses with `clawmem` context by `scenario_id`.

Write failures to `clawmem` are treated as request failures. Status enrichment is best-effort and reports degraded memory context rather than failing the whole status endpoint.

## Boundary decisions

- `clawbot-server` remains the owner of shared control-plane logic
- `clawmem` remains the owner of memory internals
- `clawbot-trust-lab` owns trust-lab domain evolution, trust artifact workflows, replay archival, and orchestration-facing app code
