import type {
  BenchmarkRound,
  DetectionResult,
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
    important_findings: ["Baseline challenger set fully caught."]
  },
  scenario_results: [],
  promotion_results: [],
  delta: [],
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
    ]
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
    memory_context_present: true
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

export const promotionDetail: PromotionDetail = {
  round_id: currentRound.id,
  promotion: currentRound.promotion_results[0],
  review: initialPromotionReview,
  detection_result: detectionResult,
  scenario_result: currentRound.scenario_results[0]
};

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
