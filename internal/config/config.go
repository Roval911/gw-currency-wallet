package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

const (
	CONFIG_DIR  = "configs"
	CONFIG_FILE = "main"
)

type Config struct {
	DB Postgres

	Server struct {
		Port int `envconfig:"SERVER_PORT" default:"8080"`
	}
}

type Postgres struct {
	Host     string `envconfig:"DB_HOST" required:"true"`
	Port     int    `envconfig:"DB_PORT" required:"true"`
	Username string `envconfig:"DB_USERNAME" required:"true"`
	Name     string `envconfig:"DB_NAME" required:"true"`
	SSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
}

func New() (*Config, error) {
	// Загружаем переменные из config.env (если файл существует)
	if err := godotenv.Load(); err != nil {
		log.Println("No config.env file found, using system environment variables")
	}

	cfg := new(Config)

	// Заполняем структуру из переменных окружения
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
