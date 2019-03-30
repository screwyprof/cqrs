package eventhandler

import "github.com/screwyprof/cqrs/examples/bank/pkg/report"

type AccountDetailsProjector struct {
	accountReporter report.AccountReporting
}

func NewAccountDetailsProjector(accountReporter report.AccountReporting) *AccountDetailsProjector {
	return &AccountDetailsProjector{accountReporter: accountReporter}
}
