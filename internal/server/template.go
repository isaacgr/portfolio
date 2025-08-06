package server

import (
	"html/template"
	"io"
	"path/filepath"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(
	w io.Writer,
	name string,
	data any,
) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplateRenderer(
	baseDir string,
) (*TemplateRenderer, error) {
	funcMap := template.FuncMap{
		"safe": RenderSafeComments,
	}
	root := template.New("").Funcs(funcMap)
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
	r := &TemplateRenderer{
		templates: root,
	}
	return r, nil
}

func RenderSafeComments(s string) template.HTML {
	return template.HTML(s)
}
