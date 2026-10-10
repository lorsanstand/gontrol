package config

import (
	"errors"
	"io/fs"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load(cfg interface{}) error {
	if err := godotenv.Load(".env"); err != nil && errors.Is(err, fs.ErrNotExist) {
		return err
	}

	err := env.Parse(cfg)
	if err != nil {
		return err
	}

	return nil
}
