# DRQ Week-Run Final Assessment

## Outcome
The tuned Phase B detector outperformed the Phase A baseline and should be considered the stronger candidate configuration.

## Phase A (baseline, v1.0.0-9)
- Rounds: 13
- Total promotions: 39
- Avg promotions/round: 3.00
- Avg replay pass rate: 0.08

## Phase B (tuned, v1.0.0-10)
- Rounds: 14
- Total promotions: 0
- Avg promotions/round: 0.00
- Avg replay pass rate: 1.00

## Targeted improvements confirmed
The following recurring weak cases were corrected in the tuned detector window:
- commerce-v2-expired-inactive-mandate
- commerce-v3-approval-removed
- commerce-s3-approval-removed-after-authorization

Each reached step_up_required and passed in the latest tuned round.

## Final recommendation
Adopt the tuned detector image as the preferred Trust Lab benchmark candidate and use Phase B as the reference state for future replay and tuning work.