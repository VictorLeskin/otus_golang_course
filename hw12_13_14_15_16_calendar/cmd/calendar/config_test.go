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
			err := validateLogLevel(tc.level)

			if tc.res {
				assert.NotNil(t, err)
			}
			if !tc.res {
				assert.Nil(t, err)
			}
		})
	}
}
func TestLoadConfig(t *testing.T) {
	// Тест 1: Пустой JSON - должны остаться дефолты
	emptyJSON := `{}`
	tmpFile := createTempFile(t, emptyJSON)
	defer os.Remove(tmpFile.Name())

	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Проверяем, что остались дефолтные значения
	if cfg.Logger.Level != "info" {
		t.Errorf("Expected default level 'info', got '%s'", cfg.Logger.Level)
	}
	if cfg.Logger.File != "calendar.log" {
		t.Errorf("Expected default file 'calendar.log', got '%s'", cfg.Logger.File)
	}

	// Тест 2: Частичный JSON - часть дефолтов, часть из файла
	partialJSON := `{
        "logger": {
            "level": "debug"
            // file не указан - должен остаться дефолтный
        }
    }`
	tmpFile2 := createTempFile(t, partialJSON)
	defer os.Remove(tmpFile2.Name())

	cfg2, err := LoadConfig(tmpFile2.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// level из JSON, file из дефолтов
	if cfg2.Logger.Level != "debug" {
		t.Errorf("Expected level from JSON 'debug', got '%s'", cfg2.Logger.Level)
	}
	if cfg2.Logger.File != "calendar.log" {
		t.Errorf("Expected default file 'calendar.log', got '%s'", cfg2.Logger.File)
	}

	// Тест 3: Полный JSON - полностью переопределяем дефолты
	fullJSON := `{
        "logger": {
            "level": "error",
            "file": "prod.log"
        }
    }`
	tmpFile3 := createTempFile(t, fullJSON)
	defer os.Remove(tmpFile3.Name())

	cfg3, err := LoadConfig(tmpFile3.Name())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Все значения из JSON
	if cfg3.Logger.Level != "error" {
		t.Errorf("Expected level from JSON 'error', got '%s'", cfg3.Logger.Level)
	}
	if cfg3.Logger.File != "prod.log" {
		t.Errorf("Expected file from JSON 'prod.log', got '%s'", cfg3.Logger.File)
	}
}

func createTempFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Cannot create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Cannot write to temp file: %v", err)
	}
	tmpFile.Close()

	// Открываем заново для чтения
	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Cannot reopen temp file: %v", err)
	}

	return tmpFile
}
