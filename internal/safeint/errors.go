package safeint

import (
	. "inaz2/GoSplit/internal/gerrors"

	"errors"
)

// ErrSafeInt represents a error in this package.
var ErrSafeInt = errors.New("safeint")

// SafeIntErrorf returns a new Gerror from ErrSafeInt.
func SafeIntErrorf(format string, a ...any) Gerror {
	return GErrorf(ErrSafeInt, format, a...)
}

var ErrOverflow = errors.New("integer overflow occured")
var ErrDivisionByZero = errors.New("division by zero")
