package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {
	SetupCommadLineParameters()

	params, err := ParseCommandLine()
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	address := params.host + ":" + strconv.Itoa(params.port)
	var in io.ReadCloser
	var out io.Writer

	tc := NewTelnetClient(address, params.timeout, in, out)

	_ = tc
	os.Exit(0)
}
