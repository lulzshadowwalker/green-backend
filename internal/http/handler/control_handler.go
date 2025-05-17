package handler

import (
	"log/slog"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Control struct {
	service ControlService
}

type ControlService interface {
	//
}

func NewControlHandler(c ControlService) *Control {
	return &Control{
		service: c,
	}
}

func (c *Control) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/control", c.Index)
	e.POST("/api/control", c.Toggle)
}

var automatic bool = false

func (c *Control) Toggle(ctx echo.Context) error {
	start := time.Now()
	reqID := ctx.Response().Header().Get(echo.HeaderXRequestID)
	slog.Info("Toggle request received",
		"method", ctx.Request().Method,
		"path", ctx.Path(),
		"remote_addr", ctx.RealIP(),
		"request_id", reqID,
	)

	automatic = !automatic

	slog.Info("Toggled automatic mode",
		"automatic", automatic,
		"duration_ms", time.Since(start).Milliseconds(),
		"request_id", reqID,
	)

	return ctx.JSON(http.StatusOK, echo.Map{
		"automatic": automatic,
	})
}

func (c *Control) Index(ctx echo.Context) error {
	start := time.Now()
	reqID := ctx.Response().Header().Get(echo.HeaderXRequestID)
	slog.Info("Index request received",
		"method", ctx.Request().Method,
		"path", ctx.Path(),
		"remote_addr", ctx.RealIP(),
		"request_id", reqID,
	)

	fan := randomInt()
	pump := randomBool()
	door := randomBool()
	heat := randomInt()
	light := randomInt()

	slog.Info("Returning control status",
		"automatic", automatic,
		"fan", fan,
		"pump", pump,
		"door", door,
		"heat", heat,
		"light", light,
		"duration_ms", time.Since(start).Milliseconds(),
		"request_id", reqID,
	)

	return ctx.JSON(http.StatusOK, echo.Map{
		"automatic": automatic,
		"fan":       fan,
		"pump":      pump,
		"door":      door,
		"heat":      heat,
		"light":     light,
	})
}

func randomInt() int {
	if rand.Intn(2) == 0 {
		return 0
	}

	return 255
}

func randomBool() bool {
	return rand.Int31n(2) == 1
}
