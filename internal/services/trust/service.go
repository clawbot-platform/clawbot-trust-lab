package trust

import (
	"fmt"

	domaintrust "clawbot-trust-lab/internal/domain/trust"
)

type Store interface {
	PutMandate(domaintrust.Mandate)
	PutProvenance(domaintrust.ProvenanceRecord)
	PutApproval(domaintrust.ApprovalRecord)
	PutDecision(domaintrust.TrustDecision)
	ListTrustDecisions() []domaintrust.TrustDecision
	GetTrustDecision(string) (domaintrust.TrustDecision, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) RecordMandate(item domaintrust.Mandate) domaintrust.Mandate {
	s.store.PutMandate(item)
	return item
}

func (s *Service) RecordProvenance(item domaintrust.ProvenanceRecord) domaintrust.ProvenanceRecord {
	s.store.PutProvenance(item)
	return item
}

func (s *Service) RecordApproval(item domaintrust.ApprovalRecord) domaintrust.ApprovalRecord {
	s.store.PutApproval(item)
	return item
}

func (s *Service) RecordDecision(item domaintrust.TrustDecision) (domaintrust.TrustDecision, error) {
	if item.ID == "" {
		return domaintrust.TrustDecision{}, fmt.Errorf("trust decision id is required")
	}
	if item.EntityType == "" || item.EntityID == "" {
		return domaintrust.TrustDecision{}, fmt.Errorf("trust decision entity is required")
	}
	s.store.PutDecision(item)
	return item, nil
}

func (s *Service) ListDecisions() []domaintrust.TrustDecision {
	return s.store.ListTrustDecisions()
}

func (s *Service) GetDecision(id string) (domaintrust.TrustDecision, error) {
	return s.store.GetTrustDecision(id)
}
