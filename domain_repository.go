package cqrs

import "github.com/google/uuid"

type DomainRepository interface {
	ByID(ID uuid.UUID, aggregateType string) (ComplexAggregate, error)
	Add(aggregateRoot ComplexAggregate)
}
