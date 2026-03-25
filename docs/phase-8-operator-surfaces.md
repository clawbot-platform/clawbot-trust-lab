# Phase 8 Operator Surfaces

Phase 8 makes the outputs of Phases 5 through 7 usable by a human operator.

This phase does not add new detector logic or new challenger generation. It adds visibility and lightweight review workflows.

## What became real in Phase 8

- operator review API endpoints
- round comparison support
- promotion review actions
- in-app report browsing
- a small React operator UI in `web/`

## Main operator journey

1. Open the rounds page.
2. Inspect a round summary and its robustness outcome.
3. Open promoted challenger cases from that round.
4. Inspect the linked detection result and scenario context.
5. Record a lightweight review action.
6. Compare the round with a previous round.
7. Open the generated reports without leaving the app.

## Scope boundary

Phase 8 is not a generic admin panel.

It focuses on:

- rounds
- promotions
- detection inspection
- comparison
- report review

It does not attempt to become:

- a case-management suite
- an analytics warehouse
- an auth-heavy enterprise dashboard

## Review actions

Operators can mark a promotion as:

- accepted
- duplicate
- needs_follow_up
- false_signal

Each review can include a short note. The storage is intentionally lightweight.

## Why this prepares later phases

Later operator surfaces can build on Phase 8 because the repo now has:

- stable URLs for review workflows
- report browsing over existing artifacts
- lightweight operator state for promotions
- round comparison support for human analysis
