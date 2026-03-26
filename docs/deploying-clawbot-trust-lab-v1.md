# Deploying Clawbot Trust Lab Version 1

This is the single-document guide for installing and validating Clawbot Trust Lab Version 1.

## What Version 1 is

Version 1 is the current supported release mode of `clawbot-trust-lab`.

It is:

- a self-sufficient DRQ-style trust lab
- a replay-driven adversarial regression harness for agentic commerce fraud and trust controls
- a `shadow` / `recommendation_only` evaluation system

It ships with:

- a built-in scenario catalog
- challenger variants
- replay preservation
- round reporting
- recommendations
- an operator UI

## What Version 1 is not

Version 1 is not:

- a production fraud engine
- a replacement for an incumbent platform
- the planned enterprise sidekick mode

Planned Version 2 is future work and is documented separately as a roadmap direction.

## Required components

Version 1 depends on:

- `clawbot-server`
- `clawmem`
- `clawbot-trust-lab`
- the trust-lab operator UI
- PostgreSQL for `clawbot-server`

The Docker workflow in this repo builds and runs all of those pieces.

## Repository layout assumption

The Docker stack assumes these repositories are checked out side by side:

```text
clawbot-platform/
  clawbot-server/
  clawbot-trust-lab/
  clawmem/
```

The compose file builds sibling services from that adjacent-checkout layout.

## Docker prerequisites

- Docker Engine
- Docker Compose v2 plugin

## Step 1: prepare the Version 1 Docker env file

From the `clawbot-trust-lab` repo root:

```bash
cd /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab
cp docker-compose.v1.env.example docker-compose.v1.env
```

Default values are suitable for a local lab run.

The main values you may want to change are:

- `V1_POSTGRES_PASSWORD`
- `V1_CONTROL_PLANE_PORT`
- `V1_CLAWMEM_PORT`
- `V1_TRUST_LAB_PORT`
- `V1_OPERATOR_UI_PORT`
- scheduler settings for longer runs

## Step 2: start the Version 1 stack

```bash
docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml up -d --build
```

This starts:

- `postgres`
- `control-plane`
- `clawmem`
- `trust-lab`
- `trust-lab-ui`

## Step 3: verify container state

```bash
docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml ps
```

Expected outcome:

- all five services are present
- `control-plane`, `clawmem`, `trust-lab`, and `trust-lab-ui` become healthy

## Step 4: verify health endpoints

```bash
curl http://127.0.0.1:8090/healthz
curl http://127.0.0.1:8090/readyz
curl http://127.0.0.1:8081/healthz
curl http://127.0.0.1:8088/healthz
```

Expected outcome:

- trust-lab returns `{"status":"ok"}`
- trust-lab returns `{"status":"ready"}`
- control-plane health responds
- clawmem health responds

## Step 5: open the operator UI

```text
http://127.0.0.1:8091/
```

The UI is a thin operator surface that proxies API requests to the trust-lab backend.

## Step 6: run Version 1 validation

The repo keeps `scripts/phase9_validation_report.py` for continuity, but it now functions as the Version 1 validation/report script.

Run it like this:

```bash
python3 ./scripts/version1_validation_report.py \
  --deployment-mode docker \
  --compose-file docker-compose.v1.yml \
  --compose-env-file docker-compose.v1.env \
  --run-round \
  --output-dir ./version1-validation-output
```

This validates:

- release docs and key files
- Docker compose service state
- trust-lab health and readiness
- scenario and benchmark APIs
- recommendation and trend APIs
- report artifact presence under `reports/`

Outputs:

- `version1-validation-output/version1-validation-report.md`
- `version1-validation-output/version1-validation-report.html`

Compatibility copies are also written as:

- `version1-validation-output/phase9-validation-report.md`
- `version1-validation-output/phase9-validation-report.html`

## Step 7: confirm Trust Lab is working

Run one round manually:

```bash
ROUND_ID=$(curl -s -X POST http://127.0.0.1:8090/api/v1/benchmark/rounds/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce"}' | jq -r '.data.id')
echo "$ROUND_ID"
```

Inspect operator-facing outputs:

```bash
curl http://127.0.0.1:8090/api/v1/operator/rounds
curl http://127.0.0.1:8090/api/v1/operator/rounds/$ROUND_ID
curl http://127.0.0.1:8090/api/v1/operator/promotions
curl http://127.0.0.1:8090/api/v1/operator/recommendations
curl http://127.0.0.1:8090/api/v1/operator/reports/$ROUND_ID
```

Expected artifacts under the host-mounted reports directory:

- `reports/<round-id>/round-summary.json`
- `reports/<round-id>/round-summary.md`
- `reports/<round-id>/promotion-report.json`
- `reports/<round-id>/detection-delta.json`
- `reports/<round-id>/executive-summary.md`
- `reports/<round-id>/recommendation-report.json`

## Runtime directories

The compose stack mounts these host paths:

- [`reports`](../reports) for round/report artifacts
- [`var/replay-archive`](../var/replay-archive) for replay archive files
- [`var/docker/clawmem`](../var/docker/clawmem) for clawmem persisted storage

That means reports and replay outputs remain visible from the host filesystem.

## Useful commands

Start:

```bash
docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml up -d --build
```

Stop:

```bash
docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml down
```

Logs:

```bash
docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml logs -f trust-lab
```

## How Version 1 differs from planned Version 2

Version 1:

- uses built-in scenarios and challenger variants
- proves the replay-driven adversarial regression loop end to end
- is the current supported mode

Planned Version 2:

- uses incumbent-provided scenarios, features, and data
- is more integration-oriented
- is a future enterprise sidekick direction, not the current release
