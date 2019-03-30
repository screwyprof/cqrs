package mock

import (
	"errors"

	"github.com/screwyprof/cqrs/pkg/cqrs"
)

var (
	ErrItCanHappenOnceOnly  = errors.New("some business rule error occurred")
	ErrMakeSomethingHandlerNotFound  = errors.New("handler for MakeSomethingHappen command is not found")
	ErrOnSomethingHappenedApplierNotFound  = errors.New("event applier for OnSomethingHappened event is not found")

	TestAggregateType = "mock.TestAggregate"
)

type StringIdentifier string
func (i StringIdentifier) String() string {
	return string(i)
}

// TestAggregate a pure aggregate (has no external dependencies or dark magic method) used for testing.
type TestAggregate struct {
	id cqrs.Identifier
	version int
	aggType string

	alreadyHappened bool
}

// NewTestAggregate creates a new instance of TestAggregate.
func NewTestAggregate(ID cqrs.Identifier) *TestAggregate {
	return &TestAggregate{id:ID}
}

// AggregateID implements cqrs.Aggregate interface.
func (a *TestAggregate) AggregateID() cqrs.Identifier {
	return a.id
}

// AggregateType implements cqrs.Aggregate interface.
func (a *TestAggregate) AggregateType() string {
	return "mock.TestAggregate"
}

func (a *TestAggregate) MakeSomethingHappen(c MakeSomethingHappen) ([]cqrs.DomainEvent, error) {
	if a.alreadyHappened {
		return nil, ErrItCanHappenOnceOnly
	}
	return []cqrs.DomainEvent{SomethingHappened{}}, nil
}

func (a *TestAggregate) OnSomethingHappened(e SomethingHappened) {
	a.alreadyHappened = true
}

func (a *TestAggregate) OnSomethingElseHappened(e SomethingElseHappened) {

}
