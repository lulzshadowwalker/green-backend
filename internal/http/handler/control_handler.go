package handler

import (
	"math/rand"
	"net/http"

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
	automatic = !automatic
	return ctx.JSON(http.StatusOK, echo.Map{
		"automatic": automatic,
	})
}

func (c *Control) Index(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, echo.Map{
		"automatic": automatic,
		"fan":       randomInt(),
		"pump":      randomBool(),
		"door":      randomBool(),
		"heat":      randomInt(),
		"light":     randomInt(),
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
