package aggregate_test

import (
	"testing"

	"github.com/go-faker/faker/v4"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/aggregate"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
)

// ensure that factory implements cqrs.AggregateFactory interface.
var _ cqrs.AggregateFactory = (*aggregate.Factory)(nil)

func TestNewFactory(t *testing.T) {
	t.Run("ItReturnsNewFactoryInstance", func(t *testing.T) {
		f := aggregate.NewFactory()
		assert.True(t, f != nil)
	})
}

func TestFactoryCreateAggregate(t *testing.T) {
	t.Run("ItPanicsIfTheAggregateIsNotRegistered", func(t *testing.T) {
		f := aggregate.NewFactory()

		_, err := f.CreateAggregate(mock.TestAggregateType, mock.StringIdentifier(faker.UUIDHyphenated()))

		assert.Equals(t, mock.ErrAggIsNotRegistered, err)
	})
}

func TestFactoryRegisterAggregate(t *testing.T) {
	t.Run("ItRegistersAnAggregateFactory", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		agg := mock.NewTestAggregate(ID)

		commandHandler := aggregate.NewCommandHandler()
		commandHandler.RegisterHandlers(agg)

		eventApplier := aggregate.NewEventApplier()
		eventApplier.RegisterAppliers(agg)

		expected := aggregate.NewAdvanced(
			agg,
			commandHandler,
			eventApplier,
		)

		f := aggregate.NewFactory()

		// act
		f.RegisterAggregate(func(ID cqrs.Identifier) cqrs.AdvancedAggregate {
			return expected
		})
		newAgg, err := f.CreateAggregate(mock.TestAggregateType, ID)

		// assert
		assert.Ok(t, err)
		assert.Equals(t, expected, newAgg)
	})
}
