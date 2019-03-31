package account

import "github.com/screwyprof/cqrs/examples/bank/pkg/domain"

// Aggregate handles operations with an account.
type Aggregate struct {
	id domain.Identifier
}

// NewAggregate creates a new instance of *Aggregate.
func NewAggregate(ID domain.Identifier) *Aggregate {
	if ID == nil {
		panic("ID required")
	}
	return &Aggregate{id: ID}
}

// AggregateID returns aggregate ID.
func (a *Aggregate) AggregateID() domain.Identifier {
	return a.id
}

// AggregateType return aggregate type.
func (*Aggregate) AggregateType() string {
	panic("implement me")
}
