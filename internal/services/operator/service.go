package operator

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	domainbenchmark "clawbot-trust-lab/internal/domain/benchmark"
	detectionmodel "clawbot-trust-lab/internal/domain/detection"
)

type BenchmarkReader interface {
	ListRounds() []domainbenchmark.BenchmarkRound
	GetRound(string) (domainbenchmark.BenchmarkRound, error)
	GetRoundReports(string) (domainbenchmark.ReportIndex, error)
}

type DetectionReader interface {
	GetResult(string) (detectionmodel.DetectionResult, error)
}

type ReviewStore interface {
	PutReview(domainbenchmark.PromotionReview)
	GetReview(string) (domainbenchmark.PromotionReview, bool)
	ListReviews() []domainbenchmark.PromotionReview
}

type ReviewInput struct {
	Status string `json:"status"`
	Note   string `json:"note"`
}

type PromotionRecord struct {
	RoundID   string                            `json:"round_id"`
	Promotion domainbenchmark.PromotionDecision `json:"promotion"`
	Review    *domainbenchmark.PromotionReview  `json:"review,omitempty"`
}

type PromotionDetail struct {
	RoundID         string                            `json:"round_id"`
	Promotion       domainbenchmark.PromotionDecision `json:"promotion"`
	Review          *domainbenchmark.PromotionReview  `json:"review,omitempty"`
	DetectionResult detectionmodel.DetectionResult    `json:"detection_result"`
	ScenarioResult  *domainbenchmark.ScenarioResult   `json:"scenario_result,omitempty"`
}

type ReportContent struct {
	Descriptor domainbenchmark.ReportDescriptor `json:"descriptor"`
	Content    string                           `json:"content"`
}

type Service struct {
	benchmark BenchmarkReader
	detection DetectionReader
	reviews   ReviewStore
	now       func() time.Time
}

func NewService(benchmark BenchmarkReader, detection DetectionReader, reviews ReviewStore) *Service {
	return &Service{
		benchmark: benchmark,
		detection: detection,
		reviews:   reviews,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) ListRounds() []domainbenchmark.BenchmarkRound {
	return s.benchmark.ListRounds()
}

func (s *Service) GetRound(id string) (domainbenchmark.BenchmarkRound, error) {
	return s.benchmark.GetRound(id)
}

func (s *Service) CompareRounds(currentID string, previousID string) (domainbenchmark.RoundComparison, error) {
	current, err := s.benchmark.GetRound(strings.TrimSpace(currentID))
	if err != nil {
		return domainbenchmark.RoundComparison{}, err
	}
	previous, err := s.benchmark.GetRound(strings.TrimSpace(previousID))
	if err != nil {
		return domainbenchmark.RoundComparison{}, err
	}

	return domainbenchmark.RoundComparison{
		CurrentRoundID:          current.ID,
		PreviousRoundID:         previous.ID,
		CurrentRobustness:       current.Summary.RobustnessOutcome,
		PreviousRobustness:      previous.Summary.RobustnessOutcome,
		PromotionsCountDelta:    current.Summary.PromotionCount - previous.Summary.PromotionCount,
		ReplayPassRateDelta:     current.Summary.ReplayPassRate - previous.Summary.ReplayPassRate,
		ChallengerCountDelta:    current.Summary.ChallengerCount - previous.Summary.ChallengerCount,
		ImportantFindingsAdded:  diffStrings(current.Summary.ImportantFindings, previous.Summary.ImportantFindings),
		ImportantFindingsClosed: diffStrings(previous.Summary.ImportantFindings, current.Summary.ImportantFindings),
		DetectionDeltaCount:     len(current.Delta),
	}, nil
}

func (s *Service) ListPromotions(statusFilter string) []PromotionRecord {
	statusFilter = strings.TrimSpace(statusFilter)
	rounds := s.benchmark.ListRounds()
	records := make([]PromotionRecord, 0)
	for _, round := range rounds {
		for _, item := range round.PromotionResults {
			record := PromotionRecord{
				RoundID:   round.ID,
				Promotion: item,
			}
			if review, ok := s.reviews.GetReview(item.ID); ok {
				record.Review = &review
			}
			if statusFilter != "" {
				if record.Review == nil || string(record.Review.Status) != statusFilter {
					continue
				}
			}
			records = append(records, record)
		}
	}
	return records
}

func (s *Service) GetPromotion(id string) (PromotionDetail, error) {
	for _, round := range s.benchmark.ListRounds() {
		for _, item := range round.PromotionResults {
			if item.ID != id {
				continue
			}
			detectionResult, err := s.detection.GetResult(item.DetectionResultRef)
			if err != nil {
				return PromotionDetail{}, err
			}
			detail := PromotionDetail{
				RoundID:         round.ID,
				Promotion:       item,
				DetectionResult: detectionResult,
			}
			if review, ok := s.reviews.GetReview(item.ID); ok {
				detail.Review = &review
			}
			for _, result := range round.ScenarioResults {
				if result.ID == item.ScenarioResultRef {
					copied := result
					detail.ScenarioResult = &copied
					break
				}
			}
			return detail, nil
		}
	}
	return PromotionDetail{}, fmt.Errorf("promotion %s not found", id)
}

func (s *Service) ReviewPromotion(id string, input ReviewInput) (domainbenchmark.PromotionReview, error) {
	status, err := parseReviewStatus(input.Status)
	if err != nil {
		return domainbenchmark.PromotionReview{}, err
	}
	if _, err := s.GetPromotion(id); err != nil {
		return domainbenchmark.PromotionReview{}, err
	}

	review := domainbenchmark.PromotionReview{
		PromotionID: id,
		Status:      status,
		UpdatedAt:   s.now(),
	}
	if note := strings.TrimSpace(input.Note); note != "" {
		review.Note = &domainbenchmark.OperatorNote{
			ID:        "note-" + strings.ReplaceAll(id, "promo-", ""),
			Body:      note,
			CreatedAt: s.now(),
		}
	}
	s.reviews.PutReview(review)
	return review, nil
}

func (s *Service) GetDetectionResult(id string) (detectionmodel.DetectionResult, error) {
	return s.detection.GetResult(id)
}

func (s *Service) GetReports(roundID string) ([]domainbenchmark.ReportDescriptor, error) {
	index, err := s.benchmark.GetRoundReports(roundID)
	if err != nil {
		return nil, err
	}
	descriptors := make([]domainbenchmark.ReportDescriptor, 0, len(index.Artifacts))
	for _, artifact := range index.Artifacts {
		descriptors = append(descriptors, domainbenchmark.ReportDescriptor{
			RoundID:      roundID,
			ArtifactName: artifact.Name,
			Path:         artifact.Path,
			Kind:         artifact.Kind,
		})
	}
	return descriptors, nil
}

func (s *Service) GetReportArtifact(roundID string, artifactName string) (ReportContent, error) {
	reports, err := s.GetReports(roundID)
	if err != nil {
		return ReportContent{}, err
	}
	for _, item := range reports {
		if item.ArtifactName != artifactName {
			continue
		}
		body, err := os.ReadFile(item.Path)
		if err != nil {
			return ReportContent{}, fmt.Errorf("read report artifact %s: %w", item.Path, err)
		}
		return ReportContent{
			Descriptor: item,
			Content:    string(body),
		}, nil
	}
	return ReportContent{}, fmt.Errorf("report artifact %s for round %s not found", artifactName, roundID)
}

func parseReviewStatus(value string) (domainbenchmark.PromotionReviewStatus, error) {
	status := domainbenchmark.PromotionReviewStatus(strings.TrimSpace(value))
	switch status {
	case domainbenchmark.PromotionReviewAccepted,
		domainbenchmark.PromotionReviewDuplicate,
		domainbenchmark.PromotionReviewNeedsFollowUp,
		domainbenchmark.PromotionReviewFalseSignal:
		return status, nil
	default:
		return "", fmt.Errorf("invalid promotion review status %q", value)
	}
}

func diffStrings(current []string, previous []string) []string {
	out := make([]string, 0)
	for _, item := range current {
		if !slices.Contains(previous, item) {
			out = append(out, item)
		}
	}
	return out
}
