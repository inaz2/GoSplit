package gosplit

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// goSplitError represents error and stacktrace.
type goSplitError struct {
	err   error
	stack []byte
}

// GoSplitErr represents a goSplitError value for errors.Is().
var GoSplitErr *goSplitError

// GoSplitErrorf returns a new goSplitError.
func GoSplitErrorf(format string, a ...any) error {
	err := fmt.Errorf(format, a...)

	var t *goSplitError
	if errors.As(err, &t) {
		// keep original stacktrace
		return &goSplitError{err: err, stack: t.stack}
	} else {
		return &goSplitError{err: err, stack: debug.Stack()}
	}
}

// Error implemenrts error.Error.
func (e *goSplitError) Error() string {
	return e.err.Error()
}

// Is returns true if target is goSplitError.
func (e *goSplitError) Is(target error) bool {
	_, ok := target.(*goSplitError)
	return ok
}

// Unwrap returns the wrapped error.
func (e *goSplitError) Unwrap() error {
	return e.err
}

// Format implements fmt.Formatter, extending "%+v" appends stacktrace.
func (e *goSplitError) Format(f fmt.State, verb rune) {
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
