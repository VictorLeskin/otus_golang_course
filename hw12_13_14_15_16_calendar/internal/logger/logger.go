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
	loggingLevel int
}

func New(config LoggerConfig) *Logger {
	return &Logger{
		loggingLevel: validLevels[config.Level],
	}
}

var validLevels = map[string]int{
	"debug":   0,
	"info":    1,
	"warning": 2,
	"error":   3, // Логирует ошибку	Ошибка в бизнес-логике.
	"fatal":   4, // Логирует и завершает программу	Невозможно продолжать работу.
	"panic":   5, // Логирует и вызывает panic	Программная ошибка, баг.
}

// validateLogLevel проверяет корректность уровня логирования.
func ValidateLogLevel(level string) error {
	lowerLevel := strings.ToLower(level)
	if _, exist := validLevels[lowerLevel]; !exist {
		return fmt.Errorf("invalid log level: %s. Valid values: debug, info, warning, error, fatal, panic", level)
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
