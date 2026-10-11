package server

import (
	"github.com/isaacgr/portfolio/internal/webserver"
)

type Server struct {
	webServer *webserver.WebServer
}

func NewServer(
	webServer *webserver.WebServer,
) *Server {
	return &Server{
		webServer: webServer,
	}
}

func (s *Server) Start() {
	s.webServer.Start()
}
