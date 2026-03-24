package replay

import (
	"fmt"
	"strings"
	"time"
)

type Service struct {
	store ArchiveStore
}

type ArchiveStore interface {
	Create(ReplayCase) (ReplayCase, error)
	List() []ReplayCase
}

type CreateCaseInput struct {
	ScenarioID              string   `json:"scenario_id"`
	TrustArtifactRefs       []string `json:"trust_artifact_refs"`
	BenchmarkRoundRef       string   `json:"benchmark_round_ref"`
	OutcomeSummary          string   `json:"outcome_summary"`
	PromotionRecommendation string   `json:"promotion_recommendation"`
	PromotionReason         string   `json:"promotion_reason"`
}

func NewService(replayStore ArchiveStore) *Service {
	return &Service{store: replayStore}
}

func (s *Service) CreateCase(input CreateCaseInput) (ReplayCase, error) {
	if strings.TrimSpace(input.ScenarioID) == "" {
		return ReplayCase{}, fmt.Errorf("scenario_id is required")
	}
	if strings.TrimSpace(input.OutcomeSummary) == "" {
		return ReplayCase{}, fmt.Errorf("outcome_summary is required")
	}

	item := ReplayCase{
		ID:                "rc-" + input.ScenarioID + "-" + time.Now().UTC().Format("20060102150405"),
		ScenarioID:        input.ScenarioID,
		TrustArtifactRefs: append([]string(nil), input.TrustArtifactRefs...),
		BenchmarkRoundRef: input.BenchmarkRoundRef,
		OutcomeSummary:    input.OutcomeSummary,
		ArchiveRef: ReplayArchiveRef{
			Bucket: "local-replay-archive",
			Key:    "",
		},
		Promotion: ReplayPromotionDecision{
			Status:   input.PromotionRecommendation,
			Reason:   input.PromotionReason,
			Promoted: strings.EqualFold(input.PromotionRecommendation, "promote"),
		},
		RecordedAt: time.Now().UTC(),
	}
	item.ArchiveRef.Key = item.ID + ".json"

	return s.store.Create(item)
}

func (s *Service) ListCases() []ReplayCase {
	return s.store.List()
}
