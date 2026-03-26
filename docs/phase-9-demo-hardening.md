# Phase 9 Demo Hardening

Phase 9 turns the earlier Red Queen MVP into a stronger portfolio and demo asset.

The point of this phase is not “more simulation.” The point is to make the trust lab look like something a fraud team could actually trial in shadow mode beside an incumbent control stack.

## What a reviewer should understand quickly

After Phase 9, this repo can show five things clearly:

1. the scenario catalog now looks like real commerce and refund behavior instead of toy cases
2. the detector distinguishes human and delegated or agentic behavior using mostly ordinary fraud-stack data
3. the benchmark loop keeps learning from replay rather than forgetting prior failures
4. the output is framed as recommendations, not unsafe production blocking
5. the operator surface can explain what changed over time

## Demo walkthrough

Inspect the stable and challenger packs first:

```bash
curl http://127.0.0.1:8090/api/v1/scenarios/packs/commerce-pack
curl http://127.0.0.1:8090/api/v1/scenarios/packs/challenger-pack
```

Run one round:

```bash
curl -X POST http://127.0.0.1:8090/api/v1/benchmark/rounds/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce"}'
```

Then show the operator-facing outputs:

```bash
curl http://127.0.0.1:8090/api/v1/operator/promotions
curl http://127.0.0.1:8090/api/v1/operator/recommendations
curl http://127.0.0.1:8090/api/v1/operator/trends/summary
```

If you only have a few minutes, that sequence tells the whole story:

- what gets tested
- what slipped through or regressed
- what got promoted into replay
- what the system recommends next
- how the trend line is moving over time

## What changed

- the scenario catalog now covers human baseline, benign agentic, suspicious, and challenger flows from `docs/phase-9-scenario-catalog.md`
- feature extraction is explicitly tiered as Tier A, Tier B, and Tier C
- the detector still works with Tier A plus Tier B alone
- benchmark rounds now emit structured recommendations
- round summaries now carry shadow-mode / recommendation-only framing
- a bounded scheduler loop can execute repeat rounds for multi-day demo runs
- long-run summaries expose promotions over time, replay pass rate over time, blind spots, regressions, and recurring patterns

## Why the catalog matters

Phase 9 moves the benchmark from “clever adversarial demo” toward “credible fraud-control validation harness.”

The catalog now includes:

- human baseline behavior that should not look suspicious by default
- benign delegated or agentic behavior with valid controls
- suspicious refund and delegated-purchase patterns that should remain explainably catchable
- challenger variants that probe whether the baseline detector preserves its gains

## Stable and living sets

The stable set now emphasizes:

- direct human and human refund baselines
- benign delegated and agent-assisted flows
- suspicious refund flows that should remain explainably catchable

The living set now emphasizes:

- weak provenance
- expired mandate coverage
- removed approval evidence
- actor switch from human to agent
- repeat attempt escalation
- merchant/category scope drift
- high-value delegated purchase

## Recommendation layer

Phase 9 recommendations are intentionally operational and sidecar-friendly:

- `add_to_replay_stable_set`
- `tighten_refund_review_rule`
- `require_step_up_for_delegated_refunds`
- `require_provenance_for_delegated_purchase`
- `investigate_repeat_refund_pattern`
- `monitor_in_shadow_mode`

They are generated from round promotions, replay outcomes, and rule-hit patterns rather than from opaque analytics.

That makes them easy to explain in a demo:

- this case should be added to replay
- this refund rule probably needs tightening
- this delegated-refund path should require step-up
- this purchase path should require provenance
- this pattern is worth monitoring in shadow mode before changing production controls

## Tier A / B / C in practice

- Tier A uses commerce, refund, order, merchant, and PSP-adjacent signals that most existing stacks already have
- Tier B uses light historical aggregation, such as repeat attempts or recent history
- Tier C uses optional agentic overlays like mandate or provenance details

The important point for the Phase 9 story is that the harness still works with Tier A plus Tier B alone. Tier C improves differentiation, but it is not required to make the demo meaningful.

## 7-day demo mode

The lightweight scheduler uses:

- `BENCHMARK_SCHEDULER_ENABLED`
- `BENCHMARK_SCHEDULER_SCENARIO_FAMILY`
- `BENCHMARK_SCHEDULER_INTERVAL`
- `BENCHMARK_SCHEDULER_MAX_RUNS`
- `BENCHMARK_SCHEDULER_DRY_RUN`

For a week-long demo, set:

- interval to `24h`
- max runs to `7`
- dry run to `true`

The rounds still write real reports, recommendations, promotions, and long-run summary data. "Dry run" here means the harness remains recommendation-only; it does not block live commerce decisions.

## What to call out in a week-long run

- promotions should accumulate only when the harness finds something worth preserving
- replay pass rate should show whether prior gains held
- recurring evasion patterns should become visible in trend summaries
- recommendations should stay operational and bounded, not turn into vague analytics output
