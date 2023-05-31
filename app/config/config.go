package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host string `env:"HOST" envDefault:"127.0.0.1" json:"host"`
	Port string `env:"PORT" envDefault:"27312" json:"port"`
}

func (c *Config) Init() {
	err := env.Parse(c)
	if err != nil {
		panic(err)
	}
}

func NewConfig() *Config {
	return &Config{}
}
