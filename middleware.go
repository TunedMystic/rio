package rio

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/tunedmystic/rio/logger"
	"github.com/tunedmystic/rio/utils"
)

// logResponseWriter allows us to capture the response status code.
// .
type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := &logResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Defer the logging call.
		defer func(start time.Time) {
			logger.Info(
				"request",
				slog.Int("status", ww.status),
				slog.String("method", r.Method),
				slog.String("url", r.URL.RequestURI()),
				slog.Duration("time", time.Since(start)),
			)
		}(time.Now())

		// Call the next handler
		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}

func RecoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				utils.Http500(w, err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func SecureHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
