// Package config contains configuration structs.
package config

import (
	"errors"
	"fmt"
	"os"
)

var errUndefinedEnvVar = errors.New("undefined environment variable")

// Config hold the service config.
type Config struct {
	Database    *Database
	LoggerLevel string
}

// New is a constructor function for Config.
func New() (*Config, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}

	loggerLevel, defined := os.LookupEnv("LOGGER_LEVEL")
	if !defined {
		return nil, fmt.Errorf("%w: LOGGER_LEVEL", errUndefinedEnvVar)
	}

	return &Config{
		Database:    db,
		LoggerLevel: loggerLevel,
	}, nil
}
