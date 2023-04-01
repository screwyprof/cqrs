package bank_test

import (
	"fmt"
	"os"

	"github.com/go-faker/faker/v4"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/aggregate"
	"github.com/screwyprof/cqrs/aggregate/aggtest"
	"github.com/screwyprof/cqrs/examples/bank/domain/account"
	"github.com/screwyprof/cqrs/examples/bank/domain/command"
	eh "github.com/screwyprof/cqrs/examples/bank/eventhandler"
	"github.com/screwyprof/cqrs/examples/bank/reporting"
	"github.com/screwyprof/cqrs/examples/bank/ui"
	"github.com/screwyprof/cqrs/x/aggstore"
	"github.com/screwyprof/cqrs/x/dispatcher"
	"github.com/screwyprof/cqrs/x/eventbus"
	"github.com/screwyprof/cqrs/x/eventhandler"
	"github.com/screwyprof/cqrs/x/eventstore"
)

func Example() {
	ID := aggtest.StringIdentifier(faker.UUIDHyphenated())
	AccNumber := "ACC777"

	accountReporter := reporting.NewInMemoryAccountReporter()

	d := createDispatcher(accountReporter)

	failCommandOnError(d.Handle(command.OpenAccount{ID: ID, Number: AccNumber}))
	failCommandOnError(d.Handle(command.DepositMoney{ID: ID, Amount: 1000}))
	failCommandOnError(d.Handle(command.WithdrawMoney{ID: ID, Amount: 100}))
	failCommandOnError(d.Handle(command.DepositMoney{ID: ID, Amount: 500}))

	printer := ui.NewConsolePrinter(os.Stdout, accountReporter)
	failOnError(printer.PrintAccountStatement(ID))

	// Output:
	// Account #ACC777:
	// # |   Amount |  Balance
	// 1 |  1000.00 |  1000.00
	// 2 |  -100.00 |   900.00
	// 3 |   500.00 |  1400.00
}

func createDispatcher(accountReporter eh.AccountReporting) *dispatcher.Dispatcher {
	aggregateFactory := aggregate.NewFactory()
	aggregateFactory.RegisterAggregate("account.Aggregate", createAggregate)

	accountDetailsProjector := eventhandler.New()
	accountDetailsProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

	eventPublisher := eventbus.NewInMemoryEventBus()
	eventPublisher.Register(accountDetailsProjector)

	aggregateStore := aggstore.NewStore(
		eventstore.NewInInMemoryEventStore(eventPublisher),
		aggregateFactory,
	)

	return dispatcher.NewDispatcher(aggregateStore)
}

func createAggregate(ID cqrs.Identifier) cqrs.ESAggregate {
	acc := account.NewAggregate(ID)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(acc)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(acc)

	return aggregate.New(acc, commandHandler, eventApplier)
}

func failCommandOnError(_ []cqrs.DomainEvent, err error) {
	failOnError(err)
}

func failOnError(err error) {
	if err != nil {
		fmt.Printf("an error occurred: %v", err)
		os.Exit(1)
	}
}
