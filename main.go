package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	bHelp    bool
	bVersion bool
	nLines   int
	nNumber  int
	nBytes   int
)

// generateOutFileName returns n-th output file name with prefix
//
// only support 2-character suffix; aa , ab, ..., zz
func generateOutFileName(prefix string, number int) (string, error) {
	table := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	if number >= len(table)*len(table) {
		return "", fmt.Errorf("output file suffixes exhausted")
	}

	n0 := number % len(table)
	number = number / len(table)
	n1 := number % len(table)
	suffix := table[n1] + table[n0]

	outFileName := prefix + suffix
	return outFileName, nil
}

// openFileOrStdin opens filePath or returns os.Stdin
func openFileOrStdin(filePath string) (*os.File, error) {
	if filePath == "-" {
		return os.Stdin, nil
	} else {
		rFile, err := os.Open(filePath)
		return rFile, err
	}
}

// GoSplit holds filePath and prefix
type GoSplit struct {
	filePath string
	prefix string
}

// ByLines splits the content of rFile by nLines
func (g GoSplit) ByLines(nLines int) error {
	if g.prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nLines <= 0 {
		return fmt.Errorf("nLines must be larger than zero")
	}

	rFile, err := openFileOrStdin(g.filePath)
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}
	defer rFile.Close()

	scanner := bufio.NewScanner(rFile)

OuterLoop:
	for i := 0; ; i++ {
		outFileName, err := generateOutFileName(g.prefix, i)
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
func (g GoSplit) ByNumber(nNumber int) error {
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
		outFileName, err := generateOutFileName(g.prefix, i)
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
func (g GoSplit) ByBytes(nBytes int64) error {
	if g.prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nBytes <= 0 {
		return fmt.Errorf("nBytes must be larger than zero")
	}

	rFile, err := openFileOrStdin(g.filePath)
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}
	defer rFile.Close()

	for i := 0; ; i++ {
		outFileName, err := generateOutFileName(g.prefix, i)
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

func init() {
	flag.BoolVar(&bHelp, "help", false, "display this help and exit")
	flag.BoolVar(&bVersion, "version", false, "output version information and exit")
	flag.IntVar(&nLines, "l", 0, "put NUMBER lines/records per output file")
	flag.IntVar(&nNumber, "n", 0, "split into N files based on size of input")
	flag.IntVar(&nBytes, "b", 0, "put SIZE bytes per output file")
}

func main() {
	var (
		filePath string
		prefix   string
	)

	flag.Parse()

	switch flag.NArg() {
	case 0:
		filePath = "-"
		prefix = "x"
	case 1:
		filePath = flag.Args()[0]
		prefix = "x"
	default:
		filePath = flag.Args()[0]
		prefix = flag.Args()[1]
	}

	gosplit := GoSplit{filePath, prefix}

	switch {
	case bHelp:
		usageFormat := `Usage: %s [OPTION]... [FILE [PREFIX]]
Output pieces of FILE to PREFIXaa, PREFIXab, ...;
default size is 1000 lines, and default PREFIX is 'x'.

With no FILE, or when FILE is -, read standard input.
`
		fmt.Printf(usageFormat, os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	case bVersion:
		fmt.Println("inaz2/GoSplit 1.0.0")
		os.Exit(0)
	case nLines > 0:
		gosplit := GoSplit{filePath, prefix}
		err := gosplit.ByLines(nLines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case nNumber > 0:
		gosplit := GoSplit{filePath, prefix}
		err := gosplit.ByNumber(nNumber)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case nBytes > 0:
		gosplit := GoSplit{filePath, prefix}
		err := gosplit.ByBytes(int64(nBytes))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		nLines = 1000
		err := gosplit.ByLines(nLines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
