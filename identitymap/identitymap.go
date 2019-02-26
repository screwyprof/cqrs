package identitymap

import (
	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
)

type IdentityMap struct {
	identityMap map[uuid.UUID]cqrs.ComplexAggregate
}

func NewIdentityMap() *IdentityMap {
	return &IdentityMap{
		identityMap: map[uuid.UUID]cqrs.ComplexAggregate{},
	}
}

func (m *IdentityMap) ByID(ID uuid.UUID) cqrs.ComplexAggregate {
	return m.identityMap[ID]
}

func (m *IdentityMap) Add(aggregateRoot cqrs.ComplexAggregate) {
	m.identityMap[aggregateRoot.AggregateID()] = aggregateRoot
}

func (m *IdentityMap) Remove(ID uuid.UUID) {
	delete(m.identityMap, ID)
}
