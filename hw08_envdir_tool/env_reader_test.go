package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

type TestOpSystemER struct {
	retEnviron []string
}

func (os TestOpSystemER) ReadDir(_ string) (t []os.DirEntry, e error) {
	return t, nil
}

func (os TestOpSystemER) ReadFile(_ string) ([]byte, error) {
	return []byte{}, nil
}

func (os TestOpSystemER) Environ() []string {
	return os.retEnviron
}

func (os TestOpSystemER) Run(_ *exec.Cmd) error {
	return fmt.Errorf("Not implemented")
}

func Test_EnvironmentReader_Read(t *testing.T) {
	var os TestOpSystemER
	os.retEnviron = []string{
		"ABD=99",
		"Q=\"A string\"",
	}
	t0 := EnvironmentReader{}
	t0.SetOs(os)

	t0.Read()
	require.Equal(t, []string{
		"ABD=99",
		"Q=\"A string\"",
	}, t0.currentVariables)

	assert.Equal(t, 2, len(t0.mapCurrentVariables))
	require.Equal(t, "99", t0.mapCurrentVariables["ABD"])
	require.Equal(t, "\"A string\"", t0.mapCurrentVariables["Q"])
}

func Test_EnvironmentReader_convertVariablesToMap(t *testing.T) {
	// Place your code here
	t0 := EnvironmentReader{}
	t0.currentVariables = []string{
		"ABD=99",
		"Q=\"A string\"",
	}
	t0.convertVariablesToMap()
	assert.Equal(t, 2, len(t0.mapCurrentVariables))
	require.Equal(t, "99", t0.mapCurrentVariables["ABD"])
	require.Equal(t, "\"A string\"", t0.mapCurrentVariables["Q"])
}

func Test_EnvironmentReader_replaceVariables(t *testing.T) {
	t0 := EnvironmentReader{}
	t0.mapCurrentVariables = map[string]string{
		"A": "0",
		"B": "1",
		"C": "2",
		"D": "3",
	}

	res := t0.replaceVariables(map[string]string{
		"A": "",
		"C": "100",
		"D": "",
		"K": "There is not such variables",
	})

	assert.Equal(t, 3, len(t0.mapCurrentVariables))
	assert.Equal(t, "100", t0.mapCurrentVariables["C"])
	assert.Equal(t, "1", t0.mapCurrentVariables["B"])
	assert.Equal(t, "There is not such variables", t0.mapCurrentVariables["K"])

	sort.Strings(res)
	require.Equal(t, res, []string{
		"B=1",
		"C=100",
		"K=There is not such variables",
	})
}

func Test_EnvironmentReader_makeNewEnvironment(t *testing.T) {
	t0 := EnvironmentReader{}
	t0.mapCurrentVariables = map[string]string{
		"A": "A string",
		"B": "1",
		"C": "9999",
	}

	res := t0.makeNewEnvironment()
	sort.Strings(res)
	require.Equal(t, res, []string{
		"A=A string",
		"B=1",
		"C=9999",
	})
}
