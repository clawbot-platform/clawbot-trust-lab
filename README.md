# clawbot-trust-lab

`clawbot-trust-lab` is the flagship trust-lab vertical for the `clawbot-platform` organization.

Phase 5 adds the first real commerce-world baseline on top of the earlier trust, replay, and `clawmem` integration work. The repo can now execute deterministic commerce scenarios that update local state, emit transaction and trust events, record trust decisions, and persist replay plus memory outputs.

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

Run tests:

```bash
go test ./...
```

## Repo layout

- `cmd/trust-lab/` contains the service entrypoint
- `internal/app/` wires config, bootstrap, router, and graceful shutdown
- `internal/domain/` contains trust-lab domain types and Phase 3 plus 4.1 services
- `internal/services/` contains the Phase 5 commerce, event, trust-decision, and scenario execution layer
- `internal/clients/` contains external service clients, including the live `clawmem` HTTP client
- `internal/platform/loader/` loads versioned scenario packs from disk
- `internal/platform/store/` contains local in-memory and file-backed stores for the trust-lab slice
- `configs/scenario-packs/` contains starter scenario pack data
- `docs/` contains Phase 2 through 5 architecture and contributor docs

## Docs

- [Architecture](./docs/architecture.md)
- [API](./docs/api.md)
- [Commerce model](./docs/commerce-model.md)
- [Development](./docs/development.md)
- [Event model](./docs/event-model.md)
- [Repo layout](./docs/repo-layout.md)
- [Domain model](./docs/domain-model.md)
- [Phase 2 trust lab](./docs/phase-2-trust-lab.md)
- [Phase 3 first flow](./docs/phase-3-first-flow.md)
- [Phase 4.1 clawmem integration](./docs/phase-4-1-clawmem-integration.md)
- [Phase 5 commerce world](./docs/phase-5-commerce-world.md)
- [Scenario pack format](./docs/scenario-pack-format.md)
