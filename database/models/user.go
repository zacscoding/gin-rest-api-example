package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"column:email;unique_index"`
	Username string `json:"username" gorm:"column:username;unique_index"`
	Password string `json:"password" gorm:"column:password; not null"`
	Bio      string `json:"bio" gorm:"column:bio;size:1024"`
	Image    string `json:"image"`
}

type Follow struct {
	Follower    User
	FollowerID  uint
	Following   User
	FollowingID uint
	CreatedAt   time.Time
}

