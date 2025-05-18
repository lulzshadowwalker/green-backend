package service

import (
	"context"
	"time"

	"github.com/lulzshadowwalker/green-backend/internal"
)

// SensorControlsService defines the interface for managing sensor control modes.
type SensorControlsService interface {
	GetAllSensorControls(ctx context.Context) ([]internal.SensorControl, error)
	GetSensorControlByType(ctx context.Context, sensorType string) (internal.SensorControl, error)
	SetSensorControlMode(ctx context.Context, sensorType, mode string, manualUntil *time.Time) (internal.SensorControl, error)
	SetSensorControlModeWithValue(ctx context.Context, sensorType, mode string, manualUntil *time.Time, manualIntValue *int, manualBoolValue *bool) (internal.SensorControl, error)
}

// SensorControlsStore defines the data access interface for sensor controls.
type SensorControlsStore interface {
	GetAllSensorControls(ctx context.Context) ([]internal.SensorControl, error)
	GetSensorControlByType(ctx context.Context, sensorType string) (internal.SensorControl, error)
	UpdateSensorControlMode(ctx context.Context, sensorType, mode string, manualUntil *time.Time, manualIntValue *int, manualBoolValue *bool) (internal.SensorControl, error)
	InsertOrUpdateSensorControl(ctx context.Context, sensorType, mode string, manualUntil *time.Time, manualIntValue *int, manualBoolValue *bool) (internal.SensorControl, error)
}

// sensorControlsService is the concrete implementation of SensorControlsService.
type sensorControlsService struct {
	store SensorControlsStore
}

// NewSensorControlsService creates a new SensorControlsService.
func NewSensorControlsService(store SensorControlsStore) SensorControlsService {
	return &sensorControlsService{store: store}
}

func (s *sensorControlsService) GetAllSensorControls(ctx context.Context) ([]internal.SensorControl, error) {
	return s.store.GetAllSensorControls(ctx)
}

func (s *sensorControlsService) GetSensorControlByType(ctx context.Context, sensorType string) (internal.SensorControl, error) {
	return s.store.GetSensorControlByType(ctx, sensorType)
}

func (s *sensorControlsService) SetSensorControlMode(ctx context.Context, sensorType, mode string, manualUntil *time.Time) (internal.SensorControl, error) {
	// Use InsertOrUpdate to ensure the row exists for the sensor type.
	return s.store.InsertOrUpdateSensorControl(ctx, sensorType, mode, manualUntil, nil, nil)
}

func (s *sensorControlsService) SetSensorControlModeWithValue(ctx context.Context, sensorType, mode string, manualUntil *time.Time, manualIntValue *int, manualBoolValue *bool) (internal.SensorControl, error) {
	return s.store.InsertOrUpdateSensorControl(ctx, sensorType, mode, manualUntil, manualIntValue, manualBoolValue)
}
