package gosplit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
)

// GoSplit provides the methods for splitting the file.
type GoSplit struct {
	filePath string
	prefix   string
	outDir   string
}

// New returns a new GoSplit struct.
func New(filePath string, prefix string) *GoSplit {
	return &GoSplit{
		filePath: filePath,
		prefix:   prefix,
		outDir:   "./",
	}
}

// SetOutDir changes the directory of output files.
//
// This method is mainly for testing.
func (g *GoSplit) SetOutDir(outDir string) {
	g.outDir = outDir
}

// ParseSize converts strSize to nBytes, e.g. "10K" -> 10 * 1024.
func (g *GoSplit) ParseSize(strSize string) (int64, error) {
	re := regexp.MustCompile(`^(\d+)(b|(\w)(iB|B)?)?$`)
	m := re.FindStringSubmatch(strSize)
	if m == nil {
		return 0, fmt.Errorf("invalid number of bytes: %#v", strSize)
	}

	x, err := strconv.ParseInt(string(m[1]), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number of bytes: %#v", strSize)
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
			return 0, fmt.Errorf("invalid number of bytes: %#v", strSize)
		}
		multiplier, err = safePowInt64(base, exponent)
		if err != nil {
			// integer overflow occured
			return 0, fmt.Errorf("invalid number of bytes: %#v: Value too large for defined data type", strSize)
		}
	}

	n, err := safeMulInt64(x, multiplier)
	if err != nil {
		// integer overflow occured
		return 0, fmt.Errorf("invalid number of bytes: %#v: Value too large for defined data type", strSize)
	}

	if n <= 0 {
		return 0, fmt.Errorf("invalid number of bytes: %#v: Numerical result out of range", strSize)
	}

	return n, nil
}

// ByLines splits the content of filePath by nLines.
func (g *GoSplit) ByLines(nLines int) error {
	if nLines <= 0 {
		return fmt.Errorf("invalid number of lines: %#v", nLines)
	}

	var rFile *os.File
	if g.filePath == "-" {
		rFile = os.Stdin
	} else {
		f, err := os.Open(g.filePath)
		if err != nil {
			return fmt.Errorf("failed to open: %w", err)
		}
		defer f.Close()

		rFile = f
		if _, err := g.checkFileSize(rFile); err != nil {
			return err
		}
	}

	if err := g.byLinesInternal(rFile, nLines); err != nil {
		return err
	}

	return nil
}

// ByNumber splits the content of filePath into nNumber files.
func (g *GoSplit) ByNumber(nNumber int) error {
	if nNumber <= 0 {
		return fmt.Errorf("invalid number of chunks: %#v", nNumber)
	}

	var rFile *os.File
	if g.filePath == "-" {
		// print error message when filePath is stdin
		return fmt.Errorf("cannot determine file size")
	} else {
		f, err := os.Open(g.filePath)
		if err != nil {
			return fmt.Errorf("failed to open: %w", err)
		}
		defer f.Close()

		rFile = f
	}

	fileSize, err := g.checkFileSize(rFile)
	if err != nil {
		return err
	}

	if err := g.byNumberInternal(rFile, fileSize, nNumber); err != nil {
		return err
	}

	return nil
}

// ByBytes splits the content of filePath by nBytes.
func (g *GoSplit) ByBytes(nBytes int64) error {
	if nBytes <= 0 {
		return fmt.Errorf("invalid number of bytes: %#v", nBytes)
	}

	var rFile *os.File
	if g.filePath == "-" {
		rFile = os.Stdin
	} else {
		f, err := os.Open(g.filePath)
		if err != nil {
			return fmt.Errorf("failed to open: %w", err)
		}
		defer f.Close()

		rFile = f
		if _, err := g.checkFileSize(rFile); err != nil {
			return err
		}
	}

	if err := g.byBytesInternal(rFile, nBytes); err != nil {
		return err
	}

	return nil
}

// safePowInt64 returns b**k with checking integer overflow
func safePowInt64(b int64, k int64) (int64, error) {
	var err error

	if k < 0 {
		return 0, nil
	}

	result := int64(1)
	x := b
	for {
		if k&1 == 1 {
			result, err = safeMulInt64(result, x)
			if err != nil {
				return 0, err
			}
		}
		k >>= 1
		if k <= 0 {
			break
		}
		x, err = safeMulInt64(x, x)
		if err != nil {
			return 0, err
		}
	}
	return result, nil
}

// safeMulInt64 return x*y with checking integer overflow
func safeMulInt64(x int64, y int64) (int64, error) {
	z := x * y
	if y != 0 && z/y != x {
		return 0, fmt.Errorf("integer overflow occured: %#v * %#v -> %#v", x, y, z)
	}
	return z, nil
}

// checkFileSize returns fileSize with checking disk free space for output files.
func (g *GoSplit) checkFileSize(rFile *os.File) (int64, error) {
	fileInfo, err := rFile.Stat()
	if err != nil {
		return 0, fmt.Errorf("failed to stat: %w", err)
	}
	fileSize := fileInfo.Size()

	freeBytesAvailable, err := getDiskFreeSpace(g.outDir)
	if err != nil {
		return 0, fmt.Errorf("failed to getDiskFreeSpace: %w", err)
	}

	if uint64(fileSize) > freeBytesAvailable {
		return 0, fmt.Errorf("no free space available")
	}
	return fileSize, nil
}

// generateOutFilePath returns n-th output file name with prefix.
//
// only support 2-character suffix; aa , ab, ..., zz.
func (g *GoSplit) generateOutFilePath(number int) (string, error) {
	table := []byte("abcdefghijklmnopqrstuvwxyz")
	if number >= len(table)*len(table) {
		return "", fmt.Errorf("output file suffixes exhausted")
	}

	n0 := number % len(table)
	number = number / len(table)
	n1 := number % len(table)
	suffix := string([]byte{table[n1], table[n0]})

	outFileName := g.prefix + suffix
	outFilePath := path.Join(g.outDir, outFileName)
	return outFilePath, nil
}

// byLinesInternal splits the content from io.Reader by nLines.
func (g *GoSplit) byLinesInternal(r io.Reader, nLines int) error {
	scanner := bufio.NewScanner(r)

OuterLoop:
	for i := 0; ; i++ {
		outFilePath, err := g.generateOutFilePath(i)
		if err != nil {
			return fmt.Errorf("failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFilePath)
		if err != nil {
			return fmt.Errorf("failed to create: %w", err)
		}
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
		return fmt.Errorf("failed to read: %w", err)
	}

	return nil
}

// byNumberInternal splits the content from io.Reader into nNumber files.
func (g *GoSplit) byNumberInternal(r io.Reader, fileSize int64, nNumber int) error {
	chunkSize := fileSize / int64(nNumber)

	for i := 0; i < nNumber; i++ {
		outFilePath, err := g.generateOutFilePath(i)
		if err != nil {
			return fmt.Errorf("failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFilePath)
		if err != nil {
			return fmt.Errorf("failed to create: %w", err)
		}
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
				return fmt.Errorf("failed to write: %w", err)
			}
		} else {
			if _, err := io.Copy(wFile, r); err != nil {
				return fmt.Errorf("failed to write: %w", err)
			}
		}
	}

	return nil
}

// byBytesInternal splits the content from io.Reader by nBytes.
func (g *GoSplit) byBytesInternal(r io.Reader, nBytes int64) error {
	for i := 0; ; i++ {
		outFilePath, err := g.generateOutFilePath(i)
		if err != nil {
			return fmt.Errorf("failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFilePath)
		if err != nil {
			return fmt.Errorf("failed to create: %w", err)
		}
		written, err := io.CopyN(wFile, r, nBytes)
		if written < nBytes {
			if written == 0 {
				defer os.Remove(outFilePath)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
	}

	return nil
}
