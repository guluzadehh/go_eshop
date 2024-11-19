package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/guluzadehh/go_eshop/services/user/internal/config"
	"github.com/joho/godotenv"
)

type logger struct {
	sl      *slog.Logger
	verbose bool
}

func newLogger(v bool) *logger {
	return &logger{
		sl: slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}),
		),
		verbose: v,
	}
}

func (l *logger) Printf(format string, v ...interface{}) {
	l.sl.Info(fmt.Sprintf(format, v))
}

func (l *logger) Verbose() bool {
	return true
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s\n", err.Error())
	}

	var action string

	flag.StringVar(&action, "action", "up", "migration action")
	flag.Parse()

	config := config.MustLoad()

	dbOpts := config.Postgresql.Options
	dbOpts = append(dbOpts, fmt.Sprintf("x-migrations-table=%s", config.Migrations.TableName))

	m, err := migrate.New(
		"file://"+config.Migrations.Path,
		config.Postgresql.DSN(dbOpts),
	)
	if err != nil {
		panic(err)
	}

	m.Log = newLogger(true)

	switch action {
	case "up":
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no migrations to apply")
				return
			}

			panic(err)
		}

		fmt.Println("migrations have been applied successfully")

	case "down":
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				fmt.Println("no changes")
			}

			panic(err)
		}

		fmt.Println("migrations are down")
	}

}
