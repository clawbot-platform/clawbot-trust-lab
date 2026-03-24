# Phase 4.1 clawmem Integration

Phase 4.1 replaces the stubbed memory path from Phase 3 with a real HTTP integration to `clawmem`.

## What changed

- `clawbot-trust-lab` now has a real `clawmem` HTTP client
- trust artifact creation writes a trust memory record to `clawmem`
- replay case creation writes a replay memory record to `clawmem`
- readiness now checks both the control plane and `clawmem`
- trust and replay status endpoints can include `clawmem`-backed context

## Error handling choice

Create flows treat `clawmem` as a required dependency.

- if `clawmem` write succeeds, the trust-lab write continues
- if `clawmem` write fails, the request returns `502 Bad Gateway`
- trust-lab does not silently swallow the failure

Status enrichment is softer:

- if `clawmem` context loads, the endpoint includes memory-backed details
- if `clawmem` context fails, the endpoint still returns `200` with a degraded status field

## Local validation flow

1. Start `clawbot-server`
2. Start `clawmem`
3. Start `clawbot-trust-lab`
4. Create a trust artifact
5. Create a replay case
6. Query `clawmem` directly to verify the memory records were written
