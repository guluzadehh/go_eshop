package auth

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/domain/models"
)

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

type AuthService struct {
	log          *slog.Logger
	config       *config.Config
	userProvider UserProvider
}

func New(log *slog.Logger, config *config.Config, userProvider UserProvider) *AuthService {
	return &AuthService{
		log:          log,
		config:       config,
		userProvider: userProvider,
	}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (access string, refresh string, err error) {
	return "", "", nil
}
