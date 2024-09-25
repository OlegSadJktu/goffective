package config

import (
	"errors"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres struct {
		Password string `env:"PG_PASSWORD" env-default:"goffective"`
		User     string `env:"PG_USER" env-default:"goffective"`
		Database string `env:"PG_DATABASE" env-default:"goffective"`
		Port     string `env:"PG_PORT" env-default:"5432"`
		Host     string `env:"PG_HOST" env-default:"localhost"`
	}
	Server struct {
		Port string `env:"PORT" env-default:"8080"`
	}
	Env string `env:"ENV" env-default:"dev"`
}

func MustLoad(path string) *Config {
	if len(path) == 0 {
		log.Fatal("invalid path")
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("invalid path: %s", err)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
