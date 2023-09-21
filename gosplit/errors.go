package gosplit

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// genericError represents kind, error and stacktrace.
type genericError struct {
	kind  error
	err   error
	stack []byte
}

// ErrGoSplit represents a general GoSplit error.
var ErrGoSplit = errors.New("gosplit error")

// GoSplitErrorf returns a new genericError with ErrGoSplit.
func GoSplitErrorf(format string, a ...any) error {
	kind := ErrGoSplit
	err := fmt.Errorf(format, a...)

	var stack []byte
	var e *genericError
	if errors.As(err, &e) {
		// keep original stacktrace
		stack = e.stack
	} else {
		stack = debug.Stack()
	}

	return &genericError{kind: kind, err: err, stack: stack}
}

// Error implemenrts error.Error.
func (e *genericError) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped errors.
func (e *genericError) Unwrap() []error {
	return []error{e.kind, e.err}
}

// Format implements fmt.Formatter, extending "%+v" as error with stacktrace.
func (e *genericError) Format(f fmt.State, verb rune) {
	format := fmt.FormatString(f, verb)
	msg := fmt.Sprintf(format, e.err)
	if verb == 'v' && f.Flag('+') {
		msg += "\n" + string(e.stack)
	}
	fmt.Fprint(f, msg)
}
