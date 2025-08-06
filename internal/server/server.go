package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type WebServer struct {
	Handler *http.ServeMux
	server  *http.Server
	done    chan struct{}
	context context.Context
	Errors  chan error
	log     *slog.Logger
}

func NewWebServer(
	host string,
	port int,
	context context.Context,
	logger *slog.Logger,
) *WebServer {
	mux := http.NewServeMux()
	addr := fmt.Sprintf("%s:%d", host, port)

	s := &WebServer{
		Handler: mux,
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		done:    make(chan struct{}),
		Errors:  make(chan error, 1024),
		context: context,
		log:     logger,
	}
	return s
}

func (s *WebServer) Start() {
	go s.startServer()
	go s.dispatch()
}

func (s *WebServer) Stop() {
	if err := s.server.Shutdown(s.context); err != nil {
		s.log.Error("HTTP shutdown error", "Error", err)
	}
	s.log.Info("Server shutdown complete")
	close(s.done)
}

func (s *WebServer) WaitUntilClosed() {
	<-s.done
}

func (s *WebServer) startServer() {
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.log.Error("Unable to start server", "Error", err)
	}

}
func (s *WebServer) dispatch() {
	for {
		select {
		case err, ok := <-s.Errors:
			if !ok {
				s.Stop()
				return
			}
			s.handleError(err)
		}
	}
}

func (s *WebServer) handleError(err error) {
	s.log.Error(
		"Web server encountered an error",
		"Error",
		err,
	)
}
