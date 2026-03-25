# Development

## Local workflow

1. Copy `.env.example` to `.env`.
2. Start `clawbot-server` so `/readyz` can pass the control-plane health check.
3. Start `clawmem` so trust and replay creation can persist memory records.
4. Run `go run ./cmd/trust-lab` or `make run`.
5. Execute a Phase 5 scenario through `/api/v1/scenarios/execute`, or run a Phase 7 round through `/api/v1/benchmark/rounds/run`.
6. Use `go test ./...` or `make test` for validation.
7. Scenario packs are loaded from `configs/scenario-packs/`.
8. Round reports are written under `reports/<round-id>/`.

## Commands

- `make run`
- `make test`
- `make lint`
- `make security`
- `go run ./cmd/trust-lab`

## Required env

- `CONTROL_PLANE_BASE_URL`
- `CLAWMEM_BASE_URL`
- `REPORTS_DIR`

`CLAWMEM_TIMEOUT` defaults to `5s`. A legacy `MEMORY_BASE_URL` fallback is still accepted for compatibility, but new setup should use `CLAWMEM_BASE_URL`.
