// Package gerrors implements generalized error with stacktrace.
package gerrors

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

// Errorf returns a new errBase error with formatting. The error string of errBase is discarded.
func Errorf(errBase error, format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	return Link(err, errBase)
}

// Link returns a new error linked to err1. The error string of err1 is discarded.
func Link(err2 error, err1 error) error {
	// prepend err1 by zero-length formatting "%.w"
	err := fmt.Errorf("%.w%w", err1, err2)

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
	return fmt.Sprintf("&gerrors.ErrorWithStack{err: %#v, stack: %#v}", e.err, e.stack)
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
