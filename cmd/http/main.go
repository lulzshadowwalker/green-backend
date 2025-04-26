package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"

	application "github.com/lulzshadowwalker/green-backend/internal/http/app"
	_ "github.com/lulzshadowwalker/green-backend/internal/logging"
	"github.com/lulzshadowwalker/green-backend/internal/psql"
	"github.com/lulzshadowwalker/green-backend/internal/psql/db"
)

func main() {
	//  TODO: Use a config package
	pool, err := psql.Connect(psql.ConnectionParams{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		slog.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}

  q := db.New(pool)
  hello, err := q.GetHello(context.Background())
  if err != nil {
    slog.Error("failed to get hello", "err", err)
    os.Exit(1)
  }

  slog.Info("hello", "value", hello)

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
