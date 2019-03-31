package eventhandler

import (
	"github.com/screwyprof/cqrs/examples/bank/pkg/event"
	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

// AccountDetailsProjector projects account details to the read side.
type AccountDetailsProjector struct {
	accountReporter report.AccountReporting
}

// NewAccountDetailsProjector creates new instance of AccountDetailsProjector.
func NewAccountDetailsProjector(accountReporter report.AccountReporting) *AccountDetailsProjector {
	if accountReporter == nil {
		panic("accountReporter is required")
	}
	return &AccountDetailsProjector{accountReporter: accountReporter}
}

// OnAccountOpened handles AccountOpened event.
func (p *AccountDetailsProjector) OnAccountOpened(e event.AccountOpened) error {
	p.accountReporter.Save(&report.Account{
		ID:     e.ID,
		Number: e.Number,
	})

	return nil
}

// OnMoneyDeposited handles MoneyDeposited event.
func (p *AccountDetailsProjector) OnMoneyDeposited(e event.MoneyDeposited) error {
	return p.addLedger(Ledger{
		ID:      e.ID,
		Action:  "deposit",
		Balance: e.Balance,
		Amount:  e.Amount,
	})
}

// OnMoneyWithdrawn handles MoneyWithdrawn event.
func (p *AccountDetailsProjector) OnMoneyWithdrawn(e event.MoneyWithdrawn) error {
	return p.addLedger(Ledger{
		ID:      e.ID,
		Action:  "withdraw",
		Balance: e.Balance,
		Amount:  e.Amount,
	})
}

func (p *AccountDetailsProjector) addLedger(l Ledger) error {
	acc, err := p.accountReporter.AccountDetailsFor(l.ID)
	if err != nil {
		return err
	}

	acc.Balance = l.Balance
	acc.Ledgers = append(acc.Ledgers, report.Ledger{Action: l.Action, Amount: l.Amount})

	p.accountReporter.Save(acc)
	return nil
}

// Ledger holds ledger info.
type Ledger struct {
	ID      report.Identifier
	Action  string
	Balance int
	Amount  int
}
