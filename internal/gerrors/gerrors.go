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

// Format implements fmt.Formatter, extending "%+v" and "%#+v" as error with stacktrace.
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

// Wrapper provides the methods for wrapping an error with base error.
//
// Intended to use Wrapper.Errorf instead of fmt.Errorf.
type Wrapper struct {
	errBase error
}

// NewWrapper returns a new Wrapper of err. The error string of base error is discarded.
func NewWrapper(err error) *Wrapper {
	return &Wrapper{errBase: err}
}

// Errorf returns a new Error by formatting.
func (w *Wrapper) Errorf(format string, a ...any) Error {
	err := fmt.Errorf(format, a...)
	return w.Link(err, w.errBase)
}

// Link returns a new Error linked to errOld. The error string of errOld is discarded.
func (w *Wrapper) Link(errNew error, errOld error) Error {
	// append errOld by zero-length format specifier "%.w"
	// because errNew is expected to be handled earlier
	err := fmt.Errorf("%w%.w", errNew, errOld)

	// prepend a base error if err is not wrapped by it
	// because a base error is expected to be handled earlier
	if !errors.Is(err, w.errBase) {
		err = fmt.Errorf("%.w%w", w.errBase, err)
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
