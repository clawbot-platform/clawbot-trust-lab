# clawbot-trust-lab

[![ci](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/ci.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/ci.yml)
[![quality](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/quality.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/quality.yml)
[![security](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/security.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/security.yml)
[![docker-compose-validate](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/docker-compose-validate.yml/badge.svg)](https://github.com/clawbot-platform/clawbot-trust-lab/actions/workflows/docker-compose-validate.yml)
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

## DRQ week-run result

The repository now includes a completed week-run benchmark outcome for the commerce scenario family.

### Run window

- March 28, 2026 to April 4, 2026
- image-based distributed homelab deployment
- Docker-only runtime hosts using published GHCR images

### Phase split

- **Phase A** — baseline detector (`v1.0.0-9-g45d296f`)
- **Phase B** — tuned detector (`v1.0.0-10-g796421c`)

### Headline result

| Metric                 | Phase A | Phase B |
|------------------------|--------:|--------:|
| Rounds                 |      13 |      14 |
| Total promotions       |      39 |       0 |
| Avg promotions / round |    3.00 |    0.00 |
| Avg replay pass rate   |    0.08 |    1.00 |
| Zero-promotion rounds  |       0 |      14 |
| Perfect replay rounds  |       1 |      14 |

### What Phase A achieved

Phase A established the baseline benchmark behavior and proved that the harness could repeatedly surface meaningful weaknesses over time rather than only in a one-time run.

It repeatedly exposed three recurring weak cases:

- `commerce-v2-expired-inactive-mandate`
- `commerce-v3-approval-removed`
- `commerce-s3-approval-removed-after-authorization`

Those cases were promoted repeatedly, showing that the baseline detector was underscoring important delegated-commerce risk patterns.

### What Phase B achieved

Phase B used a narrow tuning pass that targeted only the recurring weak cases from Phase A.

The tuned detector:

- reduced promotions from 39 total to 0
- improved replay pass rate from 0.08 average to 1.00
- moved the three targeted weak cases from `suspicious / passed=false` to `step_up_required / passed=true`
- preserved those gains across repeated rounds rather than only one tuned execution

### Weekly run reports

The finished run artifacts are documented in:

- [`docs/drq-week-run-summary.md`](./docs/drq-week-run-summary.md)
- [`docs/DRQ Week-Run Assessment-Executive Report.md`](./docs/DRQ%20Week-Run%20Assessment-Executive%20Report.md)
- [`docs/DRQ Week-Run Effort Review.md`](./docs/DRQ%20Week-Run%20Effort%20Review.md)
- [`docs/DRQ Week-Run Final Assessment.md`](./docs/DRQ%20Week-Run%20Final%20Assessment.md)

These documents are now the reference summary of what Version 1 achieved in a real scheduled run.

## Docker installability

This repo now includes a repo-native Version 1 Docker workflow:

- a core compose stack under [`deploy/compose/docker-compose.yml`](./deploy/compose/docker-compose.yml)
- a local-development override under [`deploy/compose/docker-compose.override.yml`](./deploy/compose/docker-compose.override.yml)
- an optional overlay file under [`deploy/compose/docker-compose.optional.yml`](./deploy/compose/docker-compose.optional.yml)

Version 1 now supports two honest Docker paths:

- local developer builds from adjacent source checkouts
- runtime-host pulls from GHCR

Published images:

- `ghcr.io/clawbot-platform/clawbot-server`
- `ghcr.io/clawbot-platform/clawmem`
- `ghcr.io/clawbot-platform/clawbot-trust-lab`
- `ghcr.io/clawbot-platform/clawbot-trust-lab-ui`

Tag examples:

- baseline: `drq-v1-baseline-20260329`
- tuned: `drq-v1-tuned-20260401`

Runtime hosts such as `optiplex-1`, `optiplex-2`, and `thinkpad-p50` should pull published images instead of building locally. They do not need Go, npm, or other developer tooling for deployment.

Publish from GitHub Actions with the `publish-image` workflow and a `release_tag` input such as:

- `drq-v1-baseline-20260329`
- `drq-v1-tuned-20260401`

If a stale GHCR package already exists from an older manual or CLI push and is not linked to the repository, fix that before relying on the publish workflow:

- connect the package to the repository in GitHub Packages, or
- delete the stale package and republish from Actions, or
- publish once to a temporary image name if cleanup has to be staged

## Core vs optional stack

The default Version 1 stack is intentionally lean:

- `postgres`
- `control-plane`
- `clawmem`
- `trust-lab`
- `trust-lab-ui`

Optional services are separated from the default path. Today, the optional overlay exists for future extensions, but Version 1 does not require any extra services beyond the core stack.

### Local bind-mount mode

For local development and dry-run review, you can enable host-visible outputs with:

- `deploy/compose/docker-compose.local-bind.yml`

This overlays bind mounts for:
- `var/docker/clawmem`
- `reports`
- `var/replay-archive`

CI and default Compose usage should continue using the named-volume core stack.

## Quick start with Docker

### Local developer build path

1. Copy the shared env file:

```bash
cd /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab
cp .env.example .env
```

2. Start the core Version 1 stack:

```bash
make up
```

3. Verify the stack:

```bash
make ps
curl http://127.0.0.1:8090/healthz
curl http://127.0.0.1:8090/readyz
curl http://127.0.0.1:8091/
```

4. Run the core smoke flow:

```bash
make smoke
```

### Runtime-host image pull path

Set the GHCR tags in `.env`, then pull and start:

```bash
CONTROL_PLANE_IMAGE_TAG=drq-v1-baseline-20260329
CLAWMEM_IMAGE_TAG=drq-v1-baseline-20260329
TRUST_LAB_IMAGE_TAG=drq-v1-baseline-20260329
TRUST_LAB_UI_IMAGE_TAG=drq-v1-baseline-20260329
```

```bash
docker compose --env-file .env \
  -f deploy/compose/docker-compose.yml \
  -f deploy/compose/docker-compose.override.yml \
  pull

docker compose --env-file .env \
  -f deploy/compose/docker-compose.yml \
  -f deploy/compose/docker-compose.override.yml \
  up -d
```

5. Run runtime validation against the deployed stack:

```bash
make validate-v1-runtime
```

The script is the Version 1 validation and readiness tool. DRQ run reporting lives inside `clawbot-trust-lab` itself.

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
- round reports under `reports/<round-id>/`
- 24-hour dry-run reports under `reports/daily/<window>/`
- 1-week management reports under `reports/management/<window>/`
- promotion decisions
- replay regressions
- recommendation reports
- executive and machine-readable report artifacts
- historical round reload across restart
- operator-facing review surfaces

## DRQ reporting

Version 1 supports three DRQ report types:

- round report for a single benchmark round
- 24-hour dry-run report for the last 24 hours or an explicit time window
- 1-week management report for the last 7 days or an explicit time window

Generate them from the service:

```bash
go run ./cmd/trust-lab report round --round-id <round-id>
go run ./cmd/trust-lab report dry-run --last 24h
go run ./cmd/trust-lab report management --last 168h
```

Docker-friendly equivalents:

```bash
docker compose --env-file .env -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec trust-lab \
  clawbot-trust-lab report round --round-id <round-id>

docker compose --env-file .env -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec trust-lab \
  clawbot-trust-lab report dry-run --last 24h

docker compose --env-file .env -f deploy/compose/docker-compose.yml -f deploy/compose/docker-compose.override.yml exec trust-lab \
  clawbot-trust-lab report management --last 168h
```

These are separate from `scripts/version1_validation_report.py`. The validator checks whether Version 1 is healthy and installable; the DRQ report commands summarize benchmark evidence and stakeholder-ready findings.

## Validation script

[`scripts/version1_validation_report.py`](scripts/version1_validation_report.py) is the Version 1 validation/readiness script.

It supports two explicit modes:

- `--mode developer` for CI, release workstations, and full repo-quality validation
- `--mode runtime` for appliance-style deployments that only need deployed-system checks

Developer mode can validate:

- docs and release-surface files
- backend and web quality checks
- Docker compose state
- health and readiness endpoints
- round execution
- recommendation and trend endpoints
- presence of expected report artifacts

Runtime mode validates only the deployed system:

- health and readiness endpoints
- `/version` build metadata
- Docker compose state when `--deployment-mode docker` is used
- benchmark, operator, recommendation, promotion, report, and scheduler endpoints
- optional round execution when `--run-round` is supplied

Runtime mode does not require local developer tools like `go`, `golangci-lint`, `gosec`, `govulncheck`, or `npm`.

It writes both:

- `version1-validation-report.md`
- `version1-validation-report.html`

Example developer-mode validation:

```bash
python3 ./scripts/version1_validation_report.py \
  --mode developer \
  --deployment-mode docker \
  --compose-file deploy/compose/docker-compose.yml \
  --compose-override-file deploy/compose/docker-compose.override.yml \
  --compose-env-file .env \
  --run-round \
  --output-dir ./version1-validation-output
```

Example runtime-mode validation:

```bash
make validate-v1-runtime
```

## Local source run

Docker Compose is the current supported deployment path for Version 1.

If you want a non-Docker source run for development:

```bash
cp .env.example .env
make run
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
- [Version 1 product brief](./docs/version-1-deployment-instructions.md)
- [Version 1 scenario catalog](./docs/version-1-scenario-catalog.md)
- [Production bridge](./docs/production-bridge.md)
- [Benchmark model](./docs/benchmark-model.md)
- [Commerce model](./docs/commerce-model.md)

Weekly run references:

- [DRQ Week-Run Summary](./docs/drq-week-run-summary.md)
- [DRQ Week-Run Assessment — Executive Report](./docs/DRQ%20Week-Run%20Assessment-Executive%20Report.md)
- [DRQ Week-Run Effort Review](./docs/DRQ%20Week-Run%20Effort%20Review.md)
- [DRQ Week-Run Final Assessment](./docs/DRQ%20Week-Run%20Final%20Assessment.md)

Historical `docs/phase-*` files remain as implementation history and archive material, not as the main onboarding surface.

## Next steps

The week run established the tuned detector as the stronger benchmark candidate state. The next step is to build from that Phase B baseline rather than restarting from the original baseline.

### Near-term next steps

- expand challenger coverage beyond the initial weak-case family
- strengthen replay-stable governance and regression gating
- improve automated phase comparison and long-run reporting
- make blind-spot history versus current blind spots easier to distinguish
- preserve the tuned detector as the new replay reference point

### LLM roadmap

Version 1 intentionally kept the detector core deterministic.

If an LLM is introduced in a future version, the recommended order is:

1. report and explanation sidecar
2. recommendation clustering and analyst narrative support
3. optional shadow-mode advisory scoring
4. only later, if justified, any deeper detector-path role

The next version should not collapse deterministic benchmark logic and LLM reasoning into one opaque path.

## Planned Version 2

Planned Version 2 is the enterprise sidekick release direction.

It is expected to add:

- incumbent-provided scenarios and data
- richer feature mapping to existing fraud stacks
- more configurable ingestion and evaluation workflows
- enterprise-oriented integration posture

It is not implemented in this repository as the current supported mode.
Version 1 remains the documented and supported release surface today.

## License

This repository is released under the [MIT License](./LICENSE).

