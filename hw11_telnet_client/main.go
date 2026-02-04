package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	SetupCommandLineParameters()

	params, err := ParseCommandLine()
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	client := NewTelnetClient(
		params.host+":"+strconv.Itoa(params.port),
		params.timeout,
		nil,
		nil)

	if c, ok := client.(*MyTelnetClient); ok {
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
