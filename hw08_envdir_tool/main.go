package main

import (
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
	fmt.Println("Utiltity to run a program with a specified set of enviroment variables.")
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

/*

func makeEnvVars(oldEnvVars map[string]string) []string {
	var env []string
	for key, value := range oldEnvVars {
		env = append(env, key+"="+value)
	}
	return env
}

func executeCommand(ret CommanLineParameter, env []string) error {
	cmd := exec.Command(ret.command, ret.arguments...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env

	return cmd.Run()
}

func Exectute(ret CommanLineParameter) error {

	if vars, err := processDir(ret.dirName); err == nil {
		currentEnv := os.Environ() // enviriment variables  "KEY=VALUE"

		envVars := envToMap(currentEnv) // convert them to map

		k := &envVars.variables

		replaceEnvVars(k, vars)

		newEnv := makeEnvVars(envVars.variables)

		return executeCommand(ret, newEnv)
	} else {
		return err
	}
}

*/

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
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "envdir: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
