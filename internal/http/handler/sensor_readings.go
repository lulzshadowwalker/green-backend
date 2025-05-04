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
	Temperature  float64 `json:"temperature" validate:"number"`
	Humidity     float64 `json:"humidity" validate:"number"`
	LightLevel   float64 `json:"lightLevel" validate:"number"`
	WaterLevel   float64 `json:"waterLevel" validate:"number"`
	SoilMoisture float64 `json:"soilMoisture" validate:"number"`
}

func (sr *SensorReadings) Create(c echo.Context) error {
	var req CreateSensorReadingRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	//  TODO: Move this into the service class with a transaction
	readings := make([]internal.SensorReading, 0)

	if req.Temperature != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "temperature",
			Value:      req.Temperature,
		})
		if err != nil {
			return err
		}

		readings = append(readings, m)
	}

	if req.Humidity != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "humidity",
			Value:      req.Humidity,
		})
		if err != nil {
			return err
		}

		readings = append(readings, m)
	}

	if req.LightLevel != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "light",
			Value:      req.LightLevel,
		})
		if err != nil {
			return err
		}

		readings = append(readings, m)
	}

	if req.WaterLevel != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "water",
			Value:      req.WaterLevel,
		})
		if err != nil {
			return err
		}

		readings = append(readings, m)
	}

	if req.SoilMoisture != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "soil",
			Value:      req.SoilMoisture,
		})
		if err != nil {
			return err
		}

		readings = append(readings, m)
	}

	return c.JSON(http.StatusOK, sr.collection(readings))
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
