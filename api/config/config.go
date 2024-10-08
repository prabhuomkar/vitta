package config

import "time"

// Config for the API
type Config struct {
	LogLevel        string        `env:"LOG_LEVEL"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT"`
	DatabaseURL     string        `env:"DATABASE_URL"`
	DatabaseTimeout time.Duration `env:"DATABASE_TIMEOUT"`
}
