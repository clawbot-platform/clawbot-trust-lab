package store

import (
	"sort"
	"sync"

	"clawbot-trust-lab/internal/domain/trust"
)

type TrustArtifactStore interface {
	Create(trust.TrustArtifact) error
	List() []trust.TrustArtifact
}

type InMemoryTrustArtifactStore struct {
	mu        sync.RWMutex
	artifacts map[string]trust.TrustArtifact
}

func NewInMemoryTrustArtifactStore() *InMemoryTrustArtifactStore {
	return &InMemoryTrustArtifactStore{artifacts: map[string]trust.TrustArtifact{}}
}

func (s *InMemoryTrustArtifactStore) Create(artifact trust.TrustArtifact) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.artifacts[artifact.ID] = artifact
	return nil
}

func (s *InMemoryTrustArtifactStore) List() []trust.TrustArtifact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]trust.TrustArtifact, 0, len(s.artifacts))
	for _, artifact := range s.artifacts {
		items = append(items, artifact)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}
