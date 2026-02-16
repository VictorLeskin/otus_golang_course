package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Config настройки логгера.
type Config struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type Logger struct { // TODO
	loggingLevel int
	output       io.Writer
	file         *os.File // храним файл отдельно, если это файл
}

func getWriter(fileName string) (output io.Writer, file *os.File, err error) {
	if fileName == "" {
		output = os.Stdout
	} else {
		f, err0 := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0o644)
		if err0 != nil {
			err = err0
		} else {
			output = f
			file = f
		}
	}

	return output, file, err
}

func New(config Config) *Logger {
	output, file, err := getWriter(config.File)
	if err != nil {
		panic(fmt.Errorf("cannot open log file: %w", err))
	}
	return &Logger{
		loggingLevel: validLevels[config.Level],
		output:       output,
		file:         file,
	}
}

func NewWriterLogger(level string, output io.Writer) *Logger {
	return &Logger{
		loggingLevel: validLevels[level],
		output:       output,
	}
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

const (
	LevelDebug   int = iota // 0
	LevelInfo               // 1
	LevelWarning            // 2
	LevelError              // 3
	LevelFatal              // 4
	LevelPanic              // 5
)

var validLevels = map[string]int{
	"debug":   LevelDebug,
	"info":    LevelInfo,
	"warning": LevelWarning,
	"error":   LevelError, // Логирует ошибку	Ошибка в бизнес-логике.
	"fatal":   LevelFatal, // Логирует и завершает программу	Невозможно продолжать работу.
	"panic":   LevelPanic, // Логирует и вызывает panic	Программная ошибка, баг.
}

var headerLevels = []string{
	"[D]",
	"[I]",
	"[W]",
	"[E]",
	"[F]",
	"[P]",
}

// validateLogLevel проверяет корректность уровня логирования.
func ValidateLogLevel(level string) error {
	lowerLevel := strings.ToLower(level)
	if _, exist := validLevels[lowerLevel]; !exist {
		return fmt.Errorf("invalid log level: %s. Valid values: debug, info, warning, error, fatal, panic", level)
	}

	return nil
}

func (l Logger) println0(lvl string, msg string) {
	fmt.Fprintln(l.output, lvl, msg)
}

func (l Logger) Println(lvl int, msg string) {
	if lvl >= l.loggingLevel {
		l.println0(headerLevels[lvl], msg)
	}
}

func (l Logger) Printf(lvl int, format string, args ...interface{}) {
	if lvl >= l.loggingLevel {
		msg := fmt.Sprintf(format, args...)
		l.println0(headerLevels[lvl], msg)
	}
}

func (l Logger) Debug(msg string) {
	l.Println(LevelDebug, msg)
}

func (l Logger) Info(msg string) {
	l.Println(LevelInfo, msg)
}

func (l Logger) Warning(msg string) {
	l.Println(LevelWarning, msg)
}

func (l Logger) Error(msg string) {
	l.Println(LevelError, msg)
}

func (l Logger) Fatal(msg string) {
	l.Println(LevelFatal, msg)
}

func (l Logger) Panic(msg string) {
	l.Println(LevelPanic, msg)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.Printf(LevelDebug, format, args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	l.Printf(LevelInfo, format, args...)
}

func (l Logger) Warningf(format string, args ...interface{}) {
	l.Printf(LevelWarning, format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Printf(LevelError, format, args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	l.Printf(LevelFatal, format, args...)
}

func (l Logger) Panicf(format string, args ...interface{}) {
	l.Printf(LevelPanic, format, args...)
}
