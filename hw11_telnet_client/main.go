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
		fmt.Fprintln(os.Stderr, err.Error())
		flag.Usage()
		os.Exit(1)
	}

	client := NewTelnetClient(
		params.host+":"+strconv.Itoa(params.port),
		params.timeout,
		os.Stdin,
		os.Stdout)

	if c, ok := client.(*MyTelnetClient); ok {
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
