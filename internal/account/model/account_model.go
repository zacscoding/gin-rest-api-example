package model

import (
	"encoding/json"
	"fmt"
)

type Account struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (a Account) String() string {
	return fmt.Sprintf("{%d %s %s %s}", a.ID, a.Username, a.Email, "[PROTECTED]")
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
