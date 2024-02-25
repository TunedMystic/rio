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

func NewLogger(w io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(w, nil))
}

func LogDebug(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

func LogInfo(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

func LogWarn(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

func LogError(msg string, attrs ...slog.Attr) {
	defaultLogger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}
