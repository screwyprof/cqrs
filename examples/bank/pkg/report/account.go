package report

// Account an Account representation.
type Account struct {
	ID      Identifier
	Number  string
	Balance int64
	Ledgers []Ledger
}
