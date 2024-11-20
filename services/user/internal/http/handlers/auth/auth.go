package authhttp

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/http/handlers"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (access string, refresh string, err error)
	SetLog(log *slog.Logger)
}

type AuthHandler struct {
	*handlers.Handler
	cfg  *config.Config
	srvc AuthService
}

func New(log *slog.Logger, config *config.Config, srvc AuthService) *AuthHandler {
	return &AuthHandler{
		srvc:    srvc,
		cfg:     config,
		Handler: handlers.NewHandler(log),
	}
}
