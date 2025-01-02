package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Config struct {
	DB Postgres

	Server struct {
		Port      int    `envconfig:"SERVER_PORT" default:"8080"`
		JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
	}
	ExchangeService struct {
		Address string `envconfig:"EXCHANGE_SERVICE_ADDRESS" required:"true"`
	}
	Redis struct {
		Address  string `envconfig:"REDIS_ADDRESS" required:"true"`
		Password string `envconfig:"REDIS_PASSWORD" default:""`
		DB       int    `envconfig:"REDIS_DB" default:"0"`
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
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("No config.env file found, using system environment variables")
	}

	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
