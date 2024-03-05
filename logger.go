package rio

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// ------------------------------------------------------------------
//
//
// Default Logger
//
//
// ------------------------------------------------------------------

var defaultLogger = NewLogger(os.Stdout)

// Logger sets the default logger to the given slog.Logger.
func Logger(l *slog.Logger) {
	defaultLogger = l
}

// NewLogger constructs and returns a new *slog.Logger.
func NewLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, nil))
}

// LogDebug logs a debug message.
func LogDebug(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

// LogInfo logs an info message.
func LogInfo(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

// LogWarn logs a warning message.
func LogWarn(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

// LogError logs an error.
func LogError(err error, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelError, err.Error(), attrs...)
}
