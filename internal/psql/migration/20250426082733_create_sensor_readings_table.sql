-- +goose Up
-- +goose StatementBegin
CREATE TABLE sensor_readings (
  id          BIGSERIAL     PRIMARY KEY,
  sensor_type TEXT          NOT NULL,                     
  value       DOUBLE PRECISION NOT NULL,
  timestamp   TIMESTAMPTZ   NOT NULL  DEFAULT now()
);

CREATE INDEX idx_readings_dev_type_time
  ON sensor_readings (sensor_type, timestamp DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_readings_dev_type_time;
DROP TABLE IF EXISTS sensor_readings;
-- +goose StatementEnd
