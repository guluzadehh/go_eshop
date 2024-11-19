package authhttp

import "github.com/guluzadehh/go_eshop/services/user/internal/lib/api"

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	api.Response
	Data *Data `json:"data,omitempty"`
}

type Data struct {
	Token string `json:"access_token"`
}
