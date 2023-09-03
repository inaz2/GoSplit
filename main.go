package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	bHelp    bool
	bVersion bool
	nLines   int
	nNumber  int
	nBytes   int64
)

func init() {
	flag.BoolVar(&bHelp, "help", false, "display this help and exit")
	flag.BoolVar(&bVersion, "version", false, "output version information and exit")
	flag.IntVar(&nLines, "l", 0, "put NUMBER lines/records per output file")
	flag.IntVar(&nNumber, "n", 0, "split into N files based on size of input")
	flag.Int64Var(&nBytes, "b", 0, "put SIZE bytes per output file")
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
		err := gosplit.ByLines(nLines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case nNumber > 0:
		err := gosplit.ByNumber(nNumber)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case nBytes > 0:
		err := gosplit.ByBytes(nBytes)
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
