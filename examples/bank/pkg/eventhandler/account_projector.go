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

func (p *AccountDetailsProjector) OnAccountOpened(e event.AccountOpened) error {
	p.accountReporter.Save(&report.Account{
		ID:     e.ID,
		Number: e.Number,
	})

	return nil
}

func (p *AccountDetailsProjector) OnMoneyDeposited(e event.MoneyDeposited) error {
	acc, err := p.accountReporter.AccountDetailsFor(e.ID)
	if err != nil {
		return err
	}

	acc.Balance = e.Balance
	acc.Ledgers = append(acc.Ledgers, report.Ledger{Action: "deposit", Amount: e.Amount})

	p.accountReporter.Save(acc)
	return nil
}

func (p *AccountDetailsProjector) OnMoneyWithdrawn(e event.MoneyWithdrawn) error {
	acc, err := p.accountReporter.AccountDetailsFor(e.ID)
	if err != nil {
		return err
	}

	acc.Balance = e.Balance
	acc.Ledgers = append(acc.Ledgers, report.Ledger{Action: "withdraw", Amount: e.Amount})

	p.accountReporter.Save(acc)
	return nil
}
