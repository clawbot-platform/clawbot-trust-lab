Phase 9 validation runner

Files:
- phase9_validation_report.py
- phase9_validation-report.md / .html (generated when run)

What it does:
- runs backend static checks
- runs web checks (lint/test/build/e2e) if web/ exists
- checks docs for Phase 9 API/persistence coverage
- checks report artifact directories for expected files
- calls live API endpoints for scenarios, rounds, promotions, recommendations, trends, and scheduler status
- optionally triggers a fresh benchmark round
- writes both Markdown and HTML reports

Example usage:
1) From the clawbot-trust-lab repo root, with services already running:

   python3 /mnt/data/phase9_validation_report.py \
     --repo-root . \
     --api-base http://127.0.0.1:8090 \
     --output-dir ./phase9-validation-output

2) To also trigger a fresh benchmark round during validation:

   python3 /mnt/data/phase9_validation_report.py \
     --repo-root . \
     --api-base http://127.0.0.1:8090 \
     --run-round \
     --output-dir ./phase9-validation-output

3) To skip slower sections:

   python3 /mnt/data/phase9_validation_report.py \
     --repo-root . \
     --skip-web \
     --output-dir ./phase9-validation-output

Outputs:
- phase9-validation-output/phase9-validation-report.md
- phase9-validation-output/phase9-validation-report.html

Exit code:
- 0 if all checks passed
- 1 if any checks failed
