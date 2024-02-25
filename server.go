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

type Server struct {
	mux        *http.ServeMux
	middleware []func(http.Handler) http.Handler
}

func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.Use(RecoverPanic, LogRequest, SecureHeaders)
	return s
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) Use(middleware ...func(http.Handler) http.Handler) {
	s.middleware = append(s.middleware, middleware...)
}

func (s *Server) Handler() http.Handler {
	slices.Reverse(s.middleware)
	var h http.Handler = s.mux

	for i := range s.middleware {
		m := s.middleware[i]
		h = m(h)
	}
	return h
}

func (s *Server) Serve(addr string) error {
	LogInfo("starting server", slog.String("port", addr))
	return Serve(addr, s.Handler())
}

// ------------------------------------------------------------------
//
//
// Serve Helper Function
//
//
// ------------------------------------------------------------------

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
