package config

import (
	"log/slog"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config for the API.
type Config struct {
	Port               uint16        `default:"8080"                                            env:"PORT"`
	LogLevel           slog.Level    `default:"DEBUG"                                           env:"LOG_LEVEL"`
	UploadMemoryLimit  int64         `default:"67108864"                                        env:"UPLOAD_MEMORY_LIMIT"`
	ShutdownTimeout    time.Duration `default:"5s"                                              env:"SHUTDOWN_TIMEOUT"`
	ReadTimeout        time.Duration `default:"5s"                                              env:"READ_TIMEOUT"`
	WriteTimeout       time.Duration `default:"5s"                                              env:"WRITE_TIMEOUT"`
	DatabaseURL        string        `default:"postgres://vdbuser:vdbpass@localhost:5432/vitta" env:"DATABASE_URL"`
	DatabaseTimeout    time.Duration `default:"5s"                                              env:"DATABASE_TIMEOUT"`
	AdminUsername      string        `default:"vitta"                                           env:"ADMIN_USERNAME"`
	AdminPassword      string        `default:"vittaT3st!"                                      env:"ADMIN_PASSWORD"`
	AdaptersConfigPath string        `default:"adapters.csv"                                    env:"ADAPTERS_PATH"`
}

func New() (*Config, error) {
	var cfg Config

	err := envconfig.Process("vitta", &cfg)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return &cfg, nil
}
