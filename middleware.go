package rio

import (
	"log/slog"
	"net/http"
)

func LogRequest(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func RecoverPanic(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			// Create a deferred function (which will always be run in the event
			// of a panic as Go uwinds the stack).
			defer func() {

				// Use the builtin recover function to check
				// if there has been a panic or not.
				if err := recover(); err != nil {
					w.Header().Set("Connection", "close")

					logger.Error("Panic recovered")
					status := http.StatusInternalServerError
					http.Error(w, http.StatusText(status), status)
				}
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
