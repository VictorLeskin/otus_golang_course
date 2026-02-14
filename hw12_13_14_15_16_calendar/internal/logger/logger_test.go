package logger

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLogLevel(t *testing.T) {
	testCases := []struct {
		level string
		res   bool
	}{
		{"debug", true},
		{"info", true},
		{"warning", true},
		{"error", true},
		{"fatal", true},
		{"panic", true},
		{"invalid", false},
		{"", false},
		{"trace", false},
	}

	for _, tc := range testCases {
		t.Run(tc.level, func(t *testing.T) {
			err := ValidateLogLevel(tc.level)

			if tc.res {
				assert.Nil(t, err)
			}
			if !tc.res {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestLogger_println0(t *testing.T) {
	var buf strings.Builder // уже имеет Write([]byte) (int, error)

	t0 := Logger{output: &buf}

	t0.println0("[A]", "Wello, horld!")

	assert.Equal(t, "[A] Wello, horld!\n", buf.String())
}

func TestLogger_Println(t *testing.T) {
	var buf strings.Builder // уже имеет Write([]byte) (int, error)

	t0 := Logger{output: &buf, loggingLevel: 1}

	t0.Println(0, "Level 0")
	t0.Println(1, "Level 1")
	t0.Println(2, "Level 2")

	assert.Equal(t, "[I] Level 1\n[W] Level 2\n", buf.String())
}

func TestLogger_Printf(t *testing.T) {

	var buf strings.Builder // уже имеет Write([]byte) (int, error)
	t0 := Logger{output: &buf, loggingLevel: 1}

	t0.Printf(0, "Level %d %s", 0, "A")
	t0.Printf(1, "Level %d %s", 1, "B")
	t0.Printf(2, "Level %d %s", 2, "C")

	assert.Equal(t, "[I] Level 1 B\n[W] Level 2 C\n", buf.String())
}

func TestLogger_Message(t *testing.T) {
	tests := []struct {
		loggingLever int
		funcName     string
		result       string
	}{
		{0, "Debug", "[D] message\n"},
		{1, "Debug", ""},
		{0, "Info", "[I] message\n"},
		{1, "Info", "[I] message\n"},
		{2, "Info", ""},
		{1, "Warning", "[W] message\n"},
		{2, "Warning", "[W] message\n"},
		{3, "Warning", ""},
		{2, "Error", "[E] message\n"},
		{3, "Error", "[E] message\n"},
		{4, "Error", ""},
		{3, "Fatal", "[F] message\n"},
		{4, "Fatal", "[F] message\n"},
		{5, "Fatal", ""},
		{4, "Panic", "[P] message\n"},
		{5, "Panic", "[P] message\n"},
	}

	for _, tc := range tests {
		var buf strings.Builder
		t0 := Logger{output: &buf, loggingLevel: tc.loggingLever}
		switch tc.funcName {
		case "Debug":
			t0.Debug("message")
		case "Info":
			t0.Info("message")
		case "Warning":
			t0.Warning("message")
		case "Error":
			t0.Error("message")
		case "Fatal":
			t0.Fatal("message")
		case "Panic":
			t0.Panic("message")
		default:
			panic(fmt.Sprintf("Not such function %s", tc.funcName))
		}

		assert.Equal(t, tc.result, buf.String())
	}
}

func TestLogger_Messagef(t *testing.T) {
	tests := []struct {
		loggingLever int
		funcName     string
		result       string
	}{
		{0, "Debug", "[D] message 4 Q\n"},
		{1, "Debug", ""},
		{0, "Info", "[I] message 4 Q\n"},
		{1, "Info", "[I] message 4 Q\n"},
		{2, "Info", ""},
		{1, "Warning", "[W] message 4 Q\n"},
		{2, "Warning", "[W] message 4 Q\n"},
		{3, "Warning", ""},
		{2, "Error", "[E] message 4 Q\n"},
		{3, "Error", "[E] message 4 Q\n"},
		{4, "Error", ""},
		{3, "Fatal", "[F] message 4 Q\n"},
		{4, "Fatal", "[F] message 4 Q\n"},
		{5, "Fatal", ""},
		{4, "Panic", "[P] message 4 Q\n"},
		{5, "Panic", "[P] message 4 Q\n"},
	}
	i := 4
	s := "Q"

	for _, tc := range tests {
		var buf strings.Builder
		t0 := Logger{output: &buf, loggingLevel: tc.loggingLever}
		switch tc.funcName {
		case "Debug":
			t0.Debugf("message %d %s", i, s)
		case "Info":
			t0.Infof("message %d %s", i, s)
		case "Warning":
			t0.Warningf("message %d %s", i, s)
		case "Error":
			t0.Errorf("message %d %s", i, s)
		case "Fatal":
			t0.Fatalf("message %d %s", i, s)
		case "Panic":
			t0.Panicf("message %d %s", i, s)
		default:
			panic(fmt.Sprintf("Not such function %s", tc.funcName))
		}

		assert.Equal(t, tc.result, buf.String())
	}
}
