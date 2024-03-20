package config

import (
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	Name string `json:"name" yaml:"name"`
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

type DB struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Name     string `json:"name" yaml:"name"`
}

type Gateway struct {
	Server Server
}

type User struct {
	Server Server
	DB     DB
}

type Auth struct {
	Server Server
}

type Article struct {
	Server Server
	DB     DB
}

type Bookmark struct {
	Server Server
	DB     DB
}

type Comment struct {
	Server Server
	DB     DB
}

type JWT struct {
	Secret  string        `json:"secret" yaml:"secret"`
	Expires time.Duration `json:"expires" yaml:"expires"`
}

type Config struct {
	Gateway  Gateway
	User     User
	Auth     Auth
	Article  Article
	Bookmark Bookmark
	Comment  Comment
	JWT      JWT
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.JWT.Expires = conf.JWT.Expires * time.Second

	return conf, nil
}
