package account

import "gin-rest-api-example/internal/account/model"

type UserResponse struct {
	User User `json:"user"`
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

func NewUserResponse(acc *model.Account) *UserResponse {
	return &UserResponse{
		User: User{
			Username: acc.Username,
			Email:    acc.Email,
			Bio:      acc.Bio,
			Image:    acc.Image,
		},
	}
}
