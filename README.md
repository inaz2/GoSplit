# GoSplit

Implemented -l, -n, -b options based on GNU coreutils' split behavior.


## CAUTION AND/OR DISCLAIMER

To simplify implementation, the following behaviors are different to the original one.

Regarding CHUNKS specified with the -n option, detailed specifications such as K/N and l/N are omitted, and only the number of output files N is specified.

The suffix of the output file name is limited to two characters aa-zz, omitting the length extension when the number increases.
Therefore, the process exits with an error after the 676th output.


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
* Input file of size 0
* Prefix string of length 0
* The number of lines, output files and bytes less than 1


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
$ go run . --help
Usage: /tmp/go-build2320408667/b001/exe/GoSplit [OPTION]... [FILE [PREFIX]]
Output pieces of FILE to PREFIXaa, PREFIXab, ...;
default size is 1000 lines, and default PREFIX is 'x'.

With no FILE, or when FILE is -, read standard input.

The SIZE argument is an integer and optional unit (example: 10K is 10*1024).
Units are K,M,G,T,P,E,Z,Y (powers of 1024) or KB,MB,... (powers of 1000).
Binary prefixes can be used, too: KiB=K, MiB=M, and so on.

  -b string
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
$ go test -v ./...
```

* Test methods is created for each purpose for each option.
* Test that the pairs of each output file name and the number of lines/bytes matches.
* Create output files in temporary directories, and remove at the end of the test.
