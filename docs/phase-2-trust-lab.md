# Phase 2 Trust Lab

## What Phase 2 includes

- a Go trust-lab service shell
- health, readiness, and version endpoints
- versioned `/api/v1` trust-lab status endpoints
- typed trust-lab domain packages
- a real control-plane client boundary for `clawbot-server`
- a stub memory client boundary for future `clawmem` integration

## What Phase 2 excludes

- the full scenario engine
- the full risk engine
- Red Queen mutation logic
- full `clawmem` storage or retrieval internals
- a production UI

## How this sets up later phases

Phase 2 gives later phases a stable place to add:

- scenario execution flows
- replay evaluation logic
- trust artifact pipelines
- benchmark orchestration
- deeper memory retrieval and storage integration
