package config

import "github.com/caarlos0/env/v11"

type Config struct {
	postgresUser string `env:"POSTGRES_USER,required"`
	postgresPass string `env:"POSTGRES_PASSWORD,required"`
	postgresDB   string `env:"POSTGRES_DB,required"`
	postgresHost string `env:"POSTGRES_HOST,required"`
	postgresPort string `env:"POSTGRES_PORT,required"`
	serviceHost  string `env:"SERVICE_HOST,required"`
	servicePort  string `env:"SERVICE_PORT,required"`
}

func MustLoad() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic(err)
	}

	return cfg
}
