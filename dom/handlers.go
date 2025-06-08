package dom

import "net/http"

// HandlerFunc is a custom http handler signature which accepts
// an http.ResponseWriter, *http.Request and returns a Node.
// HandlerFuncs must be converted into an http.Handler with the Handler middleware.
type HandlerFunc func(http.ResponseWriter, *http.Request) Node

// Handler is a middleware which converts a dom.HandlerFunc to an http.Handler.
// It centralizes the error handling and rendering of the Node.
func Handler(next HandlerFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		node := next(w, r)

		if node == nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)

		} else if err := node.Render(w); err != nil {
			status := http.StatusInternalServerError
			http.Error(w, http.StatusText(status), status)
		}
	}
	return http.HandlerFunc(fn)
}
