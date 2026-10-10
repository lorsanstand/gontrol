package config

import (
	"log/slog"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     int    `env:"PORT" envDefault:"5467"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"WARN"`
}

func Load() (Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return Config{}, err
	}

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) Level() slog.Level {
	switch c.LogLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}
