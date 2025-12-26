package std

import (
	"context"
	"log/slog"
	"os"

	"ecom-internship/internal/logger"
)

type StdLogger struct {
	logger *slog.Logger
}

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

func (l *StdLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *StdLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *StdLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *StdLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *StdLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *StdLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *StdLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *StdLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

func (l *StdLogger) With(args ...any) logger.Logger {
	return &StdLogger{
		logger: l.logger.With(args...),
	}
}
