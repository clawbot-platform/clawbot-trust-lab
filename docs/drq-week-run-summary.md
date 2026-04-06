# DRQ Week-Run Summary

## Run segmentation
- Root used: `/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab/docs/drq-week-run-reports/reports`
- Phase B starts at round: `round-20260401070127`

## Parse diagnostics
- Round directories discovered: 27
- Parsed rounds: 27
- Phase A rounds classified: 13
- Phase B rounds classified: 14
- Skipped round directories: 0

## Headline comparison

| Metric                 |                Phase A |                Phase B |
|------------------------|-----------------------:|-----------------------:|
| Rounds                 |                     13 |                     14 |
| First round            | `round-20260328202225` | `round-20260401070127` |
| Last round             | `round-20260331203605` | `round-20260404130127` |
| Total promotions       |                     39 |                      0 |
| Avg promotions / round |                   3.00 |                   0.00 |
| Avg replay pass rate   |                   0.08 |                   1.00 |
| Zero-promotion rounds  |                      0 |                     14 |
| Perfect replay rounds  |                      1 |                     14 |

## Detector versions seen
- Phase A: v1.0.0-9-g45d296f
- Phase B: v1.0.0-10-g796421c

## Robustness outcomes
- Phase A: {'new_blind_spot_discovered': 13}
- Phase B: {'improved': 14}

## Targeted weak-case status in latest round of each phase

### commerce-v2-expired-inactive-mandate
- Phase A latest: status=`suspicious`, passed=`False`, rules=`['missing_mandate_delegated_action', 'prior_step_up_decision']`
- Phase B latest: status=`step_up_required`, passed=`True`, rules=`['expired_inactive_mandate', 'missing_mandate_delegated_action', 'prior_step_up_decision']`

### commerce-v3-approval-removed
- Phase A latest: status=`suspicious`, passed=`False`, rules=`['agent_refund_without_approval', 'prior_step_up_decision']`
- Phase B latest: status=`step_up_required`, passed=`True`, rules=`['agent_refund_without_approval', 'approval_removed_after_authorization', 'prior_step_up_decision']`

### commerce-s3-approval-removed-after-authorization
- Phase A latest: status=`suspicious`, passed=`False`, rules=`['agent_refund_without_approval', 'prior_step_up_decision']`
- Phase B latest: status=`step_up_required`, passed=`True`, rules=`['agent_refund_without_approval', 'approval_removed_after_authorization', 'prior_step_up_decision']`

## Recommendation type totals
- Phase A: {'monitor_in_shadow_mode': 13, 'add_to_replay_stable_set': 13, 'tighten_refund_review_rule': 13, 'require_step_up_for_delegated_refunds': 13, 'require_provenance_for_delegated_purchase': 13, 'investigate_repeat_refund_pattern': 13}
- Phase B: {'monitor_in_shadow_mode': 14, 'tighten_refund_review_rule': 14, 'require_step_up_for_delegated_refunds': 14, 'require_provenance_for_delegated_purchase': 14, 'investigate_repeat_refund_pattern': 14}

## Conclusion
- Phase B improved detector performance versus Phase A.
- Promotions fell materially in the tuned window.
- Replay pass rate improved in the tuned window.
