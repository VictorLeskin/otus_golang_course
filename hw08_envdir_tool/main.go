package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		arguments: os.Args[3:], // вот это самая простая форма
	}, nil
}

func processFileContent(content []byte) string {
	if len(content) == 0 {
		return "" // Empty file
	}

	// first line
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if !scanner.Scan() {
		return ""
	}

	firstLine := scanner.Text()

	firstLine = strings.TrimRight(firstLine, " \t") // strip a end of line
	firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")

	return firstLine
}

func processDir(dirPath string) (map[string]string, error) {
	envVars := make(map[string]string)

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			continue // skip files with '=' in name
		}

		content, err := os.ReadFile(filepath.Join(dirPath, file.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "File reading error: %s: %v\n", file.Name(), err)
			continue
		}

		envVars[file.Name()] = processFileContent(content)
	}

	return envVars, nil
}

func replaceEnvVars(oldEnvVars *map[string]string, replacements map[string]string) {
	for key, value := range replacements {
		if value == "" {
			delete(*oldEnvVars, key)
		} else {
			(*oldEnvVars)[key] = value
		}
	}
}

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

		replaceEnvVars(&envVars, vars)

		newEnv := makeEnvVars(envVars)

		return executeCommand(ret, newEnv)
	} else {
		return err
	}
}

func main() {
	SetupCommadLineParameters()

	params, err := ParseCommadLine()
	if err != nil {
		fmt.Println(err.Error())
		flag.Usage()
		os.Exit(1)
	}

	err = Exectute(params)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "envdir: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
