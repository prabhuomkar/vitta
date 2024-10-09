package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config for the API.
type Config struct {
	Port            uint16        `default:"8080"                                            env:"PORT"`
	LogLevel        slog.Level    `default:"INFO"                                            env:"LOG_LEVEL"`
	ShutdownTimeout time.Duration `default:"5s"                                              env:"SHUTDOWN_TIMEOUT"`
	ReadTimeout     time.Duration `default:"5s"                                              env:"READ_TIMEOUT"`
	WriteTimeout    time.Duration `default:"5s"                                              env:"WRITE_TIMEOUT"`
	DatabaseURL     string        `default:"postgres://vdbuser:vdbpass@localhost:5432/vitta" env:"DATABASE_URL"`
	DatabaseTimeout time.Duration `default:"5s"                                              env:"DATABASE_TIMEOUT"`
	AdminUsername   string        `default:"vitta"                                           env:"ADMIN_USERNAME"`
	AdminPassword   string        `default:"vittaT3st!"                                      env:"ADMIN_PASSWORD"`
}

func New() *Config {
	var cfg Config

	err := envconfig.Process("vitta", &cfg)
	if err != nil {
		slog.Error("error reading config", "error", err)
		os.Exit(1)
	}

	return &cfg
}
