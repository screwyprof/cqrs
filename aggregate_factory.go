package cqrs

import (
	"github.com/google/uuid"
)

type AggregateFactory func(uuid.UUID) ComplexAggregate

var aggregatefactories = make(map[string]AggregateFactory)

func RegisterAggregate(factory AggregateFactory) {
	agg := factory(uuid.New())
	aggregatefactories[agg.AggregateType()] = factory
}

func CreateAggregate(aggregateType string, id uuid.UUID) ComplexAggregate {
	factory, ok := aggregatefactories[aggregateType]
	if !ok {
		panic(aggregateType + " is not registered")
	}
	return factory(id)
}
