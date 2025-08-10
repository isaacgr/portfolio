package web

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Status       int
	ErrorMessage string
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
	s.Server.GET("/contact", Contact)
	s.Server.POST("/contact", ContactSubmit)
	s.Server.Static("/static", "web/static")
}

func Index(c echo.Context) error {
	data := map[string]any{
		"Title":    "Integrated Concepts",
		"Sitename": "Integrated Concepts",
	}
	return c.Render(http.StatusOK, "base", data)
}

func Contact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact", nil)
}

func ContactSubmit(c echo.Context) error {
	_, err := c.FormParams()
	if err != nil {
		log.Error("Unable to parse contact form. ", "Error", err.Error())
		return c.Render(
			http.StatusInternalServerError,
			"error",
			Error{
				Status:       http.StatusInternalServerError,
				ErrorMessage: "Unable to submit contact info.",
			},
		)
	}
	log.Info("Got contact request.")

	return c.Render(http.StatusOK, "success", nil)
}
