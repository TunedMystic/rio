package rio

import (
	"errors"
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
type logResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *logResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// LogRequest is a middleware which logs the http request and response status.
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

// SkipLogger is a middleware which logs the http request and response status
// if the request url does not match the given path.
func SkipLogger(excludePath string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// If the url matches the excludePath,
			// the request will not be logged.
			if r.URL.Path == excludePath {
				next.ServeHTTP(w, r)
				return
			}

			// If the url does not match the excludePath,
			// the request will be logged.
			LogRequest(next).ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
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
func RecoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// The deferred function will always run,
		// even in the event of a panic.
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				LogError(err.(error))
				Http500(w)
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
func SecureHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
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
func CacheControl(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=172800")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// CacheControlWithAge is a middleware which sets the caching policy for assets.
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

// ------------------------------------------------------------------
//
//
// NotFound Middleware
//
//
// ------------------------------------------------------------------

// NotFound is a middleware which returns a 404 Not Found error
// if the request path is not "/".
func NotFound(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			Http404(w, http.StatusText(http.StatusNotFound))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// NotFound is a middleware which returns a 404 Not Found json error
// if the request path is not "/".
func NotFoundJson(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			Json404(w, http.StatusText(http.StatusNotFound))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// ------------------------------------------------------------------
//
//
// MakeHandler Middleware
//
//
// ------------------------------------------------------------------

// HandlerFunc is a custom http handler signature which accepts
// an http.ResponseWriter, *http.Request and returns an error.
// HandlerFuncs must be converted into an http.Handler with the MakeHandler middleware.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// MakeHandler is a middleware which converts a rio.HandlerFunc to an http.Handler.
// It centralizes the error handling with the custom AppError error type.
func MakeHandler(next HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Run the handler and check for errors.
		if err := next(w, r); err != nil {
			// If the error is an AppError, then write it to the ResponseWriter.
			var appErr AppError
			if errors.As(err, &appErr) {
				if writeErr := appErr.WriteTo(w); writeErr != nil {
					LogError(writeErr)
					Http500(w)
				}
				return
			}
			// If the error is NOT an AppError, then log it
			// and return a generic Http 500.
			LogError(err)
			Http500(w)
		}
	}
	return http.HandlerFunc(fn)
}
