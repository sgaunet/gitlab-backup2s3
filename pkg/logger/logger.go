package logger

import (
	"io"
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// NewLogger creates a new logger
// logLevel is the level of logging
// Possible values of logLevel are: "debug", "info", "warn", "error"
// Default value is "info"
func NewLogger(logLevel string) Logger {
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: false,
	})
	logger := slog.New(logHandler)
	return logger
}

// NoLogger creates a logger that does not log anything
func NoLogger() Logger {
	noLogger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))
	return noLogger
}
