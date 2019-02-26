package cqrs

import (
	"github.com/google/uuid"
)

type AggregateStore interface {
	UnitOfWork

	ByID(ID uuid.UUID, aggregateType string) (ComplexAggregate, error)
	Add(aggregateRoot ComplexAggregate)
}
