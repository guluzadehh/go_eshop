package authhttp

import (
	"context"
	"log/slog"
	"net/http"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (access string, refresh string, err error)
}

type AuthHandler struct {
	log  *slog.Logger
	srvc AuthService
}

func New(log *slog.Logger, srvc AuthService) *AuthHandler {
	return &AuthHandler{
		log:  log,
		srvc: srvc,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}
