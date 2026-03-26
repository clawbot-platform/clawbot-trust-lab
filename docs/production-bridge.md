# Production Bridge

`clawbot-trust-lab` is intentionally positioned as a production-bridge system, not a production replacement for the incumbent fraud stack.

## One-sentence positioning

This project is a shadow-mode adversarial regression harness for fraud controls in agentic commerce, not a new primary decisioning engine.

## Operating posture

- `evaluation_mode = shadow`
- `blocking_mode = recommendation_only`

That means the harness:

- executes deterministic human and agentic commerce scenarios
- evaluates them with the current baseline detector
- promotes meaningful failures into replay
- emits operator-facing recommendations
- does not directly block production traffic

## Why this matters

Most fraud teams already have:

- PSP and checkout telemetry
- order and refund history
- merchant/category information
- existing fraud review queues and controls

Phase 9 is built so it can sit beside that stack:

- Tier A uses common commerce and PSP signals
- Tier B uses light aggregations from existing history
- Tier C is optional agentic overlay context

The demo works without Tier C. Tier C only improves differentiation and explanation when it is available.

## What this proves beside an existing fraud stack

A team does not need to replace its current controls to get value from this harness. The production-bridge story is:

1. keep the incumbent stack in charge
2. run trust-lab beside it in shadow mode
3. use replay, promotions, and trend summaries to find blind spots and regressions
4. feed the recommendations back into human review or targeted control tuning

That is a safer and more believable adoption story than claiming a greenfield replacement.

## Suggested adoption path

1. Run the trust-lab benchmark in shadow mode beside the incumbent stack.
2. Compare recommendations against existing fraud review outcomes.
3. Promote meaningful failures into replay.
4. Use round comparisons and long-run summaries to show preserved gains or regressions.
5. Trial targeted rule changes in shadow mode before any production enforcement changes.

## Demo-friendly artifacts

When presenting the project, the most useful artifacts are:

- round summaries that show stable vs living performance
- promoted cases that show concrete blind spots
- replay pass rate that shows whether gains stuck
- recommendation reports that translate findings into next actions
- operator trend summaries that make multi-day runs understandable
