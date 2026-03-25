package store

import (
	"testing"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
)

func TestBenchmarkStoreMergesHistoricalAndLiveRounds(t *testing.T) {
	store := NewBenchmarkStore()
	historical := benchmark.BenchmarkRound{
		ID:          "round-1",
		CompletedAt: time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
		ReportDir:   "./reports/round-1",
		Reports: benchmark.ReportIndex{
			RoundID:   "round-1",
			Directory: "./reports/round-1",
			Artifacts: []benchmark.ReportArtifact{{Name: "round-summary.json", Path: "./reports/round-1/round-summary.json", Kind: "json"}},
		},
	}
	live := benchmark.BenchmarkRound{
		ID:          "round-1",
		CompletedAt: time.Date(2026, 3, 25, 12, 5, 0, 0, time.UTC),
		Summary:     benchmark.RoundSummary{RoundID: "round-1", PromotionCount: 1},
	}

	store.PutHistorical(historical)
	store.Put(live)

	items := store.List()
	if len(items) != 1 {
		t.Fatalf("expected 1 merged round, got %d", len(items))
	}
	if items[0].Summary.PromotionCount != 1 {
		t.Fatalf("expected live summary to win, got %#v", items[0].Summary)
	}
	if items[0].ReportDir != historical.ReportDir {
		t.Fatalf("expected historical report dir to be preserved, got %s", items[0].ReportDir)
	}
	if len(items[0].Reports.Artifacts) != 1 {
		t.Fatalf("expected historical report artifacts to be preserved, got %#v", items[0].Reports.Artifacts)
	}
}

func TestBenchmarkStoreListsNewestRoundFirst(t *testing.T) {
	store := NewBenchmarkStore()
	store.PutHistorical(benchmark.BenchmarkRound{
		ID:          "round-older",
		CompletedAt: time.Date(2026, 3, 25, 12, 0, 0, 0, time.UTC),
	})
	store.Put(benchmark.BenchmarkRound{
		ID:          "round-newer",
		CompletedAt: time.Date(2026, 3, 25, 12, 10, 0, 0, time.UTC),
	})

	items := store.List()
	if len(items) != 2 {
		t.Fatalf("expected 2 rounds, got %d", len(items))
	}
	if items[0].ID != "round-newer" || items[1].ID != "round-older" {
		t.Fatalf("unexpected order: %#v", items)
	}

	latest, ok := store.Latest()
	if !ok || latest.ID != "round-newer" {
		t.Fatalf("expected latest round-newer, got %#v", latest)
	}
}
