package config

import (
	"fmt"
	"log/slog"

	"github.com/lorsanstand/gontrol/internal/utils/config"
)

type Config struct {
	Port       int        `env:"PORT" envDefault:"5467"`
	LogLevel   slog.Level `env:"LOG_LEVEL" envDefault:"WARN"`
	DBPort     int        `env:"DB_PORT,required"`
	DBHost     string     `env:"DB_HOST,required"`
	DBPassword string     `env:"DB_PASS,required"`
	DBUser     string     `env:"DB_USER,required"`
	DBName     string     `env:"DB_NAME,required"`
	DBModeSSL  string     `env:"DB_SSL_MODE" envDefault:"disable"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := config.Load(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBModeSSL,
	)
}
