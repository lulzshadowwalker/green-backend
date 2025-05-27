package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Threshold struct {
	//
}

type ThresholdValues struct {
	LightMin int `json:"light_min"`
	Lightmax int `json:"light_max"`
	SoilMin  int `json:"soil_min"`
	SoilMax  int `json:"soil_max"`
	TempMin  int `json:"temp_min"`
	TempMax  int `json:"temp_max"`
}

func NewThresholdHandler() *Threshold {
	return &Threshold{}
}

func (t *Threshold) RegisterRoutes(a *echo.Echo) {
	a.GET("/api/thresholds", t.Index)
}

func (t *Threshold) Index(c echo.Context) error {
	start := time.Now()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	slog.Info("Threshold Index request received",
		"method", c.Request().Method,
		"path", c.Path(),
		"remote_addr", c.RealIP(),
		"request_id", reqID,
	)

	sample := ThresholdValues{
		LightMin: 100,
		Lightmax: 800,
		SoilMin:  200,
		SoilMax:  600,
		TempMin:  15,
		TempMax:  30,
	}

	slog.Info("Returning threshold values",
		"thresholds", sample,
		"request_id", reqID,
		"duration_ms", time.Since(start).Milliseconds(),
	)

	return c.JSON(http.StatusOK, t.resource(sample))
}

func (t *Threshold) resource(r ThresholdValues) echo.Map {
	return echo.Map{
		"light_min": r.LightMin,
		"light_max": r.Lightmax,
		"soil_min":  r.SoilMin,
		"soil_max":  r.SoilMax,
		"temp_min":  r.TempMin,
		"temp_max":  r.TempMax,
	}
}

func (t *Threshold) collection(r []ThresholdValues) echo.Map {
	res := make([]echo.Map, len(r))
	for i, rr := range r {
		res[i] = t.resource(rr)
	}

	return echo.Map{
		"data": res,
	}
}
