package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DataBaseConfig
	JWT      JWTSConfig
}

type ServerConfig struct {
	Port string
}

type DataBaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTSConfig struct {
	Secret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},

		Database: DataBaseConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_DB"),
			SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
		},

		JWT: JWTSConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}
	return cfg, nil
}
