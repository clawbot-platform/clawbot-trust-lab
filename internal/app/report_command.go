package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"clawbot-trust-lab/internal/config"
	"clawbot-trust-lab/internal/platform/bootstrap"
	"clawbot-trust-lab/internal/services/reporting"
)

var buildReportDependencies = bootstrap.Build

func RunReportCommand(ctx context.Context, cfg config.Config, logger *slog.Logger, writer io.Writer, args []string) error {
	deps, err := buildReportDependencies(cfg, logger)
	if err != nil {
		return err
	}

	health := reportHealthSummary(ctx, deps)
	return dispatchReportCommand(writer, deps, health, args)
}

func dispatchReportCommand(writer io.Writer, deps bootstrap.Dependencies, health reporting.OperationalHealthSummary, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing report type: expected round, dry-run, or management")
	}
	switch args[0] {
	case "round":
		return runRoundReport(writer, deps, args[1:])
	case "dry-run":
		return runWindowReport(writer, deps, health, "dry-run", 24*time.Hour, args[1:])
	case "management":
		return runWindowReport(writer, deps, health, "management", 7*24*time.Hour, args[1:])
	default:
		return fmt.Errorf("unsupported report type %q", args[0])
	}
}

func runRoundReport(writer io.Writer, deps bootstrap.Dependencies, args []string) error {
	fs := flag.NewFlagSet("round", flag.ContinueOnError)
	roundID := fs.String("round-id", "", "round id to render")
	fs.SetOutput(io.Discard)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*roundID) == "" {
		return fmt.Errorf("--round-id is required")
	}

	round, err := deps.Benchmark.GetRound(strings.TrimSpace(*roundID))
	if err != nil {
		return err
	}
	index, err := deps.Reporting.Generate(round)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(writer, "generated round report for %s in %s\n", round.ID, index.Directory)
	for _, item := range index.Artifacts {
		if item.Name == reporting.ArtifactRoundReportJSON || item.Name == reporting.ArtifactRoundReportMD {
			_, _ = fmt.Fprintf(writer, "- %s\n", item.Path)
		}
	}
	return nil
}

func runWindowReport(writer io.Writer, deps bootstrap.Dependencies, health reporting.OperationalHealthSummary, kind string, defaultLast time.Duration, args []string) error {
	fs := flag.NewFlagSet(kind, flag.ContinueOnError)
	last := fs.String("last", defaultLast.String(), "duration window, for example 24h or 168h")
	from := fs.String("from", "", "window start in RFC3339")
	to := fs.String("to", "", "window end in RFC3339")
	fs.SetOutput(io.Discard)
	if err := fs.Parse(args); err != nil {
		return err
	}

	window, err := parseReportWindow(time.Now().UTC(), *last, *from, *to)
	if err != nil {
		return err
	}

	var generated reporting.GeneratedReport
	switch kind {
	case "dry-run":
		generated, err = deps.Reporting.GenerateDryRunReport(window, deps.Benchmark.ListRounds(), health)
	case "management":
		generated, err = deps.Reporting.GenerateManagementReport(window, deps.Benchmark.ListRounds(), health)
	default:
		return fmt.Errorf("unsupported report kind %q", kind)
	}
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(writer, "generated %s report for %s to %s in %s\n", kind, window.Start.Format(time.RFC3339), window.End.Format(time.RFC3339), generated.Directory)
	for _, item := range generated.Artifacts {
		_, _ = fmt.Fprintf(writer, "- %s\n", item.Path)
	}
	return nil
}

func parseReportWindow(now time.Time, lastValue string, fromValue string, toValue string) (reporting.ReportWindow, error) {
	lastValue = strings.TrimSpace(lastValue)
	fromValue = strings.TrimSpace(fromValue)
	toValue = strings.TrimSpace(toValue)

	switch {
	case fromValue != "" || toValue != "":
		if fromValue == "" || toValue == "" {
			return reporting.ReportWindow{}, fmt.Errorf("--from and --to must be provided together")
		}
		start, err := time.Parse(time.RFC3339, fromValue)
		if err != nil {
			return reporting.ReportWindow{}, fmt.Errorf("parse --from: %w", err)
		}
		end, err := time.Parse(time.RFC3339, toValue)
		if err != nil {
			return reporting.ReportWindow{}, fmt.Errorf("parse --to: %w", err)
		}
		if end.Before(start) {
			return reporting.ReportWindow{}, fmt.Errorf("--to must be after --from")
		}
		return reporting.ReportWindow{
			Label:       reportWindowLabel(start.UTC(), end.UTC()),
			Start:       start.UTC(),
			End:         end.UTC(),
			GeneratedAt: now.UTC(),
		}, nil
	default:
		duration, err := time.ParseDuration(lastValue)
		if err != nil {
			return reporting.ReportWindow{}, fmt.Errorf("parse --last: %w", err)
		}
		start := now.UTC().Add(-duration)
		end := now.UTC()
		return reporting.ReportWindow{
			Label:       reportWindowLabel(start, end),
			Start:       start,
			End:         end,
			GeneratedAt: now.UTC(),
		}, nil
	}
}

func reportWindowLabel(start time.Time, end time.Time) string {
	if start.Format("2006-01-02") == end.Format("2006-01-02") {
		return end.Format("2006-01-02")
	}
	return start.Format("2006-01-02") + "_to_" + end.Format("2006-01-02")
}

func reportHealthSummary(ctx context.Context, deps bootstrap.Dependencies) reporting.OperationalHealthSummary {
	summary := reporting.OperationalHealthSummary{
		TrustLabStatus:         "ok",
		ControlPlaneStatus:     "ok",
		MemoryStatus:           "ok",
		HealthHistoryAvailable: false,
		Note:                   "Version 1 does not yet persist a degraded-period or recovery timeline. This report captures a generation-time health snapshot plus benchmark outcomes.",
	}
	if err := deps.ControlPlane.Health(ctx); err != nil {
		summary.ControlPlaneStatus = "degraded"
		summary.DegradedPeriods = append(summary.DegradedPeriods, "control-plane unreachable at report generation time")
	}
	if err := deps.Memory.Health(ctx); err != nil {
		summary.MemoryStatus = "degraded"
		summary.DegradedPeriods = append(summary.DegradedPeriods, "clawmem unreachable at report generation time")
	}
	return summary
}
