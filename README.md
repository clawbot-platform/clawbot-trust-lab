# clawbot-trust-lab

`clawbot-trust-lab` is the flagship trust-lab vertical for the `clawbot-platform` organization.

Phase 2 establishes the first trust-lab service shell, typed domain foundation, control-plane client boundary, and future `clawmem` integration contract.

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
make run
```

Run tests:

```bash
go test ./...
```

## Repo layout

- `cmd/trust-lab/` contains the service entrypoint
- `internal/app/` wires config, bootstrap, router, and graceful shutdown
- `internal/domain/` contains trust-lab domain types
- `internal/clients/` contains external service clients and stubs
- `internal/platform/` contains trust-lab specific bootstrap and store scaffolding
- `docs/` contains Phase 2 architecture and contributor docs

## Docs

- [Architecture](./docs/architecture.md)
- [Development](./docs/development.md)
- [Repo layout](./docs/repo-layout.md)
- [Domain model](./docs/domain-model.md)
- [Phase 2 trust lab](./docs/phase-2-trust-lab.md)
