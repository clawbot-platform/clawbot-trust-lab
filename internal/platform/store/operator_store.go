package store

import (
	"sync"

	"clawbot-trust-lab/internal/domain/benchmark"
)

type OperatorStore struct {
	mu      sync.RWMutex
	reviews map[string]benchmark.PromotionReview
}

func NewOperatorStore() *OperatorStore {
	return &OperatorStore{
		reviews: map[string]benchmark.PromotionReview{},
	}
}

func (s *OperatorStore) PutReview(review benchmark.PromotionReview) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reviews[review.PromotionID] = review
}

func (s *OperatorStore) GetReview(promotionID string) (benchmark.PromotionReview, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	review, ok := s.reviews[promotionID]
	return review, ok
}

func (s *OperatorStore) ListReviews() []benchmark.PromotionReview {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]benchmark.PromotionReview, 0, len(s.reviews))
	for _, review := range s.reviews {
		items = append(items, review)
	}
	return items
}
