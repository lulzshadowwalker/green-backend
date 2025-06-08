package stores

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	rows, err := sr.q.GetSensorReadingsPastDays(ctx, 7)
	if err != nil {
		return nil, err
	}

	res := make([]internal.SensorReading, len(rows))
	for i, rr := range rows {
		res[i] = sr.toEntity(rr)
	}

	return res, nil
}

// GetSensorReadingsSince returns all sensor readings since the given time.
func (sr *SensorReadings) GetSensorReadingsSince(ctx context.Context, sinceTime time.Time) ([]internal.SensorReading, error) {
	// Use GetSensorReadingsByTime with sinceTime and now
	now := time.Now().UTC()
	limit := int32(1000) // Arbitrary large limit; adjust as needed
	offset := int32(0)
	rows, err := sr.q.GetSensorReadingsByTime(ctx, struct {
		Timestamp   pgtype.Timestamptz
		Timestamp_2 pgtype.Timestamptz
		Limit       int32
		Offset      int32
	}{
		Timestamp:   pgtype.Timestamptz{Time: sinceTime, Valid: true},
		Timestamp_2: pgtype.Timestamptz{Time: now, Valid: true},
		Limit:       limit,
		Offset:      offset,
	})
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
