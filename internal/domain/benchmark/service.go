package benchmark

import (
	"context"
	"sync"
)

type ControlPlaneClient interface {
	RegisterBenchmarkMetadata(context.Context, RegistrationRequest) (RegistrationResult, error)
}

type Service struct {
	client  ControlPlaneClient
	mu      sync.RWMutex
	history []RegistrationResult
}

func NewService(client ControlPlaneClient) *Service {
	return &Service{client: client}
}

func (s *Service) RegisterRound(ctx context.Context, request RegistrationRequest) (RegistrationResult, error) {
	result, err := s.client.RegisterBenchmarkMetadata(ctx, request)
	if err != nil {
		return RegistrationResult{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.history = append(s.history, result)
	return result, nil
}

func (s *Service) Status() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lastStatus := "idle"
	if len(s.history) > 0 {
		lastStatus = s.history[len(s.history)-1].Status
	}

	return map[string]any{
		"registrations": len(s.history),
		"last_status":   lastStatus,
	}
}
