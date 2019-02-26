package cqrs

import "github.com/google/uuid"

type Command interface {
	AggregateID() uuid.UUID
	AggregateType() string
}
