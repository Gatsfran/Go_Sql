package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type PostgresConfig struct {
	Host     string `envconfig:"DB_HOST" default:"localhost"`
	Port     string `envconfig:"DB_PORT" default:"5400"`
	Username string `envconfig:"DB_USERNAME" default:"postgres"`
	Password string `envconfig:"DB_PASSWORD" default:"docker"`
	Database string `envconfig:"DB_NAME" default:"postgres"`
}
type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT" default:"8080"`
}
type Config struct {
	Postgres PostgresConfig
	Server   ServerConfig
}

func NewCfg() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения по умолчанию")
	}
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при парсинге переменных окружения: %v", err)
	}
	return cfg
}
