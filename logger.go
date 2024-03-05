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

// Logger sets the default logger to the provided slog.Logger.
// .
func Logger(l *slog.Logger) {
	defaultLogger = l
}

// NewLogger constructs a new *slog.Logger with the given writer.
// .
func NewLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, nil))
}

// LogDebug logs a debug message using the default logger.
// .
func LogDebug(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

// LogDebug logs an info message using the default logger.
// .
func LogInfo(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

// LogDebug logs a warning message using the default logger.
// .
func LogWarn(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

// LogDebug logs an error using the default logger.
// .
func LogError(err error, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelError, err.Error(), attrs...)
}
