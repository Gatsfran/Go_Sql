package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5400"`
	Username string `envconfig:"DB_USERNAME" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"docker"`
	Database string `envconfig:"DB_NAME" default:"postgres"`
}

func NewCfg() Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при парсинге переменных окружения: %v", err)
	}
	return cfg
}
