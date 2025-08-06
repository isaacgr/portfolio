package articles

import (
	"errors"
)

type Article struct {
	Title       string
	Description string
	Date        string
	Content     []byte
}

func ConvertArticles(articles []Article) ([][]byte, error) {
	return nil, errors.New("Unable to convert articles")
}
