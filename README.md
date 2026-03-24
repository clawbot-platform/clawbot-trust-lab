# clawbot-trust-lab

`clawbot-trust-lab` is the flagship trust-lab vertical for the `clawbot-platform` organization.

Phase 3 adds the first real end-to-end trust-lab slice: load a scenario pack, create a trust artifact, archive a replay case, register a benchmark round through the control-plane client, and invoke the memory client in a real flow.

## What belongs here

- trust-lab specific service code
- scenario, trust, replay, benchmark, and agent domain models
- integrations with `clawbot-server`
- future-facing memory client contracts for `clawmem`
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

Run tests:

```bash
go test ./...
```

## Repo layout

- `cmd/trust-lab/` contains the service entrypoint
- `internal/app/` wires config, bootstrap, router, and graceful shutdown
- `internal/domain/` contains trust-lab domain types and Phase 3 services
- `internal/clients/` contains external service clients and stubs
- `internal/platform/loader/` loads versioned scenario packs from disk
- `internal/platform/store/` contains local in-memory and file-backed stores for the Phase 3 slice
- `configs/scenario-packs/` contains starter scenario pack data
- `docs/` contains Phase 3 architecture and contributor docs

## Docs

- [Architecture](./docs/architecture.md)
- [API](./docs/api.md)
- [Development](./docs/development.md)
- [Repo layout](./docs/repo-layout.md)
- [Domain model](./docs/domain-model.md)
- [Phase 2 trust lab](./docs/phase-2-trust-lab.md)
- [Phase 3 first flow](./docs/phase-3-first-flow.md)
- [Scenario pack format](./docs/scenario-pack-format.md)
