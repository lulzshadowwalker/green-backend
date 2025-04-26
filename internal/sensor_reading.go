package internal

import "time"

type SensorReading struct {
	ID         string
	SensorType string
	Value      float64
	Timestamp  time.Time
}

type CreateSensorReadingParams struct {
	SensorType string
	Value      float64
}
