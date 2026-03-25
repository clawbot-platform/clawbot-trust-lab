# Phase 7 Red Queen MVP

Phase 7 adds the first adversarial benchmark loop to `clawbot-trust-lab`.

This is intentionally a small Red Queen MVP, not a generic simulator or tournament framework.

## What became real in Phase 7

- benchmark round models
- explicit challenger variants
- a deterministic round runner
- replay promotion policy
- replay regression checks
- round-level report generation
- APIs to run and inspect rounds

## Stable vs living sets

Phase 7 separates:

- stable set
  - the known baseline commerce scenarios
- living set
  - explicit challenger variants designed to probe detector weaknesses
- replay regression set
  - previously promoted cases rerun in later rounds

This prevents false confidence from improving only on old cases while failing on new challenger behavior.

## Challenger variants

The MVP keeps challengers explicit and inspectable:

- weakened provenance on delegated purchase
- expired mandate on delegated purchase
- approval removed from agent refund

There is no mutation framework in this phase. The value is in the loop, not in synthetic variation breadth.

## Promotion policy

A challenger case is promoted when the current detector meaningfully underperforms, including:

- suspicious behavior evaluated as clean
- suspicious behavior scored below its expected floor
- a meaningful replay regression

Promotions produce explicit records with:

- promotion reason
- rationale
- linked scenario result
- linked detection result
- linked replay case ref

## Replay regression

Each new round retests previously promoted cases. This provides a small deterministic regression signal without requiring new infrastructure.

The round summary computes replay pass rate and uses it to help derive the robustness outcome.

## Reports

Every round writes artifacts under `reports/<round-id>/`:

- `round-summary.json`
- `round-summary.md`
- `detection-delta.json`
- `promotion-report.json`
- `executive-summary.md`

This makes the round useful to both humans and later automation.

## Why this prepares later phases

Phase 7 gives later operator surfaces and richer DRQ behavior a real foundation:

- stable and living evaluation sets
- repeatable challenger execution
- explicit promotion history
- replay regression checks
- report artifacts that can feed future dashboards or trend analysis
