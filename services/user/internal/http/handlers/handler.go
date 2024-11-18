package handlers

import (
	"log/slog"

	"github.com/guluzadehh/go_eshop/services/user/internal/lib/render"
)

type Handler struct {
	Log *slog.Logger
	*render.Responder
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		Log:       log,
		Responder: render.NewResponder(log),
	}
}
