package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email    string `gorm:"column:email;unique_index"`
	Username string `gorm:"column:username;unique_index"`
	Password string `gorm:"column:password; not null"`
	Bio      string `gorm:"column:bio;size:1024"`
	Image    string
}

type Follow struct {
	Follower    User
	FollowerID  uint
	Following   User
	FollowingID uint
	CreatedAt   time.Time
}
