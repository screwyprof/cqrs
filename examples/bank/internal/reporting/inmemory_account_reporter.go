package reporting

import (
	"errors"
	"sync"

	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

var (
	ErrAccountNotFound = errors.New("account is not found")
)

type InMemoryAccountReporter struct {
	accounts  map[report.Identifier]*report.Account
	accountMu sync.RWMutex
}

func NewInMemoryAccountReporter() *InMemoryAccountReporter {
	return &InMemoryAccountReporter{
		accounts: make(map[report.Identifier]*report.Account),
	}
}

func (r *InMemoryAccountReporter) AccountDetailsFor(ID report.Identifier) (*report.Account, error) {
	r.accountMu.RLock()
	defer r.accountMu.RUnlock()
	return r.accounts[ID], nil
}

func (r *InMemoryAccountReporter) Save(account *report.Account) {
	r.accountMu.Lock()
	defer r.accountMu.Unlock()
	r.accounts[account.ID] = account
}
