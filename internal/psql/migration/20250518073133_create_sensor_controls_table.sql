-- +goose Up
CREATE TABLE IF NOT EXISTS sensor_controls (
    id SERIAL PRIMARY KEY,
    sensor_type VARCHAR(32) NOT NULL UNIQUE,
    mode VARCHAR(16) NOT NULL DEFAULT 'automatic', -- 'automatic' or 'manual'
    manual_until TIMESTAMPTZ, -- nullable, for future timed manual mode
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

-- Insert default rows for supported sensors
INSERT INTO
    sensor_controls (sensor_type, mode)
VALUES
    ('temperature', 'automatic'),
    ('humidity', 'automatic'),
    ('lightLevel', 'automatic'),
    ('waterLevel', 'automatic'),
    ('soilMoisture', 'automatic') ON CONFLICT (sensor_type) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS sensor_controls;
