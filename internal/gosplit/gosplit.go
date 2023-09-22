// Package gosplit implements file splitting.
package gosplit

import (
	. "inaz2/GoSplit/internal/gerrors"
	"inaz2/GoSplit/internal/safeint"

	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
)

// ErrGoSplit represents a error in this package.
var ErrGoSplit = errors.New("gosplit")

// GoSplitErrorf returns a new Gerror from ErrGoSplit.
func GoSplitErrorf(format string, a ...any) Gerror {
	return GErrorf(ErrGoSplit, format, a...)
}

// GoSplit provides the methods for splitting the file.
type GoSplit struct {
	filePath         string
	prefix           string
	outDir           string
	wVerbose         io.Writer
	bNumericSuffix   bool
	bElideEmptyFiles bool
}

// New returns a new GoSplit struct.
func New(filePath string, prefix string) *GoSplit {
	return &GoSplit{
		filePath: filePath,
		prefix:   prefix,
		outDir:   "./",
		wVerbose: io.Discard,
	}
}

// SetNumericSuffix changes bVerbose flag.
func (g *GoSplit) SetVerboseWriter(w io.Writer) {
	g.wVerbose = w
}

// SetNumericSuffix changes bNumericSuffix flag.
func (g *GoSplit) SetNumericSuffix(bNumericSuffix bool) {
	g.bNumericSuffix = bNumericSuffix
}

// SetNumericSuffix changes bElideEmptyFiles flag.
func (g *GoSplit) SetElideEmptyFiles(bElideEmptyFiles bool) {
	g.bElideEmptyFiles = bElideEmptyFiles
}

// SetOutDir changes the directory of output files.
//
// This method is mainly for testing.
func (g *GoSplit) SetOutDir(outDir string) {
	g.outDir = outDir
}

// ParseSize converts strSize to nBytes, e.g. "10K" -> 10 * 1024.
func (g *GoSplit) ParseSize(strSize string) (int64, Gerror) {
	re := regexp.MustCompile(`^(\d+)(b|(\w)(iB|B)?)?$`)
	m := re.FindStringSubmatch(strSize)
	if m == nil {
		return 0, GoSplitErrorf("invalid number of bytes: %#v", strSize)
	}

	x, err := strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return 0, GoSplitErrorf("invalid number of bytes: %#v", strSize)
	}

	var (
		base       int64
		multiplier int64
	)

	switch m[4] {
	case "B":
		base = 1000
	case "iB":
		base = 1024
	default:
		base = 1024
	}

	switch m[2] {
	case "":
		multiplier = 1
	case "b":
		multiplier = 512
	default:
		exponentMap := map[string]int64{
			"E": 6, "G": 3, "K": 1, "k": 1, "M": 2, "m": 2,
			"P": 5, "Q": 10, "R": 9, "T": 4, "Y": 8, "Z": 7,
		}
		exponent, ok := exponentMap[m[3]]
		if !ok {
			return 0, GoSplitErrorf("invalid number of bytes: %#v", strSize)
		}
		multiplier, err = safeint.PowInt64(base, exponent)
		if err != nil {
			e := GoSplitErrorf("invalid number of bytes: %#v: Value too large for defined data type", strSize)
			return 0, GLink(e, err)
		}
	}

	n, err := safeint.MulInt64(x, multiplier)
	if err != nil {
		e := GoSplitErrorf("invalid number of bytes: %#v: Value too large for defined data type", strSize)
		return 0, GLink(e, err)
	}

	if n <= 0 {
		return 0, GoSplitErrorf("invalid number of bytes: %#v: Numerical result out of range", strSize)
	}

	return n, nil
}

// ByLines splits the content of filePath by nLines.
func (g *GoSplit) ByLines(nLines int) Gerror {
	if nLines <= 0 {
		return GoSplitErrorf("invalid number of lines: %#v", nLines)
	}

	var rFile *os.File
	if g.filePath == "-" {
		rFile = os.Stdin
	} else {
		f, err := os.Open(g.filePath)
		if err != nil {
			return GoSplitErrorf("failed to open: %w", err)
		}
		defer f.Close()

		rFile = f
		if _, err := g.checkFileSize(rFile); err != nil {
			return err
		}
	}

	if err := g.doByLines(rFile, nLines); err != nil {
		return err
	}

	return nil
}

// ByNumber splits the content of filePath into nNumber files.
func (g *GoSplit) ByNumber(nNumber int) Gerror {
	if nNumber <= 0 {
		return GoSplitErrorf("invalid number of chunks: %#v", nNumber)
	}

	var rFile *os.File
	if g.filePath == "-" {
		// print error message when filePath is stdin
		return GoSplitErrorf("cannot determine file size")
	} else {
		f, err := os.Open(g.filePath)
		if err != nil {
			return GoSplitErrorf("failed to open: %w", err)
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return GoSplitErrorf("failed to stat: %w", err)
		}
		mode := fi.Mode()
		if !mode.IsRegular() {
			return GoSplitErrorf("cannot determine file size")
		}

		rFile = f
	}

	fileSize, err := g.checkFileSize(rFile)
	if err != nil {
		return err
	}

	if g.bElideEmptyFiles && fileSize == 0 {
		return nil
	}

	if err := g.doByNumber(rFile, fileSize, nNumber); err != nil {
		return err
	}

	return nil
}

// ByBytes splits the content of filePath by nBytes.
func (g *GoSplit) ByBytes(nBytes int64) Gerror {
	if nBytes <= 0 {
		return GoSplitErrorf("invalid number of bytes: %#v", nBytes)
	}

	var rFile *os.File
	if g.filePath == "-" {
		rFile = os.Stdin
	} else {
		f, err := os.Open(g.filePath)
		if err != nil {
			return GoSplitErrorf("failed to open: %w", err)
		}
		defer f.Close()

		rFile = f
		if _, err := g.checkFileSize(rFile); err != nil {
			return err
		}
	}

	if err := g.doByBytes(rFile, nBytes); err != nil {
		return err
	}

	return nil
}

// checkFileSize returns fileSize with checking disk free space for output files.
func (g *GoSplit) checkFileSize(rFile *os.File) (int64, Gerror) {
	fi, err := rFile.Stat()
	if err != nil {
		return 0, GoSplitErrorf("failed to stat: %w", err)
	}
	if fi.IsDir() {
		return 0, GoSplitErrorf("is a directory")
	}
	fileSize := fi.Size()

	freeBytesAvailable, gerr := getDiskFreeSpace(g.outDir)
	if err != nil {
		return 0, gerr
	}

	if uint64(fileSize) > freeBytesAvailable {
		return 0, GoSplitErrorf("no free space available")
	}
	return fileSize, nil
}

// generateOutFilePath returns n-th output file name with prefix.
//
// only support 2-character suffix; aa , ab, ..., zz.
func (g *GoSplit) generateOutFilePath(number int) (string, Gerror) {
	var table []byte

	if g.bNumericSuffix {
		table = []byte("0123456789")
	} else {
		table = []byte("abcdefghijklmnopqrstuvwxyz")
	}

	if number >= len(table)*len(table) {
		return "", GoSplitErrorf("output file suffixes exhausted")
	}

	n0 := number % len(table)
	number = number / len(table)
	n1 := number % len(table)
	suffix := string([]byte{table[n1], table[n0]})

	outFileName := g.prefix + suffix
	outFilePath := path.Join(g.outDir, outFileName)
	return outFilePath, nil
}

// doByLines splits the content from io.Reader by nLines.
func (g *GoSplit) doByLines(r io.Reader, nLines int) Gerror {
	scanner := bufio.NewScanner(r)

OuterLoop:
	for i := 0; ; i++ {
		outFilePath, gerr := g.generateOutFilePath(i)
		if gerr != nil {
			return gerr
		}
		wFile, err := os.Create(outFilePath)
		if err != nil {
			return GoSplitErrorf("failed to create: %w", err)
		}
		fmt.Fprintf(g.wVerbose, "creating file %#v\n", outFilePath)

		for j := 0; j < nLines; j++ {
			if !scanner.Scan() {
				if j == 0 {
					defer os.Remove(outFilePath)
				}
				break OuterLoop
			}
			fmt.Fprintln(wFile, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return GoSplitErrorf("failed to read: %w", err)
	}

	return nil
}

// doByNumber splits the content from io.Reader into nNumber files.
func (g *GoSplit) doByNumber(r io.Reader, fileSize int64, nNumber int) Gerror {
	chunkSize := fileSize / int64(nNumber)

	for i := 0; i < nNumber; i++ {
		outFilePath, gerr := g.generateOutFilePath(i)
		if gerr != nil {
			return gerr
		}
		wFile, err := os.Create(outFilePath)
		if err != nil {
			return GoSplitErrorf("failed to create: %w", err)
		}
		fmt.Fprintf(g.wVerbose, "creating file %#v\n", outFilePath)

		// the last file size should be larger than or equal to chunkSize
		if i < nNumber-1 {
			written, err := io.CopyN(wFile, r, chunkSize)
			if written < chunkSize {
				if written == 0 {
					defer os.Remove(outFilePath)
				}
				break
			}
			if err != nil {
				return GoSplitErrorf("failed to write: %w", err)
			}
		} else {
			if _, err := io.Copy(wFile, r); err != nil {
				return GoSplitErrorf("failed to write: %w", err)
			}
		}
	}

	return nil
}

// doByBytes splits the content from io.Reader by nBytes.
func (g *GoSplit) doByBytes(r io.Reader, nBytes int64) Gerror {
	for i := 0; ; i++ {
		outFilePath, gerr := g.generateOutFilePath(i)
		if gerr != nil {
			return gerr
		}
		wFile, err := os.Create(outFilePath)
		if err != nil {
			return GoSplitErrorf("failed to create: %w", err)
		}
		fmt.Fprintf(g.wVerbose, "creating file %#v\n", outFilePath)

		written, err := io.CopyN(wFile, r, nBytes)
		if written < nBytes {
			if written == 0 {
				defer os.Remove(outFilePath)
			}
			break
		}
		if err != nil {
			return GoSplitErrorf("failed to write: %w", err)
		}
	}

	return nil
}
