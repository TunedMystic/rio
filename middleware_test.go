package rio

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func TestLogRequest(t *testing.T) {
	var logBuf bytes.Buffer
	logger := testLogger(&logBuf)
	handler := LogRequest(logger)(testHandlerOk)

	req := httptest.NewRequest("GET", "/test_log", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Body.String(), "OK")
	assert.True(t, strings.Contains(logBuf.String(), "level=INFO msg=request status=200 method=GET url=/test_log"))
}

func TestSkipLogger(t *testing.T) {
	var logBuf bytes.Buffer
	logger := testLogger(&logBuf)

	middleware := SkipLogger(logger, "/skip_me")
	handler := middleware(testHandlerOk)

	t.Run("path should be logged", func(t *testing.T) {
		logBuf.Reset()
		req := httptest.NewRequest("GET", "/log_me", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.True(t, strings.Contains(logBuf.String(), "level=INFO msg=request status=200 method=GET url=/log_me"))
	})

	t.Run("path should be skipped", func(t *testing.T) {
		logBuf.Reset()
		req := httptest.NewRequest("GET", "/skip_me", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, logBuf.String(), "")
	})
}

func TestRecoverPanic(t *testing.T) {
	req := httptest.NewRequest("GET", "/panic_test", nil)

	t.Run("handler panics", func(t *testing.T) {
		handlerWithPanic := RecoverPanic(testHandlerPanic)
		rr := httptest.NewRecorder()

		handlerWithPanic.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusInternalServerError)
		assert.Equal(t, rr.Body.String(), "Internal Server Error\n")
		assert.Equal(t, rr.Header().Get("Connection"), "close")
	})

	t.Run("handler does not panic", func(t *testing.T) {
		handlerWithoutPanic := RecoverPanic(testHandlerOk)
		rr := httptest.NewRecorder()

		handlerWithoutPanic.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, rr.Body.String(), "OK")
	})
}

func TestSecureHeaders(t *testing.T) {
	handler := SecureHeaders(testHandlerOk)
	req := httptest.NewRequest("GET", "/secure", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	headers := rr.Header()
	assert.Equal(t, headers.Get("Referrer-Policy"), "origin-when-cross-origin")
	assert.Equal(t, headers.Get("Strict-Transport-Security"), "max-age=31536000; includeSubDomains; preload")
	assert.Equal(t, headers.Get("X-Content-Type-Options"), "nosniff")
	assert.Equal(t, headers.Get("X-Frame-Options"), "deny")
	assert.Equal(t, headers.Get("X-XSS-Protection"), "0")
}

func TestCacheControl(t *testing.T) {
	handler := CacheControl(testHandlerOk)
	req := httptest.NewRequest("GET", "/cache_default", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Header().Get("Cache-Control"), "max-age=172800")
}

func TestCacheControlWithAge(t *testing.T) {
	age := 3600
	middleware := CacheControlWithAge(age)
	handler := middleware(testHandlerOk)

	req := httptest.NewRequest("GET", "/cache_custom_age", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, rr.Header().Get("Cache-Control"), fmt.Sprintf("max-age=%d", age))
}

func TestNotFound(t *testing.T) {
	handler404 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Custom Not Found"))
	})
	handlerRoot := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Root OK"))
	})

	middleware := NotFound(handler404, handlerRoot)

	t.Run("path is /", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, rr.Body.String(), "Root OK")
	})

	t.Run("path is not /", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/some_other_path", nil)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusNotFound)
		assert.Equal(t, rr.Body.String(), "Custom Not Found")
	})
}

// ------------------------------------------------------------------
//
// Test Helpers
//
// ------------------------------------------------------------------

// testLogger creates a logger that writes to a buffer and
// removes the timestamp for predictable log output in tests.
func testLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
}

// testHandlerOk is an http.Handler that writes HTTP 200 OK.
var testHandlerOk = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
})

// testHandlerPanic is an http.Handler that always panics.
var testHandlerPanic = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	panic(errors.New("test panic"))
})
