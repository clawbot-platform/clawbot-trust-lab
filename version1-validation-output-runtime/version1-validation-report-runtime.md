# Clawbot Trust Lab Version 1 Validation Report

Generated: 2026-03-27T10:33:23-04:00

Repo root: `/Users/piyushdaiya/Documents/projects/clawbot-platform/clawbot-trust-lab`

Validation mode: `runtime`

Deployment mode: `local`

Executed groups:
- live runtime API checks

Skipped groups:
- repo docs and release-surface checks (runtime mode)
- Go developer-tooling checks (runtime mode)
- web developer-tooling checks (runtime mode)
- local report artifact directory checks (runtime mode uses API-visible report checks instead)

Summary: **12 passed / 1 failed / 13 total**

## Checks

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

### [PASS] GET /version

- Kind: `api`
- Summary: http 200; version=v1.0.0-6-g5e6bfd0-dirty; commit=5e6bfd07d9e4; build_date=2026-03-27T14:22:17Z
- URL: `http://127.0.0.1:8090/version`

```text
{
  "version": "v1.0.0-6-g5e6bfd0-dirty",
  "commit": "5e6bfd07d9e4",
  "build_date": "2026-03-27T14:22:17Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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
      "created_at": "2026-03-27T14:22:18.141109Z"
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

### [PASS] GET /api/v1/benchmark/rounds

- Kind: `api`
- Summary: http 200; rounds=22; rounds_with_prod_bridge=22; rounds_with_tier_c_usage=18; recommendation_totals=117
- URL: `http://127.0.0.1:8090/api/v1/benchmark/rounds`

```text
{
  "data": [
    {
      "id": "round-20260327001551",
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
      "started_at": "2026-03-27T00:15:51.579522Z",
      "completed_at": "2026-03-27T00:15:52.135474Z",
      "round_status": "completed",
      "report_dir": "reports/round-20260327001551",
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
- Summary: http 200; rounds=22; rounds_with_prod_bridge=22; rounds_with_tier_c_usage=18; recommendation_totals=117
- URL: `http://127.0.0.1:8090/api/v1/operator/rounds`

```text
{
  "data": [
    {
      "id": "round-20260327001551",
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
      "started_at": "2026-03-27T00:15:51.579522Z",
      "completed_at": "2026-03-27T00:15:52.135474Z",
      "round_status": "completed",
      "report_dir": "reports/round-20260327001551",
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
- Summary: http 200; promotions=141; distinct_rounds=22
- URL: `http://127.0.0.1:8090/api/v1/operator/promotions`

```text
{
  "data": [
    {
      "round_id": "round-20260327001551",
      "promotion": {
        "id": "promo-round-20260327001551-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260327001551",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-27T00:15:52.066139Z"
      }
    },
    {
      "round_id": "round-20260327001551",
      "promotion": {
        "id": "promo-round-20260327001551-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260327001551",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-27T00:15:52.089424Z"
      }
    },
    {
      "round_id": "round-20260327001551",
      "promotion": {
        "id": "promo-round-20260327001551-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260327001551",
        "scenario_id": "commerce-v2-expired-inactive-mandate",
        "challenger_variant_id": "variant-v2-expired-mandate",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v2-expired-inactive-mandate",
        "replay_case_ref": "rc-v2-expired-inactive-mandate",
        "scenario_result_ref": "sr-replay_regression-commerce-v2-expired-inactive-mandate",
        "promoted": true,
        "created_at": "2026-03-27T00:15:52.113008Z"
      }
    },
    {
      "round_id": "round-20260327001551",
      "promotion": {
        "id": "promo-round-20260327001551-commerce-v3-approval-removed-regression",
        "round_id": "round-20260327001551",
        "scenario_id": "commerce-v3-approval-removed",
        "challenger_variant_id": "variant-v3-approval-removed",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v3-approval-removed",
        "replay_case_ref": "rc-v3-approval-removed",
        "scenario_result_ref": "sr-replay_regression-commerce-v3-approval-removed",
        "promoted": true,
        "created_at": "2026-03-27T00:15:52.135263Z"
      }
    },
    {
      "round_id": "round-20260326235545",
      "promotion": {
        "id": "promo-round-20260326235545-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260326235545",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-26T23:55:46.019675Z"
      }
    },
    {
      "round_id": "round-20260326235545",
      "promotion": {
        "id": "promo-round-20260326235545-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260326235545",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-26T23:55:46.043533Z"
      }
    },
    {
      "round_id": "round-20260326235545",
      "promotion": {
        "id": "promo-round-20260326235545-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260326235545",
        "scenario_id": "commerce-v2-expired-inactive-mandate",
        "challenger_variant_id": "variant-v2-expired-mandate",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v2-expired-inactive-mandate",
        "replay_case_ref": "rc-v2-expired-inactive-mandate",
        "scenario_result_ref": "sr-replay_regression-commerce-v2-expired-inactive-mandate",
        "promoted": true,
        "created_at": "2026-03-26T23:55:46.066293Z"
      }
    },
    {
      "round_id": "round-20260326235545",
      "promotion": {
        "id": "promo-round-20260326235545-commerce-v3-approval-removed-regression",
        "round_id": "round-20260326235545",
        "scenario_id": "commerce-v3-approval-removed",
        "challenger_variant_id": "variant-v3-approval-removed",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v3-approval-removed",
        "replay_case_ref": "rc-v3-approval-removed",
        "scenario_result_ref": "sr-replay_regression-commerce-v3-approval-removed",
        "promoted": true,
        "created_at": "2026-03-26T23:55:46.090089Z"
      }
    },
    {
      "round_id": "round-20260326235501",
      "promotion": {
        "id": "promo-round-20260326235501-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260326235501",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-26T23:55:02.278603Z"
      }
    },
    {
      "round_id": "round-20260326235501",
      "promotion": {
        "id": "promo-round-20260326235501-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260326235501",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-26T23:55:02.299538Z"
      }
    },
    {
      "round_id": "round-20260326235501",
      "promotion": {
        "id": "promo-round-20260326235501-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260326235501",
        "scenario_id": "commerce-v2-expired-inactive-mandate",
        "challenger_variant_id": "variant-v2-expired-mandate",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v2-expired-inactive-mandate",
        "replay_case_ref": "rc-v2-expired-inactive-mandate",
        "scenario_result_ref": "sr-replay_regression-commerce-v2-expired-inactive-mandate",
        "promoted": true,
        "created_at": "2026-03-26T23:55:02.320203Z"
      }
    },
    {
      "round_id": "round-20260326235501",
      "promotion": {
        "id": "promo-round-20260326235501-commerce-v3-approval-removed-regression",
        "round_id": "round-20260326235501",
        "scenario_id": "commerce-v3-approval-removed",
        "challenger_variant_id": "variant-v3-approval-removed",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-v3-approval-removed",
        "replay_case_ref": "rc-v3-approval-removed",
        "scenario_result_ref": "sr-replay_regression-commerce-v3-approval-removed",
        "promoted": true,
        "created_at": "2026-03-26T23:55:02.341605Z"
      }
    },
    {
      "round_id": "round-20260326225239",
      "promotion": {
        "id": "promo-round-20260326225239-commerce-challenger-weakened-provenance-purchase-regression",
        "round_id": "round-20260326225239",
        "scenario_id": "commerce-challenger-weakened-provenance-purchase",
        "challenger_variant_id": "variant-weakened-provenance",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-challenger-weakened-provenance-purchase",
        "replay_case_ref": "rc-challenger-weakened-provenance-purchase",
        "scenario_result_ref": "sr-replay_regression-commerce-challenger-weakened-provenance-purchase",
        "promoted": true,
        "created_at": "2026-03-26T22:52:39.67427Z"
      }
    },
    {
      "round_id": "round-20260326225239",
      "promotion": {
        "id": "promo-round-20260326225239-commerce-s3-approval-removed-after-authorization-regression",
        "round_id": "round-20260326225239",
        "scenario_id": "commerce-s3-approval-removed-after-authorization",
        "challenger_variant_id": "variant-s3-approval-removed-after-authorization",
        "promotion_reason": "meaningful_regression",
        "rationale": "Previously promoted replay case regressed below its expected detection floor.",
        "detection_result_ref": "det-order-s3-approval-removed-after-authorization",
        "replay_case_ref": "rc-s3-approval-removed-after-authorization",
        "scenario_result_ref": "sr-replay_regression-commerce-s3-approval-removed-after-authorization",
        "promoted": true,
        "created_at": "2026-03-26T22:52:39.698772Z"
      }
    },
    {
      "round_id": "round-20260326225239",
      "promotion": {
        "id": "promo-round-20260326225239-commerce-v2-expired-inactive-mandate-regression",
        "round_id": "round-20260326225239",
        "scenario_id": "commerce-v2-expired-inactive-mandate"
```

### [PASS] GET /api/v1/benchmark/recommendations

- Kind: `api`
- Summary: http 200; recommendations=117; types=add_to_replay_stable_set, investigate_repeat_refund_pattern, monitor_in_shadow_mode, require_provenance_for_delegated_purchase, require_step_up_for_delegated_refunds, tighten_refund_review_rule
- URL: `http://127.0.0.1:8090/api/v1/benchmark/recommendations`

```text
{
  "data": [
    {
      "id": "rec-round-20260327001551-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260327001551",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260327001551-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260327001551-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260327001551-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260327001551-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260327001551-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260327001551",
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
      "id": "rec-round-20260327001551-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260327001551",
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
        "promo-round-20260327001551-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260327001551-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260327001551-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
      "linked_round_id": "round-20260327001551",
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
      "id": "rec-round-20260327001551-delegated-provenance",
      "type": "require_provenance_for_delegated_purchase",
      "rationale": "Delegated purchase paths with weak or missing provenance should not be treated as equivalent to ordinary human commerce, especially when they drift into new behavior patterns.",
      "priority": "moderate",
      "linked_round_id": "round-20260327001551",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s1-refund-weak-authorization",
        "commerce-s2-delegated-purchase-weak-provenance",
        "commerce-v1-weakened-provenance",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260327001551-commerce-challenger-weakened-provenance-purchase-regression"
      ],
      "supporting_rule_ids": [
        "missing_provenance_sensitive_action"
      ],
      "suggested_action": "Require provenance for delegated purchases or keep them in recommendation-only shadow review until the team is comfortable tightening incumbent purchase controls.",
      "existing_control_integration_note": "Useful as a sidecar recommendation beside current purchase review thresholds and PSP checks."
    },
    {
      "id": "rec-round-20260327001551-shadow",
      "type": "monitor_in_shadow_mode",
      "rationale": "This round is best used as a recommendation-only sidecar beside the incumbent fraud stack so the team can compare benchmark findings against existing review and decision outcomes.",
      "priority": "low",
      "linked_round_id": "round-20260327001551",
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
        "promo-round-20260327001551-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260327001551-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260327001551-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260327001551-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Keep the harness in shadow mode, compare its outputs with current queue and policy outcomes, and only trial control changes after replay confirms the improvement.",
      "existing_control_integration_note": "Designed to run beside existing fraud rules, queueing, and PSP controls without blocking live traffic."
    },
    {
      "id": "rec-round-20260326235545-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260326235545",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260326235545-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260326235545-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260326235545-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260326235545-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260326235545-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260326235545",
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
      "id": "rec-round-20260326235545-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260326235545",
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
        "promo-round-20260326235545-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260326235545-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260326235545-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
     
```

### [PASS] GET /api/v1/operator/recommendations

- Kind: `api`
- Summary: http 200; recommendations=117; types=add_to_replay_stable_set, investigate_repeat_refund_pattern, monitor_in_shadow_mode, require_provenance_for_delegated_purchase, require_step_up_for_delegated_refunds, tighten_refund_review_rule
- URL: `http://127.0.0.1:8090/api/v1/operator/recommendations`

```text
{
  "data": [
    {
      "id": "rec-round-20260327001551-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260327001551",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260327001551-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260327001551-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260327001551-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260327001551-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260327001551-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260327001551",
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
      "id": "rec-round-20260327001551-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260327001551",
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
        "promo-round-20260327001551-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260327001551-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260327001551-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
      "linked_round_id": "round-20260327001551",
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
      "id": "rec-round-20260327001551-delegated-provenance",
      "type": "require_provenance_for_delegated_purchase",
      "rationale": "Delegated purchase paths with weak or missing provenance should not be treated as equivalent to ordinary human commerce, especially when they drift into new behavior patterns.",
      "priority": "moderate",
      "linked_round_id": "round-20260327001551",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s1-refund-weak-authorization",
        "commerce-s2-delegated-purchase-weak-provenance",
        "commerce-v1-weakened-provenance",
        "commerce-v4-actor-switch-human-to-agent",
        "commerce-v5-repeat-attempt-escalation"
      ],
      "linked_promotion_ids": [
        "promo-round-20260327001551-commerce-challenger-weakened-provenance-purchase-regression"
      ],
      "supporting_rule_ids": [
        "missing_provenance_sensitive_action"
      ],
      "suggested_action": "Require provenance for delegated purchases or keep them in recommendation-only shadow review until the team is comfortable tightening incumbent purchase controls.",
      "existing_control_integration_note": "Useful as a sidecar recommendation beside current purchase review thresholds and PSP checks."
    },
    {
      "id": "rec-round-20260327001551-shadow",
      "type": "monitor_in_shadow_mode",
      "rationale": "This round is best used as a recommendation-only sidecar beside the incumbent fraud stack so the team can compare benchmark findings against existing review and decision outcomes.",
      "priority": "low",
      "linked_round_id": "round-20260327001551",
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
        "promo-round-20260327001551-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260327001551-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260327001551-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260327001551-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Keep the harness in shadow mode, compare its outputs with current queue and policy outcomes, and only trial control changes after replay confirms the improvement.",
      "existing_control_integration_note": "Designed to run beside existing fraud rules, queueing, and PSP controls without blocking live traffic."
    },
    {
      "id": "rec-round-20260326235545-replay",
      "type": "add_to_replay_stable_set",
      "rationale": "Promoted challenger cases exposed meaningful misses or regressions and should move into replay so future benchmark rounds preserve the gain instead of rediscovering the same failure.",
      "priority": "high",
      "linked_round_id": "round-20260326235545",
      "linked_scenario_ids": [
        "commerce-challenger-weakened-provenance-purchase",
        "commerce-s3-approval-removed-after-authorization",
        "commerce-v2-expired-inactive-mandate",
        "commerce-v3-approval-removed"
      ],
      "linked_promotion_ids": [
        "promo-round-20260326235545-commerce-challenger-weakened-provenance-purchase-regression",
        "promo-round-20260326235545-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260326235545-commerce-v2-expired-inactive-mandate-regression",
        "promo-round-20260326235545-commerce-v3-approval-removed-regression"
      ],
      "suggested_action": "Add the promoted challenger cases to the replay stable set and re-run the sidecar benchmark before changing any incumbent production rule.",
      "existing_control_integration_note": "Use replay promotion as a safe bridge between benchmark discovery and incumbent fraud-stack tuning."
    },
    {
      "id": "rec-round-20260326235545-refund-review",
      "type": "tighten_refund_review_rule",
      "rationale": "Refund scenarios still surfaced weak-authorization paths that the incumbent fraud stack should queue or review more aggressively before payout or refund completion.",
      "priority": "high",
      "linked_round_id": "round-20260326235545",
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
      "id": "rec-round-20260326235545-delegated-refund-step-up",
      "type": "require_step_up_for_delegated_refunds",
      "rationale": "Agent-driven refund flows without approval evidence remain too risky to treat like ordinary refund traffic and should be routed into explicit step-up review.",
      "priority": "high",
      "linked_round_id": "round-20260326235545",
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
        "promo-round-20260326235545-commerce-s3-approval-removed-after-authorization-regression",
        "promo-round-20260326235545-commerce-v3-approval-removed-regression"
      ],
      "supporting_rule_ids": [
        "agent_refund_without_approval"
      ],
      "suggested_action": "Require step-up or human approval evidence for delegated refunds and compare the recommendation-only sidecar output with current queue decisions.",
      "existing_control_integration_note": "Designed to augment incumbent delegated-refund controls rather than replace them."
    },
    {
      "id": "rec-round-20260326235545-repeat-refund",
      "type": "investigate_repeat_refund_pattern",
      "rationale": "Repeated suspicious refund behavior is persisting across replay and benchmark history, which makes it a good candidate for targeted investigation and queue tuning beside the current fraud stack.",
      "priority": "moderate",
     
```

### [PASS] GET /api/v1/benchmark/trends/summary

- Kind: `api`
- Summary: http 200; present_keys=5; missing=none
- URL: `http://127.0.0.1:8090/api/v1/benchmark/trends/summary`

```text
{
  "data": {
    "rounds_executed": 22,
    "promotions_over_time": [
      {
        "round_id": "round-20260327001551",
        "value": 4
      },
      {
        "round_id": "round-20260326235545",
        "value": 4
      },
      {
        "round_id": "round-20260326235501",
        "value": 4
      },
      {
        "round_id": "round-20260326225239",
        "value": 4
      },
      {
        "round_id": "round-20260326214625",
        "value": 4
      },
      {
        "round_id": "round-20260326202835",
        "value": 4
      },
      {
        "round_id": "round-20260326184056",
        "value": 4
      },
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
        "round_id": "round-20260327001551",
        "value": 0
      },
      {
        "round_id": "round-20260326235545",
        "value": 0
      },
      {
        "round_id": "round-20260326235501",
        "value": 0
      },
      {
        "round_id": "round-20260326225239",
        "value": 0
      },
      {
        "round_id": "round-20260326214625",
        "value": 0
      },
      {
        "round_id": "round-20260326202835",
        "value": 0
      },
      {
        "round_id": "round-20260326184056",
        "value": 0
      },
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
    "new_blind_spots_discovered": 22,
    "regressions_observed": 0,
    "recommendation_counts_by_type": {
      "add_to_replay_stable_set": 22,
      "investigate_repeat_refund_pattern": 11,
      "monitor_in_shadow_mode": 22,
      "require_provenance_for_delegated_purchase": 18,
      "require_step_up_for_delegated_refunds": 22,
      "tighten_refund_review_rule": 22
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
- Summary: http 200; present_keys=5; missing=none
- URL: `http://127.0.0.1:8090/api/v1/operator/trends/summary`

```text
{
  "data": {
    "rounds_executed": 22,
    "promotions_over_time": [
      {
        "round_id": "round-20260327001551",
        "value": 4
      },
      {
        "round_id": "round-20260326235545",
        "value": 4
      },
      {
        "round_id": "round-20260326235501",
        "value": 4
      },
      {
        "round_id": "round-20260326225239",
        "value": 4
      },
      {
        "round_id": "round-20260326214625",
        "value": 4
      },
      {
        "round_id": "round-20260326202835",
        "value": 4
      },
      {
        "round_id": "round-20260326184056",
        "value": 4
      },
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
        "round_id": "round-20260327001551",
        "value": 0
      },
      {
        "round_id": "round-20260326235545",
        "value": 0
      },
      {
        "round_id": "round-20260326235501",
        "value": 0
      },
      {
        "round_id": "round-20260326225239",
        "value": 0
      },
      {
        "round_id": "round-20260326214625",
        "value": 0
      },
      {
        "round_id": "round-20260326202835",
        "value": 0
      },
      {
        "round_id": "round-20260326184056",
        "value": 0
      },
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
    "new_blind_spots_discovered": 22,
    "regressions_observed": 0,
    "recommendation_counts_by_type": {
      "add_to_replay_stable_set": 22,
      "investigate_repeat_refund_pattern": 11,
      "monitor_in_shadow_mode": 22,
      "require_provenance_for_delegated_purchase": 18,
      "require_step_up_for_delegated_refunds": 22,
      "tighten_refund_review_rule": 22
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
    "executed_runs": 0,
    "dry_run": true,
    "last_started_at": "0001-01-01T00:00:00Z",
    "next_run_at": "0001-01-01T00:00:00Z"
  }
}
```

### [FAIL] GET /api/v1/operator/reports/round-20260327001551

- Kind: `api`
- Summary: http 200; artifacts=0; missing=detection-delta.json,promotion-report.json,recommendation-report.json,round-report.json,round-summary.json
- URL: `http://127.0.0.1:8090/api/v1/operator/reports/round-20260327001551`

```text
{
  "data": [
    {
      "round_id": "round-20260327001551",
      "artifact_name": "detection-delta.json",
      "path": "reports/round-20260327001551/detection-delta.json",
      "kind": "json"
    },
    {
      "round_id": "round-20260327001551",
      "artifact_name": "executive-summary.md",
      "path": "reports/round-20260327001551/executive-summary.md",
      "kind": "markdown"
    },
    {
      "round_id": "round-20260327001551",
      "artifact_name": "promotion-report.json",
      "path": "reports/round-20260327001551/promotion-report.json",
      "kind": "json"
    },
    {
      "round_id": "round-20260327001551",
      "artifact_name": "recommendation-report.json",
      "path": "reports/round-20260327001551/recommendation-report.json",
      "kind": "json"
    },
    {
      "round_id": "round-20260327001551",
      "artifact_name": "round-report.json",
      "path": "reports/round-20260327001551/round-report.json",
      "kind": "json"
    },
    {
      "round_id": "round-20260327001551",
      "artifact_name": "round-report.md",
      "path": "reports/round-20260327001551/round-report.md",
      "kind": "markdown"
    },
    {
      "round_id": "round-20260327001551",
      "artifact_name": "round-summary.json",
      "path": "reports/round-20260327001551/round-summary.json",
      "kind": "json"
    },
    {
      "round_id": "round-20260327001551",
      "artifact_name": "round-summary.md",
      "path": "reports/round-20260327001551/round-summary.md",
      "kind": "markdown"
    }
  ]
}
```
