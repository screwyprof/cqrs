package cqrs

import "github.com/google/uuid"

type EventStore interface {
	LoadEventStream(aggregateID uuid.UUID) ([]DomainEvent, error)
	Store(eventProvider EventProvider) error
}
