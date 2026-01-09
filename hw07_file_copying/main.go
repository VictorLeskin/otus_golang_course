package main

import (
	"flag"
	"fmt"
	"os"
)


func main() {

	SetupCommadLineParameters()

	params, err := ParseCommadLine()
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	if err = Copy(params.input, params.output, params.offset, params.offset); err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}

	os.Exit(0)
}
