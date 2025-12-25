package logger

import (
	"context"
	"fmt"
	"time"
)

type LoggerInterface interface {
	Info(message string) string
	Warn(message string) string
	Error(message string) string
}

type Logger struct {
	Module   *string
	MinLevel LogLevels
}

type LogLevels int

const (
	INFO = iota
	WARNING
	ERROR
)

type contextKey string

const loggerContextKey contextKey = "logger"

var loggingLevels = map[LogLevels]string{
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
}

func NewLogger(module *string, minLevelToLog LogLevels) *Logger {
	return &Logger{
		Module: module,
	}
}

func (l *Logger) WithLogger(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

func GetLoggerFromContext(ctx context.Context, defaultModuleIfMissing *string) LoggerInterface {
	retrievedLogger := ctx.Value(loggerContextKey)

	if retrievedLogger == nil {
		return NewLogger(defaultModuleIfMissing, INFO)
	}

	// Try to cast to LoggerInterface
	if logger, ok := retrievedLogger.(LoggerInterface); ok {
		return logger
	}

	return NewLogger(defaultModuleIfMissing, INFO)
}

func (l *Logger) log(level LogLevels, message string) string {
	if level < l.MinLevel {
		return ""
	}

	prefix := loggingLevels[level] + " "
	if l.Module != nil {
		prefix = *l.Module + ": "
	}

	now := time.Now()

	messageToLog := prefix + message + "\n[" + now.Format("2006-01-02 15:04:05") + "]"

	fmt.Println(messageToLog)

	return messageToLog
}

func (l *Logger) Info(message string) string {
	return l.log(INFO, message)
}

func (l *Logger) Warn(message string) string {
	return l.log(WARNING, message)
}

func (l *Logger) Error(message string) string {
	return l.log(ERROR, message)
}
