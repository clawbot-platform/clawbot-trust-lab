package events

import (
	"time"

	domainevents "clawbot-trust-lab/internal/domain/events"
)

type Store interface {
	AppendEvent(domainevents.Record)
	ListEvents() []domainevents.Record
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) RecordTransaction(id string, eventType domainevents.TransactionEventType, entityType string, entityID string, scenarioID string, actorID string, occurredAt time.Time, metadata map[string]any) domainevents.TransactionEvent {
	item := domainevents.TransactionEvent{
		Record: domainevents.Record{
			ID:         id,
			Category:   domainevents.EventCategoryTransaction,
			EventType:  string(eventType),
			EntityType: entityType,
			EntityID:   entityID,
			ScenarioID: scenarioID,
			ActorID:    actorID,
			OccurredAt: occurredAt,
			Metadata:   cloneMap(metadata),
		},
	}
	s.store.AppendEvent(item.Record)
	return item
}

func (s *Service) RecordTrust(id string, eventType domainevents.TrustEventType, entityType string, entityID string, scenarioID string, actorID string, occurredAt time.Time, metadata map[string]any) domainevents.TrustEvent {
	item := domainevents.TrustEvent{
		Record: domainevents.Record{
			ID:         id,
			Category:   domainevents.EventCategoryTrust,
			EventType:  string(eventType),
			EntityType: entityType,
			EntityID:   entityID,
			ScenarioID: scenarioID,
			ActorID:    actorID,
			OccurredAt: occurredAt,
			Metadata:   cloneMap(metadata),
		},
	}
	s.store.AppendEvent(item.Record)
	return item
}

func (s *Service) ListEvents() []domainevents.Record {
	return s.store.ListEvents()
}

func cloneMap(input map[string]any) map[string]any {
	if input == nil {
		return map[string]any{}
	}
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
