package logger

import (
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

func TestLogger(t *testing.T) {
	// TODO
}
