package report

import "fmt"

// Identifier an object identifier.
type Identifier = fmt.Stringer

// GetAccountDetails returns detailed account info.
type GetAccountDetails interface {
	AccountDetailsFor(ID Identifier) (*Account, error)
}

// AccountSaver saves account info.
type AccountSaver interface {
	Save(account *Account)
}

// AccountReporting handles account read side.
type AccountReporting interface {
	GetAccountDetails
	AccountSaver
}
