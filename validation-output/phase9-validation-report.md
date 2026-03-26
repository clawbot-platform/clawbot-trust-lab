# Phase 9 Validation Report

Generated: 2026-03-25T19:50:28-04:00

Repo root: `/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab`

Summary: **26 passed / 0 failed / 26 total**

## Checks

### [PASS] README exists

- Kind: `file`
- Summary: exists

```text
/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/README.md
```

### [PASS] api.md exists

- Kind: `file`
- Summary: exists

```text
/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/docs/api.md
```

### [PASS] phase-9 scenario catalog exists

- Kind: `file`
- Summary: exists

```text
/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/docs/phase-9-scenario-catalog.md
```

### [PASS] api.md includes Recommendation / Trend / Scheduler schemas

- Kind: `doc`
- Summary: all patterns found

```text
/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/docs/api.md
```

### [PASS] go test ./...

- Kind: `command`
- Summary: exit=0
- Command: `go test ./...`

```text
?   	clawbot-trust-lab/cmd/trust-lab	[no test files]
?   	clawbot-trust-lab/internal/app	[no test files]
ok  	clawbot-trust-lab/internal/clients/controlplane	(cached)
ok  	clawbot-trust-lab/internal/clients/memory	(cached)
ok  	clawbot-trust-lab/internal/config	(cached)
?   	clawbot-trust-lab/internal/domain/actors	[no test files]
ok  	clawbot-trust-lab/internal/domain/agents	(cached)
ok  	clawbot-trust-lab/internal/domain/benchmark	(cached)
?   	clawbot-trust-lab/internal/domain/commerce	[no test files]
?   	clawbot-trust-lab/internal/domain/detection	[no test files]
?   	clawbot-trust-lab/internal/domain/events	[no test files]
ok  	clawbot-trust-lab/internal/domain/replay	(cached)
ok  	clawbot-trust-lab/internal/domain/scenario	(cached)
ok  	clawbot-trust-lab/internal/domain/trust	(cached)
ok  	clawbot-trust-lab/internal/http/handlers	(cached)
?   	clawbot-trust-lab/internal/http/middleware	[no test files]
?   	clawbot-trust-lab/internal/http/routes	[no test files]
ok  	clawbot-trust-lab/internal/platform/bootstrap	(cached)
ok  	clawbot-trust-lab/internal/platform/loader	(cached)
ok  	clawbot-trust-lab/internal/platform/store	(cached)
ok  	clawbot-trust-lab/internal/services/benchmark	(cached)
ok  	clawbot-trust-lab/internal/services/commerce	(cached)
ok  	clawbot-trust-lab/internal/services/detection	(cached)
ok  	clawbot-trust-lab/internal/services/events	(cached)
ok  	clawbot-trust-lab/internal/services/operator	(cached)
ok  	clawbot-trust-lab/internal/services/reporting	(cached)
ok  	clawbot-trust-lab/internal/services/scenario	(cached)
ok  	clawbot-trust-lab/internal/services/trust	(cached)
?   	clawbot-trust-lab/internal/version	[no test files]
```

### [PASS] go vet ./...

- Kind: `command`
- Summary: exit=0
- Command: `go vet ./...`

### [PASS] golangci-lint run ./...

- Kind: `command`
- Summary: exit=0
- Command: `golangci-lint run ./...`

```text
0 issues.
```

### [PASS] gosec ./...

- Kind: `command`
- Summary: exit=0
- Command: `gosec ./...`

```text
Results:


[1;36mSummary:[0m
  Gosec  : dev
  Files  : 43
  Lines  : 7581
  Nosec  : 2
  Issues : [1;32m0[0m


[gosec] 2026/03/25 19:50:14 Including rules: default
[gosec] 2026/03/25 19:50:14 Excluding rules: default
[gosec] 2026/03/25 19:50:14 Including analyzers: default
[gosec] 2026/03/25 19:50:14 Excluding analyzers: default
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/cmd/trust-lab
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/operator
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/benchmark
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/actors
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/commerce
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/detection
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/reporting
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/trust
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/agents
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/events
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/detection
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/replay
[gosec] 2026/03/25 19:50:14 Checking package: agents
[gosec] 2026/03/25 19:50:14 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/agents/models.go
[gosec] 2026/03/25 19:50:14 Checking package: actors
[gosec] 2026/03/25 19:50:14 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/actors/models.go
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/trust
[gosec] 2026/03/25 19:50:14 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/events
[gosec] 2026/03/25 19:50:15 Checking package: benchmark
[gosec] 2026/03/25 19:50:15 Checking package: commerce
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/benchmark/models.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/commerce/models.go
[gosec] 2026/03/25 19:50:15 Checking package: events
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/events/models.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/config
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/scenario
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/benchmark/service.go
[gosec] 2026/03/25 19:50:15 Checking package: detection
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/detection/models.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/middleware
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/bootstrap
[gosec] 2026/03/25 19:50:15 Checking package: operator
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/operator/service.go
[gosec] 2026/03/25 19:50:15 Checking package: reporting
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/reporting/service.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/loader
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/clients/controlplane
[gosec] 2026/03/25 19:50:15 Checking package: events
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/events/service.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/store
[gosec] 2026/03/25 19:50:15 Checking package: trust
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/trust/service.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/benchmark
[gosec] 2026/03/25 19:50:15 Checking package: detection
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/detection/service.go
[gosec] 2026/03/25 19:50:15 Checking package: replay
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/replay/models.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/replay/service.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/version
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/app
[gosec] 2026/03/25 19:50:15 Checking package: trust
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/trust/models.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/trust/service.go
[gosec] 2026/03/25 19:50:15 Checking package: version
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/version/version.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/clients/memory
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/handlers
[gosec] 2026/03/25 19:50:15 Checking package: config
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/config/config.go
[gosec] 2026/03/25 19:50:15 Checking package: scenario
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/scenario/models.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/domain/scenario/service.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/commerce
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/scenario
[gosec] 2026/03/25 19:50:15 Checking package: main
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/cmd/trust-lab/main.go
[gosec] 2026/03/25 19:50:15 Import directory: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/routes
[gosec] 2026/03/25 19:50:15 Checking package: loader
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/loader/scenario_packs.go
[gosec] 2026/03/25 19:50:15 Checking package: middleware
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/middleware/logging.go
[gosec] 2026/03/25 19:50:15 Checking package: controlplane
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/clients/controlplane/client.go
[gosec] 2026/03/25 19:50:15 Checking package: bootstrap
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/bootstrap/bootstrap.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/bootstrap/history.go
[gosec] 2026/03/25 19:50:15 Checking package: store
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/store/benchmark_store.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/store/commerce_world_store.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/store/detection_store.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/store/operator_store.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/store/replay_store.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/platform/store/trust_store.go
[gosec] 2026/03/25 19:50:15 Checking package: app
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/app/app.go
[gosec] 2026/03/25 19:50:15 Checking package: benchmark
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/benchmark/scheduler.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/benchmark/service.go
[gosec] 2026/03/25 19:50:15 Checking package: commerce
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/services/commerce/service.go
[gosec] 2026/03/25 19:50:15 Checking package: handlers
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/handlers/common.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/handlers/operator.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/handlers/system.go
[gosec] 2026/03/25 19:50:15 Checking file: /Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/internal/http/handlers/trustlab.go
[gosec] 2026/03/25 19:50:15 Checking p
```

### [PASS] govulncheck ./...

- Kind: `command`
- Summary: exit=0
- Command: `govulncheck ./...`

```text
No vulnerabilities found.
```

### [PASS] npm run lint

- Kind: `command`
- Summary: exit=0
- Command: `npm run lint`

```text
> clawbot-trust-lab-operator@0.1.0 lint
> tsc --noEmit -p tsconfig.app.json && tsc --noEmit -p tsconfig.node.json
```

### [PASS] npm run test

- Kind: `command`
- Summary: exit=0
- Command: `npm run test`

```text
> clawbot-trust-lab-operator@0.1.0 test
> vitest run


[1m[46m RUN [49m[22m [36mv4.1.1 [39m[90m/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/web[39m

 [32m✓[39m src/app/App.test.tsx [2m([22m[2m1 test[22m[2m)[22m[32m 21[2mms[22m[39m
 [32m✓[39m src/pages/RecommendationsPage.test.tsx [2m([22m[2m1 test[22m[2m)[22m[32m 63[2mms[22m[39m
 [32m✓[39m src/pages/ReportsPage.test.tsx [2m([22m[2m1 test[22m[2m)[22m[32m 76[2mms[22m[39m
 [32m✓[39m src/pages/PromotionsPage.test.tsx [2m([22m[2m1 test[22m[2m)[22m[32m 79[2mms[22m[39m
 [32m✓[39m src/pages/RoundsPage.test.tsx [2m([22m[2m1 test[22m[2m)[22m[32m 84[2mms[22m[39m
 [32m✓[39m src/pages/PromotionDetailPage.test.tsx [2m([22m[2m1 test[22m[2m)[22m[32m 87[2mms[22m[39m
 [32m✓[39m src/pages/RoundDetailPage.test.tsx [2m([22m[2m1 test[22m[2m)[22m[32m 89[2mms[22m[39m

[2m Test Files [22m [1m[32m7 passed[39m[22m[90m (7)[39m
[2m      Tests [22m [1m[32m7 passed[39m[22m[90m (7)[39m
[2m   Start at [22m 19:50:18
[2m   Duration [22m 884ms[2m (transform 445ms, setup 1.13s, import 581ms, tests 499ms, environment 2.78s)[22m
```

### [PASS] npm run build

- Kind: `command`
- Summary: exit=0
- Command: `npm run build`

```text
> clawbot-trust-lab-operator@0.1.0 build
> tsc --noEmit -p tsconfig.app.json && tsc --noEmit -p tsconfig.node.json && vite build

vite v8.0.2 building client environment for production...
[2K
transforming...✓ 29 modules transformed.
rendering chunks...
computing gzip size...
dist/index.html                   0.41 kB │ gzip:  0.27 kB
dist/assets/index-K22m4cvL.css    4.22 kB │ gzip:  1.47 kB
dist/assets/index-BCO1siEP.js   219.73 kB │ gzip: 69.41 kB

✓ built in 78ms
```

### [PASS] npm run test:e2e

- Kind: `command`
- Summary: exit=0
- Command: `npm run test:e2e`

```text
> clawbot-trust-lab-operator@0.1.0 test:e2e
> playwright test


Running 2 tests using 2 workers

  ✓  1 tests/e2e/promotion-review.spec.ts:5:1 › operator can inspect a promotion and save a review action (948ms)
  ✓  2 tests/e2e/round-review.spec.ts:5:1 › operator can inspect a round, compare it, and open a report artifact (1.1s)

  2 passed (5.9s)
```

### [PASS] reports directory contains expected artifacts

- Kind: `file`
- Summary: round_dirs=14; missing_rounds=0; legacy_reconstructible=4

```text
round-20260325141726: legacy reconstructible gap for recommendation-report.json (scenario_results present for deterministic bootstrap reconstruction)
round-20260325143447: legacy reconstructible gap for recommendation-report.json (scenario_results present for deterministic bootstrap reconstruction)
round-20260325172529: legacy reconstructible gap for recommendation-report.json (scenario_results present for deterministic bootstrap reconstruction)
round-20260325180315: legacy reconstructible gap for recommendation-report.json (scenario_results present for deterministic bootstrap reconstruction)
```

### [PASS] GET /healthz

- Kind: `api`
- Summary: http 200
- URL: `http://127.0.0.1:8090/healthz`

```text
{
  "status": "ok"
}
```

### [PASS] GET /readyz

- Kind: `api`
- Summary: http 200
- URL: `http://127.0.0.1:8090/readyz`

```text
{
  "status": "ready"
}
```

### [PASS] GET /api/v1/scenarios

- Kind: `api`
- Summary: http 200; scenarios=24; missing_groups=none
- URL: `http://127.0.0.1:8090/api/v1/scenarios`

```text
{
  "data": [
    {
      "id": "commerce-s2-delegated-purchase-weak-provenance",
      "code": "S2",
      "name": "Delegated Purchase with Weak Provenance",
      "type": "commerce_purchase",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v1-weakened-provenance",
      "description": "A delegated purchase carries provenance, but the provenance evidence is materially weak.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "merchant"
      ],
      "trust_signals": [
        "active-mandate",
        "weak-provenance",
        "delegation-visible"
      ],
      "expected_outcomes": [
        "order-created",
        "payment-authorized",
        "candidate-replay-promotion"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "suspicious",
        "weak-provenance"
      ],
      "feature_model": {
        "tier_a": [
          "amount",
          "merchant_category",
          "delegated_indicator"
        ],
        "tier_b": [
          "merchant_scope_history",
          "buyer_history"
        ],
        "tier_c": [
          "provenance_confidence",
          "mandate_status",
          "delegation_mode"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-s3-approval-removed-after-authorization",
      "code": "S3",
      "name": "Approval Removed After Initial Authorization",
      "type": "commerce_refund_review",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v3-approval-removed",
      "description": "A refund begins with valid authority but approval evidence disappears before execution.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "support-reviewer"
      ],
      "trust_signals": [
        "active-mandate",
        "approval-removed",
        "agent-refund"
      ],
      "expected_outcomes": [
        "refund-requested",
        "trust-decision-step-up",
        "candidate-replay-promotion"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "approval-removed"
      ],
      "feature_model": {
        "tier_a": [
          "refund_indicator",
          "amount",
          "delegated_indicator"
        ],
        "tier_b": [
          "approval_history",
          "repeat_attempt_count"
        ],
        "tier_c": [
          "approval_evidence",
          "mandate_status",
          "delegation_mode"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-s5-merchant-scope-drift-delegated-action",
      "code": "S5",
      "name": "Merchant or Category Scope Drift Under Delegated Action",
      "type": "commerce_purchase",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v6-merchant-scope-drift",
      "description": "A delegated purchase attempts to move outside the buyer's prior merchant or category scope.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "merchant"
      ],
      "trust_signals": [
        "delegation-visible",
        "merchant-scope-drift",
        "category-drift"
      ],
      "expected_outcomes": [
        "order-created",
        "trust-decision-review-required",
        "candidate-replay-promotion"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "scope-drift"
      ],
      "feature_model": {
        "tier_a": [
          "amount",
          "merchant_category",
          "delegated_indicator"
        ],
        "tier_b": [
          "merchant_scope_history",
          "category_history",
          "recent_attempt_count"
        ],
        "tier_c": [
          "mandate_status",
          "provenance_confidence",
          "delegation_mode"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-v1-weakened-provenance",
      "code": "V1",
      "name": "Variant Weakened Provenance",
      "type": "commerce_purchase",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v1-weakened-provenance",
      "description": "Variant that weakens provenance while keeping the rest of a delegated purchase intact.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "merchant"
      ],
      "trust_signals": [
        "weak-provenance"
      ],
      "expected_outcomes": [
        "order-created",
        "trust-gap-low-provenance"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "variant",
        "weakened-provenance"
      ],
      "feature_model": {
        "tier_a": [
          "amount",
          "merchant_category",
          "delegated_indicator"
        ],
        "tier_b": [
          "buyer_history"
        ],
        "tier_c": [
          "provenance_confidence"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-v2-expired-inactive-mandate",
      "code": "V2",
      "name": "Variant Expired or Inactive Mandate",
      "type": "commerce_purchase",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v2-expired-mandate",
      "description": "Variant that expires mandate coverage before a delegated action executes.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "merchant"
      ],
      "trust_signals": [
        "expired-mandate"
      ],
      "expected_outcomes": [
        "order-created",
        "trust-decision-review-required"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "variant",
        "expired-mandate"
      ],
      "feature_model": {
        "tier_a": [
          "amount",
          "merchant_category",
          "delegated_indicator"
        ],
        "tier_b": [
          "recent_attempt_count"
        ],
        "tier_c": [
          "mandate_status",
          "delegation_mode"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-v3-approval-removed",
      "code": "V3",
      "name": "Variant Approval Removed",
      "type": "commerce_refund_review",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v3-approval-removed",
      "description": "Variant that removes approval evidence from an agent-driven refund.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "support-reviewer"
      ],
      "trust_signals": [
        "approval-removed",
        "agent-refund"
      ],
      "expected_outcomes": [
        "refund-requested",
        "trust-decision-step-up"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "variant",
        "approval-removed"
      ],
      "feature_model": {
        "tier_a": [
          "refund_indicator",
          "delegated_indicator"
        ],
        "tier_b": [
          "approval_history"
        ],
        "tier_c": [
          "approval_evidence",
          "delegation_mode"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-v4-actor-switch-human-to-agent",
      "code": "V4",
      "name": "Variant Actor Switch from Human to Agent",
      "type": "commerce_refund_review",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v4-actor-switch",
      "description": "Variant that flips a previously human refund path into an agent-driven refund without strengthening controls.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "support-reviewer"
      ],
      "trust_signals": [
        "actor-switch",
        "delegation-visible"
      ],
      "expected_outcomes": [
        "refund-requested",
        "trust-decision-review-required"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "variant",
        "actor-switch"
      ],
      "feature_model": {
        "tier_a": [
          "refund_indicator",
          "delegated_indicator"
        ],
        "tier_b": [
          "recent_attempt_count",
          "approval_history"
        ],
        "tier_c": [
          "delegation_mode",
          "approval_evidence"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-v5-repeat-attempt-escalation",
      "code": "V5",
      "name": "Variant Repeat Attempt Escalation",
      "type": "commerce_refund_review",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v5-repeat-attempt-escalation",
      "description": "Variant that escalates the number of prior similar refund attempts.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "support-reviewer"
      ],
      "trust_signals": [
        "repeat-refund-pattern",
        "agent-refund"
      ],
      "expected_outcomes": [
        "refund-requested",
        "trust-decision-step-up"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "variant",
        "repeat-attempt"
      ],
      "feature_model": {
        "tier_a": [
          "refund_indicator",
          "amount",
          "delegated_indicator"
        ],
        "tier_b": [
          "repeat_attempt_count",
          "historical_refund_rate"
        ],
        "tier_c": [
          "delegation_mode"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-v6-merchant-scope-drift",
      "code": "V6",
      "name": "Variant Merchant Scope Drift",
      "type": "commerce_purchase",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v6-merchant-scope-drift",
      "description": "Variant that moves a delegated purchase into a new merchant scope.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "merchant"
      ],
      "trust_signals": [
        "merchant-scope-drift",
        "delegation-visible"
      ],
      "expected_outcomes": [
        "order-created",
        "trust-decision-review-required"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "challenger",
        "variant",
        "scope-drift"
      ],
      "feature_model": {
        "tier_a": [
          "amount",
          "merchant_category",
          "delegated_indicator"
        ],
        "tier_b": [
          "merchant_scope_history",
          "category_history"
        ],
        "tier_c": [
          "delegation_mode",
          "mandate_status"
        ]
      },
      "created_at": "2026-03-25T22:52:59.250392Z"
    },
    {
      "id": "commerce-v7-high-value-delegated-purchase",
      "code": "V7",
      "name": "Variant High-Value Delegated Purchase",
      "type": "commerce_purchase",
      "family": "commerce",
      "set_role": "living",
      "variant_id": "variant-v7-high-value-delegated-purchase",
      "description": "Variant that materially increases delegated purchase value beyond the usual buyer pattern.",
      "pack_id": "challenger-pack",
      "version": "v2",
      "actors": [
        "buyer",
        "agent",
        "merchant"
      ],
      "trust_signals": [
        "high-value-purchase",
        "delegation-visible"
      ],
      "expected_outcomes": [
        "order-created",
        "trust-decision-review-required"
      ],
      "tags": [
        "commerce",
        "phase-9",
        "chall
```

### [PASS] POST /api/v1/benchmark/rounds/run

- Kind: `api`
- Summary: http 201
- URL: `http://127.0.0.1:8090/api/v1/benchmark/rounds/run`

```text
{
  "data": {
    "id": "round-20260325235027",
    "scenario_family": "commerce",
    "detector_version": "dev",
    "stable_scenario_refs": [
      "commerce-a1-agent-assisted-purchase-valid-controls",
      "commerce-a2-fully-delegated-replenishment-purchase",
      "commerce-a3-agent-assisted-refund-approval-evidence",
      "commerce-h1-direct-human-purchase",
      "commerce-h2-human-refund-valid-history",
      "commerce-s1-refund-weak-authorization",
      "commerce-s4-repeated-agent-refund-attempts"
    ],
    "challenger_variant_refs": [
      "variant-v1-weakened-provenance",
      "variant-v2-expired-mandate",
      "variant-v3-approval-removed",
      "variant-v4-actor-switch",
      "variant-v5-repeat-attempt-escalation",
      "variant-v6-merchant-scope-drift",
      "variant-v7-high-value-delegated-purchase",
      "variant-s2-weak-provenance",
      "variant-s3-approval-removed-after-authorization",
      "variant-s5-scope-drift"
    ],
    "replay_case_refs": [
      "rc-s3-approval-removed-after-authorization",
      "rc-v2-expired-inactive-mandate",
      "rc-v3-approval-removed"
    ],
    "started_at": "2026-03-25T23:50:27.907712Z",
    "completed_at": "2026-03-25T23:50:28.385017Z",
    "round_status": "completed",
    "report_dir": "reports/round-20260325235027",
    "scenario_results": [
      {
        "id": "sr-stable-commerce-a1-agent-assisted-purchase-valid-controls",
        "scenario_id": "commerce-a1-agent-assisted-purchase-valid-controls",
        "scenario_family": "commerce",
        "set_kind": "stable",
        "execution_status": "completed",
        "order_refs": [
          "order-a1-agent-assisted-purchase-valid-controls"
        ],
        "refund_refs": null,
        "trust_decision_refs": [
          "decision-a1-agent-assisted-purchase-valid-controls"
        ],
        "replay_case_refs": [
          "rc-a1-agent-assisted-purchase-valid-controls"
        ],
        "memory_record_refs": [
          "ta-a1-agent-assisted-purchase-valid-controls",
          "rc-a1-agent-assisted-purchase-valid-controls"
        ],
        "detection_result_ref": "det-order-a1-agent-assisted-purchase-valid-controls",
        "final_detection_status": "clean",
        "final_recommendation": "allow",
        "triggered_rule_ids": null,
        "promoted_to_replay": false,
        "expected_minimum_status": "clean",
        "passed": true,
        "notes": [
          "tier_c_used",
          "stable baseline met its expected detector outcome"
        ]
      },
      {
        "id": "sr-stable-commerce-a2-fully-delegated-replenishment-purchase",
        "scenario_id": "commerce-a2-fully-delegated-replenishment-purchase",
        "scenario_family": "commerce",
        "set_kind": "stable",
        "execution_status": "completed",
        "order_refs": [
          "order-a2-fully-delegated-replenishment-purchase"
        ],
        "refund_refs": null,
        "trust_decision_refs": [
          "decision-a2-fully-delegated-replenishment-purchase"
        ],
        "replay_case_refs": [
          "rc-a2-fully-delegated-replenishment-purchase"
        ],
        "memory_record_refs": [
          "ta-a2-fully-delegated-replenishment-purchase",
          "rc-a2-fully-delegated-replenishment-purchase"
        ],
        "detection_result_ref": "det-order-a2-fully-delegated-replenishment-purchase",
        "final_detection_status": "clean",
        "final_recommendation": "allow",
        "triggered_rule_ids": null,
        "promoted_to_replay": false,
        "expected_minimum_status": "clean",
        "passed": true,
        "notes": [
          "tier_c_used",
          "stable baseline met its expected detector outcome"
        ]
      },
      {
        "id": "sr-stable-commerce-a3-agent-assisted-refund-approval-evidence",
        "scenario_id": "commerce-a3-agent-assisted-refund-approval-evidence",
        "scenario_family": "commerce",
        "set_kind": "stable",
        "execution_status": "completed",
        "order_refs": [
          "order-a3-agent-assisted-refund-approval-evidence"
        ],
        "refund_refs": [
          "refund-a3-agent-assisted-refund-approval-evidence"
        ],
        "trust_decision_refs": [
          "decision-a3-agent-assisted-refund-approval-evidence"
        ],
        "replay_case_refs": [
          "rc-a3-agent-assisted-refund-approval-evidence"
        ],
        "memory_record_refs": [
          "ta-a3-agent-assisted-refund-approval-evidence",
          "rc-a3-agent-assisted-refund-approval-evidence"
        ],
        "detection_result_ref": "det-order-a3-agent-assisted-refund-approval-evidence",
        "final_detection_status": "clean",
        "final_recommendation": "allow",
        "triggered_rule_ids": null,
        "promoted_to_replay": false,
        "expected_minimum_status": "clean",
        "passed": true,
        "notes": [
          "tier_c_used",
          "stable baseline met its expected detector outcome"
        ]
      },
      {
        "id": "sr-stable-commerce-h1-direct-human-purchase",
        "scenario_id": "commerce-h1-direct-human-purchase",
        "scenario_family": "commerce",
        "set_kind": "stable",
        "execution_status": "completed",
        "order_refs": [
          "order-h1-direct-human-purchase"
        ],
        "refund_refs": null,
        "trust_decision_refs": [
          "decision-h1-direct-human-purchase"
        ],
        "replay_case_refs": [
          "rc-h1-direct-human-purchase"
        ],
        "memory_record_refs": [
          "ta-h1-direct-human-purchase",
          "rc-h1-direct-human-purchase"
        ],
        "detection_result_ref": "det-order-h1-direct-human-purchase",
        "final_detection_status": "clean",
        "final_recommendation": "allow",
        "triggered_rule_ids": null,
        "promoted_to_replay": false,
        "expected_minimum_status": "clean",
        "passed": true,
        "notes": [
          "stable baseline met its expected detector outcome"
        ]
      },
      {
        "id": "sr-stable-commerce-h2-human-refund-valid-history",
        "scenario_id": "commerce-h2-human-refund-valid-history",
        "scenario_family": "commerce",
        "set_kind": "stable",
        "execution_status": "completed",
        "order_refs": [
          "order-h2-human-refund-valid-history"
        ],
        "refund_refs": [
          "refund-h2-human-refund-valid-history"
        ],
        "trust_decision_refs": [
          "decision-h2-human-refund-valid-history"
        ],
        "replay_case_refs": [
          "rc-h2-human-refund-valid-history"
        ],
        "memory_record_refs": [
          "ta-h2-human-refund-valid-history",
          "rc-h2-human-refund-valid-history"
        ],
        "detection_result_ref": "det-order-h2-human-refund-valid-history",
        "final_detection_status": "clean",
        "final_recommendation": "allow",
        "triggered_rule_ids": null,
        "promoted_to_replay": false,
        "expected_minimum_status": "clean",
        "passed": true,
        "notes": [
          "stable baseline met its expected detector outcome"
        ]
      },
      {
        "id": "sr-stable-commerce-s1-refund-weak-authorization",
        "scenario_id": "commerce-s1-refund-weak-authorization",
        "scenario_family": "commerce",
        "set_kind": "stable",
        "execution_status": "completed",
        "order_refs": [
          "order-s1-refund-weak-authorization"
        ],
        "refund_refs": [
          "refund-s1-refund-weak-authorization"
        ],
        "trust_decision_refs": [
          "decision-s1-refund-weak-authorization"
        ],
        "replay_case_refs": [
          "rc-s1-refund-weak-authorization"
        ],
        "memory_record_refs": [
          "ta-s1-refund-weak-authorization",
          "rc-s1-refund-weak-authorization"
        ],
        "detection_result_ref": "det-order-s1-refund-weak-authorization",
        "final_detection_status": "step_up_required",
        "final_recommendation": "step_up",
        "triggered_rule_ids": [
          "agent_refund_without_approval",
          "missing_mandate_delegated_action",
          "missing_provenance_sensitive_action",
          "prior_step_up_decision",
          "refund_weak_authorization"
        ],
        "promoted_to_replay": false,
        "expected_minimum_status": "step_up_required",
        "passed": true,
        "notes": [
          "tier_c_used",
          "stable baseline met its expected detector outcome"
        ]
      },
      {
        "id": "sr-stable-commerce-s4-repeated-agent-refund-attempts",
        "scenario_id": "commerce-s4-repeated-agent-refund-attempts",
        "scenario_family": "commerce",
        "set_kind": "stable",
        "execution_status": "completed",
        "order_refs": [
          "order-s4-repeated-agent-refund-attempts"
        ],
        "refund_refs": [
          "refund-s4-repeated-agent-refund-attempts"
        ],
        "trust_decision_refs": [
          "decision-s4-repeated-agent-refund-attempts"
        ],
        "replay_case_refs": [
          "rc-s4-repeated-agent-refund-attempts"
        ],
        "memory_record_refs": [
          "ta-s4-repeated-agent-refund-attempts",
          "rc-s4-repeated-agent-refund-attempts"
        ],
        "detection_result_ref": "det-order-s4-repeated-agent-refund-attempts",
        "final_detection_status": "suspicious",
        "final_recommendation": "review",
        "triggered_rule_ids": [
          "prior_step_up_decision",
          "repeat_suspicious_context"
        ],
        "promoted_to_replay": false,
        "expected_minimum_status": "suspicious",
        "passed": true,
        "notes": [
          "tier_c_used",
          "stable baseline met its expected detector outcome"
        ]
      },
      {
        "id": "sr-living-commerce-v1-weakened-provenance",
        "scenario_id": "commerce-v1-weakened-provenance",
        "scenario_family": "commerce",
        "set_kind": "living",
        "challenger_variant_id": "variant-v1-weakened-provenance",
        "execution_status": "completed",
        "order_refs": [
          "order-v1-weakened-provenance"
        ],
        "refund_refs": null,
        "trust_decision_refs": [
          "decision-v1-weakened-provenance"
        ],
        "replay_case_refs": [
          "rc-v1-weakened-provenance"
        ],
        "memory_record_refs": [
          "ta-v1-weakened-provenance",
          "rc-v1-weakened-provenance"
        ],
        "detection_result_ref": "det-order-v1-weakened-provenance",
        "final_detection_status": "suspicious",
        "final_recommendation": "review",
        "triggered_rule_ids": [
          "missing_provenance_sensitive_action"
        ],
        "promoted_to_replay": false,
        "expected_minimum_status": "suspicious",
        "passed": true,
        "notes": [
          "tier_c_used",
          "challenger variant met the expected minimum detector posture"
        ]
      },
      {
        "id": "sr-living-commerce-v2-expired-inactive-mandate",
        "scenario_id": "commerce-v2-expired-inactive-mandate",
        "scenario_family": "commerce",
        "set_kind": "living",
        "challenger_variant_id": "variant-v2-expired-mandate",
        "execution_status": "completed",
        "order_refs": [
          "order-v2-expired-inactive-mandate"
        ],
        "refund_refs": null,
        "trust_decision_refs": [
          "decision-v2-expired-inactive-mandate"
        ],
        "replay_case_refs": [
          "rc-v2-expired-inactive-mandate"
        ],
        "memory_record_refs": [
          "ta-v2-expired-inactive-mandate",
          "rc-v2-expired-inactive-mandate"
        ],
        "detection_result_ref": "det-order-v2-expired-inactive-mandate",
        "final_detection_status": "suspicious",
        "final_recommendation": "review",
        "triggered_rule_ids": [
          "missing_mandate_delegated_action",

```

### [PASS] GET /api/v1/benchmark/rounds

- Kind: `api`
- Summary: http 200; rounds=15; rounds_with_prod_bridge=11; rounds_with_tier_c_usage=11; recommendation_totals=59
- URL: `http://127.0.0.1:8090/api/v1/benchmark/rounds`

```text
{
  "data": [
    {
      "id": "round-20260325235027",
      "scenario_family": "commerce",
      "detector_version": "dev",
      "stable_scenario_refs": [
        "commerce-a1-agent-assisted-purchase-valid-controls",
        "commerce-a2-fully-delegated-replenishment-purchase",
        "commerce-a3-agent-assisted-refund-approval-evidence",
        "commerce-h1-direct-human-purchase",
        "commerce-h2-human-refund-valid-history",
        "commerce-s1-refund-weak-authorization",
        "commerce-s4-repeated-agent-refund-attempts"
      ],
      "challenger_variant_refs": [
        "variant-v1-weakened-provenance",
        "variant-v2-expired-mandate",
        "variant-v3-approval-removed",
        "variant-v4-actor-switch",
        "variant-v5-repeat-attempt-escalation",
        "variant-v6-merchant-scope-drift",
        "variant-v7-high-value-delegated-purchase",
        "variant-s2-weak-provenance",
        "variant-s3-approval-removed-after-authorization",
        "variant-s5-scope-drift"
      ],
      "replay_case_refs": [
        "rc-s3-approval-removed-after-authorization",
        "rc-v2-expired-inactive-mandate",
        "rc-v3-approval-removed"
      ],
      "started_at": "2026-03-25T23:50:27.907712Z",
      "completed_at": "2026-03-25T23:50:28.385017Z",
      "round_status": "completed",
      "report_dir": "reports/round-20260325235027",
      "scenario_results": [
        {
          "id": "sr-stable-commerce-a1-agent-assisted-purchase-valid-controls",
          "scenario_id": "commerce-a1-agent-assisted-purchase-valid-controls",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-a1-agent-assisted-purchase-valid-controls"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-a1-agent-assisted-purchase-valid-controls"
          ],
          "replay_case_refs": [
            "rc-a1-agent-assisted-purchase-valid-controls"
          ],
          "memory_record_refs": [
            "ta-a1-agent-assisted-purchase-valid-controls",
            "rc-a1-agent-assisted-purchase-valid-controls"
          ],
          "detection_result_ref": "det-order-a1-agent-assisted-purchase-valid-controls",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-a2-fully-delegated-replenishment-purchase",
          "scenario_id": "commerce-a2-fully-delegated-replenishment-purchase",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-a2-fully-delegated-replenishment-purchase"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-a2-fully-delegated-replenishment-purchase"
          ],
          "replay_case_refs": [
            "rc-a2-fully-delegated-replenishment-purchase"
          ],
          "memory_record_refs": [
            "ta-a2-fully-delegated-replenishment-purchase",
            "rc-a2-fully-delegated-replenishment-purchase"
          ],
          "detection_result_ref": "det-order-a2-fully-delegated-replenishment-purchase",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-a3-agent-assisted-refund-approval-evidence",
          "scenario_id": "commerce-a3-agent-assisted-refund-approval-evidence",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-a3-agent-assisted-refund-approval-evidence"
          ],
          "refund_refs": [
            "refund-a3-agent-assisted-refund-approval-evidence"
          ],
          "trust_decision_refs": [
            "decision-a3-agent-assisted-refund-approval-evidence"
          ],
          "replay_case_refs": [
            "rc-a3-agent-assisted-refund-approval-evidence"
          ],
          "memory_record_refs": [
            "ta-a3-agent-assisted-refund-approval-evidence",
            "rc-a3-agent-assisted-refund-approval-evidence"
          ],
          "detection_result_ref": "det-order-a3-agent-assisted-refund-approval-evidence",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-h1-direct-human-purchase",
          "scenario_id": "commerce-h1-direct-human-purchase",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-h1-direct-human-purchase"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-h1-direct-human-purchase"
          ],
          "replay_case_refs": [
            "rc-h1-direct-human-purchase"
          ],
          "memory_record_refs": [
            "ta-h1-direct-human-purchase",
            "rc-h1-direct-human-purchase"
          ],
          "detection_result_ref": "det-order-h1-direct-human-purchase",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-h2-human-refund-valid-history",
          "scenario_id": "commerce-h2-human-refund-valid-history",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-h2-human-refund-valid-history"
          ],
          "refund_refs": [
            "refund-h2-human-refund-valid-history"
          ],
          "trust_decision_refs": [
            "decision-h2-human-refund-valid-history"
          ],
          "replay_case_refs": [
            "rc-h2-human-refund-valid-history"
          ],
          "memory_record_refs": [
            "ta-h2-human-refund-valid-history",
            "rc-h2-human-refund-valid-history"
          ],
          "detection_result_ref": "det-order-h2-human-refund-valid-history",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-s1-refund-weak-authorization",
          "scenario_id": "commerce-s1-refund-weak-authorization",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-s1-refund-weak-authorization"
          ],
          "refund_refs": [
            "refund-s1-refund-weak-authorization"
          ],
          "trust_decision_refs": [
            "decision-s1-refund-weak-authorization"
          ],
          "replay_case_refs": [
            "rc-s1-refund-weak-authorization"
          ],
          "memory_record_refs": [
            "ta-s1-refund-weak-authorization",
            "rc-s1-refund-weak-authorization"
          ],
          "detection_result_ref": "det-order-s1-refund-weak-authorization",
          "final_detection_status": "step_up_required",
          "final_recommendation": "step_up",
          "triggered_rule_ids": [
            "agent_refund_without_approval",
            "missing_mandate_delegated_action",
            "missing_provenance_sensitive_action",
            "prior_step_up_decision",
            "refund_weak_authorization"
          ],
          "promoted_to_replay": false,
          "expected_minimum_status": "step_up_required",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-s4-repeated-agent-refund-attempts",
          "scenario_id": "commerce-s4-repeated-agent-refund-attempts",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-s4-repeated-agent-refund-attempts"
          ],
          "refund_refs": [
            "refund-s4-repeated-agent-refund-attempts"
          ],
          "trust_decision_refs": [
            "decision-s4-repeated-agent-refund-attempts"
          ],
          "replay_case_refs": [
            "rc-s4-repeated-agent-refund-attempts"
          ],
          "memory_record_refs": [
            "ta-s4-repeated-agent-refund-attempts",
            "rc-s4-repeated-agent-refund-attempts"
          ],
          "detection_result_ref": "det-order-s4-repeated-agent-refund-attempts",
          "final_detection_status": "suspicious",
          "final_recommendation": "review",
          "triggered_rule_ids": [
            "prior_step_up_decision",
            "repeat_suspicious_context"
          ],
          "promoted_to_replay": false,
          "expected_minimum_status": "suspicious",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-living-commerce-v1-weakened-provenance",
          "scenario_id": "commerce-v1-weakened-provenance",
          "scenario_family": "commerce",
          "set_kind": "living",
          "challenger_variant_id": "variant-v1-weakened-provenance",
          "execution_status": "completed",
          "order_refs": [
            "order-v1-weakened-provenance"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-v1-weakened-provenance"
          ],
          "replay_case_refs": [
            "rc-v1-weakened-provenance"
          ],
          "memory_record_refs": [
            "ta-v1-weakened-provenance",
            "rc-v1-weakened-provenance"
          ],
          "detection_result_ref": "det-order-v1-weakened-provenance",
          "final_detection_status": "suspicious",
          "final_recommendation": "review",
          "triggered_rule_ids": [
            "missing_provenance_sensitive_action"
          ],
          "promoted_to_replay": false,
          "expected_minimum_status": "suspicious",
          "passed": true,
          "notes": [
            "tier_c_used",
            "challenger variant met the expected minimum detector posture"
          ]
        },
        {
          "id": "sr-living-commerce-v2-expired-inactive-mandate",
          "scenario_id": "commerce-v2-expired-inactive-mandate",
          "scenario_family": "commerce",
          "set_kind": "living",
          "challenger_variant_id": "variant-v2-expired-mandate",
          "execution_status": "completed",
          "order_refs"
```

### [PASS] GET /api/v1/operator/rounds

- Kind: `api`
- Summary: http 200; rounds=15; rounds_with_prod_bridge=11; rounds_with_tier_c_usage=11; recommendation_totals=59
- URL: `http://127.0.0.1:8090/api/v1/operator/rounds`

```text
{
  "data": [
    {
      "id": "round-20260325235027",
      "scenario_family": "commerce",
      "detector_version": "dev",
      "stable_scenario_refs": [
        "commerce-a1-agent-assisted-purchase-valid-controls",
        "commerce-a2-fully-delegated-replenishment-purchase",
        "commerce-a3-agent-assisted-refund-approval-evidence",
        "commerce-h1-direct-human-purchase",
        "commerce-h2-human-refund-valid-history",
        "commerce-s1-refund-weak-authorization",
        "commerce-s4-repeated-agent-refund-attempts"
      ],
      "challenger_variant_refs": [
        "variant-v1-weakened-provenance",
        "variant-v2-expired-mandate",
        "variant-v3-approval-removed",
        "variant-v4-actor-switch",
        "variant-v5-repeat-attempt-escalation",
        "variant-v6-merchant-scope-drift",
        "variant-v7-high-value-delegated-purchase",
        "variant-s2-weak-provenance",
        "variant-s3-approval-removed-after-authorization",
        "variant-s5-scope-drift"
      ],
      "replay_case_refs": [
        "rc-s3-approval-removed-after-authorization",
        "rc-v2-expired-inactive-mandate",
        "rc-v3-approval-removed"
      ],
      "started_at": "2026-03-25T23:50:27.907712Z",
      "completed_at": "2026-03-25T23:50:28.385017Z",
      "round_status": "completed",
      "report_dir": "reports/round-20260325235027",
      "scenario_results": [
        {
          "id": "sr-stable-commerce-a1-agent-assisted-purchase-valid-controls",
          "scenario_id": "commerce-a1-agent-assisted-purchase-valid-controls",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-a1-agent-assisted-purchase-valid-controls"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-a1-agent-assisted-purchase-valid-controls"
          ],
          "replay_case_refs": [
            "rc-a1-agent-assisted-purchase-valid-controls"
          ],
          "memory_record_refs": [
            "ta-a1-agent-assisted-purchase-valid-controls",
            "rc-a1-agent-assisted-purchase-valid-controls"
          ],
          "detection_result_ref": "det-order-a1-agent-assisted-purchase-valid-controls",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-a2-fully-delegated-replenishment-purchase",
          "scenario_id": "commerce-a2-fully-delegated-replenishment-purchase",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-a2-fully-delegated-replenishment-purchase"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-a2-fully-delegated-replenishment-purchase"
          ],
          "replay_case_refs": [
            "rc-a2-fully-delegated-replenishment-purchase"
          ],
          "memory_record_refs": [
            "ta-a2-fully-delegated-replenishment-purchase",
            "rc-a2-fully-delegated-replenishment-purchase"
          ],
          "detection_result_ref": "det-order-a2-fully-delegated-replenishment-purchase",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-a3-agent-assisted-refund-approval-evidence",
          "scenario_id": "commerce-a3-agent-assisted-refund-approval-evidence",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-a3-agent-assisted-refund-approval-evidence"
          ],
          "refund_refs": [
            "refund-a3-agent-assisted-refund-approval-evidence"
          ],
          "trust_decision_refs": [
            "decision-a3-agent-assisted-refund-approval-evidence"
          ],
          "replay_case_refs": [
            "rc-a3-agent-assisted-refund-approval-evidence"
          ],
          "memory_record_refs": [
            "ta-a3-agent-assisted-refund-approval-evidence",
            "rc-a3-agent-assisted-refund-approval-evidence"
          ],
          "detection_result_ref": "det-order-a3-agent-assisted-refund-approval-evidence",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-h1-direct-human-purchase",
          "scenario_id": "commerce-h1-direct-human-purchase",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-h1-direct-human-purchase"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-h1-direct-human-purchase"
          ],
          "replay_case_refs": [
            "rc-h1-direct-human-purchase"
          ],
          "memory_record_refs": [
            "ta-h1-direct-human-purchase",
            "rc-h1-direct-human-purchase"
          ],
          "detection_result_ref": "det-order-h1-direct-human-purchase",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-h2-human-refund-valid-history",
          "scenario_id": "commerce-h2-human-refund-valid-history",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-h2-human-refund-valid-history"
          ],
          "refund_refs": [
            "refund-h2-human-refund-valid-history"
          ],
          "trust_decision_refs": [
            "decision-h2-human-refund-valid-history"
          ],
          "replay_case_refs": [
            "rc-h2-human-refund-valid-history"
          ],
          "memory_record_refs": [
            "ta-h2-human-refund-valid-history",
            "rc-h2-human-refund-valid-history"
          ],
          "detection_result_ref": "det-order-h2-human-refund-valid-history",
          "final_detection_status": "clean",
          "final_recommendation": "allow",
          "triggered_rule_ids": null,
          "promoted_to_replay": false,
          "expected_minimum_status": "clean",
          "passed": true,
          "notes": [
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-s1-refund-weak-authorization",
          "scenario_id": "commerce-s1-refund-weak-authorization",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-s1-refund-weak-authorization"
          ],
          "refund_refs": [
            "refund-s1-refund-weak-authorization"
          ],
          "trust_decision_refs": [
            "decision-s1-refund-weak-authorization"
          ],
          "replay_case_refs": [
            "rc-s1-refund-weak-authorization"
          ],
          "memory_record_refs": [
            "ta-s1-refund-weak-authorization",
            "rc-s1-refund-weak-authorization"
          ],
          "detection_result_ref": "det-order-s1-refund-weak-authorization",
          "final_detection_status": "step_up_required",
          "final_recommendation": "step_up",
          "triggered_rule_ids": [
            "agent_refund_without_approval",
            "missing_mandate_delegated_action",
            "missing_provenance_sensitive_action",
            "prior_step_up_decision",
            "refund_weak_authorization"
          ],
          "promoted_to_replay": false,
          "expected_minimum_status": "step_up_required",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-stable-commerce-s4-repeated-agent-refund-attempts",
          "scenario_id": "commerce-s4-repeated-agent-refund-attempts",
          "scenario_family": "commerce",
          "set_kind": "stable",
          "execution_status": "completed",
          "order_refs": [
            "order-s4-repeated-agent-refund-attempts"
          ],
          "refund_refs": [
            "refund-s4-repeated-agent-refund-attempts"
          ],
          "trust_decision_refs": [
            "decision-s4-repeated-agent-refund-attempts"
          ],
          "replay_case_refs": [
            "rc-s4-repeated-agent-refund-attempts"
          ],
          "memory_record_refs": [
            "ta-s4-repeated-agent-refund-attempts",
            "rc-s4-repeated-agent-refund-attempts"
          ],
          "detection_result_ref": "det-order-s4-repeated-agent-refund-attempts",
          "final_detection_status": "suspicious",
          "final_recommendation": "review",
          "triggered_rule_ids": [
            "prior_step_up_decision",
            "repeat_suspicious_context"
          ],
          "promoted_to_replay": false,
          "expected_minimum_status": "suspicious",
          "passed": true,
          "notes": [
            "tier_c_used",
            "stable baseline met its expected detector outcome"
          ]
        },
        {
          "id": "sr-living-commerce-v1-weakened-provenance",
          "scenario_id": "commerce-v1-weakened-provenance",
          "scenario_family": "commerce",
          "set_kind": "living",
          "challenger_variant_id": "variant-v1-weakened-provenance",
          "execution_status": "completed",
          "order_refs": [
            "order-v1-weakened-provenance"
          ],
          "refund_refs": null,
          "trust_decision_refs": [
            "decision-v1-weakened-provenance"
          ],
          "replay_case_refs": [
            "rc-v1-weakened-provenance"
          ],
          "memory_record_refs": [
            "ta-v1-weakened-provenance",
            "rc-v1-weakened-provenance"
          ],
          "detection_result_ref": "det-order-v1-weakened-provenance",
          "final_detection_status": "suspicious",
          "final_recommendation": "review",
          "triggered_rule_ids": [
            "missing_provenance_sensitive_action"
          ],
          "promoted_to_replay": false,
          "expected_minimum_status": "suspicious",
          "passed": true,
          "notes": [
            "tier_c_used",
            "challenger variant met the expected minimum detector posture"
          ]
        },
        {
          "id": "sr-living-commerce-v2-expired-inactive-mandate",
          "scenario_id": "commerce-v2-expired-inactive-mandate",
          "scenario_family": "commerce",
          "set_kind": "living",
          "challenger_variant_id": "variant-v2-expired-mandate",
          "execution_status": "completed",
          "order_refs"
```

### [PASS] GET /api/v1/operator/promotions

- Kind: `api`
- Summary: http 200; promotions=113; distinct_rounds=15
- URL: `http://127.0.0.1:8090/api/v1/operator/promotions`

```text
{
  "data": [
    {
      "round_id": "round-20260325235027",
      "promotion": {
        "id": "promo-round-20260325235027-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260325235027",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-25T23:50:28.315036Z"
      }
    },
    {
      "round_id": "round-20260325235027",
      "promotion": {
        "id": "promo-round-20260325235027-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260325235027",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-25T23:50:28.338378Z"
      }
    },
    {
      "round_id": "round-20260325235027",
      "promotion": {
        "id": "promo-round-20260325235027-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260325235027",
        "scenario_id": "commerce-v2-expired-inactive-mandate",
        "challenger_variant_id": "variant-v2-expired-mandate",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v2-expired-inactive-mandate",
        "replay_case_ref": "rc-v2-expired-inactive-mandate",
        "scenario_result_ref": "sr-replay_regression-commerce-v2-expired-inactive-mandate",
        "promoted": true,
        "created_at": "2026-03-25T23:50:28.361305Z"
      }
    },
    {
      "round_id": "round-20260325235027",
      "promotion": {
        "id": "promo-round-20260325235027-commerce-v3-approval-removed-regression",
        "round_id": "round-20260325235027",
        "scenario_id": "commerce-v3-approval-removed",
        "challenger_variant_id": "variant-v3-approval-removed",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v3-approval-removed",
        "replay_case_ref": "rc-v3-approval-removed",
        "scenario_result_ref": "sr-replay_regression-commerce-v3-approval-removed",
        "promoted": true,
        "created_at": "2026-03-25T23:50:28.384971Z"
      }
    },
    {
      "round_id": "round-20260325233731",
      "promotion": {
        "id": "promo-round-20260325233731-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260325233731",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-25T23:37:31.774437Z"
      }
    },
    {
      "round_id": "round-20260325233731",
      "promotion": {
        "id": "promo-round-20260325233731-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260325233731",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-25T23:37:31.796096Z"
      }
    },
    {
      "round_id": "round-20260325233731",
      "promotion": {
        "id": "promo-round-20260325233731-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260325233731",
        "scenario_id": "commerce-v2-expired-inactive-mandate",
        "challenger_variant_id": "variant-v2-expired-mandate",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v2-expired-inactive-mandate",
        "replay_case_ref": "rc-v2-expired-inactive-mandate",
        "scenario_result_ref": "sr-replay_regression-commerce-v2-expired-inactive-mandate",
        "promoted": true,
        "created_at": "2026-03-25T23:37:31.817602Z"
      }
    },
    {
      "round_id": "round-20260325233731",
      "promotion": {
        "id": "promo-round-20260325233731-commerce-v3-approval-removed-regression",
        "round_id": "round-20260325233731",
        "scenario_id": "commerce-v3-approval-removed",
        "challenger_variant_id": "variant-v3-approval-removed",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v3-approval-removed",
        "replay_case_ref": "rc-v3-approval-removed",
        "scenario_result_ref": "sr-replay_regression-commerce-v3-approval-removed",
        "promoted": true,
        "created_at": "2026-03-25T23:37:31.839747Z"
      }
    },
    {
      "round_id": "round-20260325230108",
      "promotion": {
        "id": "promo-round-20260325230108-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260325230108",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-25T23:01:08.953422Z"
      }
    },
    {
      "round_id": "round-20260325230108",
      "promotion": {
        "id": "promo-round-20260325230108-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260325230108",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-25T23:01:08.973746Z"
      }
    },
    {
      "round_id": "round-20260325230108",
      "promotion": {
        "id": "promo-round-20260325230108-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260325230108",
        "scenario_id": "commerce-v2-expired-inactive-mandate",
        "challenger_variant_id": "variant-v2-expired-mandate",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v2-expired-inactive-mandate",
        "replay_case_ref": "rc-v2-expired-inactive-mandate",
        "scenario_result_ref": "sr-replay_regression-commerce-v2-expired-inactive-mandate",
        "promoted": true,
        "created_at": "2026-03-25T23:01:08.993788Z"
      }
    },
    {
      "round_id": "round-20260325230108",
      "promotion": {
        "id": "promo-round-20260325230108-commerce-v3-approval-removed-regression",
        "round_id": "round-20260325230108",
        "scenario_id": "commerce-v3-approval-removed",
        "challenger_variant_id": "variant-v3-approval-removed",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v3-approval-removed",
        "replay_case_ref": "rc-v3-approval-removed",
        "scenario_result_ref": "sr-replay_regression-commerce-v3-approval-removed",
        "promoted": true,
        "created_at": "2026-03-25T23:01:09.014524Z"
      }
    },
    {
      "round_id": "round-20260325230022",
      "promotion": {
        "id": "promo-round-20260325230022-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260325230022",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-25T23:00:23.316481Z"
      }
    },
    {
      "round_id": "round-20260325230022",
      "promotion": {
        "id": "promo-round-20260325230022-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260325230022",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-25T23:00:23.335544Z"
      }
    },
    {
      "round_id": "round-20260325230022",
      "promotion": {
        "id": "promo-round-20260325230022-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260325230022",
        "scenario_id": "commerce-v2-expired-inactive-mandate
```

### [PASS] GET /api/v1/benchmark/recommendations

- Kind: `api`
- Summary: http 200; recommendations=59; types=add_to_replay_stable_set, investigate_repeat_refund_pattern, monitor_in_shadow_mode, require_provenance_for_delegated_purchase, require_step_up_for_delegated_refunds, tighten_refund_review_rule
- URL: `http://127.0.0.1:8090/api/v1/benchmark/recommendations`

```text
{
  "data": [
    {
      "id": "rec-round-20260325235027-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260325235027-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325235027-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260325235027-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260325235027-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "supporting_rule_ids": [
        "refund_weak_authorization"
      ],
      "suggested_action": "Tighten refund review thresholds for weak or missing authority signals and compare the sidecar recommendation rate with current manual-review outcomes.",
      "existing_control_integration_note": "Fits naturally beside existing refund review queues and PSP-side refund controls."
    },
    {
      "id": "rec-round-20260325235027-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v3-approval-removed",
        "commerce-v3-approval-removed",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325235027-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260325235027-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-s4-repeated-agent-refund-attempts",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "supporting_rule_ids": [
        "repeat_suspicious_context"
      ],
      "suggested_action": "Investigate repeat refund patterns, compare them with incumbent case outcomes, and tune queueing logic in shadow mode before any blocking change.",
      "existing_control_integration_note": "Best used as an investigative sidecar signal that feeds existing fraud-review workflows."
    },
    {
      "id": "rec-round-20260325235027-delegated-provenance",
      "type": "require_provenance_for_delegated_purchase",
      "rationale": "Delegated purchase paths with weak or missing provenance should not be treated as equivalent to ordinary human commerce, especially when they drift into new behavior patterns.",
      "priority": "moderate",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s1-refund-weak-authorization",
        "commerce-s2-delegated-purchase-weak-provenance",
        "commerce-v1-weakened-provenance",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-challenger-weakened-provenance-purchase-regression"
      ],
      "supporting_rule_ids": [
        "missing_provenance_sensitive_action"
      ],
      "suggested_action": "Require provenance for delegated purchases or keep them in recommendation-only shadow review until the team is comfortable tightening incumbent purchase controls.",
      "existing_control_integration_note": "Useful as a sidecar recommendation beside current purchase review thresholds and PSP checks."
    },
    {
      "id": "rec-round-20260325235027-shadow",
      "type": "monitor_in_shadow_mode",
      "rationale": "This round is best used as a recommendation-only sidecar beside the incumbent fraud stack so the team can compare benchmark findings against existing review and decision outcomes.",
      "priority": "low",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-a1-agent-assisted-purchase-valid-controls",
        "commerce-a2-fully-delegated-replenishment-purchase",
        "commerce-a3-agent-assisted-refund-approval-evidence",
        "commerce-h1-direct-human-purchase",
        "commerce-h2-human-refund-valid-history",
        "commerce-s1-refund-weak-authorization",
        "commerce-s2-delegated-purchase-weak-provenance",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-s4-repeated-agent-refund-attempts",
        "commerce-s5-merchant-scope-drift-delegated-action",
        "commerce-v1-weakened-provenance",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation",
        "commerce-v6-merchant-scope-drift",
        "commerce-v7-high-value-delegated-purchase"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260325235027-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325235027-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260325235027-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Keep the harness in shadow mode, compare its outputs with current queue and policy outcomes, and only trial control changes after replay confirms the improvement.",
      "existing_control_integration_note": "Designed to run beside existing fraud rules, queueing, and PSP controls without blocking live traffic."
    },
    {
      "id": "rec-round-20260325233731-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260325233731",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325233731-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260325233731-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325233731-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260325233731-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260325233731-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260325233731",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "supporting_rule_ids": [
        "refund_weak_authorization"
      ],
      "suggested_action": "Tighten refund review thresholds for weak or missing authority signals and compare the sidecar recommendation rate with current manual-review outcomes.",
      "existing_control_integration_note": "Fits naturally beside existing refund review queues and PSP-side refund controls."
    },
    {
      "id": "rec-round-20260325233731-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260325233731",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v3-approval-removed",
        "commerce-v3-approval-removed",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325233731-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325233731-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260325233731-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
     
```

### [PASS] GET /api/v1/operator/recommendations

- Kind: `api`
- Summary: http 200; recommendations=59; types=add_to_replay_stable_set, investigate_repeat_refund_pattern, monitor_in_shadow_mode, require_provenance_for_delegated_purchase, require_step_up_for_delegated_refunds, tighten_refund_review_rule
- URL: `http://127.0.0.1:8090/api/v1/operator/recommendations`

```text
{
  "data": [
    {
      "id": "rec-round-20260325235027-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260325235027-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325235027-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260325235027-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260325235027-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "supporting_rule_ids": [
        "refund_weak_authorization"
      ],
      "suggested_action": "Tighten refund review thresholds for weak or missing authority signals and compare the sidecar recommendation rate with current manual-review outcomes.",
      "existing_control_integration_note": "Fits naturally beside existing refund review queues and PSP-side refund controls."
    },
    {
      "id": "rec-round-20260325235027-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v3-approval-removed",
        "commerce-v3-approval-removed",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325235027-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260325235027-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-s4-repeated-agent-refund-attempts",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "supporting_rule_ids": [
        "repeat_suspicious_context"
      ],
      "suggested_action": "Investigate repeat refund patterns, compare them with incumbent case outcomes, and tune queueing logic in shadow mode before any blocking change.",
      "existing_control_integration_note": "Best used as an investigative sidecar signal that feeds existing fraud-review workflows."
    },
    {
      "id": "rec-round-20260325235027-delegated-provenance",
      "type": "require_provenance_for_delegated_purchase",
      "rationale": "Delegated purchase paths with weak or missing provenance should not be treated as equivalent to ordinary human commerce, especially when they drift into new behavior patterns.",
      "priority": "moderate",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s1-refund-weak-authorization",
        "commerce-s2-delegated-purchase-weak-provenance",
        "commerce-v1-weakened-provenance",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-challenger-weakened-provenance-purchase-regression"
      ],
      "supporting_rule_ids": [
        "missing_provenance_sensitive_action"
      ],
      "suggested_action": "Require provenance for delegated purchases or keep them in recommendation-only shadow review until the team is comfortable tightening incumbent purchase controls.",
      "existing_control_integration_note": "Useful as a sidecar recommendation beside current purchase review thresholds and PSP checks."
    },
    {
      "id": "rec-round-20260325235027-shadow",
      "type": "monitor_in_shadow_mode",
      "rationale": "This round is best used as a recommendation-only sidecar beside the incumbent fraud stack so the team can compare benchmark findings against existing review and decision outcomes.",
      "priority": "low",
      "linked_round_id": "round-20260325235027",
      "linked_scenario_ids": [
        "commerce-a1-agent-assisted-purchase-valid-controls",
        "commerce-a2-fully-delegated-replenishment-purchase",
        "commerce-a3-agent-assisted-refund-approval-evidence",
        "commerce-h1-direct-human-purchase",
        "commerce-h2-human-refund-valid-history",
        "commerce-s1-refund-weak-authorization",
        "commerce-s2-delegated-purchase-weak-provenance",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-s4-repeated-agent-refund-attempts",
        "commerce-s5-merchant-scope-drift-delegated-action",
        "commerce-v1-weakened-provenance",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation",
        "commerce-v6-merchant-scope-drift",
        "commerce-v7-high-value-delegated-purchase"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325235027-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260325235027-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325235027-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260325235027-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Keep the harness in shadow mode, compare its outputs with current queue and policy outcomes, and only trial control changes after replay confirms the improvement.",
      "existing_control_integration_note": "Designed to run beside existing fraud rules, queueing, and PSP controls without blocking live traffic."
    },
    {
      "id": "rec-round-20260325233731-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260325233731",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325233731-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260325233731-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325233731-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260325233731-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260325233731-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260325233731",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "supporting_rule_ids": [
        "refund_weak_authorization"
      ],
      "suggested_action": "Tighten refund review thresholds for weak or missing authority signals and compare the sidecar recommendation rate with current manual-review outcomes.",
      "existing_control_integration_note": "Fits naturally beside existing refund review queues and PSP-side refund controls."
    },
    {
      "id": "rec-round-20260325233731-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260325233731",
      "linked_scenario_ids": [
        "commerce-s1-refund-weak-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v3-approval-removed",
        "commerce-v3-approval-removed",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260325233731-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260325233731-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260325233731-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
     
```

### [PASS] GET /api/v1/benchmark/trends/summary

- Kind: `api`
- Summary: http 200; present_keys=4; missing=recommendation_counts
- URL: `http://127.0.0.1:8090/api/v1/benchmark/trends/summary`

```text
{
  "data": {
    "rounds_executed": 15,
    "promotions_over_time": [
      {
        "round_id": "round-20260325235027",
        "value": 4
      },
      {
        "round_id": "round-20260325233731",
        "value": 4
      },
      {
        "round_id": "round-20260325230108",
        "value": 4
      },
      {
        "round_id": "round-20260325230022",
        "value": 4
      },
      {
        "round_id": "round-20260325220029",
        "value": 22
      },
      {
        "round_id": "round-20260325220013",
        "value": 19
      },
      {
        "round_id": "round-20260325215958",
        "value": 16
      },
      {
        "round_id": "round-20260325215157",
        "value": 13
      },
      {
        "round_id": "round-20260325215142",
        "value": 10
      },
      {
        "round_id": "round-20260325215127",
        "value": 7
      },
      {
        "round_id": "round-20260325214203",
        "value": 4
      },
      {
        "round_id": "round-20260325180315",
        "value": 2
      },
      {
        "round_id": "round-20260325172529",
        "value": 1
      },
      {
        "round_id": "round-20260325143447",
        "value": 2
      },
      {
        "round_id": "round-20260325141726",
        "value": 1
      }
    ],
    "replay_pass_rate_over_time": [
      {
        "round_id": "round-20260325235027",
        "value": 0
      },
      {
        "round_id": "round-20260325233731",
        "value": 0
      },
      {
        "round_id": "round-20260325230108",
        "value": 0
      },
      {
        "round_id": "round-20260325230022",
        "value": 0
      },
      {
        "round_id": "round-20260325220029",
        "value": 0
      },
      {
        "round_id": "round-20260325220013",
        "value": 0
      },
      {
        "round_id": "round-20260325215958",
        "value": 0
      },
      {
        "round_id": "round-20260325215157",
        "value": 0
      },
      {
        "round_id": "round-20260325215142",
        "value": 0
      },
      {
        "round_id": "round-20260325215127",
        "value": 0
      },
      {
        "round_id": "round-20260325214203",
        "value": 0.5
      },
      {
        "round_id": "round-20260325180315",
        "value": 0
      },
      {
        "round_id": "round-20260325172529",
        "value": 1
      },
      {
        "round_id": "round-20260325143447",
        "value": 0
      },
      {
        "round_id": "round-20260325141726",
        "value": 1
      }
    ],
    "new_blind_spots_discovered": 15,
    "regressions_observed": 0,
    "recommendation_counts_by_type": {
      "add_to_replay_stable_set": 11,
      "investigate_repeat_refund_pattern": 4,
      "monitor_in_shadow_mode": 11,
      "require_provenance_for_delegated_purchase": 11,
      "require_step_up_for_delegated_refunds": 11,
      "tighten_refund_review_rule": 11
    },
    "top_recurring_evasion_patterns": [
      "No previously promoted replay cases were available for regression retest.",
      "Replay regression pass rate fell to 0.00.",
      "commerce-challenger-weakened-provenance-purchase promoted because Previously promoted replay case regressed below its expected detection floor..",
      "commerce-challenger-weakened-provenance-purchase promoted because Suspicious challenger behavior evaluated as clean..",
      "commerce-s3-approval-removed-after-authorization promoted because Challenger behavior scored below its expected minimum detector posture..",
      "commerce-s3-approval-removed-after-authorization promoted because Previously promoted replay case regressed below its expected detection floor..",
      "commerce-v2-expired-inactive-mandate promoted because Challenger behavior scored below its expected minimum detector posture..",
      "commerce-v2-expired-inactive-mandate promoted because Previously promoted replay case regressed below its expected detection floor..",
      "commerce-v3-approval-removed promoted because Challenger behavior scored below its expected minimum detector posture..",
      "commerce-v3-approval-removed promoted because Previously promoted replay case regressed below its expected detection floor.."
    ]
  }
}
```

### [PASS] GET /api/v1/operator/trends/summary

- Kind: `api`
- Summary: http 200; present_keys=4; missing=recommendation_counts
- URL: `http://127.0.0.1:8090/api/v1/operator/trends/summary`

```text
{
  "data": {
    "rounds_executed": 15,
    "promotions_over_time": [
      {
        "round_id": "round-20260325235027",
        "value": 4
      },
      {
        "round_id": "round-20260325233731",
        "value": 4
      },
      {
        "round_id": "round-20260325230108",
        "value": 4
      },
      {
        "round_id": "round-20260325230022",
        "value": 4
      },
      {
        "round_id": "round-20260325220029",
        "value": 22
      },
      {
        "round_id": "round-20260325220013",
        "value": 19
      },
      {
        "round_id": "round-20260325215958",
        "value": 16
      },
      {
        "round_id": "round-20260325215157",
        "value": 13
      },
      {
        "round_id": "round-20260325215142",
        "value": 10
      },
      {
        "round_id": "round-20260325215127",
        "value": 7
      },
      {
        "round_id": "round-20260325214203",
        "value": 4
      },
      {
        "round_id": "round-20260325180315",
        "value": 2
      },
      {
        "round_id": "round-20260325172529",
        "value": 1
      },
      {
        "round_id": "round-20260325143447",
        "value": 2
      },
      {
        "round_id": "round-20260325141726",
        "value": 1
      }
    ],
    "replay_pass_rate_over_time": [
      {
        "round_id": "round-20260325235027",
        "value": 0
      },
      {
        "round_id": "round-20260325233731",
        "value": 0
      },
      {
        "round_id": "round-20260325230108",
        "value": 0
      },
      {
        "round_id": "round-20260325230022",
        "value": 0
      },
      {
        "round_id": "round-20260325220029",
        "value": 0
      },
      {
        "round_id": "round-20260325220013",
        "value": 0
      },
      {
        "round_id": "round-20260325215958",
        "value": 0
      },
      {
        "round_id": "round-20260325215157",
        "value": 0
      },
      {
        "round_id": "round-20260325215142",
        "value": 0
      },
      {
        "round_id": "round-20260325215127",
        "value": 0
      },
      {
        "round_id": "round-20260325214203",
        "value": 0.5
      },
      {
        "round_id": "round-20260325180315",
        "value": 0
      },
      {
        "round_id": "round-20260325172529",
        "value": 1
      },
      {
        "round_id": "round-20260325143447",
        "value": 0
      },
      {
        "round_id": "round-20260325141726",
        "value": 1
      }
    ],
    "new_blind_spots_discovered": 15,
    "regressions_observed": 0,
    "recommendation_counts_by_type": {
      "add_to_replay_stable_set": 11,
      "investigate_repeat_refund_pattern": 4,
      "monitor_in_shadow_mode": 11,
      "require_provenance_for_delegated_purchase": 11,
      "require_step_up_for_delegated_refunds": 11,
      "tighten_refund_review_rule": 11
    },
    "top_recurring_evasion_patterns": [
      "No previously promoted replay cases were available for regression retest.",
      "Replay regression pass rate fell to 0.00.",
      "commerce-challenger-weakened-provenance-purchase promoted because Previously promoted replay case regressed below its expected detection floor..",
      "commerce-challenger-weakened-provenance-purchase promoted because Suspicious challenger behavior evaluated as clean..",
      "commerce-s3-approval-removed-after-authorization promoted because Challenger behavior scored below its expected minimum detector posture..",
      "commerce-s3-approval-removed-after-authorization promoted because Previously promoted replay case regressed below its expected detection floor..",
      "commerce-v2-expired-inactive-mandate promoted because Challenger behavior scored below its expected minimum detector posture..",
      "commerce-v2-expired-inactive-mandate promoted because Previously promoted replay case regressed below its expected detection floor..",
      "commerce-v3-approval-removed promoted because Challenger behavior scored below its expected minimum detector posture..",
      "commerce-v3-approval-removed promoted because Previously promoted replay case regressed below its expected detection floor.."
    ]
  }
}
```

### [PASS] GET /api/v1/benchmark/scheduler/status

- Kind: `api`
- Summary: http 200; scheduler_mode=unknown; interval=24h0m0s
- URL: `http://127.0.0.1:8090/api/v1/benchmark/scheduler/status`

```text
{
  "data": {
    "enabled": false,
    "running": false,
    "scenario_family": "commerce",
    "interval": "24h0m0s",
    "max_runs": 7,
    "executed_runs": 4,
    "dry_run": false,
    "last_round_id": "round-20260325235027",
    "last_started_at": "2026-03-25T23:50:27.907712Z",
    "next_run_at": "0001-01-01T00:00:00Z"
  }
}
```
