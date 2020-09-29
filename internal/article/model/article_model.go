package model

import (
	accountModel "gin-rest-api-example/internal/account/model"
	"time"
)

type Article struct {
	ID            uint      `gorm:"column:id"`
	Slug          string    `gorm:"column:slug"`
	Title         string    `gorm:"column:title"`
	Body          string    `gorm:"column:body"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	DeletedAtUnix int64     `gorm:"column:deleted_at_unix"`
	Author        accountModel.Account
	AuthorID      uint
	Tags          []*Tag `gorm:"many2many:article_tags;association_autocreate:false"`
}

type Tag struct {
	ID        uint      `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:created_at"`
	Articles  []Article `gorm:"many2many:article_tags;"`
}

type Comment struct {
	ID        uint   `gorm:"column:id"`
	Body      string `gorm:"column:body"`
	Slug      string `gorm:"column:slug"`
	Author    accountModel.Account
	AuthorID  uint
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}
