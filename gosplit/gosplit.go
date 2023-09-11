package gosplit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
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

// ByLines splits the content of filePath by nLines.
func (g *GoSplit) ByLines(nLines int) error {
	if g.prefix == "" {
		return fmt.Errorf("prefix must not be empty string")
	}
	if nLines <= 0 {
		return fmt.Errorf("nLines must be larger than zero")
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
	if g.prefix == "" {
		return fmt.Errorf("prefix must not be empty string")
	}
	if nNumber <= 0 {
		return fmt.Errorf("nNumber must be larger than zero")
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
	if g.prefix == "" {
		return fmt.Errorf("prefix must not be empty string")
	}
	if nBytes <= 0 {
		return fmt.Errorf("nBytes must be larger than zero")
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
