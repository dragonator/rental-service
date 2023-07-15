package config

import (
	"fmt"
	"os"
)

type Server struct {
	Port string
}

// NewDatabase is a constructor function for db config.
func NewServer() (*Server, error) {
	port, defined := os.LookupEnv("SERVER_PORT")
	if !defined {
		return nil, fmt.Errorf("%w: SERVER_PORT", _errUndefinedEnvVar)
	}

	return &Server{
		Port: port,
	}, nil
}
