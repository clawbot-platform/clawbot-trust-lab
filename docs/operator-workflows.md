# Operator Workflows

## Workflow 1: Review benchmark rounds

Use the rounds surface to:

- identify the latest rounds
- inspect promotions count
- inspect replay pass rate
- inspect robustness outcome
- inspect whether a round relied only on Tier A and Tier B or also used Tier C overlays
- inspect production-bridge recommendations without leaving the round detail

The round detail page is the pivot point into promotions and reports.

Historical rounds now survive trust-lab restart because operator views are bootstrapped from `reports/<round-id>/round-summary.json` and related report artifacts at startup.

## Workflow 2: Review promoted challenger cases

Use the promotions surface to:

- find newly promoted cases
- keep seeing historical promotions after restart
- filter by review status
- open a promotion detail
- inspect the linked detection result
- record a lightweight operator action
- understand whether the detection result relied on optional Tier C context

## Workflow 3: Compare rounds

From round detail, compare a round against a previous round to inspect:

- robustness outcome changes
- promotion count changes
- replay pass-rate changes
- challenger-count changes
- important finding changes

This is intentionally small and operator-readable.

## Workflow 4: Review recommendations and long-run trends

Use the rounds surface and round detail to:

- inspect recommendation counts and recommendation types
- see whether promotions are feeding replay growth
- spot recurring evasion patterns over multiple rounds
- position the harness as a shadow-mode sidecar beside the incumbent fraud stack

## Workflow 5: Browse reports in-app

Use the reports page to:

- list report artifacts for a round
- open Markdown summaries
- inspect JSON artifacts
- move back to round detail without leaving the operator app

This works for both live rounds from the current process and historical rounds reconstructed from disk.

## Review status meaning

- `accepted`: the promotion looks valid and should inform future replay review
- `duplicate`: the promotion does not add net-new value
- `needs_follow_up`: the promotion is interesting but requires more investigation
- `false_signal`: the promotion does not represent a meaningful detector gap
