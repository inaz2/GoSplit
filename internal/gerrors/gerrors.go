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

// Error represents the interface extending error. Formatting "%+v" as error with stacktrace.
//
// Intended to use Error instead of error for type checking.
type Error interface {
	error
	value() *errorWithStack
}

// Error implements errors.error interface.
func (e *errorWithStack) Error() string {
	return e.err.Error()
}

// value implements Error interface, requires that its type is *errorWithStack.
func (e *errorWithStack) value() *errorWithStack {
	return e
}

// Unwrap returns the wrapped errors.
func (e *errorWithStack) Unwrap() error {
	return e.err
}

// GoString implements fmt.GoStringer.
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

// Wrapper provides the methods for wrapping an error with base error. The error string of base error is discarded.
//
// Intended to use Wrapper.Errorf instead of fmt.Errorf.
type Wrapper struct {
	err error
}

// NewWrapper returns a new Wrapper.
func NewWrapper(err error) *Wrapper {
	return &Wrapper{err: err}
}

// Errorf returns a new Error by formatting.
func (f *Wrapper) Errorf(format string, a ...any) Error {
	err := fmt.Errorf(format, a...)
	return f.Link(err, f.err)
}

// Link returns a new Error linked to errOld. The error string of errOld is discarded.
func (f *Wrapper) Link(errNew error, errOld error) Error {
	// prepend errOld by zero-length formatting "%.w"
	err := fmt.Errorf("%.w%w", errOld, errNew)

	// ensure that a error wraps the base error
	if !errors.Is(err, f.err) {
		err = fmt.Errorf("%w%.w", err, f.err)
	}

	var stack []byte
	var tmp *errorWithStack
	if errors.As(err, &tmp) {
		// keep original stacktrace
		stack = tmp.stack
	} else {
		stack = debug.Stack()
	}

	return &errorWithStack{err: err, stack: stack}
}
