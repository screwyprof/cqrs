package cqrs

type Transactional interface {
	BeginTransaction()
	Commit() error
	Rollback() error
}

type TransactionalEventStorage interface {
	EventStore
	Transactional
}
