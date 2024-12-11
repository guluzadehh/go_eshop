package profile

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

type ProfileProvider interface {
	ProfileById(ctx context.Context, id int64) (*models.Profile, error)
	SaveProfile(ctx context.Context, id int64, firstName, lastName, phone string) (*models.Profile, error)
}

type ProfileDeleter interface {
	DeleteProfile(ctx context.Context, id int64) error
}

type ProfileService struct {
	log             *slog.Logger
	config          *config.Config
	profileProvider ProfileProvider
	profileDeleter  ProfileDeleter
}

func New(log *slog.Logger, config *config.Config, profileProvider ProfileProvider, profileDeleter ProfileDeleter) *ProfileService {
	return &ProfileService{
		log:             log,
		config:          config,
		profileProvider: profileProvider,
		profileDeleter:  profileDeleter,
	}
}

func (s *ProfileService) SetLog(log *slog.Logger) {
	s.log = log
}

func (s *ProfileService) GetProfile(ctx context.Context, userId int64) (*models.Profile, error) {
	const op = "services.profile.GetProfile"

	log := s.log.With(slog.String("op", op))

	profile, err := s.profileProvider.ProfileById(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ProfileNotFound) {
			log.Info("user not found", slog.Int64("user_id", userId))
			return nil, service.ErrProfileNotFound
		}

		log.Error("couldn't get the user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return profile, nil
}

func (s *ProfileService) DeleteProfile(ctx context.Context, userId int64) error {
	const op = "services.profile.DeleteProfile"

	log := s.log.With(slog.String("op", op))

	if err := s.profileDeleter.DeleteProfile(ctx, userId); err != nil {
		if errors.Is(err, storage.UserNotFound) {
			return service.ErrUserNotFound
		}

		log.Error("failed to delete user", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user has been deleted", slog.Int64("user_id", userId))

	return nil
}

func (s *ProfileService) MakeProfile(
	ctx context.Context,
	userId int64,
	firstName string,
	lastName string,
	phone string,
) (*models.Profile, error) {
	const op = "services.profile.MakeProfile"

	log := s.log.With(slog.String("op", op))

	profile, err := s.profileProvider.SaveProfile(ctx, userId, firstName, lastName, phone)
	if err != nil {
		if errors.Is(err, storage.UserNotFound) {
			log.Info("user for profile doesn't exist")
			return nil, service.ErrUserNotFound
		}

		log.Error("failed to create profile for user", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("profile has been created")

	return profile, nil
}
