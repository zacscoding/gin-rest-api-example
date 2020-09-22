package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Account struct {
	ID       uint   `gorm:"column:id"`
	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
	Bio      string `gorm:"column:bio"`
	Image    string `gorm:"column:image"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	Disabled  bool      `gorm:"column:disabled"`
}

func (a Account) String() string {
	return fmt.Sprintf("Account{id:%d, username:%s, password:%s, bio:%s, image:%s, createdAt:%v, updatedAt:%v, disabled:%v",
		a.ID, a.Username, "[PROTECTED]", a.Bio, a.Image, a.CreatedAt, a.UpdatedAt, a.Disabled)
}

func (a *Account) UnmarshalJSON(b []byte) error {
	var tmp struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	a.ID = tmp.ID
	a.Username = tmp.Username
	a.Email = tmp.Email
	a.Password = tmp.Password
	return nil
}
