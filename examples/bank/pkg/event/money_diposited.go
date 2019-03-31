package event

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type MoneyDeposited struct {
	ID      domain.Identifier
	Amount  int
	Balance int
}

func (e MoneyDeposited) EventType() string {
	return "MoneyDeposited"
}
