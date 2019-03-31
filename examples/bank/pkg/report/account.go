package report

// Account An Account representation.
type Account struct {
	ID      Identifier
	Number  string
	Balance int
	Ledgers []Ledger
}
