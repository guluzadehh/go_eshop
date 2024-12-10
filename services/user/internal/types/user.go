package types

import (
	"time"

	"github.com/guluzadehh/go_eshop/services/user/internal/domain/models"
)

type UserView struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

func NewUser(u *models.User) *UserView {
	if u == nil {
		return nil
	}

	return &UserView{
		Id:        u.Id,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		IsActive:  u.IsActive,
	}
}
