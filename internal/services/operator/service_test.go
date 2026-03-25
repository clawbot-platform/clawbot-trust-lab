package operator

import (
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

func operatorRoundsFixture(t *testing.T) map[string]domainbenchmark.BenchmarkRound {
	t.Helper()
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
					{Name: "executive-summary.md", Path: t.TempDir() + "/executive-summary.md", Kind: "markdown"},
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
					{Name: "executive-summary.md", Path: t.TempDir() + "/executive-summary.md", Kind: "markdown"},
				},
			},
		},
	}
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
