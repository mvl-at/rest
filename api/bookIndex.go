package api

import "github.com/mvl-at/model"

type BookIndex struct {
	Books []Book
}

type Book struct {
	Title string
	Scores []*model.Archive
}
