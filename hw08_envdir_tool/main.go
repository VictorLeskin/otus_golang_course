package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
)

type CommanLineParameter struct {
	dirName, command string
	arguments        []string
}

func Usage() {
	fmt.Println("Utiltity to run a program with a specified set of environment variables.")
	flag.PrintDefaults()
}

func SetupCommadLineParameters() {
	flag.Usage = Usage
}

func ParseCommadLine() (ret CommanLineParameter, err error) {
	if len(os.Args) < 3 {
		return CommanLineParameter{}, fmt.Errorf("usage: /path/to/env/dir command [args...]")
	}

	return CommanLineParameter{
		dirName:   os.Args[1],
		command:   os.Args[2],
		arguments: os.Args[3:],
	}, nil
}

func main() {
	SetupCommadLineParameters()

	params, err := ParseCommadLine()
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	mExec := NewExecutor(params)

	var myOS OpSystem
	mExec.SetOs(myOS)

	if err := mExec.Execute(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			fmt.Printf("Command exited with code: %d\n", exitErr.ExitCode())
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "envdir: %v\n", err)
		os.Exit(2)
	}

	os.Exit(0)
}
