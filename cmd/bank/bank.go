package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/screwyprof/cqrs"
	"github.com/screwyprof/cqrs/commandhandler/bus"
	eventbus "github.com/screwyprof/cqrs/eventbus/memory"
	eventstore "github.com/screwyprof/cqrs/eventstore/memory"
	"github.com/screwyprof/cqrs/identitymap"
	"github.com/screwyprof/cqrs/middleware/commandhandler/transactional"
	repository "github.com/screwyprof/cqrs/repository/eventsourced"

	"github.com/screwyprof/cqrs/example/bank/command"
	"github.com/screwyprof/cqrs/example/bank/domain/account"
)

func main() {
	cqrs.RegisterAggregate(
		func(ID uuid.UUID) cqrs.ComplexAggregate {
			return account.Construct(ID)
		},
	)

	eventBus := eventbus.NewEventBus()
	identityMap := identitymap.NewIdentityMap()
	domainEventStorage := eventstore.NewEventStore()
	domainRepository := repository.NewRepository(domainEventStorage, identityMap, eventBus)

	commandBus := bus.NewCommandHandler(domainRepository)
	commandBus = cqrs.UseCommandHandlerMiddleware(commandBus, transactional.NewMiddleware(domainRepository))

	accID := uuid.New()
	err := commandBus.Handle(command.OpenAccount{AggID: accID, Number: "ACC777"})
	failOnError(err)

	err = commandBus.Handle(command.DepositMoney{AggID: accID, Amount: 500})
	failOnError(err)

	//acc, err := domainRepository.ByID(accID, "account.Account")
	//failOnError(err)
	//spew.Dump(acc)
}

func failOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
