package benchmark

import (
	"time"

	detectionmodel "clawbot-trust-lab/internal/domain/detection"
)

type BenchmarkRoundRef struct {
	ID         string    `json:"id"`
	Label      string    `json:"label"`
	ScenarioID string    `json:"scenario_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type RunInput struct {
	ScenarioFamily string `json:"scenario_family"`
}

type SchedulerControlInput struct {
	ScenarioFamily string `json:"scenario_family"`
	Interval       string `json:"interval,omitempty"`
	MaxRuns        int    `json:"max_runs,omitempty"`
	DryRun         bool   `json:"dry_run,omitempty"`
}

type StableSuiteRef struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type LivingSuiteRef struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	MutationPolicy string `json:"mutation_policy"`
}

type RegistrationRequest struct {
	StableSuite         StableSuiteRef `json:"stable_suite"`
	LivingSuite         LivingSuiteRef `json:"living_suite"`
	ScenarioPackID      string         `json:"scenario_pack_id"`
	ScenarioPackVersion string         `json:"scenario_pack_version"`
	ReplayCaseRefs      []string       `json:"replay_case_refs"`
	Notes               string         `json:"notes"`
}

type RegistrationResult struct {
	RegistrationID string    `json:"registration_id"`
	Status         string    `json:"status"`
	RegisteredAt   time.Time `json:"registered_at"`
}

type RoundStatus string
type ScenarioSetKind string
type PromotionReason string
type RobustnessOutcome string
type RecommendationType string
type RecommendationPriority string

const (
	RoundStatusRunning   RoundStatus = "running"
	RoundStatusCompleted RoundStatus = "completed"
)

const (
	ScenarioSetStable ScenarioSetKind = "stable"
	ScenarioSetLiving ScenarioSetKind = "living"
	ScenarioSetReplay ScenarioSetKind = "replay_regression"
)

const (
	PromotionReasonDetectorMiss             PromotionReason = "detector_miss"
	PromotionReasonSuspiciousBehaviorTooLow PromotionReason = "suspicious_behavior_scored_too_low"
	PromotionReasonNewTrustGapPattern       PromotionReason = "new_trust_gap_pattern"
	PromotionReasonMeaningfulRegression     PromotionReason = "meaningful_regression"
	PromotionReasonNovelEvasiveVariation    PromotionReason = "novel_evasive_variation"
)

const (
	RobustnessOutcomeImproved               RobustnessOutcome = "improved"
	RobustnessOutcomeMixed                  RobustnessOutcome = "mixed"
	RobustnessOutcomeRegressed              RobustnessOutcome = "regressed"
	RobustnessOutcomeNewBlindSpotDiscovered RobustnessOutcome = "new_blind_spot_discovered"
)

const (
	RecommendationTypeAddToReplayStableSet              RecommendationType = "add_to_replay_stable_set"
	RecommendationTypeTightenRefundReviewRule           RecommendationType = "tighten_refund_review_rule"
	RecommendationTypeRequireStepUpForDelegatedRefunds  RecommendationType = "require_step_up_for_delegated_refunds"
	RecommendationTypeRequireProvenanceForDelegatedBuys RecommendationType = "require_provenance_for_delegated_purchase"
	RecommendationTypeInvestigateRepeatRefundPattern    RecommendationType = "investigate_repeat_refund_pattern"
	RecommendationTypeMonitorInShadowMode               RecommendationType = "monitor_in_shadow_mode"
)

const (
	RecommendationPriorityLow      RecommendationPriority = "low"
	RecommendationPriorityModerate RecommendationPriority = "moderate"
	RecommendationPriorityHigh     RecommendationPriority = "high"
)

type ChallengerVariant struct {
	ID                     string                         `json:"id"`
	ScenarioID             string                         `json:"scenario_id"`
	Title                  string                         `json:"title"`
	Description            string                         `json:"description"`
	ChangeSet              []string                       `json:"change_set"`
	ExpectedMinimumStatus  detectionmodel.DetectionStatus `json:"expected_minimum_status"`
	ExpectedRecommendation detectionmodel.Recommendation  `json:"expected_recommendation"`
}

type ScenarioResult struct {
	ID                    string                         `json:"id"`
	ScenarioID            string                         `json:"scenario_id"`
	ScenarioFamily        string                         `json:"scenario_family"`
	SetKind               ScenarioSetKind                `json:"set_kind"`
	ChallengerVariantID   string                         `json:"challenger_variant_id,omitempty"`
	ExecutionStatus       string                         `json:"execution_status"`
	OrderRefs             []string                       `json:"order_refs"`
	RefundRefs            []string                       `json:"refund_refs"`
	TrustDecisionRefs     []string                       `json:"trust_decision_refs"`
	ReplayCaseRefs        []string                       `json:"replay_case_refs"`
	MemoryRecordRefs      []string                       `json:"memory_record_refs"`
	DetectionResultRef    string                         `json:"detection_result_ref"`
	FinalDetectionStatus  detectionmodel.DetectionStatus `json:"final_detection_status"`
	FinalRecommendation   detectionmodel.Recommendation  `json:"final_recommendation"`
	TriggeredRuleIDs      []string                       `json:"triggered_rule_ids"`
	PromotedToReplay      bool                           `json:"promoted_to_replay"`
	ExpectedMinimumStatus detectionmodel.DetectionStatus `json:"expected_minimum_status,omitempty"`
	Passed                bool                           `json:"passed"`
	Notes                 []string                       `json:"notes"`
}

type PromotionDecision struct {
	ID                  string          `json:"id"`
	RoundID             string          `json:"round_id"`
	ScenarioID          string          `json:"scenario_id"`
	ChallengerVariantID string          `json:"challenger_variant_id,omitempty"`
	PromotionReason     PromotionReason `json:"promotion_reason"`
	Rationale           string          `json:"rationale"`
	DetectionResultRef  string          `json:"detection_result_ref"`
	ReplayCaseRef       string          `json:"replay_case_ref,omitempty"`
	ScenarioResultRef   string          `json:"scenario_result_ref"`
	Promoted            bool            `json:"promoted"`
	CreatedAt           time.Time       `json:"created_at"`
}

type DetectionDelta struct {
	ScenarioID            string                         `json:"scenario_id"`
	SetKind               ScenarioSetKind                `json:"set_kind"`
	PreviousRoundID       string                         `json:"previous_round_id,omitempty"`
	PreviousStatus        detectionmodel.DetectionStatus `json:"previous_status"`
	CurrentStatus         detectionmodel.DetectionStatus `json:"current_status"`
	ScoreDelta            int                            `json:"score_delta"`
	NewlyTriggeredRules   []string                       `json:"newly_triggered_rules"`
	ClearedRules          []string                       `json:"cleared_rules"`
	RecommendationChanged bool                           `json:"recommendation_changed"`
}

type StableSetResult struct {
	ScenarioRefs        []string `json:"scenario_refs"`
	DetectionResultRefs []string `json:"detection_result_refs"`
	PassedCount         int      `json:"passed_count"`
	TotalCount          int      `json:"total_count"`
}

type LivingSetResult struct {
	VariantRefs         []string `json:"variant_refs"`
	DetectionResultRefs []string `json:"detection_result_refs"`
	CaughtCount         int      `json:"caught_count"`
	PromotionCount      int      `json:"promotion_count"`
	TotalCount          int      `json:"total_count"`
}

type ReportArtifact struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Kind string `json:"kind"`
}

type ReportIndex struct {
	RoundID   string           `json:"round_id"`
	Directory string           `json:"directory"`
	Artifacts []ReportArtifact `json:"artifacts"`
}

type ReportDescriptor struct {
	RoundID      string `json:"round_id"`
	ArtifactName string `json:"artifact_name"`
	Path         string `json:"path"`
	Kind         string `json:"kind"`
}

type RoundSummary struct {
	RoundID             string            `json:"round_id"`
	ScenarioFamily      string            `json:"scenario_family"`
	StableScenarioCount int               `json:"stable_scenario_count"`
	ChallengerCount     int               `json:"challenger_count"`
	ReplayRetestCount   int               `json:"replay_retest_count"`
	PromotionCount      int               `json:"promotion_count"`
	ReplayPassRate      float64           `json:"replay_pass_rate"`
	RobustnessOutcome   RobustnessOutcome `json:"robustness_outcome"`
	ImportantFindings   []string          `json:"important_findings"`
	EvaluationMode      string            `json:"evaluation_mode"`
	BlockingMode        string            `json:"blocking_mode"`
	ExistingControlNote string            `json:"existing_control_integration_note"`
	RecommendedFollowUp string            `json:"recommended_follow_up"`
	Recommendations     int               `json:"recommendations"`
	TierCUsageCount     int               `json:"tier_c_usage_count"`
}

type Recommendation struct {
	ID                             string                 `json:"id"`
	Type                           RecommendationType     `json:"type"`
	Rationale                      string                 `json:"rationale"`
	Priority                       RecommendationPriority `json:"priority"`
	LinkedRoundID                  string                 `json:"linked_round_id"`
	LinkedScenarioIDs              []string               `json:"linked_scenario_ids"`
	LinkedPromotionIDs             []string               `json:"linked_promotion_ids,omitempty"`
	SupportingRuleIDs              []string               `json:"supporting_rule_ids,omitempty"`
	SuggestedAction                string                 `json:"suggested_action"`
	ExistingControlIntegrationNote string                 `json:"existing_control_integration_note,omitempty"`
}

type RecommendationReport struct {
	RoundID                        string           `json:"round_id"`
	EvaluationMode                 string           `json:"evaluation_mode"`
	BlockingMode                   string           `json:"blocking_mode"`
	ExistingControlIntegrationNote string           `json:"existing_control_integration_note"`
	RecommendedFollowUp            string           `json:"recommended_follow_up"`
	Recommendations                []Recommendation `json:"recommendations"`
}

type LongRunSummary struct {
	RoundsExecuted         int                        `json:"rounds_executed"`
	PromotionsOverTime     []MetricPoint              `json:"promotions_over_time"`
	ReplayPassRateOverTime []MetricPoint              `json:"replay_pass_rate_over_time"`
	NewBlindSpots          int                        `json:"new_blind_spots_discovered"`
	RegressionsObserved    int                        `json:"regressions_observed"`
	RecommendationCounts   map[RecommendationType]int `json:"recommendation_counts_by_type"`
	TopRecurringPatterns   []string                   `json:"top_recurring_evasion_patterns"`
}

type MetricPoint struct {
	RoundID string  `json:"round_id"`
	Value   float64 `json:"value"`
}

type SchedulerStatus struct {
	Enabled        bool      `json:"enabled"`
	Running        bool      `json:"running"`
	ScenarioFamily string    `json:"scenario_family"`
	Interval       string    `json:"interval"`
	MaxRuns        int       `json:"max_runs"`
	ExecutedRuns   int       `json:"executed_runs"`
	DryRun         bool      `json:"dry_run"`
	LastRoundID    string    `json:"last_round_id,omitempty"`
	LastStartedAt  time.Time `json:"last_started_at,omitempty"`
	NextRunAt      time.Time `json:"next_run_at,omitempty"`
}

type BenchmarkRound struct {
	ID                    string              `json:"id"`
	ScenarioFamily        string              `json:"scenario_family"`
	DetectorVersion       string              `json:"detector_version"`
	StableScenarioRefs    []string            `json:"stable_scenario_refs"`
	ChallengerVariantRefs []string            `json:"challenger_variant_refs"`
	ReplayCaseRefs        []string            `json:"replay_case_refs"`
	StartedAt             time.Time           `json:"started_at"`
	CompletedAt           time.Time           `json:"completed_at"`
	RoundStatus           RoundStatus         `json:"round_status"`
	ReportDir             string              `json:"report_dir"`
	ScenarioResults       []ScenarioResult    `json:"scenario_results"`
	PromotionResults      []PromotionDecision `json:"promotion_results"`
	Delta                 []DetectionDelta    `json:"delta"`
	StableSet             StableSetResult     `json:"stable_set"`
	LivingSet             LivingSetResult     `json:"living_set"`
	Summary               RoundSummary        `json:"summary"`
	Reports               ReportIndex         `json:"reports"`
	Recommendations       []Recommendation    `json:"recommendations"`
}

type PromotionReviewStatus string

const (
	PromotionReviewAccepted      PromotionReviewStatus = "accepted"
	PromotionReviewDuplicate     PromotionReviewStatus = "duplicate"
	PromotionReviewNeedsFollowUp PromotionReviewStatus = "needs_follow_up"
	PromotionReviewFalseSignal   PromotionReviewStatus = "false_signal"
)

type OperatorNote struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type PromotionReview struct {
	PromotionID string                `json:"promotion_id"`
	Status      PromotionReviewStatus `json:"status"`
	Note        *OperatorNote         `json:"note,omitempty"`
	UpdatedAt   time.Time             `json:"updated_at"`
}

type RoundComparison struct {
	CurrentRoundID          string            `json:"current_round_id"`
	PreviousRoundID         string            `json:"previous_round_id"`
	CurrentRobustness       RobustnessOutcome `json:"current_robustness"`
	PreviousRobustness      RobustnessOutcome `json:"previous_robustness"`
	PromotionsCountDelta    int               `json:"promotions_count_delta"`
	ReplayPassRateDelta     float64           `json:"replay_pass_rate_delta"`
	ChallengerCountDelta    int               `json:"challenger_count_delta"`
	ImportantFindingsAdded  []string          `json:"important_findings_added"`
	ImportantFindingsClosed []string          `json:"important_findings_closed"`
	DetectionDeltaCount     int               `json:"detection_delta_count"`
}
