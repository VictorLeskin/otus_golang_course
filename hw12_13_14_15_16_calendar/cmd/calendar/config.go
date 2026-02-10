package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

// Config основная структура конфигурации.
type Config struct {
	Logger LoggerConfig `json:"logger"`
	// add confings for other subparts of project.
}

// NewDefaultConfig возвращает конфиг со значениями по умолчанию.
func NewDefaultConfig() *Config {
	return &Config{
		Logger: LoggerConfig{
			Level: "info",
			File:  "calendar.log",
		},
	}
}

// LoggerConfig настройки логгера.
type LoggerConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
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

	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	// Валидация уровня логирования.
	if err := validateLogLevel(cfg.Logger.Level); err != nil {
		return err
	}
	return nil
}

// validateLogLevel проверяет корректность уровня логирования.
func validateLogLevel(level string) error {
	validLevels := map[string]bool{
		"debug":   true,
		"info":    true,
		"warning": true,
		"error":   true, // Логирует ошибку	Ошибка в бизнес-логике.
		"fatal":   true, // Логирует и завершает программу	Невозможно продолжать работу.
		"panic":   true, // Логирует и вызывает panic	Программная ошибка, баг.
	}

	lowerLevel := strings.ToLower(level)
	if !validLevels[lowerLevel] {
		return fmt.Errorf("invalid log level: %s. Valid values: debug, info, warning, error", level)
	}

	return nil
}
