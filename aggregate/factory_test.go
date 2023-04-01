package aggregate_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
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

		_, err := f.CreateAggregate(aggtest.TestAggregateType, aggtest.StringIdentifier(faker.UUIDHyphenated()))

		assert.Equal(t, aggtest.ErrAggIsNotRegistered, err)
	})
}

func TestFactoryRegisterAggregate(t *testing.T) {
	t.Run("ItRegistersAnAggregateFactory", func(t *testing.T) {
		// arrange
		ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
		agg := aggtest.NewTestAggregate(ID)

		commandHandler := aggregate.NewCommandHandler()
		commandHandler.RegisterHandlers(agg)

		eventApplier := aggregate.NewEventApplier()
		eventApplier.RegisterAppliers(agg)

		expected := aggregate.New(
			agg,
			commandHandler,
			eventApplier,
		)

		f := aggregate.NewFactory()

		// act
		f.RegisterAggregate(func(ID cqrs.Identifier) cqrs.ESAggregate {
			return expected
		})
		newAgg, err := f.CreateAggregate(aggtest.TestAggregateType, ID)

		// assert
		assert.NoError(t, err)
		assert.Equal(t, expected, newAgg)
	})
}
