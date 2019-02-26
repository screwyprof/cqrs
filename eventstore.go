package cqrs

import "github.com/google/uuid"

type EventStore interface {
	LoadEventStream(aggregateID uuid.UUID) ([]DomainEvent, error)
	Store(aggregateID uuid.UUID, version uint64, events []DomainEvent) error
}
