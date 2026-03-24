# Architecture

`clawbot-trust-lab` is the first vertical repository built on top of the `clawbot-platform` shared foundation.

## Upstream dependencies

- `clawbot-server` provides the shared foundation and control-plane APIs
- future `clawmem` integration will provide memory-oriented capabilities

## Phase 2 topology

Phase 2 keeps the runtime intentionally small:

- one Go HTTP service
- typed trust-lab domain packages
- a control-plane client abstraction
- a stub memory client abstraction
- no heavy simulation or execution engine yet

## Boundary decisions

- `clawbot-server` remains the owner of shared control-plane logic
- `clawmem` remains the owner of memory internals
- `clawbot-trust-lab` owns trust-lab domain evolution and orchestration-facing app code
