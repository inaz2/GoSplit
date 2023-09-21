package gosplit

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// GeneralizedError represents kind, error and stacktrace.
type GeneralizedError struct {
	kind  error
	err   error
	stack []byte
}

// NewErrorf returns a new GeneralizedError.
func NewErrorf(kind error, format string, a ...any) error {
	err := fmt.Errorf(format, a...)

	var stack []byte
	var e *GeneralizedError
	if errors.As(err, &e) {
		// keep original stacktrace
		stack = e.stack
	} else {
		stack = debug.Stack()
	}

	return &GeneralizedError{kind: kind, err: err, stack: stack}
}

// ErrGoSplit represents a general GoSplit error.
var ErrGoSplit = errors.New("gosplit")

// GoSplitErrorf returns a new GeneralizedError of ErrGoSplit.
func GoSplitErrorf(format string, a ...any) error {
	return NewErrorf(ErrGoSplit, format, a...)
}

// Error implemenrts error.Error.
func (e *GeneralizedError) Error() string {
	return e.err.Error()
}

// Unwrap returns the wrapped errors.
func (e *GeneralizedError) Unwrap() []error {
	return []error{e.kind, e.err}
}

// GoString implemenrts fmt.GoStringer.
func (e *GeneralizedError) GoString() string {
	return fmt.Sprintf("&gosplit.GeneralizedError{kind: %#v, err: %#v, stack: %#v}", e.kind, e.err, e.stack)
}

// Format implements fmt.Formatter, extending "%+v" as error with stacktrace.
func (e *GeneralizedError) Format(f fmt.State, verb rune) {
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
