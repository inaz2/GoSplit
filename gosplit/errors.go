package gosplit

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// errGoSplit represents error and stacktrace.
type errGoSplit struct {
	err   error
	stack []byte
}

// ErrGoSplit represents a errGoSplit value for errors.Is().
var ErrGoSplit *errGoSplit

// GoSplitErrorf returns a new errGoSplit.
func GoSplitErrorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...)

	var t *errGoSplit
	if errors.As(err, &t) {
		// keep original stacktrace
		return &errGoSplit{err: err, stack: t.stack}
	} else {
		return &errGoSplit{err: err, stack: debug.Stack()}
	}
}

// Error implemenrts error.Error.
func (e *errGoSplit) Error() string {
	return e.err.Error()
}

// Is returns true if target is errGoSplit.
func (e *errGoSplit) Is(target error) bool {
	_, ok := target.(*errGoSplit)
	return ok
}

// Unwrap returns the wrapped error.
func (e *errGoSplit) Unwrap() error {
	return e.err
}

// Format implements fmt.Formatter, extending "%+v" appends stacktrace.
func (e *errGoSplit) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		var (
			message    string
			stacktrace string
		)

		if f.Flag('#') {
			message = fmt.Sprintf("%#v", e.err)
		} else {
			message = fmt.Sprintf("%v", e.err)
		}
		if f.Flag('+') {
			stacktrace = "\n" + string(e.stack)
		} else {
			stacktrace = ""
		}

		fmt.Fprint(f, message+stacktrace)
	default:
		fmt.Fprintf(f, "%"+string(verb), e.err)
	}
}
