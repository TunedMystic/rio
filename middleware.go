package rio

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// ------------------------------------------------------------------
//
//
// LogRequest Middleware
//
//
// ------------------------------------------------------------------

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

// Logger is a middleware which logs the http request and response status.
// .
func LogRequest(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ww := &logResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Defer the logging call.
		defer func(start time.Time) {
			LogInfo(
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

// ------------------------------------------------------------------
//
//
// RecoverPanic Middleware
//
//
// ------------------------------------------------------------------

// RecoverPanic is a middleware which recovers from panics and
// logs a HTTP 500 (Internal Server Error) if possible.
// .
func RecoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				Http500(w, err.(error))
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// ------------------------------------------------------------------
//
//
// SecureHeaders Middleware
//
//
// ------------------------------------------------------------------

// SecureHeaders is a middleware which adds HTTP security headers
// to every response, inline with current OWASP guidance.
// .
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

// ------------------------------------------------------------------
//
//
// CacheControl Middleware
//
//
// ------------------------------------------------------------------

// CacheControl is a middleware which sets the caching policy for assets.
// Defaults to 2 days.
// .
func CacheControl(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=172800")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// CacheControlWithAge is a middleware which sets the caching policy for assets.
// .
func CacheControlWithAge(age int) func(http.Handler) http.Handler {
	maxAge := fmt.Sprintf("max-age=%d", age)

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", maxAge)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
