# clawbot-trust-lab

[![ci](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/ci.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/ci.yml)
[![quality](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/quality.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/quality.yml)
[![security](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/security.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/security.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=clawbot-platform_clawbot-trust-lab&metric=alert_status&token=abb591daa9f6778dcdc919142fe123aa30947073)](https://sonarcloud.io/summary/new_code?id=clawbot-platform_clawbot-trust-lab)

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)
![Vite](https://img.shields.io/badge/Vite-8-646CFF?logo=vite)
![Vitest](https://img.shields.io/badge/Vitest-Tested-6E9F18?logo=vitest)
![Playwright](https://img.shields.io/badge/Playwright-E2E-2EAD33?logo=playwright)
![Shadow Mode](https://img.shields.io/badge/Mode-Shadow%20%2F%20Recommendation--Only-0F766E)

`clawbot-trust-lab` is the replay-driven adversarial regression harness for fraud controls in agentic commerce inside `clawbot-platform`.

Phase 9 hardens the Red Queen demo into a production-bridge story:

- deterministic human and agentic commerce scenarios
- explainable baseline detection over ordinary commerce and fraud-stack data
- replay-preserving benchmark rounds with challenger variants
- operator-facing promotions, reports, recommendations, and long-run trends
- shadow-mode outputs that can sit beside an existing fraud stack without blocking production traffic

## Why this is different from other Clawbot projects

This repo is not an assistant shell and not a generic agent demo. It is the domain/control layer for validating fraud controls against evolving human, delegated, and agent-driven commerce behaviors.

The core value proposition is:

- continuously test incumbent fraud controls against new evasions
- preserve prior gains through replay regression
- surface operator-reviewable blind spots and recommendations
- run in recommendation-only shadow mode beside existing fraud systems

## Why ZeroClaw

ZeroClaw is used as the runtime substrate across the wider platform because it gives the project a boring, inspectable runtime foundation instead of pushing runtime logic into this repo. `clawbot-trust-lab` stays focused on commerce scenarios, trust surfaces, detection, replay, and benchmark orchestration while ZeroClaw remains the execution substrate underneath the platform.

## What belongs here

- trust-lab specific service code
- scenario, trust, replay, benchmark, and agent domain models
- commerce-world and actor/delegation domain models
- integrations with `clawbot-server`
- the real HTTP integration client for `clawmem`
- trust-lab docs, tests, and CI/security automation

## What does not belong here

- shared platform control-plane logic already owned by `clawbot-server`
- full `clawmem` internals
- the full scenario engine
- the full risk engine
- Red Queen mutation logic
- reimplementations of ZeroClaw runtime behavior

## Quick start

```bash
cp .env.example .env
go run ./cmd/trust-lab
```

Expected local startup order:

1. Start `clawbot-server` so the control-plane health check passes.
2. Start `clawmem` so trust-lab can persist memory records.
3. Start `clawbot-trust-lab`.

## Demo in 5 minutes

If you are showing this repo to a hiring manager, staff engineer, or fraud leader, the clearest story is:

1. inspect a realistic commerce scenario catalog
2. run one benchmark round in shadow mode
3. inspect the promoted blind spots and recommendations
4. show that the harness can keep running on a schedule without becoming a production blocker

Inspect the catalog that drives the demo:

```bash
curl http://127.0.0.1:8090/api/v1/scenarios/packs/commerce-pack
curl http://127.0.0.1:8090/api/v1/scenarios/packs/challenger-pack
```

Run one Phase 9 benchmark round and keep the returned round id:

```bash
ROUND_ID=$(curl -s -X POST http://127.0.0.1:8090/api/v1/benchmark/rounds/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce"}' | jq -r '.data.id')
echo "$ROUND_ID"
```

That response shows the core claim of the project in one payload:

- stable-set results for known-good and known-bad cases
- living-set challenger results for evasive variants
- promotion decisions for newly meaningful failures
- replay regression results for previously preserved gains
- recommendations phrased for shadow-mode adoption

Inspect the latest round like an operator would:

```bash
curl http://127.0.0.1:8090/api/v1/operator/rounds
curl http://127.0.0.1:8090/api/v1/operator/rounds/$ROUND_ID
curl http://127.0.0.1:8090/api/v1/operator/promotions
curl http://127.0.0.1:8090/api/v1/operator/recommendations
curl http://127.0.0.1:8090/api/v1/operator/reports/$ROUND_ID
curl http://127.0.0.1:8090/api/v1/operator/trends/summary
```

Run a short scheduled demo loop:

```bash
curl -X POST http://127.0.0.1:8090/api/v1/benchmark/scheduler/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce","interval":"1s","max_runs":2,"dry_run":true}'
```

Then inspect the scheduler and accumulated long-run summary:

```bash
curl http://127.0.0.1:8090/api/v1/benchmark/scheduler/status
curl http://127.0.0.1:8090/api/v1/benchmark/trends/summary
```

## What this demo proves

- it can evaluate ordinary human commerce and delegated or agentic commerce in the same benchmark loop
- it can discover blind spots without requiring exotic lab-only telemetry
- it can preserve gains through replay instead of rediscovering the same failure every round
- it can produce reviewable recommendations in `shadow` / `recommendation_only` mode beside an incumbent fraud stack
- it can keep producing usable operator artifacts over time, not just one-off benchmark output

## Tier A / B / C

- Tier A: common PSP, checkout, and commerce data already present in most fraud stacks
- Tier B: historical and aggregation features derivable with light engineering
- Tier C: optional agentic overlay fields such as mandate or provenance details

The Phase 9 detector works with Tier A plus Tier B alone. Tier C improves differentiation and explainability, but the demo does not require exotic lab-only telemetry.

## What the week-long demo proves

- the benchmark loop can keep executing over time
- stable and living sets remain explicit and deterministic
- new blind spots and regressions accumulate into replay, reports, and recommendations
- the system can operate in `shadow` / `recommendation_only` mode beside an incumbent fraud stack

Historical rounds are reloaded from `REPORTS_DIR` on startup. The structured source of truth is:

- `reports/<round-id>/round-summary.json`
- `reports/<round-id>/promotion-report.json`
- `reports/<round-id>/detection-delta.json`
- `reports/<round-id>/recommendation-report.json`

`recommendation-report.json` is generated for new rounds. For legacy rounds that predate it, trust-lab reconstructs and backfills the artifact during historical bootstrap from the persisted round data. The backfill is deterministic and idempotent. The Phase 9 validation runner also treats a legacy round that is missing only `recommendation-report.json` as an explicit reconstructible legacy case instead of a silent artifact failure.

If a round exists both in the live in-memory store and on disk, trust-lab prefers the live round data and keeps the persisted report artifact metadata.

Run the Phase 8 operator UI:

```bash
cd web
npm install
npm run dev
```

Run tests:

```bash
go test ./...
```

```bash
cd web
npm run lint
npm run test
npm run test:coverage
npm run build
```

Run the optional Phase 8 operator E2E smoke tests:

```bash
cd web
npm run test:e2e
```

## Quality and coverage

Version 1 is wired to report both backend and web quality:

- Go coverage is generated from `go test -covermode=atomic -coverprofile=coverage.out ./...`
- operator UI coverage is generated as LCOV via `web/coverage/lcov.info`
- SonarCloud ingests both coverage reports and enforces a quality gate in CI
- Playwright remains part of the CI path so the one-week benchmark run keeps a minimal operator smoke test around the core review workflow

Local quality commands:

```bash
go test ./...
go vet ./...
golangci-lint run ./...
make coverage
make security
```

```bash
cd web
npm run lint
npm run test
npm run test:coverage
npm run build
npm run test:e2e
```

SonarCloud is configured for:

- organization: `clawbot-platform`
- project key: `clawbot-platform_clawbot-trust-lab`
- project page: [SonarCloud overview](https://sonarcloud.io/project/overview?id=clawbot-platform_clawbot-trust-lab)

## Repo layout

- `cmd/trust-lab/` contains the service entrypoint
- `internal/app/` wires config, bootstrap, router, and graceful shutdown
- `internal/domain/` contains trust-lab domain types and Phase 3 plus 4.1 services
- `internal/services/` contains the commerce, event, trust-decision, scenario execution, detection, benchmark-round, and reporting layers
- `internal/services/benchmark/` now also owns the lightweight scheduled round loop and long-run trend aggregation
- `internal/clients/` contains external service clients, including the live `clawmem` HTTP client
- `internal/platform/loader/` loads versioned scenario packs from disk
- `internal/platform/store/` contains local in-memory and file-backed stores for the trust-lab slice
- `configs/scenario-packs/` contains stable and challenger scenario packs
- `reports/` is the Phase 7 report output root
- `reports/` is also the persisted source of truth for historical rounds and Phase 9 report artifacts
- `web/` contains the thin operator UI
- `docs/` contains Phase 2 through 8 architecture and contributor docs

## Web test coverage

- component tests cover the operator pages with mocked API responses for rounds, round detail and comparison, promotions, promotion review submission, and report browsing
- Playwright covers two end-to-end journeys only:
  - round review from rounds list through report viewing
  - promotion review from promotions list through saved review state

The E2E suite stays route-mocked and deterministic. It exercises the real React app without requiring live backend services.

## Docs

- [Architecture](./docs/architecture.md)
- [API](./docs/api.md)
- [Benchmark model](./docs/benchmark-model.md)
- [Commerce model](./docs/commerce-model.md)
- [Development](./docs/development.md)
- [Detection model](./docs/detection-model.md)
- [Event model](./docs/event-model.md)
- [Operator workflows](./docs/operator-workflows.md)
- [Repo layout](./docs/repo-layout.md)
- [Domain model](./docs/domain-model.md)
- [Phase 2 trust lab](./docs/phase-2-trust-lab.md)
- [Phase 3 first flow](./docs/phase-3-first-flow.md)
- [Phase 4.1 clawmem integration](./docs/phase-4-1-clawmem-integration.md)
- [Phase 5 commerce world](./docs/phase-5-commerce-world.md)
- [Phase 6 detection baseline](./docs/phase-6-detection-baseline.md)
- [Phase 7 Red Queen MVP](./docs/phase-7-red-queen-mvp.md)
- [Phase 8 operator surfaces](./docs/phase-8-operator-surfaces.md)
- [Phase 8.1 historical round persistence](./docs/phase-8-1-historical-round-persistence.md)
- [Phase 9 demo hardening](./docs/phase-9-demo-hardening.md)
- [Phase 9 scenario catalog](./docs/phase-9-scenario-catalog.md)
- [Production bridge](./docs/production-bridge.md)
- [Reporting spec](./docs/reporting-spec.md)
- [Scenario pack format](./docs/scenario-pack-format.md)
- [UI architecture](./docs/ui-architecture.md)
