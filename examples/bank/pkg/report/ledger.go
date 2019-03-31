package report

import "fmt"

type Ledger struct {
	Action  string
	Amount  int64
	Balance int64
}

func FormatLedger(no int, l Ledger) string {
	amount := l.Amount
	if l.Action == "withdrawal" {
		amount *= -1
	}
	return fmt.Sprintf("%d | %8.2f | %8.2f\n", no, float32(amount), float32(l.Balance))
}
