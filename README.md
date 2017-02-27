
## Installation
requires Go to be installed and in $PATH
```
$ go get -u github.com/stvpalumbo/be_exam_candidate_sl
```

## Usage
requires $GOPATH/bin to be in $PATH
```
$ be_exam_candidate_sl [-in=<input-dir>] [-out=<output-dir>] [-err=<error-dir>]
```

## Assumptions
* files initially in our input directory will be considered new and processed
* files not ending in .csv will be ignored
* if not specified in options, default input/output/error directories will be created relative to $PWD
* input records can fit into memory

## Platforms
has been tested on Linux and Windows
