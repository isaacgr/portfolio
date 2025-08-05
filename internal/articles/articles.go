package internal

import "errors"

type Article struct {
	Title   string
	Summary string
	Date    string
	Content []byte
}

func FindArticles() ([]Article, error) {
	return nil, errors.New("Error")
}

func ConvertPosts() {}
