package internalhttp

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server http.Server
	app    Application
}

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.StatusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello, World!"))
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(ctx context.Context, id, title string) error
}

func NewServer(addr string, logger Logger, app Application) *Server {
	mux := http.NewServeMux()

	helloHandler := http.HandlerFunc(hello)
	mux.Handle("/", loggingMiddleware(helloHandler, logger))

	return &Server{
		server: http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
