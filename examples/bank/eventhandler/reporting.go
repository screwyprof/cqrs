package eventhandler

import "github.com/screwyprof/cqrs/examples/bank/report"

// GetAccountDetails returns detailed account info.
type GetAccountDetails interface {
	AccountDetailsFor(ID report.Identifier) (*report.Account, error)
}

// AccountSaver saves account info.
type AccountSaver interface {
	Save(account *report.Account)
}

// AccountReporting handles account read side.
type AccountReporting interface {
	GetAccountDetails
	AccountSaver
}
