package replay

import (
	"context"
	"fmt"
	"strings"
	"time"

	"clawbot-trust-lab/internal/clients/memory"
)

type Service struct {
	store  ArchiveStore
	memory memory.Client
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

type MemorySyncError struct {
	Err error
}

func (e *MemorySyncError) Error() string {
	return "clawmem replay write failed: " + e.Err.Error()
}

func (e *MemorySyncError) Unwrap() error {
	return e.Err
}

func NewService(replayStore ArchiveStore, memoryClient memory.Client) *Service {
	return &Service{store: replayStore, memory: memoryClient}
}

func (s *Service) CreateCase(ctx context.Context, input CreateCaseInput) (ReplayCase, error) {
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

	if err := s.memory.StoreReplayCase(ctx, memory.StoreReplayCaseRequest{
		ReplayCaseID: item.ID,
		ScenarioID:   item.ScenarioID,
		Summary:      item.OutcomeSummary,
		Metadata: map[string]any{
			"benchmark_round_ref": item.BenchmarkRoundRef,
			"archive_ref":         item.ArchiveRef,
			"promotion":           item.Promotion,
			"trust_artifact_refs": item.TrustArtifactRefs,
		},
		Tags: []string{"replay-case", item.Promotion.Status},
	}); err != nil {
		return ReplayCase{}, &MemorySyncError{Err: err}
	}

	return s.store.Create(item)
}

func (s *Service) ListCases() []ReplayCase {
	return s.store.List()
}

func (s *Service) SimilarCases(ctx context.Context, scenarioID string) (memory.FetchSimilarCasesResponse, error) {
	return s.memory.FetchSimilarCases(ctx, memory.FetchSimilarCasesRequest{ScenarioID: strings.TrimSpace(scenarioID)})
}
