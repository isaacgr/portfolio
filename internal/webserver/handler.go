package webserver

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/isaacgr/portfolio/internal/blog"
)

const (
	Title       = "Integrated Concepts"
	Sitename    = "irowell.io"
	Description = "Portfolio and blog for Isaac Rowell, written with Go+HTMX"
)

type RouteRegistrar interface {
	Register(mux *http.ServeMux)
}

type Handler struct {
	mux       *http.ServeMux
	responder *responder
	domains   []RouteRegistrar
}

type RootHandler struct {
	responder       Responder
	log             *slog.Logger
	uri             string
	renderer        *Renderer
	blogFinder      *blog.BlogFinder
	mailClient      *MailClient
	recaptchaParams ReCaptchaParams
}

func (h *Handler) MountAll() {
	for _, d := range h.domains {
		d.Register(h.mux)
	}
}

func NewHandler(
	r Responder,
	logger *slog.Logger,
	renderer *Renderer,
	blogFinder *blog.BlogFinder,
	mailClient *MailClient,
	recaptchaParams ReCaptchaParams,
) *RootHandler {
	return &RootHandler{
		responder:       r,
		log:             logger,
		renderer:        renderer,
		uri:             "/",
		blogFinder:      blogFinder,
		mailClient:      mailClient,
		recaptchaParams: recaptchaParams,
	}
}

func (h *RootHandler) Register(mux *http.ServeMux) {
	staticFs, err := fs.Sub(Assets, "html/static")
	if err != nil {
		h.log.Error("Cannot parse static assets")
		return
	}

	mux.HandleFunc("/", h.handleIndex)
	mux.HandleFunc("/home", h.handleHome)
	mux.HandleFunc("/experience", h.handleExperience)
	mux.HandleFunc("/contact", h.handleContact)
	mux.HandleFunc("/blog", h.handleBlog)
	mux.HandleFunc("/blog/{slug}", h.handleBlogPost)
	// I want to support the old blogs that all reference _resources
	mux.HandleFunc("/blog/{slug}/resources/", h.handleBlogPostResources)
	mux.HandleFunc("/blog/{slug}/_resources/", h.handleBlogPostResources)
	mux.Handle(
		"/static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.FS(staticFs)),
		),
	)
}

func (s *RootHandler) handleIndex(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	var buf bytes.Buffer
	var err error
	homeData := map[string]any{
		"Articles": s.blogFinder.Articles,
	}
	err = s.renderer.Render(&buf, "home", homeData)
	data := map[string]any{
		"Title":       Title,
		"Sitename":    Sitename,
		"Description": Description,
		"Body":        template.HTML(buf.String()),
	}
	err = s.renderer.Render(w, "base", data)
	if err != nil {
		s.respondError(w, r, err, http.StatusInternalServerError, "")
	}
}

func (s *RootHandler) handleHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *RootHandler) handleExperience(w http.ResponseWriter, r *http.Request) {
	h := r.Header.Get("HX-Request")
	var buf bytes.Buffer
	var err error
	if h != "" {
		err = s.renderer.Render(w, "experience", nil)
	} else {
		err = s.renderer.Render(&buf, "experience", nil)
		data := map[string]any{
			"Title":    "Integrated Concepts",
			"Sitename": "Integrated Concepts",
			"Body":     template.HTML(buf.String()),
		}
		err = s.renderer.Render(w, "base", data)
	}
	if err != nil {
		s.respondError(w, r, err, http.StatusInternalServerError, "")
	}
}

func (s *RootHandler) handleContact(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h := r.Header.Get("HX-Request")
		var buf bytes.Buffer
		var err error
		if h != "" {
			err = s.renderer.Render(w, "contact", nil)
		} else {
			err = s.renderer.Render(&buf, "contact", nil)
			data := map[string]any{
				"Title":       Title,
				"Sitename":    Sitename,
				"Description": Description,
				"Body":        template.HTML(buf.String()),
			}
			err = s.renderer.Render(w, "base", data)
		}
		if err != nil {
			s.respondError(w, r, err, http.StatusInternalServerError, "")
		}
	case http.MethodPost:
		s.contactSubmit(w, r)
	}
}

func (s *RootHandler) handleBlog(w http.ResponseWriter, r *http.Request) {
	posts := map[string]any{
		"Articles": s.blogFinder.Articles,
	}
	h := r.Header.Get("HX-Request")
	var buf bytes.Buffer
	var err error
	if h != "" {
		err = s.renderer.Render(w, "blog", posts)
	} else {
		err = s.renderer.Render(&buf, "blog", posts)
		data := map[string]any{
			"Title":       Title,
			"Sitename":    Sitename,
			"Description": Description,
			"Body":        template.HTML(buf.String()),
		}
		err = s.renderer.Render(w, "base", data)
	}
	if err != nil {
		s.respondError(w, r, err, http.StatusInternalServerError, "")
	}
}

func (s *RootHandler) handleBlogPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	slug = url.PathEscape(slug)
	article, ok := s.blogFinder.ArticlesBySlug[slug]
	if !ok {
		e := Error{
			Code: http.StatusBadRequest,
			Msg:  "Post not found",
		}
		s.renderer.Render(w, "error", e)
		return
	}
	post := map[string]any{
		"Title":     article.FrontMatter.Title,
		"Description":     article.FrontMatter.Description,
		"Tags":      article.FrontMatter.TagList,
		"Post":      template.HTML(article.HTML),
		"Archived":  article.Archived,
		"CreatedAt": article.FrontMatter.CreatedAt,
	}
	var buf bytes.Buffer
	var err error
	h := r.Header.Get("HX-Request")
	if h != "" {
		err = s.renderer.Render(w, "article", post)
	} else {
		err = s.renderer.Render(&buf, "article", post)
		data := map[string]any{
			"Title":       Title,
			"Sitename":    Sitename,
			"Description": Description,
			"Body":        template.HTML(buf.String()),
		}
		err = s.renderer.Render(w, "base", data)
	}
	if err != nil {
		s.respondError(w, r, err, http.StatusInternalServerError, "")
	}
}

func (s *RootHandler) handleBlogPostResources(
	w http.ResponseWriter,
	r *http.Request,
) {
	slug := r.PathValue("slug")
	slug = url.PathEscape(slug)
	article, ok := s.blogFinder.ArticlesBySlug[slug]
	if !ok {
		s.respondError(
			w,
			r,
			errors.New("Post not found"),
			http.StatusBadRequest,
			"",
		)
		return
	}
	// some/path/to/markdown.md but we need some/path/to/resources/
	filepath := article.Filepath
	rp := "/resources/"
	if strings.Contains(r.URL.Path, "/_resources/") {
		rp = "/_resources/"
	}
	resourcePath := strings.Split(filepath, fmt.Sprintf(
		"/%s",
		article.Filename,
	))[0] + rp
	buf, err := os.ReadFile(resourcePath + path.Base(r.URL.Path))

	if err != nil {
		s.log.Error("Unable to serve resources", "Error", err)
		http.NotFound(w, r)
		return
	}

	ext := path.Ext(r.URL.Path)
	ftSplit := strings.Split(ext, ".")

	if len(ftSplit) < 2 {
		s.respondError(
			w,
			r,
			errors.New("File not found"),
			http.StatusBadRequest,
			r.URL.Path,
		)
		return
	}
	ft := ftSplit[1]

	w.Header().Set("Content-Type", fmt.Sprintf("image/%s", ft))

	if ft == "svg" {
		w.Header().Set("Content-Type", "image/svg+xml")
	}

	w.Write(buf)
}

func (s *RootHandler) contactSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		s.log.Error("Unable to parse contact form. ", "Error", err.Error())
		s.respondError(
			w,
			r,
			errors.New("Unable to submit contact info."),
			http.StatusInternalServerError,
			"",
		)
	} else {
		s.log.Info("Got contact request.")
		if r.FormValue("contact_me_by_fax_only") != "" {
			s.log.Info("Got bot contact request. Discarding.")
			s.renderer.Render(w, "success", nil)
			return
		}
		captchaToken := r.FormValue("g-recaptcha-response")
		score, err := CreateAssessment(
			s.recaptchaParams.ProjectID,
			s.recaptchaParams.Key,
			captchaToken,
			s.recaptchaParams.Action,
			s.recaptchaParams.Credentials,
		)
		if err != nil {
			s.log.Error("Unable to generate reCAPTCHA score", "Error", err)
			s.respondError(
				w,
				r,
				errors.New("Unable to verify reCAPTCHA"),
				http.StatusBadRequest,
				"",
			)
			return
		}
		if score < 0.8 {
			s.log.Error("reCAPTCHA score failed", "Value", score)
			s.respondError(
				w,
				r,
				errors.New("reCAPTCHA score failed"),
				http.StatusBadRequest,
				"",
			)
			return
		}
		cd := NewContact(
			r.FormValue("name"),
			r.FormValue("email"),
			r.FormValue("subject"),
			r.FormValue("message"),
		)
		err = s.mailClient.SendEmail(cd)
		if err != nil {
			s.log.Error("Unable to send email.", "Error", err)
			s.respondError(
				w,
				r,
				err,
				http.StatusBadRequest,
				"",
			)
		} else {
			s.log.Info("Email sent successfully.")
			s.renderer.Render(w, "success", nil)
		}
	}
}

func (s *RootHandler) respondError(
	w http.ResponseWriter,
	r *http.Request,
	err error,
	status int,
	data string,
) {
	var errBuf bytes.Buffer
	s.renderer.Render(&errBuf, "error", Error{
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
