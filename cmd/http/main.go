package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	application "github.com/lulzshadowwalker/green-backend/internal/http/app"
	_ "github.com/lulzshadowwalker/green-backend/internal/logging"
	"github.com/lulzshadowwalker/green-backend/internal/psql"
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

	app, err := application.New(application.WithDB(pool))
	if err != nil {
		slog.Error("app creation failed", "err", err)
		return
	}

	slog.Info("server started", "addr", app.Addr(), "timeout", app.Timeout())

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- app.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("shutdown signal received", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Echo.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "err", err)
	}
	app.Close()

	err = <-serverErr
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server shutdown", "err", err)
		return
	}
}
