# clawbot-trust-lab

[![ci](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/ci.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/ci.yml)
[![quality](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/quality.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/quality.yml)
[![security](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/security.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/security.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=clawbot-platform_clawbot-trust-lab&metric=alert_status&token=abb591daa9f6778dcdc919142fe123aa30947073)](https://sonarcloud.io/summary/new_code?id=clawbot-platform_clawbot-trust-lab)

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)
![Docker](https://img.shields.io/badge/Docker-Version_1_Stack-2496ED?logo=docker)
![Shadow Mode](https://img.shields.io/badge/Mode-Shadow%20%2F%20Recommendation--Only-0F766E)

`clawbot-trust-lab` Version 1 is a self-sufficient DRQ-style trust lab for agentic commerce fraud and trust-control benchmarking.

It is the current supported mode of the repo. Version 1 runs its own scenario catalog, challenger variants, replay loop, recommendations, reports, and operator UI in `shadow` / `recommendation_only` mode.

Planned Version 2 is different:

- enterprise sidekick mode
- incumbent-provided scenarios, features, and data
- more configurable ingestion and integration workflows

Version 2 is future work. It is not the current release surface.

## What Version 1 is

Version 1 is:

- a self-running adversarial regression harness
- a replay-preserving benchmark loop for fraud controls
- a recommendation-only shadow evaluator
- a container-installable lab stack for review, demos, homelab runs, and internal evaluation

It is not a replacement for an incumbent fraud engine, and it is not a generic assistant shell.

## Why this repo exists

This repository is the vertical domain layer on top of:

- [`clawbot-server`](../clawbot-server) for the reusable control-plane foundation
- [`clawmem`](../clawmem) for memory, replay, and historical context persistence

The trust lab owns:

- commerce-world scenarios
- trust and replay workflows
- explainable detection
- benchmark rounds and scheduled execution
- promotions, recommendations, and reports
- the thin operator UI

## Docker installability

This repo now includes a real Version 1 Docker workflow:

- a Docker image build path for the trust-lab service
- a Docker image build path for the operator UI
- a compose stack that also builds the required sibling `clawbot-server` and `clawmem` services from adjacent checkouts

Current Docker assumption:

- `clawbot-server`, `clawbot-trust-lab`, and `clawmem` are checked out side by side under the same parent directory

That is the supported Docker install path documented here.

## Quick start with Docker

1. Copy the Version 1 compose env file:

```bash
cd /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab
cp docker-compose.v1.env.example docker-compose.v1.env
```

2. Build and start the Version 1 stack:

```bash
docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml up -d --build
```

3. Verify the stack:

```bash
curl http://127.0.0.1:8090/healthz
curl http://127.0.0.1:8090/readyz
curl http://127.0.0.1:8091/
```

4. Run the Version 1 validation script:

```bash
python3 ./scripts/version1_validation_report.py \
  --deployment-mode docker \
  --compose-file docker-compose.v1.yml \
  --compose-env-file docker-compose.v1.env \
  --run-round \
  --output-dir ./version1-validation-output
```

The script name is preserved for continuity, but it now serves as the Version 1 validation and reporting tool.

## What successful validation looks like

A healthy Version 1 validation run should show:

- trust-lab health and readiness passing
- benchmark and operator APIs reachable
- at least one benchmark round runnable end to end
- reports present under [`reports`](./reports)
- promotions, recommendations, and trend summaries available
- a Markdown and HTML validation report written under `version1-validation-output/`

## Primary outputs

Version 1 produces:

- benchmark rounds
- promotion decisions
- replay regressions
- recommendation reports
- executive and machine-readable report artifacts
- historical round reload across restart
- operator-facing review surfaces

## Validation script

[`scripts/phase9_validation_report.py`](scripts/version1_validation_report.py) is the Version 1 validation/report script.

It can validate:

- docs and release-surface files
- backend and web quality checks
- Docker compose state
- health and readiness endpoints
- round execution
- recommendation and trend endpoints
- presence of expected report artifacts

It writes both:

- `version1-validation-report.md`
- `version1-validation-report.html`

and keeps the older `phase9-validation-report.*` filenames as compatibility outputs.

## Local source run

If you want a non-Docker source run:

```bash
cp .env.example .env
go run ./cmd/trust-lab
```

The expected local startup order is:

1. start `clawbot-server`
2. start `clawmem`
3. start `clawbot-trust-lab`
4. optionally start `web/` with `npm run dev`

## Quality

Backend:

```bash
go test ./...
go vet ./...
golangci-lint run ./...
make coverage
make security
```

Web:

```bash
cd web
npm run lint
npm run test
npm run test:coverage
npm run build
npm run test:e2e
```

SonarCloud ingests both Go and web coverage and enforces the quality gate in CI.

## Documentation

Start here:

- [Deploying Clawbot Trust Lab Version 1](./docs/deploying-clawbot-trust-lab-v1.md)
- [API](./docs/api.md)
- [Architecture](./docs/architecture.md)
- [Benchmark model](./docs/benchmark-model.md)
- [Reporting spec](./docs/reporting-spec.md)
- [Production bridge](./docs/production-bridge.md)

Supporting reference docs:

- [Version 1 product brief](./docs/version-1-deployment-instructions.md)
- [Planned Version 2](./docs/version-2-deployment-instructions.md)
- [Phase 9 scenario catalog](./docs/phase-9-scenario-catalog.md)

Historical `docs/phase-*` files remain as implementation history and archive material, not as the main onboarding surface.

## Planned Version 2

Planned Version 2 is the enterprise sidekick release direction.

It is expected to add:

- incumbent-provided scenarios and data
- richer feature mapping to existing fraud stacks
- more configurable ingestion and evaluation workflows
- enterprise-oriented integration posture

It is not implemented in this repository as the current supported mode.
Version 1 remains the documented and supported release surface today.
