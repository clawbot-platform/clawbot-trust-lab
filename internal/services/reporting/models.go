package reporting

import (
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
)

type GeneratedReport struct {
	Directory string                     `json:"directory"`
	Artifacts []benchmark.ReportArtifact `json:"artifacts"`
}

type ReportWindow struct {
	Label       string    `json:"label"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	GeneratedAt time.Time `json:"generated_at"`
}

type OperationalHealthSummary struct {
	TrustLabStatus         string   `json:"trust_lab_status"`
	ControlPlaneStatus     string   `json:"control_plane_status"`
	MemoryStatus           string   `json:"memory_status"`
	HealthHistoryAvailable bool     `json:"health_history_available"`
	DegradedPeriods        []string `json:"degraded_periods,omitempty"`
	RecoveryNotes          []string `json:"recovery_notes,omitempty"`
	Note                   string   `json:"note"`
}

type ProductionBridgeSummary struct {
	EvaluationMode                 string `json:"evaluation_mode"`
	BlockingMode                   string `json:"blocking_mode"`
	ExistingControlIntegrationNote string `json:"existing_control_integration_note"`
	RecommendedFollowUp            string `json:"recommended_follow_up"`
}

type TierUsageSummary struct {
	ResultsEvaluated    int      `json:"results_evaluated"`
	TierAAvailableCount int      `json:"tier_a_available_count"`
	TierBAvailableCount int      `json:"tier_b_available_count"`
	TierCCapableCount   int      `json:"tier_c_capable_count"`
	TierCUsedCount      int      `json:"tier_c_used_count"`
	TierCOptional       bool     `json:"tier_c_optional"`
	UnknownScenarioIDs  []string `json:"unknown_scenario_ids,omitempty"`
	InterpretationNote  string   `json:"interpretation_note"`
}

type RecommendationTheme struct {
	Type          benchmark.RecommendationType `json:"type"`
	Count         int                          `json:"count"`
	ExampleAction string                       `json:"example_action,omitempty"`
}

type ReplayWorthyCaseSummary struct {
	ScenarioID       string                      `json:"scenario_id"`
	LatestRoundID    string                      `json:"latest_round_id"`
	PromotionCount   int                         `json:"promotion_count"`
	PromotionReasons []benchmark.PromotionReason `json:"promotion_reasons"`
	Rationale        string                      `json:"rationale"`
}

type ScenarioIssueSummary struct {
	ScenarioID string   `json:"scenario_id"`
	Count      int      `json:"count"`
	Reasons    []string `json:"reasons,omitempty"`
}

type RoundReport struct {
	ReportType              string                               `json:"report_type"`
	GeneratedAt             time.Time                            `json:"generated_at"`
	RoundID                 string                               `json:"round_id"`
	ScenarioFamily          string                               `json:"scenario_family"`
	StartedAt               time.Time                            `json:"started_at"`
	CompletedAt             time.Time                            `json:"completed_at"`
	ScenariosExecuted       int                                  `json:"scenarios_executed"`
	Summary                 benchmark.RoundSummary               `json:"summary"`
	RecommendationCounts    map[benchmark.RecommendationType]int `json:"recommendation_counts_by_type"`
	PromotionCounts         map[benchmark.PromotionReason]int    `json:"promotion_counts_by_reason"`
	Promotions              []benchmark.PromotionDecision        `json:"promotions"`
	Recommendations         []benchmark.Recommendation           `json:"recommendations"`
	Regressions             []benchmark.DetectionDelta           `json:"regressions"`
	NotableChallengerCases  []benchmark.PromotionDecision        `json:"notable_challenger_cases"`
	TierUsage               TierUsageSummary                     `json:"tier_usage"`
	ProductionBridgeSummary ProductionBridgeSummary              `json:"production_bridge_summary"`
}

type DryRunReport struct {
	ReportType                    string                               `json:"report_type"`
	GeneratedAt                   time.Time                            `json:"generated_at"`
	Window                        ReportWindow                         `json:"window"`
	RoundIDs                      []string                             `json:"round_ids"`
	TotalRounds                   int                                  `json:"total_rounds"`
	TotalPromotions               int                                  `json:"total_promotions"`
	TotalRecommendations          int                                  `json:"total_recommendations"`
	NewBlindSpotsDiscovered       int                                  `json:"new_blind_spots_discovered"`
	RegressionsObserved           int                                  `json:"regressions_observed"`
	RecommendationCounts          map[benchmark.RecommendationType]int `json:"recommendation_counts_by_type"`
	RecurringRecommendationThemes []RecommendationTheme                `json:"recurring_recommendation_themes"`
	NewReplayWorthyCases          []ReplayWorthyCaseSummary            `json:"new_replay_worthy_cases"`
	RecurringIssueScenarios       []ScenarioIssueSummary               `json:"recurring_issue_scenarios"`
	RobustnessOutcomeCounts       map[benchmark.RobustnessOutcome]int  `json:"robustness_outcome_counts"`
	OperationalHealth             OperationalHealthSummary             `json:"operational_health"`
	ProductionBridgeSummary       ProductionBridgeSummary              `json:"production_bridge_summary"`
	NotableFindings               []string                             `json:"notable_findings"`
}

type ManagementReport struct {
	ReportType                    string                    `json:"report_type"`
	GeneratedAt                   time.Time                 `json:"generated_at"`
	Window                        ReportWindow              `json:"window"`
	TotalRounds                   int                       `json:"total_rounds"`
	ExecutiveSummary              string                    `json:"executive_summary"`
	DRQValueFindings              []string                  `json:"drq_value_findings"`
	ConsistentIssueScenarios      []ScenarioIssueSummary    `json:"consistent_issue_scenarios"`
	ReplayBaselineCandidates      []ReplayWorthyCaseSummary `json:"replay_baseline_candidates"`
	OperationalHealth             OperationalHealthSummary  `json:"operational_health"`
	RecommendedNextProductionStep []string                  `json:"recommended_next_production_steps"`
	StakeholderNotes              []string                  `json:"stakeholder_notes"`
}
