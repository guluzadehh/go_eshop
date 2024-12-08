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
	ErrEmailExists        = errors.New("email is already taken")
)

type UserProvider interface {
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserSaver interface {
	CreateUser(ctx context.Context, email string, password string) (*models.User, error)
}

type AuthService struct {
	log          *slog.Logger
	config       *config.Config
	userProvider UserProvider
	userSaver    UserSaver
}

func New(log *slog.Logger, config *config.Config, userProvider UserProvider, userSaver UserSaver) *AuthService {
	return &AuthService{
		log:          log,
		config:       config,
		userProvider: userProvider,
		userSaver:    userSaver,
	}
}

func (s *AuthService) SetLog(log *slog.Logger) {
	s.log = log
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

func (s *AuthService) Signup(ctx context.Context, email string, password string) (*models.User, error) {
	const op = "services.auth.Signup"

	log := s.log.With(slog.String("op", op))

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	user, err := s.userSaver.CreateUser(ctx, email, string(bytes))
	if err != nil {
		if errors.Is(err, storage.UserExists) {
			log.Info("email is taken")
			return nil, ErrEmailExists
		}

		log.Error("couldn't save the user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user has been created")

	return user, nil
}
