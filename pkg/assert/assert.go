package assert

import (
	"reflect"
	"testing"
)

// True fails the test if the condition is false.
func True(tb testing.TB, condition bool) {
	tb.Helper()
	if !condition {
		tb.Errorf("\033[31mcondition is false\033[39m\n\n")
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatalf("\033[31munexpected error: %v\033[39m\n\n", err)
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(exp, act) {
		tb.Errorf("\033[31m\n\n\texp:\n\t%#+v\n\n\tgot:\n\t%#+v\033[39m", exp, act)
	}
}

// Panic fails the test if it didn't panic.
func Panic(tb testing.TB, f func()) {
	tb.Helper()
	defer func() {
		tb.Helper()
		if r := recover(); r == nil {
			tb.Errorf("\033[31mpanic is expected\033[39m")
		}
	}()
	f()
}
