package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB *DBConfig
}

type DBConfig struct {
	Dialect  string
	Username string
	Password string
	Name     string
	Charset  string
}

func GetConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &Config{
		DB: &DBConfig{
			Dialect:  os.Getenv("DB_DIALECT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Charset:  os.Getenv("DB_CHARSET"),
		},
	}
}
