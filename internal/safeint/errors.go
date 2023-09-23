package safeint

import (
	. "inaz2/GoSplit/internal/gerrors"

	"errors"
)

// ErrSafeInt represents any errors in this package.
var ErrSafeInt = errors.New("safeint")

// SafeIntErrorf returns a new Gerror from ErrSafeInt.
func SafeIntErrorf(format string, a ...any) Gerror {
	return GErrorf(ErrSafeInt, format, a...)
}

// Specific errors.
var (
	ErrOverflow       = errors.New("integer overflow occured")
	ErrDivisionByZero = errors.New("division by zero")
)
