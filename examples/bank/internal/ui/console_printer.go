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

	if accountReporter == nil {
		panic("accountReporter is required")
	}
	return &ConsolePrinter{w: w, accountReporter: accountReporter}
}

// PrintAccountStatement prints account statement to console.
//
// Sample output:
// Account #ACC777:
// # |  Amount | Balance
// 1 | 1000.00 | 1000.00
// 2 | -100.00 | 900.00
// 3 |  500.00 | 1400.00
func (p *ConsolePrinter) PrintAccountStatement() {
	//accountDetailedReport := accountReporter.AccountDetailsFor(ID)
	panic("Implement me")
}
