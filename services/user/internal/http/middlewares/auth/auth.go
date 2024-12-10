package authmdw

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/domain/models"
	requestmdw "github.com/guluzadehh/go_eshop/services/user/internal/http/middlewares/request"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/api"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/jwt"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/render"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/service"
)

type contextKey string

const userContextKey contextKey = "user"

type UserProviderService interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

func Authorize(log *slog.Logger, config *config.Config, userProviderService UserProviderService) mux.MiddlewareFunc {
	render := render.NewResponder(log)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlewares.auth.Authorize"

			log := log.With(
				slog.String("op", op),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", requestmdw.GetRequestId(r)),
			)

			log.Info("authorizing the user")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Info("missing Authorization header")
				render.JSON(w, http.StatusUnauthorized, authFailResponse())
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Info("invalid Authorization header format", slog.String("auth_header", authHeader))
				render.JSON(w, http.StatusUnauthorized, authFailResponse())
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Verify(tokenStr, config)
			if err != nil {
				log.Info("access token is invalid", slog.String("invalid_access_token", tokenStr), sl.Err(err))
				render.JSON(w, http.StatusUnauthorized, authFailResponse())
				return
			}

			// TODO: check if token is blocked

			email, err := token.Claims.GetSubject()
			if err != nil {
				log.Error("error while getting the subject from access token", sl.Err(err))
				render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
				return
			}

			user, err := userProviderService.GetUserByEmail(r.Context(), email)
			if err != nil {
				if errors.Is(err, service.ErrUserNotFound) {
					render.JSON(w, http.StatusUnauthorized, authFailResponse())
					return
				}

				log.Error("failed to get user by user_id from storage", sl.Err(err))
				render.JSON(w, http.StatusInternalServerError, api.UnexpectedError())
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func User(r *http.Request) *models.User {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok || user == nil {
		return nil
	}

	return user
}

func authFailResponse() api.Response {
	return api.Err("you are not authorized")
}
