package testdsl

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
)

// GivenFn is a test init function.
type GivenFn func() (cqrs.ESAggregate, []cqrs.DomainEvent)

// WhenFn is a command handler function.
type WhenFn func(agg cqrs.ESAggregate, err error) ([]cqrs.DomainEvent, error)

// ThenFn prepares the Checker.
type ThenFn func(t *testing.T) Checker

// Checker asserts the given results.
type Checker func(got []cqrs.DomainEvent, err error)

// AggregateTester defines an aggregate tester.
type AggregateTester func(given GivenFn, when WhenFn, then ThenFn)

// Test runs the test.
//
// Example:
//
//	 Test(t)(
//		  Given(agg),
//		  When(testdata.TestCommand{Param: "param"}),
//		  Then(testdata.TestEvent{Data: "param"}),
//	 )
func Test(t *testing.T) AggregateTester {
	return func(given GivenFn, when WhenFn, then ThenFn) {
		t.Helper()
		then(t)(when(applyEvents(given)))
	}
}

// Given prepares the given aggregate for testing.
func Given(agg cqrs.ESAggregate, events ...cqrs.DomainEvent) GivenFn {
	return func() (cqrs.ESAggregate, []cqrs.DomainEvent) {
		return agg, events
	}
}

// When prepares the command handler for the given command.
func When(c cqrs.Command) WhenFn {
	return func(agg cqrs.ESAggregate, err error) ([]cqrs.DomainEvent, error) {
		if err != nil {
			return nil, err
		}
		return agg.Handle(c)
	}
}

// Then asserts that the expected events are applied.
func Then(want ...cqrs.DomainEvent) ThenFn {
	return func(t *testing.T) Checker {
		return func(got []cqrs.DomainEvent, err error) {
			t.Helper()
			assert.NoError(t, err)
			assert.Equal(t, want, got)
		}
	}
}

// ThenFailWith asserts that the expected error occurred.
func ThenFailWith(want error) ThenFn {
	return func(t *testing.T) Checker {
		return func(got []cqrs.DomainEvent, err error) {
			t.Helper()
			assert.Equal(t, want, err)
		}
	}
}

func applyEvents(given GivenFn) (cqrs.ESAggregate, error) {
	agg, events := given()
	err := agg.Apply(events...)
	if err != nil {
		return nil, err
	}

	return agg, nil
}
