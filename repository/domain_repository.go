package repository

import (
	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
)

type DomainRepository struct {
	identityMap cqrs.IdentityMap
	eventStore  cqrs.AggregateStore
}

func NewDomainRepository(identityMap cqrs.IdentityMap, eventStore cqrs.AggregateStore) *DomainRepository {
	return &DomainRepository{identityMap, eventStore}
}

func (r *DomainRepository) ByID(ID uuid.UUID, aggregateType string) (cqrs.ComplexAggregate, error) {
	agg := r.identityMap.ByID(ID)
	if agg != nil {
		r.eventStore.Add(agg)
		return agg, nil
	}

	return r.eventStore.ByID(ID, aggregateType)
}

func (r *DomainRepository) Add(aggregateRoot cqrs.ComplexAggregate) {
	r.eventStore.Add(aggregateRoot)
}
