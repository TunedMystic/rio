package dom

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tunedmystic/rio/internal/assert"
)

func TestHandlerMiddleware(t *testing.T) {
	t.Run("DomHandler", func(t *testing.T) {
		fn := Handler(func(w http.ResponseWriter, r *http.Request) Node {
			return Div(Text("Hello, World!"))
		})

		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		fn.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusOK)
		assert.Equal(t, rr.Body.String(), `<div>Hello, World!</div>`)
	})

	t.Run("DomHandler with nil Node", func(t *testing.T) {
		fn := Handler(func(w http.ResponseWriter, r *http.Request) Node {
			return nil
		})

		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		fn.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusInternalServerError)
		assert.Equal(t, rr.Body.String(), "Internal Server Error\n")
	})

	t.Run("DomHandler with error", func(t *testing.T) {
		fn := Handler(func(w http.ResponseWriter, r *http.Request) Node {
			return errorNode{}
		})

		req, _ := http.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		fn.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, http.StatusInternalServerError)
		assert.Equal(t, rr.Body.String(), "Internal Server Error\n")
	})
}
