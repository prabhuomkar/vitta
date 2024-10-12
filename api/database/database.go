package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Config for the database connection.
type Config struct {
	URL     string
	Timeout time.Duration
}

type DBIface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

func New(cfg *Config) (DBIface, error) { //nolint: ireturn
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, cfg.URL)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return conn, nil
}
