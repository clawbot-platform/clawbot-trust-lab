import type {
  BenchmarkRound,
  BenchmarkRecommendation,
  DetectionResult,
  LongRunSummary,
  PromotionDetail,
  PromotionRecord,
  PromotionReview,
  ReportContent,
  ReportDescriptor,
  RoundComparison
} from "../types/api";

export const previousRound: BenchmarkRound = {
  id: "round-20260324120000",
  scenario_family: "commerce",
  round_status: "completed",
  stable_set: {
    total_count: 2,
    passed_count: 2
  },
  living_set: {
    total_count: 2,
    caught_count: 2,
    promotion_count: 0
  },
  summary: {
    round_id: "round-20260324120000",
    scenario_family: "commerce",
    stable_scenario_count: 2,
    challenger_count: 2,
    replay_retest_count: 1,
    promotion_count: 0,
    replay_pass_rate: 1,
    robustness_outcome: "mixed",
    important_findings: ["Baseline challenger set fully caught."],
    evaluation_mode: "shadow",
    blocking_mode: "recommendation_only",
    existing_control_integration_note: "Run beside an incumbent fraud stack.",
    recommended_follow_up: "Continue observing in shadow mode.",
    recommendations: 1,
    tier_c_usage_count: 0
  },
  scenario_results: [],
  promotion_results: [],
  delta: [],
  recommendations: [],
  reports: {
    round_id: "round-20260324120000",
    directory: "./reports/round-20260324120000",
    artifacts: []
  }
};

export const currentRound: BenchmarkRound = {
  id: "round-20260325120000",
  scenario_family: "commerce",
  round_status: "completed",
  stable_set: {
    total_count: 2,
    passed_count: 2
  },
  living_set: {
    total_count: 3,
    caught_count: 2,
    promotion_count: 1
  },
  summary: {
    round_id: "round-20260325120000",
    scenario_family: "commerce",
    stable_scenario_count: 2,
    challenger_count: 3,
    replay_retest_count: 1,
    promotion_count: 1,
    replay_pass_rate: 0.67,
    robustness_outcome: "new_blind_spot_discovered",
    important_findings: [
      "Weakened provenance challenger still evaluated as clean.",
      "Replay regression remained stable for the refund case."
    ],
    evaluation_mode: "shadow",
    blocking_mode: "recommendation_only",
    existing_control_integration_note: "Run beside an incumbent fraud stack.",
    recommended_follow_up: "Add promoted challenger cases into replay and keep monitoring in shadow mode.",
    recommendations: 2,
    tier_c_usage_count: 1
  },
  scenario_results: [
    {
      id: "scenario-result-1",
      scenario_id: "commerce-challenger-weakened-provenance-purchase",
      set_kind: "living",
      challenger_variant_id: "variant-weakened-provenance",
      detection_result_ref: "det-1",
      final_detection_status: "clean",
      final_recommendation: "allow",
      triggered_rule_ids: ["delegated_actor_present"],
      replay_case_refs: ["rc-1"],
      memory_record_refs: ["mem-trust-1", "mem-replay-1"],
      notes: ["Detector missed the challenger intent."]
    },
    {
      id: "scenario-result-replay-1",
      scenario_id: "commerce-suspicious-refund-attempt",
      set_kind: "replay_regression",
      detection_result_ref: "det-2",
      final_detection_status: "step_up_required",
      final_recommendation: "step_up",
      triggered_rule_ids: ["refund_weak_authorization", "agent_refund_without_approval"],
      replay_case_refs: ["rc-refund-1"],
      memory_record_refs: ["mem-replay-refund-1"],
      notes: ["Prior promoted replay remained caught."]
    }
  ],
  promotion_results: [
    {
      id: "promo-1",
      round_id: "round-20260325120000",
      scenario_id: "commerce-challenger-weakened-provenance-purchase",
      challenger_variant_id: "variant-weakened-provenance",
      promotion_reason: "detector_miss",
      rationale: "Suspicious challenger behavior evaluated as clean despite degraded provenance confidence.",
      detection_result_ref: "det-1",
      replay_case_ref: "rc-1",
      scenario_result_ref: "scenario-result-1",
      promoted: true,
      created_at: "2026-03-25T12:00:00Z"
    },
    {
      id: "promo-2",
      round_id: "round-20260325120000",
      scenario_id: "commerce-s4-repeated-agent-refund-attempts",
      challenger_variant_id: "variant-repeat-attempt",
      promotion_reason: "meaningful_regression",
      rationale: "Repeated refund attempts remained too permissive across replay context.",
      detection_result_ref: "det-3",
      replay_case_ref: "rc-2",
      scenario_result_ref: "scenario-result-2",
      promoted: true,
      created_at: "2026-03-25T12:03:00Z"
    },
    {
      id: "promo-3",
      round_id: "round-20260325120000",
      scenario_id: "commerce-s5-merchant-scope-drift",
      challenger_variant_id: "variant-scope-drift",
      promotion_reason: "new_trust_gap_pattern",
      rationale: "Delegated purchase drifted outside merchant scope without enough friction.",
      detection_result_ref: "det-4",
      replay_case_ref: "rc-3",
      scenario_result_ref: "scenario-result-3",
      promoted: true,
      created_at: "2026-03-25T12:04:00Z"
    },
    {
      id: "promo-4",
      round_id: "round-20260324120000",
      scenario_id: "commerce-v7-high-value-delegated-purchase",
      challenger_variant_id: "variant-high-value",
      promotion_reason: "suspicious_behavior_scored_too_low",
      rationale: "High-value delegated purchase still landed below the expected review threshold.",
      detection_result_ref: "det-5",
      replay_case_ref: "rc-4",
      scenario_result_ref: "scenario-result-4",
      promoted: true,
      created_at: "2026-03-24T12:04:00Z"
    }
  ],
  delta: [
    {
      scenario_id: "commerce-challenger-weakened-provenance-purchase",
      set_kind: "living",
      previous_round_id: "round-20260324120000",
      previous_status: "suspicious",
      current_status: "clean",
      score_delta: -15,
      newly_triggered_rules: [],
      cleared_rules: ["missing_provenance_sensitive_action"],
      recommendation_changed: true
    }
  ],
  recommendations: [
    {
      id: "rec-round-20260325120000-replay",
      type: "add_to_replay_stable_set",
      rationale: "Promoted challenger cases should move into replay.",
      priority: "high",
      linked_round_id: "round-20260325120000",
      linked_scenario_ids: ["commerce-challenger-weakened-provenance-purchase"],
      linked_promotion_ids: ["promo-1"],
      supporting_rule_ids: ["missing_provenance_sensitive_action"],
      suggested_action: "Add the promoted challenger into replay coverage.",
      existing_control_integration_note: "Use this as a recommendation beside the incumbent refund and delegated-purchase rules."
    },
    {
      id: "rec-round-20260325120000-shadow",
      type: "monitor_in_shadow_mode",
      rationale: "Keep the harness beside the incumbent fraud stack.",
      priority: "low",
      linked_round_id: "round-20260325120000",
      linked_scenario_ids: ["commerce-challenger-weakened-provenance-purchase"],
      supporting_rule_ids: ["repeat_suspicious_context"],
      suggested_action: "Continue comparing sidecar recommendations with production outcomes.",
      existing_control_integration_note: "No production decline is required; compare recommendation drift against current review queues."
    },
    {
      id: "rec-round-20260325120000-refund",
      type: "require_step_up_for_delegated_refunds",
      rationale: "Repeated agent refunds still need stronger manual friction when approval evidence is missing.",
      priority: "medium",
      linked_round_id: "round-20260325120000",
      linked_scenario_ids: ["commerce-s4-repeated-agent-refund-attempts"],
      linked_promotion_ids: ["promo-2"],
      supporting_rule_ids: ["refund_weak_authorization", "agent_refund_without_approval"],
      suggested_action: "Route delegated refund retries into step-up or analyst review instead of straight-through handling.",
      existing_control_integration_note: "Add this as a sidecar recommendation before changing production refund policy."
    }
  ],
  reports: {
    round_id: "round-20260325120000",
    directory: "./reports/round-20260325120000",
    artifacts: [
      {
        name: "executive-summary.md",
        path: "./reports/round-20260325120000/executive-summary.md",
        kind: "markdown"
      },
      {
        name: "round-summary.json",
        path: "./reports/round-20260325120000/round-summary.json",
        kind: "json"
      }
    ]
  }
};

export const rounds = [currentRound, previousRound];

export const roundComparison: RoundComparison = {
  current_round_id: currentRound.id,
  previous_round_id: previousRound.id,
  current_robustness: currentRound.summary.robustness_outcome,
  previous_robustness: previousRound.summary.robustness_outcome,
  promotions_count_delta: 1,
  replay_pass_rate_delta: -0.33,
  challenger_count_delta: 1,
  important_findings_added: ["Weakened provenance blind spot promoted into replay."],
  important_findings_closed: [],
  detection_delta_count: 1
};

export const longRunSummary: LongRunSummary = {
  rounds_executed: 2,
  promotions_over_time: [
    { round_id: previousRound.id, value: 0 },
    { round_id: currentRound.id, value: 1 }
  ],
  replay_pass_rate_over_time: [
    { round_id: previousRound.id, value: 1 },
    { round_id: currentRound.id, value: 0.67 }
  ],
  new_blind_spots_discovered: 1,
  regressions_observed: 0,
  recommendation_counts_by_type: {
    add_to_replay_stable_set: 1,
    monitor_in_shadow_mode: 1
  },
  top_recurring_evasion_patterns: ["Weakened provenance challenger still evaluated as clean."]
};

export const detectionResult: DetectionResult = {
  id: "det-1",
  scenario_id: "commerce-challenger-weakened-provenance-purchase",
  order_id: "order-1",
  status: "clean",
  score: 12,
  grade: "low",
  reason_codes: ["delegated_actor_present", "provenance_present"],
  recommendation: "allow",
  replay_case_refs: ["rc-1"],
  trust_decision_refs: ["trust-decision-1"],
  metadata: {
    memory_context_present: true,
    tier_profile: {
      tier_c_used: true
    }
  }
};

export const initialPromotionReview: PromotionReview = {
  promotion_id: "promo-1",
  status: "needs_follow_up",
  note: {
    id: "note-1",
    body: "Investigate whether the provenance confidence threshold should be stricter.",
    created_at: "2026-03-25T12:05:00Z"
  },
  updated_at: "2026-03-25T12:05:00Z"
};

export const updatedPromotionReview: PromotionReview = {
  promotion_id: "promo-1",
  status: "accepted",
  note: {
    id: "note-2",
    body: "Promote this challenger into replay coverage.",
    created_at: "2026-03-25T12:10:00Z"
  },
  updated_at: "2026-03-25T12:10:00Z"
};

export const promotionRecord: PromotionRecord = {
  round_id: currentRound.id,
  promotion: currentRound.promotion_results[0],
  review: initialPromotionReview
};

export const promotionRecords: PromotionRecord[] = [
  {
    round_id: currentRound.id,
    promotion: currentRound.promotion_results[0],
    review: initialPromotionReview
  },
  {
    round_id: currentRound.id,
    promotion: currentRound.promotion_results[1],
    review: {
      promotion_id: "promo-2",
      status: "accepted",
      note: {
        id: "note-3",
        body: "Keep this in replay until refund retry behavior stabilizes.",
        created_at: "2026-03-25T12:12:00Z"
      },
      updated_at: "2026-03-25T12:12:00Z"
    }
  },
  {
    round_id: currentRound.id,
    promotion: currentRound.promotion_results[2]
  },
  {
    round_id: previousRound.id,
    promotion: currentRound.promotion_results[3],
    review: {
      promotion_id: "promo-4",
      status: "duplicate",
      note: {
        id: "note-4",
        body: "Covered by an earlier high-value delegated purchase replay case.",
        created_at: "2026-03-24T12:10:00Z"
      },
      updated_at: "2026-03-24T12:10:00Z"
    }
  }
];

export const promotionDetail: PromotionDetail = {
  round_id: currentRound.id,
  promotion: currentRound.promotion_results[0],
  review: initialPromotionReview,
  detection_result: detectionResult,
  scenario_result: currentRound.scenario_results[0]
};

export const recommendations: BenchmarkRecommendation[] = currentRound.recommendations;

export const reportDescriptors: ReportDescriptor[] = currentRound.reports.artifacts.map((artifact) => ({
  round_id: currentRound.id,
  artifact_name: artifact.name,
  path: artifact.path,
  kind: artifact.kind
}));

export const executiveSummaryReport: ReportContent = {
  descriptor: reportDescriptors[0],
  content: "# Executive Summary\n\n- Blind spot persists for weakened provenance.\n- Refund replay remains caught.\n"
};

export const roundSummaryReport: ReportContent = {
  descriptor: reportDescriptors[1],
  content: JSON.stringify(
    {
      round_id: currentRound.id,
      promotion_count: currentRound.summary.promotion_count,
      replay_pass_rate: currentRound.summary.replay_pass_rate,
      robustness_outcome: currentRound.summary.robustness_outcome
    },
    null,
    2
  )
};
