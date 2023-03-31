package mock

import (
	"github.com/screwyprof/cqrs"
)

type MakeSomethingHappen struct {
	AggID cqrs.Identifier
}

func (c MakeSomethingHappen) AggregateID() cqrs.Identifier {
	return c.AggID
}

func (c MakeSomethingHappen) AggregateType() string {
	return "mock.TestAggregate"
}

func (c MakeSomethingHappen) CommandType() string {
	return "MakeSomethingHappen"
}
