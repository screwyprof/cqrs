package event

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

type MoneyWithdrawn struct {
	ID      domain.Identifier
	Amount  int
	Balance int
}

func (e MoneyWithdrawn) EventType() string {
	return "MoneyWithdrawn"
}
