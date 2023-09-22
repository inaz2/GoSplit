package gosplit

import (
	. "inaz2/GoSplit/internal/gerrors"

	"errors"
)

// ErrGoSplit represents a error in this package.
var ErrGoSplit = errors.New("gosplit")

// GoSplitErrorf returns a new Gerror from ErrGoSplit.
func GoSplitErrorf(format string, a ...any) Gerror {
	return GErrorf(ErrGoSplit, format, a...)
}

var ErrInvalidBytes = errors.New("invalid number of bytes")
var ErrInvalidLines = errors.New("invalid number of lines")
var ErrInvalidNumber = errors.New("invalid number of chunks")
var ErrUnknownSize = errors.New("cannot determine file size")
var ErrIsDirectory = errors.New("is a directory")
var ErrNoFreeSpace = errors.New("no free space available")
var ErrSuffixExhausted = errors.New("output file suffixes exhausted")
