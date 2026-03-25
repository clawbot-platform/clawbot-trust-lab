# Reporting Spec

Phase 7 reporting exists to make each Red Queen round measurable and explainable.

## Required outputs

Each round writes these files under `reports/<round-id>/`:

1. `round-summary.json`
2. `round-summary.md`
3. `detection-delta.json`
4. `promotion-report.json`
5. `executive-summary.md`

## Required content

At minimum the reports must capture:

- round id
- scenario family
- stable scenario count
- challenger count
- replay retest count
- promotions count
- replay pass rate
- robustness outcome
- important findings
- promoted case rationale

## Detection delta

The delta report compares the current round with the previous round when one exists.

It should surface:

- status changes
- score or severity changes
- newly triggered rules
- cleared rules
- recommendation changes

## Promotion report

Every promotion must include:

- scenario id
- challenger variant id
- promotion reason
- rationale
- linked detection result
- linked replay case

Promotion cannot be implicit.

## Executive summary

The executive summary should stay short and answer:

- did the detector improve, regress, or expose a new blind spot
- how many cases were promoted
- what the most important finding was
- what the next action should be

## Stable vs living reporting rule

Round reporting must keep stable-set and living-set performance visibly separate.

This is the minimum contract that later operator surfaces can build on.
