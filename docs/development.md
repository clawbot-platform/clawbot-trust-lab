# Development

## Local workflow

1. Copy `.env.example` to `.env`.
2. Start `clawbot-server` separately if you want `/readyz` to pass against the control-plane health check.
3. Run `go run ./cmd/trust-lab` or `make run`.
4. Use `go test ./...` or `make test` for validation.
5. Scenario packs are loaded from `configs/scenario-packs/`.

## Commands

- `make run`
- `make test`
- `make lint`
- `make security`
- `go run ./cmd/trust-lab`
