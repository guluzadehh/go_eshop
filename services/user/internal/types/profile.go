package types

import "github.com/guluzadehh/go_eshop/services/user/internal/domain/models"

type ProfileView struct {
	*UserView
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Picture   string `json:"picture"`
}

func NewProfile(p *models.Profile) *ProfileView {
	if p == nil {
		return nil
	}

	return &ProfileView{
		UserView:  NewUser(&p.User),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Phone:     p.Phone,
		Picture:   p.Picture,
	}
}
