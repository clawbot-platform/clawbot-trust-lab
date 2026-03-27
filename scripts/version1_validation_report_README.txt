Version 1 validation runner

File:
- version1_validation_report.py

What it does:
- validates repo docs and release-surface files in developer mode
- optionally runs backend and web quality checks in developer mode
- checks repo-native Docker Compose state for the Version 1 core stack
- calls live trust-lab APIs for health, rounds, promotions, recommendations, trends, and scheduler status
- optionally triggers a fresh benchmark round
- writes both Markdown and HTML validation reports

Modes:
- `--mode developer`
  - full repo-quality validation plus runtime checks
  - intended for CI and development workstations
- `--mode runtime`
  - deployed-system checks only
  - intended for OptiPlex-style runtime hosts and dry-run appliances
  - does not require local Go, npm, lint, or security tooling

Current Docker defaults:
- core compose: `deploy/compose/docker-compose.yml`
- local override: `deploy/compose/docker-compose.override.yml`
- env file: `.env`

Example usage:

1. Validate the running core Docker stack in runtime mode:

   make validate-v1-runtime

2. Also trigger a fresh benchmark round in runtime mode:

   python3 ./scripts/version1_validation_report.py \
     --mode runtime \
     --deployment-mode docker \
     --compose-file deploy/compose/docker-compose.yml \
     --compose-override-file deploy/compose/docker-compose.override.yml \
     --compose-env-file .env \
     --run-round \
     --output-dir ./version1-validation-output

3. Run the full developer-mode validation:

   python3 ./scripts/version1_validation_report.py \
     --mode developer \
     --deployment-mode docker \
     --compose-file deploy/compose/docker-compose.yml \
     --compose-override-file deploy/compose/docker-compose.override.yml \
     --compose-env-file .env \
     --run-round \
     --output-dir ./version1-validation-output

Outputs:
- version1-validation-output/version1-validation-report.md
- version1-validation-output/version1-validation-report.html

Exit code:
- `0` if all checks passed
- `1` if any checks failed
