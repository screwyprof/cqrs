package event

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type MoneyDeposited struct {
	ID      domain.Identifier
	Amount  int64
	Balance int64
}

func (e MoneyDeposited) EventType() string {
	return "MoneyDeposited"
}
