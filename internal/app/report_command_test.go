package app

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"clawbot-trust-lab/internal/clients/controlplane"
	"clawbot-trust-lab/internal/clients/memory"
	"clawbot-trust-lab/internal/config"
	"clawbot-trust-lab/internal/domain/benchmark"
	"clawbot-trust-lab/internal/platform/bootstrap"
	"clawbot-trust-lab/internal/platform/store"
	servicebenchmark "clawbot-trust-lab/internal/services/benchmark"
	"clawbot-trust-lab/internal/services/reporting"
)

func TestParseReportWindowSupportsLastDuration(t *testing.T) {
	now := time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC)
	window, err := parseReportWindow(now, "24h", "", "")
	if err != nil {
		t.Fatalf("parseReportWindow() error = %v", err)
	}
	if got, want := window.Start, now.Add(-24*time.Hour); !got.Equal(want) {
		t.Fatalf("expected start %s, got %s", want, got)
	}
	if window.Label != "2026-03-25_to_2026-03-26" {
		t.Fatalf("unexpected label %q", window.Label)
	}
}

func TestParseReportWindowSupportsExplicitRange(t *testing.T) {
	now := time.Date(2026, 3, 26, 15, 0, 0, 0, time.UTC)
	window, err := parseReportWindow(now, "", "2026-03-19T00:00:00Z", "2026-03-26T00:00:00Z")
	if err != nil {
		t.Fatalf("parseReportWindow() error = %v", err)
	}
	if window.Label != "2026-03-19_to_2026-03-26" {
		t.Fatalf("unexpected label %q", window.Label)
	}
}

func TestParseReportWindowRejectsHalfRange(t *testing.T) {
	if _, err := parseReportWindow(time.Now().UTC(), "", "2026-03-19T00:00:00Z", ""); err == nil {
		t.Fatal("expected parseReportWindow() to reject a half-specified range")
	}
}

func TestDispatchReportCommandRejectsMissingType(t *testing.T) {
	err := dispatchReportCommand(&bytes.Buffer{}, bootstrap.Dependencies{}, reporting.OperationalHealthSummary{}, nil)
	if err == nil {
		t.Fatal("expected dispatchReportCommand() to reject missing args")
	}
}

func TestRunReportCommandUsesBuildDependencies(t *testing.T) {
	original := buildReportDependencies
	t.Cleanup(func() { buildReportDependencies = original })

	deps := reportCommandDeps(t, []benchmark.BenchmarkRound{{
		ID:             "round-1",
		ScenarioFamily: "commerce",
		CompletedAt:    time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC),
		Summary: benchmark.RoundSummary{
			RoundID:             "round-1",
			ScenarioFamily:      "commerce",
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Run beside incumbent controls.",
			RecommendedFollowUp: "Review promoted cases.",
		},
	}})
	deps.ControlPlane = controlPlaneHealthStub{}
	deps.Memory = memoryHealthStub{}

	buildReportDependencies = func(config.Config, *slog.Logger) (bootstrap.Dependencies, error) {
		return deps, nil
	}

	var buf bytes.Buffer
	if err := RunReportCommand(context.Background(), config.Config{}, slog.Default(), &buf, []string{"round", "--round-id", "round-1"}); err != nil {
		t.Fatalf("RunReportCommand() error = %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("generated round report")) {
		t.Fatalf("expected round report output, got %s", buf.String())
	}
}

func TestRunRoundReportWritesArtifactPaths(t *testing.T) {
	deps := reportCommandDeps(t, []benchmark.BenchmarkRound{{
		ID:             "round-1",
		ScenarioFamily: "commerce",
		CompletedAt:    time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC),
		Summary: benchmark.RoundSummary{
			RoundID:             "round-1",
			ScenarioFamily:      "commerce",
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Run beside incumbent controls.",
			RecommendedFollowUp: "Review promoted cases.",
		},
	}})

	var buf bytes.Buffer
	if err := runRoundReport(&buf, deps, []string{"--round-id", "round-1"}); err != nil {
		t.Fatalf("runRoundReport() error = %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("round-report.json")) || !bytes.Contains(buf.Bytes(), []byte("round-report.md")) {
		t.Fatalf("expected round report artifact paths in output, got %s", buf.String())
	}
}

func TestRunWindowReportWritesArtifacts(t *testing.T) {
	deps := reportCommandDeps(t, []benchmark.BenchmarkRound{{
		ID:             "round-1",
		ScenarioFamily: "commerce",
		CompletedAt:    time.Date(2026, 3, 26, 12, 0, 0, 0, time.UTC),
		PromotionResults: []benchmark.PromotionDecision{{
			ID:              "promo-1",
			RoundID:         "round-1",
			ScenarioID:      "commerce-v1-weakened-provenance",
			PromotionReason: benchmark.PromotionReasonDetectorMiss,
			Rationale:       "challenger missed",
			Promoted:        true,
		}},
		ScenarioResults: []benchmark.ScenarioResult{{
			ID:         "result-1",
			ScenarioID: "commerce-v1-weakened-provenance",
			Passed:     false,
			Notes:      []string{"challenger variant exposed a detector weakness"},
		}},
		Recommendations: []benchmark.Recommendation{{
			ID:              "rec-1",
			Type:            benchmark.RecommendationTypeAddToReplayStableSet,
			LinkedRoundID:   "round-1",
			SuggestedAction: "Add the promoted case into replay coverage.",
		}},
		Summary: benchmark.RoundSummary{
			RoundID:             "round-1",
			ScenarioFamily:      "commerce",
			PromotionCount:      1,
			RobustnessOutcome:   benchmark.RobustnessOutcomeNewBlindSpotDiscovered,
			EvaluationMode:      "shadow",
			BlockingMode:        "recommendation_only",
			ExistingControlNote: "Run beside incumbent controls.",
			RecommendedFollowUp: "Review replay promotions.",
		},
	}})

	var buf bytes.Buffer
	health := reporting.OperationalHealthSummary{
		TrustLabStatus:     "ok",
		ControlPlaneStatus: "ok",
		MemoryStatus:       "ok",
		Note:               "snapshot only",
	}

	if err := runWindowReport(&buf, deps, health, "dry-run", 24*time.Hour, []string{"--from", "2026-03-25T00:00:00Z", "--to", "2026-03-26T23:59:59Z"}); err != nil {
		t.Fatalf("runWindowReport(dry-run) error = %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("dry-run-report.json")) {
		t.Fatalf("expected dry-run artifact output, got %s", buf.String())
	}

	buf.Reset()
	if err := runWindowReport(&buf, deps, health, "management", 7*24*time.Hour, []string{"--last", "168h"}); err != nil {
		t.Fatalf("runWindowReport(management) error = %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("management-report.json")) {
		t.Fatalf("expected management artifact output, got %s", buf.String())
	}
}

func TestDispatchReportCommandRejectsUnsupportedType(t *testing.T) {
	err := dispatchReportCommand(&bytes.Buffer{}, bootstrap.Dependencies{}, reporting.OperationalHealthSummary{}, []string{"weekly"})
	if err == nil {
		t.Fatal("expected unsupported report type error")
	}
}

func TestRunRoundReportRejectsMissingRoundID(t *testing.T) {
	err := runRoundReport(&bytes.Buffer{}, reportCommandDeps(t, nil), nil)
	if err == nil {
		t.Fatal("expected --round-id validation error")
	}
}

func TestRunWindowReportRejectsInvalidWindow(t *testing.T) {
	err := runWindowReport(&bytes.Buffer{}, reportCommandDeps(t, nil), reporting.OperationalHealthSummary{}, "dry-run", 24*time.Hour, []string{"--last", "not-a-duration"})
	if err == nil {
		t.Fatal("expected invalid --last error")
	}
}

func TestReportWindowLabelSameDay(t *testing.T) {
	start := time.Date(2026, 3, 26, 1, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 26, 23, 0, 0, 0, time.UTC)
	if got := reportWindowLabel(start, end); got != "2026-03-26" {
		t.Fatalf("unexpected same-day label %q", got)
	}
}

func TestReportHealthSummaryReflectsDependencyFailures(t *testing.T) {
	deps := bootstrap.Dependencies{
		ControlPlane: controlPlaneHealthStub{err: errors.New("down")},
		Memory:       memoryHealthStub{err: errors.New("down")},
	}

	summary := reportHealthSummary(context.Background(), deps)
	if summary.ControlPlaneStatus != "degraded" || summary.MemoryStatus != "degraded" {
		t.Fatalf("expected degraded dependency snapshot, got %#v", summary)
	}
	if len(summary.DegradedPeriods) != 2 {
		t.Fatalf("expected degraded notes, got %#v", summary.DegradedPeriods)
	}
}

type controlPlaneHealthStub struct {
	err error
}

func (s controlPlaneHealthStub) Health(context.Context) error { return s.err }
func (s controlPlaneHealthStub) ListRuns(context.Context) ([]controlplane.RunRef, error) {
	return nil, nil
}
func (s controlPlaneHealthStub) CreateRun(context.Context, controlplane.CreateRunRequest) (controlplane.RunRef, error) {
	return controlplane.RunRef{}, nil
}
func (s controlPlaneHealthStub) ListPolicies(context.Context) ([]controlplane.PolicyRef, error) {
	return nil, nil
}
func (s controlPlaneHealthStub) CreatePolicy(context.Context, controlplane.CreatePolicyRequest) (controlplane.PolicyRef, error) {
	return controlplane.PolicyRef{}, nil
}
func (s controlPlaneHealthStub) RegisterBenchmarkMetadata(context.Context, benchmark.RegistrationRequest) (benchmark.RegistrationResult, error) {
	return benchmark.RegistrationResult{}, nil
}

type memoryHealthStub struct {
	err error
}

func (s memoryHealthStub) Health(context.Context) error { return s.err }
func (s memoryHealthStub) StoreReplayCase(context.Context, memory.StoreReplayCaseRequest) error {
	return nil
}
func (s memoryHealthStub) FetchSimilarCases(context.Context, memory.FetchSimilarCasesRequest) (memory.FetchSimilarCasesResponse, error) {
	return memory.FetchSimilarCasesResponse{}, nil
}
func (s memoryHealthStub) StoreTrustArtifact(context.Context, memory.StoreTrustArtifactRequest) error {
	return nil
}
func (s memoryHealthStub) LoadScenarioContext(context.Context, memory.LoadScenarioContextRequest) (memory.LoadScenarioContextResponse, error) {
	return memory.LoadScenarioContextResponse{}, nil
}

func reportCommandDeps(t *testing.T, rounds []benchmark.BenchmarkRound) bootstrap.Dependencies {
	t.Helper()
	reporter := reporting.NewService(t.TempDir())
	benchmarkStore := store.NewBenchmarkStore()
	for _, round := range rounds {
		benchmarkStore.Put(round)
	}
	benchmarkService := servicebenchmark.NewService(nil, nil, nil, nil, benchmarkStore, reporter)
	return bootstrap.Dependencies{
		Benchmark: benchmarkService,
		Reporting: reporter,
	}
}
