package main

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

// Mock DirEntry implementation for testing
type mockDirEntry struct {
	name  string
	isDir bool
}

type Test_OpSystem struct {
	readDirError string
	entries      []mockDirEntry
	fileContent  map[string][]byte
}

func (m mockDirEntry) Name() string               { return m.name }
func (m mockDirEntry) IsDir() bool                { return m.isDir }
func (m mockDirEntry) Type() fs.FileMode          { return 0 }
func (m mockDirEntry) Info() (fs.FileInfo, error) { return nil, nil }

func (os Test_OpSystem) ReadDir(name string) (t []os.DirEntry, e error) {
	if os.readDirError != "" {
		return nil, fmt.Errorf(os.readDirError)
	}

	for _, en := range os.entries {
		t = append(t, en)
	}
	return t, nil
}

func (os Test_OpSystem) ReadFile(name string) ([]byte, error) {
	// Check if key exists.
	if value, exists := os.fileContent[name]; exists {
		return value, nil
	} else {
		return nil, fmt.Errorf("File reading error")
	}
}

func (os Test_OpSystem) Environ() []string {
	return []string{}
}

func TestRunCmd(t *testing.T) {
	// Place your code here
}

func Test_Executor_ConvertDirectoryToStrings(t *testing.T) {

	// REading file, skip directory and file with '=' in name.
	{
		var os Test_OpSystem

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
		var os Test_OpSystem
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

func Test_Executor_makeNewEnviromentVariables(t *testing.T) {
	t0 := Executor{}
	t0.dirContent = map[string][]byte{
		"ABC": []byte("bar\nPLEASE IGNORE SECOND LINE"),
		"DEF": {},
	}

	t0.makeNewEnviromentVariables()

	assert.Equal(t, 2, len(t0.newEnviromentVariables))
	require.Equal(t, "bar", t0.newEnviromentVariables["ABC"])
	require.Equal(t, "", t0.newEnviromentVariables["DEF"])
}
