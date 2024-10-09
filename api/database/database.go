package database

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// Config for the database connection.
type Config struct {
	URL     string
	Timeout time.Duration
}

func New(cfg *Config) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, cfg.URL)
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		os.Exit(1) //nolint: gocritic
	}

	return conn
}
