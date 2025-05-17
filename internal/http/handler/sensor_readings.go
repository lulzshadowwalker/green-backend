package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

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
	a.GET("/api/loli", sr.Index)
	a.POST("/api/readings", sr.Create)
}

func (sr *SensorReadings) Index(c echo.Context) error {
	start := time.Now()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	slog.Info("SensorReadings Index request received",
		"method", c.Request().Method,
		"path", c.Path(),
		"remote_addr", c.RealIP(),
		"request_id", reqID,
	)

	m, err := sr.service.GetSensorReadings(c.Request().Context())
	if err != nil {
		slog.Error("Failed to get sensor readings",
			"error", err,
			"request_id", reqID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	slog.Info("Returning sensor readings",
		"count", len(m),
		"request_id", reqID,
		"duration_ms", time.Since(start).Milliseconds(),
	)

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
	start := time.Now()
	reqID := c.Response().Header().Get(echo.HeaderXRequestID)
	slog.Info("SensorReadings Create request received",
		"method", c.Request().Method,
		"path", c.Path(),
		"remote_addr", c.RealIP(),
		"request_id", reqID,
	)

	var req CreateSensorReadingRequest
	if err := c.Bind(&req); err != nil {
		slog.Error("Failed to bind request body",
			"error", err,
			"request_id", reqID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	if err := c.Validate(req); err != nil {
		slog.Error("Validation failed for request body",
			"error", err,
			"request_id", reqID,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return err
	}

	//  TODO: Move this into the service class with a transaction
	readings := make([]internal.SensorReading, 0)
	createdTypes := []string{}

	if req.Temperature != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "temperature",
			Value:      req.Temperature,
		})
		if err != nil {
			slog.Error("Failed to create temperature reading",
				"error", err,
				"request_id", reqID,
			)
			return err
		}
		readings = append(readings, m)
		createdTypes = append(createdTypes, "temperature")
	}

	if req.Humidity != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "humidity",
			Value:      req.Humidity,
		})
		if err != nil {
			slog.Error("Failed to create humidity reading",
				"error", err,
				"request_id", reqID,
			)
			return err
		}
		readings = append(readings, m)
		createdTypes = append(createdTypes, "humidity")
	}

	if req.LightLevel != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "light",
			Value:      req.LightLevel,
		})
		if err != nil {
			slog.Error("Failed to create light reading",
				"error", err,
				"request_id", reqID,
			)
			return err
		}
		readings = append(readings, m)
		createdTypes = append(createdTypes, "light")
	}

	if req.WaterLevel != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "water",
			Value:      req.WaterLevel,
		})
		if err != nil {
			slog.Error("Failed to create water reading",
				"error", err,
				"request_id", reqID,
			)
			return err
		}
		readings = append(readings, m)
		createdTypes = append(createdTypes, "water")
	}

	if req.SoilMoisture != 0 {
		m, err := sr.service.CreateSensorReading(c.Request().Context(), internal.CreateSensorReadingParams{
			SensorType: "soil",
			Value:      req.SoilMoisture,
		})
		if err != nil {
			slog.Error("Failed to create soil reading",
				"error", err,
				"request_id", reqID,
			)
			return err
		}
		readings = append(readings, m)
		createdTypes = append(createdTypes, "soil")
	}

	slog.Info("Created sensor readings",
		"types", createdTypes,
		"count", len(readings),
		"request_id", reqID,
		"duration_ms", time.Since(start).Milliseconds(),
	)

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
