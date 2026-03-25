package store

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"clawbot-trust-lab/internal/domain/benchmark"
)

type BenchmarkStore struct {
	mu         sync.RWMutex
	live       map[string]benchmark.BenchmarkRound
	historical map[string]benchmark.BenchmarkRound
}

func NewBenchmarkStore() *BenchmarkStore {
	return &BenchmarkStore{
		live:       map[string]benchmark.BenchmarkRound{},
		historical: map[string]benchmark.BenchmarkRound{},
	}
}

func (s *BenchmarkStore) Put(round benchmark.BenchmarkRound) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.live[round.ID] = round
}

func (s *BenchmarkStore) PutHistorical(round benchmark.BenchmarkRound) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.historical[round.ID] = round
}

func (s *BenchmarkStore) List() []benchmark.BenchmarkRound {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]benchmark.BenchmarkRound, 0, len(s.live)+len(s.historical))
	for _, id := range s.sortedIDs() {
		item, ok := s.mergedRound(id)
		if ok {
			items = append(items, item)
		}
	}
	return items
}

func (s *BenchmarkStore) Get(id string) (benchmark.BenchmarkRound, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	round, ok := s.mergedRound(id)
	if !ok {
		return benchmark.BenchmarkRound{}, fmt.Errorf("benchmark round %s not found", id)
	}
	return round, nil
}

func (s *BenchmarkStore) Latest() (benchmark.BenchmarkRound, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ids := s.sortedIDs()
	if len(ids) == 0 {
		return benchmark.BenchmarkRound{}, false
	}
	item, ok := s.mergedRound(ids[0])
	return item, ok
}

func (s *BenchmarkStore) sortedIDs() []string {
	index := make(map[string]struct{}, len(s.live)+len(s.historical))
	for id := range s.historical {
		index[id] = struct{}{}
	}
	for id := range s.live {
		index[id] = struct{}{}
	}

	ids := make([]string, 0, len(index))
	for id := range index {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool {
		left, _ := s.mergedRound(ids[i])
		right, _ := s.mergedRound(ids[j])
		leftTime := roundSortTime(left)
		rightTime := roundSortTime(right)
		if leftTime.Equal(rightTime) {
			return left.ID > right.ID
		}
		return leftTime.After(rightTime)
	})

	return ids
}

func (s *BenchmarkStore) mergedRound(id string) (benchmark.BenchmarkRound, bool) {
	live, hasLive := s.live[id]
	historical, hasHistorical := s.historical[id]

	switch {
	case hasLive && hasHistorical:
		return mergeRound(live, historical), true
	case hasLive:
		return live, true
	case hasHistorical:
		return historical, true
	default:
		return benchmark.BenchmarkRound{}, false
	}
}

func mergeRound(live benchmark.BenchmarkRound, historical benchmark.BenchmarkRound) benchmark.BenchmarkRound {
	merged := live
	if merged.ReportDir == "" {
		merged.ReportDir = historical.ReportDir
	}
	if merged.Reports.Directory == "" || len(merged.Reports.Artifacts) == 0 {
		merged.Reports = historical.Reports
	}
	if merged.Summary.RoundID == "" {
		merged.Summary.RoundID = historical.Summary.RoundID
	}
	return merged
}

func roundSortTime(item benchmark.BenchmarkRound) time.Time {
	if !item.CompletedAt.IsZero() {
		return item.CompletedAt
	}
	if !item.StartedAt.IsZero() {
		return item.StartedAt
	}
	return time.Time{}
}
