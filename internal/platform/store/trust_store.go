package store

import (
	"sync"

	"clawbot-trust-lab/internal/domain/trust"
)

type TrustArtifactStore interface {
	Create(trust.TrustArtifact) error
	List() []trust.TrustArtifact
}

type InMemoryTrustArtifactStore struct {
	mu        sync.RWMutex
	artifacts []trust.TrustArtifact
}

func NewInMemoryTrustArtifactStore() *InMemoryTrustArtifactStore {
	return &InMemoryTrustArtifactStore{}
}

func (s *InMemoryTrustArtifactStore) Create(artifact trust.TrustArtifact) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.artifacts = append(s.artifacts, artifact)
	return nil
}

func (s *InMemoryTrustArtifactStore) List() []trust.TrustArtifact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]trust.TrustArtifact, len(s.artifacts))
	copy(items, s.artifacts)
	return items
}
