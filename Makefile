SHELL := /bin/sh

ENV_FILE := .env
GO_ENV := GOCACHE=$(CURDIR)/.cache/go-build GOMODCACHE=$(CURDIR)/.cache/go-mod
COVERAGE_FILE := coverage.out

.PHONY: help check-env check-v1-docker-env run test lint coverage coverage-html security ui-dev ui-build ui-test ui-coverage ui-e2e docker-build-v1 docker-up-v1 docker-down-v1 docker-ps-v1 validate-v1

help: ## Show available targets.
	@awk 'BEGIN {FS = ": ## "}; /^[a-zA-Z0-9_.-]+: ## / {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

check-env:
	@test -f $(ENV_FILE) || { echo "Missing $(ENV_FILE). Copy .env.example to $(ENV_FILE) first."; exit 1; }

check-v1-docker-env:
	@test -f docker-compose.v1.env || { echo "Missing docker-compose.v1.env. Copy docker-compose.v1.env.example to docker-compose.v1.env first."; exit 1; }

run: check-env ## Run the trust-lab service locally.
	@mkdir -p .cache/go-build .cache/go-mod
	@set -a; . ./.env; set +a; $(GO_ENV) go run ./cmd/trust-lab

test: ## Run the Go test suite.
	@mkdir -p .cache/go-build .cache/go-mod
	$(GO_ENV) go test ./...

lint: ## Run formatting, go vet, and golangci-lint when installed.
	@mkdir -p .cache/go-build .cache/go-mod
	@fmt_out=$$(find cmd internal -name '*.go' -print | xargs gofmt -l); \
	if [ -n "$$fmt_out" ]; then \
		echo "$$fmt_out"; \
		echo "gofmt reported unformatted files"; \
		exit 1; \
	fi
	$(GO_ENV) go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; else echo "golangci-lint not installed; skipping"; fi

coverage: ## Run the Go test suite with coverage output.
	@mkdir -p .cache/go-build .cache/go-mod
	$(GO_ENV) go test -covermode=atomic -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

coverage-html: coverage ## Render an HTML coverage report at coverage.html.
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html

security: ## Run local security checks when the tools are installed.
	@if command -v gosec >/dev/null 2>&1; then gosec ./...; else echo "gosec not installed; skipping"; fi
	@if command -v govulncheck >/dev/null 2>&1; then govulncheck ./...; else echo "govulncheck not installed; skipping"; fi
	@if command -v gitleaks >/dev/null 2>&1; then gitleaks detect --no-banner --redact; else echo "gitleaks not installed; skipping"; fi
	@if command -v trivy >/dev/null 2>&1; then trivy fs --exit-code 1 --severity HIGH,CRITICAL .; else echo "trivy not installed; skipping"; fi

ui-dev: ## Run the operator UI dev server.
	cd web && npm run dev

ui-build: ## Build the operator UI.
	cd web && npm run build

ui-test: ## Run the operator UI tests.
	cd web && npm run test

ui-coverage: ## Run operator UI tests with LCOV coverage output.
	cd web && npm run test:coverage

ui-e2e: ## Run the operator UI Playwright smoke tests.
	cd web && npm run test:e2e

docker-build-v1: check-v1-docker-env ## Build the Version 1 Docker images.
	docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml build

docker-up-v1: check-v1-docker-env ## Start the Version 1 Docker stack.
	@mkdir -p reports var/replay-archive var/docker/clawmem
	docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml up -d --build

docker-down-v1: check-v1-docker-env ## Stop the Version 1 Docker stack.
	docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml down

docker-ps-v1: check-v1-docker-env ## Show the Version 1 Docker stack state.
	docker compose --env-file docker-compose.v1.env -f docker-compose.v1.yml ps

validate-v1: check-v1-docker-env ## Run the Version 1 validation report against the Docker stack.
	python3 ./scripts/phase9_validation_report.py --deployment-mode docker --compose-file docker-compose.v1.yml --compose-env-file docker-compose.v1.env --run-round --output-dir ./version1-validation-output
