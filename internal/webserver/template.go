package webserver

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"time"
)

type Renderer struct {
	log           *slog.Logger
	foundFiles    []File
	templates     *template.Template
	templateCache map[string][]byte
}

func (t *Renderer) Render(
	w io.Writer,
	name string,
	data any,
) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func RenderSafeComments(s string) template.HTML {
	return template.HTML(s)
}

func FormatDate(s string) string {
	if s == "" {
		return ""
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("Jan 2, 2006")
	}
	return s
}

func NewRenderer(log *slog.Logger) *Renderer {
	funcMap := template.FuncMap{
		// {{ safe <!--Some HTML Comment--> }}
		"safe":       RenderSafeComments,
		"formatDate": FormatDate,
	}
	root := template.New("").Funcs(funcMap)
	return &Renderer{
		log:           log,
		templates:     root,
		templateCache: make(map[string][]byte),
	}
}

// parseFileOrDir implements fs.WalkDirFunc to parse found files in a directory
func (r *Renderer) parseFileOrDir(
	path string,
	d fs.DirEntry,
	err error,
) error {
	if err != nil {
		r.log.Warn(
			"Unable to walk path, skipping directory",
			"Error",
			err,
		)
		return fs.SkipDir
	}
	if !d.IsDir() {
		i, err := d.Info()
		if err != nil {
			// Error reading file info, skip
			r.log.Warn("Unable to read file info, skipping", "Error", err)
			return nil
		}
		ft := filepath.Ext(i.Name())
		if slices.Contains([]string{".html"}, ft) {
			f := File{
				path:  path,
				entry: d,
			}
			r.log.Info("Found valid file", "File", d.Name())
			r.foundFiles = append(r.foundFiles, f)
			return nil
		}
	}
	return nil
}

// ParseTemplates finds and parses the template files in the embeded HTML dir
func (r *Renderer) ParseTemplates() error {
	err := fs.WalkDir(HTML, ".", r.parseFileOrDir)
	if err != nil {
		return err
	}
	if len(r.foundFiles) == 0 {
		return fmt.Errorf("No templates to parse")
	}

	for _, f := range r.foundFiles {
		if _, err := r.templates.ParseFS(HTML, f.path); err != nil {
			return err
		}
	}
	return nil
}
