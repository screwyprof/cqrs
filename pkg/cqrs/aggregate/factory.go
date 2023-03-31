package aggregate

import (
	"errors"
	"sync"

	"github.com/bxcodec/faker/v4"

	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
)

// Factory handles aggregate creation.
type Factory struct {
	factories   map[string]cqrs.FactoryFn
	factoriesMu sync.RWMutex
}

// NewFactory creates a new instance of Factory.
func NewFactory() *Factory {
	return &Factory{
		factories: make(map[string]cqrs.FactoryFn),
	}
}

// RegisterAggregate registers an aggregate factory method.
func (f *Factory) RegisterAggregate(factory cqrs.FactoryFn) {
	f.factoriesMu.Lock()
	defer f.factoriesMu.Unlock()

	agg := factory(mock.StringIdentifier(faker.UUIDHyphenated()))
	f.factories[agg.AggregateType()] = factory
}

// CreateAggregate creates an aggregate of a given type.
func (f *Factory) CreateAggregate(aggregateType string, id cqrs.Identifier) (cqrs.AdvancedAggregate, error) {
	f.factoriesMu.Lock()
	defer f.factoriesMu.Unlock()

	factory, ok := f.factories[aggregateType]
	if !ok {
		return nil, errors.New(aggregateType + " is not registered")
	}
	return factory(id), nil
}
