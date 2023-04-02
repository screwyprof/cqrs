package aggregate_test

import (
	"fmt"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
)

const AggregateType = "MyAggregate"

// Identifier is a user-defined type that implements the fmt.Stringer interface.
type Identifier = fmt.Stringer

// Event is a user-defined interface that represents events in the domain.
type Event interface {
	EventType() string
}

// MyIdentifier is a user-defined identifier that implements the Identifier interface.
type MyIdentifier string

func (id MyIdentifier) String() string {
	return string(id)
}

// SomethingHappened is a user-defined event that implements the Event interface.
type SomethingHappened struct {
	ID      MyIdentifier
	Changed bool
}

func (e SomethingHappened) EventType() string {
	return "SomethingHappened"
}

// MyAggregate is a user-defined aggregate that will handle commands and apply events.
type MyAggregate struct {
	id         MyIdentifier
	isModified bool
}

func (a *MyAggregate) AggregateID() Identifier {
	return a.id
}

func (a *MyAggregate) AggregateType() string {
	return AggregateType
}

// DoSomething is an example of a command.
type DoSomething struct {
	ID MyIdentifier
}

func (c DoSomething) AggregateID() cqrs.Identifier {
	return c.ID
}

func (c DoSomething) AggregateType() string {
	return AggregateType
}

func (c DoSomething) CommandType() string {
	return "DoSomething"
}

// DoSomething is a command handler method for the MyAggregate.
func (a *MyAggregate) DoSomething(c DoSomething) ([]Event, error) {
	e := SomethingHappened{
		ID:      c.ID,
		Changed: true,
	}

	return []Event{e}, nil
}

// OnSomethingHappened is an event applier method for the MyAggregate.
func (a *MyAggregate) OnSomethingHappened(e SomethingHappened) {
	a.isModified = e.Changed
}

// Example demonstrates the use of the aggregate package to handle a command and produce events.
func Example() {
	id := MyIdentifier("123")
	agg := &MyAggregate{id: id}
	esAgg := aggregate.FromAggregate(agg)

	events, err := esAgg.Handle(DoSomething{ID: id})
	if err != nil {
		fmt.Printf("an error occurred: %v\n", err)

		return
	}

	fmt.Printf("Produced events: %v\n", events)
	// Output: Produced events: [{123 true}]
}
