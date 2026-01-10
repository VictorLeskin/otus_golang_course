package main

import (
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// The function return enviroment variable as map
func envToMap(env []string) Environment {

	envMap := make(Environment)

	for _, envLine := range env {
		// split by first '=' and store
		if idx := strings.Index(envLine, "="); idx != -1 {
			key := envLine[:idx]
			value := envLine[idx+1:]
			envMap[key] = EnvValue{value, false}
		}
	}

	return envMap
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here

	currentEnv := os.Environ() // enviriment variables  "KEY=VALUE"

	envVars := envToMap(currentEnv) // convert them to map

	return nil, nil
}
