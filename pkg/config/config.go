package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DbSource                string        `env:"DB_SOURCE"`
	MigrationUrl            string        `env:"MIGRATION_URL"`
	HttpServerAddress       string        `env:"HTTP_SERVER_ADDRESS"`
	GrpcServerAddress       string        `env:"GRPC_SERVER_ADDRESS"`
	RedisAddress            string        `env:"REDIS_ADDRESS"`
	RedisAccessTokenDB      int           `env:"REDIS_ACCESS_TOKEN_DB"`
	RedisRefreshTokenDB     int           `env:"REDIS_REFRESH_TOKEN_DB"`
	JWTIssuer               string        `env:"JWT_ISSUER"`
	JwtSecret               string        `env:"JWT_SECRET"`
	AccessTokenExpires      time.Duration `env:"ACCESS_TOKEN_EXPIRES"`
	RefreshTokenExpires     time.Duration `env:"REFRESH_TOKEN_EXPIRES"`
	OauthGoogleState        string        `env:"OAUTH_GOOGLE_STATE"`
	OauthGoogleClientID     string        `env:"OAUTH_GOOGLE_CLIENT_ID"`
	OauthGoogleClientSecret string        `env:"OAUTH_GOOGLE_CLIENT_SECRET"`
	OauthGoogleRedirectURL  string        `env:"OAUTH_GOOGLE_REDIRECT_URL"`
}

func Load() (*Config, error) {
	conf := &Config{}
	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
