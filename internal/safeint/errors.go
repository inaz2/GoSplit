package safeint

import (
	g "inaz2/GoSplit/internal/gerrors"

	"errors"
)

// ErrSafeInt represents any errors in this package.
var ErrSafeInt = errors.New("safeint")

// Specific errors.
var (
	ErrOverflow       = errors.New("integer overflow occured")
	ErrDivisionByZero = errors.New("division by zero")
)

// wrapper is a error wrapper for this package.
var wrapper = g.NewWrapper(ErrSafeInt)
