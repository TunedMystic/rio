package rio

import (
	"log/slog"
	"net/http"
	"slices"
	"time"
)

// ------------------------------------------------------------------
//
//
// Type: Server
//
//
// ------------------------------------------------------------------

// Server is a wrapper around the standard http.ServeMux.
//
// It can register routes for Handlers and HandlerFuncs.
// It can also register middleware for the entire ServeMux.
type Server struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

// NewServer constructs and returns a new *Server.
//
// The underlying mux is the standard http.ServeMux.
// The LogRequest, RecoverPanic and SecureHeaders middleware are
// automatically registered on construction.
func NewServer(middleware ...func(http.Handler) http.Handler) *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	if middleware == nil {
		s.Use(LogRequest(defaultLogger), RecoverPanic, SecureHeaders)
	} else {
		s.Use(middleware...)
	}
	return s
}

// Handler registers the handler for the given pattern.
// It is a proxy for http.ServeMux.Handle().
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern.
// It is a proxy for http.ServeMux.HandleFunc().
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

// Use registers one or more handlers as middleware for the Server.
func (s *Server) Use(middleware ...func(http.Handler) http.Handler) {
	s.middleware = append(s.middleware, middleware...)
}

// Handler returns the Server as an http.Handler.
//
// It wraps the ServeMux with the middleware handlers, and returns
// the result as an http.Handler.
//
// The middleware is reversed, so that the earliest registered
// middleware is wrapped last.
func (s *Server) Handler() http.Handler {
	slices.Reverse(s.middleware)
	var h http.Handler = s.mux

	for i := range s.middleware {
		m := s.middleware[i]
		h = m(h)
	}
	return h
}

// Serve starts an http server on the given address.
func (s *Server) Serve(addr string) error {
	LogInfo("starting server", slog.String("port", addr))
	return Serve(addr, s.Handler())
}

// ------------------------------------------------------------------
//
//
// Type: RouteMap
//
//
// ------------------------------------------------------------------

// RouteMap defines a mapping of routes to http handlers.
type RouteMap map[string]http.Handler

// RegisterRoutes registers all routes of the map to the Server.
func (rm RouteMap) RegisterRoutes(s *Server) {
	for route := range rm {
		s.Handle(route, rm[route])
	}
}

// ------------------------------------------------------------------
//
//
// Serve Helper Function
//
//
// ------------------------------------------------------------------

// Serve starts an http server.
//
// The addr is the address to listen to. The addr assumes the format "host:port".
// The handler is the http.Handler to serve.
func Serve(addr string, handler http.Handler) error {
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        handler,
		IdleTimeout:    time.Minute,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
	}
	return httpServer.ListenAndServe()
}
