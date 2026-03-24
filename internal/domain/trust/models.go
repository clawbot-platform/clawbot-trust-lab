package trust

import "time"

type TrustArtifact struct {
	ID               string              `json:"id"`
	ArtifactFamily   string              `json:"artifact_family"`
	ArtifactType     string              `json:"artifact_type"`
	SourceScenarioID string              `json:"source_scenario_id"`
	Summary          string              `json:"summary"`
	Metadata         map[string]any      `json:"metadata"`
	Mandate          *MandateArtifact    `json:"mandate,omitempty"`
	Provenance       *ProvenanceArtifact `json:"provenance,omitempty"`
	PolicyDecision   PolicyDecisionRef   `json:"policy_decision"`
	CreatedAt        time.Time           `json:"created_at"`
}

type MandateArtifact struct {
	Source      string `json:"source"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ProvenanceArtifact struct {
	SourceRepo string `json:"source_repo"`
	Revision   string `json:"revision"`
	RecordedBy string `json:"recorded_by"`
}

type PolicyDecisionRef struct {
	PolicyID      string `json:"policy_id"`
	PolicyVersion string `json:"policy_version"`
	Outcome       string `json:"outcome"`
}

type Mandate struct {
	ID              string    `json:"id"`
	PrincipalID     string    `json:"principal_id"`
	DelegateActorID string    `json:"delegate_actor_id"`
	AllowedActions  []string  `json:"allowed_actions"`
	SpendingLimit   int64     `json:"spending_limit"`
	ExpiresAt       time.Time `json:"expires_at"`
	Status          string    `json:"status"`
}

type ProvenanceRecord struct {
	ID          string    `json:"id"`
	ActorID     string    `json:"actor_id"`
	PrincipalID string    `json:"principal_id"`
	SourceType  string    `json:"source_type"`
	SourceRef   string    `json:"source_ref"`
	Confidence  float64   `json:"confidence"`
	CreatedAt   time.Time `json:"created_at"`
}

type ApprovalRecord struct {
	ID         string    `json:"id"`
	OrderID    string    `json:"order_id"`
	ActionType string    `json:"action_type"`
	ApproverID string    `json:"approver_id"`
	Outcome    string    `json:"outcome"`
	CreatedAt  time.Time `json:"created_at"`
}

type TrustDecision struct {
	ID             string    `json:"id"`
	EntityType     string    `json:"entity_type"`
	EntityID       string    `json:"entity_id"`
	Outcome        string    `json:"outcome"`
	ReasonCodes    []string  `json:"reason_codes"`
	MandateRef     string    `json:"mandate_ref"`
	ProvenanceRef  string    `json:"provenance_ref"`
	StepUpRequired bool      `json:"step_up_required"`
	RecordedAt     time.Time `json:"recorded_at"`
}
