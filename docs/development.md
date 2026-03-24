# Development

## Local workflow

1. Copy `.env.example` to `.env`.
2. Start `clawbot-server` separately if you want `/readyz` to pass against the control-plane health check.
3. Run `make run`.
4. Use `go test ./...` or `make test` for validation.

## Commands

- `make run`
- `make test`
- `make lint`
- `make security`
