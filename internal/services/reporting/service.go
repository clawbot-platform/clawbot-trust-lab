package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
	domainscenario "clawbot-trust-lab/internal/domain/scenario"
)

type Service struct {
	baseDir      string
	now          func() time.Time
	scenarioByID map[string]domainscenario.Scenario
}

const (
	ArtifactRoundSummaryJSON   = "round-summary.json"
	ArtifactRoundSummaryMD     = "round-summary.md"
	ArtifactRoundReportJSON    = "round-report.json"
	ArtifactRoundReportMD      = "round-report.md"
	ArtifactDetectionDeltaJSON = "detection-delta.json"
	ArtifactPromotionReport    = "promotion-report.json"
	ArtifactRecommendationJSON = "recommendation-report.json"
	ArtifactExecutiveSummaryMD = "executive-summary.md"
	ArtifactDryRunReportJSON   = "dry-run-report.json"
	ArtifactDryRunReportMD     = "dry-run-report.md"
	ArtifactManagementJSON     = "management-report.json"
	ArtifactManagementMD       = "management-report.md"
	reportSectionDaily         = "daily"
	reportSectionManagement    = "management"
	artifactKindJSON           = "json"
	artifactKindMarkdown       = "markdown"
)

type reportArtifact struct {
	name    string
	kind    string
	payload any
	body    string
}

type ScenarioCatalog interface {
	ListScenarios() []domainscenario.Scenario
}

func NewService(baseDir string, catalogs ...ScenarioCatalog) *Service {
	service := &Service{
		baseDir:      baseDir,
		now:          func() time.Time { return time.Now().UTC() },
		scenarioByID: map[string]domainscenario.Scenario{},
	}
	if len(catalogs) > 0 && catalogs[0] != nil {
		for _, item := range catalogs[0].ListScenarios() {
			service.scenarioByID[item.ID] = item
		}
	}
	return service
}

func (s *Service) Generate(round benchmark.BenchmarkRound) (benchmark.ReportIndex, error) {
	reportDir := filepath.Join(s.baseDir, round.ID)
	if err := os.MkdirAll(reportDir, 0o750); err != nil {
		return benchmark.ReportIndex{}, fmt.Errorf("create report dir %s: %w", reportDir, err)
	}

	index := s.roundReportIndex(reportDir)
	round.ReportDir = reportDir
	round.Reports = index
	roundReport := s.BuildRoundReport(round)

	for _, item := range s.reportArtifacts(round, roundReport) {
		if err := writeArtifact(filepath.Join(reportDir, item.name), item); err != nil {
			return benchmark.ReportIndex{}, err
		}
	}

	return index, nil
}

func (s *Service) BuildRoundReport(round benchmark.BenchmarkRound) RoundReport {
	return RoundReport{
		ReportType:             "round_report",
		GeneratedAt:            s.now(),
		RoundID:                round.ID,
		ScenarioFamily:         round.ScenarioFamily,
		StartedAt:              round.StartedAt,
		CompletedAt:            round.CompletedAt,
		ScenariosExecuted:      len(round.ScenarioResults),
		Summary:                round.Summary,
		RecommendationCounts:   recommendationCounts(round.Recommendations),
		PromotionCounts:        promotionCounts(round.PromotionResults),
		Promotions:             append([]benchmark.PromotionDecision(nil), round.PromotionResults...),
		Recommendations:        append([]benchmark.Recommendation(nil), round.Recommendations...),
		Regressions:            regressions(round.Delta),
		NotableChallengerCases: append([]benchmark.PromotionDecision(nil), round.PromotionResults...),
		TierUsage:              s.tierUsage(round.ScenarioResults),
		ProductionBridgeSummary: ProductionBridgeSummary{
			EvaluationMode:                 round.Summary.EvaluationMode,
			BlockingMode:                   round.Summary.BlockingMode,
			ExistingControlIntegrationNote: round.Summary.ExistingControlNote,
			RecommendedFollowUp:            round.Summary.RecommendedFollowUp,
		},
	}
}

func (s *Service) roundReportIndex(reportDir string) benchmark.ReportIndex {
	return benchmark.ReportIndex{
		RoundID:   filepath.Base(reportDir),
		Directory: reportDir,
		Artifacts: []benchmark.ReportArtifact{
			{Name: ArtifactRoundSummaryJSON, Path: filepath.Join(reportDir, ArtifactRoundSummaryJSON), Kind: artifactKindJSON},
			{Name: ArtifactRoundSummaryMD, Path: filepath.Join(reportDir, ArtifactRoundSummaryMD), Kind: artifactKindMarkdown},
			{Name: ArtifactRoundReportJSON, Path: filepath.Join(reportDir, ArtifactRoundReportJSON), Kind: artifactKindJSON},
			{Name: ArtifactRoundReportMD, Path: filepath.Join(reportDir, ArtifactRoundReportMD), Kind: artifactKindMarkdown},
			{Name: ArtifactDetectionDeltaJSON, Path: filepath.Join(reportDir, ArtifactDetectionDeltaJSON), Kind: artifactKindJSON},
			{Name: ArtifactPromotionReport, Path: filepath.Join(reportDir, ArtifactPromotionReport), Kind: artifactKindJSON},
			{Name: ArtifactRecommendationJSON, Path: filepath.Join(reportDir, ArtifactRecommendationJSON), Kind: artifactKindJSON},
			{Name: ArtifactExecutiveSummaryMD, Path: filepath.Join(reportDir, ArtifactExecutiveSummaryMD), Kind: artifactKindMarkdown},
		},
	}
}

func (s *Service) reportArtifacts(round benchmark.BenchmarkRound, roundReport RoundReport) []reportArtifact {
	return []reportArtifact{
		{name: ArtifactRoundSummaryJSON, kind: artifactKindJSON, payload: round},
		{name: ArtifactRoundSummaryMD, kind: artifactKindMarkdown, body: s.roundSummaryMarkdown(round)},
		{name: ArtifactRoundReportJSON, kind: artifactKindJSON, payload: roundReport},
		{name: ArtifactRoundReportMD, kind: artifactKindMarkdown, body: s.roundReportMarkdown(roundReport)},
		{name: ArtifactDetectionDeltaJSON, kind: artifactKindJSON, payload: round.Delta},
		{name: ArtifactPromotionReport, kind: artifactKindJSON, payload: round.PromotionResults},
		{name: ArtifactRecommendationJSON, kind: artifactKindJSON, payload: BuildRecommendationReport(round)},
		{name: ArtifactExecutiveSummaryMD, kind: artifactKindMarkdown, body: s.executiveSummary(round)},
	}
}

func BuildRecommendationReport(round benchmark.BenchmarkRound) benchmark.RecommendationReport {
	return benchmark.RecommendationReport{
		RoundID:                        round.ID,
		EvaluationMode:                 round.Summary.EvaluationMode,
		BlockingMode:                   round.Summary.BlockingMode,
		ExistingControlIntegrationNote: round.Summary.ExistingControlNote,
		RecommendedFollowUp:            round.Summary.RecommendedFollowUp,
		Recommendations:                append([]benchmark.Recommendation(nil), round.Recommendations...),
	}
}

func BackfillRecommendationReport(reportDir string, round benchmark.BenchmarkRound) (bool, error) {
	path := filepath.Join(reportDir, ArtifactRecommendationJSON)
	if _, err := os.Stat(path); err == nil {
		return false, nil
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("stat report %s: %w", path, err)
	}

	if err := writeJSON(path, BuildRecommendationReport(round)); err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) GenerateDryRunReport(window ReportWindow, rounds []benchmark.BenchmarkRound, health OperationalHealthSummary) (GeneratedReport, error) {
	reportDir := filepath.Join(s.baseDir, reportSectionDaily, safeWindowLabel(window))
	filtered := filterRoundsForWindow(rounds, window)
	report := s.buildDryRunReport(window, filtered, health)
	artifacts := []reportArtifact{
		{name: ArtifactDryRunReportJSON, kind: artifactKindJSON, payload: report},
		{name: ArtifactDryRunReportMD, kind: artifactKindMarkdown, body: s.dryRunReportMarkdown(report)},
	}
	return writeGeneratedReport(reportDir, artifacts)
}

func (s *Service) GenerateManagementReport(window ReportWindow, rounds []benchmark.BenchmarkRound, health OperationalHealthSummary) (GeneratedReport, error) {
	reportDir := filepath.Join(s.baseDir, reportSectionManagement, safeWindowLabel(window))
	filtered := filterRoundsForWindow(rounds, window)
	report := s.buildManagementReport(window, filtered, health)
	artifacts := []reportArtifact{
		{name: ArtifactManagementJSON, kind: artifactKindJSON, payload: report},
		{name: ArtifactManagementMD, kind: artifactKindMarkdown, body: s.managementReportMarkdown(report)},
	}
	return writeGeneratedReport(reportDir, artifacts)
}

func (s *Service) roundSummaryMarkdown(round benchmark.BenchmarkRound) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "# Round Summary\n\n")
	fmt.Fprintf(&builder, "- Round ID: `%s`\n", round.ID)
	fmt.Fprintf(&builder, "- Scenario family: `%s`\n", round.ScenarioFamily)
	fmt.Fprintf(&builder, "- Stable scenarios: `%d`\n", round.Summary.StableScenarioCount)
	fmt.Fprintf(&builder, "- Challenger variants: `%d`\n", round.Summary.ChallengerCount)
	fmt.Fprintf(&builder, "- Replay retests: `%d`\n", round.Summary.ReplayRetestCount)
	fmt.Fprintf(&builder, "- Promotions: `%d`\n", round.Summary.PromotionCount)
	fmt.Fprintf(&builder, "- Replay pass rate: `%.2f`\n", round.Summary.ReplayPassRate)
	fmt.Fprintf(&builder, "- Robustness outcome: `%s`\n\n", round.Summary.RobustnessOutcome)
	fmt.Fprintf(&builder, "- Evaluation mode: `%s`\n", round.Summary.EvaluationMode)
	fmt.Fprintf(&builder, "- Blocking mode: `%s`\n\n", round.Summary.BlockingMode)
	fmt.Fprintf(&builder, "Production-bridge note: %s\n\n", round.Summary.ExistingControlNote)

	fmt.Fprintf(&builder, "## Important Findings\n\n")
	if len(round.Summary.ImportantFindings) == 0 {
		fmt.Fprintf(&builder, "- No notable findings were recorded.\n")
	} else {
		for _, finding := range round.Summary.ImportantFindings {
			fmt.Fprintf(&builder, "- %s\n", finding)
		}
	}

	fmt.Fprintf(&builder, "\n## Promoted Cases\n\n")
	if len(round.PromotionResults) == 0 {
		fmt.Fprintf(&builder, "- No challenger cases were promoted in this round.\n")
	} else {
		for _, item := range round.PromotionResults {
			fmt.Fprintf(&builder, "- `%s`: %s\n", item.ScenarioID, item.Rationale)
		}
	}

	fmt.Fprintf(&builder, "\n## Recommendations\n\n")
	if len(round.Recommendations) == 0 {
		fmt.Fprintf(&builder, "- No explicit recommendations were generated.\n")
	} else {
		for _, item := range round.Recommendations {
			fmt.Fprintf(&builder, "- `%s` (`%s`): %s\n", item.Type, item.Priority, item.SuggestedAction)
			fmt.Fprintf(&builder, "  Rationale: %s\n", item.Rationale)
			if len(item.LinkedScenarioIDs) > 0 {
				fmt.Fprintf(&builder, "  Linked scenarios: `%s`\n", strings.Join(item.LinkedScenarioIDs, "`, `"))
			}
			if len(item.LinkedPromotionIDs) > 0 {
				fmt.Fprintf(&builder, "  Linked promotions: `%s`\n", strings.Join(item.LinkedPromotionIDs, "`, `"))
			}
		}
	}

	fmt.Fprintf(&builder, "\nRecommended follow-up: %s\n", round.Summary.RecommendedFollowUp)

	return builder.String()
}

func (s *Service) roundReportMarkdown(report RoundReport) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "# DRQ Round Report\n\n")
	fmt.Fprintf(&builder, "- Round ID: `%s`\n", report.RoundID)
	fmt.Fprintf(&builder, "- Scenario family: `%s`\n", report.ScenarioFamily)
	fmt.Fprintf(&builder, "- Generated at: `%s`\n", report.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&builder, "- Scenarios executed: `%d`\n", report.ScenariosExecuted)
	fmt.Fprintf(&builder, "- Promotions: `%d`\n", len(report.Promotions))
	fmt.Fprintf(&builder, "- Recommendations: `%d`\n", len(report.Recommendations))
	fmt.Fprintf(&builder, "- Regressions: `%d`\n\n", len(report.Regressions))

	fmt.Fprintf(&builder, "## Production Bridge Summary\n\n")
	fmt.Fprintf(&builder, "- Evaluation mode: `%s`\n", report.ProductionBridgeSummary.EvaluationMode)
	fmt.Fprintf(&builder, "- Blocking mode: `%s`\n", report.ProductionBridgeSummary.BlockingMode)
	fmt.Fprintf(&builder, "- Existing control integration note: %s\n", report.ProductionBridgeSummary.ExistingControlIntegrationNote)
	fmt.Fprintf(&builder, "- Recommended follow-up: %s\n", report.ProductionBridgeSummary.RecommendedFollowUp)

	fmt.Fprintf(&builder, "\n## Tier Usage\n\n")
	fmt.Fprintf(&builder, "- Tier A available results: `%d`\n", report.TierUsage.TierAAvailableCount)
	fmt.Fprintf(&builder, "- Tier B available results: `%d`\n", report.TierUsage.TierBAvailableCount)
	fmt.Fprintf(&builder, "- Tier C capable results: `%d`\n", report.TierUsage.TierCCapableCount)
	fmt.Fprintf(&builder, "- Tier C used results: `%d`\n", report.TierUsage.TierCUsedCount)
	fmt.Fprintf(&builder, "- Interpretation: %s\n", report.TierUsage.InterpretationNote)

	fmt.Fprintf(&builder, "\n## Notable Challenger Cases\n\n")
	if len(report.NotableChallengerCases) == 0 {
		fmt.Fprintf(&builder, "- No challenger cases were promoted in this round.\n")
	} else {
		for _, item := range report.NotableChallengerCases {
			fmt.Fprintf(&builder, "- `%s` (`%s`): %s\n", item.ScenarioID, item.PromotionReason, item.Rationale)
		}
	}

	fmt.Fprintf(&builder, "\n## Recommendation Themes\n\n")
	if len(report.RecommendationCounts) == 0 {
		fmt.Fprintf(&builder, "- No explicit recommendations were generated.\n")
	} else {
		for _, line := range sortedRecommendationCountLines(report.RecommendationCounts) {
			fmt.Fprintf(&builder, "- %s\n", line)
		}
	}

	fmt.Fprintf(&builder, "\n## Regressions\n\n")
	if len(report.Regressions) == 0 {
		fmt.Fprintf(&builder, "- No regressions were recorded for this round.\n")
	} else {
		for _, item := range report.Regressions {
			fmt.Fprintf(&builder, "- `%s`: `%s` -> `%s` (score delta `%d`)\n", item.ScenarioID, item.PreviousStatus, item.CurrentStatus, item.ScoreDelta)
		}
	}

	return builder.String()
}

func (s *Service) executiveSummary(round benchmark.BenchmarkRound) string {
	headline := "Mixed results"
	switch round.Summary.RobustnessOutcome {
	case benchmark.RobustnessOutcomeImproved:
		headline = "Robustness improved"
	case benchmark.RobustnessOutcomeRegressed:
		headline = "Regression observed"
	case benchmark.RobustnessOutcomeNewBlindSpotDiscovered:
		headline = "New blind spot discovered"
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "# Executive Summary\n\n")
	fmt.Fprintf(&builder, "Outcome: **%s**\n\n", headline)
	fmt.Fprintf(
		&builder,
		"Round `%s` evaluated %d stable scenarios, %d challenger variants, and %d replay retests.\n\n",
		round.ID,
		round.Summary.StableScenarioCount,
		round.Summary.ChallengerCount,
		round.Summary.ReplayRetestCount,
	)

	if len(round.Summary.ImportantFindings) > 0 {
		fmt.Fprintf(&builder, "Key findings:\n")
		for _, finding := range round.Summary.ImportantFindings {
			fmt.Fprintf(&builder, "- %s\n", finding)
		}
	} else {
		fmt.Fprintf(&builder, "Key findings:\n- No material findings in this round.\n")
	}

	fmt.Fprintf(
		&builder,
		"\nOperating posture: `%s` / `%s`.\n\nRecommended next action: %s\n",
		round.Summary.EvaluationMode,
		round.Summary.BlockingMode,
		round.Summary.RecommendedFollowUp,
	)
	fmt.Fprintf(&builder, "\nExisting control integration note: %s\n", round.Summary.ExistingControlNote)

	return builder.String()
}

func (s *Service) buildDryRunReport(window ReportWindow, rounds []benchmark.BenchmarkRound, health OperationalHealthSummary) DryRunReport {
	report := DryRunReport{
		ReportType:              "dry_run_report",
		GeneratedAt:             s.now(),
		Window:                  window,
		RecommendationCounts:    map[benchmark.RecommendationType]int{},
		RobustnessOutcomeCounts: map[benchmark.RobustnessOutcome]int{},
		OperationalHealth:       health,
		ProductionBridgeSummary: defaultProductionBridgeSummary(rounds),
	}

	patternCounts := map[string]*ScenarioIssueSummary{}
	replayCandidates := map[string]*ReplayWorthyCaseSummary{}
	recExamples := map[benchmark.RecommendationType]string{}

	for _, round := range rounds {
		report.RoundIDs = append(report.RoundIDs, round.ID)
		report.TotalRounds++
		report.TotalPromotions += len(round.PromotionResults)
		report.TotalRecommendations += len(round.Recommendations)
		report.RobustnessOutcomeCounts[round.Summary.RobustnessOutcome]++
		if round.Summary.RobustnessOutcome == benchmark.RobustnessOutcomeNewBlindSpotDiscovered {
			report.NewBlindSpotsDiscovered++
		}
		if round.Summary.RobustnessOutcome == benchmark.RobustnessOutcomeRegressed {
			report.RegressionsObserved++
		}

		for _, rec := range round.Recommendations {
			report.RecommendationCounts[rec.Type]++
			if recExamples[rec.Type] == "" {
				recExamples[rec.Type] = rec.SuggestedAction
			}
		}
		for _, promo := range round.PromotionResults {
			entry, ok := replayCandidates[promo.ScenarioID]
			if !ok {
				entry = &ReplayWorthyCaseSummary{
					ScenarioID:    promo.ScenarioID,
					LatestRoundID: round.ID,
					Rationale:     promo.Rationale,
				}
				replayCandidates[promo.ScenarioID] = entry
			}
			entry.PromotionCount++
			entry.PromotionReasons = appendUniquePromotionReason(entry.PromotionReasons, promo.PromotionReason)
			entry.LatestRoundID = round.ID
			entry.Rationale = promo.Rationale
		}
		for _, result := range round.ScenarioResults {
			if result.Passed {
				continue
			}
			entry, ok := patternCounts[result.ScenarioID]
			if !ok {
				entry = &ScenarioIssueSummary{ScenarioID: result.ScenarioID}
				patternCounts[result.ScenarioID] = entry
			}
			entry.Count++
			for _, note := range result.Notes {
				if note != "" {
					entry.Reasons = appendUniqueString(entry.Reasons, note)
				}
			}
		}
	}

	report.RecurringRecommendationThemes = recommendationThemes(report.RecommendationCounts, recExamples)
	report.NewReplayWorthyCases = sortReplayCandidates(replayCandidates)
	report.RecurringIssueScenarios = sortIssueSummaries(patternCounts)
	report.NotableFindings = dryRunFindings(report, rounds)

	return report
}

func (s *Service) dryRunReportMarkdown(report DryRunReport) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "# 24-Hour DRQ Dry-Run Report\n\n")
	fmt.Fprintf(&builder, "- Window: `%s` to `%s`\n", report.Window.Start.Format(time.RFC3339), report.Window.End.Format(time.RFC3339))
	fmt.Fprintf(&builder, "- Generated at: `%s`\n", report.GeneratedAt.Format(time.RFC3339))
	fmt.Fprintf(&builder, "- Total rounds: `%d`\n", report.TotalRounds)
	fmt.Fprintf(&builder, "- Total promotions: `%d`\n", report.TotalPromotions)
	fmt.Fprintf(&builder, "- Total recommendations: `%d`\n", report.TotalRecommendations)
	fmt.Fprintf(&builder, "- New blind spots: `%d`\n", report.NewBlindSpotsDiscovered)
	fmt.Fprintf(&builder, "- Regressions observed: `%d`\n\n", report.RegressionsObserved)

	fmt.Fprintf(&builder, "## Recurring Recommendation Themes\n\n")
	if len(report.RecurringRecommendationThemes) == 0 {
		fmt.Fprintf(&builder, "- No recurring recommendation themes were recorded in this window.\n")
	} else {
		for _, item := range report.RecurringRecommendationThemes {
			fmt.Fprintf(&builder, "- `%s` x `%d`: %s\n", item.Type, item.Count, item.ExampleAction)
		}
	}

	fmt.Fprintf(&builder, "\n## New Replay-Worthy Cases\n\n")
	if len(report.NewReplayWorthyCases) == 0 {
		fmt.Fprintf(&builder, "- No new replay-worthy cases were promoted during this window.\n")
	} else {
		for _, item := range report.NewReplayWorthyCases {
			fmt.Fprintf(&builder, "- `%s` promoted `%d` time(s); latest round `%s`.\n", item.ScenarioID, item.PromotionCount, item.LatestRoundID)
		}
	}

	fmt.Fprintf(&builder, "\n## Operational Stability\n\n")
	fmt.Fprintf(&builder, "- Trust Lab status at report generation: `%s`\n", report.OperationalHealth.TrustLabStatus)
	fmt.Fprintf(&builder, "- Control plane status at report generation: `%s`\n", report.OperationalHealth.ControlPlaneStatus)
	fmt.Fprintf(&builder, "- Memory status at report generation: `%s`\n", report.OperationalHealth.MemoryStatus)
	fmt.Fprintf(&builder, "- Health history available: `%t`\n", report.OperationalHealth.HealthHistoryAvailable)
	fmt.Fprintf(&builder, "- Note: %s\n", report.OperationalHealth.Note)

	fmt.Fprintf(&builder, "\n## Dry-Run Assessment\n\n")
	if len(report.NotableFindings) == 0 {
		fmt.Fprintf(&builder, "- No material dry-run findings were recorded for this window.\n")
	} else {
		for _, item := range report.NotableFindings {
			fmt.Fprintf(&builder, "- %s\n", item)
		}
	}

	return builder.String()
}

func (s *Service) buildManagementReport(window ReportWindow, rounds []benchmark.BenchmarkRound, health OperationalHealthSummary) ManagementReport {
	dryRun := s.buildDryRunReport(window, rounds, health)
	report := ManagementReport{
		ReportType:               "management_report",
		GeneratedAt:              s.now(),
		Window:                   window,
		TotalRounds:              len(rounds),
		ConsistentIssueScenarios: dryRun.RecurringIssueScenarios,
		ReplayBaselineCandidates: dryRun.NewReplayWorthyCases,
		OperationalHealth:        health,
	}

	report.DRQValueFindings = managementValueFindings(dryRun)
	report.RecommendedNextProductionStep = managementNextSteps(dryRun)
	report.ExecutiveSummary = managementExecutiveSummary(dryRun)
	report.StakeholderNotes = managementStakeholderNotes(dryRun)
	return report
}

func (s *Service) managementReportMarkdown(report ManagementReport) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "# 1-Week DRQ Management Report\n\n")
	fmt.Fprintf(&builder, "%s\n\n", report.ExecutiveSummary)

	fmt.Fprintf(&builder, "## What DRQ Found Beyond The Baseline\n\n")
	if len(report.DRQValueFindings) == 0 {
		fmt.Fprintf(&builder, "- No material DRQ-specific findings were recorded in this window.\n")
	} else {
		for _, item := range report.DRQValueFindings {
			fmt.Fprintf(&builder, "- %s\n", item)
		}
	}

	fmt.Fprintf(&builder, "\n## Scenarios That Repeatedly Surfaced Issues\n\n")
	if len(report.ConsistentIssueScenarios) == 0 {
		fmt.Fprintf(&builder, "- No scenario repeatedly surfaced issues across this window.\n")
	} else {
		for _, item := range report.ConsistentIssueScenarios {
			fmt.Fprintf(&builder, "- `%s` in `%d` round(s).\n", item.ScenarioID, item.Count)
		}
	}

	fmt.Fprintf(&builder, "\n## Replay Cases Worth Long-Lived Baseline Promotion\n\n")
	if len(report.ReplayBaselineCandidates) == 0 {
		fmt.Fprintf(&builder, "- No replay candidates were strong enough to recommend for long-lived baseline coverage in this window.\n")
	} else {
		for _, item := range report.ReplayBaselineCandidates {
			fmt.Fprintf(&builder, "- `%s` promoted `%d` time(s); latest rationale: %s\n", item.ScenarioID, item.PromotionCount, item.Rationale)
		}
	}

	fmt.Fprintf(&builder, "\n## Operational Stability\n\n")
	fmt.Fprintf(&builder, "- Trust Lab status at report generation: `%s`\n", report.OperationalHealth.TrustLabStatus)
	fmt.Fprintf(&builder, "- Control plane status at report generation: `%s`\n", report.OperationalHealth.ControlPlaneStatus)
	fmt.Fprintf(&builder, "- Memory status at report generation: `%s`\n", report.OperationalHealth.MemoryStatus)
	fmt.Fprintf(&builder, "- Note: %s\n", report.OperationalHealth.Note)

	fmt.Fprintf(&builder, "\n## Recommended Next Production Steps\n\n")
	if len(report.RecommendedNextProductionStep) == 0 {
		fmt.Fprintf(&builder, "- Continue shadow-mode monitoring.\n")
	} else {
		for _, item := range report.RecommendedNextProductionStep {
			fmt.Fprintf(&builder, "- %s\n", item)
		}
	}

	return builder.String()
}

func writeJSON(path string, payload any) error {
	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal report %s: %w", path, err)
	}

	body = append(body, '\n')

	if err := os.WriteFile(path, body, 0o600); err != nil {
		return fmt.Errorf("write report %s: %w", path, err)
	}

	return nil
}

func writeArtifact(path string, item reportArtifact) error {
	switch item.kind {
	case artifactKindJSON:
		return writeJSON(path, item.payload)
	default:
		if err := os.WriteFile(path, []byte(item.body), 0o600); err != nil {
			return fmt.Errorf("write report %s: %w", path, err)
		}
		return nil
	}
}

func writeGeneratedReport(reportDir string, artifacts []reportArtifact) (GeneratedReport, error) {
	if err := os.MkdirAll(reportDir, 0o750); err != nil {
		return GeneratedReport{}, fmt.Errorf("create report dir %s: %w", reportDir, err)
	}

	out := GeneratedReport{Directory: reportDir}
	for _, item := range artifacts {
		path := filepath.Join(reportDir, item.name)
		if err := writeArtifact(path, item); err != nil {
			return GeneratedReport{}, err
		}
		out.Artifacts = append(out.Artifacts, benchmark.ReportArtifact{
			Name: item.name,
			Path: path,
			Kind: item.kind,
		})
	}
	return out, nil
}

func filterRoundsForWindow(rounds []benchmark.BenchmarkRound, window ReportWindow) []benchmark.BenchmarkRound {
	items := make([]benchmark.BenchmarkRound, 0)
	for _, round := range rounds {
		occurredAt := round.CompletedAt
		if occurredAt.IsZero() {
			occurredAt = round.StartedAt
		}
		if occurredAt.IsZero() {
			continue
		}
		if occurredAt.Before(window.Start) || occurredAt.After(window.End) {
			continue
		}
		items = append(items, round)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].CompletedAt.Before(items[j].CompletedAt)
	})
	return items
}

func safeWindowLabel(window ReportWindow) string {
	label := strings.TrimSpace(window.Label)
	if label == "" {
		label = window.Start.UTC().Format("20060102T150405Z") + "-to-" + window.End.UTC().Format("20060102T150405Z")
	}
	replacer := strings.NewReplacer(":", "", "/", "-", " ", "_")
	return replacer.Replace(label)
}

func recommendationCounts(items []benchmark.Recommendation) map[benchmark.RecommendationType]int {
	counts := map[benchmark.RecommendationType]int{}
	for _, item := range items {
		counts[item.Type]++
	}
	return counts
}

func promotionCounts(items []benchmark.PromotionDecision) map[benchmark.PromotionReason]int {
	counts := map[benchmark.PromotionReason]int{}
	for _, item := range items {
		counts[item.PromotionReason]++
	}
	return counts
}

func regressions(items []benchmark.DetectionDelta) []benchmark.DetectionDelta {
	out := make([]benchmark.DetectionDelta, 0)
	for _, item := range items {
		if item.ScoreDelta < 0 || item.RecommendationChanged {
			out = append(out, item)
		}
	}
	return out
}

func sortedRecommendationCountLines(counts map[benchmark.RecommendationType]int) []string {
	lines := make([]string, 0, len(counts))
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("`%s` x `%d`", key, counts[benchmark.RecommendationType(key)]))
	}
	return lines
}

func recommendationThemes(counts map[benchmark.RecommendationType]int, examples map[benchmark.RecommendationType]string) []RecommendationTheme {
	items := make([]RecommendationTheme, 0, len(counts))
	for key, count := range counts {
		items = append(items, RecommendationTheme{
			Type:          key,
			Count:         count,
			ExampleAction: examples[key],
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Type < items[j].Type
		}
		return items[i].Count > items[j].Count
	})
	return items
}

func sortReplayCandidates(items map[string]*ReplayWorthyCaseSummary) []ReplayWorthyCaseSummary {
	out := make([]ReplayWorthyCaseSummary, 0, len(items))
	for _, item := range items {
		out = append(out, *item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].PromotionCount == out[j].PromotionCount {
			return out[i].ScenarioID < out[j].ScenarioID
		}
		return out[i].PromotionCount > out[j].PromotionCount
	})
	return out
}

func sortIssueSummaries(items map[string]*ScenarioIssueSummary) []ScenarioIssueSummary {
	out := make([]ScenarioIssueSummary, 0, len(items))
	for _, item := range items {
		out = append(out, *item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count == out[j].Count {
			return out[i].ScenarioID < out[j].ScenarioID
		}
		return out[i].Count > out[j].Count
	})
	return out
}

func dryRunFindings(report DryRunReport, rounds []benchmark.BenchmarkRound) []string {
	findings := make([]string, 0)
	if report.TotalRounds == 0 {
		return []string{"No benchmark rounds completed inside the selected window."}
	}
	if report.TotalPromotions > 0 {
		findings = append(findings, fmt.Sprintf("%d promotion(s) created %d replay-worthy case candidate(s) during the selected dry-run window.", report.TotalPromotions, len(report.NewReplayWorthyCases)))
	}
	if report.RegressionsObserved > 0 {
		findings = append(findings, fmt.Sprintf("%d round(s) reported regressions in replay or rule posture.", report.RegressionsObserved))
	}
	if report.NewBlindSpotsDiscovered > 0 {
		findings = append(findings, fmt.Sprintf("%d round(s) surfaced new blind spots that the baseline stable suite did not already cover.", report.NewBlindSpotsDiscovered))
	}
	if len(findings) == 0 {
		findings = append(findings, fmt.Sprintf("%d round(s) completed without new promotions or regressions in this window.", len(rounds)))
	}
	return findings
}

func managementValueFindings(dryRun DryRunReport) []string {
	findings := make([]string, 0)
	if dryRun.TotalPromotions > 0 {
		findings = append(findings, fmt.Sprintf("DRQ promoted %d challenger outcome(s) into replay-worthy review items that the stable baseline alone would not have surfaced.", dryRun.TotalPromotions))
	}
	if dryRun.NewBlindSpotsDiscovered > 0 {
		findings = append(findings, fmt.Sprintf("%d round(s) exposed new blind spots, showing value beyond static baseline-only coverage.", dryRun.NewBlindSpotsDiscovered))
	}
	if len(dryRun.RecurringIssueScenarios) > 0 {
		top := dryRun.RecurringIssueScenarios[0]
		findings = append(findings, fmt.Sprintf("Scenario `%s` repeatedly surfaced issues across %d round(s), indicating a persistent control gap worth operational attention.", top.ScenarioID, top.Count))
	}
	return findings
}

func managementExecutiveSummary(dryRun DryRunReport) string {
	if dryRun.TotalRounds == 0 {
		return "No benchmark rounds completed inside the selected management-report window, so there is not yet enough operational evidence for a broader DRQ assessment."
	}
	if dryRun.TotalPromotions == 0 && dryRun.RegressionsObserved == 0 {
		return fmt.Sprintf("Across %d completed round(s), the Version 1 DRQ harness remained operationally stable and did not surface material new replay promotions or regressions, which supports continued shadow-mode monitoring before broader expansion.", dryRun.TotalRounds)
	}
	return fmt.Sprintf("Across %d completed round(s), the Version 1 DRQ harness surfaced %d promotion(s), %d replay/baseline regression signal(s), and %d new blind-spot round(s), providing concrete evidence for where incumbent fraud controls should be reviewed in shadow mode before any production policy change.", dryRun.TotalRounds, dryRun.TotalPromotions, dryRun.RegressionsObserved, dryRun.NewBlindSpotsDiscovered)
}

func managementNextSteps(dryRun DryRunReport) []string {
	steps := make([]string, 0)
	for _, theme := range dryRun.RecurringRecommendationThemes {
		if theme.ExampleAction != "" {
			steps = append(steps, theme.ExampleAction)
		}
	}
	if len(dryRun.NewReplayWorthyCases) > 0 {
		steps = append(steps, fmt.Sprintf("Promote %d replay-worthy case(s) into longer-lived replay coverage before changing incumbent controls.", len(dryRun.NewReplayWorthyCases)))
	}
	if len(steps) == 0 {
		steps = append(steps, "Continue recommendation-only shadow operation and collect a longer evidence window before making control changes.")
	}
	return dedupeStrings(steps)
}

func managementStakeholderNotes(dryRun DryRunReport) []string {
	return []string{
		"Version 1 remains a shadow-mode, recommendation-only harness and is not a replacement for the incumbent fraud stack.",
		"Version 1 does not persist degraded-period history yet; operational stability in this report reflects current health plus round outcomes, not a full incident timeline.",
		fmt.Sprintf("Replay and recommendation totals are derived from %d round(s) in the selected window.", dryRun.TotalRounds),
	}
}

func defaultProductionBridgeSummary(rounds []benchmark.BenchmarkRound) ProductionBridgeSummary {
	for _, round := range rounds {
		if round.Summary.EvaluationMode != "" || round.Summary.BlockingMode != "" || round.Summary.ExistingControlNote != "" || round.Summary.RecommendedFollowUp != "" {
			return ProductionBridgeSummary{
				EvaluationMode:                 round.Summary.EvaluationMode,
				BlockingMode:                   round.Summary.BlockingMode,
				ExistingControlIntegrationNote: round.Summary.ExistingControlNote,
				RecommendedFollowUp:            round.Summary.RecommendedFollowUp,
			}
		}
	}
	return ProductionBridgeSummary{
		EvaluationMode:                 "shadow",
		BlockingMode:                   "recommendation_only",
		ExistingControlIntegrationNote: "Run this harness beside the incumbent fraud stack and compare its recommendations before making policy changes.",
		RecommendedFollowUp:            "Keep operating in shadow mode while promoting the strongest replay cases into longer-lived review coverage.",
	}
}

func (s *Service) tierUsage(results []benchmark.ScenarioResult) TierUsageSummary {
	summary := TierUsageSummary{
		ResultsEvaluated:   len(results),
		TierCOptional:      true,
		InterpretationNote: "Tier usage is derived from the persisted scenario catalog plus round-level Tier C usage markers. Version 1 reports Tier C capability and observed use, not a separate Tier A/Tier B decision score.",
	}
	for _, item := range results {
		scenarioItem, ok := s.scenarioByID[item.ScenarioID]
		if !ok {
			summary.UnknownScenarioIDs = appendUniqueString(summary.UnknownScenarioIDs, item.ScenarioID)
			continue
		}
		if len(scenarioItem.FeatureModel.TierA) > 0 {
			summary.TierAAvailableCount++
		}
		if len(scenarioItem.FeatureModel.TierB) > 0 {
			summary.TierBAvailableCount++
		}
		if len(scenarioItem.FeatureModel.TierC) > 0 {
			summary.TierCCapableCount++
		}
		if slicesContain(item.Notes, "tier_c_used") {
			summary.TierCUsedCount++
		}
	}
	return summary
}

func appendUniquePromotionReason(items []benchmark.PromotionReason, value benchmark.PromotionReason) []benchmark.PromotionReason {
	for _, item := range items {
		if item == value {
			return items
		}
	}
	return append(items, value)
}

func appendUniqueString(items []string, value string) []string {
	for _, item := range items {
		if item == value {
			return items
		}
	}
	return append(items, value)
}

func dedupeStrings(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		out = appendUniqueString(out, item)
	}
	return out
}

func slicesContain(items []string, needle string) bool {
	for _, item := range items {
		if item == needle {
			return true
		}
	}
	return false
}
