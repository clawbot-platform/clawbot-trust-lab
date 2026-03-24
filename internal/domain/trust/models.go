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
