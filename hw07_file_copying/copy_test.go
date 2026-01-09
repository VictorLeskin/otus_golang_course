package main

import (
	"flag"
	"os"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	// Place your code here.
	_ = t
}

// Вспомогательная функция для сброса состояния флагов.
func resetFlag() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func TestParseCommadLine(t *testing.T) {
	type sExpectedData struct {
		input, output string
		offset, limit int64
	}

	tests := []struct {
		name         string
		args         []string
		expectError  string
		expectedData sExpectedData
	}{
		{
			name:         "valid parameters",
			args:         []string{"program", "-from", "in.txt", "-to", "out.txt", "-offset", "5", "-limit", "10"},
			expectError:  "",
			expectedData: sExpectedData{"in.txt", "out.txt", 5, 10},
		},
		{
			name:        "missing input",
			args:        []string{"program", "-to", "out.txt"},
			expectError: "there is not name of the file to read from",
		},
		{
			name:        "missing output",
			args:        []string{"program", "-from", "in.txt"},
			expectError: "there is not name of the file to copy",
		},
		{
			name:         "only required parameters",
			args:         []string{"program", "-from", "in.txt", "-to", "out.txt"},
			expectError:  "",
			expectedData: sExpectedData{"in.txt", "out.txt", 0, 0},
		},
		{
			name:         "garbage in numbers",
			args:         []string{"program", "-from", "in.txt", "-to", "out.txt", "ABCD", "EFG"},
			expectError:  "",
			expectedData: sExpectedData{"in.txt", "out.txt", 0, 0},
		},
		{
			name:        "missed mandatory parameters",
			args:        []string{"program"},
			expectError: "there is not name of the file to read from",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFlag()

			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args

			result, err := ParseCommadLine()

			if tt.expectError == "" {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedData.input, result.input)
				assert.Equal(t, tt.expectedData.output, result.output)
				assert.Equal(t, tt.expectedData.offset, result.offset)
				assert.Equal(t, tt.expectedData.limit, result.limit)
			} else {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectError, err.Error())
			}
		})
	}
}
