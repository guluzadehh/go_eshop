package authhttp

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/http/handlers"
	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/render"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/services/auth"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string) (access string, refresh string, err error)
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"

	log := sl.ForHandler(h.Log, op, requestmdw.GetRequestId(r))

	var req Request
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("can't decode json", sl.Err(err))
		h.JSON(w, http.StatusBadRequest, api.Err("failed to parse request body"))
		return
	}

	access, refresh, err := h.srvc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			h.JSON(w, http.StatusUnauthorized, api.Err("invalid credentials"))
			return
		}

		h.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
		return
	}

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     h.cfg.JWT.Refresh.CookieName,
			Value:    refresh,
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true,
			Path:     h.cfg.JWT.Refresh.Uri,
			MaxAge:   int(h.cfg.JWT.Refresh.Expire.Seconds()),
			Secure: func(env string) bool {
				if env == "prod" {
					return true
				} else {
					return false
				}
			}(h.cfg.Env),
		},
	)

	h.JSON(w, http.StatusOK, Response{
		Response: api.Ok(),
		Data: &Data{
			Token: access,
		},
	})
}
