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

func TestFactory(t *testing.T) {
	t.Parallel()

	t.Run("it panics if an aggregate is not registered", func(t *testing.T) {
		t.Parallel()

		f := aggregate.NewFactory()

		_, err := f.CreateAggregate(aggtest.TestAggregateType, aggtest.StringIdentifier(faker.UUIDHyphenated()))

		assert.ErrorIs(t, err, aggregate.ErrAggregateNotRegistered)
	})

	t.Run("it creates an aggregate", func(t *testing.T) {
		t.Parallel()

		// arrange
		id := aggtest.StringIdentifier(faker.UUIDHyphenated())
		f := aggregate.NewFactory()

		// act
		f.RegisterAggregate(aggtest.TestAggregateType, func(ID cqrs.Identifier) cqrs.ESAggregate {
			agg := aggtest.NewTestAggregate(id)

			return aggregate.FromAggregate(agg)
		})

		agg, err := f.CreateAggregate(aggtest.TestAggregateType, id)

		// assert
		assert.NoError(t, err)
		assert.Implements(t, (*cqrs.ESAggregate)(nil), agg)
	})

	t.Run("it converts a cqrs.Aggregate to a cqrs.ESAggregate", func(t *testing.T) {
		t.Parallel()

		// arrange
		id := aggtest.StringIdentifier(faker.UUIDHyphenated())
		agg := aggtest.NewTestAggregate(id)

		// act
		esAgg := aggregate.FromAggregate(agg)

		// assert
		assert.Implements(t, (*cqrs.ESAggregate)(nil), esAgg)
	})
}
