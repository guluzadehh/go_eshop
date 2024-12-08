package app

import (
	"log/slog"

	httpapp "github.com/guluzadehh/go_eshop/services/user/app/http"
	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_eshop/services/user/internal/service/auth"
	"github.com/guluzadehh/go_eshop/services/user/internal/service/user"
	"github.com/guluzadehh/go_eshop/services/user/internal/storage/postgresql"
)

type App struct {
	log       *slog.Logger
	HTTPApp   *httpapp.HTTPApp
	pgStorage *postgresql.Storage
}

func New(log *slog.Logger, config *config.Config) *App {
	pgStorage, err := postgresql.New(config.Postgresql.DSN(nil))
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, config, pgStorage, pgStorage)
	userService := user.New(log, config, pgStorage)

	httpApp := httpapp.New(log, config, authService, userService)

	return &App{
		log:       log,
		HTTPApp:   httpApp,
		pgStorage: pgStorage,
	}
}

func (a *App) Start() {
	a.log.Info("running user http app")
	a.HTTPApp.MustRun()
}

func (a *App) Stop() {
	if err := a.HTTPApp.Stop(); err != nil {
		a.log.Error("error while shutdown the HTTP server", sl.Err(err))
	}

	if err := a.pgStorage.Close(); err != nil {
		a.log.Error("error while closing the postgres db connection", sl.Err(err))
	}

	a.log.Info("App has been gracefully stopped")
}
