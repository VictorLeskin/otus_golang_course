package hw02unpackstring

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDebugPackage(_ *testing.T) {
	result, err := Unpack(`\a`)
	if err != nil {
		fmt.Print("Ups..")
	}

	if result != `` {
		fmt.Print("Ups..#2")
	}
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "aaф0b", expected: "aab"},
		{input: `✅ Use strings.Repeat("=", n) 🎯 `, expected: `✅ Use strings.Repeat("=", n) 🎯 `},
		// uncomment if task with backslash completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `\\\3`, expected: `\3`},
		{input: `\\3`, expected: `\\\`},
		{input: `\0`, expected: "0"},
		{input: `\\0`, expected: ""},
		{input: `\\\0`, expected: `\0`},
		{input: `\\\\0`, expected: `\`},
		{input: `\\\\\0`, expected: `\\0`},
		{input: "\n", expected: "\n"},
		{input: "\n2", expected: "\n\n"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `\a`, `\\\a`, `\\\a`, "12345", "1a2"}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
