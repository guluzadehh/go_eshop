package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/domain/models"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/service"
	"github.com/guluzadehh/go_eshop/services/user/internal/storage"
)

type UserProvider interface {
	UserById(ctx context.Context, id int) (*models.User, error)
	UserByEmail(ctx context.Context, email string) (*models.User, error)
}

type UserService struct {
	log          *slog.Logger
	config       *config.Config
	userProvider UserProvider
}

func New(log *slog.Logger, config *config.Config, userProvider UserProvider) *UserService {
	return &UserService{
		log:          log,
		config:       config,
		userProvider: userProvider,
	}
}

func (s *UserService) SetLog(log *slog.Logger) {
	s.log = log
}

func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	const op = "services.user.GetUser"

	log := s.log.With(slog.String("op", op))

	user, err := s.userProvider.UserById(ctx, id)
	if err != nil {
		if errors.Is(err, storage.UserNotFound) {
			log.Info("user not found", slog.Int("user_id", id))
			return nil, service.ErrUserNotFound
		}

		log.Error("couldn't get the user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "services.user.GetUserByEmail"

	log := s.log.With(slog.String("op", op))

	user, err := s.userProvider.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.UserNotFound) {
			log.Info("user not found")
			return nil, service.ErrUserNotFound
		}

		log.Error("couldn't get the user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
