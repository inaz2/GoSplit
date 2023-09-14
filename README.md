# GoSplit

Implemented -l, -n, -b options based on GNU coreutils' split behavior.


## CAUTION AND/OR DISCLAIMER

To simplify implementation, the following behaviors are different to the original one.

Regarding CHUNKS specified with the -n option, only the number of output files N is accepted.
Other modes such as K/N, l/N and r/N are not supported.

The suffix of the output file name is limited to two characters aa-zz, and the process exits with an error after the 676th output (with -d option, 00-99 and the 100th output).
The suffix length extension is not implemented for safety.

If the disk free space is less than the input file size, the process exits with an error.

## LICENSE

**Choose the one of these.**

* BSD 3-Clause, as the same as https://github.com/golang/go.
* GNU AFFERO GENERAL PUBLIC LICENSE Version 3

<ins>_Hereby my GoSplit has been started to develop._</ins>


## Supported options

* -l, -n, -b
* -d, -e, --verbose
* --help, --version


## Supported irregular input

* Standard input (including argument string "-")
* Input file of size 0
* The number of lines and output files less than 1
* The size less than 1
* The size with invalid format such as "-1", "1.5" and "2X"


## Performance notice

* Read lines efficiently using bufio.Scanner()
* Read and write files efficiently using io.CopyN()/io.Copy()
* Obtain the file size in advance using (*os.File).Stat()


## Usage

The code quality is poor \
Not intended to use in production \
Not to build and exec(3) through fork(2) \
You could try to run as below

```
$ go run . -help
Usage: /tmp/go-build3934742828/b001/exe/GoSplit [OPTION]... [FILE [PREFIX]]
Output pieces of FILE to PREFIXaa, PREFIXab, ...;
default size is 1000 lines, and default PREFIX is 'x'.

With no FILE, or when FILE is -, read standard input.

  -b string
    	put SIZE bytes per output file
  -d	use numeric suffixes starting at 0, not alphabetic
  -e	do not generate empty output files with '-n'
  -help
    	display this help and exit
  -l int
    	put NUMBER lines/records per output file
  -n int
    	split into N files based on size of input
  -verbose
    	print a diagnostic just before each output file is opened
  -version
    	output version information and exit

The SIZE argument is an integer and optional unit (example: 10K is 10*1024).
Units are K,M,G,T,P,E,Z,Y (powers of 1024) or KB,MB,... (powers of 1000).
Binary prefixes can be used, too: KiB=K, MiB=M, and so on.
```


## Testing

```
$ go test -v ./...
```

* Test methods is created for each purpose for each option.
* Test that the pairs of each output file name and the number of lines/bytes matches.
* Create output files in temporary directories, and remove at the end of the test.


## Debugging

To output debug information, set any value to `DEBUG` environment variable.
