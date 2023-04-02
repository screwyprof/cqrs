package cqrs_test

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/screwyprof/cqrs"
)

type testEvent struct{}

func (e testEvent) EventType() string {
	return "testEvent"
}

func TestMatcher(t *testing.T) {
	t.Parallel()

	t.Run("it matches any event", func(t *testing.T) {
		t.Parallel()

		t.Run("no matches found", func(t *testing.T) {
			t.Parallel()

			m := cqrs.MatchAnyEventOf(faker.Word(), faker.Word())

			assert.False(t, m(testEvent{}))
		})

		t.Run("matches found", func(t *testing.T) {
			t.Parallel()

			m := cqrs.MatchAnyEventOf("testEvent")

			assert.True(t, m(testEvent{}))
		})
	})
}
