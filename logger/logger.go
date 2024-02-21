package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

var defaultLogger = NewTextLogger(os.Stdout)

func Debug(msg string, attrs ...slog.Attr) {
	defaultLogger.Debug(msg, attrs...)
}

func Info(msg string, attrs ...slog.Attr) {
	defaultLogger.Info(msg, attrs...)
}

func Warn(msg string, attrs ...slog.Attr) {
	defaultLogger.Warn(msg, attrs...)
}

func Error(msg string, attrs ...slog.Attr) {
	defaultLogger.Error(msg, attrs...)
}

type Logger struct {
	logger *slog.Logger
}

func NewTextLogger(w io.Writer) *Logger {
	return &Logger{
		logger: slog.New(slog.NewTextHandler(w, nil)),
	}
}

func NewJsonLogger(w io.Writer) *Logger {
	return &Logger{
		logger: slog.New(slog.NewJSONHandler(w, nil)),
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
