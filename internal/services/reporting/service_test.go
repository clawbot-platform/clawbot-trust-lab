package reporting

import (
	"encoding/json"
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
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Run this harness beside the incumbent fraud stack and compare recommendations before changing production policy.",
			RecommendedFollowUp: "Review promoted cases.",
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
		Recommendations: []benchmark.Recommendation{{
			ID:                             "rec-round-20260325120000-replay",
			Type:                           benchmark.RecommendationTypeAddToReplayStableSet,
			Rationale:                      "Promoted challenger should move into replay coverage.",
			Priority:                       benchmark.RecommendationPriorityHigh,
			LinkedRoundID:                  "round-20260325120000",
			LinkedScenarioIDs:              []string{"commerce-challenger-weakened-provenance-purchase"},
			LinkedPromotionIDs:             []string{"promo-1"},
			SuggestedAction:                "Add promoted challenger into replay.",
			ExistingControlIntegrationNote: "Use replay to validate proposed sidecar changes before touching production controls.",
		}},
	}

	index, err := service.Generate(round)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if len(index.Artifacts) != 6 {
		t.Fatalf("expected 6 artifacts, got %d", len(index.Artifacts))
	}

	for _, artifact := range index.Artifacts {
		if _, err := os.Stat(artifact.Path); err != nil {
			t.Fatalf("expected artifact %s to exist: %v", artifact.Path, err)
		}
		if filepath.Dir(artifact.Path) != index.Directory {
			t.Fatalf("expected artifact %s under %s", artifact.Path, index.Directory)
		}
	}

	var recommendationPath string
	for _, artifact := range index.Artifacts {
		if artifact.Name == "recommendation-report.json" {
			recommendationPath = artifact.Path
			break
		}
	}
	if recommendationPath == "" {
		t.Fatal("expected recommendation-report.json artifact")
	}

	body, err := os.ReadFile(recommendationPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var report benchmark.RecommendationReport
	if err := json.Unmarshal(body, &report); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if report.RoundID != round.ID || report.EvaluationMode != "shadow" || report.BlockingMode != "recommendation_only" {
		t.Fatalf("unexpected recommendation report %#v", report)
	}
	if len(report.Recommendations) != 1 || report.Recommendations[0].LinkedRoundID != round.ID {
		t.Fatalf("expected structured recommendations in report, got %#v", report)
	}
}

func TestBackfillRecommendationReportIsIdempotent(t *testing.T) {
	baseDir := t.TempDir()
	roundDir := filepath.Join(baseDir, "round-20260325120000")
	if err := os.MkdirAll(roundDir, 0o750); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	round := benchmark.BenchmarkRound{
		ID: "round-20260325120000",
		Summary: benchmark.RoundSummary{
			RoundID:             "round-20260325120000",
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Run as a sidecar.",
			RecommendedFollowUp: "Keep monitoring.",
		},
		Recommendations: []benchmark.Recommendation{{
			ID:                "rec-1",
			Type:              benchmark.RecommendationTypeMonitorInShadowMode,
			Priority:          benchmark.RecommendationPriorityLow,
			LinkedRoundID:     "round-20260325120000",
			LinkedScenarioIDs: []string{"commerce-h1-direct-human-purchase"},
			Rationale:         "Continue sidecar monitoring.",
			SuggestedAction:   "Keep the recommendation-only harness beside current controls.",
		}},
	}

	written, err := BackfillRecommendationReport(roundDir, round)
	if err != nil {
		t.Fatalf("BackfillRecommendationReport() error = %v", err)
	}
	if !written {
		t.Fatal("expected first backfill call to write the artifact")
	}

	path := filepath.Join(roundDir, "recommendation-report.json")
	firstBody, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	written, err = BackfillRecommendationReport(roundDir, round)
	if err != nil {
		t.Fatalf("BackfillRecommendationReport() second call error = %v", err)
	}
	if written {
		t.Fatal("expected second backfill call to be idempotent")
	}

	secondBody, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() second call error = %v", err)
	}
	if string(firstBody) != string(secondBody) {
		t.Fatal("expected recommendation report content to remain unchanged on rerun")
	}
}
