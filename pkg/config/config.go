// Package config contains configuration structs.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

var _errUndefinedEnvVar = errors.New("undefined environment variable")

// Config hold the service config.
type Config struct {
	Database            *Database
	ServerPort          string
	LoggerLevel         string
	NearThresholdRadius int
}

// New is a constructor function for Config.
func New() (*Config, error) {
	db, err := NewDatabase()
	if err != nil {
		return nil, err
	}

	serverPort, defined := os.LookupEnv("SERVER_PORT")
	if !defined {
		return nil, fmt.Errorf("%w: SERVER_PORT", _errUndefinedEnvVar)
	}

	loggerLevel, defined := os.LookupEnv("LOGGER_LEVEL")
	if !defined {
		return nil, fmt.Errorf("%w: LOGGER_LEVEL", _errUndefinedEnvVar)
	}

	nearThresholdRadius, defined := os.LookupEnv("NEAR_THRESHOLD_RADIUS_IN_MILES")
	if !defined {
		return nil, fmt.Errorf("%w: NEAR_THRESHOLD_RADIUS_IN_MILES", _errUndefinedEnvVar)
	}

	nearThresholdRadiusInMiles, err := strconv.Atoi(nearThresholdRadius)
	if err != nil {
		panic(fmt.Errorf("invalid value for NEAR_THRESHOLD_RADIUS_IN_MILES"))
	}

	return &Config{
		Database:            db,
		ServerPort:          serverPort,
		LoggerLevel:         loggerLevel,
		NearThresholdRadius: nearThresholdRadiusInMiles,
	}, nil
}
