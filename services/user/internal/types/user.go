package types

import "github.com/guluzadehh/go_eshop/services/user/internal/domain/models"

type UserView struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
}

func NewUser(u *models.User) *UserView {
	if u == nil {
		return nil
	}

	return &UserView{
		Id:    u.Id,
		Email: u.Email,
	}
}
