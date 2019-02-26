package cqrs

import "github.com/google/uuid"

type IdentityMap interface {
	ByID(ID uuid.UUID) ComplexAggregate
	Add(aggregateRoot ComplexAggregate)
	Remove(ID uuid.UUID)
}
