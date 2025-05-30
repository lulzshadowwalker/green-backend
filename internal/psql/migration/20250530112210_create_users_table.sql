-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    password_hash VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
);

INSERT INTO
    users (username, password_hash)
VALUES
    (
        'user',
        '$2a$10$pNY950b/7F7kdtuIKyQuv.5K21zbd/vfY9NPGVHxBnfBQyqioTfl.' -- "password"
    ) ON CONFLICT (username) DO NOTHING;

-- +goose Down
DROP TABLE IF EXISTS users;
