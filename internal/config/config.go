package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	PostgresConfig
	ServiceConfig
	ApiConfig
}

type PostgresConfig struct {
	PostgresUser string `env:"POSTGRES_USER,required"`
	PostgresPass string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB   string `env:"POSTGRES_DB,required"`
	PostgresHost string `env:"POSTGRES_HOST,required"`
	PostgresPort string `env:"POSTGRES_PORT,required"`
}

type ServiceConfig struct {
	ServiceHost string `env:"SERVICE_HOST,required"`
	ServicePort string `env:"SERVICE_PORT,required"`
}

type ApiConfig struct {
	ApiHost string `env:"API_HOST,required"`
	ApiPort string `env:"API_PORT,required"`
}

func MustLoad() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg
}
