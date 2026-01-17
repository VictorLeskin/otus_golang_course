package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Executor struct {
	parameters CommanLineParameter

	dirContent             map[string][]byte
	newEnvironmentVariables map[string]string

	command   string
	arguments []string

	os iOpSystem // to use a real OS or to emulate it in tests.
}

func NewExecutor(parameters CommanLineParameter) *Executor {
	return &Executor{parameters: parameters}
}

func (ex Executor) ExecuteInEnvironment(env []string) error {
	cmd := exec.Command(ex.command, ex.arguments...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = env

	return ex.os.Run(cmd)
}

func (ex Executor) Execute() error {
	if err := ex.ConvertDirectoryToStrings(); err != nil {
		return err
	}

	ex.makeNewEnvironmentVariables()

	er := EnvironmentReader{}

	er.SetOs(ex.os)
	er.Read()

	newEnvVars := er.replaceVariables(ex.newEnvironmentVariables)

	if err := ex.ExecuteInEnvironment(newEnvVars); err != nil {
		return err
	}

	return nil
}

func (ex *Executor) SetOs(os iOpSystem) {
	ex.os = os
}

func (ex *Executor) ConvertDirectoryToStrings() error {
	files, err := ex.os.ReadDir(ex.parameters.dirName)
	if err != nil {
		return err
	}

	ex.dirContent = make(map[string][]byte)
	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			fmt.Fprintf(os.Stderr, "There is '=' in file name: '%s'\n", file.Name())
			continue // skip files with '=' in name
		}

		content, err := ex.os.ReadFile(filepath.Join(ex.parameters.dirName, file.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "File reading error: %s: %v\n", file.Name(), err)
			continue
		}

		ex.dirContent[file.Name()] = content
	}

	return nil
}

func (ex *Executor) processFileContent(content []byte) string {
	if len(content) == 0 {
		return "" // Empty file.
	}

	// first line
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if !scanner.Scan() {
		return "" // Empty first line.
	}

	firstLine := scanner.Text()

	firstLine = strings.TrimRight(firstLine, " \t") // strip a end of line
	firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")

	return firstLine
}

func (ex *Executor) makeNewEnvironmentVariables() {
	ex.newEnvironmentVariables = make(map[string]string)

	for envName, fileContent := range ex.dirContent {
		ex.newEnvironmentVariables[envName] = ex.processFileContent(fileContent)
	}
}
