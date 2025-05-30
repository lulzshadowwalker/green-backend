-- name: GetSensorReadings :many
SELECT * from sensor_readings;

-- name: GetSensorReading :one
SELECT * from sensor_readings
WHERE id = $1;

-- name: GetSensorReadingsByType :many
SELECT * from sensor_readings
WHERE sensor_type = $1
ORDER BY timestamp DESC
LIMIT $2 OFFSET $3;

-- name: GetSensorReadingsByTypeAndTime :many
SELECT * from sensor_readings
WHERE sensor_type = $1
  AND timestamp >= $2
  AND timestamp <= $3
ORDER BY timestamp DESC
LIMIT $4 OFFSET $5;

-- name: GetSensorReadingsPastDays :many
SELECT * from sensor_readings
WHERE timestamp >= NOW() - INTERVAL '1 day' * $1
ORDER BY timestamp DESC;

-- name: GetSensorReadingsByTime :many
SELECT * from sensor_readings
WHERE timestamp >= $1
  AND timestamp <= $2
ORDER BY timestamp DESC
LIMIT $3 OFFSET $4;

-- name: CreateSensorReading :one
INSERT INTO sensor_readings (sensor_type, value)
VALUES ($1, $2)
RETURNING *;
