package handler

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lulzshadowwalker/green-backend/internal/service"
)

type LLMHandler struct {
	service service.LLMService
}

func NewLLMHandler(s service.LLMService) *LLMHandler {
	return &LLMHandler{service: s}
}

func (h *LLMHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/llm/plant-advice", h.StreamPlantAdvice)
}

func (h *LLMHandler) StreamPlantAdvice(c echo.Context) error {
	plant := c.QueryParam("plant")
	if plant == "" {
		plant = "strawberry"
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/plain")
	c.Response().WriteHeader(http.StatusOK)

	pr, pw := io.Pipe()
	
	ctx, cancel := context.WithTimeout(c.Request().Context(), 60*time.Second)
	defer cancel()

	go func() {
		defer pw.Close()
		err := h.service.StreamPlantAdvice(ctx, plant, pw)
		if err != nil {
			pw.CloseWithError(err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := pr.Read(buf)
		if n > 0 {
			_, writeErr := c.Response().Write(buf[:n])
			if writeErr != nil {
				return writeErr
			}
			c.Response().Flush()
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}