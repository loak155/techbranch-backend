package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Server struct {
	Name string `env:"NAME"`
	Host string `env:"HOST"`
	Port int    `env:"PORT"`
}

type DB struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	Name     string `env:"NAME"`
}

type Gateway struct {
	Server Server `envPrefix:"SERVER_"`
}

type User struct {
	Server Server `envPrefix:"SERVER_"`
	DB     DB     `envPrefix:"DB_"`
}

type Auth struct {
	Server Server `envPrefix:"SERVER_"`
}

type Article struct {
	Server Server `envPrefix:"SERVER_"`
	DB     DB     `envPrefix:"DB_"`
}

type Bookmark struct {
	Server Server `envPrefix:"SERVER_"`
	DB     DB     `envPrefix:"DB_"`
}

type Comment struct {
	Server Server `envPrefix:"SERVER_"`
	DB     DB     `envPrefix:"DB_"`
}

type JWT struct {
	Secret  string        `env:"SECRET"`
	Expires time.Duration `env:"EXPIRES"`
}

type Google struct {
	State        string `env:"STATE"`
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURL  string `env:"REDIRECT_URL"`
}

type Oauth struct {
	Google Google `envPrefix:"GOOGLE_"`
}

type Config struct {
	Gateway  Gateway  `envPrefix:"GATEWAY_"`
	User     User     `envPrefix:"USER_"`
	Auth     Auth     `envPrefix:"AUTH_"`
	Article  Article  `envPrefix:"ARTICLE_"`
	Bookmark Bookmark `envPrefix:"BOOKMARK_"`
	Comment  Comment  `envPrefix:"COMMENT_"`
	JWT      JWT      `envPrefix:"JWT_"`
	Oauth    Oauth    `envPrefix:"OAUTH_"`
}

func Load() (*Config, error) {
	conf := &Config{}
	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
