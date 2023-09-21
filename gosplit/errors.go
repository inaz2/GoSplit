package gosplit

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// ErrorWithStack represents error and stacktrace.
type ErrorWithStack struct {
	err   error
	stack []byte
}

// WrapWithStack returns a new ErrorWithStack for err.
func WrapWithStack(err error) error {
	var stack []byte
	var e *ErrorWithStack
	if errors.As(err, &e) {
		// keep original stacktrace
		stack = e.stack
	} else {
		stack = debug.Stack()
	}
	return &ErrorWithStack{err: err, stack: stack}
}

// ErrGoSplit represents a general GoSplit error.
var ErrGoSplit = errors.New("gosplit")

// GoSplitErrorf returns a new ErrorWithStack for a error joined with ErrGoSplit.
func GoSplitErrorf(format string, a ...any) error {
	e := fmt.Errorf(format, a...)
	e = fmt.Errorf("%w: %w", ErrGoSplit, e)
	return WrapWithStack(e)
}

// Error implemenrts error.Error.
func (e *ErrorWithStack) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped errors.
func (e *ErrorWithStack) Unwrap() error {
	return e.err
}

// GoString implemenrts fmt.GoStringer.
func (e *ErrorWithStack) GoString() string {
	return fmt.Sprintf("&gosplit.ErrorWithStack{err: %#v, stack: %#v}", e.err, e.stack)
}

// Format implements fmt.Formatter, extending "%+v" as error with stacktrace.
func (e *ErrorWithStack) Format(f fmt.State, verb rune) {
	var msg string
	if verb == 'v' && f.Flag('#') {
		msg = e.GoString()
	} else {
		format := fmt.FormatString(f, verb)
		msg = fmt.Sprintf(format, e.err)
	}
	if verb == 'v' && f.Flag('+') {
		msg += "\n" + string(e.stack)
	}
	fmt.Fprint(f, msg)
}
