package aggregatestore

//existingAggID := uuid.New()
//acc := account.Construct(existingAggID)
//
//accountOpened := event.NewAccountOpened(existingAggID, "ACC777")
//accountOpened.SetVersion(1)
//
//moneyDeposited := event.NewMoneyDeposited(existingAggID, 100, 100)
//moneyDeposited.SetVersion(2)
//
//moneyDeposited2 := event.NewMoneyDeposited(existingAggID, 50, 150)
//moneyDeposited2.SetVersion(3)
//
//err := acc.LoadFromHistory([]domain.DomainEvent{
//	accountOpened,
//	moneyDeposited,
//	moneyDeposited2,
//	//event.NewMoneyDeposited(existingAggID, 100, 100),
//	//event.NewMoneyDeposited(existingAggID, 50, 150),
//})
//failOnError(err)

//spew.Dump(acc.AggregateID())
//spew.Dump(acc.Version())
//spew.Dump(acc.UncommittedChanges())

//acc = account.Construct(uuid.New())
//err = acc.Handle(command.OpenAccount{Number:"ACC777"})
//failOnError(err)
//
//err = acc.Handle(command.DepositMoney{Amount:100})
//failOnError(err)
