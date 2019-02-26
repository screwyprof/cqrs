package cqrs

type UnitOfWork interface {
	Commit() error
	Rollback() error
}
