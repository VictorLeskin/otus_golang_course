package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

// Mock DirEntry implementation for testing.
type mockDirEntry struct {
	name  string
	isDir bool
}

type TestOpSystem struct {
	readDirError string
	entries      []mockDirEntry
	fileContent  map[string][]byte

	retRun error
}

var osCommandHasBeenExecuted []string
var osRunsEnvironment []string

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func (os TestOpSystem) ReadDir(_ string) (t []os.DirEntry, e error) {
	if os.readDirError != "" {
		return nil, fmt.Errorf(os.readDirError)
	}

	for _, en := range os.entries {
		t = append(t, en)
	}
	return t, nil
}

func (os TestOpSystem) ReadFile(name string) ([]byte, error) {
	// Check if key exists.
	if value, exists := os.fileContent[name]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("File reading error")
}

func (os TestOpSystem) Environ() []string {
	return []string{}
}

func (os TestOpSystem) Run(cmd *exec.Cmd) error {
	osCommandHasBeenExecuted = cmd.Args
	osRunsEnvironment = cmd.Env
	return os.retRun
}

func Test_Executor_ConvertDirectoryToStrings(t *testing.T) {
	// REading file, skip directory and file with '=' in name.
	{
		var os TestOpSystem

		os.entries = []mockDirEntry{
			{name: "ABC", isDir: false},
			{name: "folder1", isDir: true},
			{name: "DEF = 99.go", isDir: false},
		}

		parameters := CommanLineParameter{dirName: "testDir"}

		os.fileContent = make(map[string][]byte)
		os.fileContent["testDir\\ABC"] = []byte("bar\nPLEASE IGNORE SECOND LINE")

		t0 := Executor{parameters: parameters, os: os}

		err := t0.ConvertDirectoryToStrings()

		assert.Equal(t, 1, len(t0.dirContent))
		assert.Equal(t, nil, err)

		value, exist := t0.dirContent["ABC"]
		assert.Equal(t, true, exist)
		assert.Equal(t, []byte("bar\nPLEASE IGNORE SECOND LINE"), value)
	}

	// ReadDir error.
	{
		var os TestOpSystem
		os.readDirError = "Dir reading error"

		parameters := CommanLineParameter{dirName: "testDir"}

		os.fileContent = make(map[string][]byte)
		os.fileContent["testDir\\ABC"] = []byte("bar\nPLEASE IGNORE SECOND LINE")

		t0 := Executor{parameters: parameters, os: os}

		err := t0.ConvertDirectoryToStrings()

		assert.Equal(t, "Dir reading error", err.Error())
		assert.Equal(t, 0, len(t0.dirContent))
	}
}

func Test_Executor_processFileContent(t *testing.T) {
	t0 := Executor{}

	require.Equal(t, "bar", t0.processFileContent([]byte("bar\nPLEASE IGNORE SECOND LINE")))
	require.Equal(t, "", t0.processFileContent([]byte{}))
	require.Equal(t, "A\nB", t0.processFileContent([]byte{'A', 0x00, 'B'}))
	require.Equal(t, "", t0.processFileContent([]byte(" \t")))
	require.Equal(t, "bar", t0.processFileContent([]byte("bar   \nPLEASE IGNORE SECOND LINE")))
}

func Test_Executor_makeNewEnvironmentVariables(t *testing.T) {
	t0 := Executor{}
	t0.dirContent = map[string][]byte{
		"ABC": []byte("bar\nPLEASE IGNORE SECOND LINE"),
		"DEF": {},
	}

	t0.makeNewEnvironmentVariables()

	assert.Equal(t, 2, len(t0.newEnvironmentVariables))
	require.Equal(t, "bar", t0.newEnvironmentVariables["ABC"])
	require.Equal(t, "", t0.newEnvironmentVariables["DEF"])
}

func Test_Executor_ExecuteInEnvironment(t *testing.T) {
	var os TestOpSystem
	t0 := Executor{os: os}
	t0.command = "cmd"
	t0.arguments = []string{"1", "2"}

	env := []string{"A=99", "B=88"}

	ret := t0.ExecuteInEnvironment(env)

	assert.Equal(t, []string{"cmd", "1", "2"}, osCommandHasBeenExecuted)
	assert.Equal(t, []string{"A=99", "B=88"}, osRunsEnvironment)
	assert.Equal(t, true, ret == nil)
}
