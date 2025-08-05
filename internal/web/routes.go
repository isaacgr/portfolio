package web

import (
	"html/template"
	"net/http"

	internal "github.com/isaacgr/portfolio/internal/articles"
	"github.com/labstack/echo/v4"
)

type Error struct {
	Status        int
	StatusMessage string
	Body          string
}

var templ *template.Template

func (s *WebServer) RegisterRoutes() {
	r, err := newTemplateRenderer(
		"web/views",
	)
	if err != nil {
		log.Error("Unable to parse templates.", "Error", err.Error())
	}
	s.Server.Renderer = r
	s.Server.GET("/", Index)
	s.Server.GET("/blog", Blog)
	s.Server.Static("/static", "web/static")
}

func Index(c echo.Context) error {
	data := map[string]any{
		"Title":    "Integrated Concepts",
		"Sitename": "Integrated Concepts",
	}
	return c.Render(http.StatusOK, "base", data)
}

func Blog(c echo.Context) error {
	articles, err := internal.FindArticles()
	if err != nil {
		e := Error{
			Status:        500,
			StatusMessage: "Internal Server Error",
			Body:          "Unable to fetch articles",
		}
		return c.Render(
			http.StatusInternalServerError,
			"error",
			e,
		)
	}
	return c.Render(http.StatusOK, "articles", articles)
}
