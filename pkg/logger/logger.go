package logger

import (
	"log/slog"
	"os"
)

// Logger defines the logging interface
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

type logger struct {
	*slog.Logger
}

// New creates a new structured logger
func New() Logger {
	// Create a JSON handler for structured logging
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	
	return &logger{
		Logger: slog.New(handler),
	}
}

func (l *logger) Info(msg string, args ...interface{}) {
	l.Logger.Info(msg, args...)
}

func (l *logger) Error(msg string, args ...interface{}) {
	l.Logger.Error(msg, args...)
}

func (l *logger) Warn(msg string, args ...interface{}) {
	l.Logger.Warn(msg, args...)
}

func (l *logger) Debug(msg string, args ...interface{}) {
	l.Logger.Debug(msg, args...)
}

func (l *logger) Fatal(msg string, args ...interface{}) {
	l.Logger.Error(msg, args...)
	os.Exit(1)
}