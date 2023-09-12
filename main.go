package main

import (
	"inaz2/GoSplit/gosplit"
	"flag"
	"fmt"
	"os"
)

var (
	bHelp    bool
	bVersion bool
	nLines   int
	nNumber  int
	strSize  string
)

func init() {
	flag.BoolVar(&bHelp, "help", false, "display this help and exit")
	flag.BoolVar(&bVersion, "version", false, "output version information and exit")
	flag.IntVar(&nLines, "l", 0, "put NUMBER lines/records per output file")
	flag.IntVar(&nNumber, "n", 0, "split into N files based on size of input")
	flag.StringVar(&strSize, "b", "", "put SIZE bytes per output file")
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
		additionalNote := `
The SIZE argument is an integer and optional unit (example: 10K is 10*1024).
Units are K,M,G,T,P,E,Z,Y (powers of 1024) or KB,MB,... (powers of 1000).
Binary prefixes can be used, too: KiB=K, MiB=M, and so on.
`
		fmt.Printf(usageFormat, os.Args[0])
		flag.PrintDefaults()
		fmt.Print(additionalNote)
		os.Exit(0)
	case bVersion:
		fmt.Println("inaz2/GoSplit 1.0.0")
		os.Exit(0)
	case nLines != 0:
		g := gosplit.New(filePath, prefix)
		err := g.ByLines(nLines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case nNumber != 0:
		g := gosplit.New(filePath, prefix)
		err := g.ByNumber(nNumber)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case strSize != "":
		g := gosplit.New(filePath, prefix)
		nBytes, err := g.ParseSize(strSize)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = g.ByBytes(nBytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		nLines = 1000
		g := gosplit.New(filePath, prefix)
		err := g.ByLines(nLines)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
