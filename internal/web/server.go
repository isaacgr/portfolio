package web

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/isaacgr/portfolio/internal/pkg/logging"
	"github.com/labstack/echo/v4"
)

var log = logging.GetLogger("web", false)

type WebServer struct {
	Server     *echo.Echo
	Context    echo.Context
	config     *Configuration
	signalChan chan os.Signal
}

type Template struct {
	templates *template.Template
}

func NewWebServer(c *Configuration) *WebServer {
	s := echo.New()
	s.Debug = true
	return &WebServer{
		Server: s,
		config: c,
	}
}

func (s *WebServer) Start() {
	s.registerSignals()
	s.registerRenderer()
	s.RegisterRoutes()
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := s.Server.Start(address)
	if err != nil {
		log.Error("Unable to start server", "Error", err)
	}
}

func (t *Template) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func (s *WebServer) registerRenderer() {
	r, err := newTemplateRenderer(
		"web/views",
	)
	if err != nil {
		log.Error("Unable to parse templates.", "Error", err.Error())
	}
	s.Server.Renderer = r
}

func (s *WebServer) registerSignals() {
	s.signalChan = make(chan os.Signal, 1)
	signal.Notify(s.signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s.signalChan

		shutdownCtx, shutdownRelease := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer shutdownRelease()
		if err := s.Server.Shutdown(shutdownCtx); err != nil {
			log.Error("HTTP shutdown error", "Error", err)
		}

		log.Info("Server shutdown complete")
	}()
}

func newTemplateRenderer(
	baseDir string,
) (*Template, error) {
	root := template.New("")
	if _, err := root.ParseFiles(
		filepath.Join(baseDir, "base.html"),
	); err != nil {
		return nil, err
	}
	patterns := []string{
		// TODO: Handle directories with no files in them
		filepath.Join(baseDir, "components", "*.html"),
		filepath.Join(baseDir, "pages", "*.html"),
		filepath.Join(baseDir, "templates", "*.html"),
	}
	for _, pat := range patterns {
		if _, err := root.ParseGlob(pat); err != nil {
			return nil, err
		}
	}
	r := &Template{
		templates: root,
	}
	return r, nil
}
