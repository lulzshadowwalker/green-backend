package main

import (
	"errors"
	"log/slog"
	"net/http"
	application "github.com/lulzshadowwalker/green-backend/internal/http/app"
)

func main() {
  //  TODO: Use a config package 
  //  FIXME: Instead of using a .env file, we need to supply env variables to the docker container
	// if err := godotenv.Load(".env.local"); err != nil {
	// 	slog.Error("failed to load .env.local", "err", err)
	// 	os.Exit(1)
	// }

	app, err := application.New()
	if err != nil {
		slog.Error("app creation failed", "err", err)
		return
	}

	app.Echo.Logger.Info("server started", "addr", app.Addr(), "timeout", app.Timeout())
	if err := app.Start(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server shutdown", "err", err)
		} else {
			slog.Info("server crashed", "err", err)
		}
	}
	//  TODO: Graceful termination
	defer app.Close()
}
