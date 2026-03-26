SHELL := /bin/sh

COMPOSE_FILE := deploy/compose/docker-compose.yml
COMPOSE_OVERRIDE := deploy/compose/docker-compose.override.yml
COMPOSE_OPTIONAL := deploy/compose/docker-compose.optional.yml
ENV_FILE := .env
COMPOSE := docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) -f $(COMPOSE_OVERRIDE)
COMPOSE_WITH_OPTIONAL := docker compose --env-file $(ENV_FILE) -f $(COMPOSE_FILE) -f $(COMPOSE_OVERRIDE) -f $(COMPOSE_OPTIONAL)
GO_ENV := GOCACHE=$(CURDIR)/.cache/go-build GOMODCACHE=$(CURDIR)/.cache/go-mod
COVERAGE_FILE := coverage.out

.PHONY: help check-env check-env-optional up up-optional down down-optional restart ps ps-optional logs logs-optional smoke smoke-optional clean compose-validate run test lint coverage coverage-html security ui-dev ui-build ui-test ui-coverage ui-e2e validate-v1 report-round report-dry-run report-management docker-build-v1 docker-up-v1 docker-down-v1 docker-ps-v1

help: ## Show available targets.
	@grep -E '^[a-zA-Z0-9_.-]+:.*## ' $(MAKEFILE_LIST) | sed -E 's/:.*## /\t/' | awk -F '\t' '{printf "  %-18s %s\n", $$1, $$2}'

check-env: ## Validate core env values for the repo-native Version 1 stack.
	@sh ./scripts/check-env.sh $(ENV_FILE)

check-env-optional: ## Validate core + optional env values when optional services are enabled.
	@VALIDATE_OPTIONAL_STACK=1 sh ./scripts/check-env.sh $(ENV_FILE)

up: check-env ## Build and start the core Version 1 stack.
	@mkdir -p reports var/replay-archive var/docker/clawmem
	$(COMPOSE) up -d --build

up-optional: check-env-optional ## Start the core stack plus optional overlays.
	@mkdir -p reports var/replay-archive var/docker/clawmem
	$(COMPOSE_WITH_OPTIONAL) up -d --build

down: check-env ## Stop the core Version 1 stack.
	$(COMPOSE) down --remove-orphans

down-optional: check-env-optional ## Stop the core + optional Version 1 stack.
	$(COMPOSE_WITH_OPTIONAL) down --remove-orphans

restart: down up ## Restart the core Version 1 stack.

ps: check-env ## Show the core Version 1 stack state.
	$(COMPOSE) ps

ps-optional: check-env-optional ## Show the core + optional stack state.
	$(COMPOSE_WITH_OPTIONAL) ps

logs: check-env ## Tail logs for the core Version 1 stack.
	$(COMPOSE) logs -f --tail=100

logs-optional: check-env-optional ## Tail logs for the core + optional Version 1 stack.
	$(COMPOSE_WITH_OPTIONAL) logs -f --tail=100

smoke: check-env ## Run the core stack smoke validation.
	python3 ./scripts/version1_validation_report.py --deployment-mode docker --compose-file $(COMPOSE_FILE) --compose-override-file $(COMPOSE_OVERRIDE) --compose-env-file $(ENV_FILE) --skip-backend --skip-web --run-round --output-dir $${VALIDATION_OUTPUT_DIR:-./version1-validation-output}

smoke-optional: check-env-optional ## Run smoke validation with optional overlays enabled.
	python3 ./scripts/version1_validation_report.py --deployment-mode docker --compose-file $(COMPOSE_FILE) --compose-override-file $(COMPOSE_OVERRIDE) --compose-optional-file $(COMPOSE_OPTIONAL) --include-optional-stack --compose-env-file $(ENV_FILE) --skip-backend --skip-web --run-round --output-dir $${VALIDATION_OUTPUT_DIR:-./version1-validation-output}

clean: check-env ## Remove the core stack and named volumes.
	$(COMPOSE) down -v --remove-orphans

compose-validate: check-env ## Validate the rendered core Compose configuration.
	$(COMPOSE) config >/dev/null

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

validate-v1: check-env ## Run the full Version 1 validation report against the core stack.
	python3 ./scripts/version1_validation_report.py --deployment-mode docker --compose-file $(COMPOSE_FILE) --compose-override-file $(COMPOSE_OVERRIDE) --compose-env-file $(ENV_FILE) --run-round --output-dir $${VALIDATION_OUTPUT_DIR:-./version1-validation-output}

docker-build-v1: check-env ## Backward-compatible alias for building the core Version 1 stack images.
	$(COMPOSE) build

docker-up-v1: up ## Backward-compatible alias for make up.

docker-down-v1: down ## Backward-compatible alias for make down.

docker-ps-v1: ps ## Backward-compatible alias for make ps.

report-round: check-env ## Generate a round report. Set ROUND_ID=<id>.
	@test -n "$(ROUND_ID)" || { echo "Missing ROUND_ID"; exit 1; }
	@mkdir -p .cache/go-build .cache/go-mod
	@set -a; . ./.env; set +a; $(GO_ENV) go run ./cmd/trust-lab report round --round-id $(ROUND_ID)

report-dry-run: check-env ## Generate a dry-run report. Set WINDOW_LAST=24h or use WINDOW_FROM/WINDOW_TO.
	@mkdir -p .cache/go-build .cache/go-mod
	@set -a; . ./.env; set +a; if [ -n "$(WINDOW_FROM)" ] || [ -n "$(WINDOW_TO)" ]; then \
		$(GO_ENV) go run ./cmd/trust-lab report dry-run --from $(WINDOW_FROM) --to $(WINDOW_TO); \
	else \
		$(GO_ENV) go run ./cmd/trust-lab report dry-run --last $${WINDOW_LAST:-24h}; \
	fi

report-management: check-env ## Generate a management report. Set WINDOW_LAST=168h or use WINDOW_FROM/WINDOW_TO.
	@mkdir -p .cache/go-build .cache/go-mod
	@set -a; . ./.env; set +a; if [ -n "$(WINDOW_FROM)" ] || [ -n "$(WINDOW_TO)" ]; then \
		$(GO_ENV) go run ./cmd/trust-lab report management --from $(WINDOW_FROM) --to $(WINDOW_TO); \
	else \
		$(GO_ENV) go run ./cmd/trust-lab report management --last $${WINDOW_LAST:-168h}; \
	fi
