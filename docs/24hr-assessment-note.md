# DRQ V1 · 24-Hour Dry Run Assessment

## Summary

The 24-hour DRQ Version 1 dry run was an **operational success** and a **detector-quality partial success**.

The platform stayed up, completed scheduled benchmark rounds, generated round artifacts, persisted memory-backed trust and replay state, and produced repeatable promotions and recommendations.

At the same time, the run confirmed that the detector is **not yet ready to move beyond shadow-mode recommendation-only use**. The same replay-regression blind spots recurred across multiple rounds, which means the system is stable enough to evaluate, but detector tuning is still required before stronger production-readiness claims can be made.

## What worked

- `clawbot-server`, `clawmem`, and `clawbot-trust-lab` stayed healthy enough to complete the run.
- The benchmark scheduler executed **4 rounds** over the dry-run window.
- Each round produced the expected artifact set:
  - `round-summary`
  - `round-report`
  - `promotion-report`
  - `recommendation-report`
  - `executive-summary`
  - `detection-delta`
- The run produced:
  - **12 promotions**
  - **24 recommendations**
- Memory-backed trust/replay integration remained active through `clawmem`.
- The stable baseline scenarios remained consistently healthy.

## Key findings

### 1. Runtime/platform stability was proven

The dry run demonstrated that the current DRQ V1 stack can run continuously in the homelab for at least 24 hours with:

- repeated round execution
- report generation
- promotion and recommendation generation
- persistent memory integration
- no obvious platform instability in the observed runtime logs

### 2. Baseline/stable scenarios are behaving correctly

The latest round showed the stable set passing **7/7**, which is the strongest signal that the core benchmark harness and baseline scenario behavior are stable.

### 3. Detector blind spots are recurring, not random

The most important result of the dry run is that the same three cases kept recurring as promoted replay-regression problems:

- `commerce-s3-approval-removed-after-authorization`
- `commerce-v2-expired-inactive-mandate`
- `commerce-v3-approval-removed`

This means the current detector posture is not simply “noisy”; it has a repeatable weakness pattern that the benchmark harness is correctly surfacing.

### 4. Replay preservation is working, but replay improvement is not yet happening

The replay loop is functioning operationally, because promoted cases are being carried forward and retested.

However, the replay regression pass rate dropped to `0.00` in later rounds, which means the detector is not yet improving against the replay set after the initial discovery round.

### 5. Shadow-mode posture remains the correct operating mode

The recommendation outputs remained sensible and consistent across rounds:

- keep the harness in shadow mode
- add promoted cases to replay
- tighten refund review
- require step-up for delegated refunds
- require provenance for delegated purchases
- investigate repeat refund patterns

That is exactly the posture DRQ V1 should keep at this stage: recommendation-only, sidecar to incumbent controls, and focused on learning rather than blocking.

## Operational result

**Pass**

The stack completed the dry run and produced useful evidence.

## Detector-improvement result

**Not yet pass**

The run surfaced persistent replay-regression blind spots that need targeted tuning before stronger claims can be made.

## Recommendation

Use this 24-hour dry run as the baseline proof that:

1. the DRQ V1 platform is operationally viable in the homelab
2. the benchmark harness is producing useful repeatable findings
3. the next milestone should focus on detector improvement against the recurring replay-regression cases

## Next decision before a 1-week run

A 1-week run should be framed explicitly as one of the following:

### Option A — Stability over time
Use the same detector/ruleset with no tuning changes, and measure:
- uptime
- round completion
- repeated artifact generation
- recommendation consistency
- recurring blind spots over time

### Option B — Detector improvement after tuning
Apply targeted tuning for the three recurring replay-regression cases, then run again and measure:
- replay pass-rate improvement
- reduction in repeated promotions
- change in recommendation mix
- whether blind spots move from repeated regressions to stable replay coverage

The strongest overall story is to do **both**, but keep them analytically separate.
