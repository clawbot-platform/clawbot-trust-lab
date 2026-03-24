# Scenario Pack Format

Scenario packs are versioned JSON files stored under `configs/scenario-packs/`.

## Top-level fields

- `id`
- `name`
- `version`
- `description`
- `scenarios`

## Scenario fields

- `id`
- `name`
- `version`
- `scenario_type`
- `description`
- `actors`
- `trust_signals`
- `expected_outcomes`
- `tags`

## Notes

- The current loader validates required ids, names, versions, and scenario types.
- The format is intentionally simple so it is easy to explain and evolve in later phases.
- The starter pack in `configs/scenario-packs/starter-pack.json` is the reference example for Phase 3.
