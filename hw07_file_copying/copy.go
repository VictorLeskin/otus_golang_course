package main

import (
	"errors"
	"flag"
	"fmt"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

type CommanLineParameter struct {
	input, output string
	offset, limit int64
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Place your code here.
	return nil
}

func Usage() {
	fmt.Println("Copy a file or a part of it.")
	flag.PrintDefaults()
}

func SetupCommadLineParameters() {
	flag.Usage = Usage
}

func ParseCommadLine() (ret CommanLineParameter, err error) {

	// Copy a part of a file according to the operands.
	flag.StringVar(&ret.input, "from", "", "file to read from")
	flag.StringVar(&ret.output, "to", "", "file to copy")
	flag.Int64Var(&ret.offset, "offset", 0, "skip offset bytes at start of output")
	flag.Int64Var(&ret.limit, "limit", 0, "copy only 'limit' bytes")

	flag.Parse() // проанализировать аргументы

	if ret.input == "" {
		return ret, errors.New("there is not name of the file to read from")
	}
	if ret.output == "" {
		return ret, errors.New("there is not name of the file to copy")
	}
	return ret, nil
}
