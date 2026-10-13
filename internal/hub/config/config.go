package config

import (
	"log/slog"

	"github.com/lorsanstand/gontrol/internal/utils/config"
)

type Config struct {
	Port     int        `env:"PORT" envDefault:"5467"`
	LogLevel slog.Level `env:"LOG_LEVEL" envDefault:"WARN"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := config.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
