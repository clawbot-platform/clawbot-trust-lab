# Benchmark Model

Phase 7 introduces the Red Queen MVP round model.

## BenchmarkRound

Each round stores:

- `id`
- `scenario_family`
- `detector_version`
- `stable_scenario_refs`
- `challenger_variant_refs`
- `replay_case_refs`
- `started_at`
- `completed_at`
- `round_status`
- `report_dir`
- `scenario_results`
- `promotion_results`
- `delta`
- `stable_set`
- `living_set`
- `summary`
- `reports`

## Stable vs living sets

The round splits execution into:

- stable set
  - known baseline scenarios from `commerce-pack`
  - used to check that expected detector behavior still holds
- living set
  - explicit challenger variants from `challenger-pack`
  - used to probe new weaknesses
- replay regression set
  - previously promoted cases rerun in later rounds
  - used to detect regression against known failures

## ChallengerVariant

Each challenger variant is explicit. Phase 7 does not use a generic mutation engine.

Current variants are:

- weakened provenance
- expired mandate
- approval removed

Each variant stores:

- `id`
- `scenario_id`
- `title`
- `description`
- `change_set`
- `expected_minimum_status`
- `expected_recommendation`

## ScenarioResult

Each execution normalizes into a `ScenarioResult` with:

- set membership
- scenario and variant refs
- entity refs
- detection result ref
- triggered rules
- replay and memory refs
- expected minimum status
- pass/fail outcome
- promotion flag

## PromotionDecision

Promotions are explicit and tied to a reason:

- `detector_miss`
- `suspicious_behavior_scored_too_low`
- `new_trust_gap_pattern`
- `meaningful_regression`
- `novel_evasive_variation`

## RoundSummary

The summary captures:

- stable scenario count
- challenger count
- replay retest count
- promotions count
- replay pass rate
- robustness outcome
- important findings

## Robustness outcomes

- `improved`
- `mixed`
- `regressed`
- `new_blind_spot_discovered`
