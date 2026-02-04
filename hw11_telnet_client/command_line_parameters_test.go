package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_parseCommadLine(t *testing.T) {
	/*
		{
			clp, err := parseCommandLine([]string{"-invalid-flag", "192.168.1.1", "8080"})
			assert.Nil(t, err)
			assert.Equal(t, "10.2.92.212", clp.host)
			assert.Equal(t, 8888, clp.port)
			assert.Equal(t, 21*time.Second, clp.timeout)
		}
	*/

	tests := []struct {
		name        string
		args        []string
		wantHost    string
		wantPort    int
		wantTimeout time.Duration
		wantErr     bool
		wantErrStr  string // return error message
	}{
		{
			name:        "valid IPv4 with default timeout",
			args:        []string{"192.168.1.1", "8080"},
			wantHost:    "192.168.1.1",
			wantPort:    8080,
			wantTimeout: 10 * time.Second,
			wantErr:     false,
		},
		{
			name:        "valid IPv4 with custom timeout",
			args:        []string{"-timeout", "5s", "10.0.0.1", "443"},
			wantHost:    "10.0.0.1",
			wantPort:    443,
			wantTimeout: 5 * time.Second,
			wantErr:     false,
		},
		{
			name:        "min valid port",
			args:        []string{"192.168.1.1", "1"},
			wantHost:    "192.168.1.1",
			wantPort:    1,
			wantTimeout: 10 * time.Second,
			wantErr:     false,
		},
		{
			name:        "max valid port",
			args:        []string{"192.168.1.1", "65535"},
			wantHost:    "192.168.1.1",
			wantPort:    65535,
			wantTimeout: 10 * time.Second,
			wantErr:     false,
		},

		{
			name:       "missing host and port",
			args:       []string{},
			wantErr:    true,
			wantErrStr: "host and port are required",
		},
		{
			name:       "invalid host format",
			args:       []string{"not-an-ip", "8080"},
			wantErr:    true,
			wantErrStr: "wrong host address",
		},
		{
			name:       "port is not a number",
			args:       []string{"192.168.1.1", "abc"},
			wantErr:    true,
			wantErrStr: "port must be a number",
		},
		{
			name:       "port is float",
			args:       []string{"192.168.1.1", "80.5"},
			wantErr:    true,
			wantErrStr: "port must be a number",
		},
		{
			name:       "port too small",
			args:       []string{"192.168.1.1", "0"},
			wantErr:    true,
			wantErrStr: "port number must be in range [1,65535]",
		},
		{
			name:       "port too large",
			args:       []string{"192.168.1.1", "65536"},
			wantErr:    true,
			wantErrStr: "port number must be in range [1,65535]",
		},
		{
			name:       "invalid flag",
			args:       []string{"-invalid-flag", "192.168.1.1", "8080"},
			wantErr:    true,
			wantErrStr: "error parsing command line parameters:\nflag provided but not defined: -invalid-flag",
		},
		{
			name:       "host with too few octets",
			args:       []string{"192.168.1", "8080"},
			wantErr:    true,
			wantErrStr: "wrong host address",
		},
		{
			name:        "valid host but extra arguments",
			args:        []string{"192.168.1.1", "8080", "extra", "args"},
			wantHost:    "192.168.1.1",
			wantPort:    8080,
			wantTimeout: 10 * time.Second,
			wantErr:     false, // Extra arguments will be ignored
		},
		{
			name:       "invalid timeout value",
			args:       []string{"-timeout", "invalid", "192.168.1.1", "80"},
			wantErr:    true,
			wantErrStr: "error parsing command line parameters:\ninvalid value \"invalid\" for flag -timeout: parse error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCommandLine(tt.args)

			// check error if any
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Equal(t, tt.wantErrStr, err.Error())
				return
			}

			// no error check results
			if !tt.wantErr {
				assert.Equal(t, got.host, tt.wantHost)
				assert.Equal(t, got.port, tt.wantPort)
				assert.Equal(t, got.timeout, tt.wantTimeout)
			}
		})
	}
}
