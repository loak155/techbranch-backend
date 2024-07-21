package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DbSource          string `env:"DB_SOURCE"`
	MigrationUrl      string `env:"MIGRATION_URL"`
	HttpServerAddress string `env:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress string `env:"GRPC_SERVER_ADDRESS"`
}

func Load() (*Config, error) {
	conf := &Config{}
	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
