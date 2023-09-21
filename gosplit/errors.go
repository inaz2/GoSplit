package gosplit

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// errorWithStack represents kind, error and stacktrace.
type errorWithStack struct {
	err   error
	stack []byte
}

// ErrGoSplit represents a general GoSplit error.
var ErrGoSplit = errors.New("gosplit")

// GoSplitErrorf returns a new errorWithStack with ErrGoSplit.
func GoSplitErrorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	err = fmt.Errorf("%w: %w", ErrGoSplit, err)

	var stack []byte
	var e *errorWithStack
	if errors.As(err, &e) {
		// keep original stacktrace
		stack = e.stack
	} else {
		stack = debug.Stack()
	}

	return &errorWithStack{err: err, stack: stack}
}

// Error implemenrts error.Error.
func (e *errorWithStack) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped errors.
func (e *errorWithStack) Unwrap() error {
	return e.err
}

// Format implements fmt.Formatter, extending "%+v" as error with stacktrace.
func (e *errorWithStack) Format(f fmt.State, verb rune) {
	format := fmt.FormatString(f, verb)
	msg := fmt.Sprintf(format, e.err)
	if verb == 'v' && f.Flag('+') {
		msg += "\n" + string(e.stack)
	}
	fmt.Fprint(f, msg)
}
