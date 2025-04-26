package handler

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lulzshadowwalker/green-backend/internal"
)

type SensorReadings struct {
	service SensorReadingsService
}

type SensorReadingsService interface {
	GetSensorReadings(ctx context.Context) ([]internal.SensorReading, error)
	CreateSensorReading(ctx context.Context, params internal.CreateSensorReadingParams) (internal.SensorReading, error)
}

func NewSensorReadings(s SensorReadingsService) *SensorReadings {
	return &SensorReadings{s}
}

func (sr *SensorReadings) RegisterRoutes(a *echo.Echo) {
	a.GET("/api/readings", sr.Index)
	a.POST("/api/readings", sr.Create)
}

func (sr *SensorReadings) Index(c echo.Context) error {
	m, err := sr.service.GetSensorReadings(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sr.collection(m))
}

type CreateSensorReadingRequest struct {
	Data struct {
		Attributes struct {
			SensorType string  `json:"sensorType" validate:"required,oneof=temperature humidity"`
			Value      float64 `json:"value" validate:"required"`
		} `json:"attributes"`
	} `json:"data"`
}

func (sr *SensorReadings) Create(c echo.Context) error {
	var req CreateSensorReadingRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err 
	}

	attr := req.Data.Attributes
	m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
		SensorType: attr.SensorType,
		Value:      attr.Value,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sr.resource(m))
}

func (sr *SensorReadings) resource(r internal.SensorReading) echo.Map {
	return echo.Map{
		"id":   r.ID,
		"type": "sensor-reading",
		"attributes": echo.Map{
			"type":      r.SensorType,
			"value":     r.Value,
			"timestamp": r.Timestamp,
		},
		"relationships": echo.Map{},
		"includes":      echo.Map{},
		"links":         echo.Map{},
	}
}

func (sr *SensorReadings) collection(r []internal.SensorReading) echo.Map {
	res := make([]echo.Map, len(r))
	for i, rr := range r {
		res[i] = sr.resource(rr)
	}

	return echo.Map{
		"data": res,
	}
}
