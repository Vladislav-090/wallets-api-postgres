package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DataBaseConfig
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

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load .env file %w", err)
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},

		Database: DataBaseConfig{
			Host:     "localhost",
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Name:     os.Getenv("POSTGRES_DB"),
			SSLMode:  "disable",
		},
	}
	return cfg, nil
}
