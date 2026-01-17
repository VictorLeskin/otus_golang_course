package main

import (
	"os"
	"testing"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

type Test_OpSystem_ER struct {
	ret_Environ []string
}

func (os Test_OpSystem_ER) ReadDir(name string) (t []os.DirEntry, e error) {
	return t, nil
}

func (os Test_OpSystem_ER) ReadFile(name string) ([]byte, error) {
	return []byte{}, nil
}

func (os Test_OpSystem_ER) Environ() []string {
	return os.ret_Environ
}

func Test_EnviromentReader_Read(t *testing.T) {
	var os Test_OpSystem_ER
	os.ret_Environ = []string{
		"ABD=99",
		"Q=\"A string\"",
	}
	t0 := EnviromentReader{}
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

func Test_EnviromentReader_convertVariablesToMap(t *testing.T) {
	// Place your code here
	t0 := EnviromentReader{}
	t0.currentVariables = []string{
		"ABD=99",
		"Q=\"A string\"",
	}
	t0.convertVariablesToMap()
	assert.Equal(t, 2, len(t0.mapCurrentVariables))
	require.Equal(t, "99", t0.mapCurrentVariables["ABD"])
	require.Equal(t, "\"A string\"", t0.mapCurrentVariables["Q"])
}

func Test_EnviromentReader_replaceVariables(t *testing.T) {
	t0 := EnviromentReader{}
	t0.mapCurrentVariables = map[string]string{
		"A": "0",
		"B": "1",
		"C": "2",
		"D": "3",
	}

	t0.replaceVariables(map[string]string{
		"A": "",
		"C": "100",
		"D": "",
		"K": "There is not such variables",
	})

	assert.Equal(t, 3, len(t0.mapCurrentVariables))
	assert.Equal(t, "100", t0.mapCurrentVariables["C"])
	assert.Equal(t, "1", t0.mapCurrentVariables["B"])
	assert.Equal(t, "There is not such variables", t0.mapCurrentVariables["K"])
}

func Test_EnviromentReader_makeNewEnviroment(t *testing.T) {
	t0 := EnviromentReader{}
	t0.mapCurrentVariables = map[string]string{
		"A": "A string",
		"B": "1",
		"C": "9999",
	}

	res := t0.makeNewEnviroment()
	require.Equal(t, res, []string{
		"A=A string",
		"B=1",
		"C=9999",
	})
}
