package blog

import (
	"bytes"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// Parse the directory
// Find all markdown files
// Parse the markdown files for front matter
// Fallback to first h1 header (#) for title
// Fallback to file stats for created, updated etc.
// Allow specifying additional file types, unknown how to handle metadata
type BlogFinder struct {
	Articles        []Article
	ArticlesBySlug  map[string]Article
	location        string
	fileTypes       []string
	foundFiles      []File
	log             *slog.Logger
	noteChan        chan Article
	foundFilesChan  chan File
	ignoreFileNames []string
	parsedNotesChan chan Article
}

type File struct {
	path  string
	entry fs.DirEntry
}

type Article struct {
	FrontMatter FrontMatter
	Content     string
	HTML        string
	Slug        string
	Filepath    string
	Filename    string
}

type FrontMatter struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	TagList     []string // Same as TagString but converted to a list
	CreatedAt   string   `yaml:"created_at"`
	UpdatedAt   string   `yaml:"updated_at"`
	TagString   string   `yaml:"tags"` // Comma separated list of tags, improper yaml format
}

func NewBlogFinder(
	location string,
	fileTypes []string,
	ignoreFileNames []string,
	parsedNotesChan chan Article,
	logger *slog.Logger,
) *BlogFinder {
	return &BlogFinder{
		location:        location,
		fileTypes:       fileTypes,
		ignoreFileNames: ignoreFileNames,
		log:             logger,
		// NOTE: 0.2s performance increase by buffering to >= number of notes
		noteChan:        make(chan Article, 1024),
		foundFilesChan:  make(chan File, 1024),
		parsedNotesChan: parsedNotesChan,
		ArticlesBySlug:  make(map[string]Article),
	}
}

func (n *BlogFinder) Start() {
	go n.dispatch()
	n.buildNotes()

}

func (n *BlogFinder) dispatch() {
	for {
		select {
		case note, ok := <-n.noteChan:
			if !ok {
				n.log.Error("Unable to read from noteChan")
			} else {
				n.Articles = append(n.Articles, note)
				n.ArticlesBySlug[note.Slug] = note
				n.parsedNotesChan <- note
			}
		case file, ok := <-n.foundFilesChan:
			if !ok {
				n.log.Error("Unable to read from foundFilesChan")
			} else {
				go n.parseNotes(file)
			}
		}
	}
}

func (n *BlogFinder) buildNotes() error {
	err := filepath.WalkDir(n.location, n.parseFileOrDir)
	if err != nil {
		return err
	}
	if len(n.foundFiles) == 0 {
		n.log.Warn(
			"No files found matching extensions.",
			"Extensions",
			n.fileTypes,
		)
	}
	n.log.Info("Done finding files", "Found files", len(n.foundFiles))
	return nil
}

func (n *BlogFinder) parseNotes(file File) {
	rawFile, err := os.ReadFile(file.path)
	if err != nil {
		n.log.Error(
			"Unable to read file.",
			"File",
			file.path,
			"Error",
			err,
		)
		return
	}
	var fm FrontMatter
	content, err := frontmatter.Parse(bytes.NewReader(rawFile), &fm)
	if err != nil {
		n.log.Error(
			"Unable to parse front matter.",
			"File",
			file.path,
			"Error",
			err,
		)
		return
	}
	fm.TagList = strings.Split(fm.TagString, ",")
	slug := strings.ReplaceAll(fm.Title, " ", "_")

	filenameIdx := strings.LastIndex(file.path, "/")
	filename := file.path[filenameIdx+1:]

	html := n.mdToHTML(content, slug)

	note := Article{
		FrontMatter: fm,
		Content:     string(content),
		HTML:        html,
		Slug:        slug,
		Filepath:    file.path,
		Filename:    filename,
	}
	n.log.Info(
		"Successfully parsed note",
		"File",
		file.entry.Name(),
		"Title",
		fm.Title,
	)
	n.noteChan <- note
}

// Implements fs.WalkDirFunc to parse found files in a directory
func (n *BlogFinder) parseFileOrDir(
	path string,
	d fs.DirEntry,
	err error,
) error {
	if err != nil {
		n.log.Error(
			"Error accessing path",
			"Path",
			path,
			"Error",
			err,
		)
		return fs.SkipDir
	}
	if !d.IsDir() {
		// We have hit a file, is it a valid type?
		i, err := d.Info()
		if err != nil {
			n.log.Error(
				"Error reading file info. Skipping",
				"File",
				d.Name(),
				"Error",
				err,
			)
			return nil
		}
		ft := filepath.Ext(i.Name())
		if slices.Contains(n.fileTypes, ft) && !slices.Contains(
			n.ignoreFileNames,
			i.Name(),
		) {
			n.log.Info(
				"Found file",
				"File",
				i.Name(),
			)
			f := File{
				path:  path,
				entry: d,
			}
			n.foundFiles = append(n.foundFiles, f)
			n.foundFilesChan <- f
			return nil
		}
	}
	return nil
}

func (n *BlogFinder) mdToHTML(content []byte, slug string) string {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)
	doc = modifyAst(doc, slug)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	//h := markdown.ToHTML(content, nil, nil)
	h := markdown.Render(doc, renderer)
	return string(h)
}
