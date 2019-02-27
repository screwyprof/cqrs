package cqrs

import "github.com/google/uuid"

type Aggregate interface {
	AggregateID() uuid.UUID
	AggregateType() string
}
