package httpapp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guluzadehh/go_eshop/services/user/internal/config"
)

type HTTPApp struct {
	log        *slog.Logger
	httpServer *http.Server
}

func New(log *slog.Logger, config *config.Config) *HTTPApp {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", config.HTTPServer.Port),
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
	}

	router := mux.NewRouter()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User service is running"))
	})

	server.Handler = router

	return &HTTPApp{
		log:        log,
		httpServer: &server,
	}
}

func (a *HTTPApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *HTTPApp) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(slog.String("op", op))

	log.Info("starting HTTP server", slog.String("addr", a.httpServer.Addr))
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *HTTPApp) Stop() error {
	const op = "httpapp.Stop"

	log := a.log.With(slog.String("op", op))

	log.Info("stopping HTTP server", slog.String("addr", a.httpServer.Addr))

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
