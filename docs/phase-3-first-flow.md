# Phase 3 First Flow

Phase 3 implements the first real trust-lab vertical slice:

1. A scenario pack is loaded from `configs/scenario-packs/`
2. A scenario is selected by id
3. A trust artifact is created from that scenario
4. The memory client is invoked with a structured trust-artifact summary
5. A replay case is created and archived through the local replay store
6. A benchmark registration request is sent through the control-plane client abstraction

## What is real now

- scenario pack loading from disk
- typed trust artifact creation
- replay archive writes to a local file-backed store
- benchmark registration service wiring
- one real memory client usage path

## What is still not implemented

- actual scenario execution
- benchmark execution orchestration
- full memory retrieval/storage backend
- risk or adversary engines
