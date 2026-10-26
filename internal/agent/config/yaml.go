package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

var ErrCreateFile = errors.New("create new config.yaml")

type Config struct {
	HubURL   *url.URL   `yaml:"hub-url"`
	LogLevel slog.Level `yaml:"log-level"`
	AgentID  *uuid.UUID `yaml:"agent-id"`
}

// LoadConfig читает конфиг с с файла config.yaml.
// Если файла нет он будет создан и будет возращена ошибка ErrCreateFile.
func Load() (*Config, error) {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err = createFile(); err != nil {
				return nil, fmt.Errorf("create new file config: %w", err)
			}

			return nil, ErrCreateFile
		}

		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save сохраняет новый конфиг в файл config.yaml.
func (c *Config) Save() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = os.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func createFile() error {
	cfg := Config{HubURL: nil, LogLevel: slog.LevelWarn, AgentID: nil}

	return cfg.Save()
}
