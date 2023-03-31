package bank_test

import (
	"fmt"
	"os"

	"github.com/bxcodec/faker/v4"

	"github.com/screwyprof/cqrs/pkg/cqrs"
	"github.com/screwyprof/cqrs/pkg/cqrs/aggregate"
	"github.com/screwyprof/cqrs/pkg/cqrs/dispatcher"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventbus"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventhandler"
	"github.com/screwyprof/cqrs/pkg/cqrs/eventstore"
	"github.com/screwyprof/cqrs/pkg/cqrs/store"
	"github.com/screwyprof/cqrs/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/cqrs/examples/bank/internal/reporting"
	"github.com/screwyprof/cqrs/examples/bank/internal/ui"

	"github.com/screwyprof/cqrs/examples/bank/pkg/command"
	"github.com/screwyprof/cqrs/examples/bank/pkg/domain/account"
	eh "github.com/screwyprof/cqrs/examples/bank/pkg/eventhandler"
	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

func Example() {
	ID := mock.StringIdentifier(faker.UUIDHyphenated())
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

func createDispatcher(accountReporter report.AccountReporting) *dispatcher.Dispatcher {
	aggregateFactory := aggregate.NewFactory()
	aggregateFactory.RegisterAggregate(createAggregate)

	accountDetailsProjector := eventhandler.New()
	accountDetailsProjector.RegisterHandlers(eh.NewAccountDetailsProjector(accountReporter))

	eventPublisher := eventbus.NewInMemoryEventBus()
	eventPublisher.Register(accountDetailsProjector)

	aggregateStore := store.NewStore(
		eventstore.NewInInMemoryEventStore(eventPublisher),
		aggregateFactory,
	)

	return dispatcher.NewDispatcher(aggregateStore)
}

func createAggregate(ID cqrs.Identifier) cqrs.AdvancedAggregate {
	acc := account.NewAggregate(ID)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(acc)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(acc)

	return aggregate.NewAdvanced(acc, commandHandler, eventApplier)
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
