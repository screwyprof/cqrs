package ui

import (
	"io"

	"github.com/screwyprof/cqrs/examples/bank/pkg/report"
)

// ConsolePrinter prints account statement to console.
type ConsolePrinter struct {
	w               io.Writer
	accountReporter report.GetAccountDetails
}

// NewConsolePrinter creates new instance of ConsolePrinter.
func NewConsolePrinter(w io.Writer, accountReporter report.GetAccountDetails) *ConsolePrinter {
	if w == nil {
		panic("writer is required")
	}
	return &ConsolePrinter{w: w, accountReporter: accountReporter}
}

// PrintAccountStatement prints account statement to console.
func (p *ConsolePrinter) PrintAccountStatement() {
	//accountDetailedReport := accountReporter.AccountDetailsFor(ID)
	panic("Implement me")
}
