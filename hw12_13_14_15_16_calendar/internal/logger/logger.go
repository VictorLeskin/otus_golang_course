package logger

import (
	"fmt"
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

func (l Logger) Info(msg string) {
	fmt.Println(msg)
}

func (l Logger) Error(msg string) {
	// TODO
}

// TODO
