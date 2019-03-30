package cqrs

import (
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/screwyprof/cqrs/pkg/assert"
)

type testEvent struct{}

func (e testEvent) EventType() string {
	return "testEvent"
}

func TestMatchAnyEventOf(t *testing.T) {
	t.Run("ItShouldReturnFalseIfNoEventsAreMatched", func(t *testing.T) {
		m := MatchAnyEventOf(faker.Word(), faker.Word())

		assert.True(t, m(testEvent{}) == false)
	})
}
