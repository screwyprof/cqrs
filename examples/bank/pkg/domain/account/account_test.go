package account_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/screwyprof/cqrs/pkg/assert"
	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/cqrs/examples/bank/pkg/domain/account"
)

// ensure that game aggregate implements cqrs.Aggregate interface.
var _ cqrs.Aggregate = (*account.Aggregate)(nil)

func TestNewAggregate(t *testing.T) {
	t.Run("ItPanicsIfIDIsNotGiven", func(t *testing.T) {
		factory := func() {
			account.NewAggregate(nil)
		}
		assert.Panic(t, factory)
	})
}

func TestAggregateAggregateID(t *testing.T) {
	t.Run("ItReturnsAggregateID", func(t *testing.T) {
		ID := mock.StringIdentifier(faker.UUIDHyphenated())
		agg := account.NewAggregate(ID)

		assert.Equals(t, ID, agg.AggregateID())
	})
}
