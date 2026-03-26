# DRQ Reporting Spec

This document describes the Version 1 DRQ reporting surface in `clawbot-trust-lab`.

It is separate from `scripts/version1_validation_report.py`.

- `version1_validation_report.py` is the readiness and installability check
- DRQ reporting inside trust-lab summarizes benchmark evidence and operator-reviewable findings

## Report types

Version 1 supports three report types:

1. Round report
2. 24-hour dry-run report
3. 1-week management report

All report generation reuses persisted round data. It does not rerun the benchmark.

## Round report

Every benchmark round now writes these files under `reports/<round-id>/`:

- `round-summary.json`
- `round-summary.md`
- `round-report.json`
- `round-report.md`
- `detection-delta.json`
- `promotion-report.json`
- `recommendation-report.json`
- `executive-summary.md`

The round report is the richer DRQ summary for a single benchmark round.

It includes:

- round id
- scenario family
- scenarios executed
- promotions
- recommendations
- regressions derived from detection delta
- production-bridge summary
- Tier A / B / C availability and observed Tier C usage
- notable challenger cases

## 24-hour dry-run report

The dry-run report is generated on demand for the last 24 hours or an explicit time window.

Artifacts are written under `reports/daily/<window>/`:

- `dry-run-report.json`
- `dry-run-report.md`

The dry-run report is intended to answer:

- how many rounds completed in the window
- how many promotions and recommendations were produced
- which recommendation themes kept recurring
- which replay-worthy cases emerged
- whether the run looked operationally stable

Version 1 currently records a generation-time health snapshot for control-plane and memory status.
It does not yet persist a full degraded-period or recovery timeline, so the report says that explicitly instead of inventing incident history.

## 1-week management report

The management report is generated on demand for the last 7 days or an explicit time window.

Artifacts are written under `reports/management/<window>/`:

- `management-report.json`
- `management-report.md`

The management report is executive-friendly and answers:

- what DRQ found that the stable baseline alone did not
- which scenarios repeatedly surfaced issues
- which replay cases look strong enough for longer-lived baseline coverage
- whether the system remained operationally usable in shadow mode
- what next production-side actions should be considered

## Generation commands

Local source usage:

```bash
go run ./cmd/trust-lab report round --round-id <round-id>
go run ./cmd/trust-lab report dry-run --last 24h
go run ./cmd/trust-lab report management --last 168h
```

Docker usage:

```bash
docker compose --env-file .env -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec trust-lab \
  clawbot-trust-lab report round --round-id <round-id>

docker compose --env-file .env -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec trust-lab \
  clawbot-trust-lab report dry-run --last 24h

docker compose --env-file .env -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec trust-lab \
  clawbot-trust-lab report management --last 168h
```

Explicit window usage:

```bash
go run ./cmd/trust-lab report dry-run \
  --from 2026-03-25T00:00:00Z \
  --to 2026-03-26T00:00:00Z

go run ./cmd/trust-lab report management \
  --from 2026-03-19T00:00:00Z \
  --to 2026-03-26T00:00:00Z
```

## Data sources

The reporting subsystem reuses existing benchmark state:

- benchmark rounds
- round summary
- promotions
- recommendations
- detection delta
- production-bridge fields
- Tier usage markers and scenario feature catalog

It does not maintain a second source of truth for benchmark outcomes.

## Current Version 1 limits

Version 1 reporting is intentionally practical.

It does not claim:

- a full BI or analytics warehouse
- chart-heavy dashboards
- a persisted service incident timeline
- a generic replay-promotion API beyond the current Trust Lab workflow

The reports are grounded in the actual persisted round data that Version 1 already owns.
