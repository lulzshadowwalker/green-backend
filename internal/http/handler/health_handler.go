package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	//
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (hh *HealthHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/health", hh.healthCheck)
	e.HEAD("/api/health", hh.healthCheck)
	e.HEAD("/api/hello", hh.healthCheck)
}

func (hh *HealthHandler) healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
