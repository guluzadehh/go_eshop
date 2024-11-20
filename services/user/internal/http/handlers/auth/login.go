package authhttp

import (
	"errors"
	"net/http"

	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/render"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/services/auth"
)

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"

	log := sl.HandlerJob(h.Log, op, requestmdw.GetRequestId(r), h.srvc)

	var req LoginReq
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

	h.JSON(w, http.StatusOK, LoginRes{
		Response: api.Ok(),
		Data: &LoginData{
			Token: access,
		},
	})
}
