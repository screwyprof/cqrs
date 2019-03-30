package eventstore_test

import (
	"fmt"

	"github.com/bxcodec/faker/v3"

	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/aggregate"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventstore"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"
)

func ExampleInMemoryEventStore_LoadEventsFor() {
	ID := mock.StringIdentifier(faker.UUIDHyphenated())

	es := eventstore.NewInInMemoryEventStore()
	_ = es.StoreEventsFor(ID, 0, []cqrs.DomainEvent{mock.SomethingHappened{}})

	events, _ := es.LoadEventsFor(ID)
	fmt.Printf("%#v", events)

	// Output:
	// []cqrs.DomainEvent{mock.SomethingHappened{Data:""}}
}

func ExampleInMemoryEventStore_StoreEventsFor() {
	ID := mock.StringIdentifier(faker.UUIDHyphenated())

	pureAgg := mock.NewTestAggregate(ID)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(pureAgg)

	aggregate.NewAdvanced(pureAgg, commandHandler, eventApplier)

	es := eventstore.NewInInMemoryEventStore()
	err := es.StoreEventsFor(ID, 1, []cqrs.DomainEvent{mock.SomethingHappened{}})

	fmt.Printf("%v", err)

	// Output:
	// concurrency error: aggregate versions differ
}
