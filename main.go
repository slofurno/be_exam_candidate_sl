package main

import (
	"flag"
)

func main() {

	inputDir := flag.String("in", "input", "directory to watch for input")
	outputDir := flag.String("out", "output", "output directory")
	errorDir := flag.String("err", "errors", "error record directory")
	flag.Parse()

	store := newFileStore(*inputDir, *outputDir, *errorDir)
	logger := logLogger{}
	p := &recordProcessor{store, logger}
	p.Run(*inputDir)
}
