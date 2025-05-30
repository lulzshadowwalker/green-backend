-- name: GetAllSensorControls :many
SELECT sensor_type, mode, manual_until, manual_bool_value, manual_int_value FROM sensor_controls;

-- name: GetSensorControlByType :one
SELECT sensor_type, mode, manual_until, manual_bool_value, manual_int_value FROM sensor_controls
WHERE sensor_type = $1;

-- name: UpdateSensorControlMode :one
UPDATE sensor_controls
SET mode = $2,
    manual_until = $3,
    manual_bool_value = $4,
    manual_int_value = $5,
    updated_at = NOW()
WHERE sensor_type = $1
RETURNING sensor_type, mode, manual_until, manual_bool_value, manual_int_value;

-- name: InsertSensorControl :one
INSERT INTO sensor_controls (sensor_type, mode, manual_until, manual_bool_value, manual_int_value)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (sensor_type) DO UPDATE
    SET mode = EXCLUDED.mode,
        manual_until = EXCLUDED.manual_until,
        manual_bool_value = EXCLUDED.manual_bool_value,
        manual_int_value = EXCLUDED.manual_int_value,
        updated_at = NOW()
RETURNING sensor_type, mode, manual_until, manual_bool_value, manual_int_value;
