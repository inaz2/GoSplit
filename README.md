# GoSplit

Implemented -l, -n, -b options based on GNU coreutils' split behavior.


## CAUTION AND/OR DISCLAIMER

To simplify implementation, the following behaviors are different to the original one.

Regarding SIZE specified with the -b option, unit specifications such as 1K and 1KiB are omitted, and only the number of bytes is specified.
Regarding CHUNKS specified with the -n option, detailed specifications such as K/N and l/N are omitted, and only the number of output files N is specified.

The suffix of the output file name has been fixed at two characters aa-zz, omitting the extension when the number increases.
Therefore, the process ends with an error when the 676th output is executed.


## LICENSE

**Choose the one of these.**

* BSD 3-Clause, as the same as https://github.com/golang/go.
* GNU AFFERO GENERAL PUBLIC LICENSE Version 3

<ins>_Hereby my gosplit has been started to develop._</ins>


## Supported options

* -l, -n, -b
* --help, --version


## Supported irregular input

* Standard input (including argument string "-")
* input file of size 0
* prefix string of length 0
* Number of lines less than 0, number of output files, number of bytes


## Performance notice

* Read rows efficiently using bufio.Scanner()
* Efficiently read and write files using io.CopyN()/io.Copy()
* Obtained the file size in advance using (*os.File).Stat()


## Usage

The code quality is poor \
Not intended to use in production \
Not to build and exec(3) through fork(2) \
You could try to run as below

```
$ go run main.go --help
Usage: /tmp/go-build3057508251/b001/exe/main [OPTION]... [FILE [PREFIX]]
Output pieces of FILE to PREFIXaa, PREFIXab, ...;
default size is 1000 lines, and default PREFIX is 'x'.

With no FILE, or when FILE is -, read standard input.
  -b int
    	put SIZE bytes per output file
  -help
    	display this help and exit
  -l int
    	put NUMBER lines/records per output file
  -n int
    	split into N files based on size of input
  -version
    	output version information and exit
```


## Testing

```
$ go test -v
```

* Test methods were created for each purpose for each option.
* Define the expected output file name and its number of lines and bytes, and test that they all match.
* Specify a prefix that corresponds to the test method name, and delete the output file at the end of the test.
