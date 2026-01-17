package main

import (
	"strings"
)

type EnvironmentReader struct {
	os iOpSystem // to use a real OS or to emulate it in tests.

	currentVariables    []string
	mapCurrentVariables map[string]string
}

func (er *EnvironmentReader) SetOs(os iOpSystem) {
	er.os = os
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func (er *EnvironmentReader) Read() {
	er.currentVariables = er.os.Environ() // enviriment variables  "KEY=VALUE"
	er.convertVariablesToMap()            // convert to map  KEY:VALUE"
}

// The function return environment variable as map.
func (er *EnvironmentReader) convertVariablesToMap() {
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

func (er *EnvironmentReader) replaceVariables(replacements map[string]string) (env []string) {
	for key, value := range replacements {
		if value == "" {
			delete(er.mapCurrentVariables, key)
		} else {
			er.mapCurrentVariables[key] = value
		}
	}

	return er.makeNewEnvironment()
}

func (er EnvironmentReader) makeNewEnvironment() (env []string) {
	for key, value := range er.mapCurrentVariables {
		env = append(env, key+"="+value)
	}
	return env
}
