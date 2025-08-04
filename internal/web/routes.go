package web

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Article struct {
	Title   string
	Summary string
	Date    string
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
	s.Server.Static("/static", "web/static")
}

func Index(c echo.Context) error {
	article := Article{
		Title:   "Test Artcile",
		Summary: "This is a test Article",
		Date:    "Today",
	}
	data := map[string]any{
		"Title":    "Integrated Concepts",
		"Sitename": "Integrated Concepts",
		"Articles": []Article{article},
	}
	return c.Render(http.StatusOK, "base", data)
}
