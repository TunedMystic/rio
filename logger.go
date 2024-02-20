package rio

import (
	"context"
	"io"
	"log/slog"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger(w io.Writer) *Logger {
	opts := slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewTextHandler(w, &opts)

	return &Logger{
		logger: slog.New(handler),
	}
}

func (l *Logger) Debug(msg string, attrs ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

func (l *Logger) Info(msg string, attrs ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

func (l *Logger) Warn(msg string, attrs ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...slog.Attr) {
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}
