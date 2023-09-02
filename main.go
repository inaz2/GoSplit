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

// splitByLines splits the content of rFile by nLines
func splitByLines(filePath string, prefix string, nLines int) error {
	if prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nLines <= 0 {
		return fmt.Errorf("nLines must be larger than zero")
	}

	rFile, err := openFileOrStdin(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}
	defer rFile.Close()

	scanner := bufio.NewScanner(rFile)

OuterLoop:
	for i := 0; ; i++ {
		outFileName, err := generateOutFileName(prefix, i)
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

// splitByNumber splits the content of rFile into nNumber files
func splitByNumber(filePath string, prefix string, nNumber int) error {
	// print error message when filePath is stdin
	if filePath == "-" {
		return fmt.Errorf("cannot determine file size")
	}

	if prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nNumber <= 0 {
		return fmt.Errorf("nNumber must be larger than zero")
	}

	rFile, err := os.Open(filePath)
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
	chunkSize := int(float64(fileSize) / float64(nNumber))

	for i := 0; i < nNumber; i++ {
		outFileName, err := generateOutFileName(prefix, i)
		if err != nil {
			return fmt.Errorf("Failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFileName)
		if err != nil {
			return fmt.Errorf("Failed to create: %w", err)
		}
		// the last file size should be larger than or equal to chunkSize
		if i < nNumber-1 {
			written, err := io.CopyN(wFile, rFile, int64(chunkSize))
			if written < int64(chunkSize) {
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

// splitByBytes splits the content of rFile by nBytes
func splitByBytes(filePath string, prefix string, nBytes int) error {
	if prefix == "" {
		return fmt.Errorf("PREFIX must not be empty string")
	}
	if nBytes <= 0 {
		return fmt.Errorf("nBytes must be larger than zero")
	}

	rFile, err := openFileOrStdin(filePath)
	if err != nil {
		return fmt.Errorf("Failed to open: %w", err)
	}
	defer rFile.Close()

	for i := 0; ; i++ {
		outFileName, err := generateOutFileName(prefix, i)
		if err != nil {
			return fmt.Errorf("Failed to generate file name: %w", err)
		}
		wFile, err := os.Create(outFileName)
		if err != nil {
			return fmt.Errorf("Failed to create: %w", err)
		}
		written, err := io.CopyN(wFile, rFile, int64(nBytes))
		if written < int64(nBytes) {
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
		err := splitByLines(filePath, prefix, nLines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case nNumber > 0:
		err := splitByNumber(filePath, prefix, nNumber)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case nBytes > 0:
		err := splitByBytes(filePath, prefix, nBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		nLines = 1000
		err := splitByLines(filePath, prefix, nLines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
