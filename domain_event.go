package cqrs

import "github.com/google/uuid"

// DomainEvent defines an indication of a point-in-time occurrence.
type DomainEvent interface {
	EventID() uuid.UUID
	EventType() string

	AggregateID() uuid.UUID
	SetAggregateID(id uuid.UUID)

	Version() uint64
	SetVersion(version uint64)
}
