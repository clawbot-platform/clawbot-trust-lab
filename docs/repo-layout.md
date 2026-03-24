# Repo Layout

- `cmd/trust-lab/`: entrypoint for the trust-lab service
- `internal/config/`: environment-driven configuration
- `internal/http/`: versioned API shell, handlers, and middleware
- `internal/domain/`: trust-lab domain models
- `internal/clients/controlplane/`: `clawbot-server` integration boundary
- `internal/clients/memory/`: future `clawmem` contract boundary
- `internal/platform/bootstrap/`: dependency wiring
- `internal/platform/store/`: trust-lab local scaffolding such as the scenario catalog
- `docs/`: phase-specific architecture and contributor docs
- `test/`: future integration assets
