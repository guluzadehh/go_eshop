package app

import (
	"log/slog"

	httpapp "github.com/guluzadehh/go_eshop/services/user/app/http"
	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/guluzadehh/go_eshop/services/user/internal/lib/sl"
)

type App struct {
	log     *slog.Logger
	HTTPApp *httpapp.HTTPApp
}

func New(log *slog.Logger, config *config.Config) *App {
	httpApp := httpapp.New(log, config)

	return &App{
		log:     log,
		HTTPApp: httpApp,
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

	a.log.Info("App has been gracefully stopped")
}
