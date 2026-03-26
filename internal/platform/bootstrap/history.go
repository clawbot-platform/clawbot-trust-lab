package bootstrap

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	servicebenchmark "clawbot-trust-lab/internal/services/benchmark"
	servicereporting "clawbot-trust-lab/internal/services/reporting"
)

type HistoricalState struct {
	Rounds           []benchmark.BenchmarkRound
	DetectionResults []detectionmodel.DetectionResult
}

func LoadHistoricalState(reportsDir string, logger *slog.Logger) HistoricalState {
	if logger == nil {
		logger = slog.Default()
	}

	entries, err := os.ReadDir(reportsDir)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("historical round bootstrap skipped because reports directory does not exist", "reports_dir", reportsDir)
			return HistoricalState{}
		}
		logger.Warn("historical round bootstrap failed to scan reports directory", "reports_dir", reportsDir, "error", err)
		return HistoricalState{}
	}

	state := HistoricalState{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		roundID := entry.Name()
		roundDir := filepath.Join(reportsDir, roundID)

		round, ok, err := loadHistoricalRound(reportsDir, roundID, logger)
		if err != nil {
			logger.Warn("skipping malformed historical round directory", "round_dir", roundDir, "error", err)
			continue
		}
		if !ok {
			continue
		}

		state.Rounds = append(state.Rounds, round)
		state.DetectionResults = append(state.DetectionResults, reconstructDetectionResults(round)...)
	}

	sort.Slice(state.Rounds, func(i, j int) bool {
		left := historicalRoundSortKey(state.Rounds[i])
		right := historicalRoundSortKey(state.Rounds[j])
		if left.Equal(right) {
			return state.Rounds[i].ID > state.Rounds[j].ID
		}
		return left.After(right)
	})

	return state
}

func loadHistoricalRound(reportsDir, roundID string, logger *slog.Logger) (benchmark.BenchmarkRound, bool, error) {
	roundDir := filepath.Join(reportsDir, roundID)
	summaryRelPath := filepath.Join(roundID, servicereporting.ArtifactRoundSummaryJSON)

	if _, err := os.Stat(filepath.Join(reportsDir, summaryRelPath)); err != nil {
		if os.IsNotExist(err) {
			return benchmark.BenchmarkRound{}, false, nil
		}
		return benchmark.BenchmarkRound{}, false, fmt.Errorf("stat %s: %w", servicereporting.ArtifactRoundSummaryJSON, err)
	}

	var round benchmark.BenchmarkRound
	if err := readJSON(reportsDir, summaryRelPath, &round); err != nil {
		return benchmark.BenchmarkRound{}, false, fmt.Errorf("read %s: %w", servicereporting.ArtifactRoundSummaryJSON, err)
	}

	normalizeHistoricalRound(&round, roundID, roundDir)
	if err := loadHistoricalPromotions(reportsDir, roundID, &round); err != nil {
		return benchmark.BenchmarkRound{}, false, err
	}
	if err := loadHistoricalDelta(reportsDir, roundID, &round); err != nil {
		return benchmark.BenchmarkRound{}, false, err
	}
	ensureHistoricalSummaryCounts(&round)
	if err := loadHistoricalRecommendationReport(reportsDir, roundID, roundDir, &round, logger); err != nil {
		return benchmark.BenchmarkRound{}, false, err
	}

	round.Reports = listReportArtifacts(round.ID, roundDir)
	if round.Reports.RoundID == "" {
		round.Reports.RoundID = round.ID
	}

	return round, true, nil
}

func normalizeHistoricalRound(round *benchmark.BenchmarkRound, roundID, roundDir string) {
	if round == nil {
		return
	}
	if round.ID == "" {
		round.ID = roundID
	}
	round.ReportDir = roundDir
	round.Summary.RoundID = round.ID
	if round.Summary.ScenarioFamily == "" {
		round.Summary.ScenarioFamily = round.ScenarioFamily
	}
	if round.ScenarioFamily == "" {
		round.ScenarioFamily = round.Summary.ScenarioFamily
	}
	if round.Reports.RoundID == "" {
		round.Reports.RoundID = round.ID
	}
}

func loadHistoricalPromotions(reportsDir, roundID string, round *benchmark.BenchmarkRound) error {
	promotionsRelPath := filepath.Join(roundID, servicereporting.ArtifactPromotionReport)
	var promotions []benchmark.PromotionDecision
	if ok, err := readOptionalJSON(reportsDir, promotionsRelPath, &promotions); err != nil {
		return fmt.Errorf("read %s: %w", servicereporting.ArtifactPromotionReport, err)
	} else if ok {
		round.PromotionResults = promotions
	}
	for idx := range round.PromotionResults {
		if round.PromotionResults[idx].RoundID == "" {
			round.PromotionResults[idx].RoundID = round.ID
		}
	}
	return nil
}

func loadHistoricalDelta(reportsDir, roundID string, round *benchmark.BenchmarkRound) error {
	deltaRelPath := filepath.Join(roundID, servicereporting.ArtifactDetectionDeltaJSON)
	var delta []benchmark.DetectionDelta
	if ok, err := readOptionalJSON(reportsDir, deltaRelPath, &delta); err != nil {
		return fmt.Errorf("read %s: %w", servicereporting.ArtifactDetectionDeltaJSON, err)
	} else if ok {
		round.Delta = delta
	}
	return nil
}

func ensureHistoricalSummaryCounts(round *benchmark.BenchmarkRound) {
	if round.Summary.PromotionCount == 0 && len(round.PromotionResults) > 0 {
		round.Summary.PromotionCount = len(round.PromotionResults)
	}
	if round.Summary.StableScenarioCount == 0 && len(round.StableScenarioRefs) > 0 {
		round.Summary.StableScenarioCount = len(round.StableScenarioRefs)
	}
	if round.Summary.ChallengerCount == 0 && len(round.ChallengerVariantRefs) > 0 {
		round.Summary.ChallengerCount = len(round.ChallengerVariantRefs)
	}
}

func loadHistoricalRecommendationReport(reportsDir, roundID, roundDir string, round *benchmark.BenchmarkRound, logger *slog.Logger) error {
	reportRelPath := filepath.Join(roundID, servicereporting.ArtifactRecommendationJSON)
	report, format, ok, err := loadRecommendationReport(reportsDir, reportRelPath, round)
	if err != nil {
		return fmt.Errorf("read %s: %w", servicereporting.ArtifactRecommendationJSON, err)
	}
	if ok {
		applyRecommendationReport(round, report)
		if logger != nil {
			logger.Info("loaded historical recommendation report", "round_id", round.ID, "report_path", filepath.Join(roundDir, servicereporting.ArtifactRecommendationJSON), "format", format)
		}
		return nil
	}

	servicebenchmark.EnsureProductionBridgeSummary(round)
	if written, err := servicereporting.BackfillRecommendationReport(roundDir, *round); err != nil {
		if logger != nil {
			logger.Warn("historical recommendation report backfill failed", "round_id", round.ID, "round_dir", roundDir, "error", err)
		}
	} else if written && logger != nil {
		logger.Info("backfilled historical recommendation report", "round_id", round.ID, "round_dir", roundDir)
	}
	return nil
}

func loadRecommendationReport(rootDir, relPath string, round *benchmark.BenchmarkRound) (benchmark.RecommendationReport, string, bool, error) {
	body, ok, err := readOptionalJSONBytes(rootDir, relPath)
	if err != nil || !ok {
		return benchmark.RecommendationReport{}, "", ok, err
	}

	report, format, err := parseRecommendationReport(body, round)
	if err != nil {
		return benchmark.RecommendationReport{}, "", true, fmt.Errorf("unmarshal %s: %w", filepath.Clean(strings.TrimSpace(relPath)), err)
	}
	return report, format, true, nil
}

func parseRecommendationReport(body []byte, round *benchmark.BenchmarkRound) (benchmark.RecommendationReport, string, error) {
	var report benchmark.RecommendationReport
	if err := json.Unmarshal(body, &report); err == nil {
		return report, "current", nil
	}

	var recommendations []benchmark.Recommendation
	if err := json.Unmarshal(body, &recommendations); err == nil {
		return convertLegacyRecommendations(round, recommendations), "legacy_array", nil
	}

	return benchmark.RecommendationReport{}, "", fmt.Errorf("unsupported recommendation report format")
}

func convertLegacyRecommendations(round *benchmark.BenchmarkRound, recommendations []benchmark.Recommendation) benchmark.RecommendationReport {
	report := benchmark.RecommendationReport{
		Recommendations: append([]benchmark.Recommendation(nil), recommendations...),
	}
	if round == nil {
		return report
	}

	summaryRound := *round
	servicebenchmark.EnsureProductionBridgeSummary(&summaryRound)
	report.RoundID = summaryRound.ID
	report.EvaluationMode = summaryRound.Summary.EvaluationMode
	report.BlockingMode = summaryRound.Summary.BlockingMode
	report.ExistingControlIntegrationNote = summaryRound.Summary.ExistingControlNote
	report.RecommendedFollowUp = summaryRound.Summary.RecommendedFollowUp
	return report
}

func applyRecommendationReport(round *benchmark.BenchmarkRound, report benchmark.RecommendationReport) {
	if round == nil {
		return
	}
	if len(round.Recommendations) == 0 && len(report.Recommendations) > 0 {
		round.Recommendations = append([]benchmark.Recommendation(nil), report.Recommendations...)
	}
	if round.Summary.EvaluationMode == "" {
		round.Summary.EvaluationMode = report.EvaluationMode
	}
	if round.Summary.BlockingMode == "" {
		round.Summary.BlockingMode = report.BlockingMode
	}
	if round.Summary.ExistingControlNote == "" {
		round.Summary.ExistingControlNote = report.ExistingControlIntegrationNote
	}
	if round.Summary.RecommendedFollowUp == "" {
		round.Summary.RecommendedFollowUp = report.RecommendedFollowUp
	}
	if round.Summary.Recommendations == 0 && len(round.Recommendations) > 0 {
		round.Summary.Recommendations = len(round.Recommendations)
	}
}

func readJSON(rootDir, relPath string, dest any) error {
	body, err := readJSONBytes(rootDir, relPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("unmarshal %s: %w", filepath.Clean(strings.TrimSpace(relPath)), err)
	}

	return nil
}

func readOptionalJSON(rootDir, relPath string, dest any) (bool, error) {
	body, ok, err := readOptionalJSONBytes(rootDir, relPath)
	if err != nil || !ok {
		return ok, err
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return true, fmt.Errorf("unmarshal %s: %w", filepath.Clean(strings.TrimSpace(relPath)), err)
	}
	return true, nil
}

func readOptionalJSONBytes(rootDir, relPath string) ([]byte, bool, error) {
	fullPath := filepath.Join(rootDir, relPath)
	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	body, err := readJSONBytes(rootDir, relPath)
	if err != nil {
		return nil, true, err
	}
	return body, true, nil
}

func readJSONBytes(rootDir, relPath string) ([]byte, error) {
	clean := filepath.Clean(strings.TrimSpace(relPath))
	if clean == "." || clean == "" {
		return nil, fmt.Errorf("relative path is required")
	}
	if filepath.IsAbs(clean) {
		return nil, fmt.Errorf("absolute paths are not allowed: %q", relPath)
	}
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("path escapes root: %q", relPath)
	}

	root, err := os.OpenRoot(rootDir)
	if err != nil {
		return nil, fmt.Errorf("open bootstrap root %s: %w", rootDir, err)
	}
	defer func() { _ = root.Close() }()

	f, err := root.Open(clean)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", clean, err)
	}
	defer func() { _ = f.Close() }()

	body, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", clean, err)
	}

	return body, nil
}

func listReportArtifacts(roundID string, roundDir string) benchmark.ReportIndex {
	index := benchmark.ReportIndex{
		RoundID:   roundID,
		Directory: roundDir,
	}

	entries, err := os.ReadDir(roundDir)
	if err != nil {
		return index
	}

	for _, entry := range entries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		index.Artifacts = append(index.Artifacts, benchmark.ReportArtifact{
			Name: entry.Name(),
			Path: filepath.Join(roundDir, entry.Name()),
			Kind: reportKind(entry),
		})
	}

	sort.Slice(index.Artifacts, func(i, j int) bool {
		return index.Artifacts[i].Name < index.Artifacts[j].Name
	})

	return index
}

func reportKind(entry fs.DirEntry) string {
	switch strings.ToLower(filepath.Ext(entry.Name())) {
	case ".md":
		return "markdown"
	case ".json":
		return "json"
	default:
		return "file"
	}
}

func reconstructDetectionResults(round benchmark.BenchmarkRound) []detectionmodel.DetectionResult {
	results := make([]detectionmodel.DetectionResult, 0, len(round.ScenarioResults))
	for _, item := range round.ScenarioResults {
		if item.DetectionResultRef == "" {
			continue
		}

		triggered := make([]detectionmodel.RuleHit, 0, len(item.TriggeredRuleIDs))
		for _, ruleID := range item.TriggeredRuleIDs {
			triggered = append(triggered, detectionmodel.RuleHit{
				RuleID:   ruleID,
				Title:    strings.ReplaceAll(ruleID, "_", " "),
				Severity: historicalSeverity(item.FinalDetectionStatus),
				Reason:   "reconstructed from historical round artifact",
			})
		}

		results = append(results, detectionmodel.DetectionResult{
			ID:                item.DetectionResultRef,
			ScenarioID:        item.ScenarioID,
			OrderID:           firstRef(item.OrderRefs),
			RefundID:          firstRef(item.RefundRefs),
			TrustDecisionRefs: append([]string(nil), item.TrustDecisionRefs...),
			ReplayCaseRefs:    append([]string(nil), item.ReplayCaseRefs...),
			Status:            item.FinalDetectionStatus,
			Score:             historicalScore(item.FinalDetectionStatus),
			Grade:             historicalGrade(item.FinalDetectionStatus),
			TriggeredRules:    triggered,
			ReasonCodes:       append([]string(nil), item.TriggeredRuleIDs...),
			Recommendation:    item.FinalRecommendation,
			EvaluatedAt:       historicalRoundSortKey(round),
			Metadata: map[string]any{
				"historical_reconstruction": true,
				"round_id":                  round.ID,
				"scenario_result_ref":       item.ID,
			},
		})
	}
	return results
}

func historicalRoundSortKey(round benchmark.BenchmarkRound) time.Time {
	if !round.CompletedAt.IsZero() {
		return round.CompletedAt
	}
	if !round.StartedAt.IsZero() {
		return round.StartedAt
	}
	return time.Time{}
}

func historicalScore(status detectionmodel.DetectionStatus) int {
	switch status {
	case detectionmodel.DetectionStatusBlocked:
		return 80
	case detectionmodel.DetectionStatusStepUpRequired:
		return 40
	case detectionmodel.DetectionStatusSuspicious:
		return 15
	default:
		return 0
	}
}

func historicalGrade(status detectionmodel.DetectionStatus) detectionmodel.RiskGrade {
	switch status {
	case detectionmodel.DetectionStatusBlocked:
		return detectionmodel.RiskGradeCritical
	case detectionmodel.DetectionStatusStepUpRequired:
		return detectionmodel.RiskGradeHigh
	case detectionmodel.DetectionStatusSuspicious:
		return detectionmodel.RiskGradeModerate
	default:
		return detectionmodel.RiskGradeLow
	}
}

func historicalSeverity(status detectionmodel.DetectionStatus) int {
	switch status {
	case detectionmodel.DetectionStatusBlocked:
		return 30
	case detectionmodel.DetectionStatusStepUpRequired:
		return 20
	case detectionmodel.DetectionStatusSuspicious:
		return 10
	default:
		return 0
	}
}

func firstRef(items []string) string {
	if len(items) == 0 {
		return ""
	}
	return items[0]
}
