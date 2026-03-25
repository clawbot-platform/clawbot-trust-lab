package store

import (
	"fmt"
	"sync"

	"clawbot-trust-lab/internal/domain/benchmark"
)

type BenchmarkStore struct {
	mu     sync.RWMutex
	rounds map[string]benchmark.BenchmarkRound
	order  []string
}

func NewBenchmarkStore() *BenchmarkStore {
	return &BenchmarkStore{
		rounds: map[string]benchmark.BenchmarkRound{},
	}
}

func (s *BenchmarkStore) Put(round benchmark.BenchmarkRound) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.rounds[round.ID]; !exists {
		s.order = append(s.order, round.ID)
	}
	s.rounds[round.ID] = round
}

func (s *BenchmarkStore) List() []benchmark.BenchmarkRound {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]benchmark.BenchmarkRound, 0, len(s.order))
	for _, id := range s.order {
		items = append(items, s.rounds[id])
	}
	return items
}

func (s *BenchmarkStore) Get(id string) (benchmark.BenchmarkRound, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	round, ok := s.rounds[id]
	if !ok {
		return benchmark.BenchmarkRound{}, fmt.Errorf("benchmark round %s not found", id)
	}
	return round, nil
}

func (s *BenchmarkStore) Latest() (benchmark.BenchmarkRound, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.order) == 0 {
		return benchmark.BenchmarkRound{}, false
	}
	last := s.order[len(s.order)-1]
	return s.rounds[last], true
}
