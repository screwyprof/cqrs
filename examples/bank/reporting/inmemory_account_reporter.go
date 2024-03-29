package reporting

import (
	"errors"
	"sync"

	"github.com/screwyprof/cqrs/examples/bank/report"
)

var ErrAccountNotFound = errors.New("account is not found")

// InMemoryAccountReporter stores and retrieves account reports from memory.
type InMemoryAccountReporter struct {
	accounts  map[report.Identifier]*report.Account
	accountMu sync.RWMutex
}

// NewInMemoryAccountReporter creates a new instance of InMemoryAccountReporter.
func NewInMemoryAccountReporter() *InMemoryAccountReporter {
	return &InMemoryAccountReporter{
		accounts: make(map[report.Identifier]*report.Account),
	}
}

// AccountDetailsFor implements report.GetAccountDetails interface.
func (r *InMemoryAccountReporter) AccountDetailsFor(ID report.Identifier) (*report.Account, error) {
	r.accountMu.RLock()
	defer r.accountMu.RUnlock()

	if acc, ok := r.accounts[ID]; ok {
		return acc, nil
	}

	return nil, ErrAccountNotFound
}

// Save implements AccountSaver.interface.
func (r *InMemoryAccountReporter) Save(account *report.Account) {
	r.accountMu.Lock()
	defer r.accountMu.Unlock()
	r.accounts[account.ID] = account
}
