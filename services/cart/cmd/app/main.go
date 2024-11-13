package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/guluzadehh/go_eshop/services/cart/internal/config"
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
	log.Info("starting cart app", slog.String("env", config.Env))

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
