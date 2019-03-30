package mock

import "errors"

var (
	ErrAggIsNotRegistered = errors.New("mock.TestAggregate is not registered")
)