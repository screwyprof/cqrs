package eventhandler

import "github.com/screwyprof/cqrs/examples/bank/pkg/report"

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
