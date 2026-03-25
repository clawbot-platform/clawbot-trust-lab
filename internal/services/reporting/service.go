package reporting

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"clawbot-trust-lab/internal/domain/benchmark"
)

type Service struct {
	baseDir string
}

func NewService(baseDir string) *Service {
	return &Service{baseDir: baseDir}
}

func (s *Service) Generate(round benchmark.BenchmarkRound) (benchmark.ReportIndex, error) {
	reportDir := filepath.Join(s.baseDir, round.ID)
	if err := os.MkdirAll(reportDir, 0o750); err != nil {
		return benchmark.ReportIndex{}, fmt.Errorf("create report dir %s: %w", reportDir, err)
	}

	type artifact struct {
		name    string
		kind    string
		payload any
		body    string
	}

	executive := s.executiveSummary(round)
	summaryMD := s.roundSummaryMarkdown(round)

	artifacts := []artifact{
		{name: "round-summary.json", kind: "json", payload: round},
		{name: "round-summary.md", kind: "markdown", body: summaryMD},
		{name: "detection-delta.json", kind: "json", payload: round.Delta},
		{name: "promotion-report.json", kind: "json", payload: round.PromotionResults},
		{name: "executive-summary.md", kind: "markdown", body: executive},
	}

	index := benchmark.ReportIndex{
		RoundID:   round.ID,
		Directory: reportDir,
	}

	for _, item := range artifacts {
		path := filepath.Join(reportDir, item.name)

		switch item.kind {
		case "json":
			if err := writeJSON(path, item.payload); err != nil {
				return benchmark.ReportIndex{}, err
			}
		default:
			if err := os.WriteFile(path, []byte(item.body), 0o600); err != nil {
				return benchmark.ReportIndex{}, fmt.Errorf("write report %s: %w", path, err)
			}
		}

		index.Artifacts = append(index.Artifacts, benchmark.ReportArtifact{
			Name: item.name,
			Path: path,
			Kind: item.kind,
		})
	}

	return index, nil
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
		"\nRecommended next action: review %d promoted cases and compare the new delta report against the previous round.\n",
		len(round.PromotionResults),
	)

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
