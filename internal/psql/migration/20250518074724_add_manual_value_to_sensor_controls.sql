-- +goose Up
ALTER TABLE sensor_controls
  ADD COLUMN manual_bool_value BOOLEAN,
  ADD COLUMN manual_int_value INTEGER;

-- +goose Down
ALTER TABLE sensor_controls
  DROP COLUMN IF EXISTS manual_bool_value,
  DROP COLUMN IF EXISTS manual_int_value;
