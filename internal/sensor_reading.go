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

type SensorControl struct {
	SensorType      string     `json:"sensor_type"`
	Mode            string     `json:"mode"`         // "automatic" or "manual"
	ManualUntil     *time.Time `json:"manual_until,omitempty"` // optional, for future timed manual mode
	ManualBoolValue *bool      `json:"manual_bool_value,omitempty"` // optional, for boolean manual control
	ManualIntValue  *int       `json:"manual_int_value,omitempty"`  // optional, for int manual control
}
