package gosplit

import (
	. "inaz2/GoSplit/internal/gerrors"

	"errors"
)

// ErrGoSplit represents any errors in this package.
var ErrGoSplit = errors.New("gosplit")

// GoSplitErrorf returns a new Gerror from ErrGoSplit.
func GoSplitErrorf(format string, a ...any) Gerror {
	return GErrorf(ErrGoSplit, format, a...)
}

// Specific errors.
var (
	ErrInvalidBytes    = errors.New("invalid number of bytes")
	ErrInvalidLines    = errors.New("invalid number of lines")
	ErrInvalidNumber   = errors.New("invalid number of chunks")
	ErrUnknownSize     = errors.New("cannot determine file size")
	ErrIsDirectory     = errors.New("is a directory")
	ErrNoFreeSpace     = errors.New("no free space available")
	ErrSuffixExhausted = errors.New("output file suffixes exhausted")
)
