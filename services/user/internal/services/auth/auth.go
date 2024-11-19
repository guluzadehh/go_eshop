package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/domain/models"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/jwt"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
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

func (s *AuthService) Login(ctx context.Context, email string, password string) (string, string, error) {
	const op = "services.auth.Login"

	log := s.log.With(slog.String("op", op))

	user, err := s.userProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.UserNotFound) {
			log.Warn("user doesn't exist")
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Warn("invalid credentials", slog.Int64("user_id", user.Id))
			return "", "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to compare passwords", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	access, err := jwt.AccessToken(user, s.config)
	if err != nil {
		log.Error("failed to generate access token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	refresh, err := jwt.RefreshToken(user, s.config)
	if err != nil {
		log.Error("failed to generate refresh token", sl.Err(err))
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return access, refresh, nil
}
