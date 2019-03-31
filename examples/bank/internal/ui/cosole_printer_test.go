package ui_test

import (
	"bytes"
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

	t.Run("ItPanicsIfAccountReporterIsNotGiven", func(t *testing.T) {
		factory := func() {
			ui.NewConsolePrinter(&bytes.Buffer{}, nil)
		}
		assert.Panic(t, factory)
	})
}
