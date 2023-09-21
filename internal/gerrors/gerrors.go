package gerrors

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// generalizedError represents kind, error and stacktrace.
type generalizedError struct {
	kind  error
	err   error
	stack []byte
}

// Errorf returns a new generalizedError.
func Errorf(kind error, format string, a ...any) error {
	err := fmt.Errorf(format, a...)

	var stack []byte
	var e *generalizedError
	if errors.As(err, &e) {
		// keep original stacktrace
		stack = e.stack
	} else {
		stack = debug.Stack()
	}

	return &generalizedError{kind: kind, err: err, stack: stack}
}

// Error implemenrts error.Error.
func (e *generalizedError) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped errors.
func (e *generalizedError) Unwrap() []error {
	return []error{e.kind, e.err}
}

// GoString implemenrts fmt.GoStringer.
func (e *generalizedError) GoString() string {
	return fmt.Sprintf("&gerrors.generalizedError{kind: %#v, err: %#v, stack: %#v}", e.kind, e.err, e.stack)
}

// Format implements fmt.Formatter, extending "%+v" as error with stacktrace.
func (e *generalizedError) Format(f fmt.State, verb rune) {
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
