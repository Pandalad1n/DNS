package main

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ListenWebAddress string   `envconfig:"LISTEN_WEB_ADDRESS"`
	SectorID         *float64 `envconfig:"SECTOR_ID"`
	LogLevel         string   `envconfig:"LOG_LVL"`
}

func (c Config) WithEnv() Config {
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to process env.")
	}
	return c
}

func (c Config) Validate() error {
	if c.ListenWebAddress == "" {
		return errors.New("address not set")
	}
	if c.SectorID == nil {
		return errors.New("sectorID not set")
	}
	if c.LogLevel == "" {
		return errors.New("LogLevel not set")
	}
	return nil
}
