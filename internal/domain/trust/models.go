package trust

import "time"

type TrustArtifact struct {
	ID             string             `json:"id"`
	ArtifactType   string             `json:"artifact_type"`
	ScenarioID     string             `json:"scenario_id"`
	Mandate        MandateArtifact    `json:"mandate"`
	Provenance     ProvenanceArtifact `json:"provenance"`
	PolicyDecision PolicyDecisionRef  `json:"policy_decision"`
	CreatedAt      time.Time          `json:"created_at"`
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
