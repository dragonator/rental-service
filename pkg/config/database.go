package config

import (
	"fmt"
	"os"
)

// Database is a struct containing db configuration.
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// NewDatabase is a constructor function for db config.
func NewDatabase() (*Database, error) {
	host, defined := os.LookupEnv("DATABASE_HOST")
	if !defined {
		return nil, fmt.Errorf("%w: DATABASE_HOST", errUndefinedEnvVar)
	}

	port, defined := os.LookupEnv("DATABASE_PORT")
	if !defined {
		return nil, fmt.Errorf("%w: DATABASE_PORT", errUndefinedEnvVar)
	}

	user, defined := os.LookupEnv("DATABASE_USER")
	if !defined {
		return nil, fmt.Errorf("%w: DATABASE_USER", errUndefinedEnvVar)
	}

	password, defined := os.LookupEnv("DATABASE_PASSWORD")
	if !defined {
		return nil, fmt.Errorf("%w: DATABASE_PASSWORD", errUndefinedEnvVar)
	}

	name, defined := os.LookupEnv("DATABASE_NAME")
	if !defined {
		return nil, fmt.Errorf("%w: DATABASE_NAME", errUndefinedEnvVar)
	}

	return &Database{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
	}, nil
}
