package web

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type templateRenderer struct {
	templates *template.Template
}

func newTemplateRenderer(
	baseDir string,
) (*templateRenderer, error) {
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
	r := &templateRenderer{
		templates: root,
	}
	return r, nil
}

func (r *templateRenderer) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
