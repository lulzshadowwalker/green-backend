package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	application "github.com/lulzshadowwalker/green-backend/internal/http/app"
	_ "github.com/lulzshadowwalker/green-backend/internal/logging"
	"github.com/lulzshadowwalker/green-backend/internal/psql"
)

func main() {
	//  TODO: Use a config package
	pool, err := psql.Connect(psql.ConnectionParams{
		// 	Host:     os.Getenv("DB_HOST"),
		// 	Port:     os.Getenv("DB_PORT"),
		// 	Username: os.Getenv("DB_USERNAME"),
		// 	Password: os.Getenv("DB_PASSWORD"),
		// 	Name:     os.Getenv("DB_NAME"),
		// 	SSLMode:  os.Getenv("DB_SSLMODE"),

		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "example",
		Name:     "mydb",
		SSLMode:  "disable",
	})
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

	app, err := application.New(application.WithDB(pool))
	if err != nil {
		slog.Error("app creation failed", "err", err)
		return
	}

	slog.Info("server started", "addr", app.Addr(), "timeout", app.Timeout())

	//  TODO: Graceful termination
	defer app.Close()
	if err := app.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server shutdown", "err", err)
		} else {
			slog.Info("server crashed", "err", err)
		}
	}
}
