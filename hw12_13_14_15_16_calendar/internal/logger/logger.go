package logger

import (
	"fmt"
	"strings"
)

// LoggerConfig настройки логгера.
type LoggerConfig struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type Logger struct { // TODO
}

func New(config LoggerConfig) *Logger {
	return &Logger{}
}

// validateLogLevel проверяет корректность уровня логирования.
func ValidateLogLevel(level string) error {
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

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO
