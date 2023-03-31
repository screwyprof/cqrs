package report

import "fmt"

// Identifier an object identifier.
type Identifier = fmt.Stringer

// Account an Account representation.
type Account struct {
	ID      Identifier
	Number  string
	Balance int64
	Ledgers []Ledger
}
