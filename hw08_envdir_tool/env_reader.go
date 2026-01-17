package main

import (
	"strings"
)

type EnviromentReader struct {
	os iOpSystem // to use a real OS or to emulate it in tests.

	currentVariables    []string
	mapCurrentVariables map[string]string
}

func (er *EnviromentReader) SetOs(os iOpSystem) {
	er.os = os
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func (er *EnviromentReader) Read() {
	er.currentVariables = er.os.Environ() // enviriment variables  "KEY=VALUE"
	er.convertVariablesToMap()            // convert to map  KEY:VALUE"
}

// The function return enviroment variable as map
func (er *EnviromentReader) convertVariablesToMap() {
	er.mapCurrentVariables = make(map[string]string)
	for _, envLine := range er.currentVariables {
		// split by first '=' and store
		if idx := strings.Index(envLine, "="); idx != -1 {
			key := envLine[:idx]
			value := envLine[idx+1:]
			er.mapCurrentVariables[key] = value
		}
	}
}

func (er *EnviromentReader) replaceVariables(replacements map[string]string) {
	for key, value := range replacements {
		if value == "" {
			delete(er.mapCurrentVariables, key)
		} else {
			er.mapCurrentVariables[key] = value
		}
	}
}

func (er EnviromentReader) makeNewEnviroment() (env []string) {
	for key, value := range er.mapCurrentVariables {
		env = append(env, key+"="+value)
	}
	return env
}

/*
type Environment struct {
	variables map[string]EnvValue
}

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here

	currentEnv := os.Environ() // enviriment variables  "KEY=VALUE"

	envVars := envToMap(currentEnv) // convert them to map

	return envVars, nil
}
*/
