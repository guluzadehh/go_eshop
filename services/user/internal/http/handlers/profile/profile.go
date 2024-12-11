package profilehttp

import (
	"context"
	"log/slog"

	"github.com/guluzadehh/go_eshop/services/user/internal/domain/models"
	"github.com/guluzadehh/go_eshop/services/user/internal/http/handlers"
)

type ProfileService interface {
	GetProfile(ctx context.Context, id int64) (*models.Profile, error)
	MakeProfile(ctx context.Context, id int64, firstName, lastName, phone string) (*models.Profile, error)
	DeleteProfile(ctx context.Context, id int64) error
	SetLog(log *slog.Logger)
}

type ProfileHandler struct {
	*handlers.Handler
	srvc ProfileService
}

func New(log *slog.Logger, srvc ProfileService) *ProfileHandler {
	return &ProfileHandler{
		srvc:    srvc,
		Handler: handlers.NewHandler(log),
	}
}
