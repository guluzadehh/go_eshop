package authhttp

import "github.com/guluzadehh/go_eshop/services/user/internal/lib/api"

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRes struct {
	api.Response
	Data *LoginData `json:"data,omitempty"`
}

type LoginData struct {
	Token string `json:"access_token"`
}
