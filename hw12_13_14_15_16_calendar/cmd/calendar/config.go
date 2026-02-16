package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"calendar/internal/logger"
	internalhttp "calendar/internal/server/http"
	"calendar/internal/storage"
	sqlstorage "calendar/internal/storage/sql"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

// Config основная структура конфигурации.
type Config struct {
	Logger     logger.Config       `json:"logger"`
	Server     internalhttp.Config `json:"server"`
	Storage    storage.Config      `json:"storage"`
	SQLStorage sqlstorage.Config   `json:"sqlstorage"`
	// add confings for other subparts of project.
}

var ErrInvalidConfig = fmt.Errorf("not valid config")

// NewDefaultConfig возвращает конфиг со значениями по умолчанию.
func NewDefaultConfig() *Config {
	return &Config{
		Logger: logger.Config{
			Level: "info",
			File:  "calendar.log",
		},
	}
}

func ReadAll(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	return io.ReadAll(file)
}

// LoadConfig загружает конфигурацию из JSON файла.
func LoadConfig(path string) (*Config, error) {
	data, err := ReadAll(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Проверяем, что файл не пустой.
	if len(data) == 0 {
		return nil, fmt.Errorf("config file is empty")
	}

	cfg := *NewDefaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	if err = ValidateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	// Валидация уровня логирования.
	if err := logger.ValidateLogLevel(cfg.Logger.Level); err != nil {
		return ErrInvalidConfig
	}
	return nil
}
