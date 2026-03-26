package reporting

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
	domainscenario "clawbot-trust-lab/internal/domain/scenario"
)

type scenarioCatalogStub struct {
	items []domainscenario.Scenario
}

func (s scenarioCatalogStub) ListScenarios() []domainscenario.Scenario {
	return append([]domainscenario.Scenario(nil), s.items...)
}

func TestGenerateWritesRequiredArtifacts(t *testing.T) {
	baseDir := t.TempDir()
	service := NewService(baseDir, scenarioCatalogStub{items: []domainscenario.Scenario{
		{
			ID: "commerce-v1-weakened-provenance",
			FeatureModel: domainscenario.FeatureTierModel{
				TierA: []string{"amount"},
				TierB: []string{"buyer_history"},
				TierC: []string{"provenance_confidence"},
			},
		},
	}})

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
		ScenarioResults: []benchmark.ScenarioResult{{
			ID:         "result-1",
			ScenarioID: "commerce-v1-weakened-provenance",
			Notes:      []string{"tier_c_used"},
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
	if len(index.Artifacts) != 8 {
		t.Fatalf("expected 8 artifacts, got %d", len(index.Artifacts))
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
	var roundReportPath string
	for _, artifact := range index.Artifacts {
		if artifact.Name == "recommendation-report.json" {
			recommendationPath = artifact.Path
		}
		if artifact.Name == "round-report.json" {
			roundReportPath = artifact.Path
		}
	}
	if recommendationPath == "" {
		t.Fatal("expected recommendation-report.json artifact")
	}
	if roundReportPath == "" {
		t.Fatal("expected round-report.json artifact")
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

	body, err = os.ReadFile(roundReportPath)
	if err != nil {
		t.Fatalf("ReadFile(round-report.json) error = %v", err)
	}
	var roundReport RoundReport
	if err := json.Unmarshal(body, &roundReport); err != nil {
		t.Fatalf("json.Unmarshal(round-report.json) error = %v", err)
	}
	if roundReport.RoundID != round.ID || roundReport.TierUsage.TierCCapableCount != 1 {
		t.Fatalf("unexpected round report %#v", roundReport)
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

func TestBuildRecommendationReportCopiesRoundSummaryFields(t *testing.T) {
	round := benchmark.BenchmarkRound{
		ID: "round-copy",
		Summary: benchmark.RoundSummary{
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Compare sidecar recommendations to the incumbent fraud stack.",
			RecommendedFollowUp: "Review refund rule tuning.",
		},
		Recommendations: []benchmark.Recommendation{{
			ID:                "rec-copy",
			LinkedRoundID:     "round-copy",
			LinkedScenarioIDs: []string{"commerce-s1-refund-weak-authorization"},
		}},
	}

	report := BuildRecommendationReport(round)
	if report.RoundID != round.ID || report.EvaluationMode != "shadow" || report.BlockingMode != "recommendation_only" {
		t.Fatalf("unexpected report header %#v", report)
	}
	if report.ExistingControlIntegrationNote != round.Summary.ExistingControlNote || report.RecommendedFollowUp != round.Summary.RecommendedFollowUp {
		t.Fatalf("expected summary text to copy through, got %#v", report)
	}
	if len(report.Recommendations) != 1 || report.Recommendations[0].LinkedRoundID != round.ID {
		t.Fatalf("expected recommendation linkage to be preserved, got %#v", report.Recommendations)
	}
}

func TestRoundSummaryMarkdownHandlesEmptySections(t *testing.T) {
	service := NewService(t.TempDir())
	round := benchmark.BenchmarkRound{
		ID:             "round-empty",
		ScenarioFamily: "commerce",
		Summary: benchmark.RoundSummary{
			ScenarioFamily:      "commerce",
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Sidecar only.",
			RecommendedFollowUp: "No action.",
		},
	}

	body := service.roundSummaryMarkdown(round)
	for _, snippet := range []string{
		"No notable findings were recorded.",
		"No challenger cases were promoted in this round.",
		"No explicit recommendations were generated.",
		"Recommended follow-up: No action.",
	} {
		if !strings.Contains(body, snippet) {
			t.Fatalf("expected markdown to contain %q, got %s", snippet, body)
		}
	}
}

func TestExecutiveSummaryUsesRegressionHeadline(t *testing.T) {
	service := NewService(t.TempDir())
	round := benchmark.BenchmarkRound{
		ID: "round-regressed",
		Summary: benchmark.RoundSummary{
			StableScenarioCount: 1,
			ChallengerCount:     2,
			ReplayRetestCount:   3,
			RobustnessOutcome:   benchmark.RobustnessOutcomeRegressed,
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Keep as sidecar.",
			RecommendedFollowUp: "Escalate replay review.",
		},
	}

	body := service.executiveSummary(round)
	for _, snippet := range []string{
		"Regression observed",
		"evaluated 1 stable scenarios, 2 challenger variants, and 3 replay retests.",
		"Escalate replay review.",
	} {
		if !strings.Contains(body, snippet) {
			t.Fatalf("expected executive summary to contain %q, got %s", snippet, body)
		}
	}
}

func TestGenerateDryRunReportWritesArtifacts(t *testing.T) {
	baseDir := t.TempDir()
	service := NewService(baseDir)
	service.now = func() time.Time { return time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC) }

	window := ReportWindow{
		Label:       "2026-03-25_to_2026-03-26",
		Start:       time.Date(2026, 3, 25, 15, 0, 0, 0, time.UTC),
		End:         time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC),
		GeneratedAt: service.now(),
	}
	rounds := []benchmark.BenchmarkRound{
		fixtureRound("round-1", time.Date(2026, 3, 25, 20, 0, 0, 0, time.UTC), benchmark.RobustnessOutcomeNewBlindSpotDiscovered),
		fixtureRound("round-2", time.Date(2026, 3, 26, 10, 0, 0, 0, time.UTC), benchmark.RobustnessOutcomeRegressed),
	}

	out, err := service.GenerateDryRunReport(window, rounds, OperationalHealthSummary{
		TrustLabStatus:         "ok",
		ControlPlaneStatus:     "ok",
		MemoryStatus:           "degraded",
		HealthHistoryAvailable: false,
		Note:                   "snapshot only",
	})
	if err != nil {
		t.Fatalf("GenerateDryRunReport() error = %v", err)
	}
	if len(out.Artifacts) != 2 {
		t.Fatalf("expected 2 dry-run artifacts, got %d", len(out.Artifacts))
	}

	body, err := os.ReadFile(filepath.Join(out.Directory, ArtifactDryRunReportJSON))
	if err != nil {
		t.Fatalf("ReadFile(dry-run-report.json) error = %v", err)
	}
	var report DryRunReport
	if err := json.Unmarshal(body, &report); err != nil {
		t.Fatalf("json.Unmarshal(dry-run-report.json) error = %v", err)
	}
	if report.TotalRounds != 2 || report.TotalPromotions != 2 || report.RegressionsObserved != 1 {
		t.Fatalf("unexpected dry-run report %#v", report)
	}
	if len(report.NewReplayWorthyCases) == 0 || report.OperationalHealth.MemoryStatus != "degraded" {
		t.Fatalf("expected replay-worthy cases and health snapshot, got %#v", report)
	}
}

func TestGenerateManagementReportWritesArtifacts(t *testing.T) {
	baseDir := t.TempDir()
	service := NewService(baseDir)
	service.now = func() time.Time { return time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC) }

	window := ReportWindow{
		Label:       "2026-03-19_to_2026-03-26",
		Start:       time.Date(2026, 3, 19, 0, 0, 0, 0, time.UTC),
		End:         time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC),
		GeneratedAt: service.now(),
	}
	rounds := []benchmark.BenchmarkRound{
		fixtureRound("round-1", time.Date(2026, 3, 21, 20, 0, 0, 0, time.UTC), benchmark.RobustnessOutcomeNewBlindSpotDiscovered),
		fixtureRound("round-2", time.Date(2026, 3, 24, 10, 0, 0, 0, time.UTC), benchmark.RobustnessOutcomeNewBlindSpotDiscovered),
	}

	out, err := service.GenerateManagementReport(window, rounds, OperationalHealthSummary{
		TrustLabStatus:         "ok",
		ControlPlaneStatus:     "ok",
		MemoryStatus:           "ok",
		HealthHistoryAvailable: false,
		Note:                   "snapshot only",
	})
	if err != nil {
		t.Fatalf("GenerateManagementReport() error = %v", err)
	}

	body, err := os.ReadFile(filepath.Join(out.Directory, ArtifactManagementJSON))
	if err != nil {
		t.Fatalf("ReadFile(management-report.json) error = %v", err)
	}
	var report ManagementReport
	if err := json.Unmarshal(body, &report); err != nil {
		t.Fatalf("json.Unmarshal(management-report.json) error = %v", err)
	}
	if report.TotalRounds != 2 || len(report.DRQValueFindings) == 0 || len(report.RecommendedNextProductionStep) == 0 {
		t.Fatalf("unexpected management report %#v", report)
	}
	if !strings.Contains(report.ExecutiveSummary, "2 completed round") {
		t.Fatalf("expected executive summary to mention completed rounds, got %q", report.ExecutiveSummary)
	}
}

func TestGenerateWindowReportsHandleEmptyInput(t *testing.T) {
	service := NewService(t.TempDir())
	service.now = func() time.Time { return time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC) }
	window := ReportWindow{
		Label:       "2026-03-25_to_2026-03-26",
		Start:       time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
		End:         time.Date(2026, 3, 26, 0, 0, 0, 0, time.UTC),
		GeneratedAt: service.now(),
	}
	health := OperationalHealthSummary{
		TrustLabStatus:         "ok",
		ControlPlaneStatus:     "ok",
		MemoryStatus:           "ok",
		HealthHistoryAvailable: false,
		Note:                   "snapshot only",
	}

	dryRun, err := service.GenerateDryRunReport(window, nil, health)
	if err != nil {
		t.Fatalf("GenerateDryRunReport() error = %v", err)
	}
	management, err := service.GenerateManagementReport(window, nil, health)
	if err != nil {
		t.Fatalf("GenerateManagementReport() error = %v", err)
	}

	dryBody, err := os.ReadFile(filepath.Join(dryRun.Directory, ArtifactDryRunReportMD))
	if err != nil {
		t.Fatalf("ReadFile(dry-run-report.md) error = %v", err)
	}
	if !strings.Contains(string(dryBody), "No benchmark rounds completed inside the selected window.") {
		t.Fatalf("expected empty dry-run note, got %s", string(dryBody))
	}

	managementBody, err := os.ReadFile(filepath.Join(management.Directory, ArtifactManagementMD))
	if err != nil {
		t.Fatalf("ReadFile(management-report.md) error = %v", err)
	}
	if !strings.Contains(string(managementBody), "No benchmark rounds completed inside the selected management-report window") {
		t.Fatalf("expected empty management note, got %s", string(managementBody))
	}
}

func TestReportingHelpersSortAndFilterDeterministically(t *testing.T) {
	window := ReportWindow{
		Start: time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2026, 3, 26, 23, 59, 59, 0, time.UTC),
	}
	rounds := []benchmark.BenchmarkRound{
		{ID: "late", CompletedAt: time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC)},
		{ID: "early", StartedAt: time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC)},
		{ID: "outside", CompletedAt: time.Date(2026, 3, 27, 12, 0, 0, 0, time.UTC)},
		{ID: "zero"},
	}
	filtered := filterRoundsForWindow(rounds, window)
	if len(filtered) != 2 || filtered[0].ID != "early" || filtered[1].ID != "late" {
		t.Fatalf("unexpected filtered rounds %#v", filtered)
	}

	themes := recommendationThemes(
		map[benchmark.RecommendationType]int{
			benchmark.RecommendationTypeMonitorInShadowMode:               2,
			benchmark.RecommendationTypeAddToReplayStableSet:              2,
			benchmark.RecommendationTypeRequireProvenanceForDelegatedBuys: 1,
		},
		map[benchmark.RecommendationType]string{
			benchmark.RecommendationTypeAddToReplayStableSet: "Promote into replay.",
		},
	)
	if len(themes) != 3 || themes[0].Type != benchmark.RecommendationTypeAddToReplayStableSet || themes[1].Type != benchmark.RecommendationTypeMonitorInShadowMode {
		t.Fatalf("unexpected recommendation themes %#v", themes)
	}

	replayCases := sortReplayCandidates(map[string]*ReplayWorthyCaseSummary{
		"b": {ScenarioID: "b", PromotionCount: 1},
		"a": {ScenarioID: "a", PromotionCount: 1},
		"c": {ScenarioID: "c", PromotionCount: 3},
	})
	if len(replayCases) != 3 || replayCases[0].ScenarioID != "c" || replayCases[1].ScenarioID != "a" {
		t.Fatalf("unexpected replay candidates %#v", replayCases)
	}

	issueSummaries := sortIssueSummaries(map[string]*ScenarioIssueSummary{
		"b": {ScenarioID: "b", Count: 1},
		"a": {ScenarioID: "a", Count: 1},
		"c": {ScenarioID: "c", Count: 3},
	})
	if len(issueSummaries) != 3 || issueSummaries[0].ScenarioID != "c" || issueSummaries[1].ScenarioID != "a" {
		t.Fatalf("unexpected issue summaries %#v", issueSummaries)
	}
}

func TestGeneratedReportWritersHandleFallbacks(t *testing.T) {
	reportDir := filepath.Join(t.TempDir(), "reports")
	generated, err := writeGeneratedReport(reportDir, []reportArtifact{
		{name: ArtifactDryRunReportJSON, kind: artifactKindJSON, payload: map[string]string{"status": "ok"}},
		{name: ArtifactDryRunReportMD, kind: artifactKindMarkdown, body: "# Dry Run"},
	})
	if err != nil {
		t.Fatalf("writeGeneratedReport() error = %v", err)
	}
	if len(generated.Artifacts) != 2 {
		t.Fatalf("expected 2 artifacts, got %#v", generated.Artifacts)
	}

	window := ReportWindow{
		Start: time.Date(2026, 3, 25, 12, 30, 0, 0, time.UTC),
		End:   time.Date(2026, 3, 26, 12, 30, 0, 0, time.UTC),
	}
	if got := safeWindowLabel(window); got != "20260325T123000Z-to-20260326T123000Z" {
		t.Fatalf("unexpected fallback label %q", got)
	}

	bridge := defaultProductionBridgeSummary(nil)
	if bridge.EvaluationMode != "shadow" || bridge.BlockingMode != "recommendation_only" {
		t.Fatalf("unexpected default production bridge summary %#v", bridge)
	}

	if !slicesContain([]string{"tier_c_used"}, "tier_c_used") || slicesContain([]string{"tier_b_only"}, "tier_c_used") {
		t.Fatal("unexpected slicesContain result")
	}
}

func fixtureRound(id string, completedAt time.Time, outcome benchmark.RobustnessOutcome) benchmark.BenchmarkRound {
	return benchmark.BenchmarkRound{
		ID:             id,
		ScenarioFamily: "commerce",
		CompletedAt:    completedAt,
		PromotionResults: []benchmark.PromotionDecision{{
			ID:              "promo-" + id,
			RoundID:         id,
			ScenarioID:      "commerce-v1-weakened-provenance",
			PromotionReason: benchmark.PromotionReasonDetectorMiss,
			Rationale:       "Suspicious challenger behavior evaluated as clean.",
			Promoted:        true,
		}},
		ScenarioResults: []benchmark.ScenarioResult{{
			ID:                   "result-" + id,
			ScenarioID:           "commerce-v1-weakened-provenance",
			SetKind:              benchmark.ScenarioSetLiving,
			FinalDetectionStatus: "clean",
			Passed:               false,
			Notes:                []string{"challenger variant exposed a detector weakness"},
		}},
		Recommendations: []benchmark.Recommendation{{
			ID:                "rec-" + id,
			Type:              benchmark.RecommendationTypeAddToReplayStableSet,
			LinkedRoundID:     id,
			LinkedScenarioIDs: []string{"commerce-v1-weakened-provenance"},
			SuggestedAction:   "Add the promoted case into replay coverage.",
		}},
		Summary: benchmark.RoundSummary{
			RoundID:             id,
			ScenarioFamily:      "commerce",
			PromotionCount:      1,
			ReplayPassRate:      0.5,
			RobustnessOutcome:   outcome,
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Run beside the incumbent fraud stack.",
			RecommendedFollowUp: "Review replay promotions.",
		},
	}
}
