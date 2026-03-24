package replay

import "time"

type ReplayCase struct {
	ID                string                  `json:"id"`
	ScenarioID        string                  `json:"scenario_id"`
	TrustArtifactRefs []string                `json:"trust_artifact_refs"`
	BenchmarkRoundRef string                  `json:"benchmark_round_ref"`
	OutcomeSummary    string                  `json:"outcome_summary"`
	ArchiveRef        ReplayArchiveRef        `json:"archive_ref"`
	Promotion         ReplayPromotionDecision `json:"promotion"`
	RecordedAt        time.Time               `json:"recorded_at"`
}

type ReplayArchiveRef struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type ReplayPromotionDecision struct {
	Status   string `json:"status"`
	Reason   string `json:"reason"`
	Promoted bool   `json:"promoted"`
}
