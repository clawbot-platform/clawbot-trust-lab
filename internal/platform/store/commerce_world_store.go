package store

import (
	"fmt"
	"sort"
	"sync"

	"clawbot-trust-lab/internal/domain/commerce"
	"clawbot-trust-lab/internal/domain/events"
	"clawbot-trust-lab/internal/domain/trust"
)

type CommerceWorldStore struct {
	mu sync.RWMutex

	buyers        map[string]commerce.Buyer
	merchants     map[string]commerce.Merchant
	products      map[string]commerce.Product
	orders        map[string]commerce.Order
	payments      map[string]commerce.Payment
	refunds       map[string]commerce.Refund
	events        map[string]events.Record
	eventOrder    []string
	mandates      map[string]trust.Mandate
	provenance    map[string]trust.ProvenanceRecord
	approvals     map[string]trust.ApprovalRecord
	decisions     map[string]trust.TrustDecision
	decisionOrder []string
}

func NewCommerceWorldStore() *CommerceWorldStore {
	return &CommerceWorldStore{
		buyers:     map[string]commerce.Buyer{},
		merchants:  map[string]commerce.Merchant{},
		products:   map[string]commerce.Product{},
		orders:     map[string]commerce.Order{},
		payments:   map[string]commerce.Payment{},
		refunds:    map[string]commerce.Refund{},
		events:     map[string]events.Record{},
		mandates:   map[string]trust.Mandate{},
		provenance: map[string]trust.ProvenanceRecord{},
		approvals:  map[string]trust.ApprovalRecord{},
		decisions:  map[string]trust.TrustDecision{},
	}
}

func (s *CommerceWorldStore) PutBuyer(item commerce.Buyer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buyers[item.ID] = item
}

func (s *CommerceWorldStore) PutMerchant(item commerce.Merchant) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.merchants[item.ID] = item
}

func (s *CommerceWorldStore) PutProduct(item commerce.Product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.products[item.ID] = item
}

func (s *CommerceWorldStore) GetProduct(id string) (commerce.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.products[id]
	if !ok {
		return commerce.Product{}, fmt.Errorf("product %s not found", id)
	}
	return item, nil
}

func (s *CommerceWorldStore) PutOrder(item commerce.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[item.ID] = item
}

func (s *CommerceWorldStore) PutPayment(item commerce.Payment) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.payments[item.ID] = item
}

func (s *CommerceWorldStore) PutRefund(item commerce.Refund) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refunds[item.ID] = item
}

func (s *CommerceWorldStore) ListRefunds() []commerce.Refund {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]commerce.Refund, 0, len(s.refunds))
	for _, item := range s.refunds {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (s *CommerceWorldStore) GetRefund(id string) (commerce.Refund, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.refunds[id]
	if !ok {
		return commerce.Refund{}, fmt.Errorf("refund %s not found", id)
	}
	return item, nil
}

func (s *CommerceWorldStore) ListOrders() []commerce.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]commerce.Order, 0, len(s.orders))
	for _, item := range s.orders {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (s *CommerceWorldStore) GetOrder(id string) (commerce.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.orders[id]
	if !ok {
		return commerce.Order{}, fmt.Errorf("order %s not found", id)
	}
	return item, nil
}

func (s *CommerceWorldStore) PutMandate(item trust.Mandate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.mandates[item.ID] = item
}

func (s *CommerceWorldStore) PutProvenance(item trust.ProvenanceRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.provenance[item.ID] = item
}

func (s *CommerceWorldStore) PutApproval(item trust.ApprovalRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.approvals[item.ID] = item
}

func (s *CommerceWorldStore) ListApprovals() []trust.ApprovalRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]trust.ApprovalRecord, 0, len(s.approvals))
	for _, item := range s.approvals {
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items
}

func (s *CommerceWorldStore) GetMandate(id string) (trust.Mandate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.mandates[id]
	if !ok {
		return trust.Mandate{}, fmt.Errorf("mandate %s not found", id)
	}
	return item, nil
}

func (s *CommerceWorldStore) GetProvenance(id string) (trust.ProvenanceRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.provenance[id]
	if !ok {
		return trust.ProvenanceRecord{}, fmt.Errorf("provenance %s not found", id)
	}
	return item, nil
}

func (s *CommerceWorldStore) PutDecision(item trust.TrustDecision) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.decisions[item.ID]; !exists {
		s.decisionOrder = append(s.decisionOrder, item.ID)
	}
	s.decisions[item.ID] = item
}

func (s *CommerceWorldStore) ListTrustDecisions() []trust.TrustDecision {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]trust.TrustDecision, 0, len(s.decisionOrder))
	for _, id := range s.decisionOrder {
		items = append(items, s.decisions[id])
	}
	return items
}

func (s *CommerceWorldStore) GetTrustDecision(id string) (trust.TrustDecision, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.decisions[id]
	if !ok {
		return trust.TrustDecision{}, fmt.Errorf("trust decision %s not found", id)
	}
	return item, nil
}

func (s *CommerceWorldStore) AppendEvent(item events.Record) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[item.ID]; !exists {
		s.eventOrder = append(s.eventOrder, item.ID)
	}
	s.events[item.ID] = item
}

func (s *CommerceWorldStore) ListEvents() []events.Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := make([]events.Record, 0, len(s.eventOrder))
	for _, id := range s.eventOrder {
		items = append(items, s.events[id])
	}
	return items
}
