# Phase 8.1 Historical Round Persistence

## Why this phase exists

Phase 8 made rounds, promotions, and reports reviewable in the operator surface, but the benchmark/operator history still depended on in-memory process state.

That meant:

- report artifacts already existed on disk under `reports/<round-id>/`
- after a trust-lab restart, round list APIs lost that history
- operator promotions disappeared unless the same process was still alive

Phase 8.1 fixes that gap without redesigning the benchmark system.

## Source of truth

Historical reconstruction uses these files under `reports/<round-id>/`:

1. `round-summary.json`
   Required. This is the primary structured source of truth for reconstructing a historical round.
2. `promotion-report.json`
   Optional. Preferred source for historical promotion reconstruction.
3. `detection-delta.json`
   Optional. Used to enrich reconstructed round delta metadata.

Markdown files such as `round-summary.md` and `executive-summary.md` remain first-class report artifacts, but they are not the structured source of truth for round reconstruction.

## Startup bootstrap flow

On startup, trust-lab now:

1. scans `REPORTS_DIR`
2. finds subdirectories that contain `round-summary.json`
3. loads the stored round metadata
4. enriches it from `promotion-report.json` and `detection-delta.json` when present
5. rebuilds report artifact descriptors from the directory listing
6. reconstructs minimal historical detection results from stored scenario results
7. loads the reconstructed rounds into the benchmark store as historical state

Malformed report directories are skipped with explicit logging. One bad report directory does not fail the whole service.

## Merge behavior

The benchmark store now merges two sources:

- live in-memory rounds created in the current process
- historical rounds reconstructed from `reports/`

Merge rules:

- list APIs return a unified view of both
- rounds are sorted most recent first
- if the same round id exists in both places, live state wins
- persisted report metadata is preserved for that round

## Historical promotions

Historical promotions are reconstructed from the round artifacts, primarily `promotion-report.json`.

That means after restart:

- `/api/v1/operator/promotions` still shows historical promotions
- `/api/v1/operator/promotions/{id}` still resolves the linked round and scenario result
- review state is only shown if it was explicitly persisted elsewhere

Phase 8.1 does not invent historical review statuses that were never stored.

## Limitations

- reconstructed historical detection results are intentionally minimal and derive from stored scenario results
- this phase does not add a database or a second persistence system
- operator review notes still follow their existing store behavior

The goal is durable benchmark/operator history, not a broader state-management redesign.
