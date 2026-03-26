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
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	servicereporting "clawbot-trust-lab/internal/services/reporting"
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
	if len(state.Rounds[0].Reports.Artifacts) != 5 {
		t.Fatalf("expected 5 report artifacts after recommendation backfill, got %d", len(state.Rounds[0].Reports.Artifacts))
	}
	if len(state.Rounds[0].PromotionResults) != 1 {
		t.Fatalf("expected reconstructed promotions, got %#v", state.Rounds[0].PromotionResults)
	}
	if len(state.Rounds[0].Recommendations) == 0 {
		t.Fatal("expected recommendations to be reconstructed for historical round")
	}
	if _, err := os.Stat(filepath.Join(roundDir, "recommendation-report.json")); err != nil {
		t.Fatalf("expected recommendation-report.json to be backfilled: %v", err)
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

func TestLoadHistoricalStateBackfillsLegacyRecommendationReportIdempotently(t *testing.T) {
	reportsDir := t.TempDir()
	round := benchmark.BenchmarkRound{
		ID:             "round-legacy",
		ScenarioFamily: "commerce",
		StartedAt:      time.Date(2026, 3, 24, 12, 0, 0, 0, time.UTC),
		CompletedAt:    time.Date(2026, 3, 24, 12, 2, 0, 0, time.UTC),
		Summary: benchmark.RoundSummary{
			RoundID:             "round-legacy",
			ScenarioFamily:      "commerce",
			StableScenarioCount: 2,
			ChallengerCount:     1,
			PromotionCount:      1,
			ReplayPassRate:      1,
			RobustnessOutcome:   benchmark.RobustnessOutcomeNewBlindSpotDiscovered,
			ImportantFindings:   []string{"Legacy challenger should now be replay-stable."},
		},
		ScenarioResults: []benchmark.ScenarioResult{{
			ID:                   "sr-legacy-1",
			ScenarioID:           "commerce-v1-weakened-provenance",
			SetKind:              benchmark.ScenarioSetLiving,
			DetectionResultRef:   "det-legacy-1",
			FinalDetectionStatus: "clean",
			FinalRecommendation:  "allow",
			TriggeredRuleIDs:     []string{"missing_provenance_sensitive_action"},
			ReplayCaseRefs:       []string{"rc-legacy-1"},
		}},
		PromotionResults: []benchmark.PromotionDecision{{
			ID:                  "promo-legacy-1",
			RoundID:             "round-legacy",
			ScenarioID:          "commerce-v1-weakened-provenance",
			ChallengerVariantID: "variant-v1-weakened-provenance",
			PromotionReason:     benchmark.PromotionReasonDetectorMiss,
			Rationale:           "Legacy challenger evaluated as clean.",
			DetectionResultRef:  "det-legacy-1",
			ReplayCaseRef:       "rc-legacy-1",
			ScenarioResultRef:   "sr-legacy-1",
			Promoted:            true,
			CreatedAt:           time.Date(2026, 3, 24, 12, 2, 0, 0, time.UTC),
		}},
	}

	roundDir := filepath.Join(reportsDir, round.ID)
	if err := os.MkdirAll(roundDir, 0o750); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	writeFixtureJSON(t, filepath.Join(roundDir, "round-summary.json"), round)
	writeFixtureJSON(t, filepath.Join(roundDir, "promotion-report.json"), round.PromotionResults)
	writeFixtureJSON(t, filepath.Join(roundDir, "detection-delta.json"), []benchmark.DetectionDelta{{ScenarioID: round.ScenarioResults[0].ScenarioID}})

	logger := slog.New(slog.NewTextHandler(bytes.NewBuffer(nil), nil))
	state := LoadHistoricalState(reportsDir, logger)
	if len(state.Rounds) != 1 {
		t.Fatalf("expected 1 historical round, got %d", len(state.Rounds))
	}

	backfilledPath := filepath.Join(roundDir, "recommendation-report.json")
	firstBody, err := os.ReadFile(backfilledPath)
	if err != nil {
		t.Fatalf("expected recommendation-report.json after backfill: %v", err)
	}
	if len(state.Rounds[0].Recommendations) == 0 {
		t.Fatal("expected reconstructed recommendations in bootstrapped round")
	}
	if !artifactNames(state.Rounds[0].Reports).Has("recommendation-report.json") {
		t.Fatalf("expected recommendation-report.json in artifact listing, got %#v", state.Rounds[0].Reports.Artifacts)
	}

	state = LoadHistoricalState(reportsDir, logger)
	secondBody, err := os.ReadFile(backfilledPath)
	if err != nil {
		t.Fatalf("ReadFile() second pass error = %v", err)
	}
	if string(firstBody) != string(secondBody) {
		t.Fatal("expected idempotent backfill rerun to preserve recommendation-report.json content")
	}
	if len(state.Rounds) != 1 || !artifactNames(state.Rounds[0].Reports).Has("recommendation-report.json") {
		t.Fatalf("expected stable report artifact listing after rerun, got %#v", state.Rounds)
	}
}

func TestApplyRecommendationReportOnlyFillsMissingFields(t *testing.T) {
	round := benchmark.BenchmarkRound{
		ID: "round-existing",
		Summary: benchmark.RoundSummary{
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "keep existing",
			RecommendedFollowUp: "keep current",
			Recommendations:     1,
		},
		Recommendations: []benchmark.Recommendation{{
			ID: "rec-existing",
		}},
	}

	applyRecommendationReport(&round, benchmark.RecommendationReport{
		EvaluationMode:                 "ignored",
		BlockingMode:                   "ignored",
		ExistingControlIntegrationNote: "ignored",
		RecommendedFollowUp:            "ignored",
		Recommendations:                []benchmark.Recommendation{{ID: "rec-new"}},
	})

	if round.Summary.EvaluationMode != "shadow" || round.Summary.BlockingMode != "recommendation_only" {
		t.Fatalf("expected existing production-bridge fields to win, got %#v", round.Summary)
	}
	if round.Summary.ExistingControlNote != "keep existing" || round.Summary.RecommendedFollowUp != "keep current" {
		t.Fatalf("expected existing summary text to be preserved, got %#v", round.Summary)
	}
	if len(round.Recommendations) != 1 || round.Recommendations[0].ID != "rec-existing" {
		t.Fatalf("expected existing recommendations to be preserved, got %#v", round.Recommendations)
	}
}

func TestReadJSONRejectsEscapingAndAbsolutePaths(t *testing.T) {
	root := t.TempDir()

	for _, relPath := range []string{"../round-summary.json", "/tmp/round-summary.json", ""} {
		var payload map[string]any
		if err := readJSON(root, relPath, &payload); err == nil {
			t.Fatalf("expected readJSON(%q) to fail", relPath)
		}
	}
}

func TestReadOptionalJSONMissingFileReturnsFalse(t *testing.T) {
	root := t.TempDir()
	var payload map[string]any
	ok, err := readOptionalJSON(root, servicereporting.ArtifactRecommendationJSON, &payload)
	if err != nil {
		t.Fatalf("readOptionalJSON() error = %v", err)
	}
	if ok {
		t.Fatal("expected missing optional JSON file to return ok=false")
	}
}

func TestListReportArtifactsClassifiesUnknownFiles(t *testing.T) {
	roundDir := t.TempDir()
	for name, body := range map[string]string{
		servicereporting.ArtifactRoundSummaryJSON:   "{}\n",
		servicereporting.ArtifactExecutiveSummaryMD: "# Summary\n",
		"notes.txt": "plain text\n",
	} {
		if err := os.WriteFile(filepath.Join(roundDir, name), []byte(body), 0o600); err != nil {
			t.Fatalf("WriteFile(%s) error = %v", name, err)
		}
	}

	index := listReportArtifacts("round-kind-check", roundDir)
	if len(index.Artifacts) != 3 {
		t.Fatalf("expected 3 artifacts, got %d", len(index.Artifacts))
	}
	if got := index.Artifacts[2]; got.Name != "round-summary.json" && got.Kind == "" {
		t.Fatalf("unexpected artifact listing %#v", index.Artifacts)
	}
	kinds := map[string]string{}
	for _, artifact := range index.Artifacts {
		kinds[artifact.Name] = artifact.Kind
	}
	if kinds["notes.txt"] != "file" || kinds[servicereporting.ArtifactExecutiveSummaryMD] != "markdown" || kinds[servicereporting.ArtifactRoundSummaryJSON] != "json" {
		t.Fatalf("unexpected kinds %#v", kinds)
	}
}

func TestReconstructDetectionResultsMapsHistoricalStatuses(t *testing.T) {
	round := benchmark.BenchmarkRound{
		ID:          "round-1",
		CompletedAt: time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
		ScenarioResults: []benchmark.ScenarioResult{
			{
				ID:                   "sr-clean",
				ScenarioID:           "scenario-clean",
				DetectionResultRef:   "det-clean",
				FinalDetectionStatus: detectionmodel.DetectionStatusClean,
				FinalRecommendation:  detectionmodel.RecommendationAllow,
			},
			{
				ID:                   "sr-blocked",
				ScenarioID:           "scenario-blocked",
				DetectionResultRef:   "det-blocked",
				FinalDetectionStatus: detectionmodel.DetectionStatusBlocked,
				FinalRecommendation:  detectionmodel.RecommendationBlock,
				TriggeredRuleIDs:     []string{"missing_mandate_delegated_action"},
			},
		},
	}

	results := reconstructDetectionResults(round)
	if len(results) != 2 {
		t.Fatalf("expected 2 reconstructed results, got %d", len(results))
	}
	if results[0].Grade != detectionmodel.RiskGradeLow || results[0].Score != 0 {
		t.Fatalf("expected clean result to stay low-risk, got %#v", results[0])
	}
	if results[1].Grade != detectionmodel.RiskGradeCritical || results[1].Score != 80 {
		t.Fatalf("expected blocked result to reconstruct as critical, got %#v", results[1])
	}
}

func TestHistoricalHelpersCoverStatusAndSortFallbacks(t *testing.T) {
	started := time.Date(2026, 3, 25, 10, 0, 0, 0, time.UTC)
	completed := started.Add(time.Minute)

	if got := historicalRoundSortKey(benchmark.BenchmarkRound{CompletedAt: completed, StartedAt: started}); !got.Equal(completed) {
		t.Fatalf("expected completed time to win, got %s", got)
	}
	if got := historicalRoundSortKey(benchmark.BenchmarkRound{StartedAt: started}); !got.Equal(started) {
		t.Fatalf("expected started time fallback, got %s", got)
	}
	if got := historicalRoundSortKey(benchmark.BenchmarkRound{}); !got.IsZero() {
		t.Fatalf("expected zero sort key fallback, got %s", got)
	}

	cases := []struct {
		status   detectionmodel.DetectionStatus
		score    int
		grade    detectionmodel.RiskGrade
		severity int
	}{
		{status: detectionmodel.DetectionStatusSuspicious, score: 15, grade: detectionmodel.RiskGradeModerate, severity: 10},
		{status: detectionmodel.DetectionStatusStepUpRequired, score: 40, grade: detectionmodel.RiskGradeHigh, severity: 20},
		{status: detectionmodel.DetectionStatusBlocked, score: 80, grade: detectionmodel.RiskGradeCritical, severity: 30},
		{status: detectionmodel.DetectionStatusClean, score: 0, grade: detectionmodel.RiskGradeLow, severity: 0},
	}

	for _, tc := range cases {
		if got := historicalScore(tc.status); got != tc.score {
			t.Fatalf("historicalScore(%s) = %d, want %d", tc.status, got, tc.score)
		}
		if got := historicalGrade(tc.status); got != tc.grade {
			t.Fatalf("historicalGrade(%s) = %s, want %s", tc.status, got, tc.grade)
		}
		if got := historicalSeverity(tc.status); got != tc.severity {
			t.Fatalf("historicalSeverity(%s) = %d, want %d", tc.status, got, tc.severity)
		}
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

type artifactNameSet map[string]struct{}

func (s artifactNameSet) Has(name string) bool {
	_, ok := s[name]
	return ok
}

func artifactNames(index benchmark.ReportIndex) artifactNameSet {
	out := artifactNameSet{}
	for _, artifact := range index.Artifacts {
		out[artifact.Name] = struct{}{}
	}
	return out
}
