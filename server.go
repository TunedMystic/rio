package rio

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Server struct {
	logger *Logger
	mux    *http.ServeMux
}

func NewServer() *Server {
	return &Server{
		logger: NewTextLogger(os.Stdout),
		mux:    http.NewServeMux(),
	}
}

func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(pattern, handler)
}

func (s *Server) Mux() *http.ServeMux {
	return s.mux
}

func (s *Server) Logger() *Logger {
	return s.logger
}

func (s *Server) Routes() http.Handler {
	LogReq := LogRequest(s.logger)
	Recover := RecoverPanic(s.logger)
	return Recover(LogReq(SecureHeaders(s.mux)))
}

func (s *Server) Serve(addr string) error {
	s.logger.Info("starting server", slog.String("port", addr))
	return Serve(addr, s.Routes())
}

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
