package gosplit

import (
	g "inaz2/GoSplit/internal/gerrors"

	"errors"
)

// ErrGoSplit represents any errors in this package.
var ErrGoSplit = errors.New("gosplit")

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

// wrapper is a error wrapper for this package.
var wrapper = g.NewWrapper(ErrGoSplit)
