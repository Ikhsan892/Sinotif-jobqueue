package configs

import (
	"os"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func LoadDBConfig() DBConfig {
	return DBConfig{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}
}
