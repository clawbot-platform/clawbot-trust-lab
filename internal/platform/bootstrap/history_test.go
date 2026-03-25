package bootstrap

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
)

func TestLoadHistoricalStateDiscoversRoundsAndPromotions(t *testing.T) {
	reportsDir := t.TempDir()
	round := benchmark.BenchmarkRound{
		ID:             "round-20260325120000",
		ScenarioFamily: "commerce",
		StartedAt:      time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
		CompletedAt:    time.Date(2026, 3, 25, 12, 1, 0, 0, time.UTC),
		RoundStatus:    benchmark.RoundStatusCompleted,
		Summary: benchmark.RoundSummary{
			RoundID:             "round-20260325120000",
			ScenarioFamily:      "commerce",
			StableScenarioCount: 2,
			ChallengerCount:     3,
			PromotionCount:      1,
			ReplayPassRate:      1,
			RobustnessOutcome:   benchmark.RobustnessOutcomeNewBlindSpotDiscovered,
		},
		ScenarioResults: []benchmark.ScenarioResult{{
			ID:                   "sr-1",
			ScenarioID:           "commerce-challenger-weakened-provenance-purchase",
			SetKind:              benchmark.ScenarioSetLiving,
			DetectionResultRef:   "det-1",
			FinalDetectionStatus: "clean",
			FinalRecommendation:  "allow",
			TriggeredRuleIDs:     []string{"missing_provenance_sensitive_action"},
			TrustDecisionRefs:    []string{"trust-1"},
			ReplayCaseRefs:       []string{"replay-1"},
			OrderRefs:            []string{"order-1"},
		}},
		PromotionResults: []benchmark.PromotionDecision{{
			ID:                 "promo-1",
			RoundID:            "round-20260325120000",
			ScenarioID:         "commerce-challenger-weakened-provenance-purchase",
			DetectionResultRef: "det-1",
			ScenarioResultRef:  "sr-1",
			Promoted:           true,
			CreatedAt:          time.Date(2026, 3, 25, 12, 1, 0, 0, time.UTC),
		}},
	}

	roundDir := filepath.Join(reportsDir, round.ID)
	if err := os.MkdirAll(roundDir, 0o750); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	writeFixtureJSON(t, filepath.Join(roundDir, "round-summary.json"), round)
	writeFixtureJSON(t, filepath.Join(roundDir, "promotion-report.json"), round.PromotionResults)
	writeFixtureJSON(t, filepath.Join(roundDir, "detection-delta.json"), []benchmark.DetectionDelta{{ScenarioID: round.ScenarioResults[0].ScenarioID}})
	if err := os.WriteFile(filepath.Join(roundDir, "executive-summary.md"), []byte("# Executive Summary"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	state := LoadHistoricalState(reportsDir, slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil)))

	if len(state.Rounds) != 1 {
		t.Fatalf("expected 1 historical round, got %d", len(state.Rounds))
	}
	if state.Rounds[0].ReportDir != roundDir {
		t.Fatalf("expected report dir %s, got %s", roundDir, state.Rounds[0].ReportDir)
	}
	if len(state.Rounds[0].Reports.Artifacts) != 4 {
		t.Fatalf("expected 4 report artifacts, got %d", len(state.Rounds[0].Reports.Artifacts))
	}
	if len(state.Rounds[0].PromotionResults) != 1 {
		t.Fatalf("expected reconstructed promotions, got %#v", state.Rounds[0].PromotionResults)
	}
	if len(state.DetectionResults) != 1 {
		t.Fatalf("expected 1 reconstructed detection result, got %d", len(state.DetectionResults))
	}
	if state.DetectionResults[0].ID != "det-1" {
		t.Fatalf("unexpected detection result %#v", state.DetectionResults[0])
	}
}

func TestLoadHistoricalStateSkipsMalformedDirectories(t *testing.T) {
	reportsDir := t.TempDir()
	validRoundDir := filepath.Join(reportsDir, "round-valid")
	if err := os.MkdirAll(validRoundDir, 0o750); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	writeFixtureJSON(t, filepath.Join(validRoundDir, "round-summary.json"), benchmark.BenchmarkRound{
		ID:          "round-valid",
		CompletedAt: time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
		Summary:     benchmark.RoundSummary{RoundID: "round-valid", ScenarioFamily: "commerce"},
	})

	badRoundDir := filepath.Join(reportsDir, "round-bad")
	if err := os.MkdirAll(badRoundDir, 0o750); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(badRoundDir, "round-summary.json"), []byte("{not-json"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	state := LoadHistoricalState(reportsDir, logger)

	if len(state.Rounds) != 1 {
		t.Fatalf("expected only the valid round to load, got %d", len(state.Rounds))
	}
	if !strings.Contains(logs.String(), "round-bad") {
		t.Fatalf("expected malformed directory to be logged, got %s", logs.String())
	}
}

func writeFixtureJSON(t *testing.T, path string, payload any) {
	t.Helper()
	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent() error = %v", err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(path, body, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
