package psql

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionParams struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

func Connect(params ConnectionParams) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(params.Username, params.Password),
		Host:   fmt.Sprintf("%s:%s", params.Host, params.Port),
		Path:   params.Name,
	}

	q := dsn.Query()
	q.Add("sslmode", params.SSLMode)

	dsn.RawQuery = q.Encode()

	pool, err := pgxpool.New(context.Background(), dsn.String())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database because %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database because %w", err)
	}

	slog.Info("connected to database", "host", params.Host, "port", params.Port, "name", params.Name)

	return pool, nil
}
