package portfolio

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/isaacgr/portfolio/internal/api"
	"github.com/isaacgr/portfolio/internal/api/responder"
	"github.com/isaacgr/portfolio/internal/blog"
	"github.com/isaacgr/portfolio/internal/server"
)

type PortfolioApi struct {
	responder       responder.Responder
	log             *slog.Logger
	uri             string
	renderer        *server.TemplateRenderer
	blogFinder      *blog.BlogFinder
	mailClient      *server.MailClient
	recaptchaParams server.ReCaptchaParams
	staticDir       string
}

func NewPortfolioApi(
	r responder.Responder,
	logger *slog.Logger,
	renderer *server.TemplateRenderer,
	blogFinder *blog.BlogFinder,
	mailClient *server.MailClient,
	recaptchaParams server.ReCaptchaParams,
	staticDir string,
) *PortfolioApi {
	return &PortfolioApi{
		responder:       r,
		log:             logger,
		renderer:        renderer,
		uri:             "/",
		blogFinder:      blogFinder,
		mailClient:      mailClient,
		recaptchaParams: recaptchaParams,
		staticDir:       staticDir,
	}
}

func (s *PortfolioApi) getTarget(path string) string {
	baseUrl := s.uri
	return baseUrl + path
}

func (s *PortfolioApi) RegisterRoutes(handler *http.ServeMux) {
	handler.HandleFunc(s.getTarget(""), s.handleIndex)
	handler.HandleFunc("/home", s.handleHome)
	handler.HandleFunc("/experience", s.handleExperience)
	handler.HandleFunc("/contact", s.handleContact)
	handler.HandleFunc("/blog", s.handleBlog)
	handler.HandleFunc("/blog/{slug}", s.handleBlogPost)
	// I want to support the old blogs that all reference _resources
	handler.HandleFunc("/blog/{slug}/resources/", s.handleBlogPostResources)
	handler.HandleFunc("/blog/{slug}/_resources/", s.handleBlogPostResources)
	handler.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.Dir(
				s.staticDir,
			)),
		),
	)
}

func (s *PortfolioApi) respondError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	status int,
	data string,
) {
	var errBuf bytes.Buffer
	w.WriteHeader(status)
	s.renderer.Render(&errBuf, "error", api.Error{
		Msg:  err.Error(),
		Code: status,
		Data: data,
	})
	s.responder.HTML(
		w,
		http.StatusInternalServerError,
		errBuf.String(),
	)
}
