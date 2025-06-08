package service

import (
	"context"
	"time"

	"github.com/lulzshadowwalker/green-backend/internal"
)

type SensorReadings struct {
	r SensorReadingsStore
}

type SensorReadingsStore interface {
	GetSensorReadings(ctx context.Context) ([]internal.SensorReading, error)
	CreateSensorReading(ctx context.Context, params internal.CreateSensorReadingParams) (internal.SensorReading, error)
	GetSensorReadingsSince(ctx context.Context, since time.Time) ([]internal.SensorReading, error)
}

func NewSensorReadings(r SensorReadingsStore) *SensorReadings {
	return &SensorReadings{
		r: r,
	}
}

func (s SensorReadings) GetSensorReadings(ctx context.Context) ([]internal.SensorReading, error) {
	return s.r.GetSensorReadings(ctx)
}

func (s SensorReadings) CreateSensorReading(ctx context.Context, params internal.CreateSensorReadingParams) (internal.SensorReading, error) {
	m, err := s.r.CreateSensorReading(ctx, params)
	if err != nil {
		return internal.SensorReading{}, err
	}

	return m, nil
}
