// Package gerrors implements generalized error with stacktrace.
package gerrors

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// errorWithStack represents error and stacktrace.
type errorWithStack struct {
	err   error
	stack []byte
}

// Gerror represents the interface extending error. Formatting "%+v" as error with stacktrace.
//
// Intended to use Gerror instead of error for type checking.
type Gerror interface {
	error
	GError() *errorWithStack
}

// GErrorf returns a new Gerror from err by formatting. The error string of err is discarded.
func GErrorf(err error, format string, a ...any) Gerror {
	e := fmt.Errorf(format, a...)
	return GLink(e, err)
}

// GLink returns a new Gerror from err2 and links it to err1. The error string of err1 is discarded.
func GLink(err2 error, err1 error) Gerror {
	// prepend err1 by zero-length formatting "%.w"
	e := fmt.Errorf("%.w%w", err1, err2)

	var stack []byte
	var tmp *errorWithStack
	if errors.As(e, &tmp) {
		// keep original stacktrace
		stack = tmp.stack
	} else {
		stack = debug.Stack()
	}

	return &errorWithStack{err: e, stack: stack}
}

// Error implements errors.error interface.
func (e *errorWithStack) Error() string {
	return e.err.Error()
}

// GError implements Gerror interface, that requires errorWithStack.
func (e *errorWithStack) GError() *errorWithStack {
	return e
}

// Unwrap returns the wrapped errors.
func (e *errorWithStack) Unwrap() error {
	return e.err
}

// GoString implemenrts fmt.GoStringer.
func (e *errorWithStack) GoString() string {
	return fmt.Sprintf("&gerrors.errorWithStack{err: %#v, stack: %#v}", e.err, e.stack)
}

// Format implements fmt.Formatter, extending "%+v" as error with stacktrace.
func (e *errorWithStack) Format(f fmt.State, verb rune) {
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
