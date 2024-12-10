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
}

type ProfileService struct {
	log             *slog.Logger
	config          *config.Config
	profileProvider ProfileProvider
}

func New(log *slog.Logger, config *config.Config, profileProvider ProfileProvider) *ProfileService {
	return &ProfileService{
		log:             log,
		config:          config,
		profileProvider: profileProvider,
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
