export type RobustnessOutcome =
  | "improved"
  | "mixed"
  | "regressed"
  | "new_blind_spot_discovered";

export type DetectionStatus =
  | "clean"
  | "suspicious"
  | "step_up_required"
  | "blocked";

export type Recommendation = "allow" | "review" | "step_up" | "block";

export type PromotionReviewStatus =
  | "accepted"
  | "duplicate"
  | "needs_follow_up"
  | "false_signal";

export interface RoundSummary {
  round_id: string;
  scenario_family: string;
  stable_scenario_count: number;
  challenger_count: number;
  replay_retest_count: number;
  promotion_count: number;
  replay_pass_rate: number;
  robustness_outcome: RobustnessOutcome;
  important_findings: string[];
}

export interface StableSetResult {
  total_count: number;
  passed_count: number;
}

export interface LivingSetResult {
  total_count: number;
  caught_count: number;
  promotion_count: number;
}

export interface PromotionDecision {
  id: string;
  round_id: string;
  scenario_id: string;
  challenger_variant_id?: string;
  promotion_reason: string;
  rationale: string;
  detection_result_ref: string;
  replay_case_ref?: string;
  scenario_result_ref: string;
  promoted: boolean;
  created_at: string;
}

export interface ScenarioResult {
  id: string;
  scenario_id: string;
  set_kind: string;
  challenger_variant_id?: string;
  detection_result_ref: string;
  final_detection_status: DetectionStatus;
  final_recommendation: Recommendation;
  triggered_rule_ids: string[];
  replay_case_refs: string[];
  memory_record_refs: string[];
  notes: string[];
}

export interface DetectionDelta {
  scenario_id: string;
  set_kind: string;
  previous_round_id?: string;
  previous_status: DetectionStatus;
  current_status: DetectionStatus;
  score_delta: number;
  newly_triggered_rules: string[];
  cleared_rules: string[];
  recommendation_changed: boolean;
}

export interface ReportArtifact {
  name: string;
  path: string;
  kind: string;
}

export interface ReportDescriptor {
  round_id: string;
  artifact_name: string;
  path: string;
  kind: string;
}

export interface ReportIndex {
  round_id: string;
  directory: string;
  artifacts: ReportArtifact[];
}

export interface BenchmarkRound {
  id: string;
  scenario_family: string;
  round_status: string;
  stable_set: StableSetResult;
  living_set: LivingSetResult;
  summary: RoundSummary;
  scenario_results: ScenarioResult[];
  promotion_results: PromotionDecision[];
  delta: DetectionDelta[];
  reports: ReportIndex;
}

export interface OperatorNote {
  id: string;
  body: string;
  created_at: string;
}

export interface PromotionReview {
  promotion_id: string;
  status: PromotionReviewStatus;
  note?: OperatorNote;
  updated_at: string;
}

export interface PromotionRecord {
  round_id: string;
  promotion: PromotionDecision;
  review?: PromotionReview;
}

export interface DetectionResult {
  id: string;
  scenario_id: string;
  order_id?: string;
  refund_id?: string;
  status: DetectionStatus;
  score: number;
  grade: string;
  reason_codes: string[];
  recommendation: Recommendation;
  replay_case_refs: string[];
  trust_decision_refs: string[];
  metadata: Record<string, unknown>;
}

export interface PromotionDetail {
  round_id: string;
  promotion: PromotionDecision;
  review?: PromotionReview;
  detection_result: DetectionResult;
  scenario_result?: ScenarioResult;
}

export interface RoundComparison {
  current_round_id: string;
  previous_round_id: string;
  current_robustness: RobustnessOutcome;
  previous_robustness: RobustnessOutcome;
  promotions_count_delta: number;
  replay_pass_rate_delta: number;
  challenger_count_delta: number;
  important_findings_added: string[];
  important_findings_closed: string[];
  detection_delta_count: number;
}

export interface ReportContent {
  descriptor: ReportDescriptor;
  content: string;
}
