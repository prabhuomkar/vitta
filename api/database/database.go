package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config for the database connection.
type Config struct {
	URL     string
	Timeout time.Duration
}

type DBIface interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Close()
}

func New(cfg *Config) (DBIface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	conn, err := pgxpool.New(ctx, cfg.URL)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return conn, nil
}
