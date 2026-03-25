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
}
