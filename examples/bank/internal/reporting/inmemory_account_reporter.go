package reporting

import "github.com/screwyprof/cqrs/examples/bank/pkg/report"

type InMemoryAccountReporter struct{}

func NewInMemoryAccountReporter() *InMemoryAccountReporter {
	return &InMemoryAccountReporter{}
}

func (r *InMemoryAccountReporter) AccountDetailsFor(ID report.Identifier) {
	panic("implement me")
}

func (r *InMemoryAccountReporter) Update(account *report.Account) {
	panic("implement me")
}
