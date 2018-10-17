package domain

import "time"

// article type
type ArticleType struct {
	Name string `validate:"required,min=1,max=100,excludesall=/!#?@&+="`
	Title string `validate:"required,min=1,max=100"`
	Sort int `validate:"min=-99999,max=99999"`
	CreatedAt time.Time
}
