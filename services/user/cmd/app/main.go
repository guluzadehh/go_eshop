package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/guluzadehh/go_eshop/services/user/app"
	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/joho/godotenv"
)

const (
	env_local = "local"
	env_dev   = "dev"
	env_prod  = "prod"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s\n", err.Error())
	}

	config := config.MustLoad()

	log := setupLogger(config.Env)
	log.Info("starting user app", slog.String("env", config.Env))

	app := app.New(log, config)

	go app.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	app.Stop()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case env_local, env_dev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case env_prod:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
