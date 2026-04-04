# Executive Final Report
## DRQ Week-Run Assessment — Commerce Scenario Family

**Run window:** March 28, 2026 to April 4, 2026  
**Environment:** Homelab, image-based distributed deployment  
**Execution model:** Docker-only runtime hosts with public GHCR images  
**Phases:**
- **Phase A:** Baseline detector (`v1.0.0-9`)
- **Phase B:** Tuned detector (`v1.0.0-10`)

## Executive Summary

The DRQ week run was successful.

The effort demonstrated that the tuned Trust Lab detector materially outperformed the baseline detector in a stable, distributed homelab deployment. Across the tuned run window, promotions dropped from a repeated baseline pattern of **3 promotions per round** to **0 promotions per round**, while replay pass rate improved from repeated **0** during the baseline window to sustained **1.0** across the tuned window. The three targeted weak cases that drove the mid-run tuning effort were successfully corrected and met their intended minimum posture in the tuned rounds. :contentReference[oaicite:0]{index=0} :contentReference[oaicite:1]{index=1}

Operationally, the platform remained healthy throughout the run. The control plane, memory service, and Trust Lab service all stayed healthy, and the run completed using the intended image-based Docker deployment model without host-side Go or application toolchains on the runtime nodes. :contentReference[oaicite:2]{index=2}

## Key Outcome

The tuned detector should be treated as the preferred benchmark candidate configuration for future DRQ work.

## Run Design

The week run was intentionally split into two phases:

### Phase A — Baseline
Phase A established the baseline benchmark behavior using the pre-tuning detector image. During this phase, the same three weak cases repeatedly surfaced as promotions, and replay pass performance stayed poor across the scheduled rounds. :contentReference[oaicite:3]{index=3}

### Phase B — Tuned
Phase B used a narrowly tuned detector image that addressed only the known weak cases from Phase A. The tuned window showed a clean and sustained improvement without introducing operational instability. :contentReference[oaicite:4]{index=4} :contentReference[oaicite:5]{index=5}

## Headline Comparison

### Phase A — Baseline (`v1.0.0-9`)
- Rounds: **13**
- Total promotions: **39**
- Average promotions per round: **3.00**
- Replay pass rate pattern: mostly **0**
- Detector posture on targeted weak cases: below intended floor, causing promotions :contentReference[oaicite:6]{index=6}

### Phase B — Tuned (`v1.0.0-10`)
- Rounds: **14**
- Total promotions: **0**
- Average promotions per round: **0.00**
- Replay pass rate pattern: sustained **1.00**
- Detector posture on targeted weak cases: improved to intended minimum floor :contentReference[oaicite:7]{index=7} :contentReference[oaicite:8]{index=8}

## Targeted Improvements Confirmed

The Phase B tuning corrected the following three recurring weak cases:

1. **`commerce-v2-expired-inactive-mandate`**  
   Final posture in tuned round: `step_up_required`, `passed: true`, with `expired_inactive_mandate` present in triggered rules. :contentReference[oaicite:9]{index=9}

2. **`commerce-v3-approval-removed`**  
   Final posture in tuned round: `step_up_required`, `passed: true`, with `approval_removed_after_authorization` present in triggered rules. :contentReference[oaicite:10]{index=10}

3. **`commerce-s3-approval-removed-after-authorization`**  
   Final posture in tuned round: `step_up_required`, `passed: true`, with `approval_removed_after_authorization` present in triggered rules. :contentReference[oaicite:11]{index=11}

## Operational Assessment

The run also validated the deployment and operating model:

- distributed service placement across homelab nodes
- image-only runtime deployment from GHCR
- stable health across control plane, memory service, and Trust Lab
- successful report generation and retrieval throughout the run
- successful Phase A to Phase B cutover without changing deployment topology :contentReference[oaicite:12]{index=12}

This is important because it proves the DRQ workflow can be operated as a practical benchmark harness in a lightweight homelab environment, not just as a local development workflow.

## Final Recommendation

Adopt the tuned detector image as the preferred Trust Lab benchmark candidate state.

For future work:
- preserve the tuned detector behavior as the new replay reference point
- use the Phase B results as the benchmark baseline for the next iteration
- build future versions as explicit, measured improvements against this tuned state, not against the original baseline

## Final Verdict

The week run achieved its objective.

The effort produced a clear and defensible before/after result:
- baseline repeatedly exposed the same three weak cases
- tuned detector eliminated those repeated promotions
- replay performance improved materially
- the system remained operationally stable

This is a successful DRQ benchmark and tuning cycle. 
