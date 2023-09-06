package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// GoSplit holds filePath and prefix
type GoSplit struct {
	filePath string
	prefix   string
}

func NewGoSplit(filePath string, prefix string) *GoSplit {
	return &GoSplit{filePath, prefix}
}

// generateOutFileName returns n-th output file name with prefix
//
// only support 2-character suffix; aa , ab, ..., zz
func (g *GoSplit) generateOutFileName(number int) (string, error) {
	table := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	if number >= len(table)*len(table) {
		return "", fmt.Errorf("output file suffixes exhausted")
	}

	n0 := number % len(table)
	number = number / len(table)
	n1 := number % len(table)
	suffix := table[n1] + table[n0]

	outFileName := g.prefix + suffix
	return outFileName, nil
}

// openFileOrStdin opens filePath or returns os.Stdin
func (g *GoSplit) openFileOrStdin() (*os.File, error) {
	if g.filePath == "-" {
		return os.Stdin, nil
	} else {
		rFile, err := os.Open(g.filePath)
		return rFile, err
	}
}

// ByLines splits the content of rFile by nLines
func (g *GoSplit) ByLines(nLines int) error {
	if g.prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nLines <= 0 {
		return fmt.Errorf("nLines must be larger than zero")
	}

	rFile, err := g.openFileOrStdin()
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}
	defer rFile.Close()

	scanner := bufio.NewScanner(rFile)

OuterLoop:
	for i := 0; ; i++ {
		outFileName, err := g.generateOutFileName(i)
		if err != nil {
			return fmt.Errorf("Failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFileName)
		if err != nil {
			return fmt.Errorf("Failed to create: %w", err)
		}
		for j := 0; j < nLines; j++ {
			if !scanner.Scan() {
				if j == 0 {
					defer os.Remove(outFileName)
				}
				break OuterLoop
			}
			fmt.Fprintln(wFile, scanner.Text())
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("Failed to read: %w", err)
	}

	return nil
}

// ByNumber splits the content of rFile into nNumber files
func (g *GoSplit) ByNumber(nNumber int) error {
	// print error message when filePath is stdin
	if g.filePath == "-" {
		return fmt.Errorf("cannot determine file size")
	}

	if g.prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nNumber <= 0 {
		return fmt.Errorf("nNumber must be larger than zero")
	}

	rFile, err := os.Open(g.filePath)
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}
	defer rFile.Close()

	fileInfo, err := rFile.Stat()
	if err != nil {
		return fmt.Errorf("Failed to stat: %w", err)
	}
	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return nil
	}
	chunkSize := fileSize / int64(nNumber)

	for i := 0; i < nNumber; i++ {
		outFileName, err := g.generateOutFileName(i)
		if err != nil {
			return fmt.Errorf("Failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFileName)
		if err != nil {
			return fmt.Errorf("Failed to create: %w", err)
		}
		// the last file size should be larger than or equal to chunkSize
		if i < nNumber-1 {
			written, err := io.CopyN(wFile, rFile, chunkSize)
			if written < chunkSize {
				if written == 0 {
					defer os.Remove(outFileName)
				}
				break
			}
			if err != nil {
				return fmt.Errorf("Failed to write: %w", err)
			}
		} else {
			if _, err := io.Copy(wFile, rFile); err != nil {
				return fmt.Errorf("Failed to write: %w", err)
			}
		}
	}

	return nil
}

// ByBytes splits the content of rFile by nBytes
func (g *GoSplit) ByBytes(nBytes int64) error {
	if g.prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nBytes <= 0 {
		return fmt.Errorf("nBytes must be larger than zero")
	}

	rFile, err := g.openFileOrStdin()
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}
	defer rFile.Close()

	for i := 0; ; i++ {
		outFileName, err := g.generateOutFileName(i)
		if err != nil {
			return fmt.Errorf("Failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFileName)
		if err != nil {
			return fmt.Errorf("Failed to create: %w", err)
		}
		written, err := io.CopyN(wFile, rFile, nBytes)
		if written < nBytes {
			if written == 0 {
				defer os.Remove(outFileName)
			}
			break
		}
		if err != nil {
			return fmt.Errorf("Failed to write: %w", err)
		}
	}

	return nil
}
