package app

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Conf struct {
	// Addr is the port the server will listen on
	Addr string `env:"ADDR" envDefault:":8080"`

	// Verbose enables debug logging. Only set this to true in development
	Verbose bool `env:"VERBOSE"`

	// LogFormat can be either "production" or "development"
	LogFormat string `env:"LOG_FORMAT" envDefault:"production"`
}

func newConfigFromEnv() (*Conf, error) {
	cf := &Conf{}
	if err := env.ParseWithOptions(cf, env.Options{Prefix: "APP_"}); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return cf, nil
}
