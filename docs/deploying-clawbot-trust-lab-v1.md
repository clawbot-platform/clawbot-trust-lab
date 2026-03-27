# Deploying Clawbot Trust Lab Version 1

This is the single-document install and validation guide for the current supported Version 1 release.

## What Version 1 is

Version 1 is:

- a self-sufficient DRQ-style trust lab
- a replay-driven adversarial regression harness for agentic commerce fraud and trust controls
- a `shadow` / `recommendation_only` evaluator with its own scenario catalog, replay loop, reporting, and operator UI

Version 1 is not:

- a production fraud engine
- a replacement for an incumbent stack
- the planned Version 2 enterprise sidekick mode

## Supported deployment model

The supported deployment model is repo-native Docker Compose.

The core compose files are:

- [`deploy/compose/docker-compose.yml`](../deploy/compose/docker-compose.yml)
- [`deploy/compose/docker-compose.override.yml`](../deploy/compose/docker-compose.override.yml)

The optional overlay file is:

- [`deploy/compose/docker-compose.optional.yml`](../deploy/compose/docker-compose.optional.yml)

Version 1 currently supports a local-build Docker path using adjacent source checkouts of:

- `clawbot-server`
- `clawbot-trust-lab`
- `clawmem`

## Core stack

The default core stack is:

- `postgres`
- `control-plane`
- `clawmem`
- `trust-lab`
- `trust-lab-ui`

No optional services are required for:

- local evaluation
- a 24-hour dry run
- a 1-week management run

## Step 1: prepare the env file

```bash
cd /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab
cp .env.example .env
sh ./scripts/check-env.sh .env
```

## Step 2: start the core stack

```bash
make up
```

Equivalent raw Compose command:

```bash
docker compose --env-file .env \
  -f deploy/compose/docker-compose.yml \
  -f deploy/compose/docker-compose.override.yml \
  up -d --build
```

## Step 3: verify container state

```bash
make ps
```

Expected services:

- `postgres`
- `control-plane`
- `clawmem`
- `trust-lab`
- `trust-lab-ui`

## Step 4: verify health

```bash
curl http://127.0.0.1:8090/healthz
curl http://127.0.0.1:8090/readyz
curl http://127.0.0.1:8081/healthz
curl http://127.0.0.1:8088/healthz
curl http://127.0.0.1:8091/
```

## Step 5: run smoke validation

```bash
make smoke
```

This runs the Version 1 validation script in a core-stack mode:

- trust-lab health/readiness
- control-plane and clawmem health
- operator UI reachability

It is a fast readiness wait, not the full Version 1 validation report.

## Step 6: run runtime validation on the deployed stack

```bash
make validate-v1-runtime
```

For a full developer-mode validation from a CI runner or development workstation:

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

Outputs:

- `version1-validation-output/version1-validation-report.md`
- `version1-validation-output/version1-validation-report.html`

Runtime mode validates the deployed services only. Developer mode adds repo-quality checks such as Go, lint, security, and web tooling validation.

## Step 7: confirm Trust Lab is working

Run a round manually:

```bash
ROUND_ID=$(curl -s -X POST http://127.0.0.1:8090/api/v1/benchmark/rounds/run \
  -H 'Content-Type: application/json' \
  --data '{"scenario_family":"commerce"}' | jq -r '.data.id')
echo "$ROUND_ID"
```

Inspect operator-facing outputs:

```bash
curl http://127.0.0.1:8090/api/v1/operator/rounds
curl http://127.0.0.1:8090/api/v1/operator/promotions
curl http://127.0.0.1:8090/api/v1/operator/recommendations
curl http://127.0.0.1:8090/api/v1/operator/reports/$ROUND_ID
```

Expected persisted artifacts include:

- `reports/<round-id>/round-summary.json`
- `reports/<round-id>/round-report.json`
- `reports/<round-id>/promotion-report.json`
- `reports/<round-id>/detection-delta.json`
- `reports/<round-id>/recommendation-report.json`

## Step 8: generate DRQ reports

Round report:

```bash
docker compose --env-file .env \
  -f deploy/compose/docker-compose.yml \
  -f deploy/compose/docker-compose.override.yml \
  exec trust-lab clawbot-trust-lab report round --round-id $ROUND_ID
```

24-hour dry-run report:

```bash
docker compose --env-file .env \
  -f deploy/compose/docker-compose.yml \
  -f deploy/compose/docker-compose.override.yml \
  exec trust-lab clawbot-trust-lab report dry-run --last 24h
```

1-week management report:

```bash
docker compose --env-file .env \
  -f deploy/compose/docker-compose.yml \
  -f deploy/compose/docker-compose.override.yml \
  exec trust-lab clawbot-trust-lab report management --last 168h
```

## Optional overlay

Version 1 does not require optional services by default.

If an optional overlay is added later, start it explicitly:

```bash
make up-optional
```

And validate it explicitly:

```bash
VALIDATE_OPTIONAL_STACK=1 sh ./scripts/check-env.sh .env
make smoke-optional
```

## Stop the stack

```bash
make down
```

## Planned Version 2

Planned Version 2 is the enterprise sidekick direction:

- incumbent-provided scenarios
- incumbent-provided features and data
- a more configurable, integration-oriented evaluation model

That is future work. It is not the current Version 1 deployment surface.
