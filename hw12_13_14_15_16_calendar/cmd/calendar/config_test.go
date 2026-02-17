package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultConfig(t *testing.T) {
	res := NewDefaultConfig()
	assert.Equal(t, "info", res.Logger.Level)
	assert.Equal(t, "calendar.log", res.Logger.File)
}

func TestReadAll(t *testing.T) {
	{
		emptyJSON := `{}`
		tmpFile := createTempFile(t, emptyJSON)
		defer os.Remove(tmpFile.Name())

		data, err := ReadAll(tmpFile.Name())
		assert.Nil(t, err)
		assert.Equal(t, []byte("{}"), data)
	}

	{
		_, err := ReadAll("wrong file name")
		assert.NotNil(t, err)
	}
}

func TestLoadFullApplicationConfigWithSQLStorage(t *testing.T) {
	fullJSON := `
	{
		"logger": {
			"level": "info",
			"file": "calendar.log"
		},
		"server": {
			"host": "188.33.44.66",
			"port": "9876"
		},
		"storage": { 
			"type" : "sqlstorage" 
		},
		"sqlstorage": {
			"host": "localhost",
			"port": 5432,
			"database": "calendar",
			"username": "calendar_user",
			"password": "calendar_pass"
		}
	}`

	tmpFile := createTempFile(t, fullJSON)
	defer os.Remove(tmpFile.Name())

	cfg, err := LoadConfig(tmpFile.Name())

	assert.Nil(t, err)
	assert.Equal(t, "info", cfg.Logger.Level)
	assert.Equal(t, "calendar.log", cfg.Logger.File)

	assert.Equal(t, "188.33.44.66", cfg.Server.Host)
	assert.Equal(t, "9876", cfg.Server.Port)

	assert.Equal(t, "sqlstorage", cfg.Storage.Type)

	assert.Equal(t, "localhost", cfg.SQLStorage.Host)
	assert.Equal(t, 5432, cfg.SQLStorage.Port)
	assert.Equal(t, "calendar", cfg.SQLStorage.Database)
	assert.Equal(t, "calendar_user", cfg.SQLStorage.Username)
	assert.Equal(t, "calendar_pass", cfg.SQLStorage.Password)
}

func TestLoadFullApplicationConfigWithInMemoryStorage(t *testing.T) {
	fullJSON := `
	{
		"logger": {
			"level": "info",
			"file": "calendar.log"
		},
		"server": {
			"host": "188.33.44.66",
			"port": "9876"
		},
		"storage": { 
			"type" : "inmemory" 
		}
	}`

	tmpFile := createTempFile(t, fullJSON)
	defer os.Remove(tmpFile.Name())

	cfg, err := LoadConfig(tmpFile.Name())

	assert.Nil(t, err)
	assert.Equal(t, "info", cfg.Logger.Level)
	assert.Equal(t, "calendar.log", cfg.Logger.File)

	assert.Equal(t, "188.33.44.66", cfg.Server.Host)
	assert.Equal(t, "9876", cfg.Server.Port)

	assert.Equal(t, "inmemory", cfg.Storage.Type)
}

func TestLoadConfig(t *testing.T) {
	t.Run("empty JSON: expected default", func(t *testing.T) {
		emptyJSON := `{}`
		tmpFile := createTempFile(t, emptyJSON)
		defer os.Remove(tmpFile.Name())

		cfg, err := LoadConfig(tmpFile.Name())

		assert.Nil(t, err)
		assert.Equal(t, "info", cfg.Logger.Level)
		assert.Equal(t, "calendar.log", cfg.Logger.File)
	})

	t.Run("full JSON ", func(t *testing.T) {
		fullJSON := `{
			"logger": {
				"level": "error",
				"file": "prod.log"
			}
		}`
		tmpFile := createTempFile(t, fullJSON)
		defer os.Remove(tmpFile.Name())

		cfg, err := LoadConfig(tmpFile.Name())

		assert.Nil(t, err)
		assert.Equal(t, "error", cfg.Logger.Level)
		assert.Equal(t, "prod.log", cfg.Logger.File)
	})

	t.Run("empty file", func(t *testing.T) {
		empty := ``
		tmpFile := createTempFile(t, empty)
		defer os.Remove(tmpFile.Name())

		_, err := LoadConfig(tmpFile.Name())

		assert.NotNil(t, err)
	})

	t.Run("wrong JSON Format", func(t *testing.T) {
		fullJSON := `{
			"logger": {
				"level": "error",
				"file": "prod.log"
			}
		}}`
		tmpFile := createTempFile(t, fullJSON)
		defer os.Remove(tmpFile.Name())

		_, err := LoadConfig(tmpFile.Name())

		assert.NotNil(t, err)
	})
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := NewDefaultConfig()
		err := ValidateConfig(cfg)
		assert.Nil(t, err)
	})

	t.Run("not valid config", func(t *testing.T) {
		cfg := NewDefaultConfig()
		cfg.Logger.Level = "highest"
		err := ValidateConfig(cfg)
		assert.NotNil(t, err)
	})
}

func createTempFile(t *testing.T, content string) *os.File {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Cannot create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Cannot write to temp file: %v", err)
	}
	tmpFile.Close()

	// Открываем заново для чтения.
	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Cannot reopen temp file: %v", err)
	}

	return tmpFile
}
