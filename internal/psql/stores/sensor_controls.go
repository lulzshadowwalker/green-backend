package stores

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lulzshadowwalker/green-backend/internal"
	"github.com/lulzshadowwalker/green-backend/internal/psql/db"
)

type SensorControls struct {
	q *db.Queries
}

func NewSensorControls(q *db.Queries) *SensorControls {
	return &SensorControls{
		q: q,
	}
}

func (sc *SensorControls) toEntity(c db.SensorControl) internal.SensorControl {
	var manualUntil *time.Time
	if c.ManualUntil.Valid {
		t := c.ManualUntil.Time
		manualUntil = &t
	}
	var manualBoolValue *bool
	if c.ManualBoolValue.Valid {
		manualBoolValue = &c.ManualBoolValue.Bool
	}
	var manualIntValue *int
	if c.ManualIntValue.Valid {
		val := int(c.ManualIntValue.Int32)
		manualIntValue = &val
	}
	return internal.SensorControl{
		SensorType:      c.SensorType,
		Mode:            c.Mode,
		ManualUntil:     manualUntil,
		ManualBoolValue: manualBoolValue,
		ManualIntValue:  manualIntValue,
	}
}

func (sc *SensorControls) GetAllSensorControls(ctx context.Context) ([]internal.SensorControl, error) {
	rows, err := sc.q.GetAllSensorControls(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]internal.SensorControl, len(rows))
	for i, row := range rows {
		var manualUntil *time.Time
		if row.ManualUntil.Valid {
			t := row.ManualUntil.Time
			manualUntil = &t
		}
		var manualBoolValue *bool
		if row.ManualBoolValue.Valid {
			manualBoolValue = &row.ManualBoolValue.Bool
		}
		var manualIntValue *int
		if row.ManualIntValue.Valid {
			val := int(row.ManualIntValue.Int32)
			manualIntValue = &val
		}
		res[i] = internal.SensorControl{
			SensorType:      row.SensorType,
			Mode:            row.Mode,
			ManualUntil:     manualUntil,
			ManualBoolValue: manualBoolValue,
			ManualIntValue:  manualIntValue,
		}
	}
	return res, nil
}

func (sc *SensorControls) GetSensorControlByType(ctx context.Context, sensorType string) (internal.SensorControl, error) {
	row, err := sc.q.GetSensorControlByType(ctx, sensorType)
	if err != nil {
		return internal.SensorControl{}, err
	}
	var manualUntil *time.Time
	if row.ManualUntil.Valid {
		t := row.ManualUntil.Time
		manualUntil = &t
	}
	var manualBoolValue *bool
	if row.ManualBoolValue.Valid {
		manualBoolValue = &row.ManualBoolValue.Bool
	}
	var manualIntValue *int
	if row.ManualIntValue.Valid {
		val := int(row.ManualIntValue.Int32)
		manualIntValue = &val
	}
	return internal.SensorControl{
		SensorType:      row.SensorType,
		Mode:            row.Mode,
		ManualUntil:     manualUntil,
		ManualBoolValue: manualBoolValue,
		ManualIntValue:  manualIntValue,
	}, nil
}

func (sc *SensorControls) UpdateSensorControlMode(ctx context.Context, sensorType, mode string, manualUntil *time.Time, manualIntValue *int, manualBoolValue *bool) (internal.SensorControl, error) {
	var mu pgtype.Timestamptz
	if manualUntil != nil {
		mu.Valid = true
		mu.Time = *manualUntil
	} else {
		mu.Valid = false
	}
	var intVal pgtype.Int4
	if manualIntValue != nil {
		intVal.Valid = true
		intVal.Int32 = int32(*manualIntValue)
	} else {
		intVal.Valid = false
	}
	var boolVal pgtype.Bool
	if manualBoolValue != nil {
		boolVal.Valid = true
		boolVal.Bool = *manualBoolValue
	} else {
		boolVal.Valid = false
	}
	row, err := sc.q.UpdateSensorControlMode(ctx, db.UpdateSensorControlModeParams{
		SensorType:      sensorType,
		Mode:            mode,
		ManualUntil:     mu,
		ManualIntValue:  intVal,
		ManualBoolValue: boolVal,
	})
	if err != nil {
		return internal.SensorControl{}, err
	}
	var muUntil *time.Time
	if row.ManualUntil.Valid {
		t := row.ManualUntil.Time
		muUntil = &t
	}
	var manualBoolValueOut *bool
	if row.ManualBoolValue.Valid {
		manualBoolValueOut = &row.ManualBoolValue.Bool
	}
	var manualIntValueOut *int
	if row.ManualIntValue.Valid {
		val := int(row.ManualIntValue.Int32)
		manualIntValueOut = &val
	}
	return internal.SensorControl{
		SensorType:      row.SensorType,
		Mode:            row.Mode,
		ManualUntil:     muUntil,
		ManualBoolValue: manualBoolValueOut,
		ManualIntValue:  manualIntValueOut,
	}, nil
}

func (sc *SensorControls) InsertOrUpdateSensorControl(ctx context.Context, sensorType, mode string, manualUntil *time.Time, manualIntValue *int, manualBoolValue *bool) (internal.SensorControl, error) {
	var mu pgtype.Timestamptz
	if manualUntil != nil {
		mu.Valid = true
		mu.Time = *manualUntil
	} else {
		mu.Valid = false
	}
	var intVal pgtype.Int4
	if manualIntValue != nil {
		intVal.Valid = true
		intVal.Int32 = int32(*manualIntValue)
	} else {
		intVal.Valid = false
	}
	var boolVal pgtype.Bool
	if manualBoolValue != nil {
		boolVal.Valid = true
		boolVal.Bool = *manualBoolValue
	} else {
		boolVal.Valid = false
	}
	row, err := sc.q.InsertSensorControl(ctx, db.InsertSensorControlParams{
		SensorType:      sensorType,
		Mode:            mode,
		ManualUntil:     mu,
		ManualIntValue:  intVal,
		ManualBoolValue: boolVal,
	})
	if err != nil {
		return internal.SensorControl{}, err
	}
	var muUntil *time.Time
	if row.ManualUntil.Valid {
		t := row.ManualUntil.Time
		muUntil = &t
	}
	var manualBoolValueOut *bool
	if row.ManualBoolValue.Valid {
		manualBoolValueOut = &row.ManualBoolValue.Bool
	}
	var manualIntValueOut *int
	if row.ManualIntValue.Valid {
		val := int(row.ManualIntValue.Int32)
		manualIntValueOut = &val
	}
	return internal.SensorControl{
		SensorType:      row.SensorType,
		Mode:            row.Mode,
		ManualUntil:     muUntil,
		ManualBoolValue: manualBoolValueOut,
		ManualIntValue:  manualIntValueOut,
	}, nil
}
