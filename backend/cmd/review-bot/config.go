package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string

	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func NewConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process("reviewbot", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
