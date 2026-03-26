package operator

import (
	"os"
	"testing"
	"time"

	domainbenchmark "clawbot-trust-lab/internal/domain/benchmark"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
	"clawbot-trust-lab/internal/platform/store"
)

type benchmarkStub struct {
	rounds map[string]domainbenchmark.BenchmarkRound
}

func (s benchmarkStub) ListRounds() []domainbenchmark.BenchmarkRound {
	items := make([]domainbenchmark.BenchmarkRound, 0, len(s.rounds))
	for _, item := range s.rounds {
		items = append(items, item)
	}
	return items
}

func (s benchmarkStub) GetRound(id string) (domainbenchmark.BenchmarkRound, error) {
	item, ok := s.rounds[id]
	if !ok {
		return domainbenchmark.BenchmarkRound{}, errNotFound("round", id)
	}
	return item, nil
}

func (s benchmarkStub) GetRoundReports(id string) (domainbenchmark.ReportIndex, error) {
	item, err := s.GetRound(id)
	if err != nil {
		return domainbenchmark.ReportIndex{}, err
	}
	return item.Reports, nil
}

func (s benchmarkStub) ListRecommendations() []domainbenchmark.Recommendation {
	items := make([]domainbenchmark.Recommendation, 0)
	for _, round := range s.rounds {
		items = append(items, round.Recommendations...)
	}
	return items
}

func (s benchmarkStub) GetRecommendation(id string) (domainbenchmark.Recommendation, error) {
	for _, item := range s.ListRecommendations() {
		if item.ID == id {
			return item, nil
		}
	}
	return domainbenchmark.Recommendation{}, errNotFound("recommendation", id)
}

func (s benchmarkStub) LongRunSummary() domainbenchmark.LongRunSummary {
	return domainbenchmark.LongRunSummary{RoundsExecuted: len(s.rounds)}
}

type detectionStub struct {
	results map[string]detectionmodel.DetectionResult
}

func (s detectionStub) GetResult(id string) (detectionmodel.DetectionResult, error) {
	item, ok := s.results[id]
	if !ok {
		return detectionmodel.DetectionResult{}, errNotFound("detection", id)
	}
	return item, nil
}

func TestCompareRounds(t *testing.T) {
	service := NewService(benchmarkStub{rounds: operatorRoundsFixture(t)}, detectionStub{results: operatorDetectionsFixture()}, store.NewOperatorStore())

	comparison, err := service.CompareRounds("round-2", "round-1")
	if err != nil {
		t.Fatalf("CompareRounds() error = %v", err)
	}
	if comparison.PromotionsCountDelta != 1 {
		t.Fatalf("expected promotion delta 1, got %d", comparison.PromotionsCountDelta)
	}
	if comparison.DetectionDeltaCount != 1 {
		t.Fatalf("expected detection delta count 1, got %d", comparison.DetectionDeltaCount)
	}
}

func TestReviewPromotion(t *testing.T) {
	service := NewService(benchmarkStub{rounds: operatorRoundsFixture(t)}, detectionStub{results: operatorDetectionsFixture()}, store.NewOperatorStore())
	service.now = func() time.Time { return time.Date(2026, 3, 25, 13, 0, 0, 0, time.UTC) }

	review, err := service.ReviewPromotion("promo-2", ReviewInput{Status: "accepted", Note: "Promote into replay baseline."})
	if err != nil {
		t.Fatalf("ReviewPromotion() error = %v", err)
	}
	if review.Status != domainbenchmark.PromotionReviewAccepted {
		t.Fatalf("unexpected status %s", review.Status)
	}
	if review.Note == nil || review.Note.Body == "" {
		t.Fatal("expected note to be stored")
	}
}

func TestGetPromotionIncludesReview(t *testing.T) {
	reviews := store.NewOperatorStore()
	reviews.PutReview(domainbenchmark.PromotionReview{
		PromotionID: "promo-2",
		Status:      domainbenchmark.PromotionReviewNeedsFollowUp,
		UpdatedAt:   time.Date(2026, 3, 25, 13, 0, 0, 0, time.UTC),
	})

	service := NewService(benchmarkStub{rounds: operatorRoundsFixture(t)}, detectionStub{results: operatorDetectionsFixture()}, reviews)
	detail, err := service.GetPromotion("promo-2")
	if err != nil {
		t.Fatalf("GetPromotion() error = %v", err)
	}
	if detail.Review == nil || detail.Review.Status != domainbenchmark.PromotionReviewNeedsFollowUp {
		t.Fatalf("expected review on detail: %#v", detail)
	}
}

func TestListPromotionsIncludesHistoricalRoundsWithoutInventingReviewState(t *testing.T) {
	service := NewService(benchmarkStub{rounds: operatorRoundsFixture(t)}, detectionStub{results: operatorDetectionsFixture()}, store.NewOperatorStore())

	items := service.ListPromotions("")
	if len(items) != 1 {
		t.Fatalf("expected 1 promotion, got %d", len(items))
	}
	if items[0].Review != nil {
		t.Fatalf("expected historical promotion to have no reconstructed review, got %#v", items[0].Review)
	}
	if items[0].Promotion.ID != "promo-2" {
		t.Fatalf("unexpected promotion %#v", items[0].Promotion)
	}
}

func TestGetReportsAndArtifact(t *testing.T) {
	service := NewService(benchmarkStub{rounds: operatorRoundsFixture(t)}, detectionStub{results: operatorDetectionsFixture()}, store.NewOperatorStore())

	reports, err := service.GetReports("round-2")
	if err != nil {
		t.Fatalf("GetReports() error = %v", err)
	}
	if len(reports) != 1 || reports[0].ArtifactName != "executive-summary.md" {
		t.Fatalf("unexpected reports %#v", reports)
	}

	artifact, err := service.GetReportArtifact("round-2", "executive-summary.md")
	if err != nil {
		t.Fatalf("GetReportArtifact() error = %v", err)
	}
	if artifact.Content == "" {
		t.Fatal("expected report content")
	}
}

func TestRecommendationsAndTrendSummary(t *testing.T) {
	service := NewService(benchmarkStub{rounds: operatorRoundsFixture(t)}, detectionStub{results: operatorDetectionsFixture()}, store.NewOperatorStore())

	recommendations := service.ListRecommendations()
	if len(recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(recommendations))
	}

	recommendation, err := service.GetRecommendation("rec-round-2-replay")
	if err != nil {
		t.Fatalf("GetRecommendation() error = %v", err)
	}
	if recommendation.LinkedRoundID != "round-2" {
		t.Fatalf("unexpected recommendation %#v", recommendation)
	}

	summary := service.GetTrendSummary()
	if summary.RoundsExecuted != 2 {
		t.Fatalf("expected 2 executed rounds, got %d", summary.RoundsExecuted)
	}
}

func operatorRoundsFixture(t *testing.T) map[string]domainbenchmark.BenchmarkRound {
	t.Helper()
	roundOneReport := writeTempReport(t, "round-1", "executive-summary.md", "# Round 1\n")
	roundTwoReport := writeTempReport(t, "round-2", "executive-summary.md", "# Round 2\n")
	return map[string]domainbenchmark.BenchmarkRound{
		"round-1": {
			ID: "round-1",
			Summary: domainbenchmark.RoundSummary{
				PromotionCount:    0,
				ReplayPassRate:    1,
				ChallengerCount:   3,
				RobustnessOutcome: domainbenchmark.RobustnessOutcomeImproved,
				ImportantFindings: []string{"No blind spots were promoted."},
			},
			Reports: domainbenchmark.ReportIndex{
				RoundID: "round-1",
				Artifacts: []domainbenchmark.ReportArtifact{
					{Name: "executive-summary.md", Path: roundOneReport, Kind: "markdown"},
				},
			},
		},
		"round-2": {
			ID: "round-2",
			Summary: domainbenchmark.RoundSummary{
				PromotionCount:    1,
				ReplayPassRate:    0.5,
				ChallengerCount:   3,
				RobustnessOutcome: domainbenchmark.RobustnessOutcomeNewBlindSpotDiscovered,
				ImportantFindings: []string{"Weakened provenance challenger was promoted."},
			},
			PromotionResults: []domainbenchmark.PromotionDecision{{
				ID:                 "promo-2",
				ScenarioID:         "commerce-challenger-weakened-provenance-purchase",
				DetectionResultRef: "det-2",
				ScenarioResultRef:  "sr-living-commerce-challenger-weakened-provenance-purchase",
				Promoted:           true,
			}},
			ScenarioResults: []domainbenchmark.ScenarioResult{{
				ID:                 "sr-living-commerce-challenger-weakened-provenance-purchase",
				ScenarioID:         "commerce-challenger-weakened-provenance-purchase",
				DetectionResultRef: "det-2",
			}},
			Delta: []domainbenchmark.DetectionDelta{{ScenarioID: "commerce-challenger-weakened-provenance-purchase"}},
			Reports: domainbenchmark.ReportIndex{
				RoundID: "round-2",
				Artifacts: []domainbenchmark.ReportArtifact{
					{Name: "executive-summary.md", Path: roundTwoReport, Kind: "markdown"},
				},
			},
			Recommendations: []domainbenchmark.Recommendation{{
				ID:                "rec-round-2-replay",
				Type:              domainbenchmark.RecommendationTypeAddToReplayStableSet,
				LinkedRoundID:     "round-2",
				LinkedScenarioIDs: []string{"commerce-challenger-weakened-provenance-purchase"},
				SuggestedAction:   "Add the promoted case into replay.",
			}},
		},
	}
}

func writeTempReport(t *testing.T, roundID string, name string, content string) string {
	t.Helper()
	path := t.TempDir() + "/" + roundID + "-" + name
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	return path
}

func operatorDetectionsFixture() map[string]detectionmodel.DetectionResult {
	return map[string]detectionmodel.DetectionResult{
		"det-2": {
			ID:             "det-2",
			ScenarioID:     "commerce-challenger-weakened-provenance-purchase",
			Status:         detectionmodel.DetectionStatusClean,
			Recommendation: detectionmodel.RecommendationAllow,
		},
	}
}

func errNotFound(kind string, id string) error {
	return &notFoundError{kind: kind, id: id}
}

type notFoundError struct {
	kind string
	id   string
}

func (e *notFoundError) Error() string {
	return e.kind + " " + e.id + " not found"
}
