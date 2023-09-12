package gosplit

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
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

// ConvertSizeToByte converts strSize to nBytes, e.g. "1K" -> 1024.
func (g *GoSplit) ConvertSizeToByte(strSize string) (int64, error) {
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)(?:(.)(iB|B)?)?$`)
	m := re.FindStringSubmatch(strSize)
	if m == nil {
		return 0, fmt.Errorf("invalid number of bytes: %#v", strSize)
	}

	significand, err := strconv.ParseFloat(string(m[1]), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number of bytes: %#v", strSize)
	}

	var (
		base float64
		multiplier float64
	)

	switch m[3] {
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
	case "E":
		multiplier = math.Pow(base, 6)
	case "G":
		multiplier = math.Pow(base, 3)
	case "K", "k":
		multiplier = math.Pow(base, 1)
	case "M", "m":
		multiplier = math.Pow(base, 2)
	case "P":
		multiplier = math.Pow(base, 5)
	case "Q":
		multiplier = math.Pow(base, 10)
	case "R":
		multiplier = math.Pow(base, 9)
	case "T":
		multiplier = math.Pow(base, 4)
	case "Y":
		multiplier = math.Pow(base, 8)
	case "Z":
		multiplier = math.Pow(base, 7)
	default:
		return 0, fmt.Errorf("invalid number of bytes: %#v", strSize)
	}

	n := int64(significand * multiplier)
	if n <= 0 {
		return 0, fmt.Errorf("invalid number of bytes: %#v: Value too large for defined data type", strSize)
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

	fileInfo, err := rFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat: %w", err)
	}
	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return nil
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
	}

	if err := g.byBytesInternal(rFile, nBytes); err != nil {
		return err
	}

	return nil
}

// generateOutFilePath returns n-th output file name with prefix.
//
// only support 2-character suffix; aa , ab, ..., zz.
func (g *GoSplit) generateOutFilePath(number int) (string, error) {
	table := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	if number >= len(table)*len(table) {
		return "", fmt.Errorf("output file suffixes exhausted")
	}

	n0 := number % len(table)
	number = number / len(table)
	n1 := number % len(table)
	suffix := table[n1] + table[n0]

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
