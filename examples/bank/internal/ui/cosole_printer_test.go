package ui_test

import (
	"testing"

	"github.com/screwyprof/cqrs/pkg/assert"

	"github.com/screwyprof/cqrs/examples/bank/internal/ui"
)

func TestNewConsolePrinter(t *testing.T) {
	t.Run("ItPanicsIfWriterIsNotGiven", func(t *testing.T) {
		factory := func() {
			ui.NewConsolePrinter(nil, nil)
		}
		assert.Panic(t, factory)
	})
}
