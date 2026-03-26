package detection

import "time"

type DetectionStatus string
type RiskGrade string
type Recommendation string

const (
	DetectionStatusClean          DetectionStatus = "clean"
	DetectionStatusSuspicious     DetectionStatus = "suspicious"
	DetectionStatusStepUpRequired DetectionStatus = "step_up_required"
	DetectionStatusBlocked        DetectionStatus = "blocked"
)

const (
	RiskGradeLow      RiskGrade = "low"
	RiskGradeModerate RiskGrade = "moderate"
	RiskGradeHigh     RiskGrade = "high"
	RiskGradeCritical RiskGrade = "critical"
)

const (
	RecommendationAllow  Recommendation = "allow"
	RecommendationReview Recommendation = "review"
	RecommendationStepUp Recommendation = "step_up"
	RecommendationBlock  Recommendation = "block"
)

type RuleHit struct {
	RuleID   string         `json:"rule_id"`
	Title    string         `json:"title"`
	Severity int            `json:"severity"`
	Reason   string         `json:"reason"`
	Metadata map[string]any `json:"metadata"`
}

type RuleDefinition struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Severity    int    `json:"severity"`
}

type TierProfile struct {
	TierAAvailable bool     `json:"tier_a_available"`
	TierBAvailable bool     `json:"tier_b_available"`
	TierCAvailable bool     `json:"tier_c_available"`
	TierCUsed      bool     `json:"tier_c_used"`
	TierANotes     []string `json:"tier_a_notes"`
	TierBNotes     []string `json:"tier_b_notes"`
	TierCNotes     []string `json:"tier_c_notes"`
}

type DetectionContext struct {
	ScenarioID               string          `json:"scenario_id"`
	OrderID                  string          `json:"order_id,omitempty"`
	RefundID                 string          `json:"refund_id,omitempty"`
	TrustDecisionRefs        []string        `json:"trust_decision_refs"`
	ReplayCaseRefs           []string        `json:"replay_case_refs"`
	Features                 map[string]bool `json:"features"`
	Signals                  map[string]any  `json:"signals"`
	EventCount               int             `json:"event_count"`
	TrustEventCount          int             `json:"trust_event_count"`
	TrustDecisionReasonCount int             `json:"trust_decision_reason_count"`
	ReplayHistoryCount       int             `json:"replay_history_count"`
	MemoryContextPresent     bool            `json:"memory_context_present"`
	MemoryStatus             string          `json:"memory_status"`
	TierProfile              TierProfile     `json:"tier_profile"`
}

type DetectionResult struct {
	ID                string          `json:"id"`
	ScenarioID        string          `json:"scenario_id"`
	OrderID           string          `json:"order_id,omitempty"`
	RefundID          string          `json:"refund_id,omitempty"`
	TrustDecisionRefs []string        `json:"trust_decision_refs"`
	ReplayCaseRefs    []string        `json:"replay_case_refs"`
	Status            DetectionStatus `json:"status"`
	Score             int             `json:"score"`
	Grade             RiskGrade       `json:"grade"`
	TriggeredRules    []RuleHit       `json:"triggered_rules"`
	ReasonCodes       []string        `json:"reason_codes"`
	Recommendation    Recommendation  `json:"recommendation"`
	EvaluatedAt       time.Time       `json:"evaluated_at"`
	Metadata          map[string]any  `json:"metadata"`
}

type DetectionRunSummary struct {
	TotalByStatus map[DetectionStatus]int `json:"total_by_status"`
	Total         int                     `json:"total"`
	LastResultID  string                  `json:"last_result_id,omitempty"`
}
