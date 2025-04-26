package stores

import (
	"context"
	"strconv"

	"github.com/lulzshadowwalker/green-backend/internal"
	"github.com/lulzshadowwalker/green-backend/internal/psql/db"
)

type SensorReadings struct {
	q *db.Queries
}

func NewSensorReadings(q *db.Queries) *SensorReadings {
	return &SensorReadings{
		q: q,
	}
}

func (sr *SensorReadings) toEntity(r db.SensorReading) internal.SensorReading {
	return internal.SensorReading{
		ID:         strconv.Itoa(int(r.ID)),
		SensorType: r.SensorType,
		Value:      r.Value,
		Timestamp:  r.Timestamp.Time,
	}
}

func (sr *SensorReadings) GetSensorReadings(ctx context.Context) ([]internal.SensorReading, error) {
	rows, err := sr.q.GetSensorReadings(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]internal.SensorReading, len(rows))
	for i, rr := range rows {
		res[i] = sr.toEntity(rr)
	}

	return res, nil
}

func (sr *SensorReadings) CreateSensorReading(ctx context.Context, params internal.CreateSensorReadingParams) (internal.SensorReading, error) {
	arg := db.CreateSensorReadingParams{
		SensorType: params.SensorType,
		Value:      params.Value,
	}

	row, err := sr.q.CreateSensorReading(ctx, arg)
	if err != nil {
		return internal.SensorReading{}, err
	}

	return sr.toEntity(row), nil
}
