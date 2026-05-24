package blog

import (
	"strings"

	"github.com/gomarkdown/markdown/ast"
)

func modifyAst(doc ast.Node, slug string) ast.Node {
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if img, ok := node.(*ast.Image); ok && entering {
			dst := img.Destination
			if dst == nil {
				return ast.GoToNext
			}
			// Blog posts will reference images in a local folder like
			// ./resources/image.png or maybe resources/image.png
			// The handler will expect to serve the image from a url like
			// /blog/{slug}/resources/image.png
			// Update the image href to point to ./{slug}/resources/test.png
			prefix := slug + "/"
			if strings.HasPrefix(string(dst), "./") {
				s := strings.SplitAfter(string(dst), "./")
				dst = append([]byte(prefix), []byte(s[1])...)
			} else if strings.HasPrefix(string(dst), "resources") {
				dst = append([]byte(prefix), dst...)
			} else if strings.HasPrefix(string(dst), "_resources") {
				dst = append([]byte(prefix), dst...)
			}
			img.Destination = dst
		}
		return ast.GoToNext
	})
	return doc
}
