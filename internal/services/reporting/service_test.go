package reporting

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
)

func TestGenerateWritesRequiredArtifacts(t *testing.T) {
	baseDir := t.TempDir()
	service := NewService(baseDir)

	round := benchmark.BenchmarkRound{
		ID:             "round-20260325120000",
		ScenarioFamily: "commerce",
		StartedAt:      time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
		CompletedAt:    time.Date(2026, 3, 25, 12, 1, 0, 0, time.UTC),
		Summary: benchmark.RoundSummary{
			RoundID:             "round-20260325120000",
			ScenarioFamily:      "commerce",
			StableScenarioCount: 2,
			ChallengerCount:     3,
			ReplayRetestCount:   1,
			PromotionCount:      1,
			ReplayPassRate:      0.5,
			RobustnessOutcome:   benchmark.RobustnessOutcomeNewBlindSpotDiscovered,
			ImportantFindings:   []string{"Weakened provenance challenger was promoted."},
		},
		PromotionResults: []benchmark.PromotionDecision{{
			ID:              "promo-1",
			ScenarioID:      "commerce-challenger-weakened-provenance-purchase",
			PromotionReason: benchmark.PromotionReasonDetectorMiss,
			Rationale:       "Suspicious challenger behavior evaluated as clean.",
			Promoted:        true,
		}},
		Delta: []benchmark.DetectionDelta{{
			ScenarioID:      "commerce-clean-agent-assisted-purchase",
			PreviousRoundID: "round-previous",
		}},
	}

	index, err := service.Generate(round)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(index.Artifacts) != 5 {
		t.Fatalf("expected 5 artifacts, got %d", len(index.Artifacts))
	}

	for _, artifact := range index.Artifacts {
		if _, err := os.Stat(artifact.Path); err != nil {
			t.Fatalf("expected artifact %s to exist: %v", artifact.Path, err)
		}
		if filepath.Dir(artifact.Path) != index.Directory {
			t.Fatalf("expected artifact %s under %s", artifact.Path, index.Directory)
		}
	}
}
