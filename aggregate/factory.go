package aggregate

import (
	"errors"

	"github.com/go-faker/faker/v4"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
)

// Factory is responsible for creating aggregates by their type.
//
// It maintains a registry of factory functions for different aggregate types.
type Factory struct {
	factories map[string]cqrs.FactoryFn
}

// NewFactory creates a new instance of Factory and initializes its internal factory registry.
//
// It returns a pointer to the created Factory instance.
func NewFactory() *Factory {
	return &Factory{
		factories: make(map[string]cqrs.FactoryFn),
	}
}

// RegisterAggregate registers an aggregate factory method.
//
// The factory function is used to create aggregates of a specific type.
func (f *Factory) RegisterAggregate(factory cqrs.FactoryFn) {
	agg := factory(aggtest.StringIdentifier(faker.UUIDHyphenated()))
	f.factories[agg.AggregateType()] = factory
}

// CreateAggregate creates an aggregate of a given type.
//
// It returns an error if the requested aggregate type is not registered in the factory.
func (f *Factory) CreateAggregate(aggregateType string, id cqrs.Identifier) (cqrs.ESAggregate, error) {
	factory, ok := f.factories[aggregateType]
	if !ok {
		return nil, errors.New(aggregateType + " is not registered")
	}

	return factory(id), nil
}
