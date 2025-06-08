package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lulzshadowwalker/green-backend/internal"
	internalhttp "github.com/lulzshadowwalker/green-backend/internal/http"
)

type Control struct {
	service ControlService
}

type ControlService interface {
	GetAllSensorControls(ctx context.Context) ([]internal.SensorControl, error)
	SetSensorControlMode(ctx context.Context, sensorType, mode string, manualUntil *time.Time) (internal.SensorControl, error)
	SetSensorControlModeWithValue(ctx context.Context, sensorType, mode string, manualUntil *time.Time, manualIntValue *int, manualBoolValue *bool) (internal.SensorControl, error)
}

func NewControlHandler(c ControlService) *Control {
	return &Control{
		service: c,
	}
}

func (c *Control) RegisterRoutes(e *echo.Echo) {
	e.GET("/api/control", c.Index)
	e.POST("/api/control", internalhttp.JWTAuthMiddleware(c.Set))
}

type setControlRequest struct {
	SensorType      string     `json:"sensor_type"`
	Mode            string     `json:"mode"` // "automatic" or "manual"
	ManualUntil     *time.Time `json:"manual_until,omitempty"`
	ManualIntValue  *int       `json:"manual_int_value,omitempty"`
	ManualBoolValue *bool      `json:"manual_bool_value,omitempty"`
}

func (c *Control) Set(ctx echo.Context) error {
	start := time.Now()
	reqID := ctx.Response().Header().Get(echo.HeaderXRequestID)
	slog.Info("Set control request received",
		"method", ctx.Request().Method,
		"path", ctx.Path(),
		"remote_addr", ctx.RealIP(),
		"request_id", reqID,
	)

	var req setControlRequest
	if err := ctx.Bind(&req); err != nil {
		slog.Error("Failed to bind request", "error", err, "request_id", reqID)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
	}

	if req.SensorType == "" || (req.Mode != "automatic" && req.Mode != "manual") {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "sensor_type and valid mode required"})
	}

	// Pass manual values to the service layer (requires service and store updates)
	control, err := c.service.SetSensorControlModeWithValue(
		ctx.Request().Context(),
		req.SensorType,
		req.Mode,
		req.ManualUntil,
		req.ManualIntValue,
		req.ManualBoolValue,
	)
	if err != nil {
		slog.Error("Failed to set sensor control mode", "error", err, "request_id", reqID)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to set control mode"})
	}

	slog.Info("Set sensor control mode",
		"sensor_type", req.SensorType,
		"mode", req.Mode,
		"duration_ms", time.Since(start).Milliseconds(),
		"request_id", reqID,
	)

	return ctx.JSON(http.StatusOK, control)
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

	controls, err := c.service.GetAllSensorControls(ctx.Request().Context())
	if err != nil {
		slog.Error("Failed to get sensor controls", "error", err, "request_id", reqID)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to get control status"})
	}

	// Flat JSON: { "fan_mode": "manual", "fan": 255, ... }
	result := make(map[string]interface{})
	for _, ctrl := range controls {
		modeKey := ctrl.SensorType + "_mode"
		result[modeKey] = ctrl.Mode

		// Only one value per sensor, no _int_value/_bool_value suffix
		switch ctrl.SensorType {
		case "fan", "heat", "light":
			if ctrl.ManualIntValue != nil {
				result[ctrl.SensorType] = *ctrl.ManualIntValue
			} else {
				result[ctrl.SensorType] = nil
			}
		case "door", "pump":
			if ctrl.ManualBoolValue != nil {
				result[ctrl.SensorType] = *ctrl.ManualBoolValue
			} else {
				result[ctrl.SensorType] = nil
			}
		default:
			result[ctrl.SensorType] = nil
		}
	}

	slog.Info("Returning control status",
		"result", result,
		"duration_ms", time.Since(start).Milliseconds(),
		"request_id", reqID,
	)

	return ctx.JSON(http.StatusOK, result)
}
