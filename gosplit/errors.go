package gosplit

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// GoSplitError represents error and stacktrace.
type GoSplitError struct {
	err   error
	stack []byte
}

// GoSplitErrorf returns a new GoSplitError.
func GoSplitErrorf(format string, a ...any) *GoSplitError {
	err := fmt.Errorf(format, a...)

	var t *GoSplitError
	if errors.As(err, &t) {
		// keep original stacktrace
		return &GoSplitError{err: err, stack: t.stack}
	} else {
		return &GoSplitError{err: err, stack: debug.Stack()}
	}
}

// Error implemenrts error.Error.
func (e *GoSplitError) Error() string {
	return e.err.Error()
}

// Format implements fmt.Formatter, implemented "%+v" with stacktrace.
func (e *GoSplitError) Format(f fmt.State, verb rune) {
	switch verb {
	case 'v':
		if f.Flag('+') {
			fmt.Fprintf(f, "%v\n%s", e.err, e.stack)
		} else {
			fmt.Fprintf(f, "%v", e.err)
		}
	case 's', 'q', 'x', 'X':
		fmt.Fprintf(f, "%"+string(verb), e.err)
	}
}
