package aggtest

import (
	"errors"

	"github.com/screwyprof/cqrs"
)

var (
	ErrItCanHappenOnceOnly = errors.New("some business rule error occurred")

	TestAggregateType = "mock.TestAggregate" //nolint:gochecknoglobals
)

type StringIdentifier string

func (i StringIdentifier) String() string {
	return string(i)
}

// TestAggregate is a user-defined aggregate (has no external dependencies or dark magic methods) used for testing.
type TestAggregate struct {
	id Identifier

	alreadyHappened bool
}

// NewTestAggregate creates a new instance of TestAggregate.
func NewTestAggregate(id cqrs.Identifier) *TestAggregate {
	return &TestAggregate{id: id}
}

// AggregateID implements cqrs.Aggregate interface.
func (a *TestAggregate) AggregateID() Identifier {
	return a.id
}

// AggregateType implements cqrs.Aggregate interface.
func (a *TestAggregate) AggregateType() string {
	return TestAggregateType
}

func (a *TestAggregate) MakeSomethingHappen(_ MakeSomethingHappen) ([]Event, error) {
	if a.alreadyHappened {
		return nil, ErrItCanHappenOnceOnly
	}

	return []Event{SomethingHappened{}}, nil
}

func (a *TestAggregate) OnSomethingHappened(_ SomethingHappened) {
	a.alreadyHappened = true
}

func (a *TestAggregate) OnSomethingElseHappened(_ SomethingElseHappened) {
}
