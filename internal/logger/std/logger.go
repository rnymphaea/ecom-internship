// Package std provides a logger implementation based on the standard library slog.
package std

import (
	"log/slog"
	"os"

	"ecom-internship/internal/logger"
)

// StdLogger implements the Logger interface using slog.
//
//nolint:revive
type StdLogger struct {
	logger *slog.Logger
}

// New creates a new StdLogger instance with the specified log level.
func New(lvl string) *StdLogger {
	var level slog.Level
	switch lvl {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return &StdLogger{
		logger: slog.New(handler),
	}
}

// Debug logs a message at DEBUG level.
func (l *StdLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

// Info logs a message at INFO level.
func (l *StdLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

// Warn logs a message at WARN level.
func (l *StdLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

// Error logs a message at ERROR level.
func (l *StdLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

// With adds attributes to the logger.
func (l *StdLogger) With(args ...any) logger.Logger {
	return &StdLogger{
		logger: l.logger.With(args...),
	}
}
