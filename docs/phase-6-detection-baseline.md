# Phase 6 Detection Baseline

Phase 6 adds the first benchmarkable detector for the trust lab.

This phase does not attempt to be clever. It establishes a deterministic baseline that later phases can measure against.

## What became real in Phase 6

- a typed detection domain
- deterministic feature extraction from world state and event streams
- explicit baseline rules
- stored detection results
- API endpoints to evaluate and inspect results
- lightweight replay and `clawmem` context enrichment

## Why this phase exists

Phase 5 created a small but coherent commerce world. That world now produces:

- orders
- refunds
- payments
- trust decisions
- append-only events
- replay artifacts
- `clawmem` memory records

Phase 6 turns those outputs into the first explainable detector. The goal is not accuracy theater. The goal is to create a baseline that is:

- deterministic
- reproducible
- easy to debug
- easy to benchmark
- useful as a reference point for later Red Queen work

## Baseline outcomes

### Clean agent-assisted purchase

The clean purchase scenario generally evaluates to `clean` because:

- a delegated actor is present
- an active mandate exists
- provenance exists
- no refund path is involved
- no prior step-up requirement exists

This scenario demonstrates the detector can recognize a trust-complete delegated flow instead of flagging every agent action.

### Suspicious refund attempt

The suspicious refund scenario evaluates to `step_up_required` because it combines:

- an agent-driven refund action
- expired authority
- weak refund authorization
- missing approval evidence
- a prior trust decision that already required step-up

This scenario demonstrates that the baseline is event-first and trust-surface aware, not just order-state aware.

## Minimal replay and memory enrichment

Phase 6 uses replay and `clawmem` context in a bounded way:

- replay history count contributes context
- memory record presence contributes context
- repeat suspicious context can trigger an additional rule hit

There is no semantic retrieval, vector search, or hidden ranking logic in this phase.

## What remains for Phase 7

Phase 7 can now focus on Red Queen style mutation and adaptation because the platform already has:

- deterministic scenarios
- deterministic detector outputs
- replayable suspicious cases
- memory-backed context
- a stable API surface for evaluation and inspection

That means Phase 7 can measure how mutations change rule hits and outcomes, rather than inventing its evaluation surface from scratch.
