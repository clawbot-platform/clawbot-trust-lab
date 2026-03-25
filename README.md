# clawbot-trust-lab

`clawbot-trust-lab` is the flagship trust-lab vertical for the `clawbot-platform` organization.

Phase 8.1 adds restart-safe historical round persistence on top of the operator surfaces. The repo now reconstructs benchmark and operator history from `reports/<round-id>/` at startup, so prior rounds, promotions, and report artifacts remain visible after service restart.

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

Execute a Phase 5 scenario:

```bash
curl -X POST http://127.0.0.1:8090/api/v1/scenarios/execute \
  -H 'Content-Type: application/json' \
  --data '{"scenario_id":"commerce-clean-agent-assisted-purchase"}'
```

Evaluate the Phase 6 baseline detector:

```bash
curl -X POST http://127.0.0.1:8090/api/v1/detection/evaluate \
  -H 'Content-Type: application/json' \
  --data '{"scenario_id":"commerce-suspicious-refund-attempt"}'
```

Run a Phase 7 Red Queen round:

```bash
curl -X POST http://127.0.0.1:8090/api/v1/benchmark/rounds/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce"}'
```

Historical rounds are reloaded from `REPORTS_DIR` on startup. The structured source of truth is:

- `reports/<round-id>/round-summary.json`
- `reports/<round-id>/promotion-report.json`
- `reports/<round-id>/detection-delta.json`

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
npm run build
```

Run the optional Phase 8 operator E2E smoke tests:

```bash
cd web
npm run test:e2e
```

## Repo layout

- `cmd/trust-lab/` contains the service entrypoint
- `internal/app/` wires config, bootstrap, router, and graceful shutdown
- `internal/domain/` contains trust-lab domain types and Phase 3 plus 4.1 services
- `internal/services/` contains the commerce, event, trust-decision, scenario execution, detection, benchmark-round, and reporting layers
- `internal/clients/` contains external service clients, including the live `clawmem` HTTP client
- `internal/platform/loader/` loads versioned scenario packs from disk
- `internal/platform/store/` contains local in-memory and file-backed stores for the trust-lab slice
- `configs/scenario-packs/` contains stable and challenger scenario packs
- `reports/` is the Phase 7 report output root
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
- [Reporting spec](./docs/reporting-spec.md)
- [Scenario pack format](./docs/scenario-pack-format.md)
- [UI architecture](./docs/ui-architecture.md)
