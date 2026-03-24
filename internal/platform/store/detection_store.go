package store

import (
	"fmt"
	"sync"

	"clawbot-trust-lab/internal/domain/detection"
)

type DetectionStore struct {
	mu      sync.RWMutex
	results map[string]detection.DetectionResult
	order   []string
	lastID  string
}

func NewDetectionStore() *DetectionStore {
	return &DetectionStore{
		results: map[string]detection.DetectionResult{},
	}
}

func (s *DetectionStore) Put(item detection.DetectionResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.results[item.ID]; !exists {
		s.order = append(s.order, item.ID)
	}
	s.results[item.ID] = item
	s.lastID = item.ID
}

func (s *DetectionStore) List() []detection.DetectionResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]detection.DetectionResult, 0, len(s.order))
	for _, id := range s.order {
		items = append(items, s.results[id])
	}
	return items
}

func (s *DetectionStore) Get(id string) (detection.DetectionResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.results[id]
	if !ok {
		return detection.DetectionResult{}, fmt.Errorf("detection result %s not found", id)
	}
	return item, nil
}

func (s *DetectionStore) Summary() detection.DetectionRunSummary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	summary := detection.DetectionRunSummary{
		TotalByStatus: map[detection.DetectionStatus]int{},
		Total:         len(s.order),
		LastResultID:  s.lastID,
	}
	for _, item := range s.results {
		summary.TotalByStatus[item.Status]++
	}
	return summary
}
