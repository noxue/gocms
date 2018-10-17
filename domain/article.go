package domain

import "time"

type Article struct {
	Name  string `validate:"required,min=1,max=200,excludesall=!/#?@&+="`
	Title string `validate:"required,min=1,max=100"`
	// article type name
	Type        string `validate:"required,min=1,max=100,excludesall=!/#?@&+="`
	Description string
	Content     string `validate:"required,min=1"`
	Sort        int    `validate:"min=-99999,max=99999"`
	// previous article
	Prev string
	// next article
	Next string
	Good int
	Top  int
	// article tags
	Tags      []string
	Hits      int
	Author    string
	Chapters  []*Article
	CreatedAt time.Time
}