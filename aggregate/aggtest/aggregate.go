package aggtest

import (
	"errors"

	"github.com/screwyprof/cqrs"
)

var (
	ErrItCanHappenOnceOnly                = errors.New("some business rule error occurred")
	ErrMakeSomethingHandlerNotFound       = errors.New("handler for MakeSomethingHappen command is not found")
	ErrOnSomethingHappenedApplierNotFound = errors.New("event applier for OnSomethingHappened event is not found")

	TestAggregateType = "mock.TestAggregate"
)

type StringIdentifier string

func (i StringIdentifier) String() string {
	return string(i)
}

// TestAggregate a pure aggregate (has no external dependencies or dark magic method) used for testing.
type TestAggregate struct {
	id      Identifier
	version int
	aggType string

	alreadyHappened bool
}

// NewTestAggregate creates a new instance of TestAggregate.
func NewTestAggregate(ID cqrs.Identifier) *TestAggregate {
	return &TestAggregate{id: ID}
}

// AggregateID implements cqrs.Aggregate interface.
func (a *TestAggregate) AggregateID() Identifier {
	return a.id
}

// AggregateType implements cqrs.Aggregate interface.
func (a *TestAggregate) AggregateType() string {
	return "mock.TestAggregate"
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
