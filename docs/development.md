# Development

## Supported Docker workflow

The current supported container workflow is the repo-native Version 1 compose stack:

- [`deploy/compose/docker-compose.yml`](../deploy/compose/docker-compose.yml) for the core stack
- [`deploy/compose/docker-compose.override.yml`](../deploy/compose/docker-compose.override.yml) for local-development tweaks
- [`deploy/compose/docker-compose.optional.yml`](../deploy/compose/docker-compose.optional.yml) for future optional services

The default core stack is:

- `postgres`
- `control-plane`
- `clawmem`
- `trust-lab`
- `trust-lab-ui`

Version 1 currently assumes adjacent checkouts of:

- `clawbot-server`
- `clawbot-trust-lab`
- `clawmem`

## Core Docker flow

```bash
cd /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab
cp .env.example .env
make up
make ps
make smoke
make down
```

For a full developer-mode Docker validation run:

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

For an appliance-style runtime validation run that skips local developer tooling:

```bash
make validate-v1-runtime
```

## Local source workflow

Use this for development only. Docker Compose remains the supported Version 1 deployment path.

1. Copy `.env.example` to `.env`.
2. Start `clawbot-server`.
3. Start `clawmem`.
4. Run `make run`.
5. Run `cd web && npm install && npm run dev` for UI work.

Reports are written to:

- `reports/<round-id>/`
- `reports/daily/<window>/`
- `reports/management/<window>/`

## Common commands

- `make up`
- `make down`
- `make ps`
- `make logs`
- `make smoke`
- `make validate-v1`
- `make validate-v1-runtime`
- `make run`
- `make test`
- `make lint`
- `make security`
- `make report-round ROUND_ID=<round-id>`
- `make report-dry-run`
- `make report-management`
- `make ui-dev`
- `make ui-build`
- `make ui-test`

## Required env

Core env validation is handled by [`scripts/check-env.sh`](../scripts/check-env.sh).

By default it validates only the core stack variables in `.env`.

If optional services are ever enabled, validate them explicitly with:

```bash
VALIDATE_OPTIONAL_STACK=1 sh ./scripts/check-env.sh .env
```

## Version metadata

`/version` now prefers injected build metadata from the Makefile and Docker build path. If those values are not injected, the service falls back to Go build info such as embedded VCS revision and VCS build time when available.
