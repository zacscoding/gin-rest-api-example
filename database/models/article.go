package models

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
)

type Article struct {
	gorm.Model
	Slug        string `gorm:"unique_index; not null"`
	Title       string
	Description string `gorm:"size:2048"`
	Body        string `gorm:"size:2048"`
	Author      User
	AuthorID    uint
	Tags        []Tag     `gorm:"many2many:article_tags;"`
	Comment     []Comment `gorm:"ForeignKey:ArticleID"`
}

func (a *Article) UpdateSlug() {
	a.Slug = slug.Make(a.Title)
}

type ArticleFavorite struct {
	User      User
	UserID    uint
	Article   Article
	ArticleID uint
}

type Tag struct {
	gorm.Model
	Name     string    `gorm:"unique_index"`
	Articles []Article `gorm:"many2many:article_tags;"`
}

type ArticleTag struct {
	Article   Article
	ArticleID uint
	Tag       Tag
	TagID     uint
}

type Comment struct {
	gorm.Model
	Body      string `gorm:""`
	Article   Article
	ArticleID uint
	Author    User
	AuthorID  uint
}
