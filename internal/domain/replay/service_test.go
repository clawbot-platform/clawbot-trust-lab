package replay_test

import (
	"path/filepath"
	"testing"

	"clawbot-trust-lab/internal/domain/replay"
	"clawbot-trust-lab/internal/platform/store"
)

func TestCreateCaseWritesReplayArchive(t *testing.T) {
	replayStore, err := store.NewFileReplayStore(filepath.Join(t.TempDir(), "archive"))
	if err != nil {
		t.Fatalf("NewFileReplayStore() error = %v", err)
	}
	service := replay.NewService(replayStore)

	item, err := service.CreateCase(replay.CreateCaseInput{
		ScenarioID:              "scenario-1",
		TrustArtifactRefs:       []string{"ta-scenario-1"},
		BenchmarkRoundRef:       "bench-1",
		OutcomeSummary:          "baseline replay outcome",
		PromotionRecommendation: "promote",
		PromotionReason:         "matches expectations",
	})
	if err != nil {
		t.Fatalf("CreateCase() error = %v", err)
	}

	if item.ID == "" {
		t.Fatal("expected replay case id")
	}
	if len(service.ListCases()) != 1 {
		t.Fatalf("expected 1 replay case, got %d", len(service.ListCases()))
	}
}
