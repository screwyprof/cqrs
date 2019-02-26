package cqrs

type ComplexAggregate interface {
	EventProvider
	CommandHandler
}
