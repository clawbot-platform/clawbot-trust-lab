# Development

## Local workflow

1. Copy `.env.example` to `.env`.
2. Start `clawbot-server` so `/readyz` can pass the control-plane health check.
3. Start `clawmem` so trust and replay creation can persist memory records.
4. Run `go run ./cmd/trust-lab` or `make run`.
5. Run `cd web && npm install && npm run dev` for the Phase 8 operator UI.
6. Execute a Phase 5 scenario through `/api/v1/scenarios/execute`, or run a Phase 7 round through `/api/v1/benchmark/rounds/run`.
7. Use `go test ./...` or `make test` for backend validation.
8. Use `cd web && npm run test` for mocked component/page coverage.
9. Use `cd web && npm run test:e2e` for the two Playwright operator journeys.
10. Scenario packs are loaded from `configs/scenario-packs/`.
11. Round reports are written under `reports/<round-id>/`.

## Docker Version 1 workflow

The supported container path is the Version 1 stack in [`docker-compose.v1.yml`](../docker-compose.v1.yml).

It assumes adjacent checkouts of:

- `clawbot-server`
- `clawbot-trust-lab`
- `clawmem`

Basic Docker flow:

```bash
cp docker-compose.v1.env.example docker-compose.v1.env
docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml up -d --build
python3 ./scripts/version1_validation_report.py --deployment-mode docker --compose-file docker-compose.v1.yml --compose-env-file docker-compose.v1.env --run-round --output-dir ./version1-validation-output
```

## Commands

- `make run`
- `make test`
- `make lint`
- `make security`
- `make docker-build-v1`
- `make docker-up-v1`
- `make docker-down-v1`
- `make docker-ps-v1`
- `make validate-v1`
- `make ui-dev`
- `make ui-build`
- `make ui-test`
- `cd web && npm run lint`
- `cd web && npm run test`
- `cd web && npm run test:e2e`
- `go run ./cmd/trust-lab`

## Required env

- `CONTROL_PLANE_BASE_URL`
- `CLAWMEM_BASE_URL`
- `REPORTS_DIR`

`CLAWMEM_TIMEOUT` defaults to `5s`. A legacy `MEMORY_BASE_URL` fallback is still accepted for compatibility, but new setup should use `CLAWMEM_BASE_URL`.

## Phase 8 web tests

- Component tests use Vitest plus React Testing Library with mocked operator API methods. They validate page rendering and the main review/comparison actions without requiring live services.
- End-to-end tests use Playwright with route interception against the real Vite app. The suite intentionally contains only two tests so it stays credible and fast.
- `npm run test:e2e` assumes a locally available Chrome-compatible browser. If your environment differs, set `PLAYWRIGHT_CHANNEL` or install a Playwright browser separately.
- If you already have the Vite app running, you can point Playwright at it with `PLAYWRIGHT_BASE_URL=http://127.0.0.1:4173 npm run test:e2e`.
